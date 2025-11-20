package entities

import "time"

// Supplier representa un proveedor
type Supplier struct {
	ID          uint
	Name        string
	ContactName string
	Phone       string
	Email       string
	Address     string
	Notes       string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate valida los datos del proveedor
func (s *Supplier) Validate() error {
	if s.Name == "" {
		return ErrInvalidInput
	}
	return nil
}
