package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDb(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	// Ping check for db
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
