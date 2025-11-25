package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// FinancialTransactionRepository define las operaciones para gestionar transacciones financieras
type FinancialTransactionRepository interface {
	Create(ctx context.Context, transaction *entities.FinancialTransaction) error
	GetByID(ctx context.Context, id uint) (*entities.FinancialTransaction, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.FinancialTransaction, error)

	// Métodos para calcular balances
	GetTotalIncome(ctx context.Context) (float64, error)
	GetTotalExpenses(ctx context.Context) (float64, error)
	GetBalance(ctx context.Context) (float64, error)

	// Métodos por categoría
	GetTotalByCategory(ctx context.Context, transactionType entities.FinancialTransactionType, category entities.FinancialTransactionCategory) (float64, error)

	// Métodos por rango de fechas
	GetBalanceByDateRange(ctx context.Context, startDate, endDate string) (float64, error)
}
