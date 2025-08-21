package router

import (
	"handworks/api-gateway/internal/handlers"
	"handworks/api-gateway/internal/types"
	"handworks/common/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ApiGatewayRouter struct {
	s *types.ApiGatewayServices
	r *gin.RouterGroup
	l *utils.Logger
}

func NewRouter(l *utils.Logger, r *gin.RouterGroup, s *types.ApiGatewayServices) *ApiGatewayRouter {
	return &ApiGatewayRouter{
		s: s,
		r: r,
		l: l,
	}
}

func SetupRouter(l *utils.Logger, r *gin.Engine) {
	serviceApi := r.Group("/api/")
	{
		// Health check
		serviceApi.GET("/health", handlers.HealthCheck)
		// Swagger
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		apiServices := handlers.NewGatewayServices(l)
		apiRouter := NewRouter(l, serviceApi, apiServices)
		apiRouter.StartRouterServices()
	}
}
func (r *ApiGatewayRouter) StartRouterServices() {
	r.l.Info("Starting API Gateway services...")
	r.l.Info("API Gateway services started successfully.")
	// Servies here
}
