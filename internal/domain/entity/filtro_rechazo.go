// internal/domain/entity/filtro_rechazo.go
package entity

import (
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// FiltroRechazo represents a rejection filter (capacitor bank).
// Nominal current: I = KVAR / (KV × √3), where KV = Voltaje / 1000.
type FiltroRechazo struct {
	Equipo
	KVAR int
}

func NewFiltroRechazo(clave string, voltaje, kvar, itm, bornes int) (*FiltroRechazo, error) {
	if kvar <= 0 {
		return nil, fmt.Errorf("KVAR debe ser mayor que cero: %d", kvar)
	}
	if voltaje <= 0 {
		return nil, fmt.Errorf("%w: voltaje es %d", ErrDivisionPorCero, voltaje)
	}
	return &FiltroRechazo{
		Equipo: Equipo{
			Clave:   clave,
			Tipo:    TipoFiltroRechazo,
			Voltaje: voltaje,
			ITM:     itm,
			Bornes:  bornes,
		},
		KVAR: kvar,
	}, nil
}

func (fr *FiltroRechazo) CalcularCorrienteNominal() (valueobject.Corriente, error) {
	kv := float64(fr.Voltaje) / 1000.0
	denominador := kv * math.Sqrt(3)
	corriente := float64(fr.KVAR) / denominador
	return valueobject.NewCorriente(corriente)
}
