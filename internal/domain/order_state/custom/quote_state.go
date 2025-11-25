package custom

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
)

type QuoteState struct {
	*order_state.BaseState
}

func NewQuoteState() order_state.OrderState {
	return &QuoteState{
		BaseState: &order_state.BaseState{
			Status: entities.OrderStatusQuote,
			AllowedTransitions: []entities.OrderStatus{
				entities.OrderStatusApproved,
				entities.OrderStatusCancelled,
			},
		},
	}
}

func (s *QuoteState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// Publicar evento de cotización creada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventOrderStatusChanged,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusQuote,
			Timestamp: order.CreatedAt,
		})
	}
	return nil
}
