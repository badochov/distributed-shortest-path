package worker

import (
	"encoding/json"
	"fmt"
	"github.com/badochov/distributed-shortest-path/src/services/worker/service_server/api"
	"net/http"
)

type ShortestPathArgs = api.ShortestPathRequest
type ShortestPathResult = api.ShortestPathResponse

type Client interface {
	CalculateArcFlags() error
	ShortestPath(args ShortestPathArgs) (ShortestPathResult, error)
}

type client struct {
	client *http.Client
}

func (c *client) CalculateArcFlags() error {
	r, err := c.client.Get(api.CalculateArcFlagsUrl)
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

func (c *client) ShortestPath(args ShortestPathArgs) (res ShortestPathResult, err error) {
	r, err := c.client.Get(api.CalculateArcFlagsUrl)
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

// NewClient sets up client to a worker service.
func NewClient(cl *http.Client) Client {
	return &client{client: cl}
}
