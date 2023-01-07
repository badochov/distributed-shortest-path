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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(api.CalculateArcFlagsUrl), nil)
	if err != nil {
		return fmt.Errorf("error creating CalculateArcFlags request, %w", err)
	}
	r, err := c.client.Do(req)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url(api.ShortestPathUrl), bytes.NewBuffer(data))
	if err != nil {
		return ShortestPathResult{}, fmt.Errorf("error creating ShortestPath request, %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := c.client.Do(req)
	if err != nil {
		return ShortestPathResult{}, fmt.Errorf("error performing ShortestPath request, %w", err)
	}
	defer func() {
		if errClose := r.Body.Close(); errClose != nil {
			err = fmt.Errorf("error closing ShortestPath request body, %w", errClose)
		}
	}()

	if r.StatusCode != http.StatusOK {
		return res, fmt.Errorf("worker server responded to ShortestPath with code=%d, status %s", r.StatusCode, r.Status)
	}

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return ShortestPathResult{}, fmt.Errorf("error decoding ShortestPath request body, %w", err)
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
