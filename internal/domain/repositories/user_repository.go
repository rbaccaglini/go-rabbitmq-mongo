package repositories

import (
	"context"

	"api-rabbitmq/internal/domain/entities"
)

// UserRepository define a interface para operações de banco de dados
type UserRepository interface {
	Save(ctx context.Context, user *entities.ProcessedUser) (string, error)
	FindAll(ctx context.Context) ([]entities.ProcessedUser, error)
	FindByID(ctx context.Context, id string) (*entities.ProcessedUser, error)
	Close() error
}
