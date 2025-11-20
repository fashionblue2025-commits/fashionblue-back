package auth

import (
	"net/http"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/auth"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	loginUC    *auth.LoginUseCase
	registerUC *auth.RegisterUseCase
}

func NewAuthHandler(loginUC *auth.LoginUseCase, registerUC *auth.RegisterUseCase) *AuthHandler {
	return &AuthHandler{
		loginUC:    loginUC,
		registerUC: registerUC,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email     string            `json:"email" validate:"required,email"`
	Password  string            `json:"password" validate:"required,min=6"`
	FirstName string            `json:"first_name" validate:"required"`
	LastName  string            `json:"last_name" validate:"required"`
	Role      entities.UserRole `json:"role" validate:"required"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	token, user, err := h.loginUC.Execute(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.OK(c, "Login successful", map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	user := &entities.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	token, err := h.registerUC.Execute(c.Request().Context(), user, req.Password)
	if err != nil {
		return response.BadRequest(c, "Failed to register user", err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User registered successfully",
		"data": map[string]interface{}{
			"token": token,
			"user":  user,
		},
	})
}
