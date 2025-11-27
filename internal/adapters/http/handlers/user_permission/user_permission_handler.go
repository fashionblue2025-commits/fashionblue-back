package userpermission

import (
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/dto"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/adapters/http/middleware"
	userpermission "github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/user_permission"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

type UserPermissionHandler struct {
	managePermissionsUC    *userpermission.ManageUserPermissionsUseCase
	checkPermissionUC      *userpermission.CheckCategoryPermissionUseCase
	getAllowedCategoriesUC *userpermission.GetUserAllowedCategoriesUseCase
}

func NewUserPermissionHandler(
	managePermissionsUC *userpermission.ManageUserPermissionsUseCase,
	checkPermissionUC *userpermission.CheckCategoryPermissionUseCase,
	getAllowedCategoriesUC *userpermission.GetUserAllowedCategoriesUseCase,
) *UserPermissionHandler {
	return &UserPermissionHandler{
		managePermissionsUC:    managePermissionsUC,
		checkPermissionUC:      checkPermissionUC,
		getAllowedCategoriesUC: getAllowedCategoriesUC,
	}
}

// GetUserPermissions obtiene todos los permisos de un usuario
func (h *UserPermissionHandler) GetUserPermissions(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	permissions, err := h.managePermissionsUC.GetUserPermissions(c.Request().Context(), uint(userID))
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve permissions", err)
	}

	permissionDTOs := dto.ToUserCategoryPermissionDTOList(permissions)
	return response.OK(c, "Permissions retrieved successfully", permissionDTOs)
}

// SetUserPermissions establece todos los permisos de un usuario (reemplaza existentes)
func (h *UserPermissionHandler) SetUserPermissions(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	var req dto.SetUserPermissionsRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	// Convertir DTOs a entidades
	permissions := make([]entities.UserCategoryPermission, len(req.Permissions))
	for i, item := range req.Permissions {
		permissions[i] = entities.UserCategoryPermission{
			UserID:     uint(userID),
			CategoryID: item.CategoryID,
			CanView:    item.CanView,
			CanCreate:  item.CanCreate,
			CanEdit:    item.CanEdit,
			CanDelete:  item.CanDelete,
		}
	}

	if err := h.managePermissionsUC.SetUserPermissions(c.Request().Context(), uint(userID), permissions); err != nil {
		return response.InternalServerError(c, "Failed to set permissions", err)
	}

	return response.OK(c, "Permissions updated successfully", nil)
}

// AddCategoryPermission agrega o actualiza un permiso específico
func (h *UserPermissionHandler) AddCategoryPermission(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	categoryID, err := strconv.ParseUint(c.Param("categoryId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID", err)
	}

	var req dto.UpdateUserCategoryPermissionRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	permission := &entities.UserCategoryPermission{
		UserID:     uint(userID),
		CategoryID: uint(categoryID),
		CanView:    req.CanView,
		CanCreate:  req.CanCreate,
		CanEdit:    req.CanEdit,
		CanDelete:  req.CanDelete,
	}

	if err := h.managePermissionsUC.AddCategoryPermission(c.Request().Context(), permission); err != nil {
		return response.InternalServerError(c, "Failed to add permission", err)
	}

	return response.OK(c, "Permission added successfully", nil)
}

// RemoveCategoryPermission elimina un permiso específico
func (h *UserPermissionHandler) RemoveCategoryPermission(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	categoryID, err := strconv.ParseUint(c.Param("categoryId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID", err)
	}

	if err := h.managePermissionsUC.RemoveCategoryPermission(c.Request().Context(), uint(userID), uint(categoryID)); err != nil {
		return response.InternalServerError(c, "Failed to remove permission", err)
	}

	return response.OK(c, "Permission removed successfully", nil)
}

// GetAllowedCategories obtiene las categorías a las que un usuario tiene acceso
func (h *UserPermissionHandler) GetAllowedCategories(c echo.Context) error {
	// Obtener ID del usuario (puede ser "me" para el usuario actual)
	userIDParam := c.Param("id")
	var userID uint

	if userIDParam == "me" {
		// Obtener usuario del contexto (autenticado)
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return response.Unauthorized(c, "User not authenticated")
		}
		userID = user.ID
	} else {
		parsedID, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			return response.BadRequest(c, "Invalid user ID", err)
		}
		userID = uint(parsedID)
	}

	// Obtener acción (por defecto: view)
	action := c.QueryParam("action")
	if action == "" {
		action = "view"
	}

	// Validar acción
	validActions := map[string]bool{
		"view":   true,
		"create": true,
		"edit":   true,
		"delete": true,
	}
	if !validActions[action] {
		return response.BadRequest(c, "Invalid action. Must be: view, create, edit, or delete", nil)
	}

	categories, err := h.getAllowedCategoriesUC.Execute(c.Request().Context(), userID, action)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve allowed categories", err)
	}

	categoryDTOs := dto.ToCategoryDTOList(categories)
	return response.OK(c, "Allowed categories retrieved successfully", categoryDTOs)
}

// GetCategoryPermissions obtiene todos los usuarios con permisos sobre una categoría
func (h *UserPermissionHandler) GetCategoryPermissions(c echo.Context) error {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID", err)
	}

	permissions, err := h.managePermissionsUC.GetCategoryPermissions(c.Request().Context(), uint(categoryID))
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve permissions", err)
	}

	permissionDTOs := dto.ToUserCategoryPermissionDTOList(permissions)
	return response.OK(c, "Permissions retrieved successfully", permissionDTOs)
}

// CheckPermission verifica si un usuario tiene un permiso específico sobre una categoría
func (h *UserPermissionHandler) CheckPermission(c echo.Context) error {
	userIDParam := c.Param("id")
	var userID uint

	if userIDParam == "me" {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return response.Unauthorized(c, "User not authenticated")
		}
		userID = user.ID
	} else {
		parsedID, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			return response.BadRequest(c, "Invalid user ID", err)
		}
		userID = uint(parsedID)
	}

	categoryID, err := strconv.ParseUint(c.Param("categoryId"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID", err)
	}

	action := c.QueryParam("action")
	if action == "" {
		action = "view"
	}

	hasPermission, err := h.checkPermissionUC.Execute(c.Request().Context(), userID, uint(categoryID), action)
	if err != nil {
		return response.InternalServerError(c, "Failed to check permission", err)
	}

	result := map[string]interface{}{
		"has_permission": hasPermission,
		"user_id":        userID,
		"category_id":    categoryID,
		"action":         action,
	}

	return response.OK(c, "Permission checked successfully", result)
}
