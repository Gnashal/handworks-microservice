package types

import "handworks/common/grpc/booking"

type MattressCleaningDetails struct {
	ID       string
	BedType  string
	WidthCM  int32
	DepthCM  int32
	HeightCM int32
}

func (mattressCleaning MattressCleaningDetails) ToProto() *booking.MattressCleaningDetails {
	return &booking.MattressCleaningDetails{
		Id:       mattressCleaning.ID,
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
		ID:       pb.Id,
		BedType:  pb.BedType.String(),
		WidthCM:  pb.WidthCm,
		DepthCM:  pb.DepthCm,
		HeightCM: pb.HeightCm,
	}
}
