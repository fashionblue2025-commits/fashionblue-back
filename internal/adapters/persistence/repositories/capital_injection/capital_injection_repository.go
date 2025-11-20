package capital_injection

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type capitalInjectionRepository struct {
	db *gorm.DB
}

// NewCapitalInjectionRepository crea una nueva instancia del repositorio
func NewCapitalInjectionRepository(db *gorm.DB) ports.CapitalInjectionRepository {
	return &capitalInjectionRepository{db: db}
}

func (r *capitalInjectionRepository) Create(ctx context.Context, injection *entities.CapitalInjection) error {
	model := &models.CapitalInjectionModel{}
	model.FromEntity(injection)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	injection.ID = model.ID
	injection.CreatedAt = model.CreatedAt
	injection.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *capitalInjectionRepository) GetByID(ctx context.Context, id uint) (*entities.CapitalInjection, error) {
	var model models.CapitalInjectionModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *capitalInjectionRepository) List(ctx context.Context, filters map[string]interface{}) ([]entities.CapitalInjection, error) {
	var models []models.CapitalInjectionModel
	query := r.db.WithContext(ctx).Order("date DESC")

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	injections := make([]entities.CapitalInjection, len(models))
	for i, model := range models {
		injections[i] = *model.ToEntity()
	}

	return injections, nil
}

func (r *capitalInjectionRepository) GetTotal(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&models.CapitalInjectionModel{}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	return total, err
}
