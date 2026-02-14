// internal/application/port/seleccionar_temperatura.go
package port

import (
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// SeleccionarTemperaturaPort defines the contract for temperature selection
// based on NOM electrical rules.
type SeleccionarTemperaturaPort interface {
	// SeleccionarTemperatura returns the appropriate temperature column
	// based on current, conduit type, and optional override.
	SeleccionarTemperatura(
		corriente valueobject.Corriente,
		tipoCanalizacion entity.TipoCanalizacion,
		temperaturaOverride *valueobject.Temperatura,
	) valueobject.Temperatura
}
