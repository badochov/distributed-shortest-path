package worker

import (
	"container/heap"
	"context"
	"log"
	"math"
	"sync"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/worker/api"
	"github.com/badochov/distributed-shortest-path/src/services/worker/discoverer"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/link_server"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service"
	"golang.org/x/sync/errgroup"
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
	knownVertices               map[db.VertexId]bool
	arcFlags                    []db.ArcFlag
	neighbouringVerticesRegions map[db.EdgeId]db.RegionId
}

type executionData struct {
	queue     PriorityQueue
	queueItem map[db.VertexId]*Item
	processed map[db.VertexId]bool
	through   map[db.VertexId]db.VertexId
}

type worker struct {
	db             db.DB
	discoverer     discoverer.Discoverer
	generation     db.Generation
	regionId       db.RegionId
	data           workerData
	links          map[db.RegionId]link.RegionManager
	linkPort       int
	executions     map[api.RequestId]*executionData
	executionsLock sync.RWMutex
}

func (w *worker) Finish(ctx context.Context, requestId api.RequestId) error {
	w.executionsLock.Lock()
	defer w.executionsLock.Unlock()
	// TODO delete doesn't really free resources in go. Once upon a while we need to copy the map to free.
	ex := w.executions[requestId]
	*ex = executionData{}
	delete(w.executions, requestId)
	return nil
}

func (w *worker) CalculateArcFlags(ctx context.Context) error {
	//TODO implement me
	return nil
}

func (w *worker) Init(ctx context.Context, requestId api.RequestId) error {
	w.executionsLock.Lock()
	defer w.executionsLock.Unlock()
	w.executions[requestId] = &executionData{
		queueItem: make(map[db.VertexId]*Item),
		processed: make(map[db.VertexId]bool),
		through:   make(map[db.VertexId]db.VertexId),
	}
	return nil
}

func (w *worker) Step(ctx context.Context, vertexId db.VertexId, distance float64, through db.VertexId, requestId api.RequestId) (db.VertexId, float64, db.VertexId, error) {
	log.Println("[STEP] Vertex", vertexId, "distance", distance, "through", through, "Request", requestId)
	w.executionsLock.RLock()
	ex := w.executions[requestId]
	w.executionsLock.RUnlock()

	if w.data.knownVertices[vertexId] {
		ex.through[vertexId] = through
		ex.processed[vertexId] = true
	}

	for _, edge := range w.data.edges[vertexId] {
		if !ex.processed[edge.To] {
			d := distance + edge.Length
			item, inQueue := ex.queueItem[edge.To]
			if !inQueue {
				ex.through[edge.To] = vertexId
				newItem := &Item{id: edge.To, distance: d}
				ex.queueItem[edge.To] = newItem
				heap.Push(&ex.queue, newItem)
			} else if d < item.distance {
				ex.through[edge.To] = vertexId
				ex.queue.update(item, d)
			}
		}
	}

	for ex.queue.Len() > 0 {
		item := ex.queue[0] // Peak
		if !ex.processed[item.id] {
			return item.id, item.distance, ex.through[item.id], nil
		}
		_ = heap.Pop(&ex.queue)
	}
	return db.VertexId(-1), math.Inf(1), db.VertexId(-1), nil
}

func (w *worker) Reconstruct(ctx context.Context, vertexId db.VertexId, requestId api.RequestId) ([]db.VertexId, error) {
	log.Println("[RECONSTRUCT] Vertex", vertexId, "Request", requestId)
	w.executionsLock.RLock()
	ex := w.executions[requestId]
	w.executionsLock.RUnlock()

	var path []db.VertexId
	vertexId, inRegion := ex.through[vertexId]
	for inRegion {
		path = append(path, vertexId)
		vertexId, inRegion = ex.through[vertexId]
	}
	return path, nil
}

