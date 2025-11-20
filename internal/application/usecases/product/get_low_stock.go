package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetLowStockProductsUseCase struct {
	productRepo ports.ProductRepository
}

func NewGetLowStockProductsUseCase(productRepo ports.ProductRepository) *GetLowStockProductsUseCase {
	return &GetLowStockProductsUseCase{productRepo: productRepo}
}

func (uc *GetLowStockProductsUseCase) Execute(ctx context.Context) ([]entities.Product, error) {
	return uc.productRepo.GetLowStockProducts(ctx)
}
