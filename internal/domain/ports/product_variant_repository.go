package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// ProductVariantRepository define las operaciones para variantes de producto
type ProductVariantRepository interface {
	// Create crea una nueva variante
	Create(ctx context.Context, variant *entities.ProductVariant) error

	// GetByID obtiene una variante por su ID
	GetByID(ctx context.Context, id uint) (*entities.ProductVariant, error)

	// GetByProductAndAttributes busca una variante específica por producto, color y talla
	GetByProductAndAttributes(ctx context.Context, productID uint, color string, sizeID *uint) (*entities.ProductVariant, error)

	// ListByProduct lista todas las variantes de un producto
	ListByProduct(ctx context.Context, productID uint) ([]entities.ProductVariant, error)

	// Update actualiza una variante
	Update(ctx context.Context, variant *entities.ProductVariant) error

	// UpdateStock actualiza el stock de una variante (incrementa o decrementa)
	UpdateStock(ctx context.Context, variantID uint, quantity int) error

	// ReserveStock reserva stock de una variante
	ReserveStock(ctx context.Context, variantID uint, quantity int) error

	// ReleaseStock libera stock reservado de una variante
	ReleaseStock(ctx context.Context, variantID uint, quantity int) error

	// Delete elimina una variante
	Delete(ctx context.Context, id uint) error

	// GetLowStockVariants obtiene variantes con stock bajo
	GetLowStockVariants(ctx context.Context) ([]entities.ProductVariant, error)
}
