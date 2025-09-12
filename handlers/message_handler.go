package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"api-rabbitmq/models"
)

type MessageHandler struct {
	httpClient *http.Client
	repo       *models.MongoDBRepository
}

func NewMessageHandler(repo *models.MongoDBRepository) *MessageHandler {
	return &MessageHandler{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		repo: repo,
	}
}

func (h *MessageHandler) ProcessMessage(userData models.UserData) (*models.ProcessedData, error) {
	log.Printf("Processing message for user: %s, document: %s, zipCode: %s",
		userData.Name, userData.DocumentNumber, userData.ZipCode)

	// Validate document number
	isValid, err := h.validateDocument(userData.DocumentNumber)
	if err != nil {
		return nil, fmt.Errorf("document validation failed: %v", err)
	}

	log.Printf("Document validation result for %s: %t", userData.DocumentNumber, isValid)

	// Get address from zip code
	address, err := h.getAddress(userData.ZipCode)
	if err != nil {
		return nil, fmt.Errorf("address lookup failed: %v", err)
	}

	log.Printf("Address found for zipCode %s: %s, %s - %s",
		userData.ZipCode, address.Street, address.City, address.State)

	processedData := &models.ProcessedData{
		UserData:      userData,
		DocumentValid: isValid,
		Address:       *address,
		Status:        "processed",
		Message:       "Message processed successfully",
	}

	// Salvar no MongoDB apenas se o repositório estiver disponível
	if h.repo != nil {
		id, err := h.repo.SaveProcessedData(processedData)
		if err != nil {
			log.Printf("Error saving to MongoDB: %v", err)
			// Não retornamos erro aqui, apenas logamos, pois o processamento principal foi bem-sucedido
		} else {
			log.Printf("Data saved to MongoDB with ID: %s", id)
		}
	} else {
		log.Println("MongoDB repository not available, skipping save operation")
	}

	return processedData, nil
}

func (h *MessageHandler) validateDocument(documentNumber string) (bool, error) {
	url := fmt.Sprintf("http://localhost:8082/api/v1/is-document-valid/%s", documentNumber)

	log.Printf("Calling document validation API: %s", url)

	resp, err := h.httpClient.Get(url)
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

	var validationResponse models.DocumentValidationResponse
	if err := json.Unmarshal(body, &validationResponse); err != nil {
		return false, fmt.Errorf("failed to unmarshal validation response: %v", err)
	}

	return validationResponse.IsValid, nil
}

func (h *MessageHandler) getAddress(zipCode string) (*models.AddressResponse, error) {
	url := fmt.Sprintf("http://localhost:8081/api/v1/address/%s", zipCode)

	log.Printf("Calling address API: %s", url)

	resp, err := h.httpClient.Get(url)
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

	var addressResponse models.AddressResponse
	if err := json.Unmarshal(body, &addressResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal address response: %v", err)
	}

	return &addressResponse, nil
}
