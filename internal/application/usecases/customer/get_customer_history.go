package customer

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetCustomerHistoryUseCase struct {
	customerTransactionRepo ports.CustomerTransactionRepository
}

func NewGetCustomerHistoryUseCase(customerTransactionRepo ports.CustomerTransactionRepository) *GetCustomerHistoryUseCase {
	return &GetCustomerHistoryUseCase{customerTransactionRepo: customerTransactionRepo}
}

func (uc *GetCustomerHistoryUseCase) Execute(ctx context.Context, customerID uint) ([]entities.CustomerTransaction, error) {
	return uc.customerTransactionRepo.ListByCustomer(ctx, customerID)
}
