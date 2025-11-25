package entities

import (
	"errors"
	"time"
)

// OrderStatus representa el estado de una orden
type OrderStatus string

const (
	// Estados para CUSTOM (producción por demanda)
	OrderStatusQuote         OrderStatus = "QUOTE"
	OrderStatusApproved      OrderStatus = "APPROVED"
	OrderStatusManufacturing OrderStatus = "MANUFACTURING"
	OrderStatusFinished      OrderStatus = "FINISHED"
	OrderStatusDelivered     OrderStatus = "DELIVERED"
	OrderStatusCancelled     OrderStatus = "CANCELLED"

	// Estados para INVENTORY (producción para stock)
	OrderStatusPlanned OrderStatus = "PLANNED"
	// Usa: MANUFACTURING, FINISHED, CANCELLED

	// Estados para SALE (venta de existente)
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	// Usa: DELIVERED, CANCELLED

	// Deprecated: usar MANUFACTURING
	OrderStatusInProduction OrderStatus = "IN_PRODUCTION"
)

// OrderType representa el tipo de orden
type OrderType string

const (
	OrderTypeCustom    OrderType = "CUSTOM"    // Producción por demanda (cotización)
	OrderTypeInventory OrderType = "INVENTORY" // Producción para stock
	OrderTypeSale      OrderType = "SALE"      // Venta de producto existente
)

// Order representa una orden de venta
type Order struct {
	ID                    uint
	OrderNumber           string
	CustomerID            *uint // ID del cliente interno (opcional, para registro contable)
	CustomerName          string
	SellerID              uint
	Seller                *User
	Type                  OrderType
	Status                OrderStatus
	TotalAmount           float64
	Discount              float64
	Notes                 string
	OrderDate             time.Time
	EstimatedDeliveryDate *time.Time
	ActualDeliveryDate    *time.Time
	Items                 []OrderItem
	Photos                []OrderPhoto
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// Validate valida los datos de la orden
func (o *Order) Validate() error {
	if o.CustomerName == "" {
		return errors.New("customer name is required")
	}
	if o.SellerID == 0 {
		return errors.New("seller is required")
	}
	// Validar tipo de orden
	if o.Type != OrderTypeCustom && o.Type != OrderTypeInventory && o.Type != OrderTypeSale {
		return errors.New("invalid order type: must be CUSTOM, INVENTORY, or SALE")
	}
	if o.TotalAmount < 0 {
		return errors.New("total amount cannot be negative")
	}
	if o.Discount < 0 {
		return errors.New("discount cannot be negative")
	}
	return nil
}

// CanEditItems verifica si se pueden editar los items de la orden
func (o *Order) CanEditItems() bool {
	return o.Status == OrderStatusQuote || o.Status == OrderStatusApproved
}

// CanChangeStatus verifica si se puede cambiar al estado especificado
// Deprecated: usar OrderStateMachine en su lugar
func (o *Order) CanChangeStatus(newStatus OrderStatus) bool {
	if o.Status == OrderStatusCancelled {
		return false
	}

	// Mantener compatibilidad con código antiguo
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusQuote:         {OrderStatusApproved, OrderStatusCancelled},
		OrderStatusApproved:      {OrderStatusManufacturing, OrderStatusInProduction, OrderStatusCancelled},
		OrderStatusManufacturing: {OrderStatusFinished, OrderStatusCancelled},
		OrderStatusInProduction:  {OrderStatusFinished, OrderStatusCancelled},
		OrderStatusFinished:      {OrderStatusDelivered, OrderStatusCancelled},
		OrderStatusDelivered:     {},
		OrderStatusCancelled:     {},
		// Nuevos estados
		OrderStatusPlanned:   {OrderStatusManufacturing, OrderStatusCancelled},
		OrderStatusPending:   {OrderStatusConfirmed, OrderStatusCancelled},
		OrderStatusConfirmed: {OrderStatusDelivered, OrderStatusCancelled},
	}

	allowedStatuses, exists := validTransitions[o.Status]
	if !exists {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return true
		}
	}

	return false
}

// CalculateTotal calcula el total de la orden basado en los items
func (o *Order) CalculateTotal() float64 {
	total := 0.0
	for _, item := range o.Items {
		total += item.Subtotal
	}
	return total - o.Discount
}

// NeedsManufacturing determina si la orden requiere fabricación
// Retorna true si algún item necesita ser fabricado
func (o *Order) NeedsManufacturing() bool {
	for _, item := range o.Items {
		if item.NeedsManufacturing() {
			return true
		}
	}
	return false
}

// HasFullStockCoverage determina si todos los items están cubiertos por stock reservado
func (o *Order) HasFullStockCoverage() bool {
	return !o.NeedsManufacturing()
}

// IsInternalCustomer determina si la orden es para un cliente interno
// Los clientes internos requieren registro contable al completar la venta
func (o *Order) IsInternalCustomer() bool {
	return o.CustomerID != nil && *o.CustomerID > 0
}
