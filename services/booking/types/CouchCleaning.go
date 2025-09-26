package types

import (
	"encoding/json"
	"handworks/common/grpc/booking"
)

type CouchCleaningDetails struct {
	CouchType string `json:"couch_type"`
	WidthCM   int32  `json:"width_cm"`
	DepthCM   int32  `json:"depth_cm"`
	HeightCM  int32  `json:"height_cm"`
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

func (couchCleaning *CouchCleaningDetails) MarshalCouchDetails() ([]byte, error) {
	couch := CouchCleaningDetails{
		CouchType: couchCleaning.CouchType,
		WidthCM:   couchCleaning.WidthCM,
		DepthCM:   couchCleaning.DepthCM,
		HeightCM:  couchCleaning.HeightCM,
	}
	return json.Marshal(couch)
}

func UnmarshalCouchDetails(detailsOut []byte) (*CouchCleaningDetails, error) {
	var couchDetails CouchCleaningDetails
	if err := json.Unmarshal(detailsOut, &couchDetails); err != nil {
		return nil, err
	}
	return &CouchCleaningDetails{
		CouchType: couchDetails.CouchType,
		WidthCM:   couchDetails.WidthCM,
		DepthCM:   couchDetails.DepthCM,
		HeightCM:  couchDetails.HeightCM,
	}, nil
}
