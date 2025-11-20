package auth

import (
	"context"
	"errors"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"github.com/golang-jwt/jwt/v5"
)

type LoginUseCase struct {
	userRepo  ports.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewLoginUseCase(userRepo ports.UserRepository, jwtSecret string, jwtExpiry time.Duration) *LoginUseCase {
	return &LoginUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (string, *entities.User, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", nil, errors.New("user is inactive")
	}

	if !user.CheckPassword(password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := uc.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (uc *LoginUseCase) generateToken(user *entities.User) (string, error) {
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
