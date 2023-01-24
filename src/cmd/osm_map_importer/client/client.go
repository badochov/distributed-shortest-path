package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/badochov/distributed-shortest-path/src/services/manager/api"
)

type Edge = api.Edge
type Vertex = api.Vertex

type Client struct {
	client  *http.Client
	baseUrl string
}

func (c *Client) AddEdges(edges []Edge) error {
	data, err := json.Marshal(api.AddEdgesRequest{
		Edges: edges,
	})
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.url(api.AddEdgesUrl), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("AddEdges failed with status=%d, and message=%s", resp.StatusCode, resp.Status)
}

func (c *Client) AddVertices(vertices []Vertex) error {
	data, err := json.Marshal(api.AddVerticesRequest{
		Vertices: vertices,
	})
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.url(api.AddVerticesUrl), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("AddVertices failed with status=%d, and message=%s", resp.StatusCode, resp.Status)
}

func (c *Client) url(path string) string {
	return fmt.Sprintf("http://%s%s", c.baseUrl, path)
}

func New(client *http.Client, url string) *Client {
	return &Client{client: client, baseUrl: url}
}
