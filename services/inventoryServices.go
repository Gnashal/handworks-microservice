package services

import (
	"context"
	"handworks-api/types"
)

func (s *InventoryService) CreateItem(ctx context.Context, req types.CreateItemRequest) (types.InventoryItem, error) {
	return types.InventoryItem{}, nil
}

func (s *InventoryService) GetItem(ctx context.Context, id string) (types.InventoryItem, error) {
	return types.InventoryItem{}, nil
}
func (s *InventoryService) GetItems(ctx context.Context) ([]types.InventoryItem, error) {
	return nil, nil
}
func (s *InventoryService) ListItemsByType(ctx context.Context, itemType string) ([]types.InventoryItem, error) {
	return nil, nil
}

func (s *InventoryService) ListItemsByStatus(ctx context.Context, status string) ([]types.InventoryItem, error) {
	return nil, nil
}

func (s *InventoryService) ListItemsByCategory(ctx context.Context, category string) ([]types.InventoryItem, error) {
	return nil, nil
}

func (s *InventoryService) UpdateItem(ctx context.Context, req types.UpdateItemRequest) (types.InventoryItem, error) {
	return types.InventoryItem{}, nil
}

func (s *InventoryService) DeleteItem(ctx context.Context, id string) (types.InventoryItem, error) {
	return types.InventoryItem{}, nil
}
