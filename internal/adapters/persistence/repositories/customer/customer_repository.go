package customer

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) ports.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *entities.Customer) error {
	model := &models.CustomerModel{}
	model.FromEntity(customer)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*customer = *model.ToEntity()
	return nil
}

func (r *customerRepository) GetByID(ctx context.Context, id uint) (*entities.Customer, error) {
	var model models.CustomerModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *customerRepository) GetByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	var model models.CustomerModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *customerRepository) GetByDocument(ctx context.Context, documentNum string) (*entities.Customer, error) {
	var model models.CustomerModel
	err := r.db.WithContext(ctx).Where("document_num = ?", documentNum).First(&model).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *customerRepository) List(ctx context.Context, filters map[string]interface{}) ([]entities.Customer, error) {
	var modelList []models.CustomerModel
	query := r.db.WithContext(ctx)

	if customerType, ok := filters["type"].(string); ok && customerType != "" {
		query = query.Where("type = ?", customerType)
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ?", "%"+name+"%", "%"+name+"%")
	}

	if err := query.Find(&modelList).Error; err != nil {
		return nil, err
	}

	customers := make([]entities.Customer, len(modelList))
	for i, model := range modelList {
		customers[i] = *model.ToEntity()
	}

	return customers, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *entities.Customer) error {
	model := &models.CustomerModel{}
	model.FromEntity(customer)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	*customer = *model.ToEntity()
	return nil
}

func (r *customerRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.CustomerModel{}, id).Error
}

// GetUpcomingPayments obtiene clientes con pagos próximos
func (r *customerRepository) GetUpcomingPayments(ctx context.Context, daysRange int) ([]entities.Customer, error) {
	var modelList []models.CustomerModel

	// Obtener todos los clientes activos con frecuencia de pago configurada
	query := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("payment_frequency != ?", "NONE").
		Where("payment_days != ?", "")

	if err := query.Find(&modelList).Error; err != nil {
		return nil, err
	}

	// Filtrar en memoria los que tienen pagos próximos
	customers := make([]entities.Customer, 0)
	for _, model := range modelList {
		customer := model.ToEntity()
		if customer.IsPaymentDue(daysRange) {
			customers = append(customers, *customer)
		}
	}

	return customers, nil
}

// GetBalance calcula el balance actual de un cliente
// Balance = Σ(DEUDA) - Σ(ABONO)
func (r *customerRepository) GetBalance(ctx context.Context, customerID uint) (float64, error) {
	var balance float64

	err := r.db.WithContext(ctx).
		Model(&models.CustomerTransactionModel{}).
		Where("customer_id = ?", customerID).
		Select(`COALESCE(
			SUM(CASE WHEN type = 'DEUDA' THEN amount ELSE 0 END) - 
			SUM(CASE WHEN type = 'ABONO' THEN amount ELSE 0 END), 
			0
		)`).
		Scan(&balance).Error

	return balance, err
}

// CustomerTransactionRepository
type customerTransactionRepository struct {
	db *gorm.DB
}

func NewCustomerTransactionRepository(db *gorm.DB) ports.CustomerTransactionRepository {
	return &customerTransactionRepository{db: db}
}

func (r *customerTransactionRepository) Create(ctx context.Context, transaction *entities.CustomerTransaction) error {
	model := &models.CustomerTransactionModel{}
	model.FromEntity(transaction)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*transaction = *model.ToEntity()
	return nil
}

func (r *customerTransactionRepository) GetByID(ctx context.Context, id uint) (*entities.CustomerTransaction, error) {
	var model models.CustomerTransactionModel
	err := r.db.WithContext(ctx).Preload("Customer").First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *customerTransactionRepository) ListByCustomer(ctx context.Context, customerID uint) ([]entities.CustomerTransaction, error) {
	var modelList []models.CustomerTransactionModel
	err := r.db.WithContext(ctx).
		Where("customer_id = ?", customerID).
		Order("date DESC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	transactions := make([]entities.CustomerTransaction, len(modelList))
	for i, model := range modelList {
		transactions[i] = *model.ToEntity()
	}

	return transactions, nil
}

func (r *customerTransactionRepository) List(ctx context.Context, filters map[string]interface{}) ([]entities.CustomerTransaction, error) {
	var modelList []models.CustomerTransactionModel
	query := r.db.WithContext(ctx).Preload("Customer")

	if customerID, ok := filters["customer_id"].(uint); ok && customerID > 0 {
		query = query.Where("customer_id = ?", customerID)
	}

	if transactionType, ok := filters["type"].(string); ok && transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	if err := query.Order("date DESC").Find(&modelList).Error; err != nil {
		return nil, err
	}

	transactions := make([]entities.CustomerTransaction, len(modelList))
	for i, model := range modelList {
		transactions[i] = *model.ToEntity()
	}

	return transactions, nil
}
