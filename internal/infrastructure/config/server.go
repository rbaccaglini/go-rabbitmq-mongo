package config

import (
	"fmt"
	"time"
)

// ServerConfig configurações do servidor HTTP
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// LoadServerConfig carrega configurações do servidor
func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Port:         GetEnv("SERVER_PORT", "8080"),
		Host:         GetEnv("SERVER_HOST", "localhost"),
		ReadTimeout:  GetEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: GetEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:  GetEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
	}
}

// GetAddress retorna o endereço do servidor
func (s *ServerConfig) GetAddress() string {
	return s.Host + ":" + s.Port
}

// Validate valida as configurações do servidor
func (s *ServerConfig) Validate() error {
	if s.Port == "" {
		return fmt.Errorf("server port is required")
	}
	return nil
}
