package types

import (
	"handworks/common/grpc/inventory"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type DbInventoryItem struct {
	ID          string
	Name        string
	Type        string // stored as string in DB (maps to enum ItemType)
	Status      string // stored as string in DB (maps to enum ItemStatus)
	Category    string // stored as string in DB (maps to enum ItemCategory)
	Quantity    int32
	Unit        string
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (dbItem DbInventoryItem) ToProto() *inventory.InventoryItem {
	return &inventory.InventoryItem{
		Id:          dbItem.ID,
		Name:        dbItem.Name,
		Type:        inventory.ItemType(inventory.ItemType_value[dbItem.Type]),
		Status:      inventory.ItemStatus(inventory.ItemStatus_value[dbItem.Status]),
		Category:    inventory.ItemCategory(inventory.ItemCategory_value[dbItem.Category]),
		Quantity:    dbItem.Quantity,
		Unit:        dbItem.Unit,
		IsAvailable: dbItem.IsAvailable,
		CreatedAt:   timestamppb.New(dbItem.CreatedAt),
		UpdatedAt:   timestamppb.New(dbItem.UpdatedAt),
	}
}
