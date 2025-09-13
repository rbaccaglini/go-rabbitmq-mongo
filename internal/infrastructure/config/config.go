package config

// Config struct principal que agrupa todas as configurações
type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	RabbitMQ     RabbitMQConfig
	ExternalAPIs ExternalAPIsConfig
	Environment  string
}

// NewConfig cria uma nova instância de configuração
func NewConfig() *Config {
	return &Config{
		Server:       LoadServerConfig(),
		Database:     LoadDatabaseConfig(),
		RabbitMQ:     LoadRabbitMQConfig(),
		ExternalAPIs: LoadExternalAPIsConfig(),
		Environment:  GetEnv("APP_ENV", "development"),
	}
}
