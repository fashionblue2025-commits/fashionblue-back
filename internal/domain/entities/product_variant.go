package entities

import "time"

// ProductVariant representa una variante específica de un producto (color + talla)
// Cada variante tiene su propio stock independiente
// Ejemplo: "Chaqueta Básica - Negro - M" es una variante
type ProductVariant struct {
	ID            uint
	ProductID     uint     // FK al producto base
	Product       *Product // Relación con el producto base
	Color         string   // Color de esta variante (ej: "Negro", "Azul")
	SizeID        *uint    // Talla de esta variante (opcional)
	Size          *Size    // Relación con la talla
	Stock         int      // Stock disponible de esta variante específica
	ReservedStock int      // Stock reservado por órdenes aprobadas
	UnitPrice     float64  // Precio puede variar por variante
	IsActive      bool     // Si esta variante está activa
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// GetAvailableStock retorna el stock disponible (stock total - stock reservado)
func (pv *ProductVariant) GetAvailableStock() int {
	return pv.Stock - pv.ReservedStock
}

// CanReserve verifica si hay suficiente stock disponible para reservar
func (pv *ProductVariant) CanReserve(quantity int) bool {
	return pv.GetAvailableStock() >= quantity
}

// GetFullName retorna el nombre completo de la variante
// Ejemplo: "Chaqueta Básica - Negro - M"
func (pv *ProductVariant) GetFullName() string {
	if pv.Product == nil {
		return ""
	}

	fullName := pv.Product.Name

	if pv.Color != "" {
		fullName += " - " + pv.Color
	}

	if pv.Size != nil {
		fullName += " - " + pv.Size.Value
	}

	return fullName
}

// IsLowStock verifica si el stock está bajo
func (pv *ProductVariant) IsLowStock() bool {
	if pv.Product == nil {
		return false
	}
	return pv.Stock <= pv.Product.MinStock
}
