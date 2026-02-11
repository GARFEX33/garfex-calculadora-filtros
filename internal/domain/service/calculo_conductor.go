// internal/domain/service/calculo_conductor.go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

var ErrConductorNoEncontrado = errors.New("no se encontró conductor con capacidad suficiente")

// EntradaTablaConductor represents one row from NOM table 310-15(b)(16).
// Must be sorted smallest-to-largest calibre (as in the NOM table).
// Conductor holds the full physical/electrical properties needed to construct
// a Conductor value object.
type EntradaTablaConductor struct {
	Capacidad float64 // ampacity in amperes
	Conductor valueobject.ConductorParams
}

// SeleccionarConductorAlimentacion picks the smallest conductor from the NOM table
// whose ampacity >= corrienteAjustada / hilosPorFase.
func SeleccionarConductorAlimentacion(
	corrienteAjustada valueobject.Corriente,
	hilosPorFase int,
	tabla []EntradaTablaConductor,
) (valueobject.Conductor, error) {
	if len(tabla) == 0 {
		return valueobject.Conductor{}, fmt.Errorf("%w: tabla vacía", ErrConductorNoEncontrado)
	}

	if hilosPorFase < 1 {
		hilosPorFase = 1
	}

	corrientePorHilo := corrienteAjustada.Valor() / float64(hilosPorFase)

	for _, entrada := range tabla {
		if entrada.Capacidad >= corrientePorHilo {
			return valueobject.NewConductor(entrada.Conductor)
		}
	}

	return valueobject.Conductor{}, fmt.Errorf(
		"%w: corriente por hilo %.2f A excede máxima capacidad de tabla %.2f A",
		ErrConductorNoEncontrado, corrientePorHilo, tabla[len(tabla)-1].Capacidad,
	)
}
