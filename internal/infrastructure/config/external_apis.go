package config

import (
	"fmt"
	"time"
)

// ExternalAPIsConfig configurações das APIs externas
type ExternalAPIsConfig struct {
	DocumentValidationURL string
	AddressServiceURL     string
	Timeout               time.Duration
	RetryAttempts         int
	RetryDelay            time.Duration
}

// LoadExternalAPIsConfig carrega configurações das APIs externas
func LoadExternalAPIsConfig() ExternalAPIsConfig {
	return ExternalAPIsConfig{
		DocumentValidationURL: GetEnv("DOCUMENT_VALIDATION_URL", "http://localhost:8082/api/v1/is-document-valid"),
		AddressServiceURL:     GetEnv("ADDRESS_SERVICE_URL", "http://localhost:8081/api/v1/address"),
		Timeout:               GetEnvDuration("EXTERNAL_API_TIMEOUT", 30*time.Second),
		RetryAttempts:         GetEnvInt("EXTERNAL_API_RETRY_ATTEMPTS", 3),
		RetryDelay:            GetEnvDuration("EXTERNAL_API_RETRY_DELAY", 1*time.Second),
	}
}

// Validate valida as configurações das APIs externas
func (e *ExternalAPIsConfig) Validate() error {
	if e.DocumentValidationURL == "" {
		return fmt.Errorf("document validation URL is required")
	}
	if e.AddressServiceURL == "" {
		return fmt.Errorf("address service URL is required")
	}
	return nil
}
