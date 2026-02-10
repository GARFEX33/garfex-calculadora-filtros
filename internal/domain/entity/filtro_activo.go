// internal/domain/entity/filtro_activo.go
package entity

import (
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// FiltroActivo represents an active filter. Its nominal current equals
// its amperage rating directly (no formula needed).
type FiltroActivo struct {
	Equipo
	Amperaje int
}

func NewFiltroActivo(clave string, voltaje, amperaje int, itm ITM) (*FiltroActivo, error) {
	if amperaje <= 0 {
		return nil, fmt.Errorf("amperaje debe ser mayor que cero: %d", amperaje)
	}
	return &FiltroActivo{
		Equipo: Equipo{
			Clave:   clave,
			Tipo:    TipoFiltroActivo,
			Voltaje: voltaje,
			ITM:     itm,
		},
		Amperaje: amperaje,
	}, nil
}

func (fa *FiltroActivo) CalcularCorrienteNominal() (valueobject.Corriente, error) {
	return valueobject.NewCorriente(float64(fa.Amperaje))
}

// PotenciaKVA returns apparent power: I × V × √3 / 1000 [kVA]
func (fa *FiltroActivo) PotenciaKVA() float64 {
	return float64(fa.Amperaje) * float64(fa.Voltaje) * math.Sqrt(3) / 1000.0
}

// PotenciaKW returns active power. FiltroActivo has PF=1, so kW = kVA.
func (fa *FiltroActivo) PotenciaKW() float64 {
	return fa.PotenciaKVA()
}

// PotenciaKVAR returns reactive power. FiltroActivo has PF=1, so kVAR = 0.
func (fa *FiltroActivo) PotenciaKVAR() float64 {
	return 0
}
