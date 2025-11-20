package ports

import "github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"

// PaymentMethodRepository define las operaciones de persistencia para métodos de pago
type PaymentMethodRepository interface {
	Create(paymentMethod *entities.PaymentMethodOption) error
	GetByID(id uint) (*entities.PaymentMethodOption, error)
	List(activeOnly bool) ([]*entities.PaymentMethodOption, error)
	Update(paymentMethod *entities.PaymentMethodOption) error
	Delete(id uint) error
}
