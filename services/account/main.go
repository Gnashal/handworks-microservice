package main

import (
	"context"
	"handworks-services-account/db"
	"handworks-services-account/server"
	"handworks-services-account/service"
	"handworks/common/grpc/account"
	"handworks/common/utils"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pool, err := db.InitDb(ctx)
	if err != nil {
		logger.Fatal("Account DB Initialization Failed: %v", err)
	}
	logger.Info("Account DB Initialization Success")
	defer pool.Close()

	accService := &service.AccountService{
		L:                                 logger,
		DB:                                pool,
		UnimplementedAccountServiceServer: account.UnimplementedAccountServiceServer{},
	}

	if err := server.StartGrpcServer(accService, logger); err != nil {
		logger.Fatal("Account gRPC Server failed: %v", err)
	}
}
