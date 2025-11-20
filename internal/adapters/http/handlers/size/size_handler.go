package size

import (
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/size"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type SizeHandler struct {
	listSizesUC      *size.ListSizesUseCase
	getSizeUC        *size.GetSizeUseCase
	getSizesByTypeUC *size.GetSizesByTypeUseCase
}

func NewSizeHandler(
	listSizesUC *size.ListSizesUseCase,
	getSizeUC *size.GetSizeUseCase,
	getSizesByTypeUC *size.GetSizesByTypeUseCase,
) *SizeHandler {
	return &SizeHandler{
		listSizesUC:      listSizesUC,
		getSizeUC:        getSizeUC,
		getSizesByTypeUC: getSizesByTypeUC,
	}
}

// List lista todas las tallas con filtros opcionales
func (h *SizeHandler) List(c echo.Context) error {
	filters := make(map[string]interface{})

	// Filtro por tipo (SHIRT, PANTS, SHOES)
	if sizeType := c.QueryParam("type"); sizeType != "" {
		filters["type"] = sizeType
	}

	sizes, err := h.listSizesUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to list sizes", err)
	}

	return response.OK(c, "Sizes retrieved successfully", sizes)
}

// GetByID obtiene una talla por su ID
func (h *SizeHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid size ID", err)
	}

	size, err := h.getSizeUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Size not found")
	}

	return response.OK(c, "Size retrieved successfully", size)
}

// GetByType obtiene tallas por tipo
func (h *SizeHandler) GetByType(c echo.Context) error {
	sizeType := entities.SizeType(c.Param("type"))

	// Validar tipo
	if sizeType != entities.SizeTypeShirt &&
		sizeType != entities.SizeTypePants &&
		sizeType != entities.SizeTypeShoes {
		return response.BadRequest(c, "Invalid size type. Must be SHIRT, PANTS, or SHOES", nil)
	}

	sizes, err := h.getSizesByTypeUC.Execute(c.Request().Context(), sizeType)
	if err != nil {
		return response.InternalServerError(c, "Failed to get sizes by type", err)
	}

	return response.OK(c, "Sizes retrieved successfully", sizes)
}
