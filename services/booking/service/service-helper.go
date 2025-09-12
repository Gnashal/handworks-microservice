package service

import (
	"context"
	"fmt"
	"handworks-services-booking/types"
	"time"

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

func (b *BookingService) MakeBooking(c context.Context, tx pgx.Tx, customerID string, addressID string, bookingType string, dirtyScale int32, schedule time.Time, status string, notes string) (*types.DbBookings, error) {
	var createdBook types.DbBookings
	err := tx.QueryRow(c,
		`INSERT INTO booking.bookings (customer_id, address_id, booking_type, dirty_scale, schedule, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, customer_id, address_id, booking_type, dirty_scale, schedule, status, notes, created_at, updated_at`,
		customerID, addressID, bookingType, dirtyScale, schedule, status, notes,
	).Scan(
		&createdBook.ID,
		&createdBook.CustomerID,
		&createdBook.AddressID,
		&createdBook.BookingType,
		&createdBook.DirtyScale,
		&createdBook.Schedule,
		&createdBook.Status,
		&createdBook.Notes,
		&createdBook.CreatedAt,
		&createdBook.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert created booking: %w", err)
	}

	return &createdBook, nil
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

func (b *BookingService) PatchBook(c context.Context, tx pgx.Tx, ID string, addressID string, bookingType string, dirtyScale int32, schedule time.Time, status string, notes string) (*types.DbBookings, error) {
	var updatedBook types.DbBookings

	err := tx.QueryRow(c,
		`UPDATE booking.bookings
		 SET address_id   = $2,
		     booking_type = $3,
		     dirty_scale  = $4,
		     schedule     = $5,
		     status       = $6,
		     notes        = $7,
		     updated_at   = NOW()
		 WHERE id = $1
		 RETURNING id, customer_id, address_id, booking_type, dirty_scale,
		           schedule, status, notes, created_at, updated_at`,
		ID, addressID, bookingType, dirtyScale, schedule, status, notes,
	).Scan(
		&updatedBook.ID,
		&updatedBook.CustomerID,
		&updatedBook.AddressID,
		&updatedBook.BookingType,
		&updatedBook.DirtyScale,
		&updatedBook.Schedule,
		&updatedBook.Status,
		&updatedBook.Notes,
		&updatedBook.CreatedAt,
		&updatedBook.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	return &updatedBook, nil
}

func (b *BookingService) RemoveBook(c context.Context, tx pgx.Tx, ID string) (*types.DbBookings, error) {
	var deletedBook types.DbBookings

	err := tx.QueryRow(c,
		`DELETE FROM booking.bookings
		 WHERE id = $1
		 RETURNING id, customer_id, address_id, booking_type, dirty_scale,
		           schedule, status, notes, created_at, updated_at`,
		ID,
	).Scan(
		&deletedBook.ID,
		&deletedBook.CustomerID,
		&deletedBook.AddressID,
		&deletedBook.BookingType,
		&deletedBook.DirtyScale,
		&deletedBook.Schedule,
		&deletedBook.Status,
		&deletedBook.Notes,
		&deletedBook.CreatedAt,
		&deletedBook.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to delete book: %w", err)
	}

	return &deletedBook, nil
}
