package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"api-rabbitmq/models"
	"api-rabbitmq/services"
)

func main() {
	// Conectar ao MongoDB
	mongoURI := "mongodb://admin:password@localhost:27017"
	mongoRepo, err := models.NewMongoDBRepository(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoRepo.Close()

	// Conectar ao RabbitMQ
	rabbitMQURL := "amqp://admin:password@localhost:5672/"
	rabbitService, err := services.NewRabbitMQService(rabbitMQURL, mongoRepo)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ service: %v", err)
	}
	defer rabbitService.Close()

	go func() {
		if err := rabbitService.ConsumeMessages(); err != nil {
			log.Fatalf("Failed to start consuming messages: %v", err)
		}
	}()

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API is running",
		})
	})

	// Publicar mensagem
	router.POST("/publish", func(c *gin.Context) {
		var userData models.UserData
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON format",
			})
			return
		}

		if err := rabbitService.PublishMessage(userData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to publish message: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Message published successfully",
			"data":    userData,
		})
	})

	// Buscar dados processados
	router.GET("/processed-data", func(c *gin.Context) {
		data, err := mongoRepo.GetAllProcessedData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch data: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  data,
			"count": len(data),
		})
	})

	// Status das conexões
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"rabbitmq_connected": rabbitService.IsConnected(),
			"mongo_connected":    true,
			"queue":              "user_data_queue",
		})
	})

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
