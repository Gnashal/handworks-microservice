package main

import (
	"context"
	"handworks-services-account/db"
	"handworks-services-inventory/server"
	"handworks-services-inventory/service"
	"handworks/common/grpc/inventory"
	"handworks/common/utils"

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

	ctx := context.Background()
	pool, err := db.InitDb(ctx)
	if err != nil {
		logger.Fatal("Inventory DB Initialization Failed: %v", err)
	}
	logger.Info("Inventory DB Initialization Success")
	defer pool.Close()

	accService := &service.InventoryService{
		L:                                   logger,
		DB:                                  pool,
		UnimplementedInventoryServiceServer: inventory.UnimplementedInventoryServiceServer{},
	}

	if err := server.StartGrpcServer(accService, logger); err != nil {
		logger.Fatal("Inventory gRPC Server failed: %v", err)
	}
}
