package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// ProductRepository define las operaciones para productos
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id uint) (*entities.Product, error)
	GetByAttributes(ctx context.Context, name, color string, sizeID *uint) (*entities.Product, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.Product, error)
	ListByCategory(ctx context.Context, categoryID uint) ([]entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	UpdateStock(ctx context.Context, productID uint, quantity int) error
	Delete(ctx context.Context, id uint) error
	GetLowStockProducts(ctx context.Context) ([]entities.Product, error)
}

// CategoryRepository define las operaciones para categorías
type CategoryRepository interface {
	Create(ctx context.Context, category *entities.Category) error
	GetByID(ctx context.Context, id uint) (*entities.Category, error)
	GetByName(ctx context.Context, name string) (*entities.Category, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.Category, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, id uint) error
}
