package supplier

import (
	"context"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type DeleteSupplierUseCase struct {
	repo ports.SupplierRepository
}

func NewDeleteSupplierUseCase(repo ports.SupplierRepository) *DeleteSupplierUseCase {
	return &DeleteSupplierUseCase{repo: repo}
}

func (uc *DeleteSupplierUseCase) Execute(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}
