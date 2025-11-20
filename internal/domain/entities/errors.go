package entities

import "errors"

var (
	// ErrInvalidInput indica que los datos de entrada son inválidos
	ErrInvalidInput = errors.New("invalid input data")

	// ErrNotFound indica que el recurso no fue encontrado
	ErrNotFound = errors.New("resource not found")

	// ErrUnauthorized indica que el usuario no está autorizado
	ErrUnauthorized = errors.New("unauthorized")

	// ErrAlreadyExists indica que el recurso ya existe
	ErrAlreadyExists = errors.New("resource already exists")
)
