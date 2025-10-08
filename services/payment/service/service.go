package service

import (
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
