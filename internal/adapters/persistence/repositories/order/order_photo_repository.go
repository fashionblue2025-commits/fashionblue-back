package order

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type orderPhotoRepository struct {
	db *gorm.DB
}

func NewOrderPhotoRepository(db *gorm.DB) ports.OrderPhotoRepository {
	return &orderPhotoRepository{db: db}
}

func (r *orderPhotoRepository) Create(ctx context.Context, photo *entities.OrderPhoto) error {
	model := &models.OrderPhotoModel{}
	model.FromEntity(photo)
	
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	
	*photo = *model.ToEntity()
	return nil
}

func (r *orderPhotoRepository) GetByID(ctx context.Context, id uint) (*entities.OrderPhoto, error) {
	var model models.OrderPhotoModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	
	if err != nil {
		return nil, err
	}
	
	return model.ToEntity(), nil
}

func (r *orderPhotoRepository) GetByOrderID(ctx context.Context, orderID uint) ([]entities.OrderPhoto, error) {
	var models []models.OrderPhotoModel
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Find(&models).Error
	
	if err != nil {
		return nil, err
	}
	
	photos := make([]entities.OrderPhoto, len(models))
	for i, model := range models {
		photos[i] = *model.ToEntity()
	}
	
	return photos, nil
}

func (r *orderPhotoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.OrderPhotoModel{}, id).Error
}
