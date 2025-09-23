package types

import (
	"handworks/common/grpc/booking"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type BaseBookingDetails struct {
	ID                string
	CustID            string
	CustomerFirstName string
	CustomerLastName  string
	Address           Address
	Schedule          time.Time
	DirtyScale        int32
	PaymentStatus     string
	ReviewStatus      string
	Photos            []string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
}

func (baseBooking BaseBookingDetails) ToProto() *booking.BaseBookingDetails {
	var updatedAt *timestamppb.Timestamp
	if baseBooking.UpdatedAt != nil {
		updatedAt = timestamppb.New(*baseBooking.UpdatedAt)
	}

	return &booking.BaseBookingDetails{
		Id:                baseBooking.ID,
		CustId:            baseBooking.CustID,
		CustomerFirstName: baseBooking.CustomerFirstName,
		CustomerLastName:  baseBooking.CustomerLastName,
		Address:           baseBooking.Address.ToProto(),
		Schedule:          timestamppb.New(baseBooking.Schedule),
		DirtyScale:        baseBooking.DirtyScale,
		PaymentStatus:     baseBooking.PaymentStatus,
		ReviewStatus:      baseBooking.ReviewStatus,
		Photos:            baseBooking.Photos,
		CreatedAt:         timestamppb.New(baseBooking.CreatedAt),
		UpdatedAt:         updatedAt,
	}
}

func BaseBookingDetailsFromProto(pb *booking.BaseBookingDetails) BaseBookingDetails {
	if pb == nil {
		return BaseBookingDetails{}
	}

	var updatedAt *time.Time
	if pb.UpdatedAt != nil {
		t := pb.UpdatedAt.AsTime()
		updatedAt = &t
	}

	return BaseBookingDetails{
		ID:                pb.Id,
		CustID:            pb.CustId,
		CustomerFirstName: pb.CustomerFirstName,
		CustomerLastName:  pb.CustomerLastName,
		Address:           AddressFromProto(pb.Address),
		Schedule:          pb.Schedule.AsTime(),
		DirtyScale:        pb.DirtyScale,
		PaymentStatus:     pb.PaymentStatus,
		ReviewStatus:      pb.ReviewStatus,
		Photos:            pb.Photos,
		CreatedAt:         pb.CreatedAt.AsTime(),
		UpdatedAt:         updatedAt,
	}
}
