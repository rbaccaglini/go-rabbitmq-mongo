package config

import "time"

const (
	// Environment
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"

	// Default Values
	DefaultServerPort   = "8080"
	DefaultServerHost   = "localhost"
	DefaultReadTimeout  = 30 * time.Second
	DefaultWriteTimeout = 30 * time.Second
)

// IsDevelopment verifica se o ambiente é de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.Environment == EnvDevelopment
}

// IsProduction verifica se o ambiente é de produção
func (c *Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

// IsStaging verifica se o ambiente é de staging
func (c *Config) IsStaging() bool {
	return c.Environment == EnvStaging
}
