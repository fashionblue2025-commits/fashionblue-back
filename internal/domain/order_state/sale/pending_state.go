package sale

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type PendingState struct {
	*order_state.BaseState
	productVariantRepo ports.ProductVariantRepository
}

func NewPendingState(productVariantRepo ports.ProductVariantRepository) order_state.OrderState {
	return &PendingState{
		BaseState: &order_state.BaseState{
			Status: entities.OrderStatusPending,
			AllowedTransitions: []entities.OrderStatus{
				entities.OrderStatusConfirmed,
				entities.OrderStatusCancelled,
			},
		},
		productVariantRepo: productVariantRepo,
	}
}

func (s *PendingState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// Reservar stock de variantes
	if s.productVariantRepo != nil {
		for i := range order.Items {
			item := &order.Items[i]

			if item.ProductVariantID == 0 {
				return errors.New("product variant not found for item: " + item.ProductName)
			}

			variant, err := s.productVariantRepo.GetByID(ctx, item.ProductVariantID)
			if err != nil {
				return err
			}

			// Verificar stock disponible
			if !variant.CanReserve(item.Quantity) {
				return errors.New("insufficient stock for variant: " + variant.GetFullName())
			}

			// Reservar stock
			if err := s.productVariantRepo.ReserveStock(ctx, variant.ID, item.Quantity); err != nil {
				return err
			}
		}
	}

	// Publicar evento de venta pendiente
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventSalePending,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusPending,
		})

		// Publicar evento de stock reservado
		data.Publisher.Publish(events.OrderEvent{
			Type:    events.EventStockReserved,
			OrderID: order.ID,
			Order:   order,
			Data: map[string]interface{}{
				"items": order.Items,
			},
		})
	}

	return nil
}
