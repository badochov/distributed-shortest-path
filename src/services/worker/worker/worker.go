package worker

import (
	"context"
	api "github.com/badochov/distributed-shortest-path/src/libs/api/worker_api"
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/worker/discoverer"
)

type Deps struct {
	Db         db.DB
	Discoverer discoverer.Discoverer
	RegionID   db.RegionId
	Context    context.Context
}

type ShortestPathArgs = api.ShortestPathRequest
type ShortestPathResult = api.ShortestPathResponse

type ServerInterface interface {
	CalculateArcFlags() error
	ShortestPath(args ShortestPathArgs) (ShortestPathResult, error)
}

type LinkInterface interface {
	Add(a int32, b int32) (int32, error) // Example
}

// Worker All methods from link service and worker service should end up calling this interface.
type Worker interface {
	ServerInterface
	LinkInterface

	LoadRegionData(ctx context.Context) error
}

type workerData struct {
	vertices []db.VertexId
	edges    []db.Edge
	arcFlags []db.ArcFlag
}

type worker struct {
	db         db.DB
	discoverer discoverer.Discoverer
	generation db.Generation
	regionId   db.RegionId
	data       workerData
}

func (w *worker) CalculateArcFlags() error {
	//TODO implement me
	panic("implement me")
}

func (w *worker) ShortestPath(args ShortestPathArgs) (ShortestPathResult, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worker) LoadRegionData(ctx context.Context) (err error) {
	w.generation, err = w.db.GetActiveGeneration(ctx)
	if err != nil {
		return
	}
	w.data.vertices, err = w.db.GetVertexIds(ctx, w.regionId, w.generation)
	if err != nil {
		return
	}
	w.data.edges, err = w.db.GetEdgesFrom(ctx, w.data.vertices, w.generation)
	if err != nil {
		return
	}
	eIds := make([]db.EdgeId, 0, len(w.data.edges))
	for _, e := range w.data.edges {
		eIds = append(eIds, e.Id)
	}
	w.data.arcFlags, err = w.db.GetArcFlags(ctx, eIds, w.generation)
	if err != nil {
		return
	}

	return
}

func (w *worker) initDiscoverer(ctx context.Context) error {
	err := w.discoverer.Run(ctx)
	if err != nil {
		return err
	}
	go func() {
		select {
		case status := <-w.discoverer.InstanceStatuses():
			w.handleInstanceStatus(status)
		case data := <-w.discoverer.RegionDataChan():
			w.handleRegionData(data)
		case <-ctx.Done():
			return
		}
	}()
	return nil
}

func (w *worker) Add(a int32, b int32) (int32, error) {
	return a + b, nil
}

func (w *worker) handleInstanceStatus(status discoverer.WorkerInstanceStatus) {
	//TODO
	panic("IMPLEMENT ME")
}

func (w *worker) handleRegionData(data discoverer.RegionData) {
	//TODO
	panic("IMPLEMENT ME")
}

func New(deps Deps) (Worker, error) {
	w := &worker{
		db:         deps.Db,
		discoverer: deps.Discoverer,
		regionId:   deps.RegionID,
	}
	if err := w.initDiscoverer(deps.Context); err != nil {
		return nil, err
	}

	return w, nil
}
