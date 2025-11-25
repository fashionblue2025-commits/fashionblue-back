# Event Handlers (Consumers)

Este paquete contiene los event handlers (consumers) que procesan eventos de órdenes de forma asíncrona.

## Arquitectura

```
EventBus (Publisher) → Event Handlers (Consumers)
                    ↓
            [Logging, Notifications, Analytics, Audit, Webhooks]
```

## Handlers Disponibles

### 1. LoggingHandler
**Propósito:** Registra todos los eventos en los logs del sistema.

**Eventos suscritos:**
- Todos los eventos (`EventOrderStatusChanged`)

**Uso:**
```go
loggingHandler := event_handlers.NewLoggingHandler(eventBus)
loggingHandler.Start()
```

**Salida de ejemplo:**
```
📋 [EVENT] Type: order.approved | OrderID: 123 | Status: quote -> approved
✅ Order #123 has been approved
```

---

### 2. NotificationHandler
**Propósito:** Envía notificaciones a clientes sobre eventos importantes.

**Eventos suscritos:**
- `EventOrderApproved`
- `EventOrderDelivered`
- `EventOrderCancelled`
- `EventSaleConfirmed`

**Uso:**
```go
notificationHandler := event_handlers.NewNotificationHandler(eventBus)
notificationHandler.Start()
```

**Extensiones posibles:**
- Enviar emails
- Enviar SMS
- Push notifications
- Webhooks a sistemas externos

---

### 3. AnalyticsHandler
**Propósito:** Recopila métricas y estadísticas de órdenes.

**Eventos suscritos:**
- Todos los eventos (`EventOrderStatusChanged`)

**Métricas rastreadas:**
- Total de órdenes aprobadas
- Total de órdenes canceladas
- Total de órdenes entregadas
- Tasa de cancelación

**Uso:**
```go
analyticsHandler := event_handlers.NewAnalyticsHandler(eventBus)
analyticsHandler.Start()

// Obtener métricas
metrics := analyticsHandler.GetMetrics()
fmt.Printf("Approved: %d, Cancelled: %d\n", metrics.ApprovedOrders, metrics.CancelledOrders)
```

---

### 4. AuditHandler
**Propósito:** Mantiene un registro de auditoría de todos los eventos.

**Eventos suscritos:**
- Todos los eventos (`EventOrderStatusChanged`)

**Uso:**
```go
// Sin repositorio (solo logs)
auditHandler := event_handlers.NewAuditHandler(eventBus, nil)
auditHandler.Start()

// Con repositorio (guarda en BD)
auditHandler := event_handlers.NewAuditHandler(eventBus, auditRepository)
auditHandler.Start()
```

**Casos de uso:**
- Compliance y regulaciones
- Investigación de incidentes
- Trazabilidad completa

---

### 5. WebhookHandler
**Propósito:** Envía webhooks HTTP a sistemas externos.

**Eventos suscritos:**
- `EventOrderApproved`
- `EventOrderDelivered`
- `EventOrderCancelled`
- `EventStockUpdated`

**Configuración:**
```go
webhookConfig := event_handlers.WebhookConfig{
    URL:     "https://api.example.com/webhooks",
    Enabled: true,
    Secret:  "your-webhook-secret",
}
webhookHandler := event_handlers.NewWebhookHandler(eventBus, webhookConfig)
webhookHandler.Start()
```

**Payload de ejemplo:**
```json
{
  "event_type": "order.approved",
  "order_id": 123,
  "old_status": "quote",
  "new_status": "approved",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {}
}
```

---

## Flujo de Eventos

### Ejemplo: Cambio de estado de orden

