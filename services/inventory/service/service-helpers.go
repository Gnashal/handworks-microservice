package service

import (
	"context"
	"fmt"
	"handworks-services-inventory/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (i *InventoryService) withTx(
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
				i.L.Error("rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()
	return fn(tx)
}

func (i *InventoryService) CreateInventoryItem(
	c context.Context,
	tx pgx.Tx,
	name, itemType, unit, category string,
	quantity, maxQuantity int32,
) (*types.DbInventoryItem, error) {
	var item types.DbInventoryItem

	if err := tx.QueryRow(c,
		`INSERT INTO inventory.items
		 (name, type, unit, quantity, max_quantity, category ,is_available)
		 VALUES ($1, $2, $3, $4, $5, $6, true)
		 RETURNING id, name, type, status, unit, category, quantity, max_quantity, is_available, created_at, updated_at`,
		name, itemType, unit, quantity, maxQuantity, category,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Type,
		&item.Status,
		&item.Unit,
		&item.Category,
		&item.Quantity,
		&item.MaxQuantity,
		&item.IsAvailable,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("could not create inventory item: %w", err)
	}

	return &item, nil
}

func (i *InventoryService) FetchInventoryItem(
	c context.Context,
	tx pgx.Tx,
	id string,
) (*types.DbInventoryItem, error) {
	var item types.DbInventoryItem

	if err := tx.QueryRow(c,
		`SELECT id, name, type, status, unit, quantity, max_quantity, is_available, created_at, updated_at
		 FROM inventory.items
		 WHERE id = $1`,
		id,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Type,
		&item.Status,
		&item.Unit,
		&item.Quantity,
		&item.MaxQuantity,
		&item.IsAvailable,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("could not fetch inventory item with id %s: %w", id, err)
	}

	return &item, nil
}
