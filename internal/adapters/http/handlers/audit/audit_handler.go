package audit

import (
	"strconv"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

// AuditHTTPHandler maneja las peticiones HTTP de auditoría
type AuditHTTPHandler struct {
	repository ports.AuditLogRepository
}

// NewAuditHTTPHandler crea un nuevo handler HTTP de auditoría
func NewAuditHTTPHandler(repository ports.AuditLogRepository) *AuditHTTPHandler {
	return &AuditHTTPHandler{
		repository: repository,
	}
}

// GetAuditLogs obtiene logs de auditoría con filtros
// GET /api/v1/audit/logs
func (h *AuditHTTPHandler) GetAuditLogs(c echo.Context) error {
	// Parsear filtros
	filters := entities.AuditLogFilters{
		EventType: c.QueryParam("eventType"),
		Limit:     50, // Default
		Offset:    0,
	}

	// OrderID
	if orderIDStr := c.QueryParam("orderId"); orderIDStr != "" {
		if orderID, err := strconv.ParseUint(orderIDStr, 10, 32); err == nil {
			orderIDUint := uint(orderID)
			filters.OrderID = &orderIDUint
		}
	}

	// UserID
	if userIDStr := c.QueryParam("userId"); userIDStr != "" {
		if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userIDUint := uint(userID)
			filters.UserID = &userIDUint
		}
	}

	// StartDate
	if startDateStr := c.QueryParam("startDate"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = &startDate
		}
	}

	// EndDate
	if endDateStr := c.QueryParam("endDate"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = &endDate
		}
	}

	// Limit
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	// Offset
	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	// Obtener logs
	logs, err := h.repository.List(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve audit logs", err)
	}

	// Obtener total count
	total, err := h.repository.Count(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to count audit logs", err)
	}

	return response.Success(c, 200, "Audit logs retrieved successfully", map[string]interface{}{
		"logs":   logs,
		"total":  total,
		"limit":  filters.Limit,
		"offset": filters.Offset,
	})
}

// GetAuditLogByID obtiene un log específico por ID
// GET /api/v1/audit/logs/:id
func (h *AuditHTTPHandler) GetAuditLogByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid audit log ID", err)
	}

	log, err := h.repository.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Audit log not found")
	}

	return response.Success(c, 200, "Audit log retrieved successfully", log)
}

// GetAuditLogsByOrder obtiene todos los logs de una orden
// GET /api/v1/audit/orders/:orderId/logs
func (h *AuditHTTPHandler) GetAuditLogsByOrder(c echo.Context) error {
	orderID, err := strconv.ParseUint(c.Param("orderId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid order ID", err)
	}

	logs, err := h.repository.GetByOrderID(c.Request().Context(), uint(orderID))
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve audit logs", err)
	}

	return response.Success(c, 200, "Order audit logs retrieved successfully", map[string]interface{}{
		"orderId": orderID,
		"logs":    logs,
		"total":   len(logs),
	})
}

// GetAuditLogsByUser obtiene todos los logs de un usuario
// GET /api/v1/audit/users/:userId/logs
func (h *AuditHTTPHandler) GetAuditLogsByUser(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	logs, err := h.repository.GetByUserID(c.Request().Context(), uint(userID))
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve audit logs", err)
	}

	return response.Success(c, 200, "User audit logs retrieved successfully", map[string]interface{}{
		"userId": userID,
		"logs":   logs,
		"total":  len(logs),
	})
}

// GetAuditStats obtiene estadísticas de auditoría
// GET /api/v1/audit/stats
func (h *AuditHTTPHandler) GetAuditStats(c echo.Context) error {
	ctx := c.Request().Context()

	// Obtener total de logs
	total, err := h.repository.Count(ctx, entities.AuditLogFilters{})
	if err != nil {
		return response.InternalServerError(c, "Failed to get audit stats", err)
	}

	// Logs de hoy
	today := time.Now().Truncate(24 * time.Hour)
	todayLogs, err := h.repository.Count(ctx, entities.AuditLogFilters{
		StartDate: &today,
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to get today's audit stats", err)
	}

	// Logs de esta semana
	weekAgo := time.Now().AddDate(0, 0, -7)
	weekLogs, err := h.repository.Count(ctx, entities.AuditLogFilters{
		StartDate: &weekAgo,
	})
	if err != nil {
		return response.InternalServerError(c, "Failed to get week's audit stats", err)
	}

	// Eventos críticos (aprobaciones y cancelaciones)
	approvals, _ := h.repository.Count(ctx, entities.AuditLogFilters{
		EventType: "order.approved",
	})

	cancellations, _ := h.repository.Count(ctx, entities.AuditLogFilters{
		EventType: "order.cancelled",
	})

	stats := map[string]interface{}{
		"total": total,
		"today": todayLogs,
		"week":  weekLogs,
		"criticalEvents": map[string]interface{}{
			"approvals":     approvals,
			"cancellations": cancellations,
		},
	}

	return response.Success(c, 200, "Audit stats retrieved successfully", stats)
}
