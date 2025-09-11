package service

import "handworks/common/grpc/booking"

var validBookingTypes = map[string]struct{}{
	booking.BookingType_CARPET.String():            {},
	booking.BookingType_CAR_INTERIOR.String():      {},
	booking.BookingType_COUCH.String():             {},
	booking.BookingType_GENERAL.String():           {},
	booking.BookingType_MATRESS.String():           {},
	booking.BookingType_POST_CONSTRUCTION.String(): {},
	booking.BookingType_SOFA.String():              {},
}

func DetermineBookingType(booking_type string) string {
	if _, ok := validBookingTypes[booking_type]; ok {
		return booking_type
	}
	return booking.BookingType_SOFA.String()
}
