package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response estructura estándar de respuesta
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success respuesta exitosa
func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error respuesta de error
func Error(c echo.Context, statusCode int, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	return c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// BadRequest respuesta de solicitud incorrecta
func BadRequest(c echo.Context, message string, err error) error {
	return Error(c, http.StatusBadRequest, message, err)
}

// Unauthorized respuesta no autorizado
func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden respuesta prohibido
func Forbidden(c echo.Context, message string) error {
	return Error(c, http.StatusForbidden, message, nil)
}

// NotFound respuesta no encontrado
func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message, nil)
}

// InternalServerError respuesta error interno del servidor
func InternalServerError(c echo.Context, message string, err error) error {
	return Error(c, http.StatusInternalServerError, message, err)
}

// Created respuesta de recurso creado
func Created(c echo.Context, message string, data interface{}) error {
	return Success(c, http.StatusCreated, message, data)
}

// OK respuesta exitosa
func OK(c echo.Context, message string, data interface{}) error {
	return Success(c, http.StatusOK, message, data)
}

// NoContent respuesta sin contenido
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