//// shortestPathFast is more efficient version of ShortestPath. Albeit much more complicated probably completely not worth it.
//func (w *worker) shortestPathFast(ctx context.Context, args service.ShortestPathArgs) (service.ShortestPathResult, error) {
//	type newGenData struct {
//		distance float64
//		vertexId db.VertexId
//	}
//	type res struct {
//		distance float64
//		vertexId db.VertexId
//		err      error
//	}
//	executionLinks := make([]link.Link, len(w.links))
//	for i, regionManager := range w.links {
//		_, l, err := regionManager.GetLink()
//		if err != nil {
//			return service.ShortestPathResult{}, err
//		}
//		if err := l.Init(ctx, args.RequestId); err != nil {
//			return service.ShortestPathResult{}, err
//		}
//		executionLinks[i] = l
//	}
//
//	newGenChans := make([]chan newGenData, len(executionLinks))
//	resChan := make(chan res)
//	ctx, cancel := context.WithCancel(ctx)
//	for i, l := range executionLinks {
//		newGenChan := make(chan newGenData)
//		newGenChans[i] = newGenChan
//		go func(l link.Link) {
//			for data := range newGenChan {
//				linkVertexId, linkDistance, err := l.Step(ctx, data.vertexId, data.distance, args.RequestId)
//				if err != nil {
//					resChan <- res{
//						err: err,
//					}
//					return
//				}
//				resChan <- res{
//					distance: linkDistance,
//					vertexId: linkVertexId,
//				}
//			}
//		}(l)
//	}
//	// Close new gen chans to stop goroutines.
//	defer func() {
//		for _, ch := range newGenChans {
//			close(ch)
//		}
//	}()
//
//	iters := 0 // For debug purposes.
//	data := newGenData{
//		distance: 0,
//		vertexId: args.From,
//	}
//	for data.vertexId != args.To {
//		// Start Step requests.
//		for _, ch := range newGenChans {
//			ch <- data
//		}
//		data.distance = math.Inf(1)
//
//		// Get data from links.
//		var err error
//		for range newGenChans {
//			r := <-resChan
//			if r.err != nil {
//				if err == nil {
//					// If it was first error cancel the context.
//					cancel()
//				}
//				err = multierror.Append(err, r.err)
//			} else if r.distance < data.distance { // Update next vertex
//				data = newGenData{
//					distance: r.distance,
//					vertexId: r.vertexId,
//				}
//			}
//		}
//		// Check error.
//		if err != nil {
//			// Context was canceled in err handler above.
//			return service.ShortestPathResult{}, err
//		}
//
//		iters++
//		if math.IsInf(data.distance, 1) {
//			log.Println("ITERS", iters)
//			cancel()
//			return service.ShortestPathResult{NoPath: true}, nil
//		}
//	}
//	cancel()
//	return service.ShortestPathResult{Distance: data.distance}, nil
//}

func (w *worker) ShortestPath(ctx context.Context, args service.ShortestPathArgs) (service.ShortestPathResult, error) {
	executionLinks := make([]link.Link, len(w.links))
	for i, regionManager := range w.links {
		_, l, err := regionManager.GetLink()
		if err != nil {
			return service.ShortestPathResult{}, err
		}
		if err := l.Init(ctx, args.RequestId); err != nil {
			return service.ShortestPathResult{}, err
		}
		executionLinks[i] = l
	}

	vertexId, distance, through := args.From, float64(0), db.VertexId(-1)
	iters := 0
	var mutex sync.Mutex
	for vertexId != args.To {
		newVertexId, newDistance, newThrough := db.VertexId(-1), math.Inf(1), db.VertexId(-1)
		errGrp, ctx := errgroup.WithContext(ctx)
		for _, l := range executionLinks {
			l := l
			errGrp.Go(func() error {
				linkVertexId, linkDistance, linkThrough, err := l.Step(ctx, vertexId, distance, through, args.RequestId)
				if err != nil {
					return err
				}
				mutex.Lock() // More go-like would be to send the result via channel and aggregate them in main thread.
				defer mutex.Unlock()
				if linkDistance < newDistance {
					newVertexId, newDistance, newThrough = linkVertexId, linkDistance, linkThrough
				}
				return nil
			})
		}
		if err := errGrp.Wait(); err != nil {
			return service.ShortestPathResult{}, err
		}
		vertexId, distance, through = newVertexId, newDistance, newThrough
		iters++
		if math.IsInf(distance, 1) {
			log.Println("ITERS", iters)
			return service.ShortestPathResult{NoPath: true}, nil
		}
	}

	// Reconstruct the path.
	path := []db.VertexId{vertexId}
	for vertexId != -1 {
		errGrp, ctx := errgroup.WithContext(ctx)
		for _, l := range executionLinks {
			l := l
			errGrp.Go(func() error {
				linkPath, err := l.Reconstruct(ctx, vertexId, args.RequestId)
				if err != nil {
					return err
				}
				if len(linkPath) > 0 {
					mutex.Lock()
					path = append(path, linkPath...)
					mutex.Unlock()
				}
				return nil
			})
		}
		if err := errGrp.Wait(); err != nil {
			return service.ShortestPathResult{}, err
		}

		vertexId = path[len(path)-1]
	}
	path = path[:len(path)-1] // Last element is -1.
	// Path is in reverse order.
	for i := 0; i < len(path)/2; i++ {
		path[i], path[len(path)-1] = path[len(path)-1], path[i]
	}

	// Free resources.
	errGrp, ctx := errgroup.WithContext(ctx)
	for _, l := range executionLinks {
		l := l
		errGrp.Go(func() error {
			return l.Finish(ctx, args.RequestId)
		})
	}
	if err := errGrp.Wait(); err != nil {
		return service.ShortestPathResult{}, err
	}

	return service.ShortestPathResult{Distance: distance, Vertices: path}, nil
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
	w.data.knownVertices = make(map[int64]bool)
	for _, edges := range w.data.edges {
		for _, e := range edges {
			w.data.knownVertices[e.To] = true
			w.data.knownVertices[e.From] = true
		}
	}

	//eIds := make([]db.EdgeId, 0, len(w.data.edges))for _, edges := range w.data.edges {
	//	for _, e := range edges {
	//		eIds = append(eIds, e.Id)
	//
	//		eIds = append(eIds, e.Id)
	//	}
	//}
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
		log.Println(err)
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
		executions: make(map[api.RequestId]*executionData),
	}
	if err := w.initDiscoverer(deps.Context); err != nil {
		return nil, err
	}

	return w, nil
}
