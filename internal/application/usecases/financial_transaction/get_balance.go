package financial_transaction

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetBalanceUseCase struct {
	transactionRepo ports.FinancialTransactionRepository
}

func NewGetBalanceUseCase(transactionRepo ports.FinancialTransactionRepository) *GetBalanceUseCase {
	return &GetBalanceUseCase{
		transactionRepo: transactionRepo,
	}
}

// BalanceResult representa el resultado del balance
type BalanceResult struct {
	TotalIncome   float64 `json:"totalIncome"`
	TotalExpenses float64 `json:"totalExpenses"`
	Balance       float64 `json:"balance"`
}

func (uc *GetBalanceUseCase) Execute(ctx context.Context) (*BalanceResult, error) {
	income, err := uc.transactionRepo.GetTotalIncome(ctx)
	if err != nil {
		return nil, err
	}

	expenses, err := uc.transactionRepo.GetTotalExpenses(ctx)
	if err != nil {
		return nil, err
	}

	return &BalanceResult{
		TotalIncome:   income,
		TotalExpenses: expenses,
		Balance:       income - expenses,
	}, nil
}
