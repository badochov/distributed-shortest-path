package worker_client

import (
	"context"
	"github.com/badochov/distributed-shortest-path/src/libs/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

type WorkerClient interface {
	rpc.WorkerClient
	io.Closer
}

type client struct {
	rpc.WorkerClient
	conn *grpc.ClientConn
}

func (c client) Close() error {
	return c.conn.Close()
}

// New open new grpc connection to WorkerService. Must be closed after use.
func New(ctx context.Context, addr string) (WorkerClient, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return client{WorkerClient: rpc.NewWorkerClient(conn), conn: conn}, nil
}
