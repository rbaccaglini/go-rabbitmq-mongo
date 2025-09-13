package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"api-rabbitmq/internal/application/services"
	"api-rabbitmq/internal/application/usecases"
	"api-rabbitmq/internal/infrastructure/config"
	"api-rabbitmq/internal/infrastructure/database/mongodb"
	"api-rabbitmq/internal/infrastructure/http/handlers"
	"api-rabbitmq/internal/infrastructure/messagebroker/rabbitmq"
	"api-rabbitmq/internal/interfaces/api"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validar configurações
	if err := cfg.Database.Validate(); err != nil {
		log.Printf("Warning: Database configuration error: %v", err)
	}

	// Inicializar repositório
	userRepo, err := mongodb.NewUserRepository(cfg.Database.GetConnectionString())
	if err != nil {
		log.Printf("Warning: MongoDB not available: %v", err)
		userRepo = nil
	} else {
		defer userRepo.Close()
	}

	// Inicializar serviços externos
	extServices := services.NewExternalServices(&cfg.ExternalAPIs)

	// Inicializar use case
	userUseCase := usecases.NewUserUseCase(userRepo, extServices)

	// Inicializar RabbitMQ
	rabbitMQService, err := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URI, userUseCase)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQService.Close()

	// Iniciar consumo de mensagens
	go func() {
		if err := rabbitMQService.ConsumeMessages(); err != nil {
			log.Fatalf("Failed to start consuming messages: %v", err)
		}
	}()

	// Inicializar handlers
	userHandler := handlers.NewUserHandler(userUseCase, rabbitMQService)

	// Configurar router
	router := gin.Default()
	api.SetupRoutes(router, userHandler)

	// Iniciar servidor
	serverPort := cfg.Server.Port
	log.Printf("Server starting on %s", serverPort)
	if err := router.Run(":" + serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
