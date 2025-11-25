package order

import (
	"strconv"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/dto"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/order"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	createOrderUC       *order.CreateOrderUseCase
	getOrderUC          *order.GetOrderUseCase
	listOrdersUC        *order.ListOrdersUseCase
	updateOrderStatusUC *order.UpdateOrderStatusUseCase
	addOrderItemUC      *order.AddOrderItemUseCase
	updateOrderItemUC   *order.UpdateOrderItemUseCase
	removeOrderItemUC   *order.RemoveOrderItemUseCase
	changeOrderStatusUC *order.ChangeOrderStatusUseCase
}

func NewOrderHandler(
	createOrderUC *order.CreateOrderUseCase,
	getOrderUC *order.GetOrderUseCase,
	listOrdersUC *order.ListOrdersUseCase,
	updateOrderStatusUC *order.UpdateOrderStatusUseCase,
	addOrderItemUC *order.AddOrderItemUseCase,
	updateOrderItemUC *order.UpdateOrderItemUseCase,
	removeOrderItemUC *order.RemoveOrderItemUseCase,
	changeOrderStatusUC *order.ChangeOrderStatusUseCase,
) *OrderHandler {
	return &OrderHandler{
		createOrderUC:       createOrderUC,
		getOrderUC:          getOrderUC,
		listOrdersUC:        listOrdersUC,
		updateOrderStatusUC: updateOrderStatusUC,
		addOrderItemUC:      addOrderItemUC,
		updateOrderItemUC:   updateOrderItemUC,
		removeOrderItemUC:   removeOrderItemUC,
		changeOrderStatusUC: changeOrderStatusUC,
	}
}

// CreateOrder crea una nueva orden/cotización en el sistema
// Valida que los productos existan y calcula el total automáticamente
// @Request: Order
// @Response: Order
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var req struct {
		CustomerID            *uint              `json:"customerId"` // ID del cliente interno (opcional)
		CustomerName          string             `json:"customerName"`
		SellerID              uint               `json:"sellerId"`
		Type                  entities.OrderType `json:"type"`
		Discount              float64            `json:"discount"`
		Notes                 string             `json:"notes"`
		EstimatedDeliveryDate *time.Time         `json:"estimatedDeliveryDate"`
		Items                 []struct {
			ProductID   uint    `json:"productId"`   // Opcional para CUSTOM/INVENTORY
			ProductName string  `json:"productName"` // Requerido para CUSTOM/INVENTORY
			Color       string  `json:"color"`
			SizeID      *uint   `json:"sizeId"`
			CategoryID  uint    `json:"categoryId"`
			Quantity    int     `json:"quantity"`
			UnitPrice   float64 `json:"unitPrice"`
		} `json:"items"`
	}

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	// Crear orden
	orderEntity := &entities.Order{
		CustomerID:            req.CustomerID,
		CustomerName:          req.CustomerName,
		SellerID:              req.SellerID,
		Type:                  req.Type,
		Discount:              req.Discount,
		Notes:                 req.Notes,
		EstimatedDeliveryDate: req.EstimatedDeliveryDate,
		OrderDate:             time.Now(),
	}

	// Agregar items
	orderEntity.Items = make([]entities.OrderItem, len(req.Items))
	for i, item := range req.Items {
		orderEntity.Items[i] = entities.OrderItem{
			ProductVariantID: item.ProductID, // Puede ser 0 para variantes nuevas
			ProductName:      item.ProductName,
			Color:            item.Color,
			SizeID:           item.SizeID,
			CategoryID:       item.CategoryID,
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
		}
		orderEntity.Items[i].CalculateSubtotal()
	}

	if err := h.createOrderUC.Execute(c.Request().Context(), orderEntity); err != nil {
		return response.BadRequest(c, "Failed to create order", err)
	}

	return response.Created(c, "Order created successfully", dto.ToOrderDTO(orderEntity))
}

// GetOrder obtiene una orden por su ID incluyendo todos sus items
// @Response: Order
func (h *OrderHandler) GetOrder(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid order ID", err)
	}

	order, err := h.getOrderUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Order not found")
	}

	return response.OK(c, "Order retrieved successfully", dto.ToOrderDTO(order))
}

// ListOrders lista todas las órdenes con filtros opcionales
// Soporta filtros por: status, seller_id, type, start_date, end_date
// @Response: Order
func (h *OrderHandler) ListOrders(c echo.Context) error {
	filters := make(map[string]interface{})

	// Filtros opcionales
	if status := c.QueryParam("status"); status != "" {
		filters["status"] = status
	}
	if sellerID := c.QueryParam("seller_id"); sellerID != "" {
		if id, err := strconv.ParseUint(sellerID, 10, 32); err == nil {
			filters["seller_id"] = uint(id)
		}
	}
	if orderType := c.QueryParam("type"); orderType != "" {
		filters["type"] = orderType
	}
	if startDate := c.QueryParam("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_date"] = t
		}
	}
	if endDate := c.QueryParam("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filters["end_date"] = t
		}
	}

	orders, err := h.listOrdersUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to list orders", err)
	}

	return response.OK(c, "Orders retrieved successfully", dto.ToOrderDTOList(orders))
}

