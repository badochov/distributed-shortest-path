package discoverer

import "github.com/badochov/distributed-shortest-path/src/services/manager/common"

type Deps struct {
}

type Discoverer interface {
	common.Runner
}

type discoverer struct {
}

func (d *discoverer) Run() error {
	//TODO implement me
	panic("implement me")
}

func New(deps Deps) Discoverer {
	return &discoverer{}
}
