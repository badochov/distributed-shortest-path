package worker

import (
	"context"
	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/worker/discoverer"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service_server/api"
)

type Deps struct {
	Db         *db.DB
	Discoverer discoverer.Discoverer
}

type ShortestPathArgs = api.ShortestPathRequest
type ShortestPathResult = api.ShortestPathResponse

// Worker All methods from link service and worker service should end up calling this interface.
type Worker interface {
	// Server
	CalculateArcFlags() error
	ShortestPath(args ShortestPathArgs) (ShortestPathResult, error)
	
	// Link
	Add(a int32, b int32) (int32, error) // Example

	Run(ctx context.Context) error
}

type worker struct {
	db *db.DB
	discoverer discoverer.Discoverer
}

func (w *worker) CalculateArcFlags() error {
	//TODO implement me
	panic("implement me")
}

func (w *worker) ShortestPath(args ShortestPathArgs) (ShortestPathResult, error) {
	//TODO implement me
	panic("implement me")
}

func (w *worker) Run(ctx context.Context) error {
	return w.discoverer.Run(ctx)
}

func (w *worker) Add(a int32, b int32) (int32, error) {
	return a + b, nil
}

func (w *worker) AssignSegment(id int32) error {
	//TODO implement me
	panic("implement me")
}

func New(deps Deps) Worker {
	return &worker{
		db: deps.Db,
		discoverer: deps.Discoverer,
	}
}
