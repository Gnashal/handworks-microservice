package types

import (
	"handworks/common/grpc/booking"
)

type AddOns struct {
	ID            string
	ServiceDetail ServiceDetails
	Price         float32
}

func (addOns AddOns) ToProto() *booking.AddOn {
	return &booking.AddOn{
		Id:            addOns.ID,
		ServiceDetail: addOns.ServiceDetail.ToProto(),
		Price:         addOns.Price,
	}
}

func AddOnsToProto(addOns []AddOns) []*booking.AddOn {
	protos := make([]*booking.AddOn, 0, len(addOns))
	for _, addOn := range addOns {
		protos = append(protos, addOn.ToProto())
	}
	return protos
}
