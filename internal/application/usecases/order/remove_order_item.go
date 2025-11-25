package order

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type RemoveOrderItemUseCase struct {
	orderRepo     ports.OrderRepository
	orderItemRepo ports.OrderItemRepository
}

func NewRemoveOrderItemUseCase(
	orderRepo ports.OrderRepository,
	orderItemRepo ports.OrderItemRepository,
) *RemoveOrderItemUseCase {
	return &RemoveOrderItemUseCase{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (uc *RemoveOrderItemUseCase) Execute(ctx context.Context, itemID uint) error {
	// Obtener item
	item, err := uc.orderItemRepo.GetByID(ctx, itemID)
	if err != nil {
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

	// Eliminar item
	if err := uc.orderItemRepo.Delete(ctx, itemID); err != nil {
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
