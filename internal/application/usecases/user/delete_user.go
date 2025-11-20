package user

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// DeleteUserUseCase maneja la eliminación de usuarios
type DeleteUserUseCase struct {
	userRepo ports.UserRepository
}

// NewDeleteUserUseCase crea una nueva instancia del caso de uso
func NewDeleteUserUseCase(userRepo ports.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepo: userRepo,
	}
}

// Execute ejecuta el caso de uso de eliminar usuario
func (uc *DeleteUserUseCase) Execute(ctx context.Context, id uint) error {
	// Verificar que el usuario existe
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.userRepo.Delete(ctx, id)
}
