package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"api-rabbitmq/handlers"
	"api-rabbitmq/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	queueName = "user_data_queue"
)

type RabbitMQService struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	handler   *handlers.MessageHandler
	connected bool
}

func NewRabbitMQService(rabbitMQURL string, mongoRepo *models.MongoDBRepository) (*RabbitMQService, error) {
	log.Printf("Attempting to connect to RabbitMQ...")

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	_, err = channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	handler := handlers.NewMessageHandler(mongoRepo)

	log.Printf("Successfully connected to RabbitMQ and declared queue: %s", queueName)

	return &RabbitMQService{
		conn:      conn,
		channel:   channel,
		handler:   handler,
		connected: true,
	}, nil
}

func (r *RabbitMQService) ConsumeMessages() error {
	err := r.channel.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := r.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)

			var userData models.UserData
			if err := json.Unmarshal(msg.Body, &userData); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				msg.Nack(false, false)
				continue
			}

			result, err := r.handler.ProcessMessage(userData)
			if err != nil {
				log.Printf("Error processing message: %v", err)
				msg.Nack(false, true)
			} else {
				log.Printf("Message processed successfully: %+v", result)
				msg.Ack(false)
			}
		}
	}()

	log.Printf("Waiting for messages on queue: %s", queueName)
	return nil
}

func (r *RabbitMQService) PublishMessage(userData models.UserData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	err = r.channel.PublishWithContext(ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	log.Printf("Published message to queue '%s': %s", queueName, body)
	return nil
}

func (r *RabbitMQService) IsConnected() bool {
	return r.connected && !r.conn.IsClosed()
}

func (r *RabbitMQService) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	r.connected = false
	log.Println("RabbitMQ connection closed")
}
