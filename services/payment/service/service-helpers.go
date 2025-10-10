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
