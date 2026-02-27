// internal/equipos/domain/entity/errors.go
package entity

import "errors"

// Domain errors for the equipos feature.
var (
	// ErrTipoFiltroInvalido is returned when a TipoFiltro string is not recognized.
	ErrTipoFiltroInvalido = errors.New("tipo de filtro inválido")

	// ErrConexionInvalida is returned when a Conexion string is not recognized.
	ErrConexionInvalida = errors.New("tipo de conexión inválido")

	// ErrVoltajeInvalido is returned when voltaje is <= 0.
	ErrVoltajeInvalido = errors.New("el voltaje debe ser mayor que cero")

	// ErrAmperajeInvalido is returned when amperaje (qn/In) is <= 0.
	ErrAmperajeInvalido = errors.New("el amperaje debe ser mayor que cero")

	// ErrITMInvalido is returned when the ITM value is <= 0.
	ErrITMInvalido = errors.New("el ITM debe ser mayor que cero")

	// ErrTipoVoltajeInvalido is returned when a TipoVoltaje string is not recognized.
	ErrTipoVoltajeInvalido = errors.New("tipo de voltaje inválido")
)
