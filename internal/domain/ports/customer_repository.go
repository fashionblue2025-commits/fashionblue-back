package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// CustomerRepository define las operaciones para clientes
type CustomerRepository interface {
	Create(ctx context.Context, customer *entities.Customer) error
	GetByID(ctx context.Context, id uint) (*entities.Customer, error)
	GetByEmail(ctx context.Context, email string) (*entities.Customer, error)
	GetByDocument(ctx context.Context, documentNum string) (*entities.Customer, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.Customer, error)
	Update(ctx context.Context, customer *entities.Customer) error
	Delete(ctx context.Context, id uint) error
	GetUpcomingPayments(ctx context.Context, daysRange int) ([]entities.Customer, error)
	GetBalance(ctx context.Context, customerID uint) (float64, error)
}

// CustomerTransactionRepository define las operaciones para transacciones de clientes
type CustomerTransactionRepository interface {
	Create(ctx context.Context, transaction *entities.CustomerTransaction) error
	GetByID(ctx context.Context, id uint) (*entities.CustomerTransaction, error)
	ListByCustomer(ctx context.Context, customerID uint) ([]entities.CustomerTransaction, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.CustomerTransaction, error)
}
