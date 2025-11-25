package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

type OrderItemRepository interface {
	Create(ctx context.Context, item *entities.OrderItem) error
	GetByID(ctx context.Context, id uint) (*entities.OrderItem, error)
	GetByOrderID(ctx context.Context, orderID uint) ([]entities.OrderItem, error)
	Update(ctx context.Context, item *entities.OrderItem) error
	Delete(ctx context.Context, id uint) error
}
