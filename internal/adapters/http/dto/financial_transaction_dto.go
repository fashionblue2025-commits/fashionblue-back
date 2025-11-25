package dto

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

type FinancialTransactionDTO struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"`        // INCOME o EXPENSE
	Category    string    `json:"category"`    // Categoría específica
	Amount      float64   `json:"amount"`      // Siempre positivo
	Description string    `json:"description"` // Descripción detallada
	Date        time.Time `json:"date"`        // Fecha de la transacción
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func ToFinancialTransactionDTO(transaction *entities.FinancialTransaction) *FinancialTransactionDTO {
	return &FinancialTransactionDTO{
		ID:          transaction.ID,
		Type:        string(transaction.Type),
		Category:    string(transaction.Category),
		Amount:      transaction.Amount,
		Description: transaction.Description,
		Date:        transaction.Date,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}
}

func ToFinancialTransactionDTOList(transactions []entities.FinancialTransaction) []*FinancialTransactionDTO {
	var dtos []*FinancialTransactionDTO
	for _, transaction := range transactions {
		dtos = append(dtos, ToFinancialTransactionDTO(&transaction))
	}
	return dtos
}
