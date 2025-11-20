package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreateProductUseCase struct {
	productRepo ports.ProductRepository
}

func NewCreateProductUseCase(productRepo ports.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{productRepo: productRepo}
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, product *entities.Product) error {
	product.CalculateProductionCost()
	return uc.productRepo.Create(ctx, product)
}
