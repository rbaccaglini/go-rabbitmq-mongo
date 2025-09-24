package usecases

import (
	"context"
	"fmt"
	"log"
	"sync"

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

	var wg sync.WaitGroup
	var isValid bool
	var errAddr, errDoc error
	var address *entities.AddressResponse

	// Validar documento
	wg.Add(1)
	go func() {
		isValid, errDoc = uc.extServices.ValidateDocument(userData.DocumentNumber)
	}()

	// Buscar endereço
	wg.Add(1)
	go func() {
		address, errAddr = uc.extServices.GetAddress(userData.ZipCode)
	}()

	wg.Wait()

	if errDoc != nil {
		return nil, fmt.Errorf("document validation failed: %v", errDoc)
	}

	if errAddr != nil {
		return nil, fmt.Errorf("address lookup failed: %v", errAddr)
	}

	processedUser := &entities.ProcessedUser{
		Name:     userData.Name,
		Document: entities.DocumentUserProcessed{DocumentNumber: userData.DocumentNumber, IsValid: isValid},
		Address:  *address,
		Status:   "processed",
		Message:  "User processed successfully",
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
