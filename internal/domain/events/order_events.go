package events

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// OrderEvent representa un evento de orden
type OrderEvent struct {
	Type      OrderEventType
	OrderID   uint
	Order     *entities.Order
	OldStatus entities.OrderStatus
	NewStatus entities.OrderStatus
	Data      map[string]interface{} // Datos adicionales del evento
	Timestamp time.Time
}

// OrderEventType representa el tipo de evento
type OrderEventType string

const (
	// Eventos de transición de estado
	EventOrderStatusChanged OrderEventType = "order.status.changed"
	EventOrderApproved      OrderEventType = "order.approved"
	EventOrderManufacturing OrderEventType = "order.manufacturing"
	EventOrderFinished      OrderEventType = "order.finished"
	EventOrderDelivered     OrderEventType = "order.delivered"
	EventOrderCancelled     OrderEventType = "order.cancelled"

	// Eventos específicos de INVENTORY
	EventInventoryPlanned       OrderEventType = "inventory.planned"
	EventInventoryManufacturing OrderEventType = "inventory.manufacturing"
	EventInventoryFinished      OrderEventType = "inventory.finished"

	// Eventos específicos de SALE
	EventSalePending   OrderEventType = "sale.pending"
	EventSaleConfirmed OrderEventType = "sale.confirmed"
	EventSaleDelivered OrderEventType = "sale.delivered"

	// Eventos de acciones
	EventProductCreationRequired OrderEventType = "product.creation.required"
	EventStockUpdated            OrderEventType = "stock.updated"
	EventStockReserved           OrderEventType = "stock.reserved"
	EventStockReleased           OrderEventType = "stock.released"

	// Eventos contables y financieros
	EventInternalCustomerSaleCompleted OrderEventType = "internal.customer.sale.completed"
	EventSaleCompleted                 OrderEventType = "sale.completed" // Para registros financieros automáticos
)

// EventBus maneja la distribución de eventos
type EventBus struct {
	subscribers map[OrderEventType][]chan OrderEvent
}

// NewEventBus crea un nuevo bus de eventos
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[OrderEventType][]chan OrderEvent),
	}
}

// Subscribe suscribe un canal a un tipo de evento
func (eb *EventBus) Subscribe(eventType OrderEventType, ch chan OrderEvent) {
	if eb.subscribers[eventType] == nil {
		eb.subscribers[eventType] = []chan OrderEvent{}
	}
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
}

// Publish publica un evento a todos los suscriptores
func (eb *EventBus) Publish(event OrderEvent) {
	subscribers := eb.subscribers[event.Type]
	for _, ch := range subscribers {
		// Enviar de forma no bloqueante
		select {
		case ch <- event:
		default:
			// Si el canal está lleno, no bloquear
		}
	}

	// También publicar al canal genérico de cambios de estado
	if event.Type != EventOrderStatusChanged {
		genericEvent := event
		genericEvent.Type = EventOrderStatusChanged
		genericSubscribers := eb.subscribers[EventOrderStatusChanged]
		for _, ch := range genericSubscribers {
			select {
			case ch <- genericEvent:
			default:
			}
		}
	}
}

// Unsubscribe desuscribe un canal de un tipo de evento
func (eb *EventBus) Unsubscribe(eventType OrderEventType, ch chan OrderEvent) {
	subscribers := eb.subscribers[eventType]
	for i, subscriber := range subscribers {
		if subscriber == ch {
			eb.subscribers[eventType] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}
}

// Close cierra todos los canales suscritos
func (eb *EventBus) Close() {
	for _, subscribers := range eb.subscribers {
		for _, ch := range subscribers {
			close(ch)
		}
	}
}
