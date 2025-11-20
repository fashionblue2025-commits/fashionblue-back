package user

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// UpdateUserUseCase maneja la actualización de usuarios
type UpdateUserUseCase struct {
	userRepo ports.UserRepository
}

// NewUpdateUserUseCase crea una nueva instancia del caso de uso
func NewUpdateUserUseCase(userRepo ports.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: userRepo,
	}
}

// Execute ejecuta el caso de uso de actualizar usuario
func (uc *UpdateUserUseCase) Execute(ctx context.Context, user *entities.User) error {
	// Verificar que el usuario existe
	_, err := uc.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}
