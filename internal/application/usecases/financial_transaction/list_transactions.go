package financial_transaction

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListTransactionsUseCase struct {
	transactionRepo ports.FinancialTransactionRepository
}

func NewListTransactionsUseCase(transactionRepo ports.FinancialTransactionRepository) *ListTransactionsUseCase {
	return &ListTransactionsUseCase{
		transactionRepo: transactionRepo,
	}
}

func (uc *ListTransactionsUseCase) Execute(ctx context.Context, filters map[string]interface{}) ([]entities.FinancialTransaction, error) {
	return uc.transactionRepo.List(ctx, filters)
}
