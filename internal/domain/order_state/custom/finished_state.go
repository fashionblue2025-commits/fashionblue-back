package custom

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type FinishedState struct {
	*order_state.BaseState
	productRepo ports.ProductRepository
}

func NewFinishedState(productRepo ports.ProductRepository) order_state.OrderState {
	return &FinishedState{
		BaseState: &order_state.BaseState{
			Status: entities.OrderStatusFinished,
			AllowedTransitions: []entities.OrderStatus{
				entities.OrderStatusDelivered,
				entities.OrderStatusCancelled,
			},
		},
		productRepo: productRepo,
	}
}

func (s *FinishedState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// 🏭 FINISHED = Productos terminados, listos para entregar
	//
	// Lógica de stock:
	// 1. Items con stock reservado (ReservedQuantity > 0): Ya están contabilizados, NO modificar stock
	// 2. Items fabricados (Quantity - ReservedQuantity): Se fabricaron pero NO incrementar stock aún
	//    - El stock se incrementará cuando se creen las variantes nuevas si es necesario
	//    - Para variantes existentes, el stock reservado ya está contabilizado
	//
	// Cuando se entregue (DELIVERED), se liberará el stock reservado

	// Publicar evento de producción finalizada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventOrderFinished,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusFinished,
			Data: map[string]interface{}{
				"orderType":          order.Type,
				"producedQuantities": data.ProducedQuantities,
			},
		})

		// 🏭 Solo crear variantes NUEVAS (ProductVariantID == 0)
		// NO incrementar stock de variantes existentes
		data.Publisher.Publish(events.OrderEvent{
			Type:    events.EventProductCreationRequired,
			OrderID: order.ID,
			Order:   order,
			Data: map[string]interface{}{
				"orderType":          entities.OrderTypeCustom,
				"items":              order.Items,
				"producedQuantities": data.ProducedQuantities,
			},
		})
	}

	return nil
}
