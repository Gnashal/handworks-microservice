package main

import (
	"handworks/api-gateway/config"
	_ "handworks/api-gateway/docs"
	"handworks/api-gateway/internal/router"
	"handworks/common/utils"
	"os"
)

// @title Handworks API Gateway
// @version 1.0
// @description Gateway for client requests to services
// @host localhost:8080
// @BasePath /

func main() {
	logger, err := utils.NewLogger()
	var unused int
	_ = config.NewCors()
	r := router.SetupRouter(logger)
	if err != nil {
		logger.Error("Failed to create loggerger: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("Swagger on http://localhost:%s/swagger/index.html", port)
	logger.Info("Starting server on :%s...", port)
	if err := r.Run(":" + port); err != nil {
		logger.Error("Server failed: %v", err)
	}
}
