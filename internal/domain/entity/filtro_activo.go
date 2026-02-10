// internal/domain/entity/filtro_activo.go
package entity

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// FiltroActivo represents an active filter. Its nominal current equals
// its amperage rating directly (no formula needed).
type FiltroActivo struct {
	Equipo
	Amperaje int
}

func NewFiltroActivo(clave string, voltaje, amperaje, itm, bornes int) (*FiltroActivo, error) {
	if amperaje <= 0 {
		return nil, fmt.Errorf("amperaje debe ser mayor que cero: %d", amperaje)
	}
	return &FiltroActivo{
		Equipo: Equipo{
			Clave:   clave,
			Tipo:    TipoFiltroActivo,
			Voltaje: voltaje,
			ITM:     itm,
			Bornes:  bornes,
		},
		Amperaje: amperaje,
	}, nil
}

func (fa *FiltroActivo) CalcularCorrienteNominal() (valueobject.Corriente, error) {
	return valueobject.NewCorriente(float64(fa.Amperaje))
}
