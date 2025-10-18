package helpers

import (
	model "handworks-gateway/graph/generated/models"
	"handworks/common/grpc/inventory"
)

func MapInventoryItem(in *inventory.InventoryItem) *model.InventoryItem {
	if in == nil {
		return nil
	}
	return &model.InventoryItem{
		ID:          in.Id,
		Name:        in.Name,
		Type:        in.Type.String(),
		Status:      in.Status.String(),
		Category:    in.Status.String(),
		Quantity:    in.Quantity,
		MaxQuantity: &in.MaxQuantity,
		Unit:        &in.Unit,
		IsAvailable: in.IsAvailable,
		ImageURL:    &in.ImageUrl,
		CreatedAt:   timestampToTime(in.CreatedAt),
		UpdatedAt:   timestampToTime(in.UpdatedAt),
	}
}

func MapInventoryItems(itemsIn []*inventory.InventoryItem) []*model.InventoryItem {

	items := make([]*model.InventoryItem, 0, len(itemsIn))
	for _, item := range itemsIn {
		items = append(items, MapInventoryItem(item))
	}
	return items
}
