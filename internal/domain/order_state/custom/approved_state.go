package custom

import (
	"context"
	"log"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/order_state"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ApprovedState struct {
	*order_state.BaseState
	productRepo        ports.ProductRepository
	productVariantRepo ports.ProductVariantRepository
}

func NewApprovedState(productRepo ports.ProductRepository, productVariantRepo ports.ProductVariantRepository) order_state.OrderState {
	return &ApprovedState{
		BaseState: &order_state.BaseState{
			Status: entities.OrderStatusApproved,
			AllowedTransitions: []entities.OrderStatus{
				entities.OrderStatusManufacturing,
				entities.OrderStatusCancelled,
			},
		},
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
	}
}

func (s *ApprovedState) OnEnter(ctx context.Context, order *entities.Order, data order_state.StateTransitionData) error {
	// 🔒 Reservar stock de productos existentes
	if s.productVariantRepo != nil {
		if err := s.reserveStockForItems(ctx, order); err != nil {
			return err
		}
	}

	// Ajustar transiciones permitidas según si necesita fabricación
	s.updateAllowedTransitions(order)

	// Publicar evento de orden aprobada
	if data.Publisher != nil {
		data.Publisher.Publish(events.OrderEvent{
			Type:      events.EventOrderApproved,
			OrderID:   order.ID,
			Order:     order,
			NewStatus: entities.OrderStatusApproved,
			Data: map[string]interface{}{
				"orderType":          order.Type,
				"needsManufacturing": order.NeedsManufacturing(),
			},
		})
	}
	return nil
}

// reserveStockForItems reserva stock disponible para cada item de la orden
func (s *ApprovedState) reserveStockForItems(ctx context.Context, order *entities.Order) error {
	for i := range order.Items {
		item := &order.Items[i]

		// Solo reservar si la variante ya existe
		if item.IsNewVariant() {
			continue
		}

		if err := s.reserveStockForItem(ctx, item); err != nil {
			return err
		}
	}
	return nil
}

// reserveStockForItem reserva stock disponible para un item específico
func (s *ApprovedState) reserveStockForItem(ctx context.Context, item *entities.OrderItem) error {
	// Obtener la variante
	variant, err := s.productVariantRepo.GetByID(ctx, item.ProductVariantID)
	if err != nil {
		// Variante no encontrada, se creará en FINISHED
		log.Printf("⚠️  [SKIP] Variant #%d not found, will be created later", item.ProductVariantID)
		item.ReservedQuantity = 0 // No hay stock reservado
		return nil
	}

	// Calcular cuánto stock disponible tenemos
	availableStock := variant.GetAvailableStock()
	if availableStock <= 0 {
		// No hay stock disponible
		log.Printf("📦 [NO STOCK] Variant #%d: No available stock (Stock: %d, Reserved: %d)",
			variant.ID, variant.Stock, variant.ReservedStock)
		item.ReservedQuantity = 0 // No hay stock reservado
		return nil
	}

	// Reservar lo que podamos del stock existente
	reserveQty := item.Quantity
	if availableStock < item.Quantity {
		reserveQty = availableStock
		log.Printf("⚠️  [PARTIAL] Variant #%d: Requested %d, available %d, reserving %d",
			variant.ID, item.Quantity, availableStock, reserveQty)
	}

	// Reservar stock en la variante
	if err := s.productVariantRepo.ReserveStock(ctx, variant.ID, reserveQty); err != nil {
		log.Printf("❌ [ERROR] Failed to reserve stock for variant #%d: %v", variant.ID, err)
		return err
	}

	// Guardar la cantidad reservada en el item
	item.ReservedQuantity = reserveQty

	log.Printf("🔒 [RESERVED] Variant #%d: Reserved %d units (Requested: %d, To manufacture: %d)",
		variant.ID, reserveQty, item.Quantity, item.Quantity-reserveQty)

	return nil
}

// updateAllowedTransitions ajusta las transiciones permitidas según si la orden necesita fabricación
func (s *ApprovedState) updateAllowedTransitions(order *entities.Order) {
	if order.HasFullStockCoverage() {
		// Si todo está cubierto por stock, saltar MANUFACTURING
		s.AllowedTransitions = []entities.OrderStatus{
			entities.OrderStatusFinished,
			entities.OrderStatusCancelled,
		}
	} else {
		// Si necesita fabricación, flujo normal
		s.AllowedTransitions = []entities.OrderStatus{
			entities.OrderStatusManufacturing,
			entities.OrderStatusCancelled,
		}
	}
}

// DetermineNextState determina automáticamente el siguiente estado basado en el stock disponible
func (s *ApprovedState) DetermineNextState(ctx context.Context, order *entities.Order) (entities.OrderStatus, bool) {
	// Si todos los items tienen stock completo, saltar MANUFACTURING
	if order.HasFullStockCoverage() {
		// Todo está en stock, ir directo a FINISHED
		log.Printf("✅ [AUTO-TRANSITION] Order #%d has full stock coverage, skipping MANUFACTURING → FINISHED", order.ID)
		return entities.OrderStatusFinished, true
	}

	// No hay transición automática (caso raro, pero seguro)
	return "", false
}
