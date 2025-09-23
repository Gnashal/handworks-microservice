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

func (b *BookingService) MakeGeneralBooking(c context.Context, tx pgx.Tx, general *booking.GeneralCleaningDetails) (*types.GeneralCleaningDetails, error) {
	query := `
		INSERT INTO booking.general 
		(home_type, sqm)
        VALUES ($1, $2)
        RETURNING id, home_type, sqm`

	var generalBook types.GeneralCleaningDetails

	if err := tx.QueryRow(c, query, general.HomeType, general.Sqm).
		Scan(&generalBook.ID, &generalBook.HomeType, &generalBook.SQM); err != nil {
		return nil, err
	}

	return &generalBook, nil
}

func (b *BookingService) MakeCouchBooking(c context.Context, tx pgx.Tx, couch *booking.CouchCleaningDetails) (*types.CouchCleaningDetails, error) {
	query := `
		INSERT INTO booking.couch 
		(couch_type, width_cm, depth_cm, height_cm)
        VALUES ($1, $2, $3, $4)
        RETURNING id, couch_type, width_cm, depth_cm, height_cm`

	var couchBook types.CouchCleaningDetails

	if err := tx.QueryRow(c, query, couch.CouchType, couch.WidthCm, couch.DepthCm, couch.HeightCm).
		Scan(&couchBook.ID, &couchBook.CouchType, &couchBook.WidthCM, &couchBook.DepthCM, &couchBook.HeightCM); err != nil {
		return nil, err
	}

	return &couchBook, nil
}

func (b *BookingService) MakeMattressBooking(c context.Context, tx pgx.Tx, mattress *booking.MattressCleaningDetails) (*types.MattressCleaningDetails, error) {
	query := `
		INSERT INTO booking.mattress 
		(bed_type, width_cm, depth_cm, height_cm)
        VALUES ($1, $2, $3, $4)
        RETURNING id, bed_type, width_cm, depth_cm, height_cm`

	var mattressBook types.MattressCleaningDetails

	if err := tx.QueryRow(c, query, mattress.BedType, mattress.WidthCm, mattress.DepthCm, mattress.HeightCm).
		Scan(&mattressBook.ID, &mattressBook.BedType, &mattressBook.WidthCM, &mattressBook.DepthCM, &mattressBook.HeightCM); err != nil {
		return nil, err
	}

	return &mattressBook, nil
}

func (b *BookingService) MakeCarBooking(c context.Context, tx pgx.Tx, car *booking.CarCleaningDetails) (*types.CarCleaningDetails, error) {
	query := `
		INSERT INTO booking.car 
		(car_type, childseat)
        VALUES ($1, $2)
        RETURNING id, car_type, childseat`

	var carBook types.CarCleaningDetails

	if err := tx.QueryRow(c, query, car.CarType, car.ChildSeats).
		Scan(&carBook.ID, &carBook.CarType, &carBook.ChildSeats); err != nil {
		return nil, err
	}

	return &carBook, nil
}

func (b *BookingService) MakePostConstructionBooking(c context.Context, tx pgx.Tx, postConstruction *booking.PostConstructionCleaningDetails) (*types.PostConstructionDetails, error) {
	query := `
		INSERT INTO booking.postconstruction 
		(sqm)
        VALUES ($1)
        RETURNING id, sqm`

	var postConstructionBook types.PostConstructionDetails

	if err := tx.QueryRow(c, query, postConstruction.Sqm).
		Scan(&postConstructionBook.ID, &postConstructionBook.SQM); err != nil {
		return nil, err
	}

	return &postConstructionBook, nil
}

