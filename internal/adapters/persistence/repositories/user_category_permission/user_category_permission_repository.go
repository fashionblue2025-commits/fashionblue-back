package usercategorypermission

import (
	"context"
	"errors"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/persistence/models"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type userCategoryPermissionRepository struct {
	db *gorm.DB
}

// NewUserCategoryPermissionRepository crea una nueva instancia del repositorio
func NewUserCategoryPermissionRepository(db *gorm.DB) ports.UserCategoryPermissionRepository {
	return &userCategoryPermissionRepository{db: db}
}

func (r *userCategoryPermissionRepository) Create(ctx context.Context, permission *entities.UserCategoryPermission) error {
	model := &models.UserCategoryPermissionModel{}
	model.FromEntity(permission)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	permission.ID = model.ID
	permission.CreatedAt = model.CreatedAt
	permission.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *userCategoryPermissionRepository) Update(ctx context.Context, permission *entities.UserCategoryPermission) error {
	model := &models.UserCategoryPermissionModel{}
	model.FromEntity(permission)

	return r.db.WithContext(ctx).Save(model).Error
}

func (r *userCategoryPermissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.UserCategoryPermissionModel{}, id).Error
}

func (r *userCategoryPermissionRepository) GetByID(ctx context.Context, id uint) (*entities.UserCategoryPermission, error) {
	var model models.UserCategoryPermissionModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}

func (r *userCategoryPermissionRepository) GetByUserAndCategory(ctx context.Context, userID, categoryID uint) (*entities.UserCategoryPermission, error) {
	var model models.UserCategoryPermissionModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

func (r *userCategoryPermissionRepository) ListByUser(ctx context.Context, userID uint) ([]entities.UserCategoryPermission, error) {
	var models []models.UserCategoryPermissionModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Category").
		Find(&models).Error; err != nil {
		return nil, err
	}

	permissions := make([]entities.UserCategoryPermission, len(models))
	for i, model := range models {
		permissions[i] = *model.ToEntity()
	}

	return permissions, nil
}

func (r *userCategoryPermissionRepository) ListByCategory(ctx context.Context, categoryID uint) ([]entities.UserCategoryPermission, error) {
	var models []models.UserCategoryPermissionModel
	if err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Preload("User").
		Find(&models).Error; err != nil {
		return nil, err
	}

	permissions := make([]entities.UserCategoryPermission, len(models))
	for i, model := range models {
		permissions[i] = *model.ToEntity()
	}

	return permissions, nil
}

func (r *userCategoryPermissionRepository) GetAllowedCategoriesForUser(ctx context.Context, userID uint, action string) ([]uint, error) {
	var categoryIDs []uint

	query := r.db.WithContext(ctx).
		Model(&models.UserCategoryPermissionModel{}).
		Where("user_id = ?", userID).
		Select("category_id")

	switch action {
	case "view":
		query = query.Where("can_view = ?", true)
	case "create":
		query = query.Where("can_create = ?", true)
	case "edit":
		query = query.Where("can_edit = ?", true)
	case "delete":
		query = query.Where("can_delete = ?", true)
	default:
		return nil, errors.New("invalid action")
	}

	if err := query.Pluck("category_id", &categoryIDs).Error; err != nil {
		return nil, err
	}

	return categoryIDs, nil
}

func (r *userCategoryPermissionRepository) HasPermission(ctx context.Context, userID, categoryID uint, action string) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).
		Model(&models.UserCategoryPermissionModel{}).
		Where("user_id = ? AND category_id = ?", userID, categoryID)

	switch action {
	case "view":
		query = query.Where("can_view = ?", true)
	case "create":
		query = query.Where("can_create = ?", true)
	case "edit":
		query = query.Where("can_edit = ?", true)
	case "delete":
		query = query.Where("can_delete = ?", true)
	default:
		return false, errors.New("invalid action")
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userCategoryPermissionRepository) SetPermissions(ctx context.Context, userID uint, permissions []entities.UserCategoryPermission) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Eliminar permisos existentes
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserCategoryPermissionModel{}).Error; err != nil {
			return err
		}

		// Crear nuevos permisos
		for _, perm := range permissions {
			model := &models.UserCategoryPermissionModel{}
			perm.UserID = userID // Asegurar que el UserID es correcto
			model.FromEntity(&perm)

			if err := tx.Create(model).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *userCategoryPermissionRepository) DeleteByUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.UserCategoryPermissionModel{}).Error
}

func (r *userCategoryPermissionRepository) DeleteByCategory(ctx context.Context, categoryID uint) error {
	return r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Delete(&models.UserCategoryPermissionModel{}).Error
}
