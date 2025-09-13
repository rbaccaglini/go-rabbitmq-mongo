package config

import "fmt"

// RabbitMQConfig configurações do RabbitMQ
type RabbitMQConfig struct {
	URI           string
	QueueName     string
	PrefetchCount int
	PrefetchSize  int
	Durable       bool
	AutoAck       bool
}

// LoadRabbitMQConfig carrega configurações do RabbitMQ
func LoadRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		URI:           GetEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		QueueName:     GetEnv("RABBITMQ_QUEUE_NAME", "user_data_queue"),
		PrefetchCount: GetEnvInt("RABBITMQ_PREFETCH_COUNT", 1),
		PrefetchSize:  GetEnvInt("RABBITMQ_PREFETCH_SIZE", 0),
		Durable:       GetEnvBool("RABBITMQ_DURABLE", true),
		AutoAck:       GetEnvBool("RABBITMQ_AUTO_ACK", false),
	}
}

// Validate valida as configurações do RabbitMQ
func (r *RabbitMQConfig) Validate() error {
	if r.URI == "" {
		return fmt.Errorf("RabbitMQ URI is required")
	}
	if r.QueueName == "" {
		return fmt.Errorf("RabbitMQ queue name is required")
	}
	return nil
}
