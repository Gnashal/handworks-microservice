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
	var dbInv types.DbInventoryItem
	if err := i.withTx(ctx, i.DB, func(tx pgx.Tx) error {
		inv, err := i.FetchInventoryItem(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		dbInv = *inv
		return nil
	}); err != nil {
		i.L.Error("Failed fetching item: %s", err)
		return nil, err
	}
	return &inventory.GetItemResponse{
		Item: dbInv.ToProto(),
	}, nil
}
func (i *InventoryService) ListAllItems(ctx context.Context, in *inventory.Empty) (*inventory.ListItemsResponse, error) {
	var dbItems []*types.DbInventoryItem

	if err := i.withTx(ctx, i.DB, func(tx pgx.Tx) error {
		items, err := i.FetchInventoryItems(ctx, tx)
		if err != nil {
			return err
		}
		dbItems = items
		return nil
	}); err != nil {
		i.L.Error("Failed listing items: %s", err)
		return nil, err
	}
	protoItems := make([]*inventory.InventoryItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		protoItems = append(protoItems, dbItem.ToProto())
	}
	return &inventory.ListItemsResponse{
		Items: protoItems,
	}, nil
}
func (i *InventoryService) ListItemsByType(ctx context.Context, in *inventory.ListItemsByTypeRequest) (*inventory.ListItemsResponse, error) {
	var dbItems []*types.DbInventoryItem

	if err := i.withTx(ctx, i.DB, func(tx pgx.Tx) error {
		items, err := i.FetchInventoryItemsByType(ctx, tx, in.Type)
		if err != nil {
			return err
		}
		dbItems = items
		return nil
	}); err != nil {
		i.L.Error("Failed listing items: %s", err)
		return nil, err
	}
	protoItems := make([]*inventory.InventoryItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		protoItems = append(protoItems, dbItem.ToProto())
	}
	return &inventory.ListItemsResponse{
		Items: protoItems,
	}, nil
}
func (i *InventoryService) ListItemsByCategory(ctx context.Context, in *inventory.ListItemsByCategoryRequest) (*inventory.ListItemsResponse, error) {
	var dbItems []*types.DbInventoryItem

	if err := i.withTx(ctx, i.DB, func(tx pgx.Tx) error {
		items, err := i.FetchInventoryItemsByCategory(ctx, tx, in.Category)
		if err != nil {
			return err
		}
		dbItems = items
		return nil
	}); err != nil {
		i.L.Error("Failed listing items: %s", err)
		return nil, err
	}
	protoItems := make([]*inventory.InventoryItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		protoItems = append(protoItems, dbItem.ToProto())
	}
	return &inventory.ListItemsResponse{
		Items: protoItems,
	}, nil
}
func (i *InventoryService) ListItemsByStatus(ctx context.Context, in *inventory.ListItemsByStatusRequest) (*inventory.ListItemsResponse, error) {
	var dbItems []*types.DbInventoryItem

	if err := i.withTx(ctx, i.DB, func(tx pgx.Tx) error {
		items, err := i.FetchInventoryItemsByStatus(ctx, tx, in.Status)
		if err != nil {
			return err
		}
		dbItems = items
		return nil
	}); err != nil {
		i.L.Error("Failed listing items: %s", err)
		return nil, err
	}
	protoItems := make([]*inventory.InventoryItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		protoItems = append(protoItems, dbItem.ToProto())
	}
	return &inventory.ListItemsResponse{
		Items: protoItems,
	}, nil
}
func (i *InventoryService) UpdateItem(ctx context.Context, in *inventory.UpdateItemRequest) (*inventory.UpdateItemResponse, error) {
	var dbItem types.DbInventoryItem

	if err := i.withTx(ctx, i.DB, func(tx pgx.Tx) error {
		item, err := i.UpdateInventoryItem(ctx, tx, in)
		if err != nil {
			return err
		}
		dbItem = *item
		return nil
	}); err != nil {
		i.L.Error("Failed to update item: %s", err)
		return nil, err
	}
	return &inventory.UpdateItemResponse{
		Item: dbItem.ToProto(),
	}, nil
}
func (i *InventoryService) DeleteItem(ctx context.Context, in *inventory.DeleteItemRequest) (*inventory.DeleteItemResponse, error) {
	var dbInv types.DbInventoryItem
	if err := i.withTx(ctx, i.DB, func(tx pgx.Tx) error {
		inv, err := i.DeleteInventoryItem(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		dbInv = *inv
		return nil
	}); err != nil {
		i.L.Error("Failed deleting item: %s", err)
		return nil, err
	}

	return &inventory.DeleteItemResponse{
		Item: dbInv.ToProto(),
	}, nil
}
