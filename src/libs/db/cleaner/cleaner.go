package cleaner

import (
	"context"
	"fmt"

	"github.com/badochov/distributed-shortest-path/src/libs/db/conn"
	"github.com/badochov/distributed-shortest-path/src/libs/db/model"
	"github.com/badochov/distributed-shortest-path/src/libs/db/query"
	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Cleaner struct {
	q   *query.Query
	con *gorm.DB
}

func (c *Cleaner) truncateTable(ctx context.Context, table schema.Tabler) error {
	return c.con.WithContext(ctx).Exec(fmt.Sprintf(`TRUNCATE "%s" CASCADE;`, table.TableName())).Error
}

func (c *Cleaner) Clean(ctx context.Context) error {
	var err error

	if delErr := c.truncateTable(ctx, c.q.Generation); delErr != nil {
		err = multierror.Append(err, delErr)
	}
	if delErr := c.truncateTable(ctx, c.q.ArcFlag); delErr != nil {
		err = multierror.Append(err, delErr)
	}
	if delErr := c.truncateTable(ctx, c.q.Edge); delErr != nil {
		err = multierror.Append(err, delErr)
	}
	if delErr := c.truncateTable(ctx, c.q.RegionBinding); delErr != nil {
		err = multierror.Append(err, delErr)
	}
	if delErr := c.truncateTable(ctx, c.q.Vertex); delErr != nil {
		err = multierror.Append(err, delErr)
	}

	return err
}

func Connect(con *gorm.DB) (*Cleaner, error) {
	if err := con.AutoMigrate(model.List...); err != nil {
		return nil, err
	}
	return &Cleaner{q: query.Use(con), con: con}, nil
}

func ConnectToDefault() (*Cleaner, error) {
	con, err := conn.Default()
	if err != nil {
		return nil, err
	}
	return Connect(con)
}
