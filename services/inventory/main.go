package main

import (
	"context"
	"handworks-services-account/db"
	"handworks-services-inventory/server"
	"handworks-services-inventory/service"
	"handworks/common/grpc/inventory"
	"handworks/common/natsconn"
	"handworks/common/utils"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	nc := natsconn.ConnectNATS()
	defer nc.Close()
	if nc != nil {
		logger.Info("Connected to NATS")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := db.InitDb(ctx)
	if err != nil {
		logger.Fatal("Inventory DB Initialization Failed: %v", err)
	}
	logger.Info("Inventory DB Initialization Success")
	defer pool.Close()

	invService := &service.InventoryService{
		L:                                   logger,
		DB:                                  pool,
		NC:                                  nc,
		UnimplementedInventoryServiceServer: inventory.UnimplementedInventoryServiceServer{},
	}
	go func() {
		if err := invService.HandleSubscriptions(ctx); err != nil {
			logger.Error("NATS subscription error: %v", err)
			cancel()
		}
	}()

	go func() {
		if err := server.StartGrpcServer(ctx, invService, logger); err != nil {
			logger.Error("gRPC server error: %v", err)
			cancel()
		}
	}()

	// Block until cancel
	<-ctx.Done()
	logger.Warn("Shutting down Inventory service...")

	time.Sleep(1 * time.Second)
	logger.Info("Shutdown complete")
}
