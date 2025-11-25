package order

import (
	"context"
	"errors"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type DeliverOrderUseCase struct {
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
}

func NewDeliverOrderUseCase(orderRepo ports.OrderRepository, productRepo ports.ProductRepository) *DeliverOrderUseCase {
	return &DeliverOrderUseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (uc *DeliverOrderUseCase) Execute(ctx context.Context, orderID uint) error {
	// Obtener orden
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Validar que se puede entregar
	if !order.CanChangeStatus(entities.OrderStatusDelivered) {
		return errors.New("cannot deliver order in current status")
	}

	// NOTA: La liberación de stock ahora se maneja en DeliveredState.OnEnter()
	// Este código ya no es necesario porque los estados manejan la lógica de stock
	// for _, item := range order.Items {
	// 	variant, err := uc.productVariantRepo.GetByID(ctx, item.ProductVariantID)
	// 	...
	// }

	// Actualizar orden
	now := time.Now()
	order.ActualDeliveryDate = &now
	order.Status = entities.OrderStatusDelivered

	return uc.orderRepo.Update(ctx, order)
}
