package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

type OrderPhotoRepository interface {
	Create(ctx context.Context, photo *entities.OrderPhoto) error
	GetByID(ctx context.Context, id uint) (*entities.OrderPhoto, error)
	GetByOrderID(ctx context.Context, orderID uint) ([]entities.OrderPhoto, error)
	Delete(ctx context.Context, id uint) error
}
