package main

import (
	"context"
	"handworks-services-booking/db"
	"handworks-services-booking/server"
	"handworks-services-booking/service"
	"handworks/common/grpc/booking"
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

	ctx := context.Background()
	pool, err := db.InitDb(ctx)
	if err != nil {
		logger.Fatal("Booking DB Initialization Failed: %v", err)
	}
	logger.Info("Booking DB Initialization Success")
	defer pool.Close()

	bookService := service.BookingService{
		L:                                 logger,
		DB:                                pool,
		UnimplementedBookingServiceServer: booking.UnimplementedBookingServiceServer{},
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := server.StartGrpcServer(&bookService, logger); err != nil {
			logger.Fatal("Initialization of Account GRPC Server Failed: %v", err)
		}
	}()

	wg.Wait()
}
