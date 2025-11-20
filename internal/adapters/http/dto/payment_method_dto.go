package dto

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// PaymentMethodDTO representa la respuesta de un método de pago
type PaymentMethodDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToPaymentMethodDTO convierte una entidad PaymentMethodOption a DTO
func ToPaymentMethodDTO(paymentMethod *entities.PaymentMethodOption) PaymentMethodDTO {
	return PaymentMethodDTO{
		ID:        paymentMethod.ID,
		Name:      paymentMethod.Name,
		IsActive:  paymentMethod.IsActive,
		CreatedAt: paymentMethod.CreatedAt,
		UpdatedAt: paymentMethod.UpdatedAt,
	}
}

// ToPaymentMethodDTOList convierte una lista de métodos de pago a DTOs
func ToPaymentMethodDTOList(paymentMethods []entities.PaymentMethodOption) []PaymentMethodDTO {
	dtos := make([]PaymentMethodDTO, len(paymentMethods))
	for i, pm := range paymentMethods {
		dtos[i] = ToPaymentMethodDTO(&pm)
	}
	return dtos
}

// ToPaymentMethodDTOListFromPointers convierte una lista de punteros a DTOs
func ToPaymentMethodDTOListFromPointers(paymentMethods []*entities.PaymentMethodOption) []PaymentMethodDTO {
	dtos := make([]PaymentMethodDTO, len(paymentMethods))
	for i, pm := range paymentMethods {
		dtos[i] = ToPaymentMethodDTO(pm)
	}
	return dtos
}
