package executor

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/manager/api"
	"github.com/badochov/distributed-shortest-path/src/services/manager/worker"
	"github.com/badochov/distributed-shortest-path/src/services/manager/worker/service_manager"
	workerApi "github.com/badochov/distributed-shortest-path/src/services/worker/api"
	"github.com/cenkalti/backoff/v4"
	"golang.org/x/sync/errgroup"
)

// TODO [wprzytula] retries (one implemention of retires is in src/services/manager/worker/service_manager/manager.go)
// TODO [wprzytula] restarting worker service on failures
// TODO [wprzytula] bringing back sane state at startup in case of service failure
// TODO [wprzytula] redo timeouts. Current implementation with timeoutCtx at the start of almost every method looks poor.

type regionId = db.RegionId
type generation = db.Generation

type Deps struct {
	NumRegions            int
	RegionUrlTemplate     string
	Db                    db.DB
	Port                  int
	WorkerServerManager   service_manager.WorkerServiceManager
	DefaultWorkerReplicas int32
}

type Executor interface {
	ShortestPath(req api.ShortestPathRequest) (resp api.ShortestPathResponse, code int, err error)
	AddEdges(req api.AddEdgesRequest) (resp api.AddEdgesRequest, code int, err error)
	AddVertices(req api.AddVerticesRequest) (resp api.AddVerticesResponse, code int, err error)
	RecalculateDS() (resp api.RecalculateDsResponse, code int, err error)

	GetGeneration() (resp api.GetGenerationResponse, code int, err error)

	Healthz() (resp api.HealthzResponse, code int, err error)
}

type executor struct {
	generation            generation
	nextGeneration        generation
	clients               map[regionId]worker.Client
	db                    db.DB
	workerServerManager   service_manager.WorkerServiceManager
	recalculateLock       sync.RWMutex
	defaultWorkerReplicas int32
	requestId             atomic.Uint64
}

func (e *executor) GetGeneration() (resp api.GetGenerationResponse, code int, err error) {
	return api.GetGenerationResponse{Generation: e.generation}, http.StatusOK, nil
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

	ctx, can := timeoutCtx(30 * time.Second) // TODO[wprzytula] why? how much?
	defer can()

	regId, err := e.getRegion(req.From)
	if err != nil {
		return api.ShortestPathResponse{}, http.StatusInternalServerError, err
	}
	reqId := e.genReqId()

	workerReq := workerApi.ShortestPathRequest{
		RequestId: reqId,
		From:      req.From,
		To:        req.To,
	}

	// TODO [wprzytula] Think if retries should be implemented and how.
	res, err := e.clients[regId].ShortestPath(ctx, workerReq)
	if err != nil {
		return api.ShortestPathResponse{}, http.StatusInternalServerError, err
	}

	return api.ShortestPathResponse{
		Distance: res.Distance,
		Vertices: res.Vertices,
	}, http.StatusOK, nil
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

	ctx, cancel := timeoutCtx(15 * time.Second)
	defer cancel()

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
		// TODO [wprzytula] handle restart workers on failure.
		// TODO [wprzytula] handle clean state. (Set proper active generation, etc.).
		return api.RecalculateDsResponse{}, http.StatusInternalServerError, err
	}

	// Shutdown worker service.
	if err := e.rescaleAllRegions(ctx, 0); err != nil {
		return wrap(err)
	}
	if err := e.incNextGen(ctx); err != nil {
		return wrap(err) // TODO[wprzytula]: consider inc iff not yet incremented
	}
	if err := e.divideIntoRegions(ctx); err != nil {
		return wrap(err)
	}
	if err := e.setActiveGeneration(ctx, e.nextGeneration); err != nil {
		return wrap(err)
	}
	// Start worker service.
	if err := e.rescaleAllRegions(ctx, e.defaultWorkerReplicas); err != nil {
		return wrap(err)
	}
	// TODO [wprzytula] wait for workers to be alive (eg. Add Healthz method to client and wait for it to respond with success)
	if err := e.calculateArcFlags(ctx); err != nil {
		return wrap(err)
	}
	if err := e.setGenToNext(ctx); err != nil {
		return wrap(err)
	}
	// Restart worker service.
	if err := e.rescaleAllRegions(ctx, 0); err != nil {
		return wrap(err)
	}
	if err := e.rescaleAllRegions(ctx, e.defaultWorkerReplicas); err != nil {
		return wrap(err)
	}
	// TODO [wprzytula] wait for workers to be alive (eg. Add Healthz method to client and wait for it to respond with success)

	return api.RecalculateDsResponse{}, http.StatusOK, nil
}