// AddOrderItem agrega un nuevo item a una orden existente
// Valida que el producto exista y recalcula el total de la orden
// @Request: OrderItem
// @Response: OrderItem
func (h *OrderHandler) AddOrderItem(c echo.Context) error {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid order ID", err)
	}

	var req struct {
		ProductID uint    `json:"productId"`
		Color     string  `json:"color"`
		SizeID    *uint   `json:"sizeId"`
		Quantity  int     `json:"quantity"`
		UnitPrice float64 `json:"unitPrice"`
	}

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	item := &entities.OrderItem{
		OrderID:          uint(orderID),
		ProductVariantID: req.ProductID,
		Color:            req.Color,
		SizeID:           req.SizeID,
		Quantity:         req.Quantity,
		UnitPrice:        req.UnitPrice,
	}

	if err := h.addOrderItemUC.Execute(c.Request().Context(), item); err != nil {
		return response.BadRequest(c, "Failed to add order item", err)
	}

	return response.Created(c, "Order item added successfully", dto.ToOrderItemDTO(item))
}

// UpdateOrderItem actualiza un item existente de una orden
// Recalcula el subtotal y el total de la orden
// @Request: OrderItem
// @Response: OrderItem
func (h *OrderHandler) UpdateOrderItem(c echo.Context) error {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID", err)
	}

	var req struct {
		ProductID uint    `json:"productId"`
		Color     string  `json:"color"`
		SizeID    *uint   `json:"sizeId"`
		Quantity  int     `json:"quantity"`
		UnitPrice float64 `json:"unitPrice"`
	}

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	item := &entities.OrderItem{
		ID:               uint(itemID),
		ProductVariantID: req.ProductID,
		Color:            req.Color,
		SizeID:           req.SizeID,
		Quantity:         req.Quantity,
		UnitPrice:        req.UnitPrice,
	}

	if err := h.updateOrderItemUC.Execute(c.Request().Context(), item); err != nil {
		return response.BadRequest(c, "Failed to update order item", err)
	}

	return response.OK(c, "Order item updated successfully", dto.ToOrderItemDTO(item))
}

// RemoveOrderItem elimina un item de una orden
// Recalcula el total de la orden después de eliminar el item
func (h *OrderHandler) RemoveOrderItem(c echo.Context) error {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID", err)
	}

	if err := h.removeOrderItemUC.Execute(c.Request().Context(), uint(itemID)); err != nil {
		return response.BadRequest(c, "Failed to remove order item", err)
	}

	return response.OK(c, "Order item removed successfully", nil)
}

// ChangeOrderStatus cambia el estado de una orden (endpoint unificado)
// Maneja todas las transiciones de estado y ejecuta las acciones correspondientes
// Para órdenes CUSTOM en estado APPROVED, permite especificar cantidades producidas
// @Request: UpdateOrderStatusRequest
// @Response: Order
func (h *OrderHandler) ChangeOrderStatus(c echo.Context) error {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid order ID", err)
	}

	var req struct {
		Status             string       `json:"status"`
		ProducedQuantities map[uint]int `json:"producedQuantities,omitempty"` // itemID -> cantidad
	}

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if req.Status == "" {
		return response.BadRequest(c, "Status is required", nil)
	}

	newStatus := entities.OrderStatus(req.Status)

	result, err := h.changeOrderStatusUC.Execute(
		c.Request().Context(),
		uint(orderID),
		newStatus,
		req.ProducedQuantities,
	)
	if err != nil {
		return response.BadRequest(c, "Failed to change order status", err)
	}

	// Convertir orden a DTO con estados permitidos
	orderDTO := dto.ToOrderDTO(result.Order)

	// Agregar estados permitidos a la respuesta
	responseData := map[string]interface{}{
		"order":               orderDTO,
		"allowedNextStatuses": result.AllowedNextStatuses,
	}

	return response.OK(c, "Order status changed successfully", responseData)
}

// GetAllowedNextStatuses obtiene los estados permitidos para una orden
func (h *OrderHandler) GetAllowedNextStatuses(c echo.Context) error {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid order ID", err)
	}

	allowedStatuses, err := h.changeOrderStatusUC.GetAllowedNextStatuses(c.Request().Context(), uint(orderID))
	if err != nil {
		return response.InternalServerError(c, "Failed to get allowed statuses", err)
	}

	return response.OK(c, "Allowed statuses retrieved successfully", map[string]interface{}{
		"allowedNextStatuses": allowedStatuses,
	})
}
