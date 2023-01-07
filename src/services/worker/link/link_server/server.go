package link_server

import (
	"github.com/badochov/distributed-shortest-path/src/services/worker/common"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Worker interface {
	Add(ctx context.Context, a, b int32) (int32, error) // Example
}

type Deps struct {
	Listener net.Listener
	Worker   Worker
}

type linkService struct {
	proto.UnimplementedLinkServer
	worker Worker
}

func (s *linkService) Add(ctx context.Context, req *proto.AddRequest) (*proto.AddResponse, error) {
	res, err := s.worker.Add(ctx, req.A, req.B)
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
