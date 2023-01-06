package executor

import (
	api "github.com/badochov/distributed-shortest-path/src/libs/api/worker_api"
	"github.com/badochov/distributed-shortest-path/src/services/worker/common"
	"github.com/badochov/distributed-shortest-path/src/services/worker/worker"
	"net/http"
)

type Deps struct {
	Worker worker.Worker
}

type Executor interface {
	ShortestPath(req api.ShortestPathRequest) (resp api.ShortestPathResponse, code int, err error)
	CalculateArcFlags() (resp api.CalculateArcFlagsResponse, code int, err error)

	Healthz() (resp api.HealthzResponse, code int, err error)

	common.Runner
}

type executor struct {
	worker worker.Worker
}

func (e *executor) Run() error {
	return nil
}

func (e *executor) ShortestPath(req api.ShortestPathRequest) (resp api.ShortestPathResponse, code int, err error) {
	resp, err = e.worker.ShortestPath(req)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}
	code = http.StatusOK
	return
}

func (e *executor) CalculateArcFlags() (resp api.CalculateArcFlagsResponse, code int, err error) {
	err = e.worker.CalculateArcFlags()
	if err != nil {
		code = http.StatusInternalServerError
		return
	}
	code = http.StatusOK
	return
}

func (e *executor) Healthz() (resp api.HealthzResponse, code int, err error) {
	code = http.StatusOK
	return
}

func New(deps Deps) Executor {
	return &executor{
		worker: deps.Worker,
	}
}