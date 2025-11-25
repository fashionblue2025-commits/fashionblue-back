package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// ProductPhotoModel representa el modelo de persistencia para fotos de productos
type ProductPhotoModel struct {
	ID           uint          `gorm:"primaryKey"`
	ProductID    uint          `gorm:"not null;index"`
	Product      *ProductModel `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	PhotoURL     string        `gorm:"not null"`
	Description  string        `gorm:"type:text"`
	IsPrimary    bool          `gorm:"default:false;index"`
	DisplayOrder int           `gorm:"default:0"`
	UploadedAt   time.Time     `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName especifica el nombre de la tabla
func (ProductPhotoModel) TableName() string {
	return "product_photos"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *ProductPhotoModel) ToEntity() *entities.ProductPhoto {
	return &entities.ProductPhoto{
		ID:           m.ID,
		ProductID:    m.ProductID,
		PhotoURL:     m.PhotoURL,
		Description:  m.Description,
		IsPrimary:    m.IsPrimary,
		DisplayOrder: m.DisplayOrder,
		UploadedAt:   m.UploadedAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo
func (m *ProductPhotoModel) FromEntity(photo *entities.ProductPhoto) {
	m.ID = photo.ID
	m.ProductID = photo.ProductID
	m.PhotoURL = photo.PhotoURL
	m.Description = photo.Description
	m.IsPrimary = photo.IsPrimary
	m.DisplayOrder = photo.DisplayOrder
	m.UploadedAt = photo.UploadedAt
	m.CreatedAt = photo.CreatedAt
	m.UpdatedAt = photo.UpdatedAt
}
