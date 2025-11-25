package product

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *entities.Product) error {
	model := &models.ProductModel{}
	model.FromEntity(product)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*product = *model.ToEntity()
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*entities.Product, error) {
	var model models.ProductModel
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Preload("Variants.Size").
		First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

// GetByAttributes busca un producto por nombre, color y talla
// Esto permite identificar productos únicos en el inventario
func (r *productRepository) GetByAttributes(ctx context.Context, name, color string, sizeID *uint) (*entities.Product, error) {
	var model models.ProductModel
	query := r.db.WithContext(ctx).Preload("Category").Preload("Size").Where("name = ? AND color = ?", name, color)

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

func (r *productRepository) List(ctx context.Context, filters map[string]interface{}) ([]entities.Product, error) {
	var modelList []models.ProductModel
	query := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Preload("Variants.Size")

	if categoryID, ok := filters["category_id"].(uint); ok && categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if err := query.Find(&modelList).Error; err != nil {
		return nil, err
	}

	products := make([]entities.Product, len(modelList))
	for i, model := range modelList {
		products[i] = *model.ToEntity()
	}

	return products, nil
}

func (r *productRepository) ListByCategory(ctx context.Context, categoryID uint) ([]entities.Product, error) {
	var modelList []models.ProductModel
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Preload("Variants.Size").
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	products := make([]entities.Product, len(modelList))
	for i, model := range modelList {
		products[i] = *model.ToEntity()
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *entities.Product) error {
	model := &models.ProductModel{}
	model.FromEntity(product)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	*product = *model.ToEntity()
	return nil
}

func (r *productRepository) UpdateStock(ctx context.Context, productID uint, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&models.ProductModel{}).
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock + ?", quantity)).
		Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ProductModel{}, id).Error
}

func (r *productRepository) GetLowStockProducts(ctx context.Context) ([]entities.Product, error) {
	var modelList []models.ProductModel
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Preload("Variants.Size").
		Where("stock <= min_stock AND is_active = ?", true).
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	products := make([]entities.Product, len(modelList))
	for i, model := range modelList {
		products[i] = *model.ToEntity()
	}

	return products, nil
}
