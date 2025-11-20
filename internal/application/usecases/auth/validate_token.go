package auth

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"github.com/golang-jwt/jwt/v5"
)

type ValidateTokenUseCase struct {
	userRepo  ports.UserRepository
	jwtSecret string
}

func NewValidateTokenUseCase(userRepo ports.UserRepository, jwtSecret string) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (uc *ValidateTokenUseCase) Execute(ctx context.Context, tokenString string) (*entities.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(uc.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	user, err := uc.userRepo.GetByID(ctx, uint(userID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	return user, nil
}
