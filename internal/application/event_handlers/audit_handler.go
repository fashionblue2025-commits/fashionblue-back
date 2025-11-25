package event_handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

// AuditHandler maneja eventos para auditoría
type AuditHandler struct {
	eventBus   *events.EventBus
	eventChan  chan events.OrderEvent
	stopChan   chan bool
	repository ports.AuditLogRepository
}

// NewAuditHandler crea un nuevo handler de auditoría
func NewAuditHandler(eventBus *events.EventBus, repository ports.AuditLogRepository) *AuditHandler {
	handler := &AuditHandler{
		eventBus:   eventBus,
		eventChan:  make(chan events.OrderEvent, 100),
		stopChan:   make(chan bool),
		repository: repository,
	}

	// Suscribirse a todos los eventos para auditoría completa
	eventBus.Subscribe(events.EventOrderStatusChanged, handler.eventChan)

	return handler
}

// Start inicia el procesamiento de eventos
func (h *AuditHandler) Start() {
	log.Println("🔍 Audit Event Handler started")

	go func() {
		for {
			select {
			case event := <-h.eventChan:
				h.handleEvent(event)
			case <-h.stopChan:
				log.Println("🔍 Audit Event Handler stopped")
				return
			}
		}
	}()
}

// Stop detiene el procesamiento de eventos
func (h *AuditHandler) Stop() {
	h.stopChan <- true
}

// handleEvent procesa un evento para auditoría
func (h *AuditHandler) handleEvent(event events.OrderEvent) {
	// Log de auditoría en consola
	log.Printf("🔍 [AUDIT] Event: %s | Order: %d | Status: %s | Time: %s",
		event.Type,
		event.OrderID,
		event.NewStatus,
		event.Timestamp.Format("2006-01-02 15:04:05"),
	)

	// Guardar en base de datos si hay repositorio configurado
	if h.repository != nil {
		ctx := context.Background()

		// Crear el log de auditoría
		auditLog := &entities.AuditLog{
			EventType:   string(event.Type),
			OrderID:     event.OrderID,
			OldStatus:   string(event.OldStatus),
			NewStatus:   string(event.NewStatus),
			Description: h.generateDescription(event),
			CreatedAt:   event.Timestamp,
			Metadata:    "{}",
		}

		// Agregar información de la orden si está disponible
		if event.Order != nil {
			auditLog.OrderNumber = event.Order.OrderNumber
		}

		// Serializar metadata adicional si existe
		if event.Data != nil {
			if metadataJSON, err := json.Marshal(event.Data); err == nil {
				auditLog.Metadata = string(metadataJSON)
			}
		}

		// Guardar en BD
		if err := h.repository.Create(ctx, auditLog); err != nil {
			log.Printf("🔍 [AUDIT ERROR] Failed to save audit log: %v", err)
		} else {
			log.Printf("🔍 [AUDIT] ✅ Audit log saved to database (ID: %d)", auditLog.ID)
		}
	}

	// Log detallado para eventos críticos
	switch event.Type {
	case events.EventOrderApproved:
		log.Printf("🔍 [AUDIT] ⚠️  CRITICAL: Order #%d approved - requires tracking", event.OrderID)

	case events.EventOrderCancelled:
		log.Printf("🔍 [AUDIT] ⚠️  CRITICAL: Order #%d cancelled - investigate reason", event.OrderID)

	case events.EventStockUpdated:
		log.Printf("🔍 [AUDIT] 📦 Stock modification for order #%d - verify inventory", event.OrderID)

	case events.EventOrderDelivered:
		log.Printf("🔍 [AUDIT] ✅ Order #%d delivered successfully", event.OrderID)

	case events.EventProductCreationRequired:
		log.Printf("🔍 [AUDIT] 🏭 Product creation required for order #%d", event.OrderID)
	}
}

// generateDescription genera una descripción legible del evento
func (h *AuditHandler) generateDescription(event events.OrderEvent) string {
	switch event.Type {
	case events.EventOrderStatusChanged:
		return "Order status changed from " + string(event.OldStatus) + " to " + string(event.NewStatus)
	case events.EventOrderApproved:
		return "Order approved and ready for production"
	case events.EventOrderManufacturing:
		return "Order entered manufacturing phase"
	case events.EventOrderFinished:
		return "Order manufacturing completed"
	case events.EventOrderDelivered:
		return "Order delivered to customer"
	case events.EventOrderCancelled:
		return "Order cancelled"
	case events.EventInventoryPlanned:
		return "Inventory production planned"
	case events.EventInventoryManufacturing:
		return "Inventory production started"
	case events.EventInventoryFinished:
		return "Inventory production completed"
	case events.EventSalePending:
		return "Sale pending confirmation"
	case events.EventSaleConfirmed:
		return "Sale confirmed"
	case events.EventSaleDelivered:
		return "Sale delivered"
	case events.EventProductCreationRequired:
		return "Product creation required"
	case events.EventStockUpdated:
		return "Stock updated"
	case events.EventStockReserved:
		return "Stock reserved for order"
	case events.EventStockReleased:
		return "Stock released from order"
	case events.EventInternalCustomerSaleCompleted:
		return "Internal customer sale completed - transaction created"
	default:
		return "Order event: " + string(event.Type)
	}
}
