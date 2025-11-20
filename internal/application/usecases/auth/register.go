package auth

import (
	"context"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterUseCase struct {
	userRepo  ports.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewRegisterUseCase(userRepo ports.UserRepository, jwtSecret string, jwtExpiry time.Duration) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, user *entities.User, password string) (string, error) {
	if err := user.HashPassword(password); err != nil {
		return "", err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	token, err := uc.generateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *RegisterUseCase) generateToken(user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    string(user.Role),
		"exp":     time.Now().Add(uc.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
