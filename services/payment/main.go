package main

import (
	"context"
	"handworks-services-payment/db"
	"handworks-services-payment/server"
	"handworks-services-payment/service"
	"handworks/common/grpc/payment"
	"handworks/common/natsconn"
	"handworks/common/utils"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	nc := natsconn.ConnectNATS()
	defer nc.Close()
	if nc != nil {
		logger.Info("NATS connection established")
	}

	ctx := context.Background()
	pool, err := db.InitDb(ctx)
	if err != nil {
		logger.Fatal("Payment DB Initialization Failed: %v", err)
	}
	logger.Info("Payment DB Initialization Success")
	defer pool.Close()

	paymentService := service.PaymentService{
		L:                                 logger,
		DB:                                pool,
		NC:                                nc,
		UnimplementedPaymentServiceServer: payment.UnimplementedPaymentServiceServer{},
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := server.StartGrpcServer(&paymentService, logger); err != nil {
			logger.Fatal("Initialization of Payment GRPC Server Failed: %v", err)
		}
	}()

	wg.Wait()
}
