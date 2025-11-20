package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// SupplierRepository define las operaciones de persistencia para proveedores
type SupplierRepository interface {
	Create(ctx context.Context, supplier *entities.Supplier) error
	GetByID(ctx context.Context, id uint) (*entities.Supplier, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.Supplier, error)
	Update(ctx context.Context, supplier *entities.Supplier) error
	Delete(ctx context.Context, id uint) error
}
