package server

import (
	"github.com/badochov/distributed-shortest-path/src/services/manager/executor"
	"github.com/badochov/distributed-shortest-path/src/services/manager/server/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Deps struct {
	Executor executor.Executor
}

type Server interface {
	Run() error
}

type server struct {
	engine *gin.Engine
}

func (s *server) Run() error {
	return s.engine.Run()
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

func New(deps Deps) Server {
	router := gin.Default()

	h := handler{
		executor: deps.Executor,
	}

	router.POST("/shortest_path", h.ShortestPath)

	router.GET("/recalculate_ds", h.RecalculateDs)

	router.POST("/add_vertices", h.AddVertices)
	router.POST("/add_edges", h.AddEdges)

	return &server{engine: router}
}
