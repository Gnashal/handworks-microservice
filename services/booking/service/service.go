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
		mainService, err := b.createMainServiceBooking(ctx, tx, in.MainService)
		if err != nil {
			return err
		}

		baseBook, err := b.MakeBaseBooking(ctx, tx, in.Base.CustId, in.Base.CustomerFirstName, in.Base.CustomerLastName, types.AddressFromProto(in.Base.Address), in.Base.Schedule.AsTime(), in.Base.DirtyScale, in.Base.PaymentStatus, in.Base.ReviewStatus, in.Base.Photos)
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
			equipment, err := b.createEquipment(ctx, tx, equipment)
			if err != nil {
				return nil
			}
			equipments = append(equipments, *equipment)
			equipmentIDs = append(equipmentIDs, equipment.ID)
		}

		var resources []types.CleaningResources
		var resourceIDs []string
		for _, resource := range in.Resources {
			resource, err := b.createResource(ctx, tx, resource)
			if err != nil {
				return nil
			}
			resources = append(resources, *resource)
			resourceIDs = append(resourceIDs, resource.ID)
		}

		var cleanersAssigned []types.CleanerAssigned
		var cleanersAssignedIDs []string
		for _, cleanerAssigned := range in.Cleaners {
			cleanerAssigned, err := b.createCleanersAssigned(ctx, tx, cleanerAssigned)
			if err != nil {
				return nil
			}
			cleanersAssigned = append(cleanersAssigned, *cleanerAssigned)
			cleanersAssignedIDs = append(cleanersAssignedIDs, cleanerAssigned.ID)
		}

		totalPrice := float64(100.11)
		bookingID, err := b.saveBooking(ctx, tx, baseBook.ID, mainService.GetID(), addonIDs, equipmentIDs, resourceIDs, cleanersAssignedIDs, totalPrice)
		if err != nil {
			return err
		}

		createdBook = &types.Booking{
			ID:          bookingID,
			Base:        *baseBook,
			MainService: mainService,
			Addons:      addons,
			Equipment:   equipments,
			Resources:   resources,
			Cleaners:    cleanersAssigned,
			TotalPrice:  float32(totalPrice),
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
