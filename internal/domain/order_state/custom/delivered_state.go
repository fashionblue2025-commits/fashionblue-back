package custom

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
			// Solo liberar si hay stock reservado
			if item.ReservedQuantity <= 0 {
				log.Printf("ℹ️  [SKIP] OrderItem #%d has no reserved stock to release", item.ID)
				continue
			}

			// Solo procesar items con variante asignada
			if item.ProductVariantID == 0 {
				log.Printf("⚠️  [WARNING] OrderItem #%d has no ProductVariantID assigned", item.ID)
				continue
			}

			// Liberar el stock que fue reservado en APPROVED
			// ReleaseStock hace: stock -= quantity, reserved_stock -= quantity
			if err := s.productVariantRepo.ReleaseStock(ctx, item.ProductVariantID, item.ReservedQuantity); err != nil {
				log.Printf("❌ [ERROR] Failed to release stock for variant #%d: %v", item.ProductVariantID, err)
				return err
			}

			log.Printf("📦 [DELIVERED] Variant #%d: Released %d units (stock: -%d, reserved: -%d)",
				item.ProductVariantID, item.ReservedQuantity, item.ReservedQuantity, item.ReservedQuantity)
		}
	}

	// Publicar evento de orden entregada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventOrderDelivered,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusDelivered,
		})

		// Publicar evento de venta completada para registrar ingreso financiero automático
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventSaleCompleted,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusDelivered,
			Data: map[string]interface{}{
				"total_amount": order.TotalAmount,
				"order_type":   order.Type,
			},
		})
		log.Printf("💰 [SALE COMPLETED] Order #%d (Type: %s) - Total: $%.2f - Financial income will be recorded",
			order.ID, order.Type, order.TotalAmount)

		// Si es cliente interno, publicar evento adicional para registro contable
		if order.IsInternalCustomer() {
			data.Publisher.Publish(events.OrderEvent{
				Type:      events.EventInternalCustomerSaleCompleted,
				OrderID:   order.ID,
				Order:     order,
				NewStatus: entities.OrderStatusDelivered,
			})
			log.Printf("📝 [INTERNAL CUSTOMER] Order #%d delivered to customer #%d - Customer transaction will be created",
				order.ID, *order.CustomerID)
		}
	}

	return nil
}