func (e *executor) rescaleAllRegions(ctx context.Context, replicas int32) error {
	// TODO [wprzytula] fix state where one of the regions didn't rescale properly.
	errgrp, ctx := errgroup.WithContext(ctx)
	for id := range e.clients {
		errgrp.Go(func() error {
			expBackoff := backoff.NewExponentialBackOff() // TODO [wprzytula] customize timeouts.
			return backoff.Retry(func() error { return e.workerServerManager.Rescale(ctx, id, replicas) }, expBackoff)
		})
	}
	return errgrp.Wait()
}

func (e *executor) incNextGen(ctx context.Context) (err error) {
	ctx, can := context.WithTimeout(ctx, time.Second)
	defer can()

	// TODO [wprzytula] Think if retries should be implemented and how.
	if err := e.db.SetNextGeneration(ctx, e.nextGeneration+1); err != nil {
		return err
	}

	e.nextGeneration++
	return nil
}

func (e *executor) getNextGen(ctx context.Context) (generation, error) {
	ctx, can := context.WithTimeout(ctx, time.Second)
	defer can()

	// TODO [wprzytula] Think if retries should be implemented and how.
	return e.db.GetNextGeneration(ctx)
}

func (e *executor) getGen(ctx context.Context) (generation, error) {
	ctx, can := context.WithTimeout(ctx, time.Second)
	defer can()

	// TODO [wprzytula] Think if retries should be implemented and how.
	return e.db.GetCurrentGeneration(ctx)
}

func (e *executor) setGenToNext(ctx context.Context) (err error) {
	ctx, can := context.WithTimeout(ctx, time.Second)
	defer can()

	// TODO [wprzytula] Think if retries should be implemented and how.
	if err := e.db.SetCurrentGeneration(ctx, e.nextGeneration); err != nil {
		return err
	}

	e.generation = e.nextGeneration
	return nil
}

func (e *executor) setActiveGeneration(ctx context.Context, gen generation) error {
	ctx, can := context.WithTimeout(ctx, time.Second)
	defer can()

	// TODO [wprzytula] Think if retries should be implemented and how.
	if err := e.db.SetActiveGeneration(ctx, gen); err != nil {
		return err
	}
	return nil
}

func (e *executor) divideIntoRegions(ctx context.Context) error {
	ctx, can := context.WithTimeout(ctx, time.Minute)
	defer can()

	// TODO [wprzytula] Think if retries should be implemented and how.
	return e.divideIntoRegionsDoer(ctx)
}

