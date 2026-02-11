// internal/domain/service/calculo_tierra.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// EntradaTablaTierra represents one row from NOM table 250-122.
// Entries must be sorted by ITMHasta ascending.
// Conductor holds the full physical/electrical properties needed to construct
// a Conductor value object.
type EntradaTablaTierra struct {
	ITMHasta  int
	Conductor valueobject.ConductorParams
}

// SeleccionarConductorTierra selects the ground conductor from NOM table 250-122
// based on the equipment's ITM (circuit breaker) rating.
func SeleccionarConductorTierra(itm int, tabla []EntradaTablaTierra) (valueobject.Conductor, error) {
	if itm <= 0 {
		return valueobject.Conductor{}, fmt.Errorf("ITM debe ser mayor que cero: %d", itm)
	}
	if len(tabla) == 0 {
		return valueobject.Conductor{}, fmt.Errorf("%w: tabla de tierra vacía", ErrConductorNoEncontrado)
	}

	for _, entrada := range tabla {
		if itm <= entrada.ITMHasta {
			return valueobject.NewConductor(entrada.Conductor)
		}
	}

	return valueobject.Conductor{}, fmt.Errorf(
		"%w: ITM %d excede máximo de tabla %d",
		ErrConductorNoEncontrado, itm, tabla[len(tabla)-1].ITMHasta,
	)
}
