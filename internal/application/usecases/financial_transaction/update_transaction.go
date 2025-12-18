package financial_transaction

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type UpdateTransactionUseCase struct {
	transactionRepo ports.FinancialTransactionRepository
}

func NewUpdateTransactionUseCase(transactionRepo ports.FinancialTransactionRepository) *UpdateTransactionUseCase {
	return &UpdateTransactionUseCase{
		transactionRepo: transactionRepo,
	}
}

func (uc *UpdateTransactionUseCase) Execute(ctx context.Context, transaction *entities.FinancialTransaction) error {
	// Validar que la transacción existe
	existing, err := uc.transactionRepo.GetByID(ctx, transaction.ID)
	if err != nil {
		return err
	}

	if existing == nil {
		return entities.ErrNotFound
	}

	// Validar transacción
	if err := transaction.Validate(); err != nil {
		return err
	}

	// Actualizar transacción
	return uc.transactionRepo.Update(ctx, transaction)
}
