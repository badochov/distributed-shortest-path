package managerServer

import (
	"context"
	"github.com/badochov/distributed-shortest-path/src/libs/rpc"
)

type server struct {
	rpc.UnimplementedWorkerServiceServer
}

func (s *server) AssignSegment(ctx context.Context, segment *rpc.Segment) (*rpc.Ack, error) {
	// TODO
	panic("implement me")
}
