package entities

import (
	"errors"
	"time"
)

// OrderPhoto representa una foto adjunta a una orden
type OrderPhoto struct {
	ID          uint
	OrderID     uint
	PhotoURL    string
	Description string
	UploadedAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate valida los datos de la foto
func (op *OrderPhoto) Validate() error {
	if op.OrderID == 0 {
		return errors.New("order id is required")
	}
	if op.PhotoURL == "" {
		return errors.New("photo url is required")
	}
	return nil
}
