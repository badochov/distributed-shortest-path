package server

import (
	"fmt"
	"net/http"

	"github.com/badochov/distributed-shortest-path/src/services/manager/api"
	"github.com/badochov/distributed-shortest-path/src/services/manager/executor"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	Executor executor.Executor
	Port     int
}

type Server interface {
	Run() error
}

type server struct {
	engine  *gin.Engine
	handler handler
	port    int
}

func (s *server) Run() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.port))
}

type handler struct {
	executor executor.Executor
}

func (h *handler) ShortestPath(c *gin.Context) {
	var req api.ShortestPathRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, code, err := h.executor.ShortestPath(req)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func (h *handler) RecalculateDs(c *gin.Context) {
	resp, code, err := h.executor.RecalculateDS()
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func (h *handler) AddVertices(c *gin.Context) {
	var req api.AddVerticesRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, code, err := h.executor.AddVertices(req)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func (h *handler) AddEdges(c *gin.Context) {
	var req api.AddEdgesRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, code, err := h.executor.AddEdges(req)
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func (h *handler) Healthz(c *gin.Context) {
	resp, code, err := h.executor.Healthz()
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func (h *handler) GetGeneration(c *gin.Context) {
	resp, code, err := h.executor.GetGeneration()
	if err != nil {
		c.AbortWithError(code, err)
		return
	}
	c.JSON(code, resp)
}

func New(deps Deps) Server {
	router := gin.Default()

	h := handler{
		executor: deps.Executor,
	}

	router.POST(api.ShortestPathUrl, h.ShortestPath)

	router.GET(api.RecalculateDsURL, h.RecalculateDs)

	router.POST(api.AddVerticesUrl, h.AddVertices)
	router.POST(api.AddEdgesUrl, h.AddEdges)

	router.GET(api.GetGenerationUrl, h.GetGeneration)

	router.GET(api.HealthzUrl, h.Healthz)

	return &server{
		engine:  router,
		handler: h,
		port:    deps.Port,
	}
}
