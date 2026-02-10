// internal/domain/valueobject/conductor.go
package valueobject

import (
	"errors"
	"fmt"
)

var ErrConductorInvalido = errors.New("datos de conductor inválidos")

var materialesValidos = map[string]bool{
	"Cu": true,
	"Al": true,
}

// Conductor represents an electrical conductor with its physical properties. Immutable.
type Conductor struct {
	calibre         string
	material        string
	tipoAislamiento string
	seccionMM2      float64
}

func NewConductor(calibre, material, tipoAislamiento string, seccionMM2 float64) (Conductor, error) {
	if calibre == "" {
		return Conductor{}, fmt.Errorf("%w: calibre vacío", ErrConductorInvalido)
	}
	if !materialesValidos[material] {
		return Conductor{}, fmt.Errorf("%w: material '%s' no válido (Cu o Al)", ErrConductorInvalido, material)
	}
	if tipoAislamiento == "" {
		return Conductor{}, fmt.Errorf("%w: tipo de aislamiento vacío", ErrConductorInvalido)
	}
	if seccionMM2 <= 0 {
		return Conductor{}, fmt.Errorf("%w: sección debe ser mayor que cero", ErrConductorInvalido)
	}
	return Conductor{
		calibre:         calibre,
		material:        material,
		tipoAislamiento: tipoAislamiento,
		seccionMM2:      seccionMM2,
	}, nil
}

func (c Conductor) Calibre() string        { return c.calibre }
func (c Conductor) Material() string        { return c.material }
func (c Conductor) TipoAislamiento() string { return c.tipoAislamiento }
func (c Conductor) SeccionMM2() float64     { return c.seccionMM2 }
