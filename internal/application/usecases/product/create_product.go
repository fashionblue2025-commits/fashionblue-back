package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreateProductUseCase struct {
	productRepo        ports.ProductRepository
	productVariantRepo ports.ProductVariantRepository
}

func NewCreateProductUseCase(
	productRepo ports.ProductRepository,
	productVariantRepo ports.ProductVariantRepository,
) *CreateProductUseCase {
	return &CreateProductUseCase{
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
	}
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, product *entities.Product) error {
	// Calcular costo de producción
	product.CalculateProductionCost()

	// Preparar variantes antes de crear el producto
	if len(product.Variants) > 0 {
		for i := range product.Variants {
			variant := &product.Variants[i]

			// Si la variante no tiene precio, usar el del producto base
			if variant.UnitPrice == 0 {
				variant.UnitPrice = product.UnitPrice
			}
		}
	}

	// Crear el producto base (GORM creará las variantes automáticamente por la relación)
	if err := uc.productRepo.Create(ctx, product); err != nil {
		return err
	}

	return nil
}
