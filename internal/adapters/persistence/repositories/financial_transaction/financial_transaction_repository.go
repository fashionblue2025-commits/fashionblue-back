package financial_transaction

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type financialTransactionRepository struct {
	db *gorm.DB
}

// NewFinancialTransactionRepository crea una nueva instancia del repositorio
func NewFinancialTransactionRepository(db *gorm.DB) ports.FinancialTransactionRepository {
	return &financialTransactionRepository{db: db}
}

func (r *financialTransactionRepository) Create(ctx context.Context, transaction *entities.FinancialTransaction) error {
	model := &models.FinancialTransactionModel{}
	model.FromEntity(transaction)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	transaction.ID = model.ID
	transaction.CreatedAt = model.CreatedAt
	transaction.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *financialTransactionRepository) GetByID(ctx context.Context, id uint) (*entities.FinancialTransaction, error) {
	var model models.FinancialTransactionModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *financialTransactionRepository) List(ctx context.Context, filters map[string]interface{}) ([]entities.FinancialTransaction, error) {
	var models []models.FinancialTransactionModel
	query := r.db.WithContext(ctx).Order("date DESC, created_at DESC")

	// Aplicar filtros
	if transactionType, ok := filters["type"].(string); ok && transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	transactions := make([]entities.FinancialTransaction, len(models))
	for i, model := range models {
		transactions[i] = *model.ToEntity()
	}
	return transactions, nil
}

func (r *financialTransactionRepository) GetTotalIncome(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&models.FinancialTransactionModel{}).
		Where("type = ?", string(entities.FinancialTransactionTypeIncome)).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *financialTransactionRepository) GetTotalExpenses(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&models.FinancialTransactionModel{}).
		Where("type = ?", string(entities.FinancialTransactionTypeExpense)).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *financialTransactionRepository) GetBalance(ctx context.Context) (float64, error) {
	var result struct {
		Income   float64
		Expenses float64
	}

	// Subconsulta para ingresos
	err := r.db.WithContext(ctx).
		Model(&models.FinancialTransactionModel{}).
		Select(`
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) as expenses
		`, string(entities.FinancialTransactionTypeIncome), string(entities.FinancialTransactionTypeExpense)).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.Income - result.Expenses, nil
}

func (r *financialTransactionRepository) GetTotalByCategory(ctx context.Context, transactionType entities.FinancialTransactionType, category entities.FinancialTransactionCategory) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&models.FinancialTransactionModel{}).
		Where("type = ? AND category = ?", string(transactionType), string(category)).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *financialTransactionRepository) GetBalanceByDateRange(ctx context.Context, startDate, endDate string) (float64, error) {
	var result struct {
		Income   float64
		Expenses float64
	}

	query := r.db.WithContext(ctx).Model(&models.FinancialTransactionModel{})

	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	err := query.
		Select(`
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) as expenses
		`, string(entities.FinancialTransactionTypeIncome), string(entities.FinancialTransactionTypeExpense)).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.Income - result.Expenses, nil
}
