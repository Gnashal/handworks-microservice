package tasks

import (
	"context"
	"fmt"

	"handworks-api/types"
)

type BookingService interface {
	CreateBooking(ctx context.Context, evt types.CreateBookingEvent) (*types.Booking, error)
	GetBookingById(ctx context.Context, id string) (*types.Booking, error)
	UpdateBooking(ctx context.Context, id string) error
}

type BookingTasks struct {
	Svc BookingService
}

func NewBookingTasks(svc BookingService) *BookingTasks {
	return &BookingTasks{Svc: svc}
}

func (t *BookingTasks) CreateBooking(ctx context.Context, evt types.CreateBookingEvent) (*types.Booking, error) {
	if t.Svc == nil {
		return nil, fmt.Errorf("booking service (CreateBooking) is not configured")
	}
	return t.Svc.CreateBooking(ctx, evt)
}

func (t *BookingTasks) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	if t.Svc == nil {
		return nil, fmt.Errorf("booking service (GetBookingById) is not configured")
	}
	return t.Svc.GetBookingById(ctx, id)
}

func (t *BookingTasks) UpdateBooking(ctx context.Context, id string) error {
	if t.Svc == nil {
		return fmt.Errorf("booking service (UpdateBooking) is not configured")
	}
	return t.Svc.UpdateBooking(ctx, id)
}
