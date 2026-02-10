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

// calibresValidos contiene los calibres permitidos según NOM 310-15(b)(16) y 250-122.
// AWG: 18 al 4/0 (los usados en instalaciones eléctricas industriales).
// MCM: 250 al 2000 (conductores de gran capacidad).
var calibresValidos = map[string]bool{
	// AWG
	"18 AWG": true, "16 AWG": true, "14 AWG": true, "12 AWG": true,
	"10 AWG": true, "8 AWG": true, "6 AWG": true, "4 AWG": true,
	"2 AWG": true, "1/0 AWG": true, "2/0 AWG": true, "3/0 AWG": true, "4/0 AWG": true,
	// MCM
	"250 MCM": true, "300 MCM": true, "350 MCM": true, "400 MCM": true,
	"500 MCM": true, "600 MCM": true, "700 MCM": true, "750 MCM": true,
	"800 MCM": true, "900 MCM": true, "1000 MCM": true, "1250 MCM": true,
	"1500 MCM": true, "1750 MCM": true, "2000 MCM": true,
}

// Conductor represents an electrical conductor with its physical properties. Immutable.
type Conductor struct {
	calibre         string
	material        string
	tipoAislamiento string
	seccionMM2      float64
}

func NewConductor(calibre, material, tipoAislamiento string, seccionMM2 float64) (Conductor, error) {
	if !calibresValidos[calibre] {
		return Conductor{}, fmt.Errorf("%w: calibre '%s' no válido según NOM", ErrConductorInvalido, calibre)
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
