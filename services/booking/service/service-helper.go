package service

import (
	"context"
	"fmt"
	"handworks-services-booking/types"
	"handworks/common/grpc/booking"
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

func (b *BookingService) MakeBaseBooking(
	c context.Context,
	tx pgx.Tx,
	custID string,
	customerFirstName string,
	customerLastName string,
	address types.Address,
	schedule time.Time,
	dirtyScale int32,
	paymentStatus string,
	reviewStatus string,
	photos []string,
) (*types.BaseBookingDetails, error) {

	var createdBaseBook types.BaseBookingDetails

	err := tx.QueryRow(c,
		`INSERT INTO booking.basebookings (
            cust_id,
            customer_first_name,
            customer_last_name,
            address,
            schedule,
            dirty_scale,
            payment_status,
            review_status,
            photos
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, cust_id, customer_first_name, customer_last_name, address, schedule, dirty_scale, payment_status, review_status, photos, created_at, updated_at`,
		custID,
		customerFirstName,
		customerLastName,
		address,
		schedule,
		dirtyScale,
		paymentStatus,
		reviewStatus,
		photos,
	).Scan(
		&createdBaseBook.ID,
		&createdBaseBook.CustID,
		&createdBaseBook.CustomerFirstName,
		&createdBaseBook.CustomerLastName,
		&createdBaseBook.Address,
		&createdBaseBook.Schedule,
		&createdBaseBook.DirtyScale,
		&createdBaseBook.PaymentStatus,
		&createdBaseBook.ReviewStatus,
		&createdBaseBook.Photos,
		&createdBaseBook.CreatedAt,
		&createdBaseBook.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert base booking: %w", err)
	}

	return &createdBaseBook, nil
}

func (b *BookingService) MakeGeneralBooking(ctx context.Context, tx pgx.Tx, general *booking.GeneralCleaningDetails) (*booking.Services, error) {
	generalDetails := types.GeneralCleaningDetails{
		HomeType: general.HomeType.String(),
		SQM:      general.Sqm,
	}

	detailsJSON, err := generalDetails.MarshalGeneralDetails()
	if err != nil {
		return nil, err
	}

	var service types.ServiceDetails
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services 
		(service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"general", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	generalOut, err := types.UnmarshalGeneralDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = generalOut

	protoService := service.ToProto()

	return protoService, nil
}

func (b *BookingService) MakeCouchBooking(ctx context.Context, tx pgx.Tx, couch *booking.CouchCleaningDetails) (*booking.Services, error) {
	couchDetails := types.CouchCleaningDetails{
		CouchType: couch.CouchType.String(),
		WidthCM:   couch.WidthCm,
		DepthCM:   couch.DepthCm,
		HeightCM:  couch.HeightCm,
	}

	detailsJSON, err := couchDetails.MarshalCouchDetails()
	if err != nil {
		return nil, err
	}

	var service types.ServiceDetails
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services 
		(service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"couch", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	couchOut, err := types.UnmarshalCouchDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = couchOut

	protoService := service.ToProto()

	return protoService, nil
}

func (b *BookingService) MakeMattressBooking(ctx context.Context, tx pgx.Tx, mattress *booking.MattressCleaningDetails) (*booking.Services, error) {
	mattressDetails := types.MattressCleaningDetails{
		BedType:  mattress.BedType.String(),
		WidthCM:  mattress.WidthCm,
		DepthCM:  mattress.DepthCm,
		HeightCM: mattress.HeightCm,
	}

	detailsJSON, err := mattressDetails.MarshalMattressDetails()
	if err != nil {
		return nil, err
	}

	var service types.ServiceDetails
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services 
		(service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"mattress", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	mattressOut, err := types.UnmarshalMattressDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = mattressOut

	protoService := service.ToProto()
	return protoService, nil
}

func (b *BookingService) MakeCarBooking(ctx context.Context, tx pgx.Tx, car *booking.CarCleaningDetails) (*booking.Services, error) {
	carDetails := types.CarCleaningDetails{
		CarType:    car.CarType.String(),
		ChildSeats: car.ChildSeats,
	}

	detailsJSON, err := carDetails.MarshalCarDetails()
	if err != nil {
		return nil, err
	}

	var service types.ServiceDetails
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services 
		(service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"car", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	carOut, err := types.UnmarshalCarDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = carOut

	protoService := service.ToProto()
	return protoService, nil
}

func (b *BookingService) MakePostConstructionBooking(ctx context.Context, tx pgx.Tx, postConstruction *booking.PostConstructionCleaningDetails) (*booking.Services, error) {
	postDetails := types.PostConstructionDetails{
		SQM: postConstruction.Sqm,
	}

	detailsJSON, err := postDetails.MarshalPostConstructionDetails()
	if err != nil {
		return nil, err
	}

	var service types.ServiceDetails
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services 
		(service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"post", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	postOut, err := types.UnmarshalPostConstructionDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = postOut

	protoService := service.ToProto()
	return protoService, nil
}

func (b *BookingService) createMainServiceBooking(
	ctx context.Context,
	tx pgx.Tx,
	mainService *booking.ServiceDetail,
) (*booking.Services, error) {
	switch details := mainService.Type.(type) {
	case *booking.ServiceDetail_General:
		return b.MakeGeneralBooking(ctx, tx, details.General)
	case *booking.ServiceDetail_Couch:
		return b.MakeCouchBooking(ctx, tx, details.Couch)
	case *booking.ServiceDetail_Mattress:
		return b.MakeMattressBooking(ctx, tx, details.Mattress)
	case *booking.ServiceDetail_Car:
		return b.MakeCarBooking(ctx, tx, details.Car)
	case *booking.ServiceDetail_Post:
		return b.MakePostConstructionBooking(ctx, tx, details.Post)
	default:
		return nil, fmt.Errorf("unsupported main service type")
	}
}

func (b *BookingService) createAddOn(
	ctx context.Context,
	tx pgx.Tx,
	addon *booking.AddOn,
) (*types.AddOns, error) {
	addOnServiceDetails, err := b.createMainServiceBooking(ctx, tx, addon.ServiceDetail.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to create service details: %w", err)
	}

	var createdAddon types.AddOns
	err = tx.QueryRow(ctx,
		`INSERT INTO booking.addons 
		(service_id, price)
		 VALUES ($1, $2)
		 RETURNING id, service_id, price`,
		addOnServiceDetails.Id,
		addon.Price,
	).Scan(
		&createdAddon.ID,
		&createdAddon.ServiceDetail.ID,
		&createdAddon.Price,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert addon: %w", err)
	}

	return &createdAddon, nil
}

func (b *BookingService) createEquipment(
	// ctx context.Context,
	// tx pgx.Tx,
	equipment *booking.CleaningEquipment,
) (*types.CleaningEquipment, error) {
	var createdEquipment types.CleaningEquipment

	createdEquipment.ID = "d3b07384-4e6f-4f7e-bf6a-2e7f3c2c1f9a"
	createdEquipment.Name = equipment.Name
	createdEquipment.Type = equipment.Type
	createdEquipment.PhotoURL = equipment.PhotoUrl

	fmt.Print(createdEquipment)

	return &createdEquipment, nil
}

func (b *BookingService) createResource(
	// ctx context.Context,
	// tx pgx.Tx,
	resource *booking.CleaningResources,
) (*types.CleaningResources, error) {
	var createdResource types.CleaningResources

	createdResource.ID = "8f14e45f-ea9d-4c3d-9b1c-ff99a5b6d7c1"
	createdResource.Name = resource.Name
	createdResource.Type = resource.Type
	createdResource.PhotoURL = resource.PhotoUrl

	fmt.Print(createdResource)

	return &createdResource, nil
}

func (b *BookingService) createCleanersAssigned(
	// ctx context.Context,
	// tx pgx.Tx,
	cleaner *booking.CleanerAssigned,
) (*types.CleanerAssigned, error) {
	var createdCleaner types.CleanerAssigned

	createdCleaner.ID = "45c48cce-0f0e-4b8f-8c4e-9d1a7f2f5e2d"
	createdCleaner.CleanerFirstName = cleaner.CleanerFirstName
	createdCleaner.CleanerLastName = cleaner.CleanerLastName
	createdCleaner.PFPUrl = cleaner.PfpUrl

	fmt.Print(createdCleaner)

	return &createdCleaner, nil
}

func (b *BookingService) saveBooking(
	ctx context.Context,
	tx pgx.Tx,
	baseBookingID, mainServiceID string,
	addonIDs, equipmentIDs, resourceIDs, cleanerIDs []string,
	totalPrice float64,
) (string, error) {
	var id string
	query := `
		INSERT INTO booking.bookings 
		(base_booking_id, main_service_id, addon_ids, equipment_ids, resource_ids, cleaner_ids, total_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := tx.QueryRow(ctx, query,
		baseBookingID,
		mainServiceID,
		addonIDs,
		equipmentIDs,
		resourceIDs,
		cleanerIDs,
		totalPrice,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (b *BookingService) FetchBookingsByUID(
	ctx context.Context,
	tx pgx.Tx,
	userID string,
) (*types.Booking, error) {
	var booking types.Booking
	var addonIDs []string
	var equipmentIDs []string
	var resourceIDs []string
	var cleanerIDs []string

	query := `
		SELECT
			b.id AS booking_id,
			b.base_booking_id,
			b.total_price,
			b.addon_ids,
			b.equipment_ids,
			b.resource_ids,
			b.cleaner_ids,
			bb.id AS base_id,
			bb.cust_id,
			bb.customer_first_name,
			bb.customer_last_name,
			bb.address,
			bb.schedule,
			bb.dirty_scale,
			bb.payment_status,
			bb.review_status,
			bb.photos,
			bb.created_at,
			bb.updated_at
		FROM booking.bookings b
		JOIN booking.basebookings bb
		ON b.base_booking_id = bb.id
		WHERE bb.cust_id = $1;`

	err := tx.QueryRow(ctx, query,
		userID,
	).Scan(
		&booking.ID,
		&booking.Base.CustID,
		&booking.TotalPrice,
		&addonIDs,
		&equipmentIDs,
		&resourceIDs,
		&cleanerIDs,
		&booking.Base.ID,
		&booking.Base.CustID,
		&booking.Base.CustomerFirstName,
		&booking.Base.CustomerLastName,
		&booking.Base.Address,
		&booking.Base.Schedule,
		&booking.Base.DirtyScale,
		&booking.Base.PaymentStatus,
		&booking.Base.ReviewStatus,
		&booking.Base.Photos,
		&booking.Base.CreatedAt,
		&booking.Base.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	for _, id := range addonIDs {
		booking.Addons = append(booking.Addons, types.AddOns{ID: id})
	}

	for _, id := range equipmentIDs {
		booking.Equipments = append(booking.Equipments, types.CleaningEquipment{ID: id})
	}

	for _, id := range resourceIDs {
		booking.Resources = append(booking.Resources, types.CleaningResources{ID: id})
	}

	for _, id := range cleanerIDs {
		booking.Cleaners = append(booking.Cleaners, types.CleanerAssigned{ID: id})
	}

	return &booking, nil
}
