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

func (b *BookingService) CreateBook(ctx context.Context, in *booking.CreateBookRequest) (*booking.CreateBookResponse, error) {
	b.L.Info("Creating Book for UID: %s...", in.AddressId)

	var createdBook *types.DbBookings

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		book, err := b.MakeBooking(ctx, tx, in.CustomerId, in.AddressId, in.BookingType.String(), in.DirtyScale, in.Schedule.AsTime(), in.Status.String(), in.Notes)
		if err != nil {
			return err
		}
		createdBook = book
		return nil
	}); err != nil {
		return nil, err
	}
	book := createdBook.ToProto()

	return &booking.CreateBookResponse{
		Book: book,
	}, nil
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

func (b *BookingService) UpdateBook(ctx context.Context, in *booking.UpdateBookRequest) (*booking.UpdateBookResponse, error) {
	b.L.Info("Updating Book by ID: %s...", in.Id)

	var updatedBook *types.DbBookings

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		book, err := b.PatchBook(ctx, tx, in.Id, in.AddressId, in.BookingType.String(), in.DirtyScale, in.Schedule.AsTime(), in.Status.String(), in.Notes)
		if err != nil {
			return err
		}
		updatedBook = book
		return nil
	}); err != nil {
		return nil, err
	}
	book := updatedBook.ToProto()

	return &booking.UpdateBookResponse{
		Book: book,
	}, nil
}

func (b *BookingService) DeleteBook(ctx context.Context, in *booking.DeleteBookRequest) (*booking.DeleteBookResponse, error) {
	b.L.Info("Deleting Book by ID: %s...", in.Id)

	var deletedBook *types.DbBookings

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		book, err := b.RemoveBook(ctx, tx, in.Id)
		if err != nil {
			return err
		}
		deletedBook = book
		return nil
	}); err != nil {
		return nil, err
	}
	book := deletedBook.ToProto()

	return &booking.DeleteBookResponse{
		Id:   book.Id,
		Book: book,
	}, nil
}
