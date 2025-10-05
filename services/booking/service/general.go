package service

import (
	"handworks/common/grpc/booking"
)

var validHomeType = map[string]struct{}{
	booking.HomeType_CONDO_ROOM.String(): {},
	booking.HomeType_HOUSE.String():      {},
	booking.HomeType_COMMERCIAL.String(): {},
}

func DetermineHomeType(home_Type string) string {
	if _, ok := validHomeType[home_Type]; ok {
		return home_Type
	}
	return booking.HomeType_CONDO_ROOM.String()
}

var validCouchType = map[string]struct{}{
	booking.CouchType_SEATER_1.String():             {},
	booking.CouchType_SEATER_2.String():             {},
	booking.CouchType_SEATER_3.String():             {},
	booking.CouchType_SEATER_3_LTYPE_SMALL.String(): {},
	booking.CouchType_SEATER_3_LTYPE_LARGE.String(): {},
	booking.CouchType_SEATER_4_LTYPE_SMALL.String(): {},
	booking.CouchType_SEATER_4_LTYPE_LARGE.String(): {},
	booking.CouchType_SEATER_5_LTYPE.String():       {},
	booking.CouchType_SEATER_6_LTYPE.String():       {},
	booking.CouchType_OTTOMAN.String():              {},
	booking.CouchType_LAZBOY.String():               {},
	booking.CouchType_CHAIR.String():                {},
}

func DetermineCouchType(couch_Type string) string {
	if _, ok := validCouchType[couch_Type]; ok {
		return couch_Type
	}
	return booking.CouchType_SEATER_1.String()
}

var validBedType = map[string]struct{}{
	booking.BedType_KING.String():           {},
	booking.BedType_QUEEN.String():          {},
	booking.BedType_KING_HEADBAND.String():  {},
	booking.BedType_QUEEN_HEADBAND.String(): {},
	booking.BedType_DOUBLE.String():         {},
	booking.BedType_SINGLE.String():         {},
	booking.BedType_BED_PILLOW.String():     {},
}

func DetermineBedType(bed_Type string) string {
	if _, ok := validBedType[bed_Type]; ok {
		return bed_Type
	}
	return booking.BedType_KING.String()
}

var validCarType = map[string]struct{}{
	booking.CarType_SEDAN.String():     {},
	booking.CarType_MPV.String():       {},
	booking.CarType_SUV.String():       {},
	booking.CarType_PICKUP.String():    {},
	booking.CarType_VAN.String():       {},
	booking.CarType_CAR_SMALL.String(): {},
}

func DetermineCarType(car_Type string) string {
	if _, ok := validCarType[car_Type]; ok {
		return car_Type
	}
	return booking.CarType_SEDAN.String()
}

var validServiceType = map[string]struct{}{
	booking.MainServiceType_GENERAL_CLEANING.String(): {},
	booking.MainServiceType_COUCH.String():            {},
	booking.MainServiceType_CAR.String():              {},
	booking.MainServiceType_MATTRESS.String():         {},
	booking.MainServiceType_POST.String():             {},
}

func DetermineServiceType(service_Type string) string {
	if _, ok := validServiceType[service_Type]; ok {
		return service_Type
	}
	return booking.MainServiceType_SERVICE_TYPE_UNSPECIFIED.String()
}
