package service_server

import (
	"context"
	"github.com/badochov/distributed-shortest-path/src/libs/rpc"
	"github.com/badochov/distributed-shortest-path/src/services/worker/common"
	"github.com/badochov/distributed-shortest-path/src/services/worker/worker"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Deps struct {
	Listener net.Listener
	Worker   worker.Worker
}

type workerService struct {
	rpc.UnimplementedWorkerServer
	worker worker.Worker
}

func (s *workerService) AssignSegment(ctx context.Context, segment *rpc.Segment) (*rpc.Ack, error) {
	if err := s.worker.AssignSegment(segment.SegmentId); err != nil {
		return nil, err
	}
	return &rpc.Ack{}, nil
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

	rpc.RegisterWorkerServer(s, &workerService{worker: deps.Worker})
	log.Printf("server listening at %v", deps.Listener.Addr())

	return &serv{server: s, listener: deps.Listener}
}
