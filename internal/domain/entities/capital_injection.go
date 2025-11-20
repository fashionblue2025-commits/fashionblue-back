package entities

import "time"

// CapitalInjection representa una inyección de capital al negocio
type CapitalInjection struct {
	ID          uint
	Amount      float64
	Description string
	Source      string // De dónde viene el capital (ej: "Inversión personal", "Préstamo bancario")
	Date        time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate valida los datos de la inyección de capital
func (ci *CapitalInjection) Validate() error {
	if ci.Amount <= 0 {
		return ErrInvalidInput
	}
	if ci.Description == "" {
		return ErrInvalidInput
	}
	return nil
}
