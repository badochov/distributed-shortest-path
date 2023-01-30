package link

import (
	"context"
	"io"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/worker/api"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/link.proto

type Link interface {
	// TODO[wprzytula]: retries, consider backoff lib to unify retries in every method here
	Add(ctx context.Context, a, b int32) (int32, error) // Example
	Init(ctx context.Context, requestId api.RequestId) error
	Step(ctx context.Context, vertexId db.VertexId, distance float64, through db.VertexId, reqId api.RequestId) (db.VertexId, float64, db.VertexId, error)

	io.Closer
}

// New creates new link. Currently, nodes connects even to itself via gRPC as it should be fast enough.
func New(ctx context.Context, addr string) (Link, error) {
	return newRemoteLink(ctx, addr)
}
