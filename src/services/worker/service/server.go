package service

import (
	"context"
	"github.com/badochov/distributed-shortest-path/src/libs/rpc"
	"github.com/badochov/distributed-shortest-path/src/services/worker/common"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Deps struct {
	Listener net.Listener
}

type workerService struct {
	rpc.UnimplementedWorkerServer
}

func (s *workerService) AssignSegment(ctx context.Context, segment *rpc.Segment) (*rpc.Ack, error) {
	// TODO
	panic("implement me")
}

type serv struct {
	server   *grpc.Server
	listener net.Listener
}

func (s *serv) Run() error {
	return s.server.Serve(s.listener)
}

type Service interface {
	common.Runner
}

func New(deps Deps) Service {
	s := grpc.NewServer()

	rpc.RegisterWorkerServer(s, &workerService{})
	log.Printf("server listening at %v", deps.Listener.Addr())

	return &serv{server: s, listener: deps.Listener}
}
