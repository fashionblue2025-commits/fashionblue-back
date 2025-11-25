package event_handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// FinancialIncomeHandler maneja la creación automática de ingresos financieros
// cuando se completa una venta (para cualquier tipo de cliente)
type FinancialIncomeHandler struct {
	eventBus                 *events.EventBus
	eventChan                chan events.OrderEvent
	stopChan                 chan bool
	financialTransactionRepo ports.FinancialTransactionRepository
}

// NewFinancialIncomeHandler crea un nuevo handler
func NewFinancialIncomeHandler(
	eventBus *events.EventBus,
	financialTransactionRepo ports.FinancialTransactionRepository,
) *FinancialIncomeHandler {
	handler := &FinancialIncomeHandler{
		eventBus:                 eventBus,
		eventChan:                make(chan events.OrderEvent, 100),
		stopChan:                 make(chan bool),
		financialTransactionRepo: financialTransactionRepo,
	}
	eventBus.Subscribe(events.EventSaleCompleted, handler.eventChan)

	return handler
}

// Start inicia el procesamiento de eventos
func (h *FinancialIncomeHandler) Start() {
	log.Println("💵 Financial Income Handler started")

	go func() {
		for {
			select {
			case event := <-h.eventChan:
				ctx := context.Background()
				if err := h.Handle(ctx, event); err != nil {
					log.Printf("❌ [FINANCIAL INCOME ERROR] Failed to handle event: %v", err)
				}
			case <-h.stopChan:
				log.Println("💵 Financial Income Handler stopped")
				return
			}
		}
	}()
}

// Stop detiene el procesamiento de eventos
func (h *FinancialIncomeHandler) Stop() {
	h.stopChan <- true
}

// Handle procesa el evento de venta completada y crea el ingreso financiero
func (h *FinancialIncomeHandler) Handle(ctx context.Context, event events.OrderEvent) error {
	// Solo procesar si es el evento correcto
	if event.Type != events.EventSaleCompleted {
		return nil
	}

	order := event.Order
	if order == nil {
		log.Printf("⚠️  [WARNING] Order is nil in event")
		return nil
	}

	// Determinar la categoría de ingreso basada en el tipo de orden
	category := entities.FinancialTransactionCategorySales

	// Calcular el monto real (total - descuento)
	amount := order.TotalAmount - order.Discount
	if amount <= 0 {
		log.Printf("⚠️  [WARNING] Order #%d has zero or negative amount: $%.2f", order.ID, amount)
		return nil
	}

	// Crear la transacción financiera de ingreso
	transaction := &entities.FinancialTransaction{
		Type:        entities.FinancialTransactionTypeIncome,
		Category:    category,
		Amount:      amount,
		Description: buildFinancialDescription(order),
		Date:        time.Now(),
	}

	// Validar antes de guardar
	if err := transaction.Validate(); err != nil {
		log.Printf("❌ [VALIDATION ERROR] Invalid financial transaction: %v", err)
		return err
	}

	// Guardar transacción financiera
	if err := h.financialTransactionRepo.Create(ctx, transaction); err != nil {
		log.Printf("❌ [ERROR] Failed to create financial income for order #%d: %v", order.ID, err)
		return err
	}

	log.Printf("💵 [FINANCIAL INCOME] Created automatic income: Order #%d (%s) - $%.2f (Category: %s)",
		order.ID, order.Type, transaction.Amount, transaction.Category)

	return nil
}

// buildFinancialDescription construye la descripción para la transacción financiera
func buildFinancialDescription(order *entities.Order) string {
	description := fmt.Sprintf("Venta - Orden %s", order.OrderNumber)

	// Agregar tipo de orden
	switch order.Type {
	case entities.OrderTypeCustom:
		description += " (Producción por demanda)"
	case entities.OrderTypeInventory:
		description += " (Producción para stock)"
	case entities.OrderTypeSale:
		description += " (Venta de inventario)"
	}

	// Agregar información del cliente si existe
	if order.CustomerID != nil {
		description += fmt.Sprintf(" - Cliente #%d", *order.CustomerID)
	}

	// Agregar resumen de items
	if len(order.Items) > 0 {
		description += " - Items: "
		itemCount := len(order.Items)
		if itemCount > 3 {
			// Si hay muchos items, solo mostrar los primeros 3
			for i := 0; i < 3; i++ {
				if i > 0 {
					description += ", "
				}
				description += order.Items[i].ProductName
			}
			description += fmt.Sprintf(" (+%d más)", itemCount-3)
		} else {
			for i, item := range order.Items {
				if i > 0 {
					description += ", "
				}
				description += item.ProductName
				if item.Quantity > 1 {
					description += fmt.Sprintf(" x%d", item.Quantity)
				}
			}
		}
	}

	return description
}
