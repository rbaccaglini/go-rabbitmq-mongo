package api

import (
	"github.com/gin-gonic/gin"

	"api-rabbitmq/internal/infrastructure/http/handlers"
)

func SetupRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	// Health check
	router.GET("/health", userHandler.HealthCheck)

	// User routes
	userGroup := router.Group("/api/v1/users")
	{
		userGroup.POST("/publish", userHandler.PublishUser)
		userGroup.GET("/processed", userHandler.GetProcessedUsers)
	}
}
