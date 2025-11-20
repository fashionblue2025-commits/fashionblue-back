package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UpdateSupplierUseCase struct {
	repo ports.SupplierRepository
}

func NewUpdateSupplierUseCase(repo ports.SupplierRepository) *UpdateSupplierUseCase {
	return &UpdateSupplierUseCase{repo: repo}
}

func (uc *UpdateSupplierUseCase) Execute(ctx context.Context, supplier *entities.Supplier) error {
	if err := supplier.Validate(); err != nil {
		return err
	}
	return uc.repo.Update(ctx, supplier)
}
