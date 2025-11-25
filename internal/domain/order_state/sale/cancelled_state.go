package sale

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CancelledState struct {
	*order_state.BaseState
	productVariantRepo ports.ProductVariantRepository
}

func NewCancelledState(productVariantRepo ports.ProductVariantRepository) order_state.OrderState {
	return &CancelledState{
		BaseState: &order_state.BaseState{
			Status:             entities.OrderStatusCancelled,
			AllowedTransitions: []entities.OrderStatus{}, // Estado final
		},
		productVariantRepo: productVariantRepo,
	}
}

func (s *CancelledState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// Liberar stock reservado de variantes
	if s.productVariantRepo != nil {
		for _, item := range order.Items {
			if item.ProductVariantID == 0 {
				continue
			}

			// Obtener variante para saber cuánto está reservado
			variant, err := s.productVariantRepo.GetByID(ctx, item.ProductVariantID)
			if err != nil {
				continue // Si no existe, continuar
			}

			// Liberar reserva (solo decrementar reserved_stock, no el stock total)
			if variant.ReservedStock > 0 {
				variant.ReservedStock -= item.Quantity
				if variant.ReservedStock < 0 {
					variant.ReservedStock = 0
				}
				if err := s.productVariantRepo.Update(ctx, variant); err != nil {
					return err
				}
			}
		}
	}

	// Publicar evento de venta cancelada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventOrderCancelled,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusCancelled,
		})

		// Publicar evento de stock liberado
		data.Publisher.Publish(events.OrderEvent{
			Type:    events.EventStockReleased,
			OrderID: order.ID,
			Order:   order,
			Data: map[string]interface{}{
				"items": order.Items,
			},
		})
	}

	return nil
}
