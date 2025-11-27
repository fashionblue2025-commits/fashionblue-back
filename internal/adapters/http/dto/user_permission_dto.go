package dto

import "github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"

// UserCategoryPermissionDTO representa los permisos de usuario en la respuesta
type UserCategoryPermissionDTO struct {
	ID           uint         `json:"id"`
	UserID       uint         `json:"user_id"`
	CategoryID   uint         `json:"category_id"`
	CategoryName string       `json:"category_name,omitempty"`
	CanView      bool         `json:"can_view"`
	CanCreate    bool         `json:"can_create"`
	CanEdit      bool         `json:"can_edit"`
	CanDelete    bool         `json:"can_delete"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
	Category     *CategoryDTO `json:"category,omitempty"`
}

// CreateUserCategoryPermissionRequest para crear un permiso
type CreateUserCategoryPermissionRequest struct {
	UserID     uint `json:"user_id" validate:"required"`
	CategoryID uint `json:"category_id" validate:"required"`
	CanView    bool `json:"can_view"`
	CanCreate  bool `json:"can_create"`
	CanEdit    bool `json:"can_edit"`
	CanDelete  bool `json:"can_delete"`
}

// UpdateUserCategoryPermissionRequest para actualizar un permiso
type UpdateUserCategoryPermissionRequest struct {
	CanView   bool `json:"can_view"`
	CanCreate bool `json:"can_create"`
	CanEdit   bool `json:"can_edit"`
	CanDelete bool `json:"can_delete"`
}

// SetUserPermissionsRequest para establecer todos los permisos de un usuario
type SetUserPermissionsRequest struct {
	Permissions []PermissionItem `json:"permissions" validate:"required"`
}

type PermissionItem struct {
	CategoryID uint `json:"category_id" validate:"required"`
	CanView    bool `json:"can_view"`
	CanCreate  bool `json:"can_create"`
	CanEdit    bool `json:"can_edit"`
	CanDelete  bool `json:"can_delete"`
}

// ToUserCategoryPermissionDTO convierte la entidad a DTO
func ToUserCategoryPermissionDTO(permission *entities.UserCategoryPermission) UserCategoryPermissionDTO {
	dto := UserCategoryPermissionDTO{
		ID:         permission.ID,
		UserID:     permission.UserID,
		CategoryID: permission.CategoryID,
		CanView:    permission.CanView,
		CanCreate:  permission.CanCreate,
		CanEdit:    permission.CanEdit,
		CanDelete:  permission.CanDelete,
		CreatedAt:  permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if permission.Category != nil {
		dto.Category = ToCategoryDTO(permission.Category)
		dto.CategoryName = permission.Category.Name
	}

	return dto
}

// ToUserCategoryPermissionDTOList convierte lista de entidades a DTOs
func ToUserCategoryPermissionDTOList(permissions []entities.UserCategoryPermission) []UserCategoryPermissionDTO {
	dtos := make([]UserCategoryPermissionDTO, len(permissions))
	for i, perm := range permissions {
		dtos[i] = ToUserCategoryPermissionDTO(&perm)
	}
	return dtos
}
