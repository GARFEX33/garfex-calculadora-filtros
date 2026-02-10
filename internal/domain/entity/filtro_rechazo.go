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

func NewFiltroRechazo(clave string, voltaje, kvar int, itm ITM) (*FiltroRechazo, error) {
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

// PotenciaKVAR returns reactive power directly from KVAR rating [kVAR].
func (fr *FiltroRechazo) PotenciaKVAR() float64 {
	return float64(fr.KVAR)
}

// PotenciaKVA returns apparent power. FiltroRechazo is purely reactive, so kVA = kVAR.
func (fr *FiltroRechazo) PotenciaKVA() float64 {
	return float64(fr.KVAR)
}

// PotenciaKW returns active power. FiltroRechazo is purely reactive, so kW = 0.
func (fr *FiltroRechazo) PotenciaKW() float64 {
	return 0
}
