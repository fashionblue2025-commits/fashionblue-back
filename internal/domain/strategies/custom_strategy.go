package strategies

import (
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state/custom"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// CustomOrderStrategy estrategia para órdenes CUSTOM (cotización)
type CustomOrderStrategy struct {
	states             map[entities.OrderStatus]order_state.OrderState
	eventPublisher     ports.EventPublisher
	productRepo        ports.ProductRepository
	productVariantRepo ports.ProductVariantRepository
}

func NewCustomOrderStrategy(eventPublisher ports.EventPublisher, productRepo ports.ProductRepository, productVariantRepo ports.ProductVariantRepository) order_state.OrderStrategy {
	return &CustomOrderStrategy{
		eventPublisher:     eventPublisher,
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		states: map[entities.OrderStatus]order_state.OrderState{
			entities.OrderStatusQuote:         custom.NewQuoteState(),
			entities.OrderStatusApproved:      custom.NewApprovedState(productRepo, productVariantRepo),
			entities.OrderStatusManufacturing: custom.NewManufacturingState(),
			entities.OrderStatusFinished:      custom.NewFinishedState(productRepo),
			entities.OrderStatusDelivered:     custom.NewDeliveredState(productVariantRepo),
			entities.OrderStatusCancelled:     custom.NewCancelledState(productVariantRepo),
		},
	}
}

func (s *CustomOrderStrategy) GetInitialStatus() entities.OrderStatus {
	return entities.OrderStatusQuote
}

func (s *CustomOrderStrategy) GetAllowedTransitions(currentStatus entities.OrderStatus) []entities.OrderStatus {
	state := s.states[currentStatus]
	if state == nil {
		return []entities.OrderStatus{}
	}
	return state.GetAllowedTransitions()
}

func (s *CustomOrderStrategy) CanTransition(from, to entities.OrderStatus) bool {
	state := s.states[from]
	if state == nil {
		return false
	}
	return state.CanTransitionTo(to)
}

func (s *CustomOrderStrategy) GetState(status entities.OrderStatus) order_state.OrderState {
	return s.states[status]
}
