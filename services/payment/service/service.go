package service

import (
	"context"
	"fmt"
	"handworks/common/grpc/payment"
	types "handworks/common/types/payment"
	"handworks/common/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type PaymentService struct {
	L  *utils.Logger
	DB *pgxpool.Pool
	NC *nats.Conn
	payment.UnimplementedPaymentServiceServer
}

func (p *PaymentService) MakeQuotation(c context.Context, in *payment.QuoteRequest) (*payment.QuoteResponse, error) {
	var dbQuote types.DbQuote

	if in.CustId == "" {
		p.L.Info("Genearting Quote Preview (No Customer ID)")
		quote, err := p.CalculateQuotePreview(c, in)
		if err != nil {
			p.L.Error("Failed to genearte Quote Preview: %v", err)
			return nil, fmt.Errorf("failed to genearte Quote Preview: %v", err)
		}
		return quote.ToProto(), nil
	}

	if err := p.withTx(c, p.DB, func(tx pgx.Tx) error {
		quote, err := p.CreateQuote(c, tx, in)
		if err != nil {
			return err
		}
		dbQuote = *quote
		return nil
	}); err != nil {
		p.L.Error("Failed creating quote: %v", err)
		return nil, err
	}
	QuoteResponse := dbQuote.ToProto()
	return QuoteResponse, nil
}

// TODO: implement this later - Yousif
func (p *PaymentService) GetAllQuotesFromCustomer(c context.Context, in *payment.CustomerRequest) (*payment.QuotesResponse, error) {
	return nil, nil
}
