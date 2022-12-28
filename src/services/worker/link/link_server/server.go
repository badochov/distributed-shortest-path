package link_server

import (
	"github.com/badochov/distributed-shortest-path/src/services/worker/common"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Deps struct {
	Listener net.Listener
}

type linkService struct {
	proto.UnimplementedLinkServer
}

func (s *linkService) Add(ctx context.Context, req *proto.AddRequest) (*proto.AddResponse, error) {
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

type Server interface {
	common.Runner
}

func New(deps Deps) Server {
	s := grpc.NewServer()

	proto.RegisterLinkServer(s, &linkService{})
	log.Printf("server listening at %v", deps.Listener.Addr())

	return &serv{server: s, listener: deps.Listener}
}
