package financial_transaction

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type GetTransactionUseCase struct {
	transactionRepo ports.FinancialTransactionRepository
}

func NewGetTransactionUseCase(transactionRepo ports.FinancialTransactionRepository) *GetTransactionUseCase {
	return &GetTransactionUseCase{
		transactionRepo: transactionRepo,
	}
}

func (uc *GetTransactionUseCase) Execute(ctx context.Context, id uint) (*entities.FinancialTransaction, error) {
	return uc.transactionRepo.GetByID(ctx, id)
}
