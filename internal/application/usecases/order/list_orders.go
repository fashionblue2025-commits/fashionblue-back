package order

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListOrdersUseCase struct {
	orderRepo ports.OrderRepository
}

func NewListOrdersUseCase(orderRepo ports.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		orderRepo: orderRepo,
	}
}

func (uc *ListOrdersUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]entities.Order, error) {
	return uc.orderRepo.List(ctx, filters)
}
