package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// PaymentMethodModel representa el modelo de persistencia para métodos de pago
type PaymentMethodModel struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName especifica el nombre de la tabla
func (PaymentMethodModel) TableName() string {
	return "payment_methods"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *PaymentMethodModel) ToEntity() *entities.PaymentMethodOption {
	return &entities.PaymentMethodOption{
		ID:        m.ID,
		Name:      m.Name,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo
func (m *PaymentMethodModel) FromEntity(pm *entities.PaymentMethodOption) {
	m.ID = pm.ID
	m.Name = pm.Name
	m.IsActive = pm.IsActive
	m.CreatedAt = pm.CreatedAt
	m.UpdatedAt = pm.UpdatedAt
}
