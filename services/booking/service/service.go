package service

import (
	"context"
	"handworks-services-booking/types"
	"handworks/common/grpc/booking"
	"handworks/common/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingService struct {
	L  *utils.Logger
	DB *pgxpool.Pool
	booking.UnimplementedBookingServiceServer
}

func (b *BookingService) CreateBook(c context.Context, in *booking.CreateBookRequest) (*booking.CreateBookResponse, error) {

	return nil, nil
}

func (b *BookingService) GetBooksByUID(ctx context.Context, in *booking.GetBooksByUIDRequest) (*booking.GetBooksByUIDResponse, error) {
	b.L.Info("Fetching Book by UID: %s...", in.CustomerId)

	if in.CustomerId == "" {
		return &booking.GetBooksByUIDResponse{
			Message: "customer_id cannot be empty",
		}, nil
	}

	var dbBooking []types.DbBookings

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		books, err := b.FetchBooksByUIDData(ctx, tx, in.CustomerId)
		if err != nil {
			return err
		}
		dbBooking = books
		return nil
	}); err != nil {
		return nil, err
	}

	protoBooks := make([]*booking.Booking, 0, len(dbBooking))
	for _, dbBook := range dbBooking {
		protoBooks = append(protoBooks, dbBook.ToProto())
	}

	return &booking.GetBooksByUIDResponse{
		Books:   protoBooks,
		Message: "success",
	}, nil
}

func (b *BookingService) UpdateBook(c context.Context, in *booking.UpdateBookRequest) (*booking.UpdateBookResponse, error) {

	return nil, nil
}

func (b *BookingService) DeleteBook(c context.Context, in *booking.DeleteBookRequest) (*booking.DeleteBookResponse, error) {

	return nil, nil
}
