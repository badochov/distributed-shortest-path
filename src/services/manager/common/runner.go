package common

import "context"

type Runner interface {
	// Run starts component. ctx is meant to be used for startup only.
	Run(ctx context.Context) error
}
