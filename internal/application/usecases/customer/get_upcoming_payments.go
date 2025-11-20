package customer

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetUpcomingPaymentsUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewGetUpcomingPaymentsUseCase(customerRepo ports.CustomerRepository) *GetUpcomingPaymentsUseCase {
	return &GetUpcomingPaymentsUseCase{customerRepo: customerRepo}
}

type CustomerWithBalance struct {
	Customer entities.Customer
	Balance  float64
}

func (uc *GetUpcomingPaymentsUseCase) Execute(ctx context.Context, daysRange int) ([]CustomerWithBalance, error) {
	// Obtener clientes con pagos próximos
	customers, err := uc.customerRepo.GetUpcomingPayments(ctx, daysRange)
	if err != nil {
		return nil, err
	}

	// Obtener balance de cada cliente
	result := make([]CustomerWithBalance, len(customers))
	for i, customer := range customers {
		balance, err := uc.customerRepo.GetBalance(ctx, customer.ID)
		if err != nil {
			balance = 0 // Si hay error, asumir balance 0
		}

		result[i] = CustomerWithBalance{
			Customer: customer,
			Balance:  balance,
		}
	}

	return result, nil
}
