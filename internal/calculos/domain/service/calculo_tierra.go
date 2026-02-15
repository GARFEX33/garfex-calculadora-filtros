// internal/calculos/domain/service/calculo_tierra.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// SeleccionarConductorTierra selects the ground conductor from NOM table 250-122
// based on the equipment's ITM rating and the desired conductor material.
// If material is Al but Al is not available for the given ITM range (NOM restriction),
// it silently falls back to Cu.
func SeleccionarConductorTierra(
	itm int,
	material valueobject.MaterialConductor,
	tabla []valueobject.EntradaTablaTierra,
) (valueobject.Conductor, error) {
	if itm <= 0 {
		return valueobject.Conductor{}, fmt.Errorf("ITM debe ser mayor que cero: %d", itm)
	}
	if len(tabla) == 0 {
		return valueobject.Conductor{}, fmt.Errorf("%w: tabla de tierra vacía", ErrConductorNoEncontrado)
	}

	for _, entrada := range tabla {
		if itm <= entrada.ITMHasta {
			if material == valueobject.MaterialAluminio && entrada.ConductorAl != nil {
				return valueobject.NewConductor(*entrada.ConductorAl)
			}
			return valueobject.NewConductor(entrada.ConductorCu)
		}
	}

	return valueobject.Conductor{}, fmt.Errorf(
		"%w: ITM %d excede máximo de tabla %d",
		ErrConductorNoEncontrado, itm, tabla[len(tabla)-1].ITMHasta,
	)
}
