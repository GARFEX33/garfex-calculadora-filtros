// internal/application/dto/errors.go
package dto

import "errors"

// Errores de validación de input.
var (
	ErrEquipoInputInvalido = errors.New("datos de equipo inválidos")
	ErrModoInvalido        = errors.New("modo de cálculo inválido")
)
