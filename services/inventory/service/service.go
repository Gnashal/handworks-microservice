package service

import (
	"context"
	"handworks-services-inventory/types"
	"handworks/common/grpc/inventory"
	"handworks/common/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryService struct {
	L  *utils.Logger
	DB *pgxpool.Pool
	inventory.UnimplementedInventoryServiceServer
}

func (i *InventoryService) CreateItem(ctx context.Context, in *inventory.CreateItemRequest) (*inventory.CreateItemResponse, error) {
	var dbInv types.DbInventoryItem
	if err := i.withTx(ctx, i.DB, func(tx pgx.Tx) error {
		inv, err := i.CreateInventoryItem(ctx, tx, in.Name, in.Type, in.Unit, in.Category, in.Quantity, in.Quantity)
		if err != nil {
			return err
		}
		dbInv = *inv
		return nil
	}); err != nil {
		i.L.Error("Failed creating item: %s", err)
		return nil, err
	}
	return &inventory.CreateItemResponse{
		Item: dbInv.ToProto(),
	}, nil
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
