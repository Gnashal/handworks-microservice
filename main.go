package main

import (
	"context"
	"handworks-api/config"
	"handworks-api/endpoints"
	"handworks-api/handlers"
	"handworks-api/middleware"
	"handworks-api/services"
	"handworks-api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	c := context.Background()
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(cors.New(config.NewCors()))
	conn, err := config.InitDB(logger, c)
	if err != nil {
		logger.Fatal("DB init failed: %v", err)
	}
	logger.Info("Connected to Db")
	defer conn.Close()
	// public paths for Clerk middleware
	publicPaths := []string{"/signup", "/health"}
	router.Use(middleware.ClerkAuthMiddleware(publicPaths))

	accountService := services.NewAccountService(conn, logger)
	inventoryService := services.NewInventoryService(conn, logger)
	bookingService := services.NewBookingService(conn, logger)
	paymentService := services.NewPaymentService(conn, logger)

	accountHandler := handlers.NewAccountHandler(accountService, logger)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService, logger)
	bookingHandler := handlers.NewBookingHandler(bookingService, logger)
	paymentHandler := handlers.NewPaymentHandler(paymentService, logger)

	api := router.Group("/api")
	{
		endpoints.AccountEndpoint(api.Group("/account"), accountHandler)
		endpoints.InventoryEndpoint(api.Group("/inventory"), inventoryHandler)
		endpoints.BookingEndpoint(api.Group("/booking"), bookingHandler)
		endpoints.PaymentEndpoint(api.Group("/payment"), paymentHandler)
	}

	port := "8080"
	logger.Info("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Server failed: %v", err)
	}
}