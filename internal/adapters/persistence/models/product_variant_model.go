package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// ProductVariantModel representa el modelo de persistencia para variantes de producto
type ProductVariantModel struct {
	ID            uint          `gorm:"primaryKey"`
	ProductID     uint          `gorm:"not null;index"`
	Product       *ProductModel `gorm:"foreignKey:ProductID"`
	Color         string        `gorm:"type:varchar(50);not null;index"`
	SizeID        *uint         `gorm:"index"`
	Size          *SizeModel    `gorm:"foreignKey:SizeID"`
	Stock         int           `gorm:"default:0;not null"`
	ReservedStock int           `gorm:"default:0;not null"`
	UnitPrice     float64       `gorm:"not null"`
	IsActive      bool          `gorm:"default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName especifica el nombre de la tabla
func (ProductVariantModel) TableName() string {
	return "product_variants"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *ProductVariantModel) ToEntity() *entities.ProductVariant {
	variant := &entities.ProductVariant{
		ID:            m.ID,
		ProductID:     m.ProductID,
		Color:         m.Color,
		SizeID:        m.SizeID,
		Stock:         m.Stock,
		ReservedStock: m.ReservedStock,
		UnitPrice:     m.UnitPrice,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}

	// Convertir Product si existe
	if m.Product != nil {
		variant.Product = m.Product.ToEntity()
	}

	// Convertir Size si existe
	if m.Size != nil {
		variant.Size = m.Size.ToEntity()
	}

	return variant
}

// FromEntity convierte una entidad de dominio a modelo
func (m *ProductVariantModel) FromEntity(variant *entities.ProductVariant) {
	m.ID = variant.ID
	m.ProductID = variant.ProductID
	m.Color = variant.Color
	m.SizeID = variant.SizeID
	m.Stock = variant.Stock
	m.ReservedStock = variant.ReservedStock
	m.UnitPrice = variant.UnitPrice
	m.IsActive = variant.IsActive
	m.CreatedAt = variant.CreatedAt
	m.UpdatedAt = variant.UpdatedAt
}
