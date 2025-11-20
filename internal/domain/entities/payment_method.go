package entities

import (
	"errors"
	"time"
)

// PaymentMethodOption representa una opción de método de pago disponible
// (ej: NEQUI Sonia, NEQUI Jhon, Daviplata, Efectivo)
type PaymentMethodOption struct {
	ID        uint
	Name      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate valida los datos del método de pago
func (pmo *PaymentMethodOption) Validate() error {
	if pmo.Name == "" {
		return errors.New("payment method name is required")
	}
	return nil
}
