package executor

import (
	"context"
	"fmt"
	api "github.com/badochov/distributed-shortest-path/src/libs/api/manager_api"
	"github.com/badochov/distributed-shortest-path/src/libs/api/worker_api"
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/manager/common"
	"github.com/badochov/distributed-shortest-path/src/services/manager/worker"
	"github.com/badochov/distributed-shortest-path/src/services/manager/worker/service_manager"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
	"time"
)

type regionId = db.RegionId
type generation = db.Generation

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
	generation          generation
	nextGeneration      generation
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

func (e *executor) getRegion(id db.VertexId) (db.RegionId, error) {
	ctx, can := timeoutCtx(1 * time.Second)
	defer can()

	return e.db.GetVertexRegion(ctx, id, e.generation)
}

func (e *executor) ShortestPath(req api.ShortestPathRequest) (resp api.ShortestPathResponse, code int, err error) {
	if err := e.start(); err != nil {
		return api.ShortestPathResponse{}, http.StatusInternalServerError, err
	}
	defer e.finish()

	ctx, can := timeoutCtx(30 * time.Second)
	defer can()

	regId, err := e.getRegion(req.From)
	if err != nil {
		return api.ShortestPathResponse{}, http.StatusInternalServerError, err
	}

	workerReq := worker_api.ShortestPathRequest{
		From: req.From,
		To:   req.To,
	}

	err = nil

	const retries = 3
	for i := 0; i < retries; i++ {
		res, workerErr := e.clients[regId].ShortestPath(ctx, workerReq)
		if workerErr == nil {
			return api.ShortestPathResponse{
				Distance: res.Distance,
				Vertices: res.Vertices,
			}, http.StatusOK, nil
		}
		err = multierror.Append(err, workerErr)
	}

	return api.ShortestPathResponse{}, http.StatusInternalServerError, err
}

func (e *executor) AddEdges(req api.AddEdgesRequest) (resp api.AddEdgesRequest, code int, err error) {
	if err := e.start(); err != nil {
		return api.AddEdgesRequest{}, http.StatusInternalServerError, err
	}
	defer e.finish()

	ctx, can := timeoutCtx(15 * time.Second)
	defer can()

	if err := e.db.AddEdges(ctx, req.Edges, e.generation); err != nil {
		return api.AddEdgesRequest{}, http.StatusInternalServerError, err
	}
	return api.AddEdgesRequest{}, http.StatusOK, err
}

func (e *executor) AddVertices(req api.AddVerticesRequest) (resp api.AddVerticesResponse, code int, err error) {
	if err := e.start(); err != nil {
		return api.AddVerticesResponse{}, http.StatusInternalServerError, err
	}
	defer e.finish()

	ctx, can := timeoutCtx(15 * time.Second)
	defer can()

	if err := e.db.AddVertices(ctx, req.Vertices, e.generation); err != nil {
		return api.AddVerticesResponse{}, http.StatusInternalServerError, err
	}
	return api.AddVerticesResponse{}, http.StatusOK, err
}

func (e *executor) RecalculateDS() (resp api.RecalculateDsResponse, code int, err error) {
	e.recalculateLock.Lock()
	defer e.recalculateLock.Unlock()

	ctx, can := context.WithTimeout(context.Background(), 10*time.Minute)
	defer can()

	wrap := func(err error) (api.RecalculateDsResponse, int, error) {
		return api.RecalculateDsResponse{}, http.StatusInternalServerError, err
	}

	if err := e.incNextGen(ctx); err != nil {
		return wrap(err)
	}
	if err := e.divideIntoRegions(ctx); err != nil {
		return wrap(err)
	}
	if err := e.calculateArcFlags(ctx); err != nil {
		return wrap(err)
	}

	return api.RecalculateDsResponse{}, http.StatusOK, nil
}

func (e *executor) incNextGen(ctx context.Context) (err error) {
	ctx, can := context.WithTimeout(ctx, time.Second)
	defer can()

	// TODO Think if retires should be implemented and how.
	if err := e.db.SetNextGeneration(ctx, e.nextGeneration+1); err != nil {
		return err
	}

	e.nextGeneration++
	return nil
}

func (e *executor) divideIntoRegions(ctx context.Context) error {
	ctx, can := context.WithTimeout(ctx, time.Minute)
	defer can()

	// TODO Think if retires should be implemented and how.
	return e.divideIntoRegionsDoer(ctx, e.nextGeneration)
}

func (e *executor) divideIntoRegionsDoer(ctx context.Context, nextGeneration generation) error {
	// TODO Implement me
	panic("Implement me")
}

func (e *executor) start() error {
	if e.recalculateLock.TryRLock() {
		return nil
	}
	return fmt.Errorf("data structure recalculation is in progress")
}

func (e *executor) finish() {
	e.recalculateLock.RUnlock()
}

func (e *executor) Healthz() (resp api.HealthzResponse, code int, err error) {
	return api.HealthzResponse{}, http.StatusOK, nil
}

func (e *executor) calculateArcFlags(baseCtx context.Context) error {
	ctx, can := context.WithTimeout(baseCtx, 8*time.Minute)
	defer can()

	grp, grpCtx := errgroup.WithContext(ctx)

	for regId, cl := range e.clients {
		regId := regId
		cl := cl

		grp.Go(func() error {
			var err error

			const retries = 3
			for i := 0; i < retries; i++ {
				if calcErr := cl.CalculateArcFlags(grpCtx); calcErr != nil {
					err = multierror.Append(err, calcErr)
				} else {
					return nil
				}
			}
			return fmt.Errorf("error calculating flags in region %d, %w", regId, err)
		})
	}

	return grp.Wait()
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
