package main

import (
	"handworks-services-account/server"
	"handworks-services-account/service"
	"handworks/common/grpc/account"
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

	accService := service.AccountService{
		L:                                 logger,
		UnimplementedAccountServiceServer: account.UnimplementedAccountServiceServer{},
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := server.StartGrpcServer(&accService, logger); err != nil {
			logger.Fatal("Initialization of Account GRPC Server Failed: %v", err)
		}
	}()

	wg.Wait()
}
