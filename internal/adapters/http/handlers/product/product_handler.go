package product

import (
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/handlers/product/dto"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/product"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	createProductUC *product.CreateProductUseCase
	getProductUC    *product.GetProductUseCase
	listProductsUC  *product.ListProductsUseCase
	updateProductUC *product.UpdateProductUseCase
	deleteProductUC *product.DeleteProductUseCase
	getLowStockUC   *product.GetLowStockProductsUseCase
}

func NewProductHandler(
	createProductUC *product.CreateProductUseCase,
	getProductUC *product.GetProductUseCase,
	listProductsUC *product.ListProductsUseCase,
	updateProductUC *product.UpdateProductUseCase,
	deleteProductUC *product.DeleteProductUseCase,
	getLowStockUC *product.GetLowStockProductsUseCase,
) *ProductHandler {
	return &ProductHandler{
		createProductUC: createProductUC,
		getProductUC:    getProductUC,
		listProductsUC:  listProductsUC,
		updateProductUC: updateProductUC,
		deleteProductUC: deleteProductUC,
		getLowStockUC:   getLowStockUC,
	}
}

func (h *ProductHandler) Create(c echo.Context) error {
	var product dto.Product
	if err := c.Bind(&product); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := h.createProductUC.Execute(c.Request().Context(), dto.FromProductDTO(&product)); err != nil {
		return response.BadRequest(c, "Failed to create product", err)
	}

	return response.Created(c, "Product created successfully", product)
}

func (h *ProductHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	product, err := h.getProductUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Product not found")
	}

	return response.OK(c, "Product retrieved successfully", product)
}

func (h *ProductHandler) List(c echo.Context) error {
	filters := make(map[string]interface{})

	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			filters["category_id"] = uint(id)
		}
	}

	if name := c.QueryParam("name"); name != "" {
		filters["name"] = name
	}

	products, err := h.listProductsUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to list products", err)
	}

	return response.OK(c, "Products retrieved successfully", products)
}

func (h *ProductHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	var product entities.Product
	if err := c.Bind(&product); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	product.ID = uint(id)
	if err := h.updateProductUC.Execute(c.Request().Context(), &product); err != nil {
		return response.BadRequest(c, "Failed to update product", err)
	}

	return response.OK(c, "Product updated successfully", product)
}

func (h *ProductHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID", err)
	}

	if err := h.deleteProductUC.Execute(c.Request().Context(), uint(id)); err != nil {
		return response.BadRequest(c, "Failed to delete product", err)
	}

	return response.OK(c, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetLowStock(c echo.Context) error {
	products, err := h.getLowStockUC.Execute(c.Request().Context())
	if err != nil {
		return response.InternalServerError(c, "Failed to get low stock products", err)
	}

	return response.OK(c, "Low stock products retrieved successfully", products)
}

// GetStats eliminado - dependía de Sales que ya no existe
