package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type DeleteProductUseCase struct {
	productRepo ports.ProductRepository
}

func NewDeleteProductUseCase(productRepo ports.ProductRepository) *DeleteProductUseCase {
	return &DeleteProductUseCase{productRepo: productRepo}
}

func (uc *DeleteProductUseCase) Execute(ctx context.Context, id uint) error {
	_, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.productRepo.Delete(ctx, id)
}
