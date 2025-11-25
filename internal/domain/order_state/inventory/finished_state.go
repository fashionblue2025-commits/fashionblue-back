package inventory

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
			Status:             entities.OrderStatusFinished,
			AllowedTransitions: []entities.OrderStatus{}, // Estado final para INVENTORY
		},
		productRepo: productRepo,
	}
}

func (s *FinishedState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// Publicar evento de producción para inventario finalizada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventInventoryFinished,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusFinished,
		})

		// 🏭 Publicar evento para crear productos nuevos y actualizar stock de existentes
		// El event handler maneja ambos casos:
		// - Si ProductID != 0: Incrementa stock del producto existente
		// - Si ProductID == 0: Crea el producto con stock inicial
		data.Publisher.Publish(events.OrderEvent{
			Type:    events.EventProductCreationRequired,
			OrderID: order.ID,
			Order:   order,
			Data: map[string]interface{}{
				"orderType": entities.OrderTypeInventory,
				"items":     order.Items,
			},
		})

		// Publicar evento de stock actualizado
		data.Publisher.Publish(events.OrderEvent{
			Type:    events.EventStockUpdated,
			OrderID: order.ID,
			Order:   order,
			Data: map[string]interface{}{
				"forInventory": true,
			},
		})
	}

	return nil
}
