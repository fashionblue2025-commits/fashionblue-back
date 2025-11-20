package supplier

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type supplierRepository struct {
	db *gorm.DB
}

// NewSupplierRepository crea una nueva instancia del repositorio
func NewSupplierRepository(db *gorm.DB) ports.SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) Create(ctx context.Context, supplier *entities.Supplier) error {
	model := &models.SupplierModel{}
	model.FromEntity(supplier)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	supplier.ID = model.ID
	supplier.CreatedAt = model.CreatedAt
	supplier.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *supplierRepository) GetByID(ctx context.Context, id uint) (*entities.Supplier, error) {
	var model models.SupplierModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *supplierRepository) List(ctx context.Context, filters map[string]interface{}) ([]entities.Supplier, error) {
	var models []models.SupplierModel
	query := r.db.WithContext(ctx)

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	suppliers := make([]entities.Supplier, len(models))
	for i, model := range models {
		suppliers[i] = *model.ToEntity()
	}

	return suppliers, nil
}

func (r *supplierRepository) Update(ctx context.Context, supplier *entities.Supplier) error {
	model := &models.SupplierModel{}
	model.FromEntity(supplier)

	return r.db.WithContext(ctx).Save(model).Error
}

func (r *supplierRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.SupplierModel{}, id).Error
}
