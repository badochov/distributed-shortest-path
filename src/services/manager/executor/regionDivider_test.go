package executor

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
	dbTesting "github.com/badochov/distributed-shortest-path/src/libs/db/testing"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"
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
		name       string
		vertices   []db.Vertex
		numRegions int
		division   map[regionId][]db.VertexId
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
			map[regionId][]db.VertexId{
				0: {1},
				1: {2},
				2: {3},
				3: {4},
			},
		},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			db := s.newDb()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			//ctx := context.Background()
			err := db.AddVertices(ctx, tc.vertices, gen)
			//return
			s.Require().NoError(err)
			rg := regionDivider{
				db:         db,
				generation: gen,
			}
			err = rg.divideIntoRegions(ctx, 16)
			s.Require().NoError(err)

			for i := 0; i < tc.numRegions; i++ {
				ids, err := db.GetVertexIds(ctx, regionId(i), gen)
				s.Require().NoError(err)
				slices.Sort(ids)
				div := tc.division[regionId(i)]
				slices.Sort(div)
				s.Equal(div, ids)
			}
		})
	}
}
