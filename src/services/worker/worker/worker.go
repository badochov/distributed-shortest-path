package worker

import (
	"context"
	"log"
	"math/rand"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/worker/api"
	"github.com/badochov/distributed-shortest-path/src/services/worker/discoverer"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/link_server"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service"
)

type Deps struct {
	Db         db.DB
	Discoverer discoverer.Discoverer
	RegionID   db.RegionId
	Context    context.Context
	LinkPort   int
}

// Worker All methods from link service and worker service should end up calling this interface.
type Worker interface {
	service.Worker
	link_server.Worker

	LoadRegionData(ctx context.Context) error
}

type workerData struct {
	vertices                    []db.VertexId
	edges                       map[db.VertexId][]db.Edge
	arcFlags                    []db.ArcFlag
	neighbouringVerticesRegions map[db.EdgeId]db.RegionId
}

type executionData struct {
	heap       *Heap
	leftChild  link.Link // leftChild == nil iff left child does not exist
	rightChild link.Link // rightChild == nil iff right child does not exist
}

type worker struct {
	db         db.DB
	discoverer discoverer.Discoverer
	generation db.Generation
	regionId   db.RegionId
	data       workerData
	links      map[db.RegionId]link.RegionManager
	linkPort   int
	executions map[api.RequestId]executionData
}

func (w *worker) CalculateArcFlags(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func randomId(min db.RegionId, max db.RegionId) uint16 {
	return uint16(rand.Intn(int(max-min))) + min
}

func (w *worker) Init(ctx context.Context, minRegionId db.RegionId, maxRegionId db.RegionId, requestId api.RequestId) error {
	if minRegionId == maxRegionId {
		w.executions[requestId] = executionData{&Heap{}, nil, nil}
		return nil
	}
	leftChildId := randomId(minRegionId, (minRegionId+maxRegionId)/2)
	rightChildId := randomId((minRegionId+maxRegionId)/2+1, maxRegionId)

	var leftChild, rightChild link.Link
	var err error
	_, leftChild, err = w.links[leftChildId].GetLink()
	if err != nil {
		return err
	}
	_, rightChild, err = w.links[rightChildId].GetLink()
	if err != nil {
		return err
	}
	w.executions[requestId] = executionData{&Heap{}, leftChild, rightChild}
	err = leftChild.Init(ctx, minRegionId, (minRegionId+maxRegionId)/2, requestId)
	if err != nil {
		return err
	}
	err = rightChild.Init(ctx, (minRegionId+maxRegionId)/2+1, maxRegionId, requestId)
	if err != nil {
		return err
	}
	return nil
}

func (w *worker) ShortestPath(ctx context.Context, args service.ShortestPathArgs) (service.ShortestPathResult, error) {
	var err error
	err = w.Init(ctx, 0, db.RegionId(len(w.links)-1), args.RequestId)
	return service.ShortestPathResult{}, err
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
	for _, edges := range w.data.edges {
		for _, e := range edges {
			eIds = append(eIds, e.Id)
		}
	}
	w.data.arcFlags, err = w.db.GetArcFlags(ctx, eIds, w.generation)
	if err != nil {
		return
	}

	w.data.neighbouringVerticesRegions, err = w.db.GetEdgeToRegionMapping(ctx, w.regionId, w.generation)
	if err != nil {
		return
	}

	return
}

func (w *worker) initDiscoverer(ctx context.Context) error {
	if err := w.discoverer.Run(ctx); err != nil {
		return err
	}
	go func() {
		for {
			select {
			// some worker changed its status, e.g. failed (or went up)
			case status := <-w.discoverer.InstanceStatuses():
				w.handleInstanceStatus(status)
			// set of workers in region has changed
			case data := <-w.discoverer.RegionDataChan():
				w.handleRegionData(ctx, data)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (w *worker) Add(ctx context.Context, a int32, b int32) (int32, error) {
	return a + b, nil
}

func (w *worker) handleInstanceStatus(status discoverer.WorkerInstanceStatus) {
	// TODO [wprzytula]
}

func (w *worker) handleRegionData(ctx context.Context, data discoverer.RegionData) {
	log.Println("Handling region data change", data)
	l, ok := w.links[data.RegionId]
	if !ok {
		l = link.NewRegionDialer(w.linkPort)
	}
	err := l.UpdateInstances(ctx, data.Instances)
	if err != nil {
		log.Print(err)
	}
	if !ok {
		w.links[data.RegionId] = l
	}
}

func New(deps Deps) (Worker, error) {
	w := &worker{
		db:         deps.Db,
		discoverer: deps.Discoverer,
		regionId:   deps.RegionID,
		linkPort:   deps.LinkPort,
		links:      make(map[db.RegionId]link.RegionManager),
	}
	if err := w.initDiscoverer(deps.Context); err != nil {
		return nil, err
	}

	return w, nil
}
