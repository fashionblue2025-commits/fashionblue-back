package user

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	model := &models.UserModel{}
	model.FromEntity(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	// Actualizar la entidad con los datos generados (ID, timestamps)
	*user = *model.ToEntity()
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	var model models.UserModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var model models.UserModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *userRepository) List(ctx context.Context, filters map[string]interface{}) ([]entities.User, error) {
	var modelList []models.UserModel
	query := r.db.WithContext(ctx)

	if role, ok := filters["role"].(string); ok && role != "" {
		query = query.Where("role = ?", role)
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&modelList).Error; err != nil {
		return nil, err
	}

	// Convertir modelos a entidades
	users := make([]entities.User, len(modelList))
	for i, model := range modelList {
		users[i] = *model.ToEntity()
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	model := &models.UserModel{}
	model.FromEntity(user)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	// Actualizar la entidad con los datos actualizados
	*user = *model.ToEntity()
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.UserModel{}, id).Error
}
