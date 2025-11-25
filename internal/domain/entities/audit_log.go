package entities

import "time"

// AuditLog representa un registro de auditoría
type AuditLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	EventType   string    `json:"eventType" gorm:"type:varchar(100);not null;index"`
	OrderID     uint      `json:"orderId" gorm:"not null;index"`
	OrderNumber string    `json:"orderNumber" gorm:"type:varchar(50)"`
	UserID      *uint     `json:"userId" gorm:"index"`
	UserName    string    `json:"userName" gorm:"type:varchar(100)"`
	OldStatus   string    `json:"oldStatus" gorm:"type:varchar(50)"`
	NewStatus   string    `json:"newStatus" gorm:"type:varchar(50)"`
	Description string    `json:"description" gorm:"type:text"`
	Metadata    string    `json:"metadata" gorm:"type:jsonb"` // JSON con datos adicionales
	IPAddress   string    `json:"ipAddress" gorm:"type:varchar(45)"`
	UserAgent   string    `json:"userAgent" gorm:"type:varchar(255)"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// TableName especifica el nombre de la tabla
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogFilters define filtros para consultar logs de auditoría
type AuditLogFilters struct {
	EventType string
	OrderID   *uint
	UserID    *uint
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}
