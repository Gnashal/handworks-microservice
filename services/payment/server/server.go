package server

import (
	"fmt"
	"handworks-services-payment/service"
	"handworks/common/grpc/payment"
	"handworks/common/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func StartGrpcServer(s *service.PaymentService, logger *utils.Logger) error {
	host := getEnv("DEV_URL", "localhost")
	port := getEnv("PORT", "9092")

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	server := grpc.NewServer()
	payment.RegisterPaymentServiceServer(server, s)

	go gracefulShutdown(server, lis, logger)
	logger.Info("Payment Service Ready")
	logger.Info("Payment Service gRPC server listening on %s", addr)
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("gRPC server failed: %w", err)
	}
	return nil
}

func gracefulShutdown(server *grpc.Server, lis net.Listener, logger *utils.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Warn("Shutting down gRPC server...")
	server.GracefulStop()
	lis.Close()
}

// tiny helper for env lookup
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
