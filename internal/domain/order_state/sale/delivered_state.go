package sale

import (
	"context"
	"log"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type DeliveredState struct {
	*order_state.BaseState
	productVariantRepo ports.ProductVariantRepository
}

func NewDeliveredState(productVariantRepo ports.ProductVariantRepository) order_state.OrderState {
	return &DeliveredState{
		BaseState: &order_state.BaseState{
			Status:             entities.OrderStatusDelivered,
			AllowedTransitions: []entities.OrderStatus{}, // Estado final
		},
		productVariantRepo: productVariantRepo,
	}
}

func (s *DeliveredState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// Liberar stock reservado y descontar del inventario
	if s.productVariantRepo != nil {
		for _, item := range order.Items {
			// Solo procesar items con variante asignada
			if item.ProductVariantID == 0 {
				continue
			}

			// Obtener la variante para saber cuánto está reservado
			variant, err := s.productVariantRepo.GetByID(ctx, item.ProductVariantID)
			if err != nil {
				log.Printf("⚠️  [WARNING] Variant #%d not found for OrderItem #%d: %v", item.ProductVariantID, item.ID, err)
				continue
			}

			// Liberar stock reservado y descontar del inventario
			// ReleaseStock hace ambas cosas: decrementa stock Y reserved_stock
			quantityToRelease := item.Quantity
			if variant.ReservedStock < item.Quantity {
				quantityToRelease = variant.ReservedStock
			}

			if quantityToRelease > 0 {
				if err := s.productVariantRepo.ReleaseStock(ctx, variant.ID, quantityToRelease); err != nil {
					log.Printf("❌ [ERROR] Failed to release stock for variant #%d: %v", variant.ID, err)
					return err
				}
				log.Printf("📦 [DELIVERED] Variant #%d: Released and delivered %d units (stock: -%d, reserved: -%d)",
					variant.ID, quantityToRelease, quantityToRelease, quantityToRelease)
			}
		}
	}

	// Publicar evento de venta entregada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventSaleDelivered,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusDelivered,
		})

		// Si es cliente interno, publicar evento para registro contable
		if order.IsInternalCustomer() {
			data.Publisher.Publish(events.OrderEvent{
				Type:      events.EventInternalCustomerSaleCompleted,
				OrderID:   order.ID,
				Order:     order,
				NewStatus: entities.OrderStatusDelivered,
			})
			log.Printf("💰 [INTERNAL CUSTOMER] Order #%d delivered to customer #%d - Transaction will be created",
				order.ID, *order.CustomerID)
		}
	}

	return nil
}
