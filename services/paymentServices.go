package services

import (
	"context"
	"fmt"
	"handworks-api/types"

	"github.com/jackc/pgx/v5"
)
func (s *PaymentService) withTx(
	ctx context.Context,
	fn func(pgx.Tx) error,
) (err error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.Logger.Error("rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()
	return fn(tx)
}

func (s *PaymentService) MakeQuotation(ctx context.Context, req types.QuoteRequest) (*types.QuoteResponse, error) {
	var quoteResponse types.QuoteResponse

	if req.CustomerID == "" {
		s.Logger.Info("Generating Quote Preview")
		quotePrev, err := s.Tasks.CalculateQuotePreview(ctx, &req)
		if err != nil {
			s.Logger.Error("Failed to genearte Quote Preview: %v", err)
			return nil, fmt.Errorf("failed to genearte Quote Preview: %v", err)
		}
		addonsBreakdown := s.Tasks.MapAddonstoAddonBreakdown(&quotePrev.Addons)
		return &types.QuoteResponse{
			QuoteId: quotePrev.ID,
			MainServiceName: quotePrev.MainService,
			MainServiceTotal: quotePrev.TotalPrice,
			TotalPrice: quotePrev.TotalPrice,
			AddonTotal: quotePrev.AddonTotal,
			Addons: addonsBreakdown,
		}, nil
	}
	if err := s.withTx(ctx, func (tx pgx.Tx) error {
		quote, err := s.Tasks.CreateQuote(ctx, tx, &req)
		if err != nil {
			return fmt.Errorf("failed to create Quote: %v", err)
		}
		quoteResponse.QuoteId = quote.ID
		quoteResponse.MainServiceName = quote.MainService
		quoteResponse.MainServiceTotal = quote.TotalPrice
		quoteResponse.AddonTotal = quote.AddonTotal
		quoteResponse.TotalPrice = quote.TotalPrice
		quoteResponse.Addons = s.Tasks.MapAddonstoAddonBreakdown(&quote.Addons)
		return nil
	}); err != nil {
		s.Logger.Error("Failed to create Quote: %v", err)
		return nil, err
	}
	return &quoteResponse, nil
}

// TODO: implement this
func (s *PaymentService) GetAllQuotesFromCustomer(ctx context.Context, id string) (*types.QuotesResponse, error) {
	return nil, nil
}