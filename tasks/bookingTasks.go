// ...existing code...
package tasks

import (
	"context"
	"fmt"

	"handworks-api/services"
	"handworks-api/types"
)

type BookingTasks struct {
	Svc *services.BookingService
}

func NewBookingTasks(svc *services.BookingService) *BookingTasks {
	return &BookingTasks{Svc: svc}
}

func (t *BookingTasks) CreateBooking(ctx context.Context, evt types.CreateBookingEvent) (*types.Booking, error) {
	if t.Svc == nil {
		return nil, fmt.Errorf("booking service is not configured")
	}
	return t.Svc.CreateBooking(ctx, evt)
}
