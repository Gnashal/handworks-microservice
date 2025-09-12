package service

import (
	"handworks/common/grpc/booking"
)

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

var validStatus = map[string]struct{}{
	booking.Status_PENDING.String():   {},
	booking.Status_CONFIRMED.String(): {},
	booking.Status_CANCELLED.String(): {},
}

func DetermineStatusType(status_type string) string {
	if _, ok := validStatus[status_type]; ok {
		return status_type
	}
	return booking.Status_CANCELLED.String()
}
