package user

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// GetUserUseCase maneja la obtención de un usuario por ID
type GetUserUseCase struct {
	userRepo ports.UserRepository
}

// NewGetUserUseCase crea una nueva instancia del caso de uso
func NewGetUserUseCase(userRepo ports.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

// Execute ejecuta el caso de uso de obtener usuario
func (uc *GetUserUseCase) Execute(ctx context.Context, id uint) (*entities.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}
