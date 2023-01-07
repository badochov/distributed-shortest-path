package executor

import (
	"context"
	"fmt"
	api "github.com/badochov/distributed-shortest-path/src/libs/api/manager_api"
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/manager/common"
	"github.com/badochov/distributed-shortest-path/src/services/manager/worker"
	"github.com/badochov/distributed-shortest-path/src/services/manager/worker/service_manager"
	"net/http"
	"sync"
	"time"
)

type regionId int

type Deps struct {
	NumRegions          int
	RegionUrlTemplate   string
	Db                  db.DB
	Port                int
	WorkerServerManager service_manager.WorkerServiceManager
}

type Executor interface {
	ShortestPath(req api.ShortestPathRequest) (resp api.ShortestPathResponse, code int, err error)
	AddEdges(req api.AddEdgesRequest) (resp api.AddEdgesRequest, code int, err error)
	AddVertices(req api.AddVerticesRequest) (resp api.AddVerticesResponse, code int, err error)
	RecalculateDS() (resp api.RecalculateDsResponse, code int, err error)

	GetGeneration() (resp api.GetGenerationResponse, code int, err error)

	Healthz() (resp api.HealthzResponse, code int, err error)

	common.Runner
}

type executor struct {
	generation          uint16
	clients             map[regionId]worker.Client
	db                  db.DB
	workerServerManager service_manager.WorkerServiceManager
	recalculateLock     sync.RWMutex
}

func (e *executor) GetGeneration() (resp api.GetGenerationResponse, code int, err error) {
	return api.GetGenerationResponse{Generation: e.generation}, http.StatusOK, nil
}

func (e *executor) Run() error {
	return nil
}

func (e *executor) ShortestPath(req api.ShortestPathRequest) (resp api.ShortestPathResponse, code int, err error) {
	//TODO implement me
	panic("implement me")
}

func (e *executor) AddEdges(req api.AddEdgesRequest) (resp api.AddEdgesRequest, code int, err error) {
	e.recalculateLock.RLock()
	defer e.recalculateLock.RUnlock()

	ctx, can := timeoutCtx(15 * time.Second)
	defer can()

	if err := e.db.AddEdges(ctx, req.Edges, e.generation); err != nil {
		return api.AddEdgesRequest{}, http.StatusInternalServerError, err
	}
	return api.AddEdgesRequest{}, http.StatusOK, err
}

func (e *executor) AddVertices(req api.AddVerticesRequest) (resp api.AddVerticesResponse, code int, err error) {
	e.recalculateLock.RLock()
	defer e.recalculateLock.RUnlock()
	
	ctx, can := timeoutCtx(15 * time.Second)
	defer can()

	if err := e.db.AddVertices(ctx, req.Vertices, e.generation); err != nil {
		return api.AddVerticesResponse{}, http.StatusInternalServerError, err
	}
	return api.AddVerticesResponse{}, http.StatusOK, err
}

func (e *executor) RecalculateDS() (resp api.RecalculateDsResponse, code int, err error) {
	//TODO implement me
	panic("implement me")
}

func (e *executor) Healthz() (resp api.HealthzResponse, code int, err error) {
	return api.HealthzResponse{}, http.StatusOK, nil
}

func timeoutCtx(duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), duration)
}

func New(deps Deps) Executor {
	ex := &executor{
		db:                  deps.Db,
		clients:             make(map[regionId]worker.Client, deps.NumRegions),
		workerServerManager: deps.WorkerServerManager,
	}

	for i := 0; i < deps.NumRegions; i++ {
		ex.clients[regionId(i)] = worker.NewClient(worker.Deps{
			HttpClient: http.DefaultClient, // TODO customize timeouts,
			Url:        fmt.Sprintf(deps.RegionUrlTemplate+":%d", i, deps.Port),
		})
	}

	return ex
}
