package supplier

import (
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/supplier"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type SupplierHandler struct {
	createSupplierUC *supplier.CreateSupplierUseCase
	getSupplierUC    *supplier.GetSupplierUseCase
	listSuppliersUC  *supplier.ListSuppliersUseCase
	updateSupplierUC *supplier.UpdateSupplierUseCase
	deleteSupplierUC *supplier.DeleteSupplierUseCase
}

func NewSupplierHandler(
	createSupplierUC *supplier.CreateSupplierUseCase,
	getSupplierUC *supplier.GetSupplierUseCase,
	listSuppliersUC *supplier.ListSuppliersUseCase,
	updateSupplierUC *supplier.UpdateSupplierUseCase,
	deleteSupplierUC *supplier.DeleteSupplierUseCase,
) *SupplierHandler {
	return &SupplierHandler{
		createSupplierUC: createSupplierUC,
		getSupplierUC:    getSupplierUC,
		listSuppliersUC:  listSuppliersUC,
		updateSupplierUC: updateSupplierUC,
		deleteSupplierUC: deleteSupplierUC,
	}
}

func (h *SupplierHandler) Create(c echo.Context) error {
	var supplier entities.Supplier
	if err := c.Bind(&supplier); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := h.createSupplierUC.Execute(c.Request().Context(), &supplier); err != nil {
		return response.BadRequest(c, "Failed to create supplier", err)
	}

	return response.Created(c, "Supplier created successfully", supplier)
}

func (h *SupplierHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid supplier ID", err)
	}

	supplier, err := h.getSupplierUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Supplier not found")
	}

	return response.OK(c, "Supplier retrieved successfully", supplier)
}

func (h *SupplierHandler) List(c echo.Context) error {
	filters := make(map[string]interface{})

	if name := c.QueryParam("name"); name != "" {
		filters["name"] = name
	}

	suppliers, err := h.listSuppliersUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve suppliers", err)
	}

	return response.OK(c, "Suppliers retrieved successfully", suppliers)
}

func (h *SupplierHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid supplier ID", err)
	}

	var supplier entities.Supplier
	if err := c.Bind(&supplier); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	supplier.ID = uint(id)

	if err := h.updateSupplierUC.Execute(c.Request().Context(), &supplier); err != nil {
		return response.BadRequest(c, "Failed to update supplier", err)
	}

	return response.OK(c, "Supplier updated successfully", supplier)
}

func (h *SupplierHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid supplier ID", err)
	}

	if err := h.deleteSupplierUC.Execute(c.Request().Context(), uint(id)); err != nil {
		return response.InternalServerError(c, "Failed to delete supplier", err)
	}

	return response.OK(c, "Supplier deleted successfully", nil)
}
