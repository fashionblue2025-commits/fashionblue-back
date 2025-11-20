package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListProductsUseCase struct {
	productRepo ports.ProductRepository
}

func NewListProductsUseCase(productRepo ports.ProductRepository) *ListProductsUseCase {
	return &ListProductsUseCase{productRepo: productRepo}
}

func (uc *ListProductsUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]entities.Product, error) {
	return uc.productRepo.List(ctx, filters)
}
