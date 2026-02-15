// internal/shared/kernel/valueobject/charola.go
package valueobject

import (
	"errors"
	"fmt"
)

var ErrCableControlInvalido = errors.New("datos de cable de control inválidos")

// CableControl representa un cable de control o comunicación que se transporta en charola.
// Value object inmutable.
type CableControl struct {
	cantidad   int
	diametroMM float64
}

// CableControlParams contiene los parámetros para crear un CableControl.
type CableControlParams struct {
	Cantidad   int
	DiametroMM float64
}

// NewCableControl crea un CableControl value object.
func NewCableControl(p CableControlParams) (CableControl, error) {
	if p.Cantidad < 0 {
		return CableControl{}, fmt.Errorf("%w: cantidad no puede ser negativa", ErrCableControlInvalido)
	}
	if p.DiametroMM <= 0 {
		return CableControl{}, fmt.Errorf("%w: diámetro debe ser mayor que cero", ErrCableControlInvalido)
	}
	return CableControl{
		cantidad:   p.Cantidad,
		diametroMM: p.DiametroMM,
	}, nil
}

func (c CableControl) Cantidad() int       { return c.cantidad }
func (c CableControl) DiametroMM() float64 { return c.diametroMM }

// ConductorCharola representa un conductor con su diámetro exterior para cálculo de espaciado en charola.
// Value object inmutable.
type ConductorCharola struct {
	diametroMM float64
}

// ConductorCharolaParams contiene los parámetros para crear un ConductorCharola.
type ConductorCharolaParams struct {
	DiametroMM float64
}

// NewConductorCharola crea un ConductorCharola value object.
func NewConductorCharola(p ConductorCharolaParams) (ConductorCharola, error) {
	if p.DiametroMM <= 0 {
		return ConductorCharola{}, fmt.Errorf("%w: diámetro debe ser mayor que cero", ErrConductorInvalido)
	}
	return ConductorCharola{diametroMM: p.DiametroMM}, nil
}

func (c ConductorCharola) DiametroMM() float64 { return c.diametroMM }
