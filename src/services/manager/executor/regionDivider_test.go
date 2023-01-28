package executor

import (
	"context"
	"os"
	"sort"
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

type regions [][]db.VertexId

func (r regions) Len() int {
	return len(r)
}
func (r regions) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r regions) Less(i, j int) bool {
	for k := 0; k < len(r[i]) && k < len(r[j]); k++ {
		if r[i][k] < r[j][k] {
			return true
		} else if r[i][k] > r[j][k] {
			return false
		}
	}
	return len(r[i]) < len(r[j])
}

func (s *dbSuite) Test_RegionDivider_DivideIntoRegions() {
	const gen = 1
	cases := []struct {
		name       string
		vertices   []db.Vertex
		numRegions int
		division   regions
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
			regions{
				{1},
				{2},
				{3},
				{4},
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
			err = rg.divideIntoRegions(ctx, tc.numRegions)
			s.Require().NoError(err)

			division := make(regions, tc.numRegions)
			for i := 0; i < tc.numRegions; i++ {
				ids, err := db.GetVertexIds(ctx, regionId(i), gen)
				s.Require().NoError(err)
				slices.Sort(ids)
				division[i] = ids
			}

			sort.Sort(division)
			sort.Sort(tc.division)
			s.Equal(division, tc.division)
		})
	}
}
