package customer

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type DeleteCustomerUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewDeleteCustomerUseCase(customerRepo ports.CustomerRepository) *DeleteCustomerUseCase {
	return &DeleteCustomerUseCase{customerRepo: customerRepo}
}

func (uc *DeleteCustomerUseCase) Execute(ctx context.Context, id uint) error {
	_, err := uc.customerRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.customerRepo.Delete(ctx, id)
}
