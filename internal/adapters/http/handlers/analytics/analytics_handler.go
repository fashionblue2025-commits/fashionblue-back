package analytics

import (
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/event_handlers"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

// AnalyticsHTTPHandler maneja las peticiones HTTP de analytics
type AnalyticsHTTPHandler struct {
	analyticsHandler *event_handlers.AnalyticsHandler
}

// NewAnalyticsHTTPHandler crea un nuevo handler HTTP de analytics
func NewAnalyticsHTTPHandler(analyticsHandler *event_handlers.AnalyticsHandler) *AnalyticsHTTPHandler {
	return &AnalyticsHTTPHandler{
		analyticsHandler: analyticsHandler,
	}
}

// GetMetrics retorna las métricas actuales del sistema
// GET /api/analytics/metrics
func (h *AnalyticsHTTPHandler) GetMetrics(c echo.Context) error {
	metrics := h.analyticsHandler.GetMetrics()

	// Convertir a un formato amigable para el frontend
	response := map[string]interface{}{
		"ordersByStatus": map[string]int{
			"total":         metrics.TotalOrders,
			"quote":         metrics.QuoteOrders,
			"approved":      metrics.ApprovedOrders,
			"manufacturing": metrics.ManufacturingOrders,
			"finished":      metrics.FinishedOrders,
			"delivered":     metrics.DeliveredOrders,
			"cancelled":     metrics.CancelledOrders,
		},
		"ordersByType": map[string]int{
			"custom":    metrics.CustomOrders,
			"inventory": metrics.InventoryOrders,
			"sale":      metrics.SaleOrders,
		},
		"revenue": map[string]interface{}{
			"total":             metrics.TotalRevenue,
			"averageOrderValue": metrics.AverageOrderValue,
		},
		"conversionRates": map[string]float64{
			"approvalRate":     metrics.ApprovalRate,
			"cancellationRate": metrics.CancellationRate,
			"completionRate":   metrics.CompletionRate,
		},
		"manufacturing": map[string]interface{}{
			"averageTime":        metrics.AverageManufacturingTime.String(),
			"averageTimeSeconds": metrics.AverageManufacturingTime.Seconds(),
		},
	}

	return c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// GetDashboardSummary retorna un resumen optimizado para el dashboard
// GET /api/analytics/dashboard
func (h *AnalyticsHTTPHandler) GetDashboardSummary(c echo.Context) error {
	metrics := h.analyticsHandler.GetMetrics()

	// KPIs principales para el dashboard
	summary := map[string]interface{}{
		"kpis": []map[string]interface{}{
			{
				"label":  "Total Revenue",
				"value":  metrics.TotalRevenue,
				"format": "currency",
				"icon":   "dollar-sign",
			},
			{
				"label":  "Orders Delivered",
				"value":  metrics.DeliveredOrders,
				"format": "number",
				"icon":   "check-circle",
			},
			{
				"label":  "Avg Order Value",
				"value":  metrics.AverageOrderValue,
				"format": "currency",
				"icon":   "trending-up",
			},
			{
				"label":  "Completion Rate",
				"value":  metrics.CompletionRate,
				"format": "percentage",
				"icon":   "percent",
			},
		},
		"ordersByStatus": map[string]int{
			"quote":         metrics.QuoteOrders,
			"approved":      metrics.ApprovedOrders,
			"manufacturing": metrics.ManufacturingOrders,
			"finished":      metrics.FinishedOrders,
			"delivered":     metrics.DeliveredOrders,
			"cancelled":     metrics.CancelledOrders,
		},
		"ordersByType": map[string]int{
			"custom":    metrics.CustomOrders,
			"inventory": metrics.InventoryOrders,
			"sale":      metrics.SaleOrders,
		},
		"alerts": h.generateAlerts(metrics),
	}

	return response.Success(c, 200, "Dashboard summary retrieved successfully", summary)
}

// generateAlerts genera alertas basadas en las métricas
func (h *AnalyticsHTTPHandler) generateAlerts(metrics event_handlers.OrderMetricsSnapshot) []map[string]interface{} {
	alerts := []map[string]interface{}{}

	// Alerta de alta tasa de cancelación
	if metrics.CancellationRate > 20 {
		alerts = append(alerts, map[string]interface{}{
			"type":    "warning",
			"message": "High cancellation rate detected",
			"value":   metrics.CancellationRate,
			"action":  "Review cancelled orders",
		})
	}

	// Alerta de baja tasa de aprobación
	if metrics.ApprovalRate < 50 && metrics.TotalOrders > 10 {
		alerts = append(alerts, map[string]interface{}{
			"type":    "info",
			"message": "Low approval rate",
			"value":   metrics.ApprovalRate,
			"action":  "Review quote process",
		})
	}

	// Alerta de órdenes en manufactura
	if metrics.ManufacturingOrders > 10 {
		alerts = append(alerts, map[string]interface{}{
			"type":    "info",
			"message": "High number of orders in manufacturing",
			"value":   metrics.ManufacturingOrders,
			"action":  "Monitor production capacity",
		})
	}

	return alerts
}
