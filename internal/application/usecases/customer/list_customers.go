package customer

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListCustomersUseCase struct {
	customerRepo ports.CustomerRepository
}

func NewListCustomersUseCase(customerRepo ports.CustomerRepository) *ListCustomersUseCase {
	return &ListCustomersUseCase{customerRepo: customerRepo}
}

func (uc *ListCustomersUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]entities.Customer, error) {
	return uc.customerRepo.List(ctx, filters)
}
