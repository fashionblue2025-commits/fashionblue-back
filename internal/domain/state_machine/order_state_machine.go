package state_machine

import (
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// OrderStateMachine maneja las transiciones de estado de órdenes
type OrderStateMachine struct {
	transitions map[entities.OrderType]map[entities.OrderStatus][]entities.OrderStatus
}

// NewOrderStateMachine crea una nueva máquina de estados
func NewOrderStateMachine() *OrderStateMachine {
	return &OrderStateMachine{
		transitions: map[entities.OrderType]map[entities.OrderStatus][]entities.OrderStatus{
			entities.OrderTypeCustom: {
				entities.OrderStatusQuote:         {entities.OrderStatusApproved, entities.OrderStatusCancelled},
				entities.OrderStatusApproved:      {entities.OrderStatusManufacturing, entities.OrderStatusCancelled},
				entities.OrderStatusManufacturing: {entities.OrderStatusFinished, entities.OrderStatusCancelled},
				entities.OrderStatusFinished:      {entities.OrderStatusDelivered, entities.OrderStatusCancelled},
				entities.OrderStatusDelivered:     {},
				entities.OrderStatusCancelled:     {},
			},
			entities.OrderTypeInventory: {
				entities.OrderStatusPlanned:       {entities.OrderStatusManufacturing, entities.OrderStatusCancelled},
				entities.OrderStatusManufacturing: {entities.OrderStatusFinished, entities.OrderStatusCancelled},
				entities.OrderStatusFinished:      {},
				entities.OrderStatusCancelled:     {},
			},
			entities.OrderTypeSale: {
				entities.OrderStatusPending:   {entities.OrderStatusConfirmed, entities.OrderStatusCancelled},
				entities.OrderStatusConfirmed: {entities.OrderStatusDelivered, entities.OrderStatusCancelled},
				entities.OrderStatusDelivered: {},
				entities.OrderStatusCancelled: {},
			},
		},
	}
}

func (sm *OrderStateMachine) CanTransition(orderType entities.OrderType, from, to entities.OrderStatus) bool {
	allowedTransitions, exists := sm.transitions[orderType][from]
	if !exists {
		return false
	}
	for _, allowed := range allowedTransitions {
		if allowed == to {
			return true
		}
	}
	return false
}

func (sm *OrderStateMachine) ValidateTransition(order *entities.Order, newStatus entities.OrderStatus) error {
	if order.Status == newStatus {
		return errors.New("order is already in this status")
	}
	if !sm.CanTransition(order.Type, order.Status, newStatus) {
		return errors.New("invalid state transition")
	}
	return nil
}

func (sm *OrderStateMachine) GetInitialStatus(orderType entities.OrderType) entities.OrderStatus {
	switch orderType {
	case entities.OrderTypeCustom:
		return entities.OrderStatusQuote
	case entities.OrderTypeInventory:
		return entities.OrderStatusPlanned
	case entities.OrderTypeSale:
		return entities.OrderStatusPending
	default:
		return entities.OrderStatusQuote
	}
}
