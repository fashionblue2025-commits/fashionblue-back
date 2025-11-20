package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// SupplierModel representa el modelo de persistencia para proveedores
type SupplierModel struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(255);not null"`
	ContactName string `gorm:"type:varchar(255)"`
	Phone       string `gorm:"type:varchar(50)"`
	Email       string `gorm:"type:varchar(255)"`
	Address     string `gorm:"type:text"`
	Notes       string `gorm:"type:text"`
	IsActive    bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName especifica el nombre de la tabla
func (SupplierModel) TableName() string {
	return "suppliers"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *SupplierModel) ToEntity() *entities.Supplier {
	return &entities.Supplier{
		ID:          m.ID,
		Name:        m.Name,
		ContactName: m.ContactName,
		Phone:       m.Phone,
		Email:       m.Email,
		Address:     m.Address,
		Notes:       m.Notes,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo
func (m *SupplierModel) FromEntity(supplier *entities.Supplier) {
	m.ID = supplier.ID
	m.Name = supplier.Name
	m.ContactName = supplier.ContactName
	m.Phone = supplier.Phone
	m.Email = supplier.Email
	m.Address = supplier.Address
	m.Notes = supplier.Notes
	m.IsActive = supplier.IsActive
	m.CreatedAt = supplier.CreatedAt
	m.UpdatedAt = supplier.UpdatedAt
}
