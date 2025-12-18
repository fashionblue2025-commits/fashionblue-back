package dto

import "time"

// AccountStatementDraftDTO representa el borrador de una cuenta de cobro que puede ser editado
type AccountStatementDraftDTO struct {
	OrderID         uint      `json:"orderId"`
	OrderNumber     string    `json:"orderNumber"`
	StatementNumber string    `json:"statementNumber"`
	SellerName      string    `json:"sellerName"`
	SellerID        string    `json:"sellerId"`
	ClientName      string    `json:"clientName"` // Editable
	City            string    `json:"city"`       // Editable
	Date            time.Time `json:"date"`       // Editable
	Concept         string    `json:"concept"`    // Editable
	TotalAmount     float64   `json:"totalAmount"`
	BankAccount     string    `json:"bankAccount"`
}

// AccountStatementConfirmRequest representa la solicitud para confirmar y generar el PDF
type AccountStatementConfirmRequest struct {
	ClientName string    `json:"clientName" validate:"required"`
	City       string    `json:"city" validate:"required"`
	Date       time.Time `json:"date" validate:"required"`
	Concept    string    `json:"concept" validate:"required"`
}
