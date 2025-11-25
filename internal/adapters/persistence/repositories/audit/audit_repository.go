package audit

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/ports"
	"gorm.io/gorm"
)

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository crea una nueva instancia del repositorio
func NewAuditLogRepository(db *gorm.DB) ports.AuditLogRepository {
	return &auditLogRepository{db: db}
}

// Create guarda un nuevo log de auditoría
func (r *auditLogRepository) Create(ctx context.Context, log *entities.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetByID obtiene un log por su ID
func (r *auditLogRepository) GetByID(ctx context.Context, id uint) (*entities.AuditLog, error) {
	var log entities.AuditLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// List obtiene logs con filtros
func (r *auditLogRepository) List(ctx context.Context, filters entities.AuditLogFilters) ([]entities.AuditLog, error) {
	var logs []entities.AuditLog

	query := r.db.WithContext(ctx).Model(&entities.AuditLog{})

	// Aplicar filtros
	if filters.EventType != "" {
		query = query.Where("event_type = ?", filters.EventType)
	}

	if filters.OrderID != nil {
		query = query.Where("order_id = ?", *filters.OrderID)
	}

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", *filters.StartDate)
	}

	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", *filters.EndDate)
	}

	// Ordenar por fecha descendente (más recientes primero)
	query = query.Order("created_at DESC")

	// Aplicar paginación
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&logs).Error
	return logs, err
}

// GetByOrderID obtiene todos los logs de una orden
func (r *auditLogRepository) GetByOrderID(ctx context.Context, orderID uint) ([]entities.AuditLog, error) {
	var logs []entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetByUserID obtiene todos los logs de un usuario
func (r *auditLogRepository) GetByUserID(ctx context.Context, userID uint) ([]entities.AuditLog, error) {
	var logs []entities.AuditLog
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// Count obtiene el total de logs con filtros
func (r *auditLogRepository) Count(ctx context.Context, filters entities.AuditLogFilters) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.AuditLog{})

	// Aplicar los mismos filtros que en List
	if filters.EventType != "" {
		query = query.Where("event_type = ?", filters.EventType)
	}

	if filters.OrderID != nil {
		query = query.Where("order_id = ?", *filters.OrderID)
	}

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", *filters.StartDate)
	}

	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", *filters.EndDate)
	}

	err := query.Count(&count).Error
	return count, err
}

// DeleteOlderThan elimina logs más antiguos que la fecha especificada
func (r *auditLogRepository) DeleteOlderThan(ctx context.Context, date string) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", date).
		Delete(&entities.AuditLog{}).Error
}
