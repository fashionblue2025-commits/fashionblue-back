#!/bin/bash

echo "🔧 Creando handlers para Suppliers y Capital Injections..."

# Handler - Supplier
cat > internal/adapters/http/handlers/supplier/supplier_handler.go << 'EOF'
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
EOF

# Handler - Capital Injection
cat > internal/adapters/http/handlers/capital_injection/capital_injection_handler.go << 'EOF'
package capital_injection

import (
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/capital_injection"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type CapitalInjectionHandler struct {
	createInjectionUC *capital_injection.CreateInjectionUseCase
	getInjectionUC    *capital_injection.GetInjectionUseCase
	listInjectionsUC  *capital_injection.ListInjectionsUseCase
	getTotalCapitalUC *capital_injection.GetTotalCapitalUseCase
}

func NewCapitalInjectionHandler(
	createInjectionUC *capital_injection.CreateInjectionUseCase,
	getInjectionUC *capital_injection.GetInjectionUseCase,
	listInjectionsUC *capital_injection.ListInjectionsUseCase,
	getTotalCapitalUC *capital_injection.GetTotalCapitalUseCase,
) *CapitalInjectionHandler {
	return &CapitalInjectionHandler{
		createInjectionUC: createInjectionUC,
		getInjectionUC:    getInjectionUC,
		listInjectionsUC:  listInjectionsUC,
		getTotalCapitalUC: getTotalCapitalUC,
	}
}

func (h *CapitalInjectionHandler) Create(c echo.Context) error {
	var injection entities.CapitalInjection
	if err := c.Bind(&injection); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := h.createInjectionUC.Execute(c.Request().Context(), &injection); err != nil {
		return response.BadRequest(c, "Failed to create capital injection", err)
	}

	return response.Created(c, "Capital injection created successfully", injection)
}

func (h *CapitalInjectionHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid injection ID", err)
	}

	injection, err := h.getInjectionUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Capital injection not found")
	}

	return response.OK(c, "Capital injection retrieved successfully", injection)
}

func (h *CapitalInjectionHandler) List(c echo.Context) error {
	filters := make(map[string]interface{})

	injections, err := h.listInjectionsUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve capital injections", err)
	}

	return response.OK(c, "Capital injections retrieved successfully", injections)
}

func (h *CapitalInjectionHandler) GetTotal(c echo.Context) error {
	total, err := h.getTotalCapitalUC.Execute(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, "Failed to get total capital", err)
	}

	return response.OK(c, "Total capital retrieved successfully", map[string]interface{}{
		"total": total,
	})
}
EOF

echo "✅ Handlers creados exitosamente"
