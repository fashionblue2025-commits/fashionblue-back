package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// UserCategoryPermissionModel representa el modelo de persistencia para permisos
type UserCategoryPermissionModel struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"not null;index"`
	CategoryID uint `gorm:"not null;index"`
	CanView    bool `gorm:"default:true"`
	CanCreate  bool `gorm:"default:false"`
	CanEdit    bool `gorm:"default:false"`
	CanDelete  bool `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Relaciones
	User     UserModel     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Category CategoryModel `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}

// TableName especifica el nombre de la tabla
func (UserCategoryPermissionModel) TableName() string {
	return "user_category_permissions"
}

// ToEntity convierte el modelo a entidad
func (m *UserCategoryPermissionModel) ToEntity() *entities.UserCategoryPermission {
	return &entities.UserCategoryPermission{
		ID:         m.ID,
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		CanView:    m.CanView,
		CanCreate:  m.CanCreate,
		CanEdit:    m.CanEdit,
		CanDelete:  m.CanDelete,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// FromEntity convierte la entidad a modelo
func (m *UserCategoryPermissionModel) FromEntity(permission *entities.UserCategoryPermission) {
	m.ID = permission.ID
	m.UserID = permission.UserID
	m.CategoryID = permission.CategoryID
	m.CanView = permission.CanView
	m.CanCreate = permission.CanCreate
	m.CanEdit = permission.CanEdit
	m.CanDelete = permission.CanDelete
	m.CreatedAt = permission.CreatedAt
	m.UpdatedAt = permission.UpdatedAt
}
