package strategies

import (
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state/sale"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// SaleOrderStrategy estrategia para órdenes SALE
type SaleOrderStrategy struct {
	states             map[entities.OrderStatus]order_state.OrderState
	eventPublisher     ports.EventPublisher
	productRepo        ports.ProductRepository
	productVariantRepo ports.ProductVariantRepository
}

func NewSaleOrderStrategy(eventPublisher ports.EventPublisher, productRepo ports.ProductRepository, productVariantRepo ports.ProductVariantRepository) order_state.OrderStrategy {
	return &SaleOrderStrategy{
		eventPublisher:     eventPublisher,
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		states: map[entities.OrderStatus]order_state.OrderState{
			entities.OrderStatusCancelled: sale.NewCancelledState(productVariantRepo),
			entities.OrderStatusConfirmed: sale.NewConfirmedState(),
			entities.OrderStatusDelivered: sale.NewDeliveredState(productVariantRepo),
			entities.OrderStatusPending:   sale.NewPendingState(productVariantRepo),
		},
	}
}

func (s *SaleOrderStrategy) GetInitialStatus() entities.OrderStatus {
	return entities.OrderStatusPending
}

func (s *SaleOrderStrategy) GetAllowedTransitions(currentStatus entities.OrderStatus) []entities.OrderStatus {
	state := s.states[currentStatus]
	if state == nil {
		return []entities.OrderStatus{}
	}
	return state.GetAllowedTransitions()
}

func (s *SaleOrderStrategy) CanTransition(from, to entities.OrderStatus) bool {
	state := s.states[from]
	if state == nil {
		return false
	}
	return state.CanTransitionTo(to)
}

func (s *SaleOrderStrategy) GetState(status entities.OrderStatus) order_state.OrderState {
	return s.states[status]
}
