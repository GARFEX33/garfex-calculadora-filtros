// internal/calculos/domain/service/calculo_conductor.go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

var ErrConductorNoEncontrado = errors.New("no se encontró conductor con capacidad suficiente")

// SeleccionarConductorAlimentacion picks the smallest conductor from the NOM table
// whose ampacity >= corrienteAjustada / hilosPorFase.
func SeleccionarConductorAlimentacion(
	corrienteAjustada valueobject.Corriente,
	hilosPorFase int,
	tabla []valueobject.EntradaTablaConductor,
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
