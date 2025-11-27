package userpermission

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ManageUserPermissionsUseCase struct {
	permissionRepo ports.UserCategoryPermissionRepository
	categoryRepo   ports.CategoryRepository
	userRepo       ports.UserRepository
}

func NewManageUserPermissionsUseCase(
	permissionRepo ports.UserCategoryPermissionRepository,
	categoryRepo ports.CategoryRepository,
	userRepo ports.UserRepository,
) *ManageUserPermissionsUseCase {
	return &ManageUserPermissionsUseCase{
		permissionRepo: permissionRepo,
		categoryRepo:   categoryRepo,
		userRepo:       userRepo,
	}
}

// SetUserPermissions establece los permisos de un usuario (reemplaza todos los existentes)
func (uc *ManageUserPermissionsUseCase) SetUserPermissions(ctx context.Context, userID uint, permissions []entities.UserCategoryPermission) error {
	// Verificar que el usuario existe
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verificar que todas las categorías existen
	for _, perm := range permissions {
		_, err := uc.categoryRepo.GetByID(ctx, perm.CategoryID)
		if err != nil {
			return errors.New("category not found: " + string(rune(perm.CategoryID)))
		}
	}

	return uc.permissionRepo.SetPermissions(ctx, userID, permissions)
}

// AddCategoryPermission agrega o actualiza un permiso específico
func (uc *ManageUserPermissionsUseCase) AddCategoryPermission(ctx context.Context, permission *entities.UserCategoryPermission) error {
	// Verificar usuario
	_, err := uc.userRepo.GetByID(ctx, permission.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verificar categoría
	_, err = uc.categoryRepo.GetByID(ctx, permission.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	// Verificar si ya existe
	existing, err := uc.permissionRepo.GetByUserAndCategory(ctx, permission.UserID, permission.CategoryID)
	if err != nil {
		return err
	}

	if existing != nil {
		// Actualizar existente
		permission.ID = existing.ID
		return uc.permissionRepo.Update(ctx, permission)
	}

	// Crear nuevo
	return uc.permissionRepo.Create(ctx, permission)
}

// RemoveCategoryPermission elimina un permiso específico
func (uc *ManageUserPermissionsUseCase) RemoveCategoryPermission(ctx context.Context, userID, categoryID uint) error {
	existing, err := uc.permissionRepo.GetByUserAndCategory(ctx, userID, categoryID)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("permission not found")
	}

	return uc.permissionRepo.Delete(ctx, existing.ID)
}

// GetUserPermissions obtiene todos los permisos de un usuario
func (uc *ManageUserPermissionsUseCase) GetUserPermissions(ctx context.Context, userID uint) ([]entities.UserCategoryPermission, error) {
	// Verificar usuario
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return uc.permissionRepo.ListByUser(ctx, userID)
}

// GetCategoryPermissions obtiene todos los usuarios con permisos sobre una categoría
func (uc *ManageUserPermissionsUseCase) GetCategoryPermissions(ctx context.Context, categoryID uint) ([]entities.UserCategoryPermission, error) {
	// Verificar categoría
	_, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	return uc.permissionRepo.ListByCategory(ctx, categoryID)
}
