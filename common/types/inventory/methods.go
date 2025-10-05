package types

import (
	"handworks/common/grpc/inventory"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (dbItem DbInventoryItem) ToProto() *inventory.InventoryItem {
	return &inventory.InventoryItem{
		Id:          dbItem.ID,
		Name:        dbItem.Name,
		Type:        inventory.ItemType(inventory.ItemType_value[dbItem.Type]),
		Status:      inventory.ItemStatus(inventory.ItemStatus_value[dbItem.Status]),
		Category:    inventory.ItemCategory(inventory.ItemCategory_value[dbItem.Category]),
		Quantity:    dbItem.Quantity,
		MaxQuantity: dbItem.MaxQuantity,
		ImageUrl:    dbItem.ImageUrl,
		Unit:        dbItem.Unit,
		IsAvailable: dbItem.IsAvailable,
		CreatedAt:   timestamppb.New(dbItem.CreatedAt),
		UpdatedAt:   timestamppb.New(dbItem.UpdatedAt),
	}
}
