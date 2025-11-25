package order

import (
	"context"
	"fmt"
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *entities.Order) error {
	model := &models.OrderModel{}
	model.FromEntity(order)
	if model.OrderNumber == "" {
		model.OrderNumber = r.generateOrderNumber()
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*order = *model.ToEntity()
	return nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uint) (*entities.Order, error) {
	var model models.OrderModel
	err := r.db.WithContext(ctx).
		Preload("Seller").
		Preload("Items").
		Preload("Items.ProductVariant").
		Preload("Items.ProductVariant.Size").
		Preload("Items.Size").
		Preload("Photos").
		First(&model, id).Error

	if err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*entities.Order, error) {
	var model models.OrderModel
	err := r.db.WithContext(ctx).
		Preload("Seller").
		Preload("Items").
		Preload("Items.ProductVariant").
		Preload("Items.ProductVariant.Size").
		Preload("Items.Size").
		Preload("Photos").
		Where("order_number = ?", orderNumber).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *orderRepository) List(ctx context.Context, filters map[string]interface{}) ([]entities.Order, error) {
	var models []models.OrderModel
	query := r.db.WithContext(ctx).
		Preload("Seller").
		Preload("Items").
		Preload("Items.ProductVariant").
		Preload("Items.ProductVariant.Size").
		Preload("Items.Size")

	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if sellerID, ok := filters["seller_id"].(uint); ok && sellerID > 0 {
		query = query.Where("seller_id = ?", sellerID)
	}
	if orderType, ok := filters["type"].(string); ok && orderType != "" {
		query = query.Where("type = ?", orderType)
	}
	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("order_date >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("order_date <= ?", endDate)
	}

	query = query.Order("order_date DESC")

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	orders := make([]entities.Order, len(models))
	for i, model := range models {
		orders[i] = *model.ToEntity()
	}

	return orders, nil
}

func (r *orderRepository) Update(ctx context.Context, order *entities.Order) error {
	model := &models.OrderModel{}
	model.FromEntity(order)

	// Usar transacción para actualizar orden e items
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Actualizar la orden
		if err := tx.Save(model).Error; err != nil {
			return err
		}

		// Actualizar los items de la orden
		for i := range order.Items {
			itemModel := &models.OrderItemModel{}
			itemModel.FromEntity(&order.Items[i])
			itemModel.OrderID = order.ID // Asegurar que tenga el OrderID correcto

			if err := tx.Save(itemModel).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uint, status entities.OrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.OrderModel{}).
		Where("id = ?", id).
		Update("status", string(status)).Error
}

func (r *orderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.OrderModel{}, id).Error
}

func (r *orderRepository) generateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%d", now.Format("20060102"), now.Unix()%10000)
}
