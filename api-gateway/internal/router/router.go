package router

import (
	"handworks/api-gateway/internal/handlers"
	"handworks/api-gateway/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(l *utils.Logger) *gin.Engine {
	r := gin.New()

	// Health check
	r.GET("/health", handlers.HealthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define other routes here (e.g., bookings, employees, etc.)

	return r
}
