package api

import (
	"github.com/gin-gonic/gin"

	"api-rabbitmq/internal/infrastructure/http/handlers"
)

func SetupRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	// Health check
	router.GET("/health", userHandler.HealthCheck)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Users
		users := v1.Group("/users")
		{
			users.POST("/publish", userHandler.PublishUser)
			users.GET("/processed", userHandler.GetProcessedUsers)
		}
	}
}
