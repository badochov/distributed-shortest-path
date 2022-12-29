package link_server

import (
	"github.com/badochov/distributed-shortest-path/src/services/worker/common"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/proto"
	"github.com/badochov/distributed-shortest-path/src/services/worker/worker"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Deps struct {
	Listener net.Listener
	Worker   worker.Worker
}

type linkService struct {
	proto.UnimplementedLinkServer
	worker worker.Worker
}

func (s *linkService) Add(ctx context.Context, req *proto.AddRequest) (*proto.AddResponse, error) {
	res, err := s.worker.Add(req.A, req.B)
	if err != nil {
		return nil, err
	}
	return &proto.AddResponse{Res: res}, nil
}

type serv struct {
	server   *grpc.Server
	listener net.Listener
}

func (s *serv) Run() error {
	return s.server.Serve(s.listener)
}

type Server interface {
	common.Runner
}

func New(deps Deps) Server {
	s := grpc.NewServer()

	proto.RegisterLinkServer(s, &linkService{worker: deps.Worker})
	log.Printf("server listening at %v", deps.Listener.Addr())

	return &serv{server: s, listener: deps.Listener}
}
