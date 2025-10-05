package service

import (
	"context"
	"fmt"
	"handworks/common/grpc/inventory"
	types "handworks/common/types/inventory"

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
	name, itemType, unit, category, imageUrl string,
	quantity, maxQuantity int32,
) (*types.DbInventoryItem, error) {
	var item types.DbInventoryItem

	if err := tx.QueryRow(c,
		`INSERT INTO inventory.items
		 (name, type, unit, quantity, max_quantity, category ,image_url, is_available)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, true)
		 RETURNING id, name, type, status, unit, category, quantity, max_quantity, image_url, is_available, created_at, updated_at`,
		name, itemType, unit, quantity, maxQuantity, category, imageUrl,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Type,
		&item.Status,
		&item.Unit,
		&item.Category,
		&item.Quantity,
		&item.MaxQuantity,
		&item.ImageUrl,
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
func (i *InventoryService) FetchInventoryItems(
	ctx context.Context,
	tx pgx.Tx,
) ([]*types.DbInventoryItem, error) {
	rows, err := tx.Query(ctx, `
        SELECT id, name, type, status, unit, category, quantity, max_quantity, is_available, created_at, updated_at
        FROM inventory.items
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, fmt.Errorf("could not fetch inventory items: %w", err)
	}
	defer rows.Close()

	var items []*types.DbInventoryItem
	for rows.Next() {
		var item types.DbInventoryItem
		if err := rows.Scan(
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
			return nil, fmt.Errorf("could not scan inventory row: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}
func (i *InventoryService) FetchInventoryItemsByType(
	ctx context.Context,
	tx pgx.Tx,
	itemType string,
) ([]*types.DbInventoryItem, error) {
	rows, err := tx.Query(ctx, `
        SELECT id, name, type, status, unit, category, quantity, max_quantity, is_available, created_at, updated_at
        FROM inventory.items
		WHERE type = $1
        ORDER BY created_at DESC
    `, itemType)
	if err != nil {
		return nil, fmt.Errorf("could not fetch inventory items: %w", err)
	}
	defer rows.Close()

	var items []*types.DbInventoryItem
	for rows.Next() {
		var item types.DbInventoryItem
		if err := rows.Scan(
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
			return nil, fmt.Errorf("could not scan inventory row: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}
func (i *InventoryService) FetchInventoryItemsByCategory(
	ctx context.Context,
	tx pgx.Tx,
	category string,
) ([]*types.DbInventoryItem, error) {
	rows, err := tx.Query(ctx, `
        SELECT id, name, type, status, unit, category, quantity, max_quantity, is_available, created_at, updated_at
        FROM inventory.items
		WHERE category = $1
        ORDER BY created_at DESC
    `, category)
	if err != nil {
		return nil, fmt.Errorf("could not fetch inventory items: %w", err)
	}
	defer rows.Close()

	var items []*types.DbInventoryItem
	for rows.Next() {
		var item types.DbInventoryItem
		if err := rows.Scan(
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
			return nil, fmt.Errorf("could not scan inventory row: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}
func (i *InventoryService) FetchInventoryItemsByStatus(
	ctx context.Context,
	tx pgx.Tx,
	status string,
) ([]*types.DbInventoryItem, error) {
	rows, err := tx.Query(ctx, `
        SELECT id, name, type, status, unit, category, quantity, max_quantity, is_available, created_at, updated_at
        FROM inventory.items
		WHERE status = $1
        ORDER BY created_at DESC
    `, status)
	if err != nil {
		return nil, fmt.Errorf("could not fetch inventory items: %w", err)
	}
	defer rows.Close()

	var items []*types.DbInventoryItem
	for rows.Next() {
		var item types.DbInventoryItem
		if err := rows.Scan(
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
			return nil, fmt.Errorf("could not scan inventory row: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}
func (i *InventoryService) UpdateInventoryItem(
	ctx context.Context,
	tx pgx.Tx,
	in *inventory.UpdateItemRequest,
) (*types.DbInventoryItem, error) {
	args := pgx.NamedArgs{
		"id":           in.Id,
		"name":         in.Name,
		"type":         in.Type,
		"status":       in.Status,
		"category":     in.Category,
		"quantity":     in.Quantity,
		"max_quantity": in.MaxQuantity,
	}

	row := tx.QueryRow(ctx, `
		UPDATE inventory.items
		SET
			name = COALESCE(NULLIF(@name, ''), name),
			type = COALESCE(NULLIF(@type, ''), type),
			status = COALESCE(NULLIF(@status, ''), status),
			category = COALESCE(NULLIF(@category, ''), category),
			quantity = COALESCE(NULLIF(@quantity, '')::int, quantity),
			max_quantity = COALESCE(NULLIF(@max_quantity, '')::int, max_quantity),
			updated_at = NOW()
		WHERE id = @id
		RETURNING id, name, type, status, unit, category, quantity, max_quantity, is_available, created_at, updated_at
	`, args)

	var item types.DbInventoryItem
	if err := row.Scan(
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
		return nil, fmt.Errorf("could not update inventory item: %w", err)
	}

	return &item, nil
}
func (i *InventoryService) DeleteInventoryItem(
	ctx context.Context,
	tx pgx.Tx,
	id string,
) (*types.DbInventoryItem, error) {
	var item types.DbInventoryItem

	err := tx.QueryRow(ctx, `
		DELETE FROM inventory.items
		WHERE id = $1
		RETURNING id, name, type, status, unit, category, quantity, max_quantity, is_available, created_at, updated_at
	`, id).Scan(
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
	)

	if err != nil {
		return nil, fmt.Errorf("could not delete inventory item with id %s: %w", id, err)
	}

	return &item, nil
}
