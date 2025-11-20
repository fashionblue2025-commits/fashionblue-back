package user

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// ListUsersUseCase maneja el listado de usuarios
type ListUsersUseCase struct {
	userRepo ports.UserRepository
}

// NewListUsersUseCase crea una nueva instancia del caso de uso
func NewListUsersUseCase(userRepo ports.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

// Execute ejecuta el caso de uso de listar usuarios
func (uc *ListUsersUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]entities.User, error) {
	return uc.userRepo.List(ctx, filters)
}
