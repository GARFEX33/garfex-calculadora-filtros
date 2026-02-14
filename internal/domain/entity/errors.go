// internal/domain/entity/errors.go
package entity

import "errors"

var (
	ErrTipoEquipoInvalido = errors.New("tipo de equipo no válido")
	ErrDivisionPorCero    = errors.New("división por cero en cálculo de corriente")
)
