package customer

import (
	"context"
	"errors"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrInvalidInput     = errors.New("invalid input: ABONO requires payment_method_id")
)

// AddTransactionUseCase maneja la lógica para agregar movimientos manuales
type AddTransactionUseCase struct {
	transactionRepo ports.CustomerTransactionRepository
	customerRepo    ports.CustomerRepository
}

// NewAddTransactionUseCase crea una nueva instancia del caso de uso
func NewAddTransactionUseCase(
	transactionRepo ports.CustomerTransactionRepository,
	customerRepo ports.CustomerRepository,
) *AddTransactionUseCase {
	return &AddTransactionUseCase{
		transactionRepo: transactionRepo,
		customerRepo:    customerRepo,
	}
}

// TransactionInput representa los datos de entrada para un movimiento
type TransactionInput struct {
	Type            entities.TransactionType `json:"type" binding:"required"`        // DEUDA o ABONO
	Amount          float64                  `json:"amount" binding:"required,gt=0"` // Siempre positivo
	Description     string                   `json:"description" binding:"required"` // Descripción del movimiento
	PaymentMethodID *uint                    `json:"payment_method_id"`              // Requerido solo para ABONO
	Date            *time.Time               `json:"date"`                           // Opcional, default: ahora
}

// AddTransactionRequest representa la solicitud para agregar movimientos
type AddTransactionRequest struct {
	CustomerID   uint               `json:"customer_id" binding:"required"`
	Transactions []TransactionInput `json:"transactions" binding:"required,min=1,dive"`
}

// Execute ejecuta el caso de uso
func (uc *AddTransactionUseCase) Execute(ctx context.Context, req AddTransactionRequest) ([]*entities.CustomerTransaction, error) {
	// Verificar que el cliente existe
	customer, err := uc.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	// Crear las transacciones
	var transactions []*entities.CustomerTransaction
	for _, input := range req.Transactions {
		// Validar que ABONO tenga método de pago
		if input.Type == entities.TransactionTypePayment && input.PaymentMethodID == nil {
			return nil, ErrInvalidInput
		}

		// Usar fecha actual si no se proporciona
		date := time.Now()
		if input.Date != nil {
			date = *input.Date
		}

		transaction := &entities.CustomerTransaction{
			CustomerID:      req.CustomerID,
			Type:            input.Type,
			Amount:          input.Amount,
			Description:     input.Description,
			PaymentMethodID: input.PaymentMethodID,
			Date:            date,
		}

		// Crear la transacción
		err := uc.transactionRepo.Create(ctx, transaction)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
