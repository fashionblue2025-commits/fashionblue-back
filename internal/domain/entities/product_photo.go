package entities

import (
	"errors"
	"time"
)

// ProductPhoto representa una foto adjunta a un producto
type ProductPhoto struct {
	ID           uint
	ProductID    uint
	PhotoURL     string
	Description  string
	IsPrimary    bool // Indica si es la foto principal del producto
	DisplayOrder int  // Orden de visualización
	UploadedAt   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Validate valida los datos de la foto
func (pp *ProductPhoto) Validate() error {
	if pp.ProductID == 0 {
		return errors.New("product id is required")
	}
	if pp.PhotoURL == "" {
		return errors.New("photo url is required")
	}
	return nil
}
