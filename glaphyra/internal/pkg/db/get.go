package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
)

func (db Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.cluster, dest, query, args...)
}
