package order_state

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// BaseState implementación base para estados
type BaseState struct {
	order              *entities.Order
	Status             entities.OrderStatus
	AllowedTransitions []entities.OrderStatus
}

func NewBaseState(order *entities.Order, status entities.OrderStatus, allowedTransitions []entities.OrderStatus) *BaseState {
	return &BaseState{
		order:              order,
		Status:             status,
		AllowedTransitions: allowedTransitions,
	}
}

// GetStatus retorna el estado que representa
func (s *BaseState) GetStatus() entities.OrderStatus {
	return s.Status
}

// OnEnter implementación por defecto (no hace nada)
func (s *BaseState) OnEnter(ctx context.Context, order *entities.Order, data StateTransitionData) error {
	return nil
}

// OnExit implementación por defecto (no hace nada)

func (s *BaseState) OnExit(ctx context.Context, order *entities.Order, data StateTransitionData) error {
	return nil
}

// CanTransitionTo verifica si puede transicionar a otro estado
func (s *BaseState) CanTransitionTo(newStatus entities.OrderStatus) bool {
	for _, allowed := range s.AllowedTransitions {
		if allowed == newStatus {
			return true
		}
	}
	return false
}

// GetAllowedTransitions retorna los estados permitidos
func (s *BaseState) GetAllowedTransitions() []entities.OrderStatus {
	return s.AllowedTransitions
}

// DetermineNextState implementación por defecto (no hay transición automática)
// Los estados que necesiten lógica especial deben sobrescribir este método
func (s *BaseState) DetermineNextState(ctx context.Context, order *entities.Order) (entities.OrderStatus, bool) {
	// Por defecto, no hay transición automática
	return "", false
}
