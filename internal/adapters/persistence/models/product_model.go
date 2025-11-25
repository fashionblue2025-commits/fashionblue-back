package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// ProductModel representa el modelo de persistencia para productos base
type ProductModel struct {
	ID              uint                  `gorm:"primaryKey"`
	Name            string                `gorm:"not null;index"`
	Description     string                `gorm:"type:text"`
	CategoryID      uint                  `gorm:"not null"`
	Category        *CategoryModel        `gorm:"foreignKey:CategoryID"`
	MaterialCost    float64               `gorm:"not null;default:0"`
	LaborCost       float64               `gorm:"not null;default:0"`
	ProductionCost  float64               `gorm:"not null;default:0"`
	UnitPrice       float64               `gorm:"not null"`
	WholesalePrice  float64               `gorm:"not null"`
	MinWholesaleQty int                   `gorm:"default:10"`
	MinStock        int                   `gorm:"default:5"`
	IsActive        bool                  `gorm:"default:true"`
	Variants        []ProductVariantModel `gorm:"foreignKey:ProductID"` // Variantes de este producto
	// Campos de origen
	OriginType    string `gorm:"type:varchar(20);default:'INVENTORY';index"` // CUSTOM, INVENTORY
	OriginOrderID *uint  `gorm:"index"`                                      // ID de la orden que lo creó
	IsCustom      bool   `gorm:"default:false;index"`                        // Producto personalizado
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName especifica el nombre de la tabla
func (ProductModel) TableName() string {
	return "products"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *ProductModel) ToEntity() *entities.Product {
	product := &entities.Product{
		ID:              m.ID,
		Name:            m.Name,
		Description:     m.Description,
		CategoryID:      m.CategoryID,
		MaterialCost:    m.MaterialCost,
		LaborCost:       m.LaborCost,
		ProductionCost:  m.ProductionCost,
		UnitPrice:       m.UnitPrice,
		WholesalePrice:  m.WholesalePrice,
		MinWholesaleQty: m.MinWholesaleQty,
		MinStock:        m.MinStock,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}

	// Convertir Category si existe
	if m.Category != nil {
		product.Category = m.Category.ToEntity()
	}

	// Convertir Variants si existen
	if len(m.Variants) > 0 {
		product.Variants = make([]entities.ProductVariant, len(m.Variants))
		for i, variantModel := range m.Variants {
			product.Variants[i] = *variantModel.ToEntity()
		}
	}

	return product
}

// FromEntity convierte una entidad de dominio a modelo
func (m *ProductModel) FromEntity(product *entities.Product) {
	m.ID = product.ID
	m.Name = product.Name
	m.Description = product.Description
	m.CategoryID = product.CategoryID
	m.MaterialCost = product.MaterialCost
	m.LaborCost = product.LaborCost
	m.ProductionCost = product.ProductionCost
	m.UnitPrice = product.UnitPrice
	m.WholesalePrice = product.WholesalePrice
	m.MinWholesaleQty = product.MinWholesaleQty
	m.MinStock = product.MinStock
	m.IsActive = product.IsActive
	m.CreatedAt = product.CreatedAt
	m.UpdatedAt = product.UpdatedAt

	// Convertir Variants si existen
	if len(product.Variants) > 0 {
		m.Variants = make([]ProductVariantModel, len(product.Variants))
		for i, variant := range product.Variants {
			m.Variants[i].FromEntity(&variant)
		}
	}
}

// CategoryModel representa el modelo de persistencia para categorías
type CategoryModel struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	Slug        string `gorm:"uniqueIndex;not null"` // Identificador para tipo de tallas
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
		Slug:        m.Slug,
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
	m.Slug = category.Slug
	m.Description = category.Description
	m.IsActive = category.IsActive
	m.CreatedAt = category.CreatedAt
	m.UpdatedAt = category.UpdatedAt
}
