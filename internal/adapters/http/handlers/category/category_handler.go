package category

import (
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/category"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	createCategoryUC *category.CreateCategoryUseCase
	getCategoryUC    *category.GetCategoryUseCase
	listCategoriesUC *category.ListCategoriesUseCase
	updateCategoryUC *category.UpdateCategoryUseCase
	deleteCategoryUC *category.DeleteCategoryUseCase
}

func NewCategoryHandler(
	createCategoryUC *category.CreateCategoryUseCase,
	getCategoryUC *category.GetCategoryUseCase,
	listCategoriesUC *category.ListCategoriesUseCase,
	updateCategoryUC *category.UpdateCategoryUseCase,
	deleteCategoryUC *category.DeleteCategoryUseCase,
) *CategoryHandler {
	return &CategoryHandler{
		createCategoryUC: createCategoryUC,
		getCategoryUC:    getCategoryUC,
		listCategoriesUC: listCategoriesUC,
		updateCategoryUC: updateCategoryUC,
		deleteCategoryUC: deleteCategoryUC,
	}
}

func (h *CategoryHandler) Create(c echo.Context) error {
	var cat entities.Category
	if err := c.Bind(&cat); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := h.createCategoryUC.Execute(c.Request().Context(), &cat); err != nil {
		return response.BadRequest(c, "Failed to create category", err)
	}

	return response.Created(c, "Category created successfully", cat)
}

func (h *CategoryHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID", err)
	}

	cat, err := h.getCategoryUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "Category not found")
	}

	return response.OK(c, "Category retrieved successfully", cat)
}

func (h *CategoryHandler) List(c echo.Context) error {
	filters := make(map[string]interface{})

	categories, err := h.listCategoriesUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to list categories", err)
	}

	return response.OK(c, "Categories retrieved successfully", categories)
}

func (h *CategoryHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID", err)
	}

	var cat entities.Category
	if err := c.Bind(&cat); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	cat.ID = uint(id)
	if err := h.updateCategoryUC.Execute(c.Request().Context(), &cat); err != nil {
		return response.BadRequest(c, "Failed to update category", err)
	}

	return response.OK(c, "Category updated successfully", cat)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID", err)
	}

	if err := h.deleteCategoryUC.Execute(c.Request().Context(), uint(id)); err != nil {
		return response.BadRequest(c, "Failed to delete category", err)
	}

	return response.OK(c, "Category deleted successfully", nil)
}
