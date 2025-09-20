package types

import "handworks/common/grpc/booking"

type CouchCleaningDetails struct {
	ID        string
	CouchType string
	WidthCM   int32
	DepthCM   int32
	HeightCM  int32
}

func (couchCleaning CouchCleaningDetails) ToProto() *booking.CouchCleaningDetails {
	return &booking.CouchCleaningDetails{
		CouchType: booking.CouchType(booking.CouchType_value[couchCleaning.CouchType]),
		WidthCm:   couchCleaning.WidthCM,
		DepthCm:   couchCleaning.DepthCM,
		HeightCm:  couchCleaning.HeightCM,
	}
}

func CouchCleaningDetailsFromProto(pb *booking.CouchCleaningDetails) CouchCleaningDetails {
	if pb == nil {
		return CouchCleaningDetails{}
	}
	return CouchCleaningDetails{
		CouchType: pb.CouchType.String(),
		WidthCM:   pb.WidthCm,
		DepthCM:   pb.DepthCm,
		HeightCM:  pb.HeightCm,
	}
}
