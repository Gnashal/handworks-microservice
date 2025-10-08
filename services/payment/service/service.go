package service

import (
	"context"
	"handworks/common/grpc/payment"
	"handworks/common/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type PaymentService struct {
	L  *utils.Logger
	DB *pgxpool.Pool
	NC *nats.Conn
	payment.UnimplementedPaymentServiceServer
}

func (p *PaymentService) GetQuote(c context.Context, in *payment.QuoteRequest) (*payment.QuoteResponse, error) {

	return &payment.QuoteResponse{}, nil
}
