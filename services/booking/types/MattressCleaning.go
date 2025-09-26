package types

import (
	"encoding/json"
	"handworks/common/grpc/booking"
)

type MattressCleaningDetails struct {
	BedType  string `json:"bed_type"`
	WidthCM  int32  `json:"width_cm"`
	DepthCM  int32  `json:"depth_cm"`
	HeightCM int32  `json:"height_cm"`
}

func (mattressCleaning MattressCleaningDetails) ToProto() *booking.MattressCleaningDetails {
	return &booking.MattressCleaningDetails{
		BedType:  booking.BedType(booking.CarType_value[mattressCleaning.BedType]),
		WidthCm:  mattressCleaning.WidthCM,
		DepthCm:  mattressCleaning.DepthCM,
		HeightCm: mattressCleaning.HeightCM,
	}
}

func MattressCleaningDetailsFromProto(pb *booking.MattressCleaningDetails) MattressCleaningDetails {
	if pb == nil {
		return MattressCleaningDetails{}
	}
	return MattressCleaningDetails{
		BedType:  pb.BedType.String(),
		WidthCM:  pb.WidthCm,
		DepthCM:  pb.DepthCm,
		HeightCM: pb.HeightCm,
	}
}

func (mattressCleaning *MattressCleaningDetails) MarshalMattressDetails() ([]byte, error) {
	mattress := MattressCleaningDetails{
		BedType:  mattressCleaning.BedType,
		WidthCM:  mattressCleaning.WidthCM,
		DepthCM:  mattressCleaning.DepthCM,
		HeightCM: mattressCleaning.HeightCM,
	}
	return json.Marshal(mattress)
}

func UnmarshalMattressDetails(detailsOut []byte) (*MattressCleaningDetails, error) {
	var mattressDetails MattressCleaningDetails
	if err := json.Unmarshal(detailsOut, &mattressDetails); err != nil {
		return nil, err
	}
	return &MattressCleaningDetails{
		BedType:  mattressDetails.BedType,
		WidthCM:  mattressDetails.WidthCM,
		DepthCM:  mattressDetails.DepthCM,
		HeightCM: mattressDetails.HeightCM,
	}, nil
}
