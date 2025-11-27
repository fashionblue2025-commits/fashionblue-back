package category

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ports.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *entities.Category) error {
	model := &models.CategoryModel{}
	model.FromEntity(category)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*category = *model.ToEntity()
	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id uint) (*entities.Category, error) {
	var model models.CategoryModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (*entities.Category, error) {
	var model models.CategoryModel
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&model).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *categoryRepository) List(ctx context.Context, filters map[string]interface{}) ([]entities.Category, error) {
	var modelList []models.CategoryModel
	query := r.db.WithContext(ctx)

	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	// Filtro por IDs específicos (para permisos de usuario)
	if ids, ok := filters["ids"].([]uint); ok && len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}

	if err := query.Find(&modelList).Error; err != nil {
		return nil, err
	}

	categories := make([]entities.Category, len(modelList))
	for i, model := range modelList {
		categories[i] = *model.ToEntity()
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *entities.Category) error {
	model := &models.CategoryModel{}
	model.FromEntity(category)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	*category = *model.ToEntity()
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.CategoryModel{}, id).Error
}
