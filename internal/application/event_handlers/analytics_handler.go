package event_handlers

import (
	"log"
	"sync"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/events"
)

// OrderMetrics almacena métricas de órdenes (uso interno con mutex)
type OrderMetrics struct {
	// Contadores por estado
	TotalOrders         int
	QuoteOrders         int
	ApprovedOrders      int
	ManufacturingOrders int
	FinishedOrders      int
	DeliveredOrders     int
	CancelledOrders     int

	// Contadores por tipo
	CustomOrders    int
	InventoryOrders int
	SaleOrders      int

	// Métricas de ventas
	TotalRevenue      float64
	AverageOrderValue float64

	// Métricas de tiempo
	ManufacturingStartTime   map[uint]time.Time     // OrderID -> inicio manufactura
	ManufacturingDuration    map[uint]time.Duration // OrderID -> duración total
	AverageManufacturingTime time.Duration

	// Métricas de conversión
	ApprovalRate     float64 // Aprobadas / Total
	CancellationRate float64 // Canceladas / Aprobadas
	CompletionRate   float64 // Entregadas / Aprobadas

	mu sync.RWMutex
}

// OrderMetricsSnapshot representa una copia de las métricas sin el mutex (para exportar)
type OrderMetricsSnapshot struct {
	// Contadores por estado
	TotalOrders         int
	QuoteOrders         int
	ApprovedOrders      int
	ManufacturingOrders int
	FinishedOrders      int
	DeliveredOrders     int
	CancelledOrders     int

	// Contadores por tipo
	CustomOrders    int
	InventoryOrders int
	SaleOrders      int

	// Métricas de ventas
	TotalRevenue      float64
	AverageOrderValue float64

	// Métricas de tiempo
	AverageManufacturingTime time.Duration

	// Métricas de conversión
	ApprovalRate     float64
	CancellationRate float64
	CompletionRate   float64
}

// AnalyticsHandler maneja eventos para analytics y métricas
type AnalyticsHandler struct {
	eventBus  *events.EventBus
	eventChan chan events.OrderEvent
	stopChan  chan bool
	metrics   *OrderMetrics
}

// NewAnalyticsHandler crea un nuevo handler de analytics
func NewAnalyticsHandler(eventBus *events.EventBus) *AnalyticsHandler {
	handler := &AnalyticsHandler{
		eventBus:  eventBus,
		eventChan: make(chan events.OrderEvent, 100),
		stopChan:  make(chan bool),
		metrics: &OrderMetrics{
			ManufacturingStartTime: make(map[uint]time.Time),
			ManufacturingDuration:  make(map[uint]time.Duration),
		},
	}

	// Suscribirse a eventos relevantes para analytics
	eventBus.Subscribe(events.EventOrderStatusChanged, handler.eventChan)

	return handler
}

// Start inicia el procesamiento de eventos
func (h *AnalyticsHandler) Start() {
	log.Println("📊 Analytics Event Handler started")

	go func() {
		for {
			select {
			case event := <-h.eventChan:
				h.handleEvent(event)
			case <-h.stopChan:
				log.Println("📊 Analytics Event Handler stopped")
				h.printMetrics()
				return
			}
		}
	}()
}

// Stop detiene el procesamiento de eventos
func (h *AnalyticsHandler) Stop() {
	h.stopChan <- true
}

// handleEvent procesa un evento para analytics
func (h *AnalyticsHandler) handleEvent(event events.OrderEvent) {
	h.metrics.mu.Lock()
	defer h.metrics.mu.Unlock()

	order := event.Order

	// Incrementar contador total en el primer evento de la orden
	if event.NewStatus != "" {
		h.metrics.TotalOrders++
	}

	// Contadores por estado
	switch event.NewStatus {
	case "QUOTE":
		h.metrics.QuoteOrders++
	case "APPROVED":
		h.metrics.ApprovedOrders++
		log.Printf("📊 [ANALYTICS] Approved orders: %d", h.metrics.ApprovedOrders)
	case "MANUFACTURING":
		h.metrics.ManufacturingOrders++
		// Registrar inicio de manufactura
		h.metrics.ManufacturingStartTime[event.OrderID] = time.Now()
		log.Printf("📊 [ANALYTICS] Manufacturing started for order #%d", event.OrderID)
	case "FINISHED":
		h.metrics.FinishedOrders++
		// Calcular duración de manufactura
		if startTime, exists := h.metrics.ManufacturingStartTime[event.OrderID]; exists {
			duration := time.Since(startTime)
			h.metrics.ManufacturingDuration[event.OrderID] = duration
			h.calculateAverageManufacturingTime()
			log.Printf("📊 [ANALYTICS] Order #%d manufacturing completed in %s", event.OrderID, duration)
		}
	case "DELIVERED":
		h.metrics.DeliveredOrders++
		// Calcular revenue si hay información de la orden
		if order != nil {
			revenue := order.TotalAmount - order.Discount
			h.metrics.TotalRevenue += revenue
			log.Printf("📊 [ANALYTICS] Order #%d delivered - Revenue: $%.2f", event.OrderID, revenue)
		}
	case "CANCELLED":
		h.metrics.CancelledOrders++
		log.Printf("📊 [ANALYTICS] Cancelled orders: %d", h.metrics.CancelledOrders)
	}

	// Contadores por tipo de orden
	if order != nil {
		switch order.Type {
		case "CUSTOM":
			h.metrics.CustomOrders++
		case "INVENTORY":
			h.metrics.InventoryOrders++
		case "SALE":
			h.metrics.SaleOrders++
		}
	}

	// Recalcular métricas derivadas
	h.calculateRates()
}

