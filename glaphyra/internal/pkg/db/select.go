package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
)

func (db Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.cluster, dest, query, args...)
}
