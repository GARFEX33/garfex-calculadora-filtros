// internal/shared/kernel/valueobject/resistencia_reactancia.go
package valueobject

import (
	"errors"
	"fmt"
)

// ErrImpedanciaInvalida is returned when resistance or reactance values are negative.
var ErrImpedanciaInvalida = errors.New("valores de impedancia inválidos")

// ResistenciaReactancia holds the impedance values for voltage drop calculation
// per Tabla 9 of NOM-001-SEDE-2012. R and X are in Ohms per km. Immutable.
type ResistenciaReactancia struct {
	r float64 // resistance [Ω/km]
	x float64 // inductive reactance [Ω/km]
}

// NewResistenciaReactancia constructs a validated ResistenciaReactancia.
// Both R and X must be non-negative (zero is allowed for bare conductors or lossless lines).
func NewResistenciaReactancia(r, x float64) (ResistenciaReactancia, error) {
	if r < 0 {
		return ResistenciaReactancia{}, fmt.Errorf("%w: resistencia no puede ser negativa: %.4f", ErrImpedanciaInvalida, r)
	}
	if x < 0 {
		return ResistenciaReactancia{}, fmt.Errorf("%w: reactancia no puede ser negativa: %.4f", ErrImpedanciaInvalida, x)
	}
	return ResistenciaReactancia{r: r, x: x}, nil
}

// R returns the resistance value in Ohms per km.
func (rr ResistenciaReactancia) R() float64 { return rr.r }

// X returns the inductive reactance value in Ohms per km.
func (rr ResistenciaReactancia) X() float64 { return rr.x }