func (b *BookingService) createMainServiceBooking(
	ctx context.Context,
	tx pgx.Tx,
	mainService *booking.ServiceDetail,
) (types.ServiceDetail, error) {
	switch details := mainService.Type.(type) {
	case *booking.ServiceDetail_General:
		general, err := b.MakeGeneralBooking(ctx, tx, details.General)
		if err != nil {
			return types.ServiceDetail{}, err
		}
		return types.ServiceDetail{General: *general}, nil
	case *booking.ServiceDetail_Couch:
		couch, err := b.MakeCouchBooking(ctx, tx, details.Couch)
		if err != nil {
			return types.ServiceDetail{}, err
		}
		return types.ServiceDetail{Couch: *couch}, nil
	case *booking.ServiceDetail_Mattress:
		mattress, err := b.MakeMattressBooking(ctx, tx, details.Mattress)
		if err != nil {
			return types.ServiceDetail{}, err
		}
		return types.ServiceDetail{Mattress: *mattress}, nil
	case *booking.ServiceDetail_Car:
		car, err := b.MakeCarBooking(ctx, tx, details.Car)
		if err != nil {
			return types.ServiceDetail{}, err
		}
		return types.ServiceDetail{Car: *car}, nil
	case *booking.ServiceDetail_Post:
		post, err := b.MakePostConstructionBooking(ctx, tx, details.Post)
		if err != nil {
			return types.ServiceDetail{}, err
		}
		return types.ServiceDetail{Post: *post}, nil
	default:
		return types.ServiceDetail{}, fmt.Errorf("unsupported main service type")
	}
}

func (b *BookingService) createAddOn(
	ctx context.Context,
	tx pgx.Tx,
	addon *booking.AddOn,
) (*types.AddOns, error) {
	addOnServiceDetails, err := b.createMainServiceBooking(ctx, tx, addon.ServiceDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to create service details: %w", err)
	}

	var createdAddon types.AddOns
	createdAddon.ServiceDetail = addOnServiceDetails
	err = tx.QueryRow(ctx,
		`INSERT INTO booking.addons (addon_service_id, price)
		 VALUES ($1, $2)
		 RETURNING id, price`,
		addOnServiceDetails.GetID(),
		addon.Price,
	).Scan(
		&createdAddon.ID,
		&createdAddon.Price,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert addon: %w", err)
	}

	return &createdAddon, nil
}

func (b *BookingService) createEquipment(
	ctx context.Context,
	tx pgx.Tx,
	equipment *booking.CleaningEquipment,
) (*types.CleaningEquipment, error) {
	var createdEquipment types.CleaningEquipment
	err := tx.QueryRow(ctx,
		`INSERT INTO booking.cleaningequipments (name, type, photo_url)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, type, photo_url`,
		equipment.Name,
		equipment.Type,
		equipment.PhotoUrl,
	).Scan(
		&createdEquipment.ID,
		&createdEquipment.Name,
		&createdEquipment.Type,
		&createdEquipment.PhotoURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert cleaning equipment: %w", err)
	}

	return &createdEquipment, nil
}

func (b *BookingService) createResource(
	ctx context.Context,
	tx pgx.Tx,
	resource *booking.CleaningResources,
) (*types.CleaningResources, error) {
	var createdResource types.CleaningResources
	err := tx.QueryRow(ctx,
		`INSERT INTO booking.cleaningresources (name, type, photo_url)
		VALUES ($1, $2, $3)
		RETURNING id, name, type, photo_url`,
		resource.Name,
		resource.Type,
		resource.PhotoUrl,
	).Scan(
		&createdResource.ID,
		&createdResource.Name,
		&createdResource.Type,
		&createdResource.PhotoURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert cleaning resource: %w", err)
	}

	return &createdResource, nil
}

func (b *BookingService) createCleanersAssigned(
	ctx context.Context,
	tx pgx.Tx,
	cleaner *booking.CleanerAssigned,
) (*types.CleanerAssigned, error) {
	var createdCleaner types.CleanerAssigned
	err := tx.QueryRow(ctx,
		`INSERT INTO booking.cleanersassigned (cleaner_first_name, cleaner_last_name, pfp_url)
		VALUES ($1, $2, $3)
		RETURNING id, cleaner_first_name, cleaner_last_name, pfp_url`,
		cleaner.CleanerFirstName,
		cleaner.CleanerLastName,
		cleaner.PfpUrl,
	).Scan(
		&createdCleaner.ID,
		&createdCleaner.CleanerFirstName,
		&createdCleaner.CleanerLastName,
		&createdCleaner.PFPUrl,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert cleaners assinged: %w", err)
	}

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
		RETURNING id
`

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
