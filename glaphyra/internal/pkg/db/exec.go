package db

import (
	"context"

	"github.com/jackc/pgconn"
)

func (db Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}
