package link

import (
	"context"
	"fmt"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	"github.com/badochov/distributed-shortest-path/src/services/worker/api"
	"github.com/badochov/distributed-shortest-path/src/services/worker/link/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (l *remoteLink) Close() error {
	return l.conn.Close()
}

func (l *remoteLink) Add(ctx context.Context, a, b int32) (int32, error) {
	// TODO [wp] think about retries.
	resp, err := l.client.Add(ctx, &proto.AddRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	return resp.Res, nil
}

func (l *remoteLink) Init(ctx context.Context, minRegionId db.RegionId, maxRegionId db.RegionId, requestId api.RequestId) error {
	_, err := l.client.Init(ctx, &proto.InitRequest{MinRegionId: uint32(minRegionId), MaxRegionId: uint32(maxRegionId), RequestId: uint64(requestId)})
	if err != nil {
		return fmt.Errorf("message, %w", err)
	}
	return nil
}

var _ Link = &remoteLink{}

type remoteLink struct {
	client proto.LinkClient
	conn   *grpc.ClientConn
}

func newRemoteLink(ctx context.Context, addr string) (*remoteLink, error) {
	con, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &remoteLink{
		client: proto.NewLinkClient(con),
		conn:   con,
	}, nil
}
