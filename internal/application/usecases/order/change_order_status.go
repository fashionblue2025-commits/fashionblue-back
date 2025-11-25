package order

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/strategies"
)

type ChangeOrderStatusUseCase struct {
	orderRepo          ports.OrderRepository
	orderItemRepo      ports.OrderItemRepository
	productRepo        ports.ProductRepository
	productVariantRepo ports.ProductVariantRepository
	eventPublisher     ports.EventPublisher
	strategies         map[entities.OrderType]order_state.OrderStrategy
}

func NewChangeOrderStatusUseCase(
	orderRepo ports.OrderRepository,
	orderItemRepo ports.OrderItemRepository,
	productRepo ports.ProductRepository,
	productVariantRepo ports.ProductVariantRepository,
	eventPublisher ports.EventPublisher,
) *ChangeOrderStatusUseCase {
	return &ChangeOrderStatusUseCase{
		orderRepo:          orderRepo,
		orderItemRepo:      orderItemRepo,
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		eventPublisher:     eventPublisher,
		strategies: map[entities.OrderType]order_state.OrderStrategy{
			entities.OrderTypeCustom:    strategies.NewCustomOrderStrategy(eventPublisher, productRepo, productVariantRepo),
			entities.OrderTypeInventory: strategies.NewInventoryOrderStrategy(eventPublisher, productRepo, productVariantRepo),
			entities.OrderTypeSale:      strategies.NewSaleOrderStrategy(eventPublisher, productRepo, productVariantRepo),
		},
	}
}

// OrderStatusChangeResult contiene el resultado del cambio de estado
type OrderStatusChangeResult struct {
	Order               *entities.Order
	AllowedNextStatuses []entities.OrderStatus
}

// Execute cambia el estado de una orden y ejecuta las acciones correspondientes
func (uc *ChangeOrderStatusUseCase) Execute(
	ctx context.Context,
	orderID uint,
	newStatus entities.OrderStatus,
	producedQuantities map[uint]int, // itemID -> cantidad producida (para FINISHED)
) (*OrderStatusChangeResult, error) {
	// Obtener orden con items
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Obtener estrategia para el tipo de orden
	strategy := uc.getStrategy(order.Type)
	if strategy == nil {
		return nil, errors.New("unsupported order type")
	}

	// Validar que no sea el mismo estado
	if order.Status == newStatus {
		return nil, errors.New("order is already in this status")
	}

	// Obtener estado actual y nuevo
	currentState := strategy.GetState(order.Status)
	newState := strategy.GetState(newStatus)

	if newState == nil {
		return nil, errors.New("invalid target status")
	}

	// Validar transición desde el estado actual
	if currentState != nil && !currentState.CanTransitionTo(newStatus) {
		return nil, errors.New("invalid state transition: current state does not allow this transition")
	}

	// Ejecutar OnExit del estado actual
	if currentState != nil {
		if err := currentState.OnExit(ctx, order, order_state.StateTransitionData{
			Publisher:          uc.eventPublisher,
			ProducedQuantities: producedQuantities,
		}); err != nil {
			return nil, err
		}
	}

	// Actualizar estado
	oldStatus := order.Status
	order.Status = newStatus

	// Ejecutar OnEnter del nuevo estado
	if err := newState.OnEnter(ctx, order, order_state.StateTransitionData{
		Publisher:          uc.eventPublisher,
		ProducedQuantities: producedQuantities,
		OldStatus:          oldStatus,
	}); err != nil {
		return nil, err
	}

	// Guardar cambios
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Verificar si hay una transición automática
	if nextStatus, shouldTransition := newState.DetermineNextState(ctx, order); shouldTransition {
		// Transición automática detectada, ejecutar recursivamente
		return uc.Execute(ctx, orderID, nextStatus, producedQuantities)
	}

	// Obtener estados permitidos desde el nuevo estado
	allowedNextStatuses := newState.GetAllowedTransitions()

	return &OrderStatusChangeResult{
		Order:               order,
		AllowedNextStatuses: allowedNextStatuses,
	}, nil
}

// GetAllowedNextStatuses obtiene los estados permitidos para una orden sin cambiar su estado
func (uc *ChangeOrderStatusUseCase) GetAllowedNextStatuses(ctx context.Context, orderID uint) ([]entities.OrderStatus, error) {
	// Obtener orden
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Obtener estrategia para el tipo de orden
	strategy := uc.getStrategy(order.Type)
	if strategy == nil {
		return nil, errors.New("unsupported order type")
	}

	// Obtener estado actual
	currentState := strategy.GetState(order.Status)
	if currentState == nil {
		return nil, errors.New("invalid current status")
	}

	// Obtener estados permitidos desde el estado actual
	return currentState.GetAllowedTransitions(), nil
}

// getStrategy obtiene la estrategia para un tipo de orden
func (uc *ChangeOrderStatusUseCase) getStrategy(orderType entities.OrderType) order_state.OrderStrategy {
	return uc.strategies[orderType]
}
