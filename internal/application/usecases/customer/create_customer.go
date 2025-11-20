package customer

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreateCustomerUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewCreateCustomerUseCase(customerRepo ports.CustomerRepository) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{customerRepo: customerRepo}
}

func (uc *CreateCustomerUseCase) Execute(ctx context.Context, customer *entities.Customer) error {
	return uc.customerRepo.Create(ctx, customer)
}
