package service

import (
	"context"
	"handworks/common/grpc/inventory"
	"handworks/common/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryService struct {
	L  *utils.Logger
	DB *pgxpool.Pool
	inventory.UnimplementedInventoryServiceServer
}

func (i *InventoryService) CreateItem(ctx context.Context, in *inventory.CreateItemRequest) (*inventory.CreateItemResponse, error) {
	// TODO: implement
	return nil, nil
}
func (i *InventoryService) GetItem(ctx context.Context, in *inventory.GetItemRequest) (*inventory.GetItemResponse, error) {
	return nil, nil
}
func (i *InventoryService) ListItems(ctx context.Context, in *inventory.ListItemsRequest) (*inventory.ListItemsResponse, error) {
	return nil, nil
}
func (i *InventoryService) UpdateItem(ctx context.Context, in *inventory.UpdateItemRequest) (*inventory.UpdateItemResponse, error) {
	return nil, nil
}
func (i *InventoryService) DeleteItem(ctx context.Context, in *inventory.DeleteItemRequest) (*inventory.DeleteItemResponse, error) {
	return nil, nil
}
