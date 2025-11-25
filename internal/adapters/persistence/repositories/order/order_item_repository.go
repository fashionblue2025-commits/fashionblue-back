package order

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type orderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) ports.OrderItemRepository {
	return &orderItemRepository{db: db}
}

func (r *orderItemRepository) Create(ctx context.Context, item *entities.OrderItem) error {
	model := &models.OrderItemModel{}
	model.FromEntity(item)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*item = *model.ToEntity()
	return nil
}

func (r *orderItemRepository) GetByID(ctx context.Context, id uint) (*entities.OrderItem, error) {
	var model models.OrderItemModel
	err := r.db.WithContext(ctx).
		Preload("ProductVariant").
		Preload("ProductVariant.Size").
		Preload("Size").
		First(&model, id).Error

	if err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *orderItemRepository) GetByOrderID(ctx context.Context, orderID uint) ([]entities.OrderItem, error) {
	var models []models.OrderItemModel
	err := r.db.WithContext(ctx).
		Preload("ProductVariant").
		Preload("ProductVariant.Size").
		Preload("Size").
		Where("order_id = ?", orderID).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	items := make([]entities.OrderItem, len(models))
	for i, model := range models {
		items[i] = *model.ToEntity()
	}

	return items, nil
}

func (r *orderItemRepository) Update(ctx context.Context, item *entities.OrderItem) error {
	model := &models.OrderItemModel{}
	model.FromEntity(item)

	return r.db.WithContext(ctx).Save(model).Error
}

func (r *orderItemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.OrderItemModel{}, id).Error
}
