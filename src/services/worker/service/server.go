package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/badochov/distributed-shortest-path/src/services/worker/api"
	"github.com/badochov/distributed-shortest-path/src/services/worker/common"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	Worker Worker
	Port   int
}

type ShortestPathArgs = api.ShortestPathRequest
type ShortestPathResult = api.ShortestPathResponse

type Worker interface {
	CalculateArcFlags(ctx context.Context) error
	ShortestPath(ctx context.Context, args ShortestPathArgs) (ShortestPathResult, error)
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
	return s.engine.Run(fmt.Sprintf(":%d", s.port))
}

type handler struct {
	worker Worker
}

func (h *handler) ShortestPath(c *gin.Context) {
	// TODO [wprzytula] adjust timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var req api.ShortestPathRequest
	if err := c.BindJSON(&req); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resp, err := h.worker.ShortestPath(ctx, req)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *handler) CalculateArcFlags(c *gin.Context) {
	// TODO [wprzytula] adjust timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err := h.worker.CalculateArcFlags(ctx)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func (h *handler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

func New(deps Deps) Server {
	router := gin.Default()

	h := handler{
		worker: deps.Worker,
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
