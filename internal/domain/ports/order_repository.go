package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id uint) (*entities.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*entities.Order, error)
	List(ctx context.Context, filters map[string]interface{}) ([]entities.Order, error)
	Update(ctx context.Context, order *entities.Order) error
	UpdateStatus(ctx context.Context, id uint, status entities.OrderStatus) error
	Delete(ctx context.Context, id uint) error
}
