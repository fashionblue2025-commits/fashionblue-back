package customer

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetCustomerBalanceUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewGetCustomerBalanceUseCase(customerRepo ports.CustomerRepository) *GetCustomerBalanceUseCase {
	return &GetCustomerBalanceUseCase{customerRepo: customerRepo}
}

func (uc *GetCustomerBalanceUseCase) Execute(ctx context.Context, customerID uint) (float64, error) {
	return uc.customerRepo.GetBalance(ctx, customerID)
}
