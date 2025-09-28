package types

import (
	"encoding/json"
	"handworks/common/grpc/booking"
)

type PostConstructionDetails struct {
	SQM int32 `json:"sqm"`
}

func (postConstruction PostConstructionDetails) ToProto() *booking.PostConstructionCleaningDetails {
	return &booking.PostConstructionCleaningDetails{
		Sqm: postConstruction.SQM,
	}
}

func PostConstructionCleaningDetailsFromProto(pb *booking.PostConstructionCleaningDetails) PostConstructionDetails {
	if pb == nil {
		return PostConstructionDetails{}
	}
	return PostConstructionDetails{
		SQM: pb.Sqm,
	}
}

func (postConstructionCleaning *PostConstructionDetails) MarshalPostConstructionDetails() ([]byte, error) {
	postConstruction := PostConstructionDetails{
		SQM: postConstructionCleaning.SQM,
	}
	return json.Marshal(postConstruction)
}

func UnmarshalPostConstructionDetails(detailsOut []byte) (*PostConstructionDetails, error) {
	var postDetails PostConstructionDetails
	if err := json.Unmarshal(detailsOut, &postDetails); err != nil {
		return nil, err
	}
	return &postDetails, nil
}
