package customer

import (
	"context"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type CreatePaymentUseCase struct {
	transactionRepo ports.CustomerTransactionRepository
	customerRepo    ports.CustomerRepository
}

func NewCreatePaymentUseCase(
	transactionRepo ports.CustomerTransactionRepository,
	customerRepo ports.CustomerRepository,
) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{
		transactionRepo: transactionRepo,
		customerRepo:    customerRepo,
	}
}

type CreatePaymentInput struct {
	CustomerID      uint
	Amount          float64
	PaymentMethodID uint
	Concept         string
	Date            time.Time
}

func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*entities.CustomerTransaction, error) {
	// Verificar que el cliente existe
	_, err := uc.customerRepo.GetByID(ctx, input.CustomerID)
	if err != nil {
		return nil, err
	}

	// Crear la transacción de pago (tipo ABONO)
	transaction := &entities.CustomerTransaction{
		CustomerID:      input.CustomerID,
		Type:            entities.TransactionTypePayment, // ABONO
		Amount:          input.Amount,                    // Siempre positivo
		Description:     input.Concept,
		PaymentMethodID: &input.PaymentMethodID,
		Date:            input.Date,
	}

	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}
