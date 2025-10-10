package service

import (
	"context"
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

func (p *PaymentService) GetQuotation(c context.Context, in *payment.QuoteRequest) (*payment.QuoteResponse, error) {
	var dbQuote types.DbQuote

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
