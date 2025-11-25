package custom

import (
	"context"
	"log"

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
	// 🔓 Liberar stock reservado
	for _, item := range order.Items {
		if item.ProductVariantID == 0 {
			continue
		}

		// Obtener variante para saber cuánto está reservado
		variant, err := s.productVariantRepo.GetByID(ctx, item.ProductVariantID)
		if err != nil {
			log.Printf("⚠️  [WARNING] Variant #%d not found: %v", item.ProductVariantID, err)
			continue
		}

		if variant.ReservedStock > 0 {
			// Liberar solo reserved_stock (no decrementar stock total)
			variant.ReservedStock -= item.Quantity
			if variant.ReservedStock < 0 {
				variant.ReservedStock = 0
			}

			if err := s.productVariantRepo.Update(ctx, variant); err != nil {
				log.Printf("⚠️  [WARNING] Failed to release stock for variant #%d: %v", variant.ID, err)
				// No retornamos error para no bloquear la cancelación
			} else {
				log.Printf("🔓 [RELEASED] Variant #%d: Released %d units", variant.ID, item.Quantity)
			}
		}
	}

	// Publicar evento de orden cancelada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventOrderCancelled,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusCancelled,
		})
	}
	return nil
}
