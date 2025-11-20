package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// UserModel representa el modelo de persistencia para usuarios
type UserModel struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Role      string `gorm:"type:varchar(50);not null"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName especifica el nombre de la tabla
func (UserModel) TableName() string {
	return "users"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *UserModel) ToEntity() *entities.User {
	return &entities.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Role:      entities.UserRole(m.Role),
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo
func (m *UserModel) FromEntity(user *entities.User) {
	m.ID = user.ID
	m.Email = user.Email
	m.Password = user.Password
	m.FirstName = user.FirstName
	m.LastName = user.LastName
	m.Role = string(user.Role)
	m.IsActive = user.IsActive
	m.CreatedAt = user.CreatedAt
	m.UpdatedAt = user.UpdatedAt
}
