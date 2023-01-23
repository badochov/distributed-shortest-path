package link

import (
	"context"
	"io"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/link.proto

type Link interface {
	Add(ctx context.Context, a, b int32) (int32, error) // Example

	io.Closer
}

// New creates new link. Currently, nodes connects even to itself via gRPC as it should be fast enough.
func New(ctx context.Context, addr string) (Link, error) {
	return newRemoteLink(ctx, addr)
}
