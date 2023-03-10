// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:            db,
		ArcFlag:       newArcFlag(db, opts...),
		Edge:          newEdge(db, opts...),
		Generation:    newGeneration(db, opts...),
		RegionBinding: newRegionBinding(db, opts...),
		Vertex:        newVertex(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	ArcFlag       arcFlag
	Edge          edge
	Generation    generation
	RegionBinding regionBinding
	Vertex        vertex
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:            db,
		ArcFlag:       q.ArcFlag.clone(db),
		Edge:          q.Edge.clone(db),
		Generation:    q.Generation.clone(db),
		RegionBinding: q.RegionBinding.clone(db),
		Vertex:        q.Vertex.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.clone(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.clone(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:            db,
		ArcFlag:       q.ArcFlag.replaceDB(db),
		Edge:          q.Edge.replaceDB(db),
		Generation:    q.Generation.replaceDB(db),
		RegionBinding: q.RegionBinding.replaceDB(db),
		Vertex:        q.Vertex.replaceDB(db),
	}
}

type queryCtx struct {
	ArcFlag       *arcFlagDo
	Edge          *edgeDo
	Generation    *generationDo
	RegionBinding *regionBindingDo
	Vertex        *vertexDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		ArcFlag:       q.ArcFlag.WithContext(ctx),
		Edge:          q.Edge.WithContext(ctx),
		Generation:    q.Generation.WithContext(ctx),
		RegionBinding: q.RegionBinding.WithContext(ctx),
		Vertex:        q.Vertex.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	return &QueryTx{q.clone(q.db.Begin(opts...))}
}

type QueryTx struct{ *Query }

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
