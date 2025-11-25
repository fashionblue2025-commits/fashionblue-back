package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type productPhotoRepository struct {
	db *gorm.DB
}

// NewProductPhotoRepository crea una nueva instancia del repositorio
func NewProductPhotoRepository(db *gorm.DB) ports.ProductPhotoRepository {
	return &productPhotoRepository{db: db}
}

func (r *productPhotoRepository) Create(ctx context.Context, photo *entities.ProductPhoto) error {
	var model models.ProductPhotoModel
	model.FromEntity(photo)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	// Actualizar el ID de la entidad con el ID generado
	photo.ID = model.ID
	return nil
}

func (r *productPhotoRepository) GetByID(ctx context.Context, id uint) (*entities.ProductPhoto, error) {
	var model models.ProductPhotoModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *productPhotoRepository) GetByProductID(ctx context.Context, productID uint) ([]entities.ProductPhoto, error) {
	var photoModels []models.ProductPhotoModel
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("display_order ASC, created_at ASC").
		Find(&photoModels).Error; err != nil {
		return nil, err
	}

	// Convertir modelos a entidades
	photos := make([]entities.ProductPhoto, len(photoModels))
	for i, model := range photoModels {
		photos[i] = *model.ToEntity()
	}
	return photos, nil
}

func (r *productPhotoRepository) Update(ctx context.Context, photo *entities.ProductPhoto) error {
	var model models.ProductPhotoModel
	model.FromEntity(photo)
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *productPhotoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ProductPhotoModel{}, id).Error
}

func (r *productPhotoRepository) SetAsPrimary(ctx context.Context, photoID uint, productID uint) error {
	// Iniciar transacción
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Quitar is_primary de todas las fotos del producto
		if err := tx.Model(&models.ProductPhotoModel{}).
			Where("product_id = ?", productID).
			Update("is_primary", false).Error; err != nil {
			return err
		}

		// Establecer la foto seleccionada como primary
		if err := tx.Model(&models.ProductPhotoModel{}).
			Where("id = ?", photoID).
			Update("is_primary", true).Error; err != nil {
			return err
		}

		return nil
	})
}
