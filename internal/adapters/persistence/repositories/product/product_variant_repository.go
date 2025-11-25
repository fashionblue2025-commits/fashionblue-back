package product

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type productVariantRepository struct {
	db *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) ports.ProductVariantRepository {
	return &productVariantRepository{db: db}
}

func (r *productVariantRepository) Create(ctx context.Context, variant *entities.ProductVariant) error {
	model := &models.ProductVariantModel{}
	model.FromEntity(variant)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*variant = *model.ToEntity()
	return nil
}

func (r *productVariantRepository) GetByID(ctx context.Context, id uint) (*entities.ProductVariant, error) {
	var model models.ProductVariantModel
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Size").
		First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *productVariantRepository) GetByProductAndAttributes(ctx context.Context, productID uint, color string, sizeID *uint) (*entities.ProductVariant, error) {
	var model models.ProductVariantModel
	query := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Size").
		Where("product_id = ? AND color = ?", productID, color)

	// Manejar sizeID que puede ser NULL
	if sizeID != nil {
		query = query.Where("size_id = ?", *sizeID)
	} else {
		query = query.Where("size_id IS NULL")
	}

	err := query.First(&model).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *productVariantRepository) ListByProduct(ctx context.Context, productID uint) ([]entities.ProductVariant, error) {
	var modelList []models.ProductVariantModel
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Size").
		Where("product_id = ?", productID).
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	variants := make([]entities.ProductVariant, len(modelList))
	for i, model := range modelList {
		variants[i] = *model.ToEntity()
	}

	return variants, nil
}

func (r *productVariantRepository) Update(ctx context.Context, variant *entities.ProductVariant) error {
	model := &models.ProductVariantModel{}
	model.FromEntity(variant)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	*variant = *model.ToEntity()
	return nil
}

func (r *productVariantRepository) UpdateStock(ctx context.Context, variantID uint, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&models.ProductVariantModel{}).
		Where("id = ?", variantID).
		Update("stock", gorm.Expr("stock + ?", quantity)).
		Error
}

func (r *productVariantRepository) ReserveStock(ctx context.Context, variantID uint, quantity int) error {
	// Verificar que hay suficiente stock disponible
	var variant models.ProductVariantModel
	if err := r.db.WithContext(ctx).First(&variant, variantID).Error; err != nil {
		return err
	}

	availableStock := variant.Stock - variant.ReservedStock
	if availableStock < quantity {
		return errors.New("insufficient stock available")
	}

	// Incrementar stock reservado
	return r.db.WithContext(ctx).
		Model(&models.ProductVariantModel{}).
		Where("id = ?", variantID).
		Update("reserved_stock", gorm.Expr("reserved_stock + ?", quantity)).
		Error
}

func (r *productVariantRepository) ReleaseStock(ctx context.Context, variantID uint, quantity int) error {
	// Decrementar stock total y stock reservado
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Decrementar stock total
		if err := tx.Model(&models.ProductVariantModel{}).
			Where("id = ?", variantID).
			Update("stock", gorm.Expr("stock - ?", quantity)).
			Error; err != nil {
			return err
		}

		// Decrementar stock reservado
		if err := tx.Model(&models.ProductVariantModel{}).
			Where("id = ?", variantID).
			Update("reserved_stock", gorm.Expr("reserved_stock - ?", quantity)).
			Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *productVariantRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ProductVariantModel{}, id).Error
}

func (r *productVariantRepository) GetLowStockVariants(ctx context.Context) ([]entities.ProductVariant, error) {
	var modelList []models.ProductVariantModel
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Size").
		Joins("JOIN products ON products.id = product_variants.product_id").
		Where("product_variants.stock <= products.min_stock AND product_variants.is_active = ?", true).
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	variants := make([]entities.ProductVariant, len(modelList))
	for i, model := range modelList {
		variants[i] = *model.ToEntity()
	}

	return variants, nil
}
