package server

import (
	"fmt"
	"handworks-services-inventory/service"
	"handworks/common/grpc/inventory"
	"handworks/common/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func StartGrpcServer(s *service.InventoryService, logger *utils.Logger) error {
	host := getEnv("DEV_URL", "localhost")
	port := getEnv("PORT", "9090")
	addr := fmt.Sprintf("%s:%s", host, port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	grpcServer := grpc.NewServer()
	inventory.RegisterInventoryServiceServer(grpcServer, s)

	// graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		logger.Warn("Shutting down Inventory gRPC server...")
		grpcServer.GracefulStop()
		lis.Close()
	}()

	logger.Info("Inventory Service gRPC server listening on %s", addr)
	return grpcServer.Serve(lis)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
