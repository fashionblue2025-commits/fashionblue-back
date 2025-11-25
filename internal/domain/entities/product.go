package entities

import "time"

// Product representa un producto base de la empresa (sin variantes de color/talla)
// Este es el "producto maestro" que contiene información común a todas las variantes
// Ejemplo: "Chaqueta Básica" es el producto base
// Las variantes serían: "Chaqueta Básica - Negro - M", "Chaqueta Básica - Azul - L", etc.
type Product struct {
	ID              uint
	Name            string           // Nombre base del producto (ej: "Chaqueta Básica")
	Description     string           // Descripción del producto
	CategoryID      uint             // Categoría del producto
	Category        *Category        // Relación con la categoría
	MaterialCost    float64          // Costo de materiales (común a todas las variantes)
	LaborCost       float64          // Costo de mano de obra (común a todas las variantes)
	ProductionCost  float64          // Costo total de producción
	UnitPrice       float64          // Precio base (puede ser sobrescrito por variante)
	WholesalePrice  float64          // Precio mayorista base
	MinWholesaleQty int              // Cantidad mínima para precio mayorista
	MinStock        int              // Stock mínimo para alertas (aplica a variantes)
	IsActive        bool             // Si el producto está activo
	Photos          []ProductPhoto   // Fotos del producto base
	Variants        []ProductVariant // Variantes de este producto (color + talla + stock)
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// CalculateProductionCost calcula el costo total de producción
func (p *Product) CalculateProductionCost() {
	p.ProductionCost = p.MaterialCost + p.LaborCost
}

// GetUnitProfit calcula la ganancia por unidad vendida al precio unitario
func (p *Product) GetUnitProfit() float64 {
	return p.UnitPrice - p.ProductionCost
}

// GetWholesaleProfit calcula la ganancia por unidad vendida al precio mayorista
func (p *Product) GetWholesaleProfit() float64 {
	return p.WholesalePrice - p.ProductionCost
}

// GetProfitMargin calcula el margen de ganancia porcentual (precio unitario)
func (p *Product) GetProfitMargin() float64 {
	if p.ProductionCost == 0 {
		return 0
	}
	return ((p.UnitPrice - p.ProductionCost) / p.ProductionCost) * 100
}

// GetTotalStock retorna el stock total de todas las variantes
func (p *Product) GetTotalStock() int {
	total := 0
	for _, variant := range p.Variants {
		total += variant.Stock
	}
	return total
}

// GetTotalAvailableStock retorna el stock disponible total de todas las variantes
func (p *Product) GetTotalAvailableStock() int {
	total := 0
	for _, variant := range p.Variants {
		total += variant.GetAvailableStock()
	}
	return total
}

// HasLowStockVariants verifica si alguna variante tiene stock bajo
func (p *Product) HasLowStockVariants() bool {
	for _, variant := range p.Variants {
		if variant.IsLowStock() {
			return true
		}
	}
	return false
}

// GetVariantByAttributes busca una variante específica por color y talla
func (p *Product) GetVariantByAttributes(color string, sizeID *uint) *ProductVariant {
	for i := range p.Variants {
		variant := &p.Variants[i]

		// Comparar color
		if variant.Color != color {
			continue
		}

		// Comparar sizeID (manejar NULL)
		if sizeID == nil && variant.SizeID == nil {
			return variant
		}
		if sizeID != nil && variant.SizeID != nil && *sizeID == *variant.SizeID {
			return variant
		}
	}
	return nil
}
