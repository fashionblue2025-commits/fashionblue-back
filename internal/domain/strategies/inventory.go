package strategies

import (
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state/inventory"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// InventoryOrderStrategy estrategia para órdenes INVENTORY
type InventoryOrderStrategy struct {
	states             map[entities.OrderStatus]order_state.OrderState
	eventPublisher     ports.EventPublisher
	productRepo        ports.ProductRepository
	productVariantRepo ports.ProductVariantRepository
}

func NewInventoryOrderStrategy(eventPublisher ports.EventPublisher, productRepo ports.ProductRepository, productVariantRepo ports.ProductVariantRepository) order_state.OrderStrategy {
	return &InventoryOrderStrategy{
		eventPublisher:     eventPublisher,
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		states: map[entities.OrderStatus]order_state.OrderState{
			entities.OrderStatusManufacturing: inventory.NewManufacturingState(),
			entities.OrderStatusFinished:      inventory.NewFinishedState(productRepo),
			entities.OrderStatusCancelled:     inventory.NewCancelledState(productVariantRepo),
			entities.OrderStatusPlanned:       inventory.NewPlannedState(),
		},
	}
}

func (s *InventoryOrderStrategy) GetInitialStatus() entities.OrderStatus {
	return entities.OrderStatusPlanned
}

func (s *InventoryOrderStrategy) GetAllowedTransitions(currentStatus entities.OrderStatus) []entities.OrderStatus {
	state := s.states[currentStatus]
	if state == nil {
		return []entities.OrderStatus{}
	}
	return state.GetAllowedTransitions()
}

func (s *InventoryOrderStrategy) CanTransition(from, to entities.OrderStatus) bool {
	state := s.states[from]
	if state == nil {
		return false
	}
	return state.CanTransitionTo(to)
}

func (s *InventoryOrderStrategy) GetState(status entities.OrderStatus) order_state.OrderState {
	return s.states[status]
}
