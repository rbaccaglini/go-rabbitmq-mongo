package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"api-rabbitmq/internal/domain/entities"
)

type ExternalServicesImpl struct {
	httpClient *http.Client
}

func NewExternalServices() *ExternalServicesImpl {
	return &ExternalServicesImpl{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *ExternalServicesImpl) ValidateDocument(documentNumber string) (bool, error) {
	url := fmt.Sprintf("http://localhost:8082/api/v1/is-document-valid/%s", documentNumber)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to call document validation API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("document validation API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	var validationResponse entities.DocumentValidationResponse
	if err := json.Unmarshal(body, &validationResponse); err != nil {
		return false, fmt.Errorf("failed to unmarshal validation response: %v", err)
	}

	return validationResponse.IsValid, nil
}

func (s *ExternalServicesImpl) GetAddress(zipCode string) (*entities.AddressResponse, error) {
	// Converter string para int para a URL
	zipCodeInt, err := strconv.Atoi(zipCode)
	if err != nil {
		return nil, fmt.Errorf("invalid zip code: %v", err)
	}

	url := fmt.Sprintf("http://localhost:8081/api/v1/address/%d", zipCodeInt)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call address API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("address API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var addressResponse entities.AddressResponse
	if err := json.Unmarshal(body, &addressResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal address response: %v", err)
	}

	return &addressResponse, nil
}
