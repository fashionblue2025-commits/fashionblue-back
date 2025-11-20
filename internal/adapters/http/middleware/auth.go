package middleware

import (
	"strings"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/auth"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware middleware de autenticación
func AuthMiddleware(validateTokenUC *auth.ValidateTokenUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Unauthorized(c, "Missing authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Unauthorized(c, "Invalid authorization header format")
			}

			token := parts[1]
			user, err := validateTokenUC.Execute(c.Request().Context(), token)
			if err != nil {
				return response.Unauthorized(c, "Invalid or expired token")
			}

			// Guardar usuario en el contexto
			c.Set("user", user)
			return next(c)
		}
	}
}

// RequireRole middleware que requiere un rol específico
func RequireRole(roles ...entities.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*entities.User)
			if !ok {
				return response.Unauthorized(c, "User not found in context")
			}

			for _, role := range roles {
				if user.Role == role {
					return next(c)
				}
			}

			return response.Forbidden(c, "Insufficient permissions")
		}
	}
}

// GetUserFromContext obtiene el usuario del contexto
func GetUserFromContext(c echo.Context) (*entities.User, error) {
	user, ok := c.Get("user").(*entities.User)
	if !ok {
		return nil, echo.NewHTTPError(401, "User not found in context")
	}
	return user, nil
}
