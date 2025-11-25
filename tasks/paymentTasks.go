package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"handworks-api/types"
	"time"

	"github.com/jackc/pgx/v5"
)

type PaymentTasks struct {}


func CalculateGeneralCleaning(details *types.GeneralCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}
	sqm := details.SQM
	homeType := details.HomeType

	switch {
	case homeType == "CONDO_ROOM" || (sqm > 0 && sqm <= 30):
		return 2000.00
	case homeType == "HOUSE" || (sqm > 30 && sqm <= 50):
		return 2500.00
	case sqm > 50 && sqm <= 100:
		return 5000.00
	default:
		return float32(sqm * 50)
	}
}

func CalculateCarCleaning(details *types.CarCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}

	var total float32
	for _, spec := range details.CleaningSpecs {
		price := types.CarPrices[spec.CarType]
		total += price * float32(spec.Quantity)
	}

	if details.ChildSeats > 0 {
		total += float32(details.ChildSeats) * 250.00
	}

	return total
}

func CalculateCouchCleaning(details *types.CouchCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}

	var total float32
	for _, spec := range details.CleaningSpecs {
		price := types.CouchPrices[spec.CouchType]
		total += price * float32(spec.Quantity)
	}

	if details.BedPillows > 0 {
		total += float32(details.BedPillows) * 100.00
	}

	return total
}

func CalculateMattressCleaning(details *types.MattressCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}

	var total float32
	for _, spec := range details.CleaningSpecs {
		price := types.MattressPrices[spec.BedType]
		total += price * float32(spec.Quantity)
	}
	return total
}
func CalculatePostConstructionCleaning(details *types.PostConstructionDetails) float32 {
	if details == nil {
		return 0.0
	}
	return float32(details.SQM * 50.00)
}

func (t *PaymentTasks) CalculatePriceByServiceType(service *types.ServicesRequest) float32 {
	if service == nil {
		return 0
	}

	var calculatedPrice float32 = 0.00

	switch service.ServiceType {
	case types.GeneralCleaning:
		calculatedPrice = CalculateGeneralCleaning(service.Details.General)

	case types.CouchCleaning:
		calculatedPrice = CalculateCouchCleaning(service.Details.Couch)

	case types.MattressCleaning:
		calculatedPrice = CalculateMattressCleaning(service.Details.Mattress)

	case types.CarCleaning:
		calculatedPrice = CalculateCarCleaning(service.Details.Car)

	case types.PostCleaning:
		calculatedPrice = CalculatePostConstructionCleaning(service.Details.Post)

	default:
		// no default action
	}

	return calculatedPrice
}

func (t *PaymentTasks) CalculateQuotePreview(c context.Context, in *types.QuoteRequest) (*types.Quote, error) {
	var dbQuote types.Quote
	var dbAddons []*types.QuoteAddon

	mainService := &types.ServicesRequest{
		ServiceType: in.Service.ServiceType,
		Details:     in.Service.Details,
	}

	subtotal := t.CalculatePriceByServiceType(mainService)
	var addonTotal float32 = 0

	for _, addon := range in.Addons {
		addonService := &types.ServicesRequest{
			ServiceType: addon.ServiceDetail.ServiceType,
			Details:     addon.ServiceDetail.Details,
		}
		addonPrice := t.CalculatePriceByServiceType(addonService)

		serviceDetail, err := json.Marshal(addon.ServiceDetail)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal addon service: %v", err)
		}

		addonTotal += addonPrice
		dbAddon := &types.QuoteAddon{
			ServiceType:   string(addon.ServiceDetail.ServiceType),
			ServiceDetail: serviceDetail,
			AddonPrice:    addonPrice,
			CreatedAt:     time.Now(),
		}
		dbAddons = append(dbAddons, dbAddon)
	}

	dbQuote = types.Quote{
		ID: 		"",	
		CustomerID:  in.CustomerID, // will be empty
		MainService: string(in.Service.ServiceType),
		Subtotal:    subtotal,
		AddonTotal:  addonTotal,
		TotalPrice:  subtotal + addonTotal,
		IsValid:     false, // marked as preview only so di siya valid
		CreatedAt:   time.Now(),
		Addons:      dbAddons,
	}

	return &dbQuote, nil
}
func (t* PaymentTasks) MapAddonstoAddonBreakdown(addons* []*types.QuoteAddon) []types.AddOnBreakdown {
	var breakdowns []types.AddOnBreakdown
	for _, addon := range *addons {
		breakdown := types.AddOnBreakdown{
			AddonID:   addon.ID,
			AddonName: addon.ServiceType,
			Price:     float64(addon.AddonPrice),
		}
		breakdowns = append(breakdowns, breakdown)
	}
	return breakdowns
}
func (p *PaymentTasks) CreateQuote(c context.Context, tx pgx.Tx, in *types.QuoteRequest) (*types.Quote, error) {
	var dbQuote types.Quote
	var dbAddons []*types.QuoteAddon

	mainService := &types.ServicesRequest{
		ServiceType: in.Service.ServiceType,
		Details:     in.Service.Details,
	}

	// Calc subtotal for main service
	subtotal := p.CalculatePriceByServiceType(mainService)
	var addonTotal float32 = 0

	// Calculate each addon price
	for _, addon := range in.Addons {
		// i genuinely dunno why addon.ServiceDetail wont work lmao
		addonService := &types.ServicesRequest{
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

		dbAddon := &types.QuoteAddon{
			ServiceType:   string(addon.ServiceDetail.ServiceType),
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
		in.CustomerID,
		in.Service.ServiceType,
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