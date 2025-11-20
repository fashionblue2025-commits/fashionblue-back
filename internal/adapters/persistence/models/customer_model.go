package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// CustomerModel representa el modelo de persistencia para clientes
type CustomerModel struct {
	ID               uint   `gorm:"primaryKey"`
	Name             string `gorm:"not null"`
	Phone            string `gorm:"not null"`
	Address          string `gorm:"type:text"`
	RiskLevel        string `gorm:"type:varchar(20);default:'LOW'"` // LOW, MEDIUM, HIGH
	ShirtSizeID      *uint  `gorm:"index"`
	PantsSizeID      *uint  `gorm:"index"`
	ShoesSizeID      *uint  `gorm:"index"`
	Birthday         *time.Time
	Notes            string    `gorm:"type:text"`
	IsActive         bool      `gorm:"default:true"`
	PaymentFrequency string    `gorm:"type:varchar(20);default:'NONE'"` // NONE, WEEKLY, BIWEEKLY, MONTHLY
	PaymentDays      string    `gorm:"type:varchar(50)"`                // Días de pago separados por coma (ej: "2,17")
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`

	// Relaciones
	ShirtSize *SizeModel `gorm:"foreignKey:ShirtSizeID"`
	PantsSize *SizeModel `gorm:"foreignKey:PantsSizeID"`
	ShoesSize *SizeModel `gorm:"foreignKey:ShoesSizeID"`
}

// TableName especifica el nombre de la tabla
func (CustomerModel) TableName() string {
	return "customers"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *CustomerModel) ToEntity() *entities.Customer {
	if m == nil {
		return nil
	}

	customer := &entities.Customer{
		ID:               m.ID,
		Name:             m.Name,
		Phone:            m.Phone,
		Address:          m.Address,
		RiskLevel:        entities.RiskLevel(m.RiskLevel),
		ShirtSizeID:      m.ShirtSizeID,
		PantsSizeID:      m.PantsSizeID,
		ShoesSizeID:      m.ShoesSizeID,
		Birthday:         m.Birthday,
		Notes:            m.Notes,
		IsActive:         m.IsActive,
		PaymentFrequency: entities.PaymentFrequency(m.PaymentFrequency),
		PaymentDays:      m.PaymentDays,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}

	// Convertir tallas si están cargadas
	if m.ShirtSize != nil {
		customer.ShirtSize = m.ShirtSize.ToEntity()
	}
	if m.PantsSize != nil {
		customer.PantsSize = m.PantsSize.ToEntity()
	}
	if m.ShoesSize != nil {
		customer.ShoesSize = m.ShoesSize.ToEntity()
	}

	return customer
}

// FromEntity convierte una entidad de dominio a modelo
func (m *CustomerModel) FromEntity(customer *entities.Customer) {
	if customer == nil {
		return
	}

	m.ID = customer.ID
	m.Name = customer.Name
	m.Phone = customer.Phone
	m.Address = customer.Address
	m.RiskLevel = string(customer.RiskLevel)
	m.ShirtSizeID = customer.ShirtSizeID
	m.PantsSizeID = customer.PantsSizeID
	m.ShoesSizeID = customer.ShoesSizeID
	m.Birthday = customer.Birthday
	m.Notes = customer.Notes
	m.IsActive = customer.IsActive
	m.PaymentFrequency = string(customer.PaymentFrequency)
	m.PaymentDays = customer.PaymentDays
	m.CreatedAt = customer.CreatedAt
	m.UpdatedAt = customer.UpdatedAt
}

// CustomerTransactionModel representa el modelo de persistencia para transacciones
type CustomerTransactionModel struct {
	ID              uint                `gorm:"primaryKey"`
	CustomerID      uint                `gorm:"not null"`
	Customer        *CustomerModel      `gorm:"foreignKey:CustomerID"`
	Type            string              `gorm:"type:varchar(10);not null"` // DEUDA o ABONO
	Amount          float64             `gorm:"not null"`                  // Siempre positivo
	Description     string              `gorm:"type:text"`                 // Descripción del movimiento
	PaymentMethodID *uint               `gorm:"index"`                     // ID del método de pago (solo para ABONO)
	PaymentMethod   *PaymentMethodModel `gorm:"foreignKey:PaymentMethodID"`
	Date            time.Time           `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TableName especifica el nombre de la tabla
func (CustomerTransactionModel) TableName() string {
	return "customer_transactions"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *CustomerTransactionModel) ToEntity() *entities.CustomerTransaction {
	transaction := &entities.CustomerTransaction{
		ID:              m.ID,
		CustomerID:      m.CustomerID,
		Type:            entities.TransactionType(m.Type),
		Amount:          m.Amount,
		Description:     m.Description,
		PaymentMethodID: m.PaymentMethodID,
		Date:            m.Date,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}

	// Convertir método de pago si está cargado
	if m.PaymentMethod != nil {
		transaction.PaymentMethod = m.PaymentMethod.ToEntity()
	}

	return transaction
}

// FromEntity convierte una entidad de dominio a modelo
func (m *CustomerTransactionModel) FromEntity(transaction *entities.CustomerTransaction) {
	m.ID = transaction.ID
	m.CustomerID = transaction.CustomerID
	m.Type = string(transaction.Type)
	m.Amount = transaction.Amount
	m.Description = transaction.Description
	m.PaymentMethodID = transaction.PaymentMethodID
	m.Date = transaction.Date
	m.CreatedAt = transaction.CreatedAt
	m.UpdatedAt = transaction.UpdatedAt
}
