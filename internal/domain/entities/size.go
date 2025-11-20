package entities

import "time"

// SizeType representa el tipo de talla
type SizeType string

const (
	SizeTypeShirt SizeType = "SHIRT" // Camiseta
	SizeTypePants SizeType = "PANTS" // Pantalón
	SizeTypeShoes SizeType = "SHOES" // Tenis/Zapatos
)

// Size representa una talla disponible en el sistema
type Size struct {
	ID        uint      `json:"id"`
	Type      SizeType  `json:"type"`      // SHIRT, PANTS, SHOES
	Value     string    `json:"value"`     // "S", "M", "L", "XL", "28", "30", "7", "8.5", etc.
	Order     int       `json:"order"`     // Para ordenar las tallas (S=1, M=2, L=3, etc.)
	IsActive  bool      `json:"is_active"` // Si la talla está activa
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
