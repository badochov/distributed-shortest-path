package link

import (
	"context"
	"io"
)

type Link interface {
	Add(ctx context.Context, a, b int32) (int32,error) // Example

	io.Closer
}
