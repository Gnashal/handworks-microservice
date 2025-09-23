package types

import "handworks/common/grpc/booking"

type GeneralCleaningDetails struct {
	ID       string
	HomeType string
	SQM      int32
}

func (generalCleaning GeneralCleaningDetails) ToProto() *booking.GeneralCleaningDetails {
	return &booking.GeneralCleaningDetails{
		Id:       generalCleaning.ID,
		HomeType: booking.HomeType(booking.HomeType_value[generalCleaning.HomeType]),
		Sqm:      generalCleaning.SQM,
	}
}

func GeneralCleaningDetailsFromProto(pb *booking.GeneralCleaningDetails) GeneralCleaningDetails {
	if pb == nil {
		return GeneralCleaningDetails{}
	}
	return GeneralCleaningDetails{
		ID:       pb.Id,
		HomeType: pb.HomeType.String(),
		SQM:      pb.Sqm,
	}
}
