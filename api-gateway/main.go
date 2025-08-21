package main

import (
	"handworks/api-gateway/config"
	_ "handworks/api-gateway/docs"
	"handworks/api-gateway/internal/router"
	"handworks/common/utils"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Handworks API Gateway
// @version 1.0
// @description Gateway for client requests to services
// @host localhost:8080
// @BasePath /

func main() {
	l, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	err = godotenv.Load()
	if err != nil {
		l.Error("Error loading .env file: %v", err)
		return
	}
	r := gin.New()
	r.Use(cors.New(config.NewCors()))
	router.SetupRouter(l, r)
	startAndStopServer(r, l)
}

func startAndStopServer(r *gin.Engine, l *utils.Logger) {
	port := os.Getenv("PORT")
	ip := os.Getenv("IP")
	l.Info("Swagger on http://%s:%s/swagger/index.html", ip, port)
	l.Info("Starting server on :%s...", port)
	if err := r.Run(":" + port); err != nil {
		l.Error("Server failed: %v", err)
	}
	// TODO: will add graceful shutdown later
}
