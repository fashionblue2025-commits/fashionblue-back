package order

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UpdateOrderStatusUseCase struct {
	orderRepo ports.OrderRepository
}

func NewUpdateOrderStatusUseCase(orderRepo ports.OrderRepository) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{
		orderRepo: orderRepo,
	}
}

func (uc *UpdateOrderStatusUseCase) Execute(ctx context.Context, orderID uint, newStatus entities.OrderStatus) error {
	// Obtener orden
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Validar transición de estado
	if !order.CanChangeStatus(newStatus) {
		return errors.New("invalid status transition")
	}

	// Actualizar estado
	return uc.orderRepo.UpdateStatus(ctx, orderID, newStatus)
}
