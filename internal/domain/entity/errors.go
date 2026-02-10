// internal/domain/entity/errors.go
package entity

import "errors"

var (
	ErrTipoFiltroInvalido = errors.New("tipo de filtro no válido")
	ErrDivisionPorCero    = errors.New("división por cero en cálculo de corriente")
)
