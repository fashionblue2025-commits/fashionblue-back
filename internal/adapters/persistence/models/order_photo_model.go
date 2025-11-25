package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"gorm.io/gorm"
)

// OrderPhotoModel representa el modelo de persistencia de una foto de orden
type OrderPhotoModel struct {
	ID          uint   `gorm:"primaryKey"`
	OrderID     uint   `gorm:"not null;index"`
	PhotoURL    string `gorm:"not null"`
	Description string
	UploadedAt  time.Time `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relaciones
	Order *OrderModel `gorm:"foreignKey:OrderID"`
}

// TableName especifica el nombre de la tabla
func (OrderPhotoModel) TableName() string {
	return "order_photos"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *OrderPhotoModel) ToEntity() *entities.OrderPhoto {
	return &entities.OrderPhoto{
		ID:          m.ID,
		OrderID:     m.OrderID,
		PhotoURL:    m.PhotoURL,
		Description: m.Description,
		UploadedAt:  m.UploadedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromEntity convierte la entidad de dominio a modelo
func (m *OrderPhotoModel) FromEntity(photo *entities.OrderPhoto) {
	m.ID = photo.ID
	m.OrderID = photo.OrderID
	m.PhotoURL = photo.PhotoURL
	m.Description = photo.Description
	m.UploadedAt = photo.UploadedAt
}
