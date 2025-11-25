package event_handlers

import (
	"context"
	"log"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// InternalCustomerTransactionHandler maneja la creación de transacciones
// para clientes internos cuando se completa una venta
type InternalCustomerTransactionHandler struct {
	eventBus                *events.EventBus
	eventChan               chan events.OrderEvent
	stopChan                chan bool
	customerTransactionRepo ports.CustomerTransactionRepository
}

// NewInternalCustomerTransactionHandler crea un nuevo handler
func NewInternalCustomerTransactionHandler(
	eventBus *events.EventBus,
	customerTransactionRepo ports.CustomerTransactionRepository,
) *InternalCustomerTransactionHandler {
	handler := &InternalCustomerTransactionHandler{
		eventBus:                eventBus,
		eventChan:               make(chan events.OrderEvent, 100),
		stopChan:                make(chan bool),
		customerTransactionRepo: customerTransactionRepo,
	}
	eventBus.Subscribe(events.EventInternalCustomerSaleCompleted, handler.eventChan)

	return handler
}

// Start inicia el procesamiento de eventos
func (h *InternalCustomerTransactionHandler) Start() {
	log.Println("💰 Internal Customer Transaction Handler started")

	go func() {
		for {
			select {
			case event := <-h.eventChan:
				ctx := context.Background()
				if err := h.Handle(ctx, event); err != nil {
					log.Printf("❌ [ERROR] Failed to handle event: %v", err)
				}
			case <-h.stopChan:
				log.Println("💰 Internal Customer Transaction Handler stopped")
				return
			}
		}
	}()
}

// Stop detiene el procesamiento de eventos
func (h *InternalCustomerTransactionHandler) Stop() {
	h.stopChan <- true
}

// Handle procesa el evento de venta completada a cliente interno
func (h *InternalCustomerTransactionHandler) Handle(ctx context.Context, event events.OrderEvent) error {
	// Solo procesar si es el evento correcto
	if event.Type != events.EventInternalCustomerSaleCompleted {
		return nil
	}

	order := event.Order
	if order == nil {
		log.Printf("⚠️  [WARNING] Order is nil in event")
		return nil
	}

	// Verificar que sea cliente interno
	if !order.IsInternalCustomer() {
		log.Printf("⚠️  [WARNING] Order #%d is not for internal customer", order.ID)
		return nil
	}

	// Crear transacción de deuda
	transaction := &entities.CustomerTransaction{
		CustomerID:  *order.CustomerID,
		Type:        entities.TransactionTypeDebt,
		Amount:      order.TotalAmount - order.Discount,
		Description: buildTransactionDescription(order),
		Date:        time.Now(),
	}

	// Guardar transacción
	if err := h.customerTransactionRepo.Create(ctx, transaction); err != nil {
		log.Printf("❌ [ERROR] Failed to create transaction for customer #%d: %v", *order.CustomerID, err)
		return err
	}

	log.Printf("💰 [TRANSACTION] Created debt transaction for customer #%d: Order #%d - $%.2f",
		*order.CustomerID, order.ID, transaction.Amount)

	return nil
}

// buildTransactionDescription construye la descripción de la transacción
func buildTransactionDescription(order *entities.Order) string {
	description := "Venta - Orden #" + order.OrderNumber

	// Agregar tipo de orden
	switch order.Type {
	case entities.OrderTypeCustom:
		description += " (Producción por demanda)"
	case entities.OrderTypeInventory:
		description += " (Producción para stock)"
	case entities.OrderTypeSale:
		description += " (Venta de existente)"
	}

	// Agregar items
	if len(order.Items) > 0 {
		description += " - Items: "
		for i, item := range order.Items {
			if i > 0 {
				description += ", "
			}
			description += item.ProductName
			if item.Color != "" {
				description += " " + item.Color
			}
			if item.Quantity > 1 {
				description += " x" + string(rune(item.Quantity+'0'))
			}
		}
	}

	return description
}
