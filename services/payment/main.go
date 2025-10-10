package main

import (
	"context"
	"handworks-services-payment/db"
	"handworks-services-payment/server"
	"handworks-services-payment/service"
	"handworks/common/grpc/payment"
	"handworks/common/natsconn"
	"handworks/common/utils"
	"time"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := db.InitDb(ctx)
	if err != nil {
		logger.Fatal("Payment DB Initialization Failed: %v", err)
	}
	logger.Info("Payment DB Initialization Success")
	defer pool.Close()

	paymentService := &service.PaymentService{
		L:                                 logger,
		DB:                                pool,
		NC:                                nc,
		UnimplementedPaymentServiceServer: payment.UnimplementedPaymentServiceServer{},
	}
	go func() {
		if err := paymentService.HandleSubscriptions(ctx); err != nil {
			logger.Error("NATS subscription error: %v", err)
			cancel()
		}
	}()

	go func() {
		if err := server.StartGrpcServer(ctx, paymentService, logger); err != nil {
			logger.Error("gRPC server error: %v", err)
			cancel()
		}
	}()

	// Block until cancel
	<-ctx.Done()
	logger.Warn("Shutting down Account service...")

	time.Sleep(1 * time.Second)
	logger.Info("Shutdown complete")
}
