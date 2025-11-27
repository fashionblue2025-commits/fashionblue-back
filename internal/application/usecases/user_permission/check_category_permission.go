package userpermission

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CheckCategoryPermissionUseCase struct {
	permissionRepo ports.UserCategoryPermissionRepository
	userRepo       ports.UserRepository
}

func NewCheckCategoryPermissionUseCase(
	permissionRepo ports.UserCategoryPermissionRepository,
	userRepo ports.UserRepository,
) *CheckCategoryPermissionUseCase {
	return &CheckCategoryPermissionUseCase{
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

// Execute verifica si un usuario tiene permiso sobre una categoría
// Retorna true si el usuario es admin o tiene el permiso específico
func (uc *CheckCategoryPermissionUseCase) Execute(ctx context.Context, userID, categoryID uint, action string) (bool, error) {
	// Obtener usuario para verificar rol
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Super Admin tiene todos los permisos
	if user.Role == entities.RoleSuperAdmin {
		return true, nil
	}

	// Verificar permiso específico
	hasPermission, err := uc.permissionRepo.HasPermission(ctx, userID, categoryID, action)
	if err != nil {
		return false, err
	}

	return hasPermission, nil
}
