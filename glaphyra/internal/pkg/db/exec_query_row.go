package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func (db Database) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}
