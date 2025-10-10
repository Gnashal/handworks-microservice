package service

import (
	"context"
	"encoding/json"
	"fmt"

	btypes "handworks/common/types/booking"
	types "handworks/common/types/payment"

	"github.com/nats-io/nats.go"
)

func (p *PaymentService) HandleSubscriptions(ctx context.Context) error {
	if err := p.SubscribeBookingRequests(ctx); err != nil {
		p.L.Error("%v\n", err)
		return err
	}
	<-ctx.Done()
	return ctx.Err()

}

func (p *PaymentService) verifyAndGetQuotePrices(ctx context.Context, quoteId string) (types.CleaningPrices, error) {
	var prices types.CleaningPrices

	var dbQuote types.DbQuote
	err := p.DB.QueryRow(ctx, `
		SELECT total_price, is_valid
		FROM payment.quotes
		WHERE id = $1
	`, quoteId).Scan(
		&dbQuote.TotalPrice,
		&dbQuote.IsValid,
	)
	if err != nil {
		return prices, fmt.Errorf("fetch main quote: %w", err)
	}
	if !dbQuote.IsValid {
		return types.CleaningPrices{}, fmt.Errorf("quote is no longer valied")
	}
	rows, err := p.DB.Query(ctx, `
		SELECT service_type, addon_price
		FROM payment.quote_addons
		WHERE quote_id = $1
	`, quoteId)
	if err != nil {
		return prices, fmt.Errorf("fetch addons: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var addon types.DbQuoteAddon
		if err := rows.Scan(
			&addon.ServiceType,
			&addon.AddonPrice,
		); err != nil {
			return prices, fmt.Errorf("scan addon: %w", err)
		}
		dbQuote.Addons = append(dbQuote.Addons, &addon)
	}

	prices.MainServicePrice = dbQuote.Subtotal

	for _, a := range dbQuote.Addons {
		prices.AddonPrices = append(prices.AddonPrices, types.AddonCleaningPrice{
			AddonName:  a.ServiceType,
			AddonPrice: a.AddonPrice,
		})
	}
	prices.MainServicePrice = dbQuote.TotalPrice
	return prices, nil
}

func (p *PaymentService) SubscribeBookingRequests(ctx context.Context) error {
	_, err := p.NC.Subscribe("booking.created", func(msg *nats.Msg) {
		var req btypes.CreateBookingEvent
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			p.L.Error("Failed to unmarshal booking request: %v\n", err)
			return
		}
		p.L.Info("Received booking.created event: %s", req.Base.CustomerFirstName)
		prices, err := p.verifyAndGetQuotePrices(ctx, req.Base.QuoteId)
		if err != nil {
			p.L.Error("Failed to get quote prices: %v\n", err)
			return
		}
		reply := types.BookingReply{
			Source: "payment-service",
			Prices: prices,
		}
		data, err := json.Marshal(reply)
		if err != nil {
			p.L.Error("Failed to marshal reply: %v\n", err)
			return
		}

		if msg.Reply != "" {
			if err := p.NC.Publish(msg.Reply, data); err != nil {
				p.L.Error("Failed to publish reply: %v\n", err)
			} else {
				p.L.Info("Sent booking reply from payment-service")
			}
		} else {
			p.L.Info("No reply subject — skipping reply")
		}
	})
	if err != nil {
		return err
	}
	return nil
}
