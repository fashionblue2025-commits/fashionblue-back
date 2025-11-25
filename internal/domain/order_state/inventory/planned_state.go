package inventory

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
)

type PlannedState struct {
	*order_state.BaseState
}

func NewPlannedState() order_state.OrderState {
	return &PlannedState{
		BaseState: &order_state.BaseState{
			Status: entities.OrderStatusPlanned,
			AllowedTransitions: []entities.OrderStatus{
				entities.OrderStatusManufacturing,
				entities.OrderStatusCancelled,
			},
		},
	}
}

func (s *PlannedState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// Publicar evento de producción planificada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventInventoryPlanned,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusPlanned,
		})
		// Los productos se crean/actualizan en FINISHED, no aquí
	}
	return nil
}
