package link_server

import (
	"log"
	"net"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/worker/api"
	"github.com/badochov/distributed-shortest-path/src/services/worker/common"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Worker interface {
	Add(ctx context.Context, a, b int32) (int32, error) // Example
	Init(ctx context.Context, minRegionId db.RegionId, maxRegionId db.RegionId, requestId api.RequestId) error
	Min(ctx context.Context, requestId api.RequestId) (bool, float64, db.VertexId, error)
	Step(ctx context.Context, vertexId db.VertexId, destId db.VertexId, requestId api.RequestId) (bool, float64, error)
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

func (s *linkService) Init(ctx context.Context, req *proto.InitRequest) (*proto.InitResponse, error) {
	err := s.worker.Init(ctx, db.RegionId(req.MinRegionId), db.RegionId(req.MaxRegionId), api.RequestId(req.RequestId))
	if err != nil {
		return nil, err
	}
	return &proto.InitResponse{}, nil
}

func (s *linkService) Min(ctx context.Context, req *proto.MinRequest) (*proto.MinResponse, error) {
	isSet, distance, vertexId, err := s.worker.Min(ctx, api.RequestId(req.RequestId))
	if err != nil {
		return nil, err
	}
	return &proto.MinResponse{IsSet: isSet, Distance: distance, VertexId: vertexId}, nil
}

func (s *linkService) Step(ctx context.Context, req *proto.StepRequest) (*proto.StepResponse, error) {
	found, distance, err := s.worker.Step(ctx, req.VetrtexId, req.DestId, api.RequestId(req.RequestId))
	if err != nil {
		return nil, err
	}
	return &proto.StepResponse{Found: found, Distance: distance}, nil
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
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zapLogger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zapLogger),
		)),
	)

	proto.RegisterLinkServer(s, &linkService{worker: deps.Worker})
	log.Printf("server listening at %v", deps.Listener.Addr())

	return &serv{server: s, listener: deps.Listener}
}
