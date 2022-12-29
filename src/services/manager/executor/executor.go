package executor

import (
	"github.com/badochov/distributed-shortest-path/src/services/manager/common"
	"github.com/badochov/distributed-shortest-path/src/services/manager/discoverer"
	"github.com/badochov/distributed-shortest-path/src/services/manager/server/api"
	"gorm.io/gorm"
)

type Deps struct {
	Discoverer discoverer.Discoverer
	Db         *gorm.DB
}

type Executor interface {
	ShortestPath(req api.ShortestPathRequest) (resp api.ShortestPathResponse, code int, err error)
	AddEdges(req api.AddEdgesRequest) (resp api.AddEdgesRequest, code int, err error)
	AddVertices(req api.AddVerticesRequest) (resp api.AddVerticesResponse, code int, err error)
	RecalculateDS() (resp api.RecalculateDsResponse, code int, err error)

	common.Runner
}

type executor struct {
	discoverer discoverer.Discoverer
	db         *gorm.DB
}

func (e *executor) Run() error {
	if err := e.discoverer.Run(); err != nil {
		return nil
	}

	//TODO implement me
	panic("implement me")
}

func (e *executor) ShortestPath(req api.ShortestPathRequest) (resp api.ShortestPathResponse, code int, err error) {
	//TODO implement me
	panic("implement me")
}

func (e *executor) AddEdges(req api.AddEdgesRequest) (resp api.AddEdgesRequest, code int, err error) {
	//TODO implement me
	panic("implement me")
}

func (e *executor) AddVertices(req api.AddVerticesRequest) (resp api.AddVerticesResponse, code int, err error) {
	//TODO implement me
	panic("implement me")
}

func (e *executor) RecalculateDS() (resp api.RecalculateDsResponse, code int, err error) {
	//TODO implement me
	panic("implement me")
}

func New(deps Deps) Executor {
	return &executor{
		discoverer: deps.Discoverer,
		db:         deps.Db,
	}
}
