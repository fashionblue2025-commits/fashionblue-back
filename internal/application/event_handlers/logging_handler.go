package event_handlers

import (
	"log"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
)

// LoggingHandler maneja eventos para logging
type LoggingHandler struct {
	eventBus  *events.EventBus
	eventChan chan events.OrderEvent
	stopChan  chan bool
}

// NewLoggingHandler crea un nuevo handler de logging
func NewLoggingHandler(eventBus *events.EventBus) *LoggingHandler {
	handler := &LoggingHandler{
		eventBus:  eventBus,
		eventChan: make(chan events.OrderEvent, 100),
		stopChan:  make(chan bool),
	}

	// Suscribirse a todos los eventos
	eventBus.Subscribe(events.EventOrderStatusChanged, handler.eventChan)

	return handler
}

// Start inicia el procesamiento de eventos
func (h *LoggingHandler) Start() {
	log.Println("📋 Logging Event Handler started")

	go func() {
		for {
			select {
			case event := <-h.eventChan:
				h.handleEvent(event)
			case <-h.stopChan:
				log.Println("📋 Logging Event Handler stopped")
				return
			}
		}
	}()
}

// Stop detiene el procesamiento de eventos
func (h *LoggingHandler) Stop() {
	h.stopChan <- true
}

// handleEvent procesa un evento
func (h *LoggingHandler) handleEvent(event events.OrderEvent) {
	// Log básico del evento
	log.Printf("📋 [EVENT] Type: %s | OrderID: %d | Status: %s -> %s | Time: %s",
		event.Type,
		event.OrderID,
		event.OldStatus,
		event.NewStatus,
		event.Timestamp.Format("2006-01-02 15:04:05"),
	)

	// Log específico por tipo de evento
	switch event.Type {
	case events.EventOrderApproved:
		log.Printf("✅ Order #%d has been approved", event.OrderID)

	case events.EventOrderManufacturing:
		log.Printf("🏭 Order #%d is now in manufacturing", event.OrderID)

	case events.EventOrderFinished:
		log.Printf("✨ Order #%d manufacturing finished", event.OrderID)

	case events.EventOrderDelivered:
		log.Printf("📦 Order #%d has been delivered", event.OrderID)

	case events.EventOrderCancelled:
		log.Printf("❌ Order #%d has been cancelled", event.OrderID)

	case events.EventStockUpdated:
		log.Printf("📊 Stock updated for order #%d", event.OrderID)

	case events.EventStockReserved:
		log.Printf("🔒 Stock reserved for order #%d", event.OrderID)

	case events.EventStockReleased:
		log.Printf("🔓 Stock released for order #%d", event.OrderID)

	case events.EventInventoryPlanned:
		log.Printf("📅 Inventory production planned for order #%d", event.OrderID)

	case events.EventSalePending:
		log.Printf("🛒 Sale pending for order #%d", event.OrderID)

	case events.EventSaleConfirmed:
		log.Printf("✅ Sale confirmed for order #%d", event.OrderID)
	}
}
