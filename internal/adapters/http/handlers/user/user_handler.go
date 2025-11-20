package user

import (
	"strconv"

	"github.com/bryanarroyaveortiz/fashion-blue/internal/application/usecases/user"
	"github.com/bryanarroyaveortiz/fashion-blue/internal/domain/entities"
	"github.com/bryanarroyaveortiz/fashion-blue/pkg/response"
	"github.com/labstack/echo/v4"
)

// UserHandler maneja las peticiones HTTP relacionadas con usuarios
type UserHandler struct {
	createUserUC     *user.CreateUserUseCase
	getUserUC        *user.GetUserUseCase
	listUsersUC      *user.ListUsersUseCase
	updateUserUC     *user.UpdateUserUseCase
	deleteUserUC     *user.DeleteUserUseCase
	changePasswordUC *user.ChangePasswordUseCase
}

// NewUserHandler crea una nueva instancia del handler
func NewUserHandler(
	createUserUC *user.CreateUserUseCase,
	getUserUC *user.GetUserUseCase,
	listUsersUC *user.ListUsersUseCase,
	updateUserUC *user.UpdateUserUseCase,
	deleteUserUC *user.DeleteUserUseCase,
	changePasswordUC *user.ChangePasswordUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUC:     createUserUC,
		getUserUC:        getUserUC,
		listUsersUC:      listUsersUC,
		updateUserUC:     updateUserUC,
		deleteUserUC:     deleteUserUC,
		changePasswordUC: changePasswordUC,
	}
}

// CreateUserRequest representa la petición para crear un usuario
type CreateUserRequest struct {
	Email     string            `json:"email" validate:"required,email"`
	Password  string            `json:"password" validate:"required,min=6"`
	FirstName string            `json:"first_name" validate:"required"`
	LastName  string            `json:"last_name" validate:"required"`
	Role      entities.UserRole `json:"role" validate:"required"`
}

// Create maneja la creación de un usuario
func (h *UserHandler) Create(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	user := &entities.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	if err := h.createUserUC.Execute(c.Request().Context(), user, req.Password); err != nil {
		return response.BadRequest(c, "Failed to create user", err)
	}

	return response.Created(c, "User created successfully", user)
}

// GetByID maneja la obtención de un usuario por ID
func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	user, err := h.getUserUC.Execute(c.Request().Context(), uint(id))
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.OK(c, "User retrieved successfully", user)
}

// List maneja el listado de usuarios
func (h *UserHandler) List(c echo.Context) error {
	filters := make(map[string]interface{})

	if role := c.QueryParam("role"); role != "" {
		filters["role"] = role
	}

	users, err := h.listUsersUC.Execute(c.Request().Context(), filters)
	if err != nil {
		return response.InternalServerError(c, "Failed to list users", err)
	}

	return response.OK(c, "Users retrieved successfully", users)
}

// Update maneja la actualización de un usuario
func (h *UserHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	var user entities.User
	if err := c.Bind(&user); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	user.ID = uint(id)
	if err := h.updateUserUC.Execute(c.Request().Context(), &user); err != nil {
		return response.BadRequest(c, "Failed to update user", err)
	}

	return response.OK(c, "User updated successfully", user)
}

// Delete maneja la eliminación de un usuario
func (h *UserHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	if err := h.deleteUserUC.Execute(c.Request().Context(), uint(id)); err != nil {
		return response.BadRequest(c, "Failed to delete user", err)
	}

	return response.OK(c, "User deleted successfully", nil)
}

// ChangePasswordRequest representa la petición para cambiar contraseña
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// ChangePassword maneja el cambio de contraseña
func (h *UserHandler) ChangePassword(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err)
	}

	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err)
	}

	if err := h.changePasswordUC.Execute(c.Request().Context(), uint(id), req.OldPassword, req.NewPassword); err != nil {
		return response.BadRequest(c, "Failed to change password", err)
	}

	return response.OK(c, "Password changed successfully", nil)
}
