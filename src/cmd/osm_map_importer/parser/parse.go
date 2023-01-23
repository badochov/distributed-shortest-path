package parser

import (
	"github.com/badochov/distributed-shortest-path/src/cmd/osm_map_importer/client"
	"github.com/glaslos/go-osm"
)

func Parse(path string) ([]client.Vertex, []client.Edge, error) {
	data, err := osm.DecodeFile(path)
	if err != nil {
		return nil, nil, err
	}

	vertices := make([]client.Vertex, 0, len(data.Nodes))
	for _, n := range data.Nodes {
		vertices = append(vertices, client.Vertex{
			Id:        n.ID,
			Latitude:  n.Lat,
			Longitude: n.Lng,
		})
	}

	edges := make([]client.Edge, 0, len(data.Ways))
	for _, w := range data.Ways {
		currN := w.Nds[0]
		``
		for _, n := range w.Nds[1:] {
			edges = append(edges, client.Edge{
				From: currN.ID,
				To:   n.ID,
				Id:   w.ID,
			})
			currN = n
		}
	}

	return vertices, edges, nil
}
