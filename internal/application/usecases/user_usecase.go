package usecases

import (
	"context"
	"fmt"
	"log"

	"api-rabbitmq/internal/domain/entities"
	"api-rabbitmq/internal/domain/repositories"
)

// UserUseCase define os casos de uso para usuários
type UserUseCase interface {
	ProcessUser(ctx context.Context, userData entities.UserData) (*entities.ProcessedUser, error)
	GetProcessedUsers(ctx context.Context) ([]entities.ProcessedUser, error)
}

type userUseCase struct {
	userRepo    repositories.UserRepository
	extServices ExternalServices
}

// ExternalServices define as dependências externas
type ExternalServices interface {
	ValidateDocument(documentNumber string) (bool, error)
	GetAddress(zipCode string) (*entities.AddressResponse, error)
}

func NewUserUseCase(userRepo repositories.UserRepository, extServices ExternalServices) UserUseCase {
	return &userUseCase{
		userRepo:    userRepo,
		extServices: extServices,
	}
}

func (uc *userUseCase) ProcessUser(ctx context.Context, userData entities.UserData) (*entities.ProcessedUser, error) {
	log.Printf("Processing user: %s, document: %s", userData.Name, userData.DocumentNumber)

	// Validar documento
	isValid, err := uc.extServices.ValidateDocument(userData.DocumentNumber)
	if err != nil {
		return nil, fmt.Errorf("document validation failed: %v", err)
	}

	// Buscar endereço
	address, err := uc.extServices.GetAddress(userData.ZipCode)
	if err != nil {
		return nil, fmt.Errorf("address lookup failed: %v", err)
	}

	processedUser := &entities.ProcessedUser{
		UserData:      userData,
		DocumentValid: isValid,
		Address:       *address,
		Status:        "processed",
		Message:       "User processed successfully",
	}

	// Salvar no repositório
	if uc.userRepo != nil {
		_, err := uc.userRepo.Save(ctx, processedUser)
		if err != nil {
			log.Printf("Warning: failed to save user: %v", err)
		}
	}

	return processedUser, nil
}

func (uc *userUseCase) GetProcessedUsers(ctx context.Context) ([]entities.ProcessedUser, error) {
	if uc.userRepo == nil {
		return nil, fmt.Errorf("user repository not available")
	}
	return uc.userRepo.FindAll(ctx)
}
