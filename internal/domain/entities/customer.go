package entities

import "time"

// RiskLevel representa el nivel de riesgo del cliente
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "LOW"    // Bajo riesgo
	RiskLevelMedium RiskLevel = "MEDIUM" // Riesgo medio
	RiskLevelHigh   RiskLevel = "HIGH"   // Alto riesgo
)

// PaymentFrequency representa la frecuencia de pago del cliente
type PaymentFrequency string

const (
	PaymentFrequencyNone     PaymentFrequency = "NONE"     // Sin pagos recurrentes
	PaymentFrequencyWeekly   PaymentFrequency = "WEEKLY"   // Semanal
	PaymentFrequencyBiweekly PaymentFrequency = "BIWEEKLY" // Quincenal
	PaymentFrequencyMonthly  PaymentFrequency = "MONTHLY"  // Mensual
)

// Customer representa un cliente (entidad de dominio pura)
type Customer struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`          // Nombre completo
	Phone       string     `json:"phone"`         // Teléfono
	Address     string     `json:"address"`       // Dirección
	RiskLevel   RiskLevel  `json:"risk_level"`    // Nivel de riesgo (LOW, MEDIUM, HIGH)
	ShirtSizeID *uint      `json:"shirt_size_id"` // ID de talla de camiseta (opcional)
	PantsSizeID *uint      `json:"pants_size_id"` // ID de talla de pantalón (opcional)
	ShoesSizeID *uint      `json:"shoes_size_id"` // ID de talla de tenis (opcional)
	Birthday    *time.Time `json:"birthday"`      // Fecha de cumpleaños (opcional)
	Notes       string     `json:"notes"`         // Notas adicionales
	IsActive    bool       `json:"is_active"`     // Si el cliente está activo

	// Campos para pagos recurrentes
	PaymentFrequency PaymentFrequency `json:"payment_frequency"` // Frecuencia de pago
	PaymentDays      string           `json:"payment_days"`      // Días de pago separados por coma (ej: "2,17" para quincenal)

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relaciones (se cargan cuando se necesitan)
	ShirtSize *Size `json:"shirt_size,omitempty"`
	PantsSize *Size `json:"pants_size,omitempty"`
	ShoesSize *Size `json:"shoes_size,omitempty"`
}

// GetAge calcula la edad del cliente
func (c *Customer) GetAge() int {
	if c.Birthday == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - c.Birthday.Year()
	if now.YearDay() < c.Birthday.YearDay() {
		age--
	}
	return age
}

// IsBirthday verifica si hoy es el cumpleaños del cliente
func (c *Customer) IsBirthday() bool {
	if c.Birthday == nil {
		return false
	}
	now := time.Now()
	return now.Month() == c.Birthday.Month() && now.Day() == c.Birthday.Day()
}

// IsHighRisk verifica si el cliente es de alto riesgo
func (c *Customer) IsHighRisk() bool {
	return c.RiskLevel == RiskLevelHigh
}

// GetPaymentDaysAsInts convierte los días de pago de string a slice de ints
func (c *Customer) GetPaymentDaysAsInts() []int {
	if c.PaymentDays == "" {
		return []int{}
	}

	days := []int{}
	for _, dayStr := range splitString(c.PaymentDays, ",") {
		if day := parseInt(dayStr); day > 0 && day <= 31 {
			days = append(days, day)
		}
	}
	return days
}

// IsPaymentDue verifica si el cliente tiene un pago próximo (dentro de los próximos N días)
func (c *Customer) IsPaymentDue(daysRange int) bool {
	if c.PaymentFrequency == PaymentFrequencyNone || c.PaymentDays == "" {
		return false
	}

	now := time.Now()
	paymentDays := c.GetPaymentDaysAsInts()

	for i := -daysRange; i <= daysRange; i++ {
		checkDate := now.AddDate(0, 0, i)
		for _, day := range paymentDays {
			if checkDate.Day() == day {
				return true
			}
		}
	}

	return false
}

// Helper functions
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	current := ""
	for _, char := range s {
		if string(char) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func parseInt(s string) int {
	result := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		}
	}
	return result
}

// TransactionType representa el tipo de transacción
type TransactionType string

const (
	TransactionTypeDebt    TransactionType = "DEUDA" // Deuda/Cargo
	TransactionTypePayment TransactionType = "ABONO" // Abono/Pago
)

// CustomerTransaction representa una transacción o movimiento de un cliente
type CustomerTransaction struct {
	ID              uint
	CustomerID      uint
	Type            TransactionType      // DEUDA o ABONO
	Amount          float64              // Siempre positivo, el tipo define si suma o resta
	Description     string               // Descripción detallada del movimiento
	PaymentMethodID *uint                // ID del método de pago (solo para ABONO)
	PaymentMethod   *PaymentMethodOption // Relación con método de pago
	Date            time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
