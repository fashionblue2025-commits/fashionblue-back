package customer

import (
	"net/http"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/dto"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/customer"
	"github.com/labstack/echo/v4"
)

// AddTransaction maneja la solicitud para agregar movimientos manuales
func (h *CustomerHandler) AddTransaction(c echo.Context) error {
	var req customer.AddTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
	}

	// Validar que hay al menos una transacción
	if len(req.Transactions) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "At least one transaction is required",
		})
	}

	transactions, err := h.addTransactionUC.Execute(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to add transactions",
			"error":   err.Error(),
		})
	}

	// Convertir a DTOs
	transactionDTOs := dto.ToCustomerTransactionDTOList(transactions)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Transactions added successfully",
		"data":    transactionDTOs,
	})
}
