package types

import "handworks/common/grpc/booking"

type Booking struct {
	ID          string
	Base        BaseBookingDetails
	MainService ServiceDetails
	Addons      []AddOns
	Equipments  []CleaningEquipment
	Resources   []CleaningResources
	Cleaners    []CleanerAssigned
	TotalPrice  float32
}

func (book Booking) ToProto() *booking.Booking {
	return &booking.Booking{
		Id:          book.ID,
		Base:        book.Base.ToProto(),
		MainService: book.MainService.ToProto(),
		Addons:      AddOnsToProto(book.Addons),
		Equipment:   CleaningEquipmentsToProto(book.Equipments),
		Resources:   CleaningResourceToProto(book.Resources),
		Cleaners:    CleanerAssignedToProto(book.Cleaners),
		TotalPrice:  book.TotalPrice,
	}
}
