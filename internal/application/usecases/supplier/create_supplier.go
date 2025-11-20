package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreateSupplierUseCase struct {
	repo ports.SupplierRepository
}

func NewCreateSupplierUseCase(repo ports.SupplierRepository) *CreateSupplierUseCase {
	return &CreateSupplierUseCase{repo: repo}
}

func (uc *CreateSupplierUseCase) Execute(ctx context.Context, supplier *entities.Supplier) error {
	if err := supplier.Validate(); err != nil {
		return err
	}
	supplier.IsActive = true
	return uc.repo.Create(ctx, supplier)
}
