package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (a *InventoryService) withTx(
	ctx context.Context,
	db *pgxpool.Pool,
	fn func(pgx.Tx) error,
) (err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				a.L.Error("rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()
	return fn(tx)
}
