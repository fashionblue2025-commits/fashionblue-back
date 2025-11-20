package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// CapitalInjectionRepository define las operaciones de persistencia para inyecciones de capital
type CapitalInjectionRepository interface {
	Create(ctx context.Context, injection *entities.CapitalInjection) error
	GetByID(ctx context.Context, id uint) (*entities.CapitalInjection, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.CapitalInjection, error)
	GetTotal(ctx context.Context) (float64, error)
}
