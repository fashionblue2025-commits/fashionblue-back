package user

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

// CreateUserUseCase maneja la creación de usuarios
type CreateUserUseCase struct {
	userRepo ports.UserRepository
}

// NewCreateUserUseCase crea una nueva instancia del caso de uso
func NewCreateUserUseCase(userRepo ports.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
	}
}

// Execute ejecuta el caso de uso de crear usuario
func (uc *CreateUserUseCase) Execute(ctx context.Context, user *entities.User, password string) error {
	// Verificar si el email ya existe
	existingUser, err := uc.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Encriptar contraseña
	if err := user.HashPassword(password); err != nil {
		return err
	}

	// Crear usuario
	return uc.userRepo.Create(ctx, user)
}
