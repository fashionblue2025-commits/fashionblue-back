package order

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UpdateOrderItemUseCase struct {
	orderRepo     ports.OrderRepository
	orderItemRepo ports.OrderItemRepository
}

func NewUpdateOrderItemUseCase(
	orderRepo ports.OrderRepository,
	orderItemRepo ports.OrderItemRepository,
) *UpdateOrderItemUseCase {
	return &UpdateOrderItemUseCase{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (uc *UpdateOrderItemUseCase) Execute(ctx context.Context, item *entities.OrderItem) error {
	// Validar item
	if err := item.Validate(); err != nil {
		return err
	}

	// Obtener orden
	order, err := uc.orderRepo.GetByID(ctx, item.OrderID)
	if err != nil {
		return err
	}

	// Verificar que se pueden editar items
	if !order.CanEditItems() {
		return errors.New("cannot edit items in current order status")
	}

	// Calcular subtotal
	item.CalculateSubtotal()

	// Actualizar item
	if err := uc.orderItemRepo.Update(ctx, item); err != nil {
		return err
	}

	// Recalcular total de la orden
	items, err := uc.orderItemRepo.GetByOrderID(ctx, item.OrderID)
	if err != nil {
		return err
	}
	order.Items = items
	order.TotalAmount = order.CalculateTotal()

	return uc.orderRepo.Update(ctx, order)
}
