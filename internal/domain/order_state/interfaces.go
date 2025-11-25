package order_state

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// OrderStrategy define el comportamiento para cada tipo de orden
type OrderStrategy interface {
	// GetInitialStatus retorna el estado inicial para este tipo de orden
	GetInitialStatus() entities.OrderStatus

	// GetAllowedTransitions retorna las transiciones permitidas desde un estado
	GetAllowedTransitions(currentStatus entities.OrderStatus) []entities.OrderStatus

	// CanTransition verifica si una transición es válida
	CanTransition(from, to entities.OrderStatus) bool

	// GetState retorna el estado correspondiente
	GetState(status entities.OrderStatus) OrderState
}

// OrderState representa un estado específico de una orden
type OrderState interface {
	// GetStatus retorna el estado que representa
	GetStatus() entities.OrderStatus

	// OnEnter se ejecuta al entrar a este estado
	OnEnter(ctx context.Context, order *entities.Order, data StateTransitionData) error

	// OnExit se ejecuta al salir de este estado
	OnExit(ctx context.Context, order *entities.Order, data StateTransitionData) error

	// CanTransitionTo verifica si puede transicionar a otro estado
	CanTransitionTo(newStatus entities.OrderStatus) bool

	// GetAllowedTransitions retorna los estados permitidos desde este estado
	GetAllowedTransitions() []entities.OrderStatus

	// DetermineNextState determina automáticamente el siguiente estado basado en las condiciones de la orden
	// Retorna el siguiente estado y un booleano indicando si debe transicionar automáticamente
	DetermineNextState(ctx context.Context, order *entities.Order) (entities.OrderStatus, bool)
}

// StateTransitionData contiene datos para la transición de estado
type StateTransitionData struct {
	ProducedQuantities map[uint]int         // itemID -> cantidad producida
	Publisher          ports.EventPublisher // EventPublisher
	Context            context.Context      // Contexto de la transición
	Repositories       *RepositoryContainer // Repositorios necesarios
	OldStatus          entities.OrderStatus // Estado anterior (para referencia)
}

// RepositoryContainer contiene los repositorios necesarios para las transiciones
type RepositoryContainer struct {
	ProductRepo   ProductRepository
	OrderItemRepo OrderItemRepository
}

// ProductRepository define las operaciones necesarias de productos
type ProductRepository interface {
	GetByID(ctx context.Context, id uint) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Create(ctx context.Context, product *entities.Product) error
}

// OrderItemRepository define las operaciones necesarias de items de orden
type OrderItemRepository interface {
	GetByOrderID(ctx context.Context, orderID uint) ([]entities.OrderItem, error)
}
