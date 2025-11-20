package customer

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UpdateCustomerUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewUpdateCustomerUseCase(customerRepo ports.CustomerRepository) *UpdateCustomerUseCase {
	return &UpdateCustomerUseCase{customerRepo: customerRepo}
}

func (uc *UpdateCustomerUseCase) Execute(ctx context.Context, customer *entities.Customer) error {
	return uc.customerRepo.Update(ctx, customer)
}
