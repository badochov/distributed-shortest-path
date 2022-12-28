package link

import "context"

type localLink struct {
}

func (l *localLink) Close() error {
	return nil
}

func (l *localLink) Add(ctx context.Context, a, b int32) (int32, error) {
	return a + b, nil
}

var _ Link = &localLink{}

func newLocalLink() *localLink {
	return &localLink{}
}
