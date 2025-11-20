package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UpdateProductUseCase struct {
	productRepo ports.ProductRepository
}

func NewUpdateProductUseCase(productRepo ports.ProductRepository) *UpdateProductUseCase {
	return &UpdateProductUseCase{productRepo: productRepo}
}

func (uc *UpdateProductUseCase) Execute(ctx context.Context, product *entities.Product) error {
	product.CalculateProductionCost()
	return uc.productRepo.Update(ctx, product)
}
