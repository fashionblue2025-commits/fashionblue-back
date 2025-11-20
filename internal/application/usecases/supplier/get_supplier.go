package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetSupplierUseCase struct {
	repo ports.SupplierRepository
}

func NewGetSupplierUseCase(repo ports.SupplierRepository) *GetSupplierUseCase {
	return &GetSupplierUseCase{repo: repo}
}

func (uc *GetSupplierUseCase) Execute(ctx context.Context, id uint) (*entities.Supplier, error) {
	return uc.repo.GetByID(ctx, id)
}
