package types

import (
	"handworks/common/grpc/booking"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type DbBookings struct {
	ID          string
	CustomerID  string
	AddressID   string
	BookingType string
	DirtyScale  int32
	Schedule    time.Time
	Status      string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func (dbBooking DbBookings) ToProto() *booking.Booking {
	var updatedAt *timestamppb.Timestamp
	if dbBooking.UpdatedAt != nil {
		updatedAt = timestamppb.New(*dbBooking.UpdatedAt)
	}

	return &booking.Booking{
		Id:          dbBooking.ID,
		CustomerId:  dbBooking.CustomerID,
		AddressId:   dbBooking.AddressID,
		BookingType: booking.BookingType(booking.BookingType_value[dbBooking.BookingType]),
		DirtyScale:  dbBooking.DirtyScale,
		Schedule:    timestamppb.New(dbBooking.Schedule),
		Status:      booking.Status(booking.Status_value[dbBooking.Status]),
		Notes:       dbBooking.Notes,
		CreatedAt:   timestamppb.New(dbBooking.CreatedAt),
		UpdatedAt:   updatedAt,
	}
}
