package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/badochov/distributed-shortest-path/src/libs/db/model"
	dbTesting "github.com/badochov/distributed-shortest-path/src/libs/db/testing"
	"github.com/stretchr/testify/suite"
)

type dbSuite struct {
	suite.Suite
}

func Test_EdgeToRegionMapping_Suite(t *testing.T) {
	suite.Run(t, new(dbSuite))
}

func (s *dbSuite) SetupSuite() {
}

func (s *dbSuite) newDb() db {
	f, err := os.CreateTemp(os.TempDir(), "edge-to-region-mapping-test.*.db")
	s.Require().NoError(err)
	path := f.Name()
	s.Require().NoError(f.Close())
	c, err := dbTesting.NewMockConn(path)
	s.Require().NoError(err)
	cl, err := Connect(c)
	s.Require().NoError(err)
	return cl.(db)
}

func (s *dbSuite) Test_EdgeToRegionMapping() {
	const gen = 1
	cases := []struct {
		name             string
		vertices         []Vertex
		edges            []Edge
		numRegions       int
		division         map[RegionId][]VertexId
		expectedMappings map[RegionId]map[EdgeId]RegionId
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
			map[RegionId]map[EdgeId]RegionId{
				0: {
					1: 1,
					2: 0,
					3: 1,
					4: 0,
				},
				1: {
					5: 0,
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

			for _, vertex := range tc.vertices {
				err = db.q.RegionBinding.WithContext(ctx).Create(&model.RegionBinding{
					VertexID:   vertex.Id,
					Region:     uint16((vertex.Id%2 + 1) % 2),
					Generation: gen,
				})
				s.Require().NoError(err)
			}

			for i := 0; i < tc.numRegions; i++ {
				mapping, err := db.GetEdgeToRegionMapping(ctx, RegionId(i), gen)
				s.Require().NoError(err)

				s.Equal(tc.expectedMappings[RegionId(i)], mapping)
			}
		})
	}
}
