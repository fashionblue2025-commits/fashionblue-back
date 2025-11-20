package entities

import "time"

// Category representa una categoría de productos (entidad de dominio pura)
type Category struct {
	ID          uint
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
