package sale

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
)

type ConfirmedState struct {
	*order_state.BaseState
}

func NewConfirmedState() order_state.OrderState {
	return &ConfirmedState{
		BaseState: &order_state.BaseState{
			Status: entities.OrderStatusConfirmed,
			AllowedTransitions: []entities.OrderStatus{
				entities.OrderStatusDelivered,
				entities.OrderStatusCancelled,
			},
		},
	}
}

func (s *ConfirmedState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// Publicar evento de venta confirmada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventSaleConfirmed,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusConfirmed,
		})
	}
	return nil
}
