package payment_method

import (
	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type paymentMethodRepository struct {
	db *gorm.DB
}

// NewPaymentMethodRepository crea una nueva instancia del repositorio
func NewPaymentMethodRepository(db *gorm.DB) ports.PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

// Create crea un nuevo método de pago
func (r *paymentMethodRepository) Create(paymentMethod *entities.PaymentMethodOption) error {
	model := &models.PaymentMethodModel{}
	model.FromEntity(paymentMethod)

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	paymentMethod.ID = model.ID
	paymentMethod.CreatedAt = model.CreatedAt
	paymentMethod.UpdatedAt = model.UpdatedAt

	return nil
}

// GetByID obtiene un método de pago por su ID
func (r *paymentMethodRepository) GetByID(id uint) (*entities.PaymentMethodOption, error) {
	var model models.PaymentMethodModel

	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	return model.ToEntity(), nil
}

// List obtiene todos los métodos de pago
func (r *paymentMethodRepository) List(activeOnly bool) ([]*entities.PaymentMethodOption, error) {
	var models []models.PaymentMethodModel
	query := r.db

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	paymentMethods := make([]*entities.PaymentMethodOption, len(models))
	for i, model := range models {
		paymentMethods[i] = model.ToEntity()
	}

	return paymentMethods, nil
}

// Update actualiza un método de pago
func (r *paymentMethodRepository) Update(paymentMethod *entities.PaymentMethodOption) error {
	model := &models.PaymentMethodModel{}
	model.FromEntity(paymentMethod)

	return r.db.Save(model).Error
}

// Delete elimina un método de pago (soft delete)
func (r *paymentMethodRepository) Delete(id uint) error {
	return r.db.Model(&models.PaymentMethodModel{}).Where("id = ?", id).Update("is_active", false).Error
}
