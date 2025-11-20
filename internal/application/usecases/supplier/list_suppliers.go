package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListSuppliersUseCase struct {
	repo ports.SupplierRepository
}

func NewListSuppliersUseCase(repo ports.SupplierRepository) *ListSuppliersUseCase {
	return &ListSuppliersUseCase{repo: repo}
}

func (uc *ListSuppliersUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]entities.Supplier, error) {
	return uc.repo.List(ctx, filters)
}
