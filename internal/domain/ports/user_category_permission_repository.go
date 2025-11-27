package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// UserCategoryPermissionRepository define las operaciones para permisos de usuario por categoría
type UserCategoryPermissionRepository interface {
	// Create crea un nuevo permiso
	Create(ctx context.Context, permission *entities.UserCategoryPermission) error

	// Update actualiza un permiso existente
	Update(ctx context.Context, permission *entities.UserCategoryPermission) error

	// Delete elimina un permiso
	Delete(ctx context.Context, id uint) error

	// GetByID obtiene un permiso por su ID
	GetByID(ctx context.Context, id uint) (*entities.UserCategoryPermission, error)

	// GetByUserAndCategory obtiene el permiso específico de un usuario sobre una categoría
	GetByUserAndCategory(ctx context.Context, userID, categoryID uint) (*entities.UserCategoryPermission, error)

	// ListByUser obtiene todos los permisos de un usuario
	ListByUser(ctx context.Context, userID uint) ([]entities.UserCategoryPermission, error)

	// ListByCategory obtiene todos los usuarios con permisos sobre una categoría
	ListByCategory(ctx context.Context, categoryID uint) ([]entities.UserCategoryPermission, error)

	// GetAllowedCategoriesForUser obtiene IDs de categorías permitidas para un usuario
	GetAllowedCategoriesForUser(ctx context.Context, userID uint, action string) ([]uint, error)

	// HasPermission verifica si un usuario tiene un permiso específico sobre una categoría
	HasPermission(ctx context.Context, userID, categoryID uint, action string) (bool, error)

	// SetPermissions establece múltiples permisos para un usuario (reemplaza existentes)
	SetPermissions(ctx context.Context, userID uint, permissions []entities.UserCategoryPermission) error

	// DeleteByUser elimina todos los permisos de un usuario
	DeleteByUser(ctx context.Context, userID uint) error

	// DeleteByCategory elimina todos los permisos de una categoría
	DeleteByCategory(ctx context.Context, categoryID uint) error
}
