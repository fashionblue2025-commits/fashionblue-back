package order

import (
	"context"
	"errors"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/strategies"
)

type CreateOrderUseCase struct {
	orderRepo          ports.OrderRepository
	productRepo        ports.ProductRepository
	productVariantRepo ports.ProductVariantRepository
	eventPublisher     ports.EventPublisher
	strategies         map[entities.OrderType]order_state.OrderStrategy
}

func NewCreateOrderUseCase(
	orderRepo ports.OrderRepository,
	productRepo ports.ProductRepository,
	productVariantRepo ports.ProductVariantRepository,
	eventPublisher ports.EventPublisher,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:          orderRepo,
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

func (uc *CreateOrderUseCase) Execute(ctx context.Context, order *entities.Order) error {
	// Validar orden
	if err := order.Validate(); err != nil {
		return err
	}

	// Enriquecer items con información del producto si es necesario
	if err := uc.enrichOrderItems(ctx, order); err != nil {
		return err
	}

	// Obtener estrategia para el tipo de orden
	strategy := uc.getStrategy(order.Type)
	if strategy == nil {
		return errors.New("unsupported order type")
	}

	// Establecer valores por defecto
	if order.Status == "" {
		// Usar el estado inicial de la estrategia
		order.Status = strategy.GetInitialStatus()
	}
	if order.OrderDate.IsZero() {
		order.OrderDate = time.Now()
	}

	// Calcular total
	order.TotalAmount = order.CalculateTotal()

	// Crear orden
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return err
	}

	// Ejecutar OnEnter del estado inicial
	initialState := strategy.GetState(order.Status)
	if initialState != nil {
		if err := initialState.OnEnter(ctx, order, order_state.StateTransitionData{
			Publisher: uc.eventPublisher,
		}); err != nil {
			return err
		}
	}

	return nil
}

// getStrategy obtiene la estrategia para un tipo de orden
func (uc *CreateOrderUseCase) getStrategy(orderType entities.OrderType) order_state.OrderStrategy {
	return uc.strategies[orderType]
}

// enrichOrderItems enriquece los items de la orden buscando variantes existentes
// Si encuentra una variante que coincida (nombre + color + talla), asigna el ProductVariantID
func (uc *CreateOrderUseCase) enrichOrderItems(ctx context.Context, order *entities.Order) error {
	for i := range order.Items {
		item := &order.Items[i]

		// Validar que ProductName siempre esté presente
		if item.ProductName == "" {
			return errors.New("product name is required for all items")
		}

		// Buscar producto base por nombre
		products, err := uc.productRepo.List(ctx, map[string]interface{}{
			"name": item.ProductName,
		})

		// Si no existe el producto base, la variante se creará después (ProductVariantID = 0)
		if err != nil || len(products) == 0 {
			continue
		}

		product := &products[0]

		// Buscar variante específica por color y talla
		variant, err := uc.productVariantRepo.GetByProductAndAttributes(ctx, product.ID, item.Color, item.SizeID)
		if err != nil {
			// Variante no existe, se creará después (ProductVariantID = 0)
			continue
		}

		// Variante encontrada, asignar ID
		item.ProductVariantID = variant.ID

		// Si no se especificó precio, usar el precio de la variante
		if item.UnitPrice == 0 {
			item.UnitPrice = variant.UnitPrice
		}
	}
	return nil
}
