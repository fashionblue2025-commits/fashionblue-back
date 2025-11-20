package user

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// ChangePasswordUseCase maneja el cambio de contraseña
type ChangePasswordUseCase struct {
	userRepo ports.UserRepository
}

// NewChangePasswordUseCase crea una nueva instancia del caso de uso
func NewChangePasswordUseCase(userRepo ports.UserRepository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo: userRepo,
	}
}

// Execute ejecuta el caso de uso de cambiar contraseña
func (uc *ChangePasswordUseCase) Execute(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.CheckPassword(oldPassword) {
		return errors.New("invalid old password")
	}

	if err := user.HashPassword(newPassword); err != nil {
		return err
	}

	return uc.userRepo.Update(ctx, user)
}
