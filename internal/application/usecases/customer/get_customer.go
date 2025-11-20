package customer

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetCustomerUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewGetCustomerUseCase(customerRepo ports.CustomerRepository) *GetCustomerUseCase {
	return &GetCustomerUseCase{customerRepo: customerRepo}
}

func (uc *GetCustomerUseCase) Execute(ctx context.Context, id uint) (*entities.Customer, error) {
	return uc.customerRepo.GetByID(ctx, id)
}
