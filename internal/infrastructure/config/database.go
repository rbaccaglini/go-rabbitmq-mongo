package config

import (
	"fmt"
	"time"
)

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	URI               string
	DatabaseName      string
	ConnectionTimeout time.Duration
	MaxPoolSize       uint64
	MinPoolSize       uint64
	SSL               bool
}

// LoadDatabaseConfig carrega configurações do banco de dados
func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		URI:               GetEnv("MONGODB_URI", "mongodb://localhost:27017/userdb"),
		DatabaseName:      GetEnv("MONGODB_DATABASE", "userdb"),
		ConnectionTimeout: GetEnvDuration("MONGODB_CONNECTION_TIMEOUT", 10*time.Second),
		MaxPoolSize:       uint64(GetEnvInt("MONGODB_MAX_POOL_SIZE", 100)),
		MinPoolSize:       uint64(GetEnvInt("MONGODB_MIN_POOL_SIZE", 1)),
		SSL:               GetEnvBool("MONGODB_SSL", false),
	}
}

// GetConnectionString retorna a string de conexão formatada
func (d *DatabaseConfig) GetConnectionString() string {
	return d.URI
}

// Validate valida as configurações do banco de dados
func (d *DatabaseConfig) Validate() error {
	if d.URI == "" {
		return fmt.Errorf("database URI is required")
	}
	if d.DatabaseName == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}
