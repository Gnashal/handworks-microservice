package services

import (
	"context"
	"fmt"
	"handworks-api/types"

	"github.com/jackc/pgx/v5"
)
func (s *InventoryService) withTx(
	ctx context.Context,
	fn func(pgx.Tx) error,
) (err error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.Logger.Error("rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()
	return fn(tx)
}
func (s *InventoryService) CreateItem(ctx context.Context, req types.CreateItemRequest) (*types.InventoryItem, error) {
	var item types.InventoryItem
	if err := s.withTx(ctx, func (tx pgx.Tx) error {
		inv, err := s.Tasks.CreateInventoryItem(ctx, tx, req.Name, req.Type, req.Unit, req.Category, req.ImageURL, req.Quantity, req.Quantity)
		if err != nil {
			return err
		}
		item = *inv
		return nil
	}); err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *InventoryService) GetItem(ctx context.Context, id string) (*types.InventoryItem, error) {
	var item types.InventoryItem
	if err := s.withTx(ctx, func (tx pgx.Tx) error {
		inv, err := s.Tasks.FetchInventoryItem(ctx, tx, id)
		if err != nil {
			return err
		}
		item = *inv
		return nil
	}); err != nil {
		return nil, err
	}
	return &item, nil
}
func (s *InventoryService) GetItems(ctx context.Context) ([]*types.InventoryItem, error) {
	var items []*types.InventoryItem
	if err := s.withTx(ctx, func (tx pgx.Tx) error {
		invs, err := s.Tasks.FetchFilter(ctx, tx, &types.InventoryFilter{})
		if err != nil {
			return err
		}
		items = invs
		return nil
	}); err != nil {
		return nil, err
	}
	return items, nil
}
func (s *InventoryService) ListItemsByType(ctx context.Context, itemType string) ([]*types.InventoryItem, error) {
	var items []*types.InventoryItem
	if err := s.withTx(ctx, func (tx pgx.Tx) error {
		invs, err := s.Tasks.FetchFilter(ctx, tx, &types.InventoryFilter{Type: &itemType})
		if err != nil {
			return err
		}
		items = invs
		return nil
	}); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *InventoryService) ListItemsByStatus(ctx context.Context, status string) ([]*types.InventoryItem, error) {
	var items []*types.InventoryItem
	if err := s.withTx(ctx, func (tx pgx.Tx) error {
		invs, err := s.Tasks.FetchFilter(ctx, tx, &types.InventoryFilter{Status: &status})
		if err != nil {
			return err
		}
		items = invs
		return nil
	}); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *InventoryService) ListItemsByCategory(ctx context.Context, category string) ([]*types.InventoryItem, error) {
	var items []*types.InventoryItem
	if err := s.withTx(ctx, func (tx pgx.Tx) error {
		invs, err := s.Tasks.FetchFilter(ctx, tx, &types.InventoryFilter{Category: &category})
		if err != nil {
			return err
		}
		items = invs
		return nil
	}); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *InventoryService) UpdateItem(ctx context.Context, req types.UpdateItemRequest) (*types.InventoryItem, error) {
	var item types.InventoryItem
	if err := s.withTx(ctx, func (tx pgx.Tx) error {
		inv, err := s.Tasks.UpdateInventoryItem(ctx, tx, &req)
		if err != nil {
			return err
		}
		item = *inv
		return nil
	}); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *InventoryService) DeleteItem(ctx context.Context, id string) (*types.InventoryItem, error) {
	var item types.InventoryItem
	if err := s.withTx(ctx, func (tx pgx.Tx) error {
		inv, err := s.Tasks.DeleteInventoryItem(ctx, tx, id)
		if err != nil {
			return err
		}
		item = *inv
		return nil
	}); err != nil {
		return nil, err
	}
	return &item, nil
}
