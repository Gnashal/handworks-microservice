package service

import (
	"context"
	"database/sql"
	"errors"
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
            photos,
			created_at,
			updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
		time.Now(),
		time.Now(),
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

func (b *BookingService) MakeGeneralBooking(ctx context.Context, tx pgx.Tx, general *booking.GeneralCleaningDetails) (*types.ServiceDetails, error) {
	generalDetails := types.GeneralCleaningDetails{
		HomeType: general.HomeType.String(),
		SQM:      general.Sqm,
	}

	detailsJSON, err := generalDetails.MarshalGeneralDetails()
	if err != nil {
		return nil, err
	}

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services 
		(service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"GENERAL", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	generalOut, err := types.UnmarshalGeneralDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = generalOut

	return &service, nil
}

func (b *BookingService) MakeCouchBooking(ctx context.Context, tx pgx.Tx, couch *booking.CouchCleaningDetails) (*types.ServiceDetails, error) {
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

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services 
		(service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"COUCH", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	couchOut, err := types.UnmarshalCouchDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = couchOut

	return &service, nil
}

func (b *BookingService) MakeMattressBooking(ctx context.Context, tx pgx.Tx, mattress *booking.MattressCleaningDetails) (*types.ServiceDetails, error) {
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

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services 
		(service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"MATTRESS", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	mattressOut, err := types.UnmarshalMattressDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = mattressOut

	return &service, nil
}

func (b *BookingService) MakeCarBooking(ctx context.Context, tx pgx.Tx, car *booking.CarCleaningDetails) (*types.ServiceDetails, error) {
	carDetails := &types.CarCleaningDetails{
		CarType:    car.CarType.String(),
		ChildSeats: car.ChildSeats,
	}

	detailsJSON, err := carDetails.MarshalCarDetails()
	if err != nil {
		return nil, err
	}

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
        INSERT INTO booking.services 
        (service_type, details)
        VALUES ($1, $2)
        RETURNING id, service_type, details`,
		"CAR", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	carOut, err := types.UnmarshalCarDetails(rawDetails)
	if err != nil {
		return nil, err
	}

	service.Details = carOut

	return &service, nil
}

func (b *BookingService) MakePostConstructionBooking(ctx context.Context, tx pgx.Tx, postConstruction *booking.PostConstructionCleaningDetails) (*types.ServiceDetails, error) {
	postDetails := types.PostConstructionDetails{
		SQM: postConstruction.Sqm,
	}

	detailsJSON, err := postDetails.MarshalPostConstructionDetails()
	if err != nil {
		return nil, err
	}

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services 
		(service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"POST", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	postOut, err := types.UnmarshalPostConstructionDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = postOut

	return &service, nil
}

func (b *BookingService) createMainServiceBooking(
	ctx context.Context,
	tx pgx.Tx,
	mainService *booking.ServiceDetail,
) (*types.ServiceDetails, error) {
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

	createdAddon := &types.AddOns{
		ServiceDetail: *addOnServiceDetails,
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO booking.addons 
		(service_id, price)
		 VALUES ($1, $2)
		 RETURNING id, service_id, price`,
		addOnServiceDetails.ID,
		addon.Price,
	).Scan(
		&createdAddon.ID,
		&createdAddon.ServiceDetail.ID,
		&createdAddon.Price,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert addon: %w", err)
	}

	return createdAddon, nil
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

func (b *BookingService) FetchBookingsByUID(ctx context.Context, tx pgx.Tx, userID string) ([]*types.Booking, error) {
	query := `
        SELECT
		b.id                   AS booking_id,
		b.base_booking_id,
		b.main_service_id,
		b.addon_ids,
		b.equipment_ids,
		b.resource_ids,
		b.cleaner_ids,
		b.total_price,

		bb.id                  AS base_id,
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
		bb.updated_at,

		a.id                   AS addon_id,
		a.service_id           AS addon_service_id,
		a.price                AS addon_price,

		ms.service_type        AS main_service_type,
		ms.details             AS main_service_details,
		s.service_type         AS addon_service_type,
		s.details              AS addon_service_details

		FROM booking.bookings b
		JOIN booking.basebookings bb ON b.base_booking_id = bb.id
		JOIN booking.services ms      ON b.main_service_id = ms.id
		LEFT JOIN LATERAL unnest(b.addon_ids) AS addon_id(addon_id) ON true
		LEFT JOIN booking.addons a     ON addon_id.addon_id = a.id
		LEFT JOIN booking.services s   ON a.service_id = s.id
        WHERE bb.cust_id = $1;
    `

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookingsMap := make(map[string]*types.Booking)

	for rows.Next() {
		var bookingID, baseID, mainServiceID string
		var addonID, addonServiceID, mainServiceType, addonServiceType string
		var addonPrice sql.NullFloat64
		var mainRawDetails, addonRawDetails []byte
		var addonIDs, equipmentIDs, resourceIDs, cleanerIDs []string
		var base types.BaseBookingDetails
		var totalPrice float32

		if err := rows.Scan(
			&bookingID,
			&baseID,
			&mainServiceID,
			&addonIDs,
			&equipmentIDs,
			&resourceIDs,
			&cleanerIDs,
			&totalPrice,
			// Base fields
			&base.ID,
			&base.CustID,
			&base.CustomerFirstName,
			&base.CustomerLastName,
			&base.Address,
			&base.Schedule,
			&base.DirtyScale,
			&base.PaymentStatus,
			&base.ReviewStatus,
			&base.Photos,
			&base.CreatedAt,
			&base.UpdatedAt,
			// addon fields
			&addonID,
			&addonServiceID,
			&addonPrice,
			// service details
			&mainServiceType,
			&mainRawDetails,
			&addonServiceType,
			&addonRawDetails,
		); err != nil {
			return nil, err
		}

		bk, exists := bookingsMap[bookingID]
		if !exists {
			details, err := types.UnmarshalServiceDetails(mainServiceType, mainRawDetails)
			if err != nil {
				return nil, err
			}

			bk = &types.Booking{
				ID:   bookingID,
				Base: base,
				MainService: types.ServiceDetails{
					ID:          mainServiceID,
					ServiceType: mainServiceType,
					Details:     details,
				},
				Addons:     []types.AddOns{},
				Equipments: []types.CleaningEquipment{},
				Resources:  []types.CleaningResources{},
				Cleaners:   []types.CleanerAssigned{},
				TotalPrice: totalPrice,
			}

			for _, id := range equipmentIDs {
				bk.Equipments = append(bk.Equipments, types.CleaningEquipment{ID: id})
			}
			for _, id := range resourceIDs {
				bk.Resources = append(bk.Resources, types.CleaningResources{ID: id})
			}
			for _, id := range cleanerIDs {
				bk.Cleaners = append(bk.Cleaners, types.CleanerAssigned{ID: id})
			}

			bookingsMap[bookingID] = bk
		}

		if addonID != "" {
			details, err := types.UnmarshalServiceDetails(addonServiceType, addonRawDetails)
			if err != nil {
				return nil, err
			}
			bk.Addons = append(bk.Addons, types.AddOns{
				ID:    addonID,
				Price: float32(addonPrice.Float64),
				ServiceDetail: types.ServiceDetails{
					ID:          addonServiceID,
					ServiceType: addonServiceType,
					Details:     details,
				},
			})
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	bookings := make([]*types.Booking, 0, len(bookingsMap))
	for _, b := range bookingsMap {
		bookings = append(bookings, b)
	}

	return bookings, nil
}

func (b *BookingService) RemoveBooking(ctx context.Context, tx pgx.Tx, bookingID string) (bool, error) {
	var baseBookingID string
	var addonIDs []string
	var mainServiceID sql.NullString

	err := tx.QueryRow(ctx,
		`SELECT base_booking_id, addon_ids, main_service_id FROM booking.bookings WHERE id = $1`,
		bookingID,
	).Scan(&baseBookingID, &addonIDs, &mainServiceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("booking %s not found", bookingID)
		}
		return false, err
	}

	var addonServiceIDs []string
	if len(addonIDs) > 0 {
		err = tx.QueryRow(ctx,
			`SELECT COALESCE(array_agg(service_id), ARRAY[]::uuid[]) FROM booking.addons WHERE id = ANY($1::uuid[])`,
			addonIDs,
		).Scan(&addonServiceIDs)
		if err != nil {
			return false, err
		}
	}

	serviceMap := map[string]struct{}{}
	if mainServiceID.Valid {
		serviceMap[mainServiceID.String] = struct{}{}
	}
	for _, s := range addonServiceIDs {
		if s != "" {
			serviceMap[s] = struct{}{}
		}
	}
	var serviceIDs []string
	for k := range serviceMap {
		serviceIDs = append(serviceIDs, k)
	}

	if len(addonIDs) > 0 {
		if _, err = tx.Exec(ctx, `DELETE FROM booking.addons WHERE id = ANY($1::uuid[])`, addonIDs); err != nil {
			return false, err
		}
	}

	cmdB, err := tx.Exec(ctx, `DELETE FROM booking.bookings WHERE id = $1`, bookingID)
	if err != nil {
		return false, err
	}
	if cmdB.RowsAffected() == 0 {
		return false, fmt.Errorf("booking %s not found", bookingID)
	}

	cmdBB, err := tx.Exec(ctx, `DELETE FROM booking.basebookings WHERE id = $1`, baseBookingID)
	if err != nil {
		return false, err
	}
	if cmdBB.RowsAffected() == 0 {
		return false, fmt.Errorf("base booking %s not found", baseBookingID)
	}

	if len(serviceIDs) > 0 {
		if _, err = tx.Exec(ctx, `DELETE FROM booking.services WHERE id = ANY($1::uuid[])`, serviceIDs); err != nil {
			return false, err
		}
	}

	return true, nil
}
