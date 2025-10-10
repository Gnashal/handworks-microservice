package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"handworks/common/grpc/booking"
	types "handworks/common/types/booking"
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

func (b *BookingService) makeBaseBooking(
	c context.Context,
	tx pgx.Tx,
	custID string,
	customerFirstName string,
	customerLastName string,
	address types.Address,
	startSched time.Time,
	endSched time.Time,
	dirtyScale int32,
	paymentStatus string,
	reviewStatus string,
	photos []string,
	quoteId string,
) (*types.BaseBookingDetails, error) {

	var createdBaseBook types.BaseBookingDetails

	err := tx.QueryRow(c,
		`INSERT INTO booking.basebookings (
            cust_id,
            customer_first_name,
            customer_last_name,
            address,
            start_sched,
			end_sched,
            dirty_scale,
            payment_status,
            review_status,
            photos,
			created_at,
			updated_at,
			quote_id
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING id, cust_id, customer_first_name, customer_last_name, address, start_sched, end_sched, dirty_scale, payment_status, review_status, photos, created_at, updated_at, quote_id`,
		custID,
		customerFirstName,
		customerLastName,
		address,
		startSched,
		endSched,
		dirtyScale,
		paymentStatus,
		reviewStatus,
		photos,
		time.Now(),
		time.Now(),
		quoteId,
	).Scan(
		&createdBaseBook.ID,
		&createdBaseBook.CustID,
		&createdBaseBook.CustomerFirstName,
		&createdBaseBook.CustomerLastName,
		&createdBaseBook.Address,
		&createdBaseBook.StartSched,
		&createdBaseBook.EndSched,
		&createdBaseBook.DirtyScale,
		&createdBaseBook.PaymentStatus,
		&createdBaseBook.ReviewStatus,
		&createdBaseBook.Photos,
		&createdBaseBook.CreatedAt,
		&createdBaseBook.UpdatedAt,
		&createdBaseBook.QuoteId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert base booking: %w", err)
	}

	return &createdBaseBook, nil
}

