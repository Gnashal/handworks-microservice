package config

import (
	"time"

	"github.com/gin-contrib/cors"
)

func NewCors() cors.Config {
	CorsConf := cors.Config{
		AllowOrigins:     []string{"https://griita-backend-gnashal6914-x2n9tdsh.leapcell.dev", "http://localhost:4040"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return CorsConf
}
