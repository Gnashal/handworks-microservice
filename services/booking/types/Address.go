package types

import "handworks/common/grpc/booking"

type Address struct {
	AddressHuman string
	AddressLat   float64
	AddressLng   float64
}

func (address Address) ToProto() *booking.Address {
	return &booking.Address{
		AddressHuman: address.AddressHuman,
		AddressLat:   address.AddressLat,
		AddressLng:   address.AddressLng,
	}
}

func AddressFromProto(pb *booking.Address) Address {
	if pb == nil {
		return Address{}
	}
	return Address{
		AddressHuman: pb.AddressHuman,
		AddressLat:   pb.AddressLat,
		AddressLng:   pb.AddressLng,
	}
}
