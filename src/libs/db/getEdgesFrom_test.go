package db

import (
	"context"
	"time"
)

func (s *dbSuite) Test_getEdgesFrom() {
	const gen = 1
	cases := []struct {
		name          string
		vertices      []Vertex
		edges         []Edge
		numRegions    int
		division      map[RegionId][]VertexId
		expectedEdges map[VertexId][]Edge
	}{
		{
			/*
					1 -e1> 2->
					|\       |e5
				  e3| \e2    |
					4--3-e4> 5
			*/
			"Five vertices in parity regions",
			[]Vertex{
				{Id: 1},
				{Id: 2},
				{Id: 3},
				{Id: 4},
				{Id: 5},
			},
			[]Edge{
				{Id: 1, From: 1, To: 2},
				{Id: 2, From: 1, To: 3},
				{Id: 3, From: 1, To: 4},
				{Id: 4, From: 3, To: 5},
				{Id: 5, From: 2, To: 5},
			},
			2,
			map[RegionId][]VertexId{
				0: {1, 3, 5},
				1: {2, 4},
			},
			map[VertexId][]Edge{
				1: {
					{Id: 1, From: 1, To: 2},
					{Id: 2, From: 1, To: 3},
					{Id: 3, From: 1, To: 4},
				},
				2: {
					{Id: 5, From: 2, To: 5},
				},
				3: {
					{Id: 4, From: 3, To: 5},
				},
			},
		},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			db := s.newDb()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := db.AddVertices(ctx, tc.vertices, gen)
			s.Require().NoError(err)

			err = db.AddEdges(ctx, tc.edges, gen)
			s.Require().NoError(err)

			var vertices []int64
			for _, vertex := range tc.vertices {
				vertices = append(vertices, vertex.Id)
			}

			for i := 0; i < tc.numRegions; i++ {
				edges, err := db.GetEdgesFrom(ctx, vertices, gen)
				s.Require().NoError(err)

				s.Equal(tc.expectedEdges, edges)
			}
		})
	}
}
