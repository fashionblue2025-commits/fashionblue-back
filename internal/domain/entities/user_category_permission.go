package entities

import "time"

// UserCategoryPermission representa los permisos de un usuario sobre una categoría
type UserCategoryPermission struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	CategoryID uint      `json:"category_id"`
	CanView    bool      `json:"can_view"`
	CanCreate  bool      `json:"can_create"`
	CanEdit    bool      `json:"can_edit"`
	CanDelete  bool      `json:"can_delete"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relaciones (opcionales)
	User     *User     `json:"user,omitempty"`
	Category *Category `json:"category,omitempty"`
}

// HasPermission verifica si tiene un permiso específico
func (p *UserCategoryPermission) HasPermission(action string) bool {
	switch action {
	case "view":
		return p.CanView
	case "create":
		return p.CanCreate
	case "edit":
		return p.CanEdit
	case "delete":
		return p.CanDelete
	default:
		return false
	}
}

// PermissionLevel devuelve el nivel de permiso (0-4)
func (p *UserCategoryPermission) PermissionLevel() int {
	level := 0
	if p.CanView {
		level++
	}
	if p.CanCreate {
		level++
	}
	if p.CanEdit {
		level++
	}
	if p.CanDelete {
		level++
	}
	return level
}

// IsReadOnly verifica si solo tiene permisos de lectura
func (p *UserCategoryPermission) IsReadOnly() bool {
	return p.CanView && !p.CanCreate && !p.CanEdit && !p.CanDelete
}

// IsFullAccess verifica si tiene acceso completo
func (p *UserCategoryPermission) IsFullAccess() bool {
	return p.CanView && p.CanCreate && p.CanEdit && p.CanDelete
}
