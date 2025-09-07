package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDb(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := os.Getenv("DEV_DB_CONN")
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
