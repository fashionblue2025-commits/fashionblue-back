package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetProductUseCase struct {
	productRepo ports.ProductRepository
}

func NewGetProductUseCase(productRepo ports.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{productRepo: productRepo}
}

func (uc *GetProductUseCase) Execute(ctx context.Context, id uint) (*entities.Product, error) {
	return uc.productRepo.GetByID(ctx, id)
}
