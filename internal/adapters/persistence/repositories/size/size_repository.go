package size

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type sizeRepository struct {
	db *gorm.DB
}

// NewSizeRepository crea una nueva instancia del repositorio de tallas
func NewSizeRepository(db *gorm.DB) ports.SizeRepository {
	return &sizeRepository{db: db}
}

func (r *sizeRepository) Create(ctx context.Context, size *entities.Size) error {
	model := &models.SizeModel{}
	model.FromEntity(size)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*size = *model.ToEntity()
	return nil
}

func (r *sizeRepository) GetByID(ctx context.Context, id uint) (*entities.Size, error) {
	var model models.SizeModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *sizeRepository) List(ctx context.Context, filters map[string]interface{}) ([]*entities.Size, error) {
	var models []models.SizeModel
	query := r.db.WithContext(ctx).Where("is_active = ?", true)

	// Aplicar filtros
	if sizeType, ok := filters["type"].(string); ok && sizeType != "" {
		query = query.Where("type = ?", sizeType)
	}

	// Ordenar por tipo y orden
	query = query.Order("type ASC, \"order\" ASC")

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	sizes := make([]*entities.Size, len(models))
	for i, model := range models {
		sizes[i] = model.ToEntity()
	}

	return sizes, nil
}

func (r *sizeRepository) Update(ctx context.Context, size *entities.Size) error {
	model := &models.SizeModel{}
	model.FromEntity(size)

	return r.db.WithContext(ctx).Save(model).Error
}

func (r *sizeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.SizeModel{}, id).Error
}

func (r *sizeRepository) GetByType(ctx context.Context, sizeType entities.SizeType) ([]*entities.Size, error) {
	var models []models.SizeModel

	if err := r.db.WithContext(ctx).
		Where("type = ? AND is_active = ?", sizeType, true).
		Order("\"order\" ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	sizes := make([]*entities.Size, len(models))
	for i, model := range models {
		sizes[i] = model.ToEntity()
	}

	return sizes, nil
}