// calculateRates calcula las tasas de conversión
// DEBE ser llamado con el mutex ya bloqueado
func (h *AnalyticsHandler) calculateRates() {
	if h.metrics.TotalOrders > 0 {
		h.metrics.ApprovalRate = float64(h.metrics.ApprovedOrders) / float64(h.metrics.TotalOrders) * 100
	}

	if h.metrics.ApprovedOrders > 0 {
		h.metrics.CancellationRate = float64(h.metrics.CancelledOrders) / float64(h.metrics.ApprovedOrders) * 100
		h.metrics.CompletionRate = float64(h.metrics.DeliveredOrders) / float64(h.metrics.ApprovedOrders) * 100
	}

	if h.metrics.DeliveredOrders > 0 {
		h.metrics.AverageOrderValue = h.metrics.TotalRevenue / float64(h.metrics.DeliveredOrders)
	}
}

// calculateAverageManufacturingTime calcula el tiempo promedio de manufactura
// DEBE ser llamado con el mutex ya bloqueado
func (h *AnalyticsHandler) calculateAverageManufacturingTime() {
	if len(h.metrics.ManufacturingDuration) == 0 {
		return
	}

	var total time.Duration
	for _, duration := range h.metrics.ManufacturingDuration {
		total += duration
	}

	h.metrics.AverageManufacturingTime = total / time.Duration(len(h.metrics.ManufacturingDuration))
}

// GetMetrics retorna una copia de las métricas actuales sin el mutex
func (h *AnalyticsHandler) GetMetrics() OrderMetricsSnapshot {
	h.metrics.mu.RLock()
	defer h.metrics.mu.RUnlock()

	// Retornar un snapshot sin el mutex
	return OrderMetricsSnapshot{
		TotalOrders:              h.metrics.TotalOrders,
		QuoteOrders:              h.metrics.QuoteOrders,
		ApprovedOrders:           h.metrics.ApprovedOrders,
		ManufacturingOrders:      h.metrics.ManufacturingOrders,
		FinishedOrders:           h.metrics.FinishedOrders,
		DeliveredOrders:          h.metrics.DeliveredOrders,
		CancelledOrders:          h.metrics.CancelledOrders,
		CustomOrders:             h.metrics.CustomOrders,
		InventoryOrders:          h.metrics.InventoryOrders,
		SaleOrders:               h.metrics.SaleOrders,
		TotalRevenue:             h.metrics.TotalRevenue,
		AverageOrderValue:        h.metrics.AverageOrderValue,
		AverageManufacturingTime: h.metrics.AverageManufacturingTime,
		ApprovalRate:             h.metrics.ApprovalRate,
		CancellationRate:         h.metrics.CancellationRate,
		CompletionRate:           h.metrics.CompletionRate,
	}
}

// printMetrics imprime un resumen completo de métricas
func (h *AnalyticsHandler) printMetrics() {
	h.metrics.mu.RLock()
	defer h.metrics.mu.RUnlock()

	log.Println("📊 ==================== ORDER METRICS SUMMARY ====================")
	log.Println("📊")
	log.Println("📊 📋 ORDERS BY STATUS:")
	log.Printf("📊   Total Orders:        %d", h.metrics.TotalOrders)
	log.Printf("📊   Quote:               %d", h.metrics.QuoteOrders)
	log.Printf("📊   Approved:            %d", h.metrics.ApprovedOrders)
	log.Printf("📊   Manufacturing:       %d", h.metrics.ManufacturingOrders)
	log.Printf("📊   Finished:            %d", h.metrics.FinishedOrders)
	log.Printf("📊   Delivered:           %d", h.metrics.DeliveredOrders)
	log.Printf("📊   Cancelled:           %d", h.metrics.CancelledOrders)
	log.Println("📊")
	log.Println("📊 🏷️  ORDERS BY TYPE:")
	log.Printf("📊   Custom:              %d", h.metrics.CustomOrders)
	log.Printf("📊   Inventory:           %d", h.metrics.InventoryOrders)
	log.Printf("📊   Sale:                %d", h.metrics.SaleOrders)
	log.Println("📊")
	log.Println("📊 💰 REVENUE METRICS:")
	log.Printf("📊   Total Revenue:       $%.2f", h.metrics.TotalRevenue)
	log.Printf("📊   Average Order Value: $%.2f", h.metrics.AverageOrderValue)
	log.Println("📊")
	log.Println("📊 📈 CONVERSION RATES:")
	log.Printf("📊   Approval Rate:       %.2f%%", h.metrics.ApprovalRate)
	log.Printf("📊   Cancellation Rate:   %.2f%%", h.metrics.CancellationRate)
	log.Printf("📊   Completion Rate:     %.2f%%", h.metrics.CompletionRate)
	log.Println("📊")
	log.Println("📊 ⏱️  MANUFACTURING TIME:")
	log.Printf("📊   Average Duration:    %s", h.metrics.AverageManufacturingTime)
	log.Printf("📊   Total Completed:     %d orders", len(h.metrics.ManufacturingDuration))
	log.Println("📊")
	log.Println("📊 ================================================================")
}
