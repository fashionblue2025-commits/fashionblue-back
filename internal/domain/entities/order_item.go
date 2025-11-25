package entities

import (
	"errors"
	"time"
)

// OrderItem representa un item de una orden
type OrderItem struct {
	ID               uint
	OrderID          uint
	ProductVariantID uint            // Referencia a la variante específica (color + talla)
	ProductVariant   *ProductVariant // Relación con la variante
	ProductName      string          // Snapshot del nombre del producto base
	CategoryID       uint            // Snapshot de la categoría del producto
	Color            string          // Snapshot del color solicitado
	SizeID           *uint           // Snapshot de la talla solicitada
	Size             *Size           // Relación con la talla
	Quantity         int             // Cantidad total solicitada
	ReservedQuantity int             // Cantidad reservada del stock existente
	UnitPrice        float64
	Subtotal         float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Validate valida los datos del item
func (oi *OrderItem) Validate() error {
	if oi.OrderID == 0 {
		return errors.New("order id is required")
	}
	// ProductID es opcional para órdenes CUSTOM/INVENTORY (se asigna cuando se crea el producto)
	// ProductName es requerido para poder crear el producto después
	if oi.ProductName == "" {
		return errors.New("product name is required")
	}
	if oi.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if oi.UnitPrice < 0 {
		return errors.New("unit price cannot be negative")
	}
	return nil
}

// CalculateSubtotal calcula el subtotal del item
func (oi *OrderItem) CalculateSubtotal() {
	oi.Subtotal = float64(oi.Quantity) * oi.UnitPrice
}

// GetQuantityToManufacture retorna cuántas unidades faltan fabricar
// (Cantidad solicitada - Cantidad reservada del stock)
func (oi *OrderItem) GetQuantityToManufacture(cantReserved int) int {
	toManufacture := oi.Quantity - cantReserved
	if toManufacture < 0 {
		return 0
	}
	return toManufacture
}

// NeedsManufacturing determina si este item requiere fabricación
// Retorna true si:
// 1. Es una variante nueva (ProductVariantID == 0), O
// 2. La cantidad solicitada es mayor que la cantidad reservada del stock
func (oi *OrderItem) NeedsManufacturing() bool {
	// Variante nueva (no existe en el catálogo)
	if oi.ProductVariantID == 0 {
		return true
	}

	// Si la cantidad solicitada es mayor que la reservada, necesita fabricación
	return oi.Quantity > oi.ReservedQuantity
}

// IsFullyCoveredByStock determina si este item está completamente cubierto por stock reservado
// Requiere que ProductVariant esté cargado
func (oi *OrderItem) IsFullyCoveredByStock() bool {
	if oi.ProductVariantID == 0 {
		return false
	}
	if oi.ProductVariant == nil {
		return false // No podemos determinar sin la variante
	}
	return oi.ProductVariant.ReservedStock >= oi.Quantity
}

// IsNewVariant determina si este item representa una variante nueva a crear
func (oi *OrderItem) IsNewVariant() bool {
	return oi.ProductVariantID == 0
}
