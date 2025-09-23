package types

import "handworks/common/grpc/booking"

type CarCleaningDetails struct {
	ID         string
	CarType    string
	ChildSeats int32
}

func (carCleaning CarCleaningDetails) ToProto() *booking.CarCleaningDetails {
	return &booking.CarCleaningDetails{
		Id:         carCleaning.ID,
		CarType:    booking.CarType(booking.CarType_value[carCleaning.CarType]),
		ChildSeats: carCleaning.ChildSeats,
	}
}

func CarCleaningDetailsFromProto(pb *booking.CarCleaningDetails) CarCleaningDetails {
	if pb == nil {
		return CarCleaningDetails{}
	}
	return CarCleaningDetails{
		ID:         pb.Id,
		CarType:    pb.CarType.String(),
		ChildSeats: pb.ChildSeats,
	}
}
