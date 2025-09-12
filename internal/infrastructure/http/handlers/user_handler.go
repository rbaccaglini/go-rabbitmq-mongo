package handlers

import (
	"api-rabbitmq/internal/infrastructure/messagebroker/rabbitmq"
	"net/http"

	"github.com/gin-gonic/gin"

	"api-rabbitmq/internal/application/usecases"
	"api-rabbitmq/internal/domain/entities"
)

type UserHandler struct {
	userUseCase     usecases.UserUseCase
	rabbitMQService *rabbitmq.RabbitMQService
}

func NewUserHandler(userUseCase usecases.UserUseCase, rabbitMQService *rabbitmq.RabbitMQService) *UserHandler {
	return &UserHandler{
		userUseCase:     userUseCase,
		rabbitMQService: rabbitMQService,
	}
}

func (h *UserHandler) PublishUser(c *gin.Context) {
	var userData entities.UserData
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := h.rabbitMQService.PublishMessage(userData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message published successfully",
		"data":    userData,
	})
}

func (h *UserHandler) GetProcessedUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.userUseCase.GetProcessedUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"count": len(users),
	})
}

func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":          "ok",
		"rabbitmq_status": h.rabbitMQService.IsConnected(),
		"message":         "API is running",
	})
}
