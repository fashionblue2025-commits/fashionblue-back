package inventory

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
)

type ManufacturingState struct {
	*order_state.BaseState
}

func NewManufacturingState() order_state.OrderState {
	return &ManufacturingState{
		BaseState: &order_state.BaseState{
			Status: entities.OrderStatusManufacturing,
			AllowedTransitions: []entities.OrderStatus{
				entities.OrderStatusFinished,
				entities.OrderStatusCancelled,
			},
		},
	}
}

func (s *ManufacturingState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// Publicar evento de fabricación para inventario iniciada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventInventoryManufacturing,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusManufacturing,
		})
	}
	return nil
}
