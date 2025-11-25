package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"gorm.io/gorm"
)

// OrderModel representa el modelo de persistencia de una orden
type OrderModel struct {
	ID                    uint      `gorm:"primaryKey"`
	OrderNumber           string    `gorm:"uniqueIndex;not null"`
	CustomerID            *uint     `gorm:"index;default:null"` // ID del cliente interno (opcional)
	CustomerName          string    `gorm:"not null"`
	SellerID              uint      `gorm:"not null;index"`
	Type                  string    `gorm:"not null;type:varchar(20)"`               // Deprecated: usar OrderType
	OrderType             string    `gorm:"type:varchar(20);default:'CUSTOM';index"` // CUSTOM, INVENTORY, SALE
	Status                string    `gorm:"not null;type:varchar(20);index"`
	TotalAmount           float64   `gorm:"not null;default:0"`
	Discount              float64   `gorm:"not null;default:0"`
	Notes                 string    `gorm:"type:text"`
	OrderDate             time.Time `gorm:"not null;index"`
	EstimatedDeliveryDate *time.Time
	ActualDeliveryDate    *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"`

	// Relaciones
	Seller *UserModel        `gorm:"foreignKey:SellerID"`
	Items  []OrderItemModel  `gorm:"foreignKey:OrderID"`
	Photos []OrderPhotoModel `gorm:"foreignKey:OrderID"`
}

// TableName especifica el nombre de la tabla
func (OrderModel) TableName() string {
	return "orders"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *OrderModel) ToEntity() *entities.Order {
	// Usar OrderType si está disponible, sino usar Type (compatibilidad)
	orderType := m.OrderType
	if orderType == "" {
		orderType = m.Type
	}

	order := &entities.Order{
		ID:                    m.ID,
		OrderNumber:           m.OrderNumber,
		CustomerID:            m.CustomerID,
		CustomerName:          m.CustomerName,
		SellerID:              m.SellerID,
		Type:                  entities.OrderType(orderType),
		Status:                entities.OrderStatus(m.Status),
		TotalAmount:           m.TotalAmount,
		Discount:              m.Discount,
		Notes:                 m.Notes,
		OrderDate:             m.OrderDate,
		EstimatedDeliveryDate: m.EstimatedDeliveryDate,
		ActualDeliveryDate:    m.ActualDeliveryDate,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}

	// Convertir seller si existe
	if m.Seller != nil {
		order.Seller = m.Seller.ToEntity()
	}

	// Convertir items
	if len(m.Items) > 0 {
		order.Items = make([]entities.OrderItem, len(m.Items))
		for i, item := range m.Items {
			order.Items[i] = *item.ToEntity()
		}
	}

	// Convertir fotos
	if len(m.Photos) > 0 {
		order.Photos = make([]entities.OrderPhoto, len(m.Photos))
		for i, photo := range m.Photos {
			order.Photos[i] = *photo.ToEntity()
		}
	}

	return order
}

// FromEntity convierte la entidad de dominio a modelo
func (m *OrderModel) FromEntity(order *entities.Order) {
	m.ID = order.ID
	m.OrderNumber = order.OrderNumber
	m.CustomerID = order.CustomerID
	m.CustomerName = order.CustomerName
	m.SellerID = order.SellerID
	m.Type = string(order.Type)      // Mantener compatibilidad
	m.OrderType = string(order.Type) // Usar nuevo campo
	m.Status = string(order.Status)
	m.TotalAmount = order.TotalAmount
	m.Discount = order.Discount
	m.Notes = order.Notes
	m.OrderDate = order.OrderDate
	m.EstimatedDeliveryDate = order.EstimatedDeliveryDate
	m.ActualDeliveryDate = order.ActualDeliveryDate

	// Convertir items
	if len(order.Items) > 0 {
		m.Items = make([]OrderItemModel, len(order.Items))
		for i, item := range order.Items {
			m.Items[i].FromEntity(&item)
		}
	}
}
