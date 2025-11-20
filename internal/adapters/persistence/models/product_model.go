package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// ProductModel representa el modelo de persistencia para productos
type ProductModel struct {
	ID              uint           `gorm:"primaryKey"`
	Name            string         `gorm:"not null"`
	Description     string         `gorm:"type:text"`
	SKU             string         `gorm:"uniqueIndex;not null"`
	CategoryID      uint           `gorm:"not null"`
	Category        *CategoryModel `gorm:"foreignKey:CategoryID"`
	MaterialCost    float64        `gorm:"not null"`
	LaborCost       float64        `gorm:"not null"`
	ProductionCost  float64        `gorm:"not null"`
	UnitPrice       float64        `gorm:"not null"`
	WholesalePrice  float64        `gorm:"not null"`
	MinWholesaleQty int            `gorm:"default:10"`
	Stock           int            `gorm:"default:0"`
	MinStock        int            `gorm:"default:5"`
	IsActive        bool           `gorm:"default:true"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TableName especifica el nombre de la tabla
func (ProductModel) TableName() string {
	return "products"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *ProductModel) ToEntity() *entities.Product {
	return &entities.Product{
		ID:              m.ID,
		Name:            m.Name,
		Description:     m.Description,
		SKU:             m.SKU,
		CategoryID:      m.CategoryID,
		MaterialCost:    m.MaterialCost,
		LaborCost:       m.LaborCost,
		ProductionCost:  m.ProductionCost,
		UnitPrice:       m.UnitPrice,
		WholesalePrice:  m.WholesalePrice,
		MinWholesaleQty: m.MinWholesaleQty,
		Stock:           m.Stock,
		MinStock:        m.MinStock,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo
func (m *ProductModel) FromEntity(product *entities.Product) {
	m.ID = product.ID
	m.Name = product.Name
	m.Description = product.Description
	m.SKU = product.SKU
	m.CategoryID = product.CategoryID
	m.MaterialCost = product.MaterialCost
	m.LaborCost = product.LaborCost
	m.ProductionCost = product.ProductionCost
	m.UnitPrice = product.UnitPrice
	m.WholesalePrice = product.WholesalePrice
	m.MinWholesaleQty = product.MinWholesaleQty
	m.Stock = product.Stock
	m.MinStock = product.MinStock
	m.IsActive = product.IsActive
	m.CreatedAt = product.CreatedAt
	m.UpdatedAt = product.UpdatedAt
}

// CategoryModel representa el modelo de persistencia para categorías
type CategoryModel struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	Description string `gorm:"type:text"`
	IsActive    bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName especifica el nombre de la tabla
func (CategoryModel) TableName() string {
	return "categories"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *CategoryModel) ToEntity() *entities.Category {
	return &entities.Category{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo
func (m *CategoryModel) FromEntity(category *entities.Category) {
	m.ID = category.ID
	m.Name = category.Name
	m.Description = category.Description
	m.IsActive = category.IsActive
	m.CreatedAt = category.CreatedAt
	m.UpdatedAt = category.UpdatedAt
}
