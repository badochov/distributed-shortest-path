package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	api "github.com/badochov/distributed-shortest-path/src/libs/api/worker_api"
	"net/http"
)

type ShortestPathArgs = api.ShortestPathRequest
type ShortestPathResult = api.ShortestPathResponse

type Client interface {
	CalculateArcFlags(ctx context.Context) error
	ShortestPath(ctx context.Context, args ShortestPathArgs) (ShortestPathResult, error)
}

type client struct {
	client  *http.Client
	baseUrl string
}

func (c *client) CalculateArcFlags(ctx context.Context) error {
	r, err := c.client.Get(c.url(api.CalculateArcFlagsUrl))
	if err != nil {
		return fmt.Errorf("error performing CalculateArcFlags request, %w", err)
	}
	if err := r.Body.Close(); err != nil {
		return fmt.Errorf("error closing CalculateArcFlags request body, %w", err)
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("worker server responded to CalculateArcFlags with code=%d, status %s", r.StatusCode, r.Status)
	}
	return nil
}

func (c *client) ShortestPath(ctx context.Context, args ShortestPathArgs) (res ShortestPathResult, err error) {
	data, err := json.Marshal(args)
	if err != nil {
		return ShortestPathResult{}, err
	}
	r, err := c.client.Post(c.url(api.ShortestPathUrl), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return ShortestPathResult{}, fmt.Errorf("error performing CalculateArcFlags request, %w", err)
	}
	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			err = fmt.Errorf("error closing CalculateArcFlags request body, %w", errClose)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return res, fmt.Errorf("worker server responded to CalculateArcFlags with code=%d, status %s", r.StatusCode, r.Status)
	}

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return ShortestPathResult{}, fmt.Errorf("error decoding CalculateArcFlags request body, %w", err)
	}

	return
}

func (c *client) url(path string) string {
	return fmt.Sprintf("http://%s%s", c.baseUrl, path)
}

type Deps struct {
	HttpClient *http.Client
	Url        string
}

// NewClient sets up client to a worker service.
func NewClient(deps Deps) Client {
	return &client{client: deps.HttpClient, baseUrl: deps.Url}
}
