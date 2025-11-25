package ports

import (
	"context"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// AuditLogRepository define la interfaz para el repositorio de logs de auditoría
type AuditLogRepository interface {
	// Create guarda un nuevo log de auditoría
	Create(ctx context.Context, log *entities.AuditLog) error

	// GetByID obtiene un log por su ID
	GetByID(ctx context.Context, id uint) (*entities.AuditLog, error)

	// List obtiene logs con filtros
	List(ctx context.Context, filters entities.AuditLogFilters) ([]entities.AuditLog, error)

	// GetByOrderID obtiene todos los logs de una orden
	GetByOrderID(ctx context.Context, orderID uint) ([]entities.AuditLog, error)

	// GetByUserID obtiene todos los logs de un usuario
	GetByUserID(ctx context.Context, userID uint) ([]entities.AuditLog, error)

	// Count obtiene el total de logs con filtros
	Count(ctx context.Context, filters entities.AuditLogFilters) (int64, error)

	// DeleteOlderThan elimina logs más antiguos que la fecha especificada
	DeleteOlderThan(ctx context.Context, date string) error
}
