package dto

import "github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"

type Product struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	SKU             string  `json:"sku"`
	CategoryID      uint    `json:"category_id"`
	MaterialCost    float64 `json:"material_cost"`
	LaborCost       float64 `json:"labor_cost"`
	ProductionCost  float64 `json:"production_cost"`
	UnitPrice       float64 `json:"unit_price"`
	WholesalePrice  float64 `json:"wholesale_price"`
	MinWholesaleQty int     `json:"min_wholesale_qty"`
	Stock           int     `json:"stock"`
	MinStock        int     `json:"min_stock"`
	IsActive        bool    `json:"is_active"`
}

func ToProductDTO(product *entities.Product) *Product {
	return &Product{
		ID:              product.ID,
		Name:            product.Name,
		Description:     product.Description,
		SKU:             product.SKU,
		CategoryID:      product.CategoryID,
		MaterialCost:    product.MaterialCost,
		LaborCost:       product.LaborCost,
		ProductionCost:  product.ProductionCost,
		UnitPrice:       product.UnitPrice,
		WholesalePrice:  product.WholesalePrice,
		MinWholesaleQty: product.MinWholesaleQty,
		Stock:           product.Stock,
		MinStock:        product.MinStock,
		IsActive:        product.IsActive,
	}
}

func ToProductDTOList(products []entities.Product) []*Product {
	dtos := make([]*Product, len(products))
	for i, product := range products {
		dtos[i] = ToProductDTO(&product)
	}
	return dtos
}

func FromProductDTO(product *Product) *entities.Product {
	return &entities.Product{
		ID:              product.ID,
		Name:            product.Name,
		Description:     product.Description,
		SKU:             product.SKU,
		CategoryID:      product.CategoryID,
		MaterialCost:    product.MaterialCost,
		LaborCost:       product.LaborCost,
		ProductionCost:  product.ProductionCost,
		UnitPrice:       product.UnitPrice,
		WholesalePrice:  product.WholesalePrice,
		MinWholesaleQty: product.MinWholesaleQty,
		Stock:           product.Stock,
		MinStock:        product.MinStock,
		IsActive:        product.IsActive,
	}
}
