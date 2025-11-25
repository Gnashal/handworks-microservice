package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *BookingService) withTx(
	ctx context.Context,
	fn func(pgx.Tx) error,
) (err error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.Logger.Error("rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()
	return fn(tx)
}
func (s *BookingService) CreateBooking(ctx context.Context) error {
	return nil
}

func (s *BookingService) GetBookingById(ctx context.Context) error {
	return nil
}

func (s *BookingService) GetBookingByUId(ctx context.Context) error {
	return nil
}

func (s *BookingService) UpdateBooking(ctx context.Context) error {
	return nil
}

func (s *BookingService) DeleteBooking(ctx context.Context) error {
	return nil
}