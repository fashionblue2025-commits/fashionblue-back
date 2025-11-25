package order

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetOrderUseCase struct {
	orderRepo ports.OrderRepository
}

func NewGetOrderUseCase(orderRepo ports.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{
		orderRepo: orderRepo,
	}
}

func (uc *GetOrderUseCase) Execute(ctx context.Context, id uint) (*entities.Order, error) {
	return uc.orderRepo.GetByID(ctx, id)
}
