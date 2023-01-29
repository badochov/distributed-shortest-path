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
	visited    map[db.VertexId]bool
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
		w.executions[requestId] = executionData{&Heap{}, make(map[db.VertexId]bool, len(w.data.vertices)), nil, nil}
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
	w.executions[requestId] = executionData{&Heap{}, make(map[db.VertexId]bool, len(w.data.vertices)), leftChild, rightChild}
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

func minDistanceVertex(isSet1 bool, distance1 float64, vertexId1 db.VertexId,
	isSet2 bool, distance2 float64, vertexId2 db.VertexId) (bool, float64, db.VertexId) {
	if isSet1 == false {
		return isSet2, distance2, vertexId2
	} else if isSet2 == false {
		return isSet1, distance1, vertexId1
	} else if distance1 < distance2 {
		return true, distance1, vertexId1
	} else {
		return true, distance2, vertexId2
	}
}

func (w *worker) Min(ctx context.Context, requestId api.RequestId) (bool, float64, db.VertexId, error) {
	isSetLeft, isSetRight, isSet := false, false, false
	var distanceLeft, distanceRight, distance float64
	var vertexIdLeft, vertexIdRight, vertexId db.VertexId
	var err error
	if w.executions[requestId].heap.Size() > 0 {
		isSet = true
		distance, vertexId = w.executions[requestId].heap.Top()
	}
	if w.executions[requestId].leftChild != nil {
		isSetLeft, distanceLeft, vertexIdLeft, err = w.executions[requestId].leftChild.Min(ctx, requestId)
		if err != nil {
			return false, 0, 0, err
		}
		isSet, distance, vertexId = minDistanceVertex(isSet, distance, vertexId, isSetLeft, distanceLeft, vertexIdLeft)
	}
	if w.executions[requestId].rightChild != nil {
		isSetRight, distanceRight, vertexIdRight, err = w.executions[requestId].rightChild.Min(ctx, requestId)
		if err != nil {
			return false, 0, 0, err
		}
		isSet, distance, vertexId = minDistanceVertex(isSet, distance, vertexId, isSetRight, distanceRight, vertexIdRight)
	}
	return isSet, distance, vertexId, nil
}

func (w *worker) Step(ctx context.Context, vertexId db.VertexId, destId db.VertexId, requestId api.RequestId) (bool, float64, error) {
	return false, 0, nil
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
