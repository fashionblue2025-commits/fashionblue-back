package payment_method

import (
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
)

type ListPaymentMethodsUseCase struct {
	repo ports.PaymentMethodRepository
}

func NewListPaymentMethodsUseCase(repo ports.PaymentMethodRepository) *ListPaymentMethodsUseCase {
	return &ListPaymentMethodsUseCase{repo: repo}
}

func (uc *ListPaymentMethodsUseCase) Execute(activeOnly bool) ([]*entities.PaymentMethodOption, error) {
	return uc.repo.List(activeOnly)
}
