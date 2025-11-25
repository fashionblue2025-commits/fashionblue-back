package customer

import (
	"net/http"
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases"
	"github.com/labstack/echo/v4"
)

type StatementHandler struct {
	generateStatementUC *usecases.GenerateCustomerStatementUseCase
}

func NewStatementHandler(generateStatementUC *usecases.GenerateCustomerStatementUseCase) *StatementHandler {
	return &StatementHandler{
		generateStatementUC: generateStatementUC,
	}
}

// DownloadStatement godoc
// @Summary Descargar estado de cuenta en PDF
// @Description Genera y descarga un PDF con el estado de cuenta del cliente. Se puede especificar un período en días o descargar todas las transacciones.
// @Tags customers
// @Accept json
// @Produce application/pdf
// @Param id path int true "ID del cliente"
// @Param days query int false "Número de días del período (omitir para todas las transacciones)"
// @Success 200 {file} binary "PDF del estado de cuenta"
// @Failure 400 {object} map[string]interface{} "Request inválido"
// @Failure 404 {object} map[string]interface{} "Cliente no encontrado"
// @Failure 500 {object} map[string]interface{} "Error del servidor"
// @Security BearerAuth
// @Router /customers/{id}/statement [get]
func (h *StatementHandler) DownloadStatement(c echo.Context) error {
	// Obtener ID del cliente
	customerIDStr := c.Param("id")
	customerID, err := strconv.ParseUint(customerIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID de cliente inválido"})
	}

	// Obtener parámetro de días (opcional)
	var days *int
	if daysStr := c.QueryParam("days"); daysStr != "" {
		daysInt, err := strconv.Atoi(daysStr)
		if err != nil || daysInt <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Parámetro 'days' debe ser un número positivo"})
		}
		days = &daysInt
	}

	// Crear request
	req := usecases.StatementRequest{
		CustomerID: uint(customerID),
		Days:       days, // nil = todas las transacciones
	}

	// Generar PDF
	response, err := h.generateStatementUC.Execute(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	// Configurar headers para descarga del PDF
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+response.Filename)
	c.Response().Header().Set("Content-Length", strconv.Itoa(len(response.PDFBytes)))

	// Enviar el PDF
	return c.Blob(http.StatusOK, "application/pdf", response.PDFBytes)
}
