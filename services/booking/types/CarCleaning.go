package types

import (
	"encoding/json"
	"handworks/common/grpc/booking"
)

type CarCleaningDetails struct {
	CarType    string `json:"car_type"`
	ChildSeats int32  `json:"child_seats"`
}

func (carCleaning CarCleaningDetails) ToProto() *booking.CarCleaningDetails {
	return &booking.CarCleaningDetails{
		CarType:    booking.CarType(booking.CarType_value[carCleaning.CarType]),
		ChildSeats: carCleaning.ChildSeats,
	}
}

func CarCleaningDetailsFromProto(pb *booking.CarCleaningDetails) CarCleaningDetails {
	if pb == nil {
		return CarCleaningDetails{}
	}
	return CarCleaningDetails{
		CarType:    pb.CarType.String(),
		ChildSeats: pb.ChildSeats,
	}
}

func (carCleaning *CarCleaningDetails) MarshalCarDetails() ([]byte, error) {
	car := CarCleaningDetails{
		CarType:    carCleaning.CarType,
		ChildSeats: carCleaning.ChildSeats,
	}
	return json.Marshal(car)
}

func UnmarshalCarDetails(detailsOut []byte) (*CarCleaningDetails, error) {
	var carDetails CarCleaningDetails
	if err := json.Unmarshal(detailsOut, &carDetails); err != nil {
		return nil, err
	}
	return &carDetails, nil
}