```
1. Usuario cambia estado de orden (QUOTE → APPROVED)
   ↓
2. ChangeOrderStatusUseCase ejecuta OnExit del estado actual
   ↓
3. ChangeOrderStatusUseCase ejecuta OnEnter del nuevo estado
   ↓
4. Estado publica evento a través del EventBus
   ↓
5. EventBus distribuye el evento a todos los handlers suscritos
   ↓
6. Handlers procesan el evento de forma asíncrona:
   - LoggingHandler: Registra en logs
   - NotificationHandler: Envía notificación al cliente
   - AnalyticsHandler: Actualiza métricas
   - AuditHandler: Guarda en registro de auditoría
   - WebhookHandler: Envía webhook a sistema externo
```

---

## Agregar un Nuevo Handler

1. Crear archivo en `internal/application/event_handlers/`:

```go
package event_handlers

import (
    "log"
    "github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
)

type MyCustomHandler struct {
    eventBus  *events.EventBus
    eventChan chan events.OrderEvent
    stopChan  chan bool
}

func NewMyCustomHandler(eventBus *events.EventBus) *MyCustomHandler {
    handler := &MyCustomHandler{
        eventBus:  eventBus,
        eventChan: make(chan events.OrderEvent, 100),
        stopChan:  make(chan bool),
    }
    
    // Suscribirse a eventos
    eventBus.Subscribe(events.EventOrderApproved, handler.eventChan)
    
    return handler
}

func (h *MyCustomHandler) Start() {
    log.Println("My Custom Handler started")
    
    go func() {
        for {
            select {
            case event := <-h.eventChan:
                h.handleEvent(event)
            case <-h.stopChan:
                log.Println("My Custom Handler stopped")
                return
            }
        }
    }()
}

func (h *MyCustomHandler) Stop() {
    h.stopChan <- true
}

func (h *MyCustomHandler) handleEvent(event events.OrderEvent) {
    // Tu lógica aquí
    log.Printf("Processing event: %s for order %d", event.Type, event.OrderID)
}
```

2. Registrar en `cmd/api/main.go`:

```go
myCustomHandler := event_handlers.NewMyCustomHandler(eventBus)
myCustomHandler.Start()

// En shutdown:
myCustomHandler.Stop()
```

---

## Mejores Prácticas

1. **Idempotencia:** Los handlers deben ser idempotentes (procesar el mismo evento múltiples veces no debe causar problemas)

2. **Error Handling:** Siempre manejar errores sin detener el handler

3. **Buffer Size:** Ajustar el tamaño del buffer del canal según el volumen esperado

4. **Timeouts:** Implementar timeouts para operaciones externas (HTTP, DB, etc.)

5. **Logging:** Siempre loguear errores y eventos importantes

6. **Graceful Shutdown:** Implementar `Stop()` para cerrar limpiamente

---

## Testing

Ejemplo de test para un handler:

```go
func TestLoggingHandler(t *testing.T) {
    eventBus := events.NewEventBus()
    handler := event_handlers.NewLoggingHandler(eventBus)
    handler.Start()
    defer handler.Stop()
    
    // Publicar evento de prueba
    eventBus.Publish(events.OrderEvent{
        Type:      events.EventOrderApproved,
        OrderID:   123,
        NewStatus: entities.OrderStatusApproved,
        Timestamp: time.Now(),
    })
    
    // Esperar procesamiento
    time.Sleep(100 * time.Millisecond)
    
    // Verificar que se procesó (revisar logs o métricas)
}
```

---

## Configuración en Producción

### Variables de entorno recomendadas:

```bash
# Webhooks
WEBHOOK_URL=https://api.example.com/webhooks
WEBHOOK_SECRET=your-secret-key
WEBHOOK_ENABLED=true

# Analytics
ANALYTICS_ENABLED=true

# Notifications
NOTIFICATION_EMAIL_ENABLED=true
NOTIFICATION_SMS_ENABLED=false
```

### Monitoreo:

- Monitorear el tamaño de los canales de eventos
- Alertar si los handlers se detienen inesperadamente
- Rastrear latencia de procesamiento de eventos
- Monitorear tasa de errores en handlers
