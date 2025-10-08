package server

import (
	"context"
	"fmt"
	"handworks-services-payment/service"
	"handworks/common/grpc/payment"
	"handworks/common/utils"
	"net"
	"os"

	"google.golang.org/grpc"
)

func StartGrpcServer(ctx context.Context, s *service.PaymentService, logger *utils.Logger) error {
	host := getEnv("DEV_URL", "localhost")
	port := getEnv("PORT", "9090")
	addr := fmt.Sprintf("%s:%s", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	grpcServer := grpc.NewServer()
	payment.RegisterPaymentServiceServer(grpcServer, s)

	// graceful shutdown
	go func() {
		<-ctx.Done()
		logger.Warn("Gracefully stopping gRPC server...")
		grpcServer.GracefulStop()
		lis.Close()
	}()

	logger.Info("Payment Service gRPC server listening on %s", addr)
	return grpcServer.Serve(lis)
}

// tiny helper for env lookup
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
