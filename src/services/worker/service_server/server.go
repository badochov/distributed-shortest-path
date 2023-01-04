package service_server

import (
	"fmt"
	api "github.com/badochov/distributed-shortest-path/src/libs/worker_api"
	"github.com/badochov/distributed-shortest-path/src/services/worker/common"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service_server/executor"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Deps struct {
	Executor executor.Executor
	Port     int
}

type Server interface {
	common.Runner
}

type server struct {
	engine  *gin.Engine
	handler handler
	port    int
}

func (s *server) Run() error {
	if err := s.handler.Run(); err != nil {
		return err
	}

	return s.engine.Run(fmt.Sprintf(":%d", s.port))
}

type handler struct {
	executor executor.Executor
}

func (h *handler) Run() error {
	return h.executor.Run()
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

func (h *handler) CalculateArcFlags(c *gin.Context) {
	resp, code, err := h.executor.CalculateArcFlags()
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

func New(deps Deps) Server {
	router := gin.Default()

	h := handler{
		executor: deps.Executor,
	}

	router.POST(api.ShortestPathUrl, h.ShortestPath)

	router.GET(api.CalculateArcFlagsUrl, h.CalculateArcFlags)

	router.GET(api.HealthzUrl, h.Healthz)

	return &server{
		engine:  router,
		handler: h,
		port:    deps.Port,
	}
}
