package userpermission

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetUserAllowedCategoriesUseCase struct {
	permissionRepo ports.UserCategoryPermissionRepository
	categoryRepo   ports.CategoryRepository
	userRepo       ports.UserRepository
}

func NewGetUserAllowedCategoriesUseCase(
	permissionRepo ports.UserCategoryPermissionRepository,
	categoryRepo ports.CategoryRepository,
	userRepo ports.UserRepository,
) *GetUserAllowedCategoriesUseCase {
	return &GetUserAllowedCategoriesUseCase{
		permissionRepo: permissionRepo,
		categoryRepo:   categoryRepo,
		userRepo:       userRepo,
	}
}

// Execute obtiene todas las categorías a las que un usuario tiene acceso
func (uc *GetUserAllowedCategoriesUseCase) Execute(ctx context.Context, userID uint, action string) ([]entities.Category, error) {
	// Obtener usuario para verificar rol
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Super Admin puede ver todas las categorías
	if user.Role == entities.RoleSuperAdmin {
		return uc.categoryRepo.List(ctx, map[string]interface{}{})
	}

	// Obtener IDs de categorías permitidas
	categoryIDs, err := uc.permissionRepo.GetAllowedCategoriesForUser(ctx, userID, action)
	if err != nil {
		return nil, err
	}

	// Si no tiene permisos, retornar lista vacía
	if len(categoryIDs) == 0 {
		return []entities.Category{}, nil
	}

	// Obtener las categorías completas
	filters := map[string]interface{}{
		"ids": categoryIDs,
	}

	return uc.categoryRepo.List(ctx, filters)
}
