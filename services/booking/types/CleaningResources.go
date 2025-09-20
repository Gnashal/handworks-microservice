package types

import "handworks/common/grpc/booking"

type CleaningResources struct {
	ID       string
	Name     string
	Type     string
	PhotoURL string
}

func (cleaningResources CleaningResources) ToProto() *booking.CleaningResources {
	return &booking.CleaningResources{
		Id:       cleaningResources.ID,
		Name:     cleaningResources.Name,
		Type:     cleaningResources.Type,
		PhotoUrl: cleaningResources.PhotoURL,
	}
}

func CleaningResourceToProto(resources []CleaningResources) []*booking.CleaningResources {
	protos := make([]*booking.CleaningResources, 0, len(resources))
	for _, resource := range resources {
		protos = append(protos, resource.ToProto())
	}
	return protos
}
