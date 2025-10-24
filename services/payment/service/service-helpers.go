package service

import (
	"context"
	"encoding/json"
	"fmt"
	"handworks/common/grpc/booking"
	"handworks/common/grpc/payment"
	types "handworks/common/types/payment"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (p *PaymentService) withTx(
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
				p.L.Error("rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()
	return fn(tx)
}

func (p *PaymentService) CreateQuote(c context.Context, tx pgx.Tx, in *payment.QuoteRequest) (*types.DbQuote, error) {
	var dbQuote types.DbQuote
	var dbAddons []*types.DbQuoteAddon

	mainService := &booking.ServicesRequest{
		ServiceType: in.Service.ServiceType,
		Details:     in.Service.Details,
	}

	// Calc subtotal for main service
	subtotal := p.CalculatePriceByServiceType(mainService)
	var addonTotal float32 = 0

	// Calculate each addon price
	for _, addon := range in.Addons {
		// i genuinely dunno why addon.ServiceDetail wont work lmao
		addonService := &booking.ServicesRequest{
			ServiceType: addon.ServiceDetail.ServiceType,
			Details:     addon.ServiceDetail.Details,
		}
		addonPrice := p.CalculatePriceByServiceType(addonService)

		// serialize the full addon service detail
		serviceDetail, err := json.Marshal(addon.ServiceDetail)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal addon service: %v", err)
		}

		addonTotal += addonPrice

		dbAddon := &types.DbQuoteAddon{
			ServiceType:   addon.GetServiceDetail().ServiceType.String(),
			ServiceDetail: serviceDetail,
			AddonPrice:    addonPrice,
			CreatedAt:     time.Now(),
		}
		dbAddons = append(dbAddons, dbAddon)
	}

	totalPrice := subtotal + addonTotal

	// Insert into quote table
	err := tx.QueryRow(c, `
		INSERT INTO payment.quotes (customer_id, main_service_type, subtotal, addon_total, total_price, is_valid)
		VALUES ($1, $2, $3, $4, $5, TRUE)
		RETURNING id, customer_id, main_service_type, subtotal, addon_total, total_price, is_valid, created_at, updated_at
	`,
		in.CustId,
		in.Service.ServiceType.String(),
		subtotal,
		addonTotal,
		totalPrice,
	).Scan(
		&dbQuote.ID,
		&dbQuote.CustomerID,
		&dbQuote.MainService,
		&dbQuote.Subtotal,
		&dbQuote.AddonTotal,
		&dbQuote.TotalPrice,
		&dbQuote.IsValid,
		&dbQuote.CreatedAt,
		&dbQuote.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert quote: %v", err)
	}

	// insert addons
	for _, addon := range dbAddons {
		err := tx.QueryRow(c, `
			INSERT INTO payment.quote_addons (quote_id, service_type, service_detail, addon_price)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at
		`,
			dbQuote.ID,
			addon.ServiceType,
			addon.ServiceDetail,
			addon.AddonPrice,
		).Scan(&addon.ID, &addon.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to insert addon: %v", err)
		}

		addon.QuoteID = dbQuote.ID
	}

	dbQuote.Addons = dbAddons
	return &dbQuote, nil
}

func (p *PaymentService) CalculateQuotePreview(c context.Context, in *payment.QuoteRequest) (*types.DbQuote, error) {
	var dbQuote types.DbQuote
	var dbAddons []*types.DbQuoteAddon

	mainService := &booking.ServicesRequest{
		ServiceType: in.Service.ServiceType,
		Details:     in.Service.Details,
	}

	subtotal := p.CalculatePriceByServiceType(mainService)
	var addonTotal float32 = 0

	for _, addon := range in.Addons {
		addonService := &booking.ServicesRequest{
			ServiceType: addon.ServiceDetail.ServiceType,
			Details:     addon.ServiceDetail.Details,
		}
		addonPrice := p.CalculatePriceByServiceType(addonService)

		serviceDetail, err := json.Marshal(addon.ServiceDetail)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal addon service: %v", err)
		}

		addonTotal += addonPrice
		dbAddon := &types.DbQuoteAddon{
			ServiceType:   addon.GetServiceDetail().ServiceType.String(),
			ServiceDetail: serviceDetail,
			AddonPrice:    addonPrice,
			CreatedAt:     time.Now(),
		}
		dbAddons = append(dbAddons, dbAddon)
	}

	dbQuote = types.DbQuote{
		CustomerID:  in.CustId, // will be empty
		MainService: in.Service.ServiceType.String(),
		Subtotal:    subtotal,
		AddonTotal:  addonTotal,
		TotalPrice:  subtotal + addonTotal,
		IsValid:     false, // marked as preview only so di siya valid
		CreatedAt:   time.Now(),
		Addons:      dbAddons,
	}

	return &dbQuote, nil
}

// func (p *PaymentService) GetQuotesByCustomerId(
// 	c context.Context,
// 	tx pgx.Tx,
// 	custID string,
// ) ([]*types.DbQuote, error) {
// 	// Fetch base quotes
// 	rows, err := tx.Query(c, `
// 		SELECT id, main_service_name, main_service_total, addon_total, total_price
// 		FROM payment.quotes
// 		WHERE cust_id = $1
// 	`, custID)
// 	if err != nil {
// 		return nil, fmt.Errorf("fetch quotes: %w", err)
// 	}
// 	defer rows.Close()

// 	var quotes []*types.DbQuote
// 	for rows.Next() {
// 		var q types.DbQuote
// 		if err := rows.Scan(
// 			&q.ID,
// 			&q.MainService,
// 			&q.TotalPrice,
// 			&q.AddonTotal,
// 			&q.TotalPrice,
// 		); err != nil {
// 			return nil, fmt.Errorf("scan quote: %w", err)
// 		}

// 		// Fetch addon breakdown for each quote
// 		addonRows, err := tx.Query(c, `
// 			SELECT id, quote_id, service_type, service_detail, addon_price, created_at
// 			FROM payment.quote_addons
// 			WHERE quote_id = $1
// 		`, q.ID)
// 		if err != nil {
// 			return nil, fmt.Errorf("fetch addons: %w", err)
// 		}

// 		for addonRows.Next() {
// 			var a types.DbQuoteAddon
// 			if err := addonRows.Scan(&a.ID, &a.Name, &a.Price); err != nil {
// 				addonRows.Close()
// 				return nil, fmt.Errorf("scan addon: %w", err)
// 			}
// 			q.Addons = append(q.Addons, a)
// 		}
// 		addonRows.Close()

// 		quotes = append(quotes, &q)
// 	}

// 	return quotes, nil
// }
