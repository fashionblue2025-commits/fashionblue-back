package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"gorm.io/gorm"
)

// SizeModel es el modelo GORM para tallas
type SizeModel struct {
	ID        uint              `gorm:"primaryKey"`
	Type      entities.SizeType `gorm:"type:varchar(20);not null;index"`
	Value     string            `gorm:"type:varchar(10);not null"`
	Order     int               `gorm:"not null;default:0"`
	IsActive  bool              `gorm:"not null;default:true"`
	CreatedAt time.Time         `gorm:"autoCreateTime"`
	UpdatedAt time.Time         `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt    `gorm:"index"`
}

// TableName especifica el nombre de la tabla
func (SizeModel) TableName() string {
	return "sizes"
}

// ToEntity convierte el modelo GORM a entidad de dominio
func (m *SizeModel) ToEntity() *entities.Size {
	if m == nil {
		return nil
	}
	return &entities.Size{
		ID:        m.ID,
		Type:      m.Type,
		Value:     m.Value,
		Order:     m.Order,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo GORM
func (m *SizeModel) FromEntity(size *entities.Size) {
	if size == nil {
		return
	}
	m.ID = size.ID
	m.Type = size.Type
	m.Value = size.Value
	m.Order = size.Order
	m.IsActive = size.IsActive
	m.CreatedAt = size.CreatedAt
	m.UpdatedAt = size.UpdatedAt
}
