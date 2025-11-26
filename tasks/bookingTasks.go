package tasks

import (
	"context"
	"fmt"

	"handworks-api/types"
)

type BookingService interface {
	CreateBooking(ctx context.Context, evt types.CreateBookingEvent) (*types.Booking, error)
}

type BookingTasks struct {
	Svc BookingService
}

func NewBookingTasks(svc BookingService) *BookingTasks {
	return &BookingTasks{Svc: svc}
}

func (t *BookingTasks) CreateBooking(ctx context.Context, evt types.CreateBookingEvent) (*types.Booking, error) {
	if t.Svc == nil {
		return nil, fmt.Errorf("booking service is not configured")
	}
	return t.Svc.CreateBooking(ctx, evt)
}
