package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

func (db Database) ExecQuery(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := db.cluster.Query(ctx, query, args...)
	noRows := true
	if err != nil {
		return fmt.Errorf("cannot exec query %s. Err: %w", query, err)
	}

	for rows.Next() {
		noRows = false
		err = pgxscan.ScanRow(dest, rows)
		if err != nil {
			return fmt.Errorf("cannot scan rows. Err: %w", err)
		}
	}
	rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}
	if noRows {
		return pgx.ErrNoRows
	}
	return nil
}
