package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// SizeRepository define las operaciones de persistencia para tallas
type SizeRepository interface {
	// Create crea una nueva talla
	Create(ctx context.Context, size *entities.Size) error

	// GetByID obtiene una talla por su ID
	GetByID(ctx context.Context, id uint) (*entities.Size, error)

	// List lista todas las tallas con filtros opcionales
	List(ctx context.Context, filters map[string]interface{}) ([]*entities.Size, error)

	// Update actualiza una talla existente
	Update(ctx context.Context, size *entities.Size) error

	// Delete elimina una talla (soft delete)
	Delete(ctx context.Context, id uint) error

	// GetByType obtiene todas las tallas de un tipo específico
	GetByType(ctx context.Context, sizeType entities.SizeType) ([]*entities.Size, error)
}
