package dto

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// CustomerDTO representa la respuesta de un cliente
type CustomerDTO struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Phone            string    `json:"phone"`
	Address          string    `json:"address,omitempty"`
	RiskLevel        string    `json:"riskLevel"`
	ShirtSizeID      *uint     `json:"shirtSizeId,omitempty"`
	PantsSizeID      *uint     `json:"pantsSizeId,omitempty"`
	ShoesSizeID      *uint     `json:"shoesSizeId,omitempty"`
	PaymentFrequency string    `json:"paymentFrequency,omitempty"`
	PaymentDays      string    `json:"paymentDays,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// ToCustomerDTO convierte una entidad Customer a DTO
func ToCustomerDTO(customer *entities.Customer) CustomerDTO {
	return CustomerDTO{
		ID:               customer.ID,
		Name:             customer.Name,
		Phone:            customer.Phone,
		Address:          customer.Address,
		RiskLevel:        string(customer.RiskLevel),
		ShirtSizeID:      customer.ShirtSizeID,
		PantsSizeID:      customer.PantsSizeID,
		ShoesSizeID:      customer.ShoesSizeID,
		PaymentFrequency: string(customer.PaymentFrequency),
		PaymentDays:      customer.PaymentDays,
		CreatedAt:        customer.CreatedAt,
		UpdatedAt:        customer.UpdatedAt,
	}
}

// ToCustomerDTOList convierte una lista de entidades Customer a DTOs
func ToCustomerDTOList(customers []entities.Customer) []CustomerDTO {
	dtos := make([]CustomerDTO, len(customers))
	for i, customer := range customers {
		dtos[i] = ToCustomerDTO(&customer)
	}
	return dtos
}

// CustomerTransactionDTO representa la respuesta de una transacción
type CustomerTransactionDTO struct {
	ID              uint              `json:"id"`
	CustomerID      uint              `json:"customerId"`
	Type            string            `json:"type"`
	Amount          float64           `json:"amount"`
	Description     string            `json:"description"`
	PaymentMethodID *uint             `json:"paymentMethodId,omitempty"`
	PaymentMethod   *PaymentMethodDTO `json:"paymentMethod,omitempty"`
	Date            time.Time         `json:"date"`
	CreatedAt       time.Time         `json:"createdAt"`
}

// ToCustomerTransactionDTO convierte una entidad CustomerTransaction a DTO
func ToCustomerTransactionDTO(transaction *entities.CustomerTransaction) CustomerTransactionDTO {
	dto := CustomerTransactionDTO{
		ID:              transaction.ID,
		CustomerID:      transaction.CustomerID,
		Type:            string(transaction.Type),
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		PaymentMethodID: transaction.PaymentMethodID,
		Date:            transaction.Date,
		CreatedAt:       transaction.CreatedAt,
	}

	if transaction.PaymentMethod != nil {
		paymentMethodDTO := ToPaymentMethodDTO(transaction.PaymentMethod)
		dto.PaymentMethod = &paymentMethodDTO
	}

	return dto
}

// ToCustomerTransactionDTOList convierte una lista de transacciones a DTOs
func ToCustomerTransactionDTOList(transactions []*entities.CustomerTransaction) []CustomerTransactionDTO {
	dtos := make([]CustomerTransactionDTO, len(transactions))
	for i, transaction := range transactions {
		dtos[i] = ToCustomerTransactionDTO(transaction)
	}
	return dtos
}

// CustomerBalanceDTO representa el balance de un cliente
type CustomerBalanceDTO struct {
	CustomerID uint    `json:"customerId"`
	Balance    float64 `json:"balance"`
}

// CustomerWithBalanceDTO representa un cliente con su balance
type CustomerWithBalanceDTO struct {
	Customer CustomerDTO `json:"customer"`
	Balance  float64     `json:"balance"`
}

// ToCustomerWithBalanceDTO convierte un customer y balance a DTO
func ToCustomerWithBalanceDTO(customer *entities.Customer, balance float64) CustomerWithBalanceDTO {
	return CustomerWithBalanceDTO{
		Customer: ToCustomerDTO(customer),
		Balance:  balance,
	}
}

// ToCustomerTransactionDTOListFromSlice convierte []entities.CustomerTransaction a DTOs
func ToCustomerTransactionDTOListFromSlice(transactions []entities.CustomerTransaction) []CustomerTransactionDTO {
	dtos := make([]CustomerTransactionDTO, len(transactions))
	for i, transaction := range transactions {
		dtos[i] = ToCustomerTransactionDTO(&transaction)
	}
	return dtos
}
