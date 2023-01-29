package parser

import (
	"math"

	"github.com/badochov/distributed-shortest-path/src/cmd/osm_map_importer/client"
	"github.com/badochov/distributed-shortest-path/src/services/manager/api"
	"github.com/glaslos/go-osm"
)

func Parse(path string) ([]client.Vertex, []client.Edge, error) {
	vertices, edges, err := getWithoutLength(path)
	if err != nil {
		return nil, nil, err
	}
	calcLength(vertices, edges)
	return vertices, edges, nil
}

func calcLength(vertices []client.Vertex, edges []client.Edge) {
	verticesMap := make(map[api.VertexId]client.Vertex)
	for _, v := range vertices {
		verticesMap[v.Id] = v
	}

	for i := range edges {
		vFrom := verticesMap[edges[i].From]
		vTo := verticesMap[edges[i].To]
		edges[i].Length = calcDistance(vFrom, vTo)
	}
}

func calcDistance(from client.Vertex, to client.Vertex) api.EdgeLength {
	// Source: https://www.movable-type.co.uk/scripts/latlong.html.
	R := 6371e3                         // metres
	φ1 := from.Latitude * math.Pi / 180 // φ, λ in radians
	φ2 := to.Latitude * math.Pi / 180   // φ, λ in radians
	Δφ := (to.Latitude - from.Latitude) * math.Pi / 180
	Δλ := (to.Longitude - from.Longitude) * math.Pi / 180
	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c // in metres
}

func getWithoutLength(path string) ([]client.Vertex, []client.Edge, error) {
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
