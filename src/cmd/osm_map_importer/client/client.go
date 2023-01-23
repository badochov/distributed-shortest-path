package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/badochov/distributed-shortest-path/src/libs/api/manager_api"
)

type Edge = manager_api.Edge
type Vertex = manager_api.Vertex

type Client struct {
	client  *http.Client
	baseUrl string
}

func (c *Client) AddEdges(edges []Edge) error {
	data, err := json.Marshal(manager_api.AddEdgesRequest{
		Edges: edges,
	})
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.url(manager_api.AddEdgesUrl), "application/json", bytes.NewReader(data))
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
	data, err := json.Marshal(manager_api.AddVerticesRequest{
		Vertices: vertices,
	})
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.url(manager_api.AddVerticesUrl), "application/json", bytes.NewReader(data))
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
