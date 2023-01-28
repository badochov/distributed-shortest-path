package executor

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	dbTesting "github.com/badochov/distributed-shortest-path/src/libs/db/testing"
	"github.com/stretchr/testify/suite"
)

type dbSuite struct {
	suite.Suite

	fs []*os.File
}

func Test_RegionDivider_Suite(t *testing.T) {
	suite.Run(t, new(dbSuite))
}

func (s *dbSuite) SetupSuite() {
}

func (s *dbSuite) newDb() db.DB {
	path, err := os.CreateTemp(os.TempDir(), "region-divider-test.*.db")
	s.Require().NoError(err)
	s.fs = append(s.fs, path)
	c, err := dbTesting.NewMockConn(path.Name())
	s.Require().NoError(err)
	cl, err := db.Connect(c)
	s.Require().NoError(err)
	return cl
}

func (s *dbSuite) TearDownTest() {
	for _, f := range s.fs {
		err := f.Close()
		s.Require().NoError(err)
	}
}

func (s *dbSuite) Test_RegionDivider_DivideIntoRegions() {
	const gen = 1
	cases := []struct {
		name            string
		vertices        []db.Vertex
		numRegions      int
		acceptableDelta int
	}{
		{
			"Four vertices each to its own region",
			[]db.Vertex{
				{
					Id:        1,
					Latitude:  1,
					Longitude: 1,
				},
				{
					Id:        2,
					Latitude:  -1,
					Longitude: 1,
				},
				{
					Id:        3,
					Latitude:  -1,
					Longitude: -1,
				},
				{
					Id:        4,
					Latitude:  1,
					Longitude: -1,
				},
			},
			4,
			0,
		},
		{
			"Four vertices two per region",
			[]db.Vertex{
				{
					Id:        1,
					Latitude:  1,
					Longitude: 1,
				},
				{
					Id:        2,
					Latitude:  -1,
					Longitude: 1,
				},
				{
					Id:        3,
					Latitude:  -1,
					Longitude: -1,
				},
				{
					Id:        4,
					Latitude:  1,
					Longitude: -1,
				},
			},
			2,
			0,
		},
		{
			"Three vertices two regions",
			[]db.Vertex{
				{
					Id:        1,
					Latitude:  1,
					Longitude: 1,
				},
				{
					Id:        2,
					Latitude:  -1,
					Longitude: 1,
				},
				{
					Id:        4,
					Latitude:  1,
					Longitude: -1,
				},
			},
			2,
			1,
		},
		{
			"Five vertices four regions",
			[]db.Vertex{
				{
					Id:        1,
					Latitude:  1,
					Longitude: 1,
				},
				{
					Id:        2,
					Latitude:  -1,
					Longitude: 1,
				},
				{
					Id:        3,
					Latitude:  -1,
					Longitude: -1,
				},
				{
					Id:        4,
					Latitude:  1,
					Longitude: -1,
				},
			},
			2,
			1,
		},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			conn := s.newDb()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := conn.AddVertices(ctx, tc.vertices, gen)
			s.Require().NoError(err)
			rg := regionDivider{
				db:         conn,
				generation: gen,
			}
			err = rg.divideIntoRegions(ctx, tc.numRegions)
			s.Require().NoError(err)

			min, max := len(tc.vertices), 0
			var count int
			m := map[regionId][]db.VertexId{}
			for i := 0; i < tc.numRegions; i++ {
				ids, err := conn.GetVertexIds(ctx, regionId(i), gen)
				s.Require().NoError(err)

				m[regionId(i)] = ids
				if len(ids) < min {
					min = len(ids)
				} else if len(ids) > max {
					max = len(ids)
				}
				count += len(ids)
			}
			s.LessOrEqual(max-min, tc.acceptableDelta, "delta higher that expected, received division: ", m)
			s.Equal(len(tc.vertices), count, "number of vertices assignments doesn't match number of vertices, received division: ", m)
		})
	}
}
