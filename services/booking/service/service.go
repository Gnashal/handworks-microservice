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

func (b *BookingService) CreateBooking(ctx context.Context, in *booking.CreateBookingRequest) (*booking.CreateBookingResponse, error) {
	b.L.Info("Creating Book for User: %s...", in.Base.CustomerFirstName)

	var createdBook *types.Booking

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		mainService, err := b.createMainServiceBooking(ctx, tx, in.MainService.Details)
		if err != nil {
			return err
		}

		baseBook, err := b.makeBaseBooking(ctx, tx, in.Base.CustId, in.Base.CustomerFirstName, in.Base.CustomerLastName, types.AddressFromProto(in.Base.Address), in.Base.Schedule.AsTime(), in.Base.DirtyScale, in.Base.PaymentStatus, in.Base.ReviewStatus, in.Base.Photos)
		if err != nil {
			return err
		}

		var addons []types.AddOns
		var addonIDs []string
		for _, addon := range in.Addons {
			addon, err := b.createAddOn(ctx, tx, addon)
			if err != nil {
				return err
			}
			addons = append(addons, *addon)
			addonIDs = append(addonIDs, addon.ID)
		}

		var equipments []types.CleaningEquipment
		var equipmentIDs []string
		for _, equipment := range in.Equipment {
			equipment, err := b.createEquipment( /*ctx, tx,*/ equipment)
			if err != nil {
				return nil
			}
			equipments = append(equipments, *equipment)
			equipmentIDs = append(equipmentIDs, equipment.ID)
		}

		var resources []types.CleaningResources
		var resourceIDs []string
		for _, resource := range in.Resources {
			resource, err := b.createResource( /*ctx, tx,*/ resource)
			if err != nil {
				return nil
			}
			resources = append(resources, *resource)
			resourceIDs = append(resourceIDs, resource.ID)
		}

		var cleanersAssigned []types.CleanerAssigned
		var cleanersAssignedIDs []string
		for _, cleanerAssigned := range in.Cleaners {
			cleanerAssigned, err := b.createCleanersAssigned( /*ctx, tx,*/ cleanerAssigned)
			if err != nil {
				return nil
			}
			cleanersAssigned = append(cleanersAssigned, *cleanerAssigned)
			cleanersAssignedIDs = append(cleanersAssignedIDs, cleanerAssigned.ID)
		}

		totalPrice := float32(100.11)
		bookingID, err := b.saveBooking(ctx, tx, baseBook.ID, mainService.ID, addonIDs, equipmentIDs, resourceIDs, cleanersAssignedIDs, totalPrice)
		if err != nil {
			return err
		}

		createdBook = &types.Booking{
			ID:          bookingID,
			Base:        *baseBook,
			MainService: *mainService,
			Addons:      addons,
			Equipments:  equipments,
			Resources:   resources,
			Cleaners:    cleanersAssigned,
			TotalPrice:  totalPrice,
		}

		return nil
	}); err != nil {
		return nil, err
	}
	book := createdBook.ToProto()

	return &booking.CreateBookingResponse{
		Book: book,
	}, nil
}

func (b *BookingService) GetBookingByUId(ctx context.Context, in *booking.GetBookingByUIdRequest) (*booking.GetBookingByUIdResponse, error) {
	b.L.Info("Getting Book for User: %s...", in.UserId)

	var book []*types.Booking

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		gotBook, err := b.fetchBookingsByUID(ctx, tx, in.UserId)
		if err != nil {
			return err
		}
		book = gotBook

		return nil
	}); err != nil {
		return nil, err
	}

	var protoBookings []*booking.Booking
	for _, b := range book {
		protoBookings = append(protoBookings, b.ToProto())
	}
	return &booking.GetBookingByUIdResponse{
		Booking: protoBookings,
	}, nil
}

func (b *BookingService) GetBookingById(ctx context.Context, in *booking.GetBookingByIdRequest) (*booking.GetBookingByIdResponse, error) {
	b.L.Info("Deleting Book from Booking with ID: %s", in.BookingId)

	var protoBooking *types.Booking

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		book, err := b.FetchBookingsByID(ctx, tx, in.BookingId)
		if err != nil {
			return err
		}
		protoBooking = book

		return nil
	}); err != nil {
		return nil, err
	}
	book := protoBooking.ToProto()

	return &booking.GetBookingByIdResponse{
		Booking: book,
	}, nil
}

func (b *BookingService) UpdateBooking(ctx context.Context, in *booking.Booking) (*booking.UpdateBookingResponse, error) {
	b.L.Info("Updating Book with ID: %s", in.Id)

	var updatedBook *types.Booking

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		mainService, err := b.updateMainServiceBooking(ctx, tx, in.MainService.Details, in.MainService.Id)
		if err != nil {
			return err
		}

		baseBook, err := b.updateBaseBooking(ctx, tx, in.Base.Id, in.Base.CustId, in.Base.CustomerFirstName, in.Base.CustomerLastName, types.AddressFromProto(in.Base.Address), in.Base.Schedule.AsTime(), in.Base.DirtyScale, in.Base.PaymentStatus, in.Base.ReviewStatus, in.Base.Photos)
		if err != nil {
			return err
		}

		var addons []types.AddOns
		for _, addon := range in.Addons {
			addon, err := b.updateAddOn(ctx, tx, addon, addon.ServiceDetail.Id)
			if err != nil {
				return err
			}
			addons = append(addons, *addon)
		}

		// To be improved when integrating
		var equipments []types.CleaningEquipment
		for _, equipment := range in.Equipment {
			equipment, err := b.createEquipment( /*ctx, tx,*/ equipment)
			if err != nil {
				return nil
			}
			equipments = append(equipments, *equipment)
		}

		// To be improved when integrating
		var resources []types.CleaningResources
		for _, resource := range in.Resources {
			resource, err := b.createResource( /*ctx, tx,*/ resource)
			if err != nil {
				return nil
			}
			resources = append(resources, *resource)
		}

		// To be improved when integrating
		var cleanersAssigned []types.CleanerAssigned
		for _, cleanerAssigned := range in.Cleaners {
			cleanerAssigned, err := b.createCleanersAssigned( /*ctx, tx,*/ cleanerAssigned)
			if err != nil {
				return nil
			}
			cleanersAssigned = append(cleanersAssigned, *cleanerAssigned)
		}

		totalPrice := float32(100.11) // replace with actual calculation logic
		updateTotalPrice, err := b.updateBookPrice(ctx, tx, in.Id, totalPrice)
		if err != nil {
			return err
		}

		updatedBook = &types.Booking{
			ID:          in.Id,
			Base:        *baseBook,
			MainService: *mainService,
			Addons:      addons,
			Equipments:  equipments,
			Resources:   resources,
			Cleaners:    cleanersAssigned,
			TotalPrice:  updateTotalPrice,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &booking.UpdateBookingResponse{
		Booking: updatedBook.ToProto(),
	}, nil
}

func (b *BookingService) DeleteBooking(ctx context.Context, in *booking.DeleteBookingByIDRequest) (*booking.DeleteBookingByIDResponse, error) {
	b.L.Info("Deleting Book with ID: %s", in.BookingId)

	var success bool

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		isDeleted, err := b.RemoveBooking(ctx, tx, in.BookingId)
		if err != nil {
			return err
		}
		success = isDeleted

		return nil
	}); err != nil {
		return nil, err
	}

	return &booking.DeleteBookingByIDResponse{
		Success: success,
	}, nil
}
