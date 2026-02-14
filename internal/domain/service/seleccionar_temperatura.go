// internal/domain/service/seleccionar_temperatura.go
package service

import (
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// SeleccionarTemperatura determines the temperature column according to NOM rules.
//
// Rules per NOM-001-SEDE-2012:
// - Corriente <= 100A -> 60°C (or 75°C if charola triangular without 60°C column)
// - Corriente > 100A -> 75°C
// - If temperaturaOverride is provided, it takes precedence
func SeleccionarTemperatura(
	corriente valueobject.Corriente,
	tipoCanalizacion entity.TipoCanalizacion,
	override *valueobject.Temperatura,
) valueobject.Temperatura {
	// If there's an explicit override, use it
	if override != nil {
		return *override
	}

	// NOM rules based on current
	if corriente.Valor() <= 100 {
		// <= 100A -> 60°C (or 75°C if charola triangular)
		if tipoCanalizacion == entity.TipoCanalizacionCharolaCableTriangular {
			return valueobject.Temp75
		}
		return valueobject.Temp60
	}

	// > 100A -> 75°C
	return valueobject.Temp75
}
