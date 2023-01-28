package db

import (
	"context"
	"fmt"

	"github.com/badochov/distributed-shortest-path/src/libs/db/conn"
	"github.com/badochov/distributed-shortest-path/src/libs/db/model"
	"github.com/badochov/distributed-shortest-path/src/libs/db/query"
	"github.com/badochov/distributed-shortest-path/src/services/manager/api"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Vertex = api.Vertex
type Edge = api.Edge
type Generation = api.Generation

type VertexId = api.VertexId
type EdgeId = api.EdgeId

type ArcFlag uint64

type RegionId = uint16

func (af ArcFlag) IsSet(id RegionId) bool {
	return af&af.getSelector(id) != 0
}

func (af ArcFlag) set(id RegionId) uint64 {
	return uint64(af | af.getSelector(id))
}

func (af ArcFlag) getSelector(id RegionId) ArcFlag {
	return 1 << id
}

type MinMax struct {
	Min float64 // Inclusive
	Max float64 // Exclusive
}

type CoordBounds struct {
	Latitude  MinMax
	Longitude MinMax
}

type DB interface {
	AddVertices(ctx context.Context, vertices []Vertex, generation Generation) error
	AddEdges(ctx context.Context, edges []Edge, generation Generation) error
	AddArcFlags(ctx context.Context, edgeIds []EdgeId, generation Generation) error

	GetVertexIds(ctx context.Context, id RegionId, generation Generation) ([]VertexId, error)
	GetEdgesFrom(ctx context.Context, from []VertexId, generation Generation) ([]Edge, error)
	GetArcFlags(ctx context.Context, edgeIds []EdgeId, generation Generation) ([]ArcFlag, error)
	GetVertexCount(ctx context.Context, coordsBetween CoordBounds, generation Generation) (int64, error)
	GetVertexCountOnVerticalSegment(ctx context.Context, latitude MinMax, longtitude float64, generation Generation) (int64, error)
	GetVertexCountOnHorizontalSegment(ctx context.Context, latitude float64, longtitude MinMax, generation Generation) (int64, error)
	GetVertexRegion(ctx context.Context, id VertexId, generation Generation) (RegionId, error)
	GetCurrentGeneration(ctx context.Context) (Generation, error)
	GetNextGeneration(ctx context.Context) (Generation, error)
	GetActiveGeneration(ctx context.Context) (Generation, error)

	SetFlag(ctx context.Context, edgeIds []EdgeId, region RegionId, generation Generation) error
	SetRegion(ctx context.Context, coordsBetween CoordBounds, region RegionId, generation Generation) (rowsAffected int64, err error)
	SetCurrentGeneration(ctx context.Context, generation Generation) error
	SetNextGeneration(ctx context.Context, generation Generation) error
	SetActiveGeneration(ctx context.Context, generation Generation) error
}

type db struct {
	q *query.Query
}

func (d db) GetVertexIds(ctx context.Context, id RegionId, generation Generation) ([]VertexId, error) {
	rb := d.q.RegionBinding

	rbs, err := d.q.WithContext(ctx).RegionBinding.Where(rb.Generation.Eq(generation), rb.Region.Eq(id)).Find()
	if err != nil {
		return nil, err
	}

	vs := make([]VertexId, 0, generation)
	for _, rb := range rbs {
		vs = append(vs, rb.VertexID)
	}

	return vs, nil
}

func (d db) GetEdgesFrom(ctx context.Context, from []VertexId, generation Generation) ([]Edge, error) {
	e := d.q.Edge

	edges, err := d.q.WithContext(ctx).Edge.Where(e.Generation.Eq(generation), e.FromId.In(from...)).Find()
	if err != nil {
		return nil, err
	}

	es := make([]Edge, 0, generation)
	for _, edge := range edges {
		es = append(es, Edge{
			From:   edge.FromId,
			To:     edge.ToId,
			Id:     edge.ID,
			Length: edge.Length,
		})
	}

	return es, nil
}

func (d db) SetNextGeneration(ctx context.Context, generation Generation) error {
	g := d.q.Generation
	c, err := d.q.WithContext(ctx).Generation.Where(g.Next).Update(g.Next, generation)
	if err != nil {
		return err
	}
	if c.RowsAffected == 0 {
		return d.q.WithContext(ctx).Generation.Create(&model.Generation{
			Generation: generation,
			Next:       true,
		})
	}
	return nil
}

func (d db) SetActiveGeneration(ctx context.Context, generation Generation) error {
	g := d.q.Generation
	c, err := d.q.WithContext(ctx).Generation.Where(g.Active).Update(g.Next, generation)
	if err != nil {
		return err
	}
	if c.RowsAffected == 0 {
		return d.q.WithContext(ctx).Generation.Create(&model.Generation{
			Generation: generation,
			Active:     true,
		})
	}
	return nil
}

func (d db) GetCurrentGeneration(ctx context.Context) (Generation, error) {
	g := d.q.Generation
	r, err := d.q.WithContext(ctx).Generation.Where(g.Current).Attrs(g.Current, g.Generation.Zero()).FirstOrCreate()
	if err != nil {
		return 0, err
	}
	return r.Generation, nil
}

func (d db) GetNextGeneration(ctx context.Context) (Generation, error) {
	g := d.q.Generation
	r, err := d.q.WithContext(ctx).Generation.Where(g.Next).Attrs(g.Next, g.Generation.Zero()).FirstOrCreate()
	if err != nil {
		return 0, err
	}
	return r.Generation, nil
}

func (d db) GetActiveGeneration(ctx context.Context) (Generation, error) {
	g := d.q.Generation
	r, err := d.q.WithContext(ctx).Generation.Where(g.Active).Attrs(g.Next, g.Generation.Zero()).FirstOrCreate()
	if err != nil {
		return 0, err
	}
	return r.Generation, nil
}

func (d db) SetCurrentGeneration(ctx context.Context, generation Generation) error {
	g := d.q.Generation
	c, err := d.q.WithContext(ctx).Generation.Where(g.Current).Update(g.Current, generation)
	if err != nil {
		return err
	}
	if c.RowsAffected == 0 {
		return d.q.WithContext(ctx).Generation.Create(&model.Generation{
			Generation: generation,
			Current:    true,
		})
	}
	return nil
}

func (d db) DeleteNextGeneration(ctx context.Context) error {
	g := d.q.Generation
	_, err := d.q.WithContext(ctx).Generation.Where(g.Next).Attrs(g.Next, g.Generation.Zero()).Delete()
	return err
}

func (d db) GetVertexRegion(ctx context.Context, id VertexId, generation Generation) (RegionId, error) {
	rb := d.q.RegionBinding
	r, err := d.q.WithContext(ctx).RegionBinding.Where(rb.VertexID.Eq(id), rb.Generation.Eq(generation)).First()
	if err != nil {
		return 0, err
	}
	return r.Region, nil
}

func (d db) SetRegion(ctx context.Context, c CoordBounds, region RegionId, generation Generation) (int64, error) {
	v := d.q.Vertex
	const batchSize = 1024

	var results []*model.Vertex
	var rowsAffected int64

	err := d.q.WithContext(ctx).Vertex.Where(
		v.Generation.Eq(generation),
		v.Latitude.Gte(c.Latitude.Min),
		v.Latitude.Lt(c.Latitude.Max),
		v.Longitude.Gte(c.Longitude.Min),
		v.Longitude.Lt(c.Longitude.Max),
	).FindInBatches(&results, batchSize, func(_ gen.Dao, _ int) error {
		rows := len(results)
		rowsAffected += int64(rows)

		rbs := make([]*model.RegionBinding, 0, rows)
		for _, r := range results {
			rbs = append(rbs, &model.RegionBinding{
				VertexID:   r.ID,
				Region:     region,
				Generation: generation,
			})
		}
		return d.q.WithContext(ctx).RegionBinding.Create(rbs...)
	})

	return rowsAffected, err
}

func (d db) GetVertexCount(ctx context.Context, c CoordBounds, generation Generation) (int64, error) {
	v := d.q.Vertex
	return d.q.WithContext(ctx).Vertex.Where(
		v.Generation.Eq(generation),
		v.Latitude.Gte(c.Latitude.Min),
		v.Latitude.Lt(c.Latitude.Max),
		v.Longitude.Gte(c.Longitude.Min),
		v.Longitude.Lt(c.Longitude.Max),
	).Count()
}

func (d db) GetVertexCountOnVerticalSegment(ctx context.Context, latitude MinMax, longtitude float64, generation Generation) (int64, error) {
	v := d.q.Vertex
	return d.q.WithContext(ctx).Vertex.Where(
		v.Generation.Eq(generation),
		v.Latitude.Gte(latitude.Min),
		v.Latitude.Lt(latitude.Max),
		v.Longitude.Eq(longtitude),
	).Count()
}

func (d db) GetVertexCountOnHorizontalSegment(ctx context.Context, latitude float64, longtitude MinMax, generation Generation) (int64, error) {
	v := d.q.Vertex
	return d.q.WithContext(ctx).Vertex.Where(
		v.Generation.Eq(generation),
		v.Latitude.Eq(latitude),
		v.Longitude.Gte(longtitude.Min),
		v.Longitude.Lt(longtitude.Max),
	).Count()
}

func (d db) SetFlag(ctx context.Context, edgeIds []EdgeId, region RegionId, generation Generation) error {
	return d.q.Transaction(
		func(tx *query.Query) error {
			res, err := getArcFlags(ctx, tx, edgeIds, generation)
			if err != nil {
				return nil
			}
			afs := make([]*model.ArcFlag, 0, len(edgeIds))
			for _, af := range res {
				afs = append(afs, &model.ArcFlag{
					ID:         af.ID,
					EdgeId:     af.EdgeId,
					Flag:       ArcFlag(af.Flag).set(region),
					Generation: af.Generation,
				})
			}
			return tx.ArcFlag.WithContext(ctx).Create()
		})
}

func getArcFlags(ctx context.Context, q *query.Query, edgeIds []EdgeId, generation Generation) ([]*model.ArcFlag, error) {
	af := q.ArcFlag
	res, err := q.WithContext(ctx).ArcFlag.Where(af.Generation.Eq(generation), af.EdgeId.In(edgeIds...)).Find()
	if err != nil {
		return nil, err
	}
	if len(res) != len(edgeIds) {
		return nil, fmt.Errorf("not all arc flags are present")
	}
	return res, nil
}

func (d db) GetArcFlags(ctx context.Context, edgeIds []EdgeId, generation Generation) ([]ArcFlag, error) {
	res, err := getArcFlags(ctx, d.q, edgeIds, generation)
	ret := make([]ArcFlag, 0, len(res))
	for _, a := range res {
		ret = append(ret, ArcFlag(a.Flag))
	}
	return ret, err
}

func (d db) AddVertices(ctx context.Context, vertices []Vertex, generation Generation) error {
	vs := make([]*model.Vertex, 0, len(vertices))
	for _, v := range vertices {
		vs = append(vs, &model.Vertex{
			ID:         v.Id,
			Latitude:   v.Latitude,
			Longitude:  v.Longitude,
			Generation: generation,
		})
	}

	return d.q.WithContext(ctx).Vertex.Create(vs...)
}

func (d db) AddEdges(ctx context.Context, edges []Edge, generation Generation) error {
	es := make([]*model.Edge, 0, len(edges))
	for _, e := range edges {
		es = append(es, &model.Edge{
			ID:         e.Id,
			FromId:     e.From,
			ToId:       e.To,
			Generation: generation,
			Length:     e.Length,
		})
	}

	return d.q.WithContext(ctx).Edge.Create(es...)
}

func (d db) AddArcFlags(ctx context.Context, edgeIds []EdgeId, generation Generation) error {
	afs := make([]*model.ArcFlag, 0, len(edgeIds))
	for _, id := range edgeIds {
		afs = append(afs, &model.ArcFlag{
			EdgeId:     id,
			Generation: generation,
		})
	}

	return d.q.WithContext(ctx).ArcFlag.Create(afs...)
}

func Connect(con *gorm.DB) (DB, error) {
	if err := con.AutoMigrate(model.List...); err != nil {
		return nil, err
	}

	return db{q: query.Use(con)}, nil
}

func ConnectToDefault() (DB, error) {
	con, err := conn.Default()
	if err != nil {
		return nil, err
	}
	return Connect(con)
}
