package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDB(ctx context.Context, DSN string) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, DSN)
	if err != nil {
		return nil, fmt.Errorf("db conn error: %w", err)
	}
	return newDatabase(pool), nil
}
