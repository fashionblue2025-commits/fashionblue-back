package financial_transaction

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreateTransactionUseCase struct {
	transactionRepo ports.FinancialTransactionRepository
}

func NewCreateTransactionUseCase(transactionRepo ports.FinancialTransactionRepository) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		transactionRepo: transactionRepo,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, transaction *entities.FinancialTransaction) error {
	// Validar transacción
	if err := transaction.Validate(); err != nil {
		return err
	}

	// Crear transacción
	return uc.transactionRepo.Create(ctx, transaction)
}
