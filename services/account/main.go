package main

import (
	"context"
	"handworks-services-account/db"
	"handworks-services-account/server"
	"handworks-services-account/service"
	"handworks/common/grpc/account"
	"handworks/common/natsconn"
	"handworks/common/utils"
	"os"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
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
		logger.Fatal("Account DB Initialization Failed: %v", err)
	}
	logger.Info("Account DB Initialization Success")
	defer pool.Close()

	// Initialize the full Clerk client
	clerkKey := os.Getenv("CLERK_SECRET_KEY")
	clerk.SetKey(clerkKey)
	logger.Info("Clerked initialized")
	accService := &service.AccountService{
		L:                                 logger,
		DB:                                pool,
		NC:                                nc,
		UnimplementedAccountServiceServer: account.UnimplementedAccountServiceServer{},
	}

	go func() {
		if err := accService.HandleSubscriptions(ctx); err != nil {
			logger.Error("NATS subscription error: %v", err)
			cancel()
		}
	}()

	go func() {
		if err := server.StartGrpcServer(ctx, accService, logger); err != nil {
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
