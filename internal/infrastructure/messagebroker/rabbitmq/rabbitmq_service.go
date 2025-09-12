package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"api-rabbitmq/internal/application/usecases"
	"api-rabbitmq/internal/domain/entities"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	userUseCase usecases.UserUseCase
	connected   bool
}

func NewRabbitMQService(rabbitMQURL string, userUseCase usecases.UserUseCase) (*RabbitMQService, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	_, err = channel.QueueDeclare(
		"user_data_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	return &RabbitMQService{
		conn:        conn,
		channel:     channel,
		userUseCase: userUseCase,
		connected:   true,
	}, nil
}

func (s *RabbitMQService) ConsumeMessages() error {
	err := s.channel.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := s.channel.Consume(
		"user_data_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}

	go s.processMessages(msgs)
	log.Printf("Waiting for messages on queue: user_data_queue")
	return nil
}

func (s *RabbitMQService) processMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		ctx := context.Background()

		var userData entities.UserData
		if err := json.Unmarshal(msg.Body, &userData); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			msg.Nack(false, false)
			continue
		}

		_, err := s.userUseCase.ProcessUser(ctx, userData)
		if err != nil {
			log.Printf("Error processing message: %v", err)
			msg.Nack(false, true)
		} else {
			log.Printf("Message processed successfully")
			msg.Ack(false)
		}
	}
}

func (s *RabbitMQService) PublishMessage(userData entities.UserData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	err = s.channel.PublishWithContext(ctx,
		"",
		"user_data_queue",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	log.Printf("Published message to queue")
	return nil
}

func (s *RabbitMQService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	s.connected = false
}

func (s *RabbitMQService) IsConnected() bool {
	return s.connected && s.conn != nil && !s.conn.IsClosed()
}
