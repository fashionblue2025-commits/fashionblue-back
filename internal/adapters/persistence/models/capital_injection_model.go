package models

import (
	"time"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
)

// CapitalInjectionModel representa el modelo de persistencia para inyecciones de capital
type CapitalInjectionModel struct {
	ID          uint      `gorm:"primaryKey"`
	Amount      float64   `gorm:"not null"`
	Description string    `gorm:"type:text;not null"`
	Source      string    `gorm:"type:varchar(255)"`
	Date        time.Time `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName especifica el nombre de la tabla
func (CapitalInjectionModel) TableName() string {
	return "capital_injections"
}

// ToEntity convierte el modelo a entidad de dominio
func (m *CapitalInjectionModel) ToEntity() *entities.CapitalInjection {
	return &entities.CapitalInjection{
		ID:          m.ID,
		Amount:      m.Amount,
		Description: m.Description,
		Source:      m.Source,
		Date:        m.Date,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromEntity convierte una entidad de dominio a modelo
func (m *CapitalInjectionModel) FromEntity(injection *entities.CapitalInjection) {
	m.ID = injection.ID
	m.Amount = injection.Amount
	m.Description = injection.Description
	m.Source = injection.Source
	m.Date = injection.Date
	m.CreatedAt = injection.CreatedAt
	m.UpdatedAt = injection.UpdatedAt
}
