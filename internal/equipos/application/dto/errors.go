// internal/equipos/application/dto/errors.go
package dto

import "errors"

// Application-level errors for the equipos feature.
var (
	// ErrEquipoNoEncontrado is returned when an equipo with the given ID does not exist.
	ErrEquipoNoEncontrado = errors.New("equipo no encontrado")

	// ErrClaveYaExiste is returned when trying to create/update with a clave that already exists.
	ErrClaveYaExiste = errors.New("la clave ya existe")

	// ErrInputInvalido is returned when the input DTO fails validation.
	ErrInputInvalido = errors.New("datos de entrada inválidos")

	// ErrIDInvalido is returned when the provided string is not a valid UUID.
	ErrIDInvalido = errors.New("el ID proporcionado no es un UUID válido")
)