func (e *executor) divideIntoRegionsHelper(ctx context.Context, coordsBetween db.CoordsBetween,
	count int64, minRegionId db.RegionId, maxRegionId db.RegionId, direction bool) error {
	var err error
	if minRegionId == maxRegionId {
		_, err := e.db.SetRegion(ctx, coordsBetween, minRegionId, e.generation)
		return err
	}
	if direction {
		left, right := coordsBetween.Longitude.Min, coordsBetween.Longitude.Max
		mid := (left + right) / 2
		var leftPart, rightPart, midPart int64
		for {
			midPart, err = e.db.GetVertexCount(ctx,
				db.CoordsBetween{
					Latitude:  coordsBetween.Latitude,
					Longitude: db.MinMax{Min: mid, Max: mid},
				},
				e.generation)
			if err != nil {
				return err
			}
			leftPart, err = e.db.GetVertexCount(ctx,
				db.CoordsBetween{
					Latitude:  coordsBetween.Latitude,
					Longitude: db.MinMax{Min: coordsBetween.Longitude.Min, Max: mid},
				},
				e.generation)
			if err != nil {
				return err
			}
			leftPart = leftPart - midPart
			rightPart = count - leftPart // including midPart
			if leftPart > rightPart {
				right = mid
			} else if leftPart+midPart > rightPart-midPart {
				break
			} else {
				left = mid
			}
			mid = (left + right) / 2
		}
		err = e.divideIntoRegionsHelper(ctx,
			db.CoordsBetween{
				Latitude:  coordsBetween.Latitude,
				Longitude: db.MinMax{Min: coordsBetween.Longitude.Min, Max: mid},
			},
			leftPart, minRegionId, (minRegionId+maxRegionId)/2, false)
		if err != nil {
			return err
		}
		return e.divideIntoRegionsHelper(ctx,
			db.CoordsBetween{
				Latitude:  coordsBetween.Latitude,
				Longitude: db.MinMax{Min: mid, Max: coordsBetween.Longitude.Max}},
			rightPart, (minRegionId+maxRegionId)/2+1, maxRegionId, false)
	} else {
		down, up := coordsBetween.Latitude.Min, coordsBetween.Latitude.Max
		mid := (down + up) / 2
		var downPart, upPart, midPart int64
		var err error
		for {
			midPart, err = e.db.GetVertexCount(ctx,
				db.CoordsBetween{
					Latitude:  db.MinMax{Min: mid, Max: mid},
					Longitude: coordsBetween.Longitude,
				},
				e.generation)
			if err != nil {
				return err
			}
			downPart, err = e.db.GetVertexCount(ctx,
				db.CoordsBetween{
					Latitude:  db.MinMax{Min: coordsBetween.Latitude.Min, Max: mid},
					Longitude: coordsBetween.Longitude,
				},
				e.generation)
			if err != nil {
				return err
			}
			downPart = downPart - midPart
			upPart = count - downPart // including midPart
			if downPart > upPart {
				up = mid
			} else if downPart+midPart > upPart-midPart {
				break
			} else {
				down = mid
			}
			mid = (down + up) / 2
		}
		err = e.divideIntoRegionsHelper(ctx,
			db.CoordsBetween{
				Latitude:  db.MinMax{Min: coordsBetween.Latitude.Min, Max: mid},
				Longitude: coordsBetween.Longitude,
			},
			downPart, minRegionId, (minRegionId+maxRegionId)/2, true)
		if err != nil {
			return err
		}
		return e.divideIntoRegionsHelper(ctx,
			db.CoordsBetween{
				Latitude:  db.MinMax{Min: mid, Max: coordsBetween.Latitude.Max},
				Longitude: coordsBetween.Longitude,
			},
			upPart, (minRegionId+maxRegionId)/2+1, maxRegionId, true)
	}
}

func (e *executor) divideIntoRegionsDoer(ctx context.Context) error {
	bounds, err := e.db.GetCoordsBounds(ctx)
	if err != nil {
		return err
	}
	count, err := e.db.GetVertexCount(ctx, bounds, e.generation)
	if err != nil {
		return err
	}
	return e.divideIntoRegionsHelper(ctx, bounds, count, 0, regionId(len(e.clients)-1), true)
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
	ctx, can := context.WithTimeout(baseCtx, 8*time.Minute) // TODO[wprzytula]
	defer can()

	grp, grpCtx := errgroup.WithContext(ctx)

	for regId, cl := range e.clients {
		regId := regId
		cl := cl

		grp.Go(func() error {
			// TODO [wprzytula] Think if retries should be implemented and how.
			if err := cl.CalculateArcFlags(grpCtx); err != nil {
				return fmt.Errorf("error calculating flags in region %d, %w", regId, err)
			}
			return nil
		})
	}

	return grp.Wait()
}

func (e *executor) init(ctx context.Context) (err error) {
	e.generation, err = e.getGen(ctx)
	if err != nil {
		return
	}
	e.nextGeneration, err = e.getNextGen(ctx)
	// TODO [wprzytula] handle starting worker service if it was stopped by previous manager who failed.
	// TODO [wprzytula] handle clean state. (Set proper active generation, etc.).
	return
}

func (e *executor) genReqId() workerApi.RequestId {
	return workerApi.RequestId(e.requestId.Add(1))
}

func timeoutCtx(duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), duration)
}

func New(ctx context.Context, deps Deps) (Executor, error) {
	ex := &executor{
		db:                    deps.Db,
		clients:               make(map[regionId]worker.Client, deps.NumRegions),
		workerServerManager:   deps.WorkerServerManager,
		defaultWorkerReplicas: deps.DefaultWorkerReplicas,
	}

	for i := 0; i < deps.NumRegions; i++ {
		ex.clients[regionId(i)] = worker.NewClient(worker.Deps{
			HttpClient: http.DefaultClient, // TODO [wprzytula] customize timeouts,
			Url:        fmt.Sprintf(deps.RegionUrlTemplate+":%d", i, deps.Port),
		})
	}

	if err := ex.init(ctx); err != nil {
		return nil, err
	}
	return ex, nil
}
