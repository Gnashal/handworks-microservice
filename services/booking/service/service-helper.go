package service

import (
	"context"
	"fmt"
	"handworks-services-booking/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (b *BookingService) withTx(
	ctx context.Context,
	db *pgxpool.Pool,
	fn func(pgx.Tx) error,
) (err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				b.L.Error("rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()
	return fn(tx)
}

func (b *BookingService) FetchBooksByUIDData(c context.Context, tx pgx.Tx, customerID string) ([]types.DbBookings, error) {
	rows, err := tx.Query(c, `
        SELECT id, customer_id, address_id, booking_type, dirty_scale,
               schedule, status, notes, created_at, updated_at
        FROM booking.bookings
        WHERE customer_id = $1
    `, customerID)
	if err != nil {
		return nil, fmt.Errorf("could not query bookings table: %w", err)
	}
	defer rows.Close()

	var books []types.DbBookings
	for rows.Next() {
		var book types.DbBookings
		if err := rows.Scan(
			&book.ID,
			&book.CustomerID,
			&book.AddressID,
			&book.BookingType,
			&book.DirtyScale,
			&book.Schedule,
			&book.Status,
			&book.Notes,
			&book.CreatedAt,
			&book.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("could not scan booking row: %w", err)
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
