package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// FinancialTransactionModel representa el modelo de persistencia para transacciones financieras
type FinancialTransactionModel struct {
	ID          uint      `gorm:"primaryKey"`
	Type        string    `gorm:"type:varchar(20);not null;index"` // INCOME o EXPENSE
	Category    string    `gorm:"type:varchar(50);not null;index"`
	Amount      float64   `gorm:"not null"`
	Description string    `gorm:"type:text;not null"`
	Date        time.Time `gorm:"not null;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName especifica el nombre de la tabla
func (FinancialTransactionModel) TableName() string {
	return "financial_transactions"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *FinancialTransactionModel) ToEntity() *entities.FinancialTransaction {
	return &entities.FinancialTransaction{
		ID:          m.ID,
		Type:        entities.FinancialTransactionType(m.Type),
		Category:    entities.FinancialTransactionCategory(m.Category),
		Amount:      m.Amount,
		Description: m.Description,
		Date:        m.Date,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo
func (m *FinancialTransactionModel) FromEntity(transaction *entities.FinancialTransaction) {
	m.ID = transaction.ID
	m.Type = string(transaction.Type)
	m.Category = string(transaction.Category)
	m.Amount = transaction.Amount
	m.Description = transaction.Description
	m.Date = transaction.Date
	m.CreatedAt = transaction.CreatedAt
	m.UpdatedAt = transaction.UpdatedAt
}
