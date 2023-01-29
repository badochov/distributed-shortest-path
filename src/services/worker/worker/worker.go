package worker

import (
	"container/heap"
	"context"
	"log"
	"math"
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
	edgeTargets                 map[db.VertexId]bool
	arcFlags                    []db.ArcFlag
	neighbouringVerticesRegions map[db.EdgeId]db.RegionId
}

type executionData struct {
	queue     PriorityQueue
	inQueue   map[db.VertexId]*Item
	processed map[db.VertexId]bool
	// leftChild  link.Link // leftChild == nil iff left child does not exist
	// rightChild link.Link // rightChild == nil iff right child does not exist
}

type worker struct {
	db             db.DB
	discoverer     discoverer.Discoverer
	generation     db.Generation
	regionId       db.RegionId
	data           workerData
	links          map[db.RegionId]link.RegionManager
	linkPort       int
	executions     map[api.RequestId]executionData
	executionLinks map[api.RequestId][]link.Link
}

func (w *worker) CalculateArcFlags(ctx context.Context) error {
	//TODO implement me
	return nil
}

func randomId(min db.RegionId, max db.RegionId) uint16 {
	return uint16(rand.Intn(int(max-min))) + min
}

func (w *worker) Init(ctx context.Context, minRegionId db.RegionId, maxRegionId db.RegionId, requestId api.RequestId) error {

	w.executions[requestId] = executionData{make(PriorityQueue, 0), make(map[db.VertexId]*Item), make(map[db.VertexId]bool)}
	return nil

	// if minRegionId >= maxRegionId {
	// 	w.executions[requestId] = executionData{make(PriorityQueue, 0), make(map[db.VertexId]*Item), make(map[db.VertexId]bool), nil, nil}
	// 	return nil
	// }
	// leftChildId := minRegionId
	// rightChildId := (minRegionId+maxRegionId)/2 + 1

	// var leftChild, rightChild link.Link
	// var err error
	// _, leftChild, err = w.links[leftChildId].GetLink()
	// if err != nil {
	// 	return err
	// }
	// _, rightChild, err = w.links[rightChildId].GetLink()
	// if err != nil {
	// 	return err
	// }
	// w.executions[requestId] = executionData{make(PriorityQueue, 0), make(map[db.VertexId]*Item, 0), make(map[db.VertexId]bool), leftChild, rightChild}
	// err = leftChild.Init(ctx, minRegionId+1, (minRegionId+maxRegionId)/2, requestId)
	// if err != nil {
	// 	return err
	// }
	// err = rightChild.Init(ctx, (minRegionId+maxRegionId)/2+2, maxRegionId, requestId)
	// if err != nil {
	// 	return err
	// }
	// return nil
}

func (w *worker) Step(ctx context.Context, vertexId db.VertexId, distance float64, requestId api.RequestId) (db.VertexId, float64, error) {
	log.Println("vertex ", vertexId, " distance ", distance, " request id ", requestId)
	queue, inQueue, processed := w.executions[requestId].queue, w.executions[requestId].inQueue, w.executions[requestId].processed

	if w.data.edgeTargets[vertexId] {
		processed[vertexId] = true
	}

	for _, edge := range w.data.edges[vertexId] {
		if !processed[edge.To] {
			item, prs := inQueue[edge.To]
			if !prs {
				newItem := &Item{id: edge.To, distance: distance + edge.Length}
				inQueue[edge.To] = newItem
				heap.Push(&queue, newItem)
			} else if distance+edge.Length < item.distance {
				queue.update(item, distance+edge.Length)
			}
		}
	}

	outVertexId, outDistance := db.VertexId(0), math.MaxFloat64
	for queue.Len() > 0 {
		item := heap.Pop(&queue).(*Item)
		if !processed[item.id] {
			outVertexId, outDistance = item.id, item.distance
			break
		}
	}

	// if w.executions[requestId].leftChild != nil {
	// 	leftVertexId, leftDistance, err := w.executions[requestId].leftChild.Step(ctx, vertexId, distance, requestId)
	// 	if err != nil {
	// 		return 0, 0, err
	// 	}
	// 	if leftDistance < outDistance {
	// 		outVertexId, outDistance = leftVertexId, leftDistance
	// 	}
	// }
	// if w.executions[requestId].rightChild != nil {
	// 	rightVertexId, rightDistance, err := w.executions[requestId].rightChild.Step(ctx, vertexId, distance, requestId)
	// 	if err != nil {
	// 		return 0, 0, err
	// 	}
	// 	if rightDistance < outDistance {
	// 		outVertexId, outDistance = rightVertexId, rightDistance
	// 	}
	// }
	return outVertexId, outDistance, nil
}

func (w *worker) ShortestPath(ctx context.Context, args service.ShortestPathArgs) (service.ShortestPathResult, error) {
	var err error

	w.executionLinks[args.RequestId] = make([]link.Link, len(w.links))
	i := 0
	for _, regionManager := range w.links {
		_, w.executionLinks[args.RequestId][i], err = regionManager.GetLink()
		if err != nil {
			return service.ShortestPathResult{}, err
		}
		w.executionLinks[args.RequestId][i].Init(ctx, 0, 0, args.RequestId)
		i++
	}

	// queue, inQueue := w.executions[args.RequestId].queue, w.executions[args.RequestId].inQueue
	// newItem := &Item{id: args.From, distance: 0}
	// inQueue[args.From] = newItem
	// heap.Push(&queue, newItem)

	vertexId, distance := args.From, float64(0)
	for vertexId != args.To || distance != math.MaxFloat64 {
		newVertexId, newDistance := db.VertexId(0), math.MaxFloat64
		for _, link := range w.executionLinks[args.RequestId] {
			linkVertexId, linkDistance, linkErr := link.Step(ctx, vertexId, distance, args.RequestId)
			if linkErr != nil {
				return service.ShortestPathResult{}, linkErr
			}
			if linkDistance < newDistance {
				newVertexId, newDistance = linkVertexId, linkDistance
			}
		}
		vertexId, distance = newVertexId, newDistance
	}
	return service.ShortestPathResult{}, nil
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
	w.data.edgeTargets = make(map[int64]bool)
	for _, edges := range w.data.edges {
		for _, e := range edges {
			eIds = append(eIds, e.Id)
			w.data.edgeTargets[e.To] = true
		}
	}
	// w.data.arcFlags, err = w.db.GetArcFlags(ctx, eIds, w.generation)
	// if err != nil {
	// 	return
	// }

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
		db:             deps.Db,
		discoverer:     deps.Discoverer,
		regionId:       deps.RegionID,
		linkPort:       deps.LinkPort,
		links:          make(map[db.RegionId]link.RegionManager),
		executions:     make(map[api.RequestId]executionData),
		executionLinks: make(map[api.RequestId][]link.Link),
	}
	if err := w.initDiscoverer(deps.Context); err != nil {
		return nil, err
	}

	return w, nil
}
