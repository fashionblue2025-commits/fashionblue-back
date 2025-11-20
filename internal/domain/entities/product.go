package entities

import "time"

// Product representa un producto de la empresa (entidad de dominio pura)
type Product struct {
	ID              uint
	Name            string
	Description     string
	SKU             string
	CategoryID      uint
	MaterialCost    float64
	LaborCost       float64
	ProductionCost  float64
	UnitPrice       float64
	WholesalePrice  float64
	MinWholesaleQty int
	Stock           int
	MinStock        int
	IsActive        bool
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

// IsLowStock verifica si el stock está bajo
func (p *Product) IsLowStock() bool {
	return p.Stock <= p.MinStock
}
