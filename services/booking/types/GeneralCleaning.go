package types

import (
	"encoding/json"
	"handworks/common/grpc/booking"
)

type GeneralCleaningDetails struct {
	HomeType string `json:"home_type"`
	SQM      int32  `json:"sqm"`
}

func (generalCleaning GeneralCleaningDetails) ToProto() *booking.GeneralCleaningDetails {
	return &booking.GeneralCleaningDetails{
		HomeType: booking.HomeType(booking.HomeType_value[generalCleaning.HomeType]),
		Sqm:      generalCleaning.SQM,
	}
}

func GeneralCleaningDetailsFromProto(pb *booking.GeneralCleaningDetails) GeneralCleaningDetails {
	if pb == nil {
		return GeneralCleaningDetails{}
	}
	return GeneralCleaningDetails{
		HomeType: pb.HomeType.String(),
		SQM:      pb.Sqm,
	}
}

func (generalCleaning *GeneralCleaningDetails) MarshalGeneralDetails() ([]byte, error) {
	general := GeneralCleaningDetails{
		HomeType: generalCleaning.HomeType,
		SQM:      generalCleaning.SQM,
	}
	return json.Marshal(general)
}

func UnmarshalGeneralDetails(detailsOut []byte) (*GeneralCleaningDetails, error) {
	var generalDetails GeneralCleaningDetails
	if err := json.Unmarshal(detailsOut, &generalDetails); err != nil {
		return nil, err
	}
	return &GeneralCleaningDetails{
		HomeType: generalDetails.HomeType,
		SQM:      generalDetails.SQM,
	}, nil
}
