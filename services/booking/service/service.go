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

	// var createdBook *types.Booking
	var bookingID string

	if err := b.withTx(ctx, b.DB, func(tx pgx.Tx) error {
		MainServiceID, err := b.createMainServiceBooking(ctx, tx, in.MainService)
		if err != nil {
			return err
		}

		baseBookID, err := b.MakeBaseBooking(ctx, tx, in.Base.CustId, in.Base.CustomerFirstName, in.Base.CustomerLastName, types.AddressFromProto(in.Base.Address), in.Base.Schedule.AsTime(), in.Base.DirtyScale, in.Base.PaymentStatus, in.Base.ReviewStatus, in.Base.Photos)
		if err != nil {
			return err
		}

		var addonIDs []string
		for _, addon := range in.Addons {
			addonID, err := b.createAddOn(ctx, tx, addon)
			if err != nil {
				return err
			}
			addonIDs = append(addonIDs, addonID)
		}

		var equipmentIDs []string
		for _, equipment := range in.Equipment {
			equipmentID, err := b.createEquipment(ctx, tx, equipment)
			if err != nil {
				return nil
			}
			equipmentIDs = append(equipmentIDs, equipmentID)
		}

		var resourceIDs []string
		for _, resource := range in.Resources {
			resourceID, err := b.createResource(ctx, tx, resource)
			if err != nil {
				return nil
			}
			resourceIDs = append(resourceIDs, resourceID)
		}

		var cleanersAssignedIDs []string
		for _, cleanersAssigned := range in.Cleaners {
			cleanerAssignedID, err := b.createCleanersAssigned(ctx, tx, cleanersAssigned)
			if err != nil {
				return nil
			}
			cleanersAssignedIDs = append(cleanersAssignedIDs, cleanerAssignedID)
		}

		totalPrice := float64(100.11) // implement calculate price when ready
		bookingID, err = b.saveBooking(ctx, tx, baseBookID, MainServiceID, addonIDs, equipmentIDs, resourceIDs, cleanersAssignedIDs, totalPrice)
		if err != nil {
			return err
		}

		// createdBook = &types.Booking{
		// 	ID:          bookingID,
		// 	Base:        types.BaseBookingDetails{ID: baseBookID},
		// 	MainService: types.ServiceDetailFromProto(in.MainService),
		// 	Addons:      make([]types.AddOns, len(addonIDs)),
		// 	Equipment:   make([]types.CleaningEquipment, len(equipmentIDs)),
		// 	Resources:   make([]types.CleaningResources, len(resourceIDs)),
		// 	Cleaners:    make([]types.CleanerAssigned, len(cleanersAssignedIDs)),
		// 	TotalPrice:  float32(totalPrice),
		// }

		return nil
	}); err != nil {
		return nil, err
	}
	// book := createdBook.ToProto()

	return &booking.CreateBookingResponse{
		BookingId: bookingID,
	}, nil
}
