package config

import (
	"time"

	"github.com/gin-contrib/cors"
)

func NewCors() cors.Config {
	CorsConf := cors.Config{
		// AllowOrigins:     []string{"url here soon, "http://localhost:5173"}, This is for prod
		AllowAllOrigins:  true, // Allow all origins for development
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return CorsConf
}
