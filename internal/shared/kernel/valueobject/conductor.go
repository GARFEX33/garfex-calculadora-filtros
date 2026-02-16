// internal/shared/kernel/valueobject/conductor.go
package valueobject

import (
	"errors"
	"fmt"
)

var ErrConductorInvalido = errors.New("datos de conductor inválidos")

var materialesValidos = map[MaterialConductor]bool{
	MaterialCobre:    true,
	MaterialAluminio: true,
}

// calibresValidos contiene los calibres permitidos según NOM 310-15(b)(16) y 250-122.
// AWG: 14 al 4/0. MCM: 250 al 1000.
var calibresValidos = map[string]bool{
	// AWG
	"14 AWG": true, "12 AWG": true, "10 AWG": true, "8 AWG": true,
	"6 AWG": true, "4 AWG": true, "2 AWG": true,
	"1/0 AWG": true, "2/0 AWG": true, "3/0 AWG": true, "4/0 AWG": true,
	// MCM
	"250 MCM": true, "300 MCM": true, "350 MCM": true, "400 MCM": true,
	"500 MCM": true, "600 MCM": true, "750 MCM": true, "1000 MCM": true,
}

// ConductorParams holds all physical and electrical properties of a conductor
// needed for electrical memory calculations per NOM-001-SEDE-2012.
//
// Required fields: Calibre, Material, SeccionMM2.
// Optional fields (zero-value allowed): TipoAislamiento (empty = bare/desnudo),
// AreaConAislamientoMM2, DiametroMM, NumeroHilos, resistances, reactance.
// Optional fields are validated at the point of use (e.g., conduit sizing
// requires AreaConAislamientoMM2; voltage drop requires resistance values).
type ConductorParams struct {
	Calibre               string
	Material              MaterialConductor
	TipoAislamiento       string  // "" for bare conductors
	SeccionMM2            float64 // sección transversal del conductor (sin aislamiento) [mm²]
	AreaConAislamientoMM2 float64 // área total incluyendo aislamiento, para cálculo de canalización [mm²]
	DiametroMM            float64 // diámetro exterior con aislamiento [mm]
	NumeroHilos           int     // número de hilos del conductor
	ResistenciaPVCPorKm   float64 // resistencia en tubería PVC [Ω/km]
	ResistenciaAlPorKm    float64 // resistencia en tubería de aluminio [Ω/km]
	ResistenciaAceroPorKm float64 // resistencia en tubería de acero [Ω/km]
	ReactanciaPorKm       float64 // reactancia inductiva [Ω/km]
}

// Conductor represents an electrical conductor with its physical and electrical
// properties per NOM-001-SEDE-2012. Immutable.
type Conductor struct {
	calibre               string
	material              MaterialConductor
	tipoAislamiento       string
	seccionMM2            float64
	areaConAislamientoMM2 float64
	diametroMM            float64
	numeroHilos           int
	resistenciaPVCPorKm   float64
	resistenciaAlPorKm    float64
	resistenciaAceroPorKm float64
	reactanciaPorKm       float64
}

func positivoF(v float64, campo string) error {
	if v <= 0 {
		return fmt.Errorf("%w: %s debe ser mayor que cero", ErrConductorInvalido, campo)
	}
	return nil
}

// NewConductor creates a Conductor value object.
// Only Calibre, Material, and SeccionMM2 are required.
// All other fields are optional and validated at the point of use.
func NewConductor(p ConductorParams) (Conductor, error) {
	if !calibresValidos[p.Calibre] {
		return Conductor{}, fmt.Errorf("%w: calibre '%s' no válido según NOM", ErrConductorInvalido, p.Calibre)
	}
	if !materialesValidos[p.Material] {
		return Conductor{}, fmt.Errorf("%w: material '%s' no válido (Cu o Al)", ErrConductorInvalido, p.Material.String())
	}
	if err := positivoF(p.SeccionMM2, "sección"); err != nil {
		return Conductor{}, err
	}
	return Conductor{
		calibre:               p.Calibre,
		material:              p.Material,
		tipoAislamiento:       p.TipoAislamiento,
		seccionMM2:            p.SeccionMM2,
		areaConAislamientoMM2: p.AreaConAislamientoMM2,
		diametroMM:            p.DiametroMM,
		numeroHilos:           p.NumeroHilos,
		resistenciaPVCPorKm:   p.ResistenciaPVCPorKm,
		resistenciaAlPorKm:    p.ResistenciaAlPorKm,
		resistenciaAceroPorKm: p.ResistenciaAceroPorKm,
		reactanciaPorKm:       p.ReactanciaPorKm,
	}, nil
}

func (c Conductor) Calibre() string                { return c.calibre }
func (c Conductor) Material() MaterialConductor    { return c.material }
func (c Conductor) TipoAislamiento() string        { return c.tipoAislamiento }
func (c Conductor) SeccionMM2() float64            { return c.seccionMM2 }
func (c Conductor) AreaConAislamientoMM2() float64 { return c.areaConAislamientoMM2 }
func (c Conductor) DiametroMM() float64            { return c.diametroMM }
func (c Conductor) NumeroHilos() int               { return c.numeroHilos }
func (c Conductor) ResistenciaPVCPorKm() float64   { return c.resistenciaPVCPorKm }
func (c Conductor) ResistenciaAlPorKm() float64    { return c.resistenciaAlPorKm }
func (c Conductor) ResistenciaAceroPorKm() float64 { return c.resistenciaAceroPorKm }
func (c Conductor) ReactanciaPorKm() float64       { return c.reactanciaPorKm }
