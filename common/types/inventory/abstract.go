package types

import "time"

type DbInventoryItem struct {
	ID          string
	Name        string
	Type        string // stored as string in DB (maps to enum ItemType)
	Status      string // stored as string in DB (maps to enum ItemStatus)
	Category    string // stored as string in DB (maps to enum ItemCategory)
	Quantity    int32
	MaxQuantity int32
	ImageUrl    string
	Unit        string
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type CleaningEquipment struct {
	ID       string
	Name     string
	Type     string
	PhotoURL string
}
type CleaningResources struct {
	ID       string
	Name     string
	Type     string
	PhotoURL string
}

type BookingReply struct {
	Source     string              `json:"source"`
	Equipments []CleaningEquipment `json:"equipments,omitempty"`
	Resources  []CleaningResources `json:"resources,omitempty"`
	Error      string              `json:"error,omitempty"`
}
