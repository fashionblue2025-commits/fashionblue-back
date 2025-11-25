package entities

import "time"

// FinancialTransactionType representa el tipo de transacción financiera
type FinancialTransactionType string

const (
	FinancialTransactionTypeIncome  FinancialTransactionType = "INCOME"  // Ingreso/Inyección de capital
	FinancialTransactionTypeExpense FinancialTransactionType = "EXPENSE" // Gasto
)

// FinancialTransactionCategory representa la categoría de la transacción
type FinancialTransactionCategory string

// Categorías de INCOME (Ingresos)
const (
	FinancialTransactionCategoryInvestment FinancialTransactionCategory = "INVESTMENT" // Inversión personal
	FinancialTransactionCategoryLoan       FinancialTransactionCategory = "LOAN"       // Préstamo
	FinancialTransactionCategoryProfit     FinancialTransactionCategory = "PROFIT"     // Reinversión de utilidades
	FinancialTransactionCategorySales      FinancialTransactionCategory = "SALES"      // Ventas (futuro)
)

// Categorías de EXPENSE (Gastos)
const (
	FinancialTransactionCategoryOperational FinancialTransactionCategory = "OPERATIONAL" // Gastos operacionales
	FinancialTransactionCategoryPersonnel   FinancialTransactionCategory = "PERSONNEL"   // Nómina y personal
	FinancialTransactionCategoryInventory   FinancialTransactionCategory = "INVENTORY"   // Compra de inventario/materias primas
	FinancialTransactionCategoryMarketing   FinancialTransactionCategory = "MARKETING"   // Marketing y publicidad
	FinancialTransactionCategoryUtilities   FinancialTransactionCategory = "UTILITIES"   // Servicios públicos
	FinancialTransactionCategoryRent        FinancialTransactionCategory = "RENT"        // Arriendo
	FinancialTransactionCategoryOther       FinancialTransactionCategory = "OTHER"       // Otros
)

// FinancialTransaction representa una transacción financiera (ingreso o gasto)
type FinancialTransaction struct {
	ID          uint
	Type        FinancialTransactionType     // INCOME o EXPENSE
	Category    FinancialTransactionCategory // Categoría específica
	Amount      float64                      // Monto (siempre positivo)
	Description string                       // Descripción detallada
	Date        time.Time                    // Fecha de la transacción
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate valida los datos de la transacción
func (ft *FinancialTransaction) Validate() error {
	if ft.Amount <= 0 {
		return ErrInvalidInput
	}
	if ft.Description == "" {
		return ErrInvalidInput
	}
	if ft.Type == "" {
		return ErrInvalidInput
	}
	if ft.Category == "" {
		return ErrInvalidInput
	}
	return nil
}

// IsIncome retorna true si es una transacción de ingreso
func (ft *FinancialTransaction) IsIncome() bool {
	return ft.Type == FinancialTransactionTypeIncome
}

// IsExpense retorna true si es una transacción de gasto
func (ft *FinancialTransaction) IsExpense() bool {
	return ft.Type == FinancialTransactionTypeExpense
}

// GetSignedAmount retorna el monto con signo (positivo para income, negativo para expense)
func (ft *FinancialTransaction) GetSignedAmount() float64 {
	if ft.IsIncome() {
		return ft.Amount
	}
	return -ft.Amount
}
