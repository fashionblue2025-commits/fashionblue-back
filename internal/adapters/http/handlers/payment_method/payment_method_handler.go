package payment_method

import (
	"net/http"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/dto"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/payment_method"
	"github.com/labstack/echo/v4"
)

type PaymentMethodHandler struct {
	listUC *payment_method.ListPaymentMethodsUseCase
}

func NewPaymentMethodHandler(
	listUC *payment_method.ListPaymentMethodsUseCase,
) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		listUC: listUC,
	}
}

// List lista todos los métodos de pago
func (h *PaymentMethodHandler) List(c echo.Context) error {
	activeOnly := c.QueryParam("active_only") == "true"

	paymentMethods, err := h.listUC.Execute(activeOnly)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Error al obtener métodos de pago",
			"error":   err.Error(),
		})
	}

	// Convertir a DTOs
	paymentMethodDTOs := dto.ToPaymentMethodDTOListFromPointers(paymentMethods)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Métodos de pago obtenidos exitosamente",
		"data":    paymentMethodDTOs,
	})
}
