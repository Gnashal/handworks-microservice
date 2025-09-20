package types

import "handworks/common/grpc/booking"

type CleanerAssigned struct {
	ID               string
	CleanerFirstName string
	CleanerLastName  string
	PFPUrl           string
}

func (cleanerAssigned CleanerAssigned) ToProto() *booking.CleanerAssigned {
	return &booking.CleanerAssigned{
		Id:               cleanerAssigned.ID,
		CleanerFirstName: cleanerAssigned.CleanerFirstName,
		CleanerLastName:  cleanerAssigned.CleanerLastName,
		PfpUrl:           cleanerAssigned.PFPUrl,
	}
}

func CleanerAssignedToProto(cleanerAssignees []CleanerAssigned) []*booking.CleanerAssigned {
	protos := make([]*booking.CleanerAssigned, 0, len(cleanerAssignees))
	for _, cleanerAssigned := range cleanerAssignees {
		protos = append(protos, cleanerAssigned.ToProto())
	}
	return protos
}
