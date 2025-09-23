package types

import "handworks/common/grpc/booking"

type PostConstructionDetails struct {
	ID  string
	SQM int32
}

func (postConstruction PostConstructionDetails) ToProto() *booking.PostConstructionCleaningDetails {
	return &booking.PostConstructionCleaningDetails{
		Id:  postConstruction.ID,
		Sqm: postConstruction.SQM,
	}
}

func PostConstructionCleaningDetailsFromProto(pb *booking.PostConstructionCleaningDetails) PostConstructionDetails {
	if pb == nil {
		return PostConstructionDetails{}
	}
	return PostConstructionDetails{
		ID:  pb.Id,
		SQM: pb.Sqm,
	}
}
