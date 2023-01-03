package link

import (
	"context"
	"io"
)

type Link interface {
	Add(ctx context.Context, a, b int32) (int32, error) // Example

	io.Closer
}

// New creates new link. Currently, nodes connects even to itself via gRPC as it should be fast enough.
func New(ctx context.Context, addr string) (Link, error) {
	return newRemoteLink(ctx, addr)
}
