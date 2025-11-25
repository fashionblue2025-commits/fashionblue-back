package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"gorm.io/gorm"
)

// OrderItemModel representa el modelo de persistencia de un item de orden
type OrderItemModel struct {
	ID               uint    `gorm:"primaryKey"`
	OrderID          uint    `gorm:"not null;index"`
	ProductVariantID *uint   `gorm:"index;default:null"` // Nullable: se asigna cuando se crea la variante
	ProductName      string  `gorm:"not null"`           // Snapshot del nombre del producto base
	CategoryID       uint    `gorm:"not null;index"`     // Snapshot de la categoría del producto
	Color            string  `gorm:"not null"`           // Snapshot del color solicitado
	SizeID           *uint   `gorm:"index;default:null"` // Snapshot de la talla solicitada
	Quantity         int     `gorm:"not null"`           // Cantidad total solicitada
	ReservedQuantity int     `gorm:"not null;default:0"` // Cantidad reservada del stock existente
	UnitPrice        float64 `gorm:"not null"`
	Subtotal         float64 `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	// Relaciones
	Order          *OrderModel          `gorm:"foreignKey:OrderID"`
	ProductVariant *ProductVariantModel `gorm:"foreignKey:ProductVariantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Size           *SizeModel           `gorm:"foreignKey:SizeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName especifica el nombre de la tabla
func (OrderItemModel) TableName() string {
	return "order_items"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *OrderItemModel) ToEntity() *entities.OrderItem {
	item := &entities.OrderItem{
		ID:               m.ID,
		OrderID:          m.OrderID,
		ProductVariantID: 0, // Default a 0 si es nil
		ProductName:      m.ProductName,
		CategoryID:       m.CategoryID,
		Color:            m.Color,
		SizeID:           m.SizeID,
		Quantity:         m.Quantity,
		ReservedQuantity: m.ReservedQuantity,
		UnitPrice:        m.UnitPrice,
		Subtotal:         m.Subtotal,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}

	// Asignar ProductVariantID si existe
	if m.ProductVariantID != nil {
		item.ProductVariantID = *m.ProductVariantID
	}

	// Convertir ProductVariant si existe
	if m.ProductVariant != nil {
		item.ProductVariant = m.ProductVariant.ToEntity()
	}

	// Convertir Size si existe
	if m.Size != nil {
		item.Size = m.Size.ToEntity()
	}

	return item
}

// FromEntity convierte la entidad de dominio a modelo
func (m *OrderItemModel) FromEntity(item *entities.OrderItem) {
	m.ID = item.ID
	// No establecer OrderID aquí - GORM lo asignará automáticamente desde la relación
	// m.OrderID = item.OrderID

	// ProductVariantID puede ser 0 para variantes nuevas (se asigna después)
	if item.ProductVariantID != 0 {
		m.ProductVariantID = &item.ProductVariantID
	} else {
		m.ProductVariantID = nil
	}

	m.ProductName = item.ProductName
	m.CategoryID = item.CategoryID
	m.Color = item.Color
	m.SizeID = item.SizeID
	m.Quantity = item.Quantity
	m.ReservedQuantity = item.ReservedQuantity
	m.UnitPrice = item.UnitPrice
	m.Subtotal = item.Subtotal
}

// BeforeSave es un hook de GORM que se ejecuta antes de guardar
// Asegura que ProductVariantID sea NULL en lugar de 0
func (m *OrderItemModel) BeforeSave(tx *gorm.DB) error {
	// Si ProductVariantID es un puntero a 0, convertirlo a nil
	if m.ProductVariantID != nil && *m.ProductVariantID == 0 {
		m.ProductVariantID = nil
	}
	return nil
}