func (b *BookingService) makeGeneralBooking(ctx context.Context, tx pgx.Tx, general *booking.GeneralCleaningDetails) (*types.ServiceDetails, error) {
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

func (b *BookingService) makeCouchBooking(ctx context.Context, tx pgx.Tx, couch *booking.CouchCleaningDetails) (*types.ServiceDetails, error) {
	specs := make([]types.CouchCleaningSpecifications, 0, len(couch.CleaningSpecs))
	for _, s := range couch.CleaningSpecs {
		specs = append(specs, types.CouchCleaningSpecifications{
			CouchType: s.CouchType.String(),
			WidthCM:   s.WidthCm,
			DepthCM:   s.DepthCm,
			HeightCM:  s.HeightCm,
			Quantity:  s.Quantity,
		})
	}

	couchDetails := types.CouchCleaningDetails{CleaningSpecs: specs}

	detailsJSON, err := json.Marshal(couchDetails)
	if err != nil {
		return nil, err
	}

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services (service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"COUCH", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	var couchOut types.CouchCleaningDetails
	if err := json.Unmarshal(rawDetails, &couchOut); err != nil {
		return nil, err
	}
	service.Details = couchOut

	return &service, nil
}
func (b *BookingService) makeMattressBooking(ctx context.Context, tx pgx.Tx, mattress *booking.MattressCleaningDetails) (*types.ServiceDetails, error) {
	specs := make([]types.MattressCleaningSpecifications, 0, len(mattress.CleaningSpecs))
	for _, s := range mattress.CleaningSpecs {
		specs = append(specs, types.MattressCleaningSpecifications{
			BedType:  s.BedType.String(),
			WidthCM:  s.WidthCm,
			DepthCM:  s.DepthCm,
			HeightCM: s.HeightCm,
			Quantity: s.Quantity,
		})
	}

	mattressDetails := types.MattressCleaningDetails{CleaningSpecs: specs}

	detailsJSON, err := json.Marshal(mattressDetails)
	if err != nil {
		return nil, err
	}

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services (service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"MATTRESS", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	var mattressOut types.MattressCleaningDetails
	if err := json.Unmarshal(rawDetails, &mattressOut); err != nil {
		return nil, err
	}
	service.Details = mattressOut

	return &service, nil
}

func (b *BookingService) makeCarBooking(ctx context.Context, tx pgx.Tx, car *booking.CarCleaningDetails) (*types.ServiceDetails, error) {
	specs := make([]types.CarCleaningSpecifications, 0, len(car.CleaningSpecs))
	for _, s := range car.CleaningSpecs {
		specs = append(specs, types.CarCleaningSpecifications{
			CarType:  s.CarType.String(),
			Quantity: s.Quantity,
		})
	}

	carDetails := types.CarCleaningDetails{
		CleaningSpecs: specs,
		ChildSeats:    car.ChildSeats,
	}

	detailsJSON, err := json.Marshal(carDetails)
	if err != nil {
		return nil, err
	}

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		INSERT INTO booking.services (service_type, details)
		VALUES ($1, $2)
		RETURNING id, service_type, details`,
		"CAR", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	var carOut types.CarCleaningDetails
	if err := json.Unmarshal(rawDetails, &carOut); err != nil {
		return nil, err
	}
	service.Details = carOut

	return &service, nil
}

func (b *BookingService) makePostConstructionBooking(ctx context.Context, tx pgx.Tx, postConstruction *booking.PostConstructionCleaningDetails) (*types.ServiceDetails, error) {
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
		return b.makeGeneralBooking(ctx, tx, details.General)
	case *booking.ServiceDetail_Couch:
		return b.makeCouchBooking(ctx, tx, details.Couch)
	case *booking.ServiceDetail_Mattress:
		return b.makeMattressBooking(ctx, tx, details.Mattress)
	case *booking.ServiceDetail_Car:
		return b.makeCarBooking(ctx, tx, details.Car)
	case *booking.ServiceDetail_Post:
		return b.makePostConstructionBooking(ctx, tx, details.Post)
	default:
		return nil, fmt.Errorf("unsupported main service type")
	}
}

func (b *BookingService) createAddOn(
	ctx context.Context,
	tx pgx.Tx,
	addon *booking.AddOnRequest,
	addOnPrice float32,
) (*types.AddOns, error) {
	addOnServiceDetails, err := b.createMainServiceBooking(ctx, tx, addon.ServiceDetail.Details)
	if err != nil {
		return nil, fmt.Errorf("failed to create service details: %w", err)
	}
	createdAddon := &types.AddOns{
		ServiceDetail: *addOnServiceDetails,
		Price:         addOnPrice,
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO booking.addons 
		(service_id, price)
		 VALUES ($1, $2)
		 RETURNING id, service_id, price`,
		addOnServiceDetails.ID,
		addOnPrice,
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

func (b *BookingService) saveBooking(
	ctx context.Context,
	tx pgx.Tx,
	baseBookingID, mainServiceID string,
	addonIDs, equipmentIDs, resourceIDs, cleanerIDs []string,
	totalPrice float32,
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

func (b *BookingService) fetchBookingsByUID(ctx context.Context, tx pgx.Tx, userID string) ([]*types.Booking, error) {
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
		bb.start_sched,
		bb.end_sched,
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
			&base.StartSched,
			&base.EndSched,
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

func (b *BookingService) FetchBookingsByID(ctx context.Context, tx pgx.Tx, bookID string) (*types.Booking, error) {
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
		bb.start_sched,
		bb.end_sched,
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
        WHERE b.id = $1;
    `

	rows, err := tx.Query(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var booking *types.Booking

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
			&base.StartSched,
			&base.EndSched,
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

		if booking == nil {
			details, err := types.UnmarshalServiceDetails(mainServiceType, mainRawDetails)
			if err != nil {
				return nil, err
			}

			booking = &types.Booking{
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
				booking.Equipments = append(booking.Equipments, types.CleaningEquipment{ID: id})
			}
			for _, id := range resourceIDs {
				booking.Resources = append(booking.Resources, types.CleaningResources{ID: id})
			}
			for _, id := range cleanerIDs {
				booking.Cleaners = append(booking.Cleaners, types.CleanerAssigned{ID: id})
			}
		}

		if addonID != "" {
			details, err := types.UnmarshalServiceDetails(addonServiceType, addonRawDetails)
			if err != nil {
				return nil, err
			}
			booking.Addons = append(booking.Addons, types.AddOns{
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

	if booking == nil {
		return nil, fmt.Errorf("booking %s not found", bookID)
	}

	return booking, nil
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

func (b *BookingService) updateBaseBooking(
	c context.Context,
	tx pgx.Tx,
	ID string,
	custID string,
	customerFirstName string,
	customerLastName string,
	address types.Address,
	startSched time.Time,
	endSched time.Time,
	dirtyScale int32,
	paymentStatus string,
	reviewStatus string,
	photos []string,
) (*types.BaseBookingDetails, error) {

	var updatedBaseBook types.BaseBookingDetails

	err := tx.QueryRow(c,
		`UPDATE booking.basebookings 
		SET cust_id = $2,
            customer_first_name = $3,
            customer_last_name = $4,
            address = $5,
            start_sched = $6,
			end_sched = $7,
            dirty_scale = $8,
            payment_status = $9,
            review_status = $10,
            photos = $11,
			updated_at = $12
        WHERE id = $1
        RETURNING id, cust_id, customer_first_name, customer_last_name, address, start_sched, end_sched, dirty_scale, payment_status, review_status, photos, created_at, updated_at`,
		ID,
		custID,
		customerFirstName,
		customerLastName,
		address,
		startSched,
		endSched,
		dirtyScale,
		paymentStatus,
		reviewStatus,
		photos,
		time.Now(),
	).Scan(
		&updatedBaseBook.ID,
		&updatedBaseBook.CustID,
		&updatedBaseBook.CustomerFirstName,
		&updatedBaseBook.CustomerLastName,
		&updatedBaseBook.Address,
		&updatedBaseBook.StartSched,
		&updatedBaseBook.EndSched,
		&updatedBaseBook.DirtyScale,
		&updatedBaseBook.PaymentStatus,
		&updatedBaseBook.ReviewStatus,
		&updatedBaseBook.Photos,
		&updatedBaseBook.CreatedAt,
		&updatedBaseBook.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update base booking: %w", err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return &updatedBaseBook, nil
}

func (b *BookingService) updateGeneralBooking(ctx context.Context, tx pgx.Tx, general *booking.GeneralCleaningDetails, serviceId string) (*types.ServiceDetails, error) {
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
		UPDATE booking.services 
		SET service_type = $2, details = $3
		WHERE id = $1
		RETURNING id, service_type, details`,
		serviceId, "GENERAL", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	generalOut, err := types.UnmarshalGeneralDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = generalOut

	return &service, nil
}

func (b *BookingService) updateCouchBooking(ctx context.Context, tx pgx.Tx, couch *booking.CouchCleaningDetails, serviceId string) (*types.ServiceDetails, error) {
	specs := make([]types.CouchCleaningSpecifications, 0, len(couch.CleaningSpecs))
	for _, s := range couch.CleaningSpecs {
		specs = append(specs, types.CouchCleaningSpecifications{
			CouchType: s.CouchType.String(),
			WidthCM:   s.WidthCm,
			DepthCM:   s.DepthCm,
			HeightCM:  s.HeightCm,
			Quantity:  s.Quantity,
		})
	}

	couchDetails := types.CouchCleaningDetails{CleaningSpecs: specs}

	detailsJSON, err := couchDetails.MarshalCouchDetails()
	if err != nil {
		return nil, err
	}

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		UPDATE booking.services 
		SET service_type = $2, details = $3
		WHERE id = $1
		RETURNING id, service_type, details`,
		serviceId, "COUCH", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	couchOut, err := types.UnmarshalCouchDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = couchOut

	return &service, nil
}

func (b *BookingService) updateMattressBooking(ctx context.Context, tx pgx.Tx, mattress *booking.MattressCleaningDetails, serviceId string) (*types.ServiceDetails, error) {
	specs := make([]types.MattressCleaningSpecifications, 0, len(mattress.CleaningSpecs))
	for _, s := range mattress.CleaningSpecs {
		specs = append(specs, types.MattressCleaningSpecifications{
			BedType:  s.BedType.String(),
			WidthCM:  s.WidthCm,
			DepthCM:  s.DepthCm,
			HeightCM: s.HeightCm,
			Quantity: s.Quantity,
		})
	}

	mattressDetails := types.MattressCleaningDetails{CleaningSpecs: specs}

	detailsJSON, err := mattressDetails.MarshalMattressDetails()
	if err != nil {
		return nil, err
	}

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
		UPDATE booking.services 
		SET service_type = $2, details = $3
		WHERE id = $1
		RETURNING id, service_type, details`,
		serviceId, "MATTRESS", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	mattressOut, err := types.UnmarshalMattressDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = mattressOut

	return &service, nil
}

func (b *BookingService) updateCarBooking(ctx context.Context, tx pgx.Tx, car *booking.CarCleaningDetails, serviceId string) (*types.ServiceDetails, error) {
	specs := make([]types.CarCleaningSpecifications, 0, len(car.CleaningSpecs))
	for _, s := range car.CleaningSpecs {
		specs = append(specs, types.CarCleaningSpecifications{
			CarType:  s.CarType.String(),
			Quantity: s.Quantity,
		})
	}

	carDetails := types.CarCleaningDetails{
		CleaningSpecs: specs,
		ChildSeats:    car.ChildSeats,
	}

	detailsJSON, err := carDetails.MarshalCarDetails()
	if err != nil {
		return nil, err
	}

	service := types.ServiceDetails{}
	var rawDetails []byte

	err = tx.QueryRow(ctx, `
        UPDATE booking.services 
        SET service_type = $2, details = $3
		WHERE id = $1
        RETURNING id, service_type, details`,
		serviceId, "CAR", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	carOut, err := types.UnmarshalCarDetails(rawDetails)
	if err != nil {
		return nil, err
	}

	service.Details = carOut

	return &service, nil
}

func (b *BookingService) updatePostConstructionBooking(ctx context.Context, tx pgx.Tx, postConstruction *booking.PostConstructionCleaningDetails, serviceId string) (*types.ServiceDetails, error) {
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
		UPDATE booking.services
		SET service_type = $2, details = $3
		WHERE id = $1
		RETURNING id, service_type, details`,
		serviceId, "POST", detailsJSON,
	).Scan(&service.ID, &service.ServiceType, &rawDetails)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	postOut, err := types.UnmarshalPostConstructionDetails(rawDetails)
	if err != nil {
		return nil, err
	}
	service.Details = postOut

	return &service, nil
}

func (b *BookingService) updateMainServiceBooking(
	ctx context.Context,
	tx pgx.Tx,
	mainService *booking.ServiceDetail,
	serviceId string,
) (*types.ServiceDetails, error) {
	switch details := mainService.Type.(type) {
	case *booking.ServiceDetail_General:
		return b.updateGeneralBooking(ctx, tx, details.General, serviceId)
	case *booking.ServiceDetail_Couch:
		return b.updateCouchBooking(ctx, tx, details.Couch, serviceId)
	case *booking.ServiceDetail_Mattress:
		return b.updateMattressBooking(ctx, tx, details.Mattress, serviceId)
	case *booking.ServiceDetail_Car:
		return b.updateCarBooking(ctx, tx, details.Car, serviceId)
	case *booking.ServiceDetail_Post:
		return b.updatePostConstructionBooking(ctx, tx, details.Post, serviceId)
	default:
		return nil, fmt.Errorf("unsupported main service type")
	}
}

func (b *BookingService) updateAddOn(
	ctx context.Context,
	tx pgx.Tx,
	addon *booking.AddOn,
	addonID string,
) (*types.AddOns, error) {
	addOnServiceDetails, err := b.updateMainServiceBooking(ctx, tx, addon.ServiceDetail.Details, addonID)
	if err != nil {
		return nil, fmt.Errorf("failed to update service details: %w", err)
	}

	updatedAddon := &types.AddOns{
		ServiceDetail: *addOnServiceDetails,
	}

	err = tx.QueryRow(ctx,
		`UPDATE booking.addons 
		 SET service_id = $2, price = $3
		 WHERE id = $1
		 RETURNING id, service_id, price`,
		addon.Id,
		addOnServiceDetails.ID,
		addon.Price,
	).Scan(
		&updatedAddon.ID,
		&updatedAddon.ServiceDetail.ID,
		&updatedAddon.Price,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update addon: %w", err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return updatedAddon, nil
}

func (b *BookingService) updateBookPrice(
	ctx context.Context,
	tx pgx.Tx,
	ID string,
	totalPrice float32,
) (float32, error) {
	var updatedTotalPrice float64
	query := `
		UPDATE booking.bookings 
		SET total_price = $2
		WHERE id = $1
		RETURNING total_price`

	err := tx.QueryRow(ctx, query,
		ID,
		totalPrice,
	).Scan(&updatedTotalPrice)
	if err != nil {
		return 0.0, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return 0.0, err
	}

	return totalPrice, nil
}
