// internal/domain/valueobject/conductor.go
package valueobject

import (
	"encoding/json"
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
	Material              string
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
	material              string
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
		return Conductor{}, fmt.Errorf("%w: material '%s' no válido (Cu o Al)", ErrConductorInvalido, p.Material)
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
func (c Conductor) Material() string               { return c.material }
func (c Conductor) TipoAislamiento() string        { return c.tipoAislamiento }
func (c Conductor) SeccionMM2() float64            { return c.seccionMM2 }
func (c Conductor) AreaConAislamientoMM2() float64 { return c.areaConAislamientoMM2 }
func (c Conductor) DiametroMM() float64            { return c.diametroMM }
func (c Conductor) NumeroHilos() int               { return c.numeroHilos }
func (c Conductor) ResistenciaPVCPorKm() float64   { return c.resistenciaPVCPorKm }
func (c Conductor) ResistenciaAlPorKm() float64    { return c.resistenciaAlPorKm }
func (c Conductor) ResistenciaAceroPorKm() float64 { return c.resistenciaAceroPorKm }
func (c Conductor) ReactanciaPorKm() float64       { return c.reactanciaPorKm }

// MarshalJSON serializa Conductor a JSON.
func (c Conductor) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Calibre               string  `json:"calibre"`
		Material              string  `json:"material"`
		TipoAislamiento       string  `json:"tipo_aislamiento"`
		SeccionMM2            float64 `json:"seccion_mm2"`
		AreaConAislamientoMM2 float64 `json:"area_con_aislamiento_mm2,omitempty"`
		DiametroMM            float64 `json:"diametro_mm,omitempty"`
		NumeroHilos           int     `json:"numero_hilos,omitempty"`
		ResistenciaPVCPorKm   float64 `json:"resistencia_pvc_por_km,omitempty"`
		ResistenciaAlPorKm    float64 `json:"resistencia_al_por_km,omitempty"`
		ResistenciaAceroPorKm float64 `json:"resistencia_acero_por_km,omitempty"`
		ReactanciaPorKm       float64 `json:"reactancia_por_km,omitempty"`
	}{
		Calibre:               c.calibre,
		Material:              c.material,
		TipoAislamiento:       c.tipoAislamiento,
		SeccionMM2:            c.seccionMM2,
		AreaConAislamientoMM2: c.areaConAislamientoMM2,
		DiametroMM:            c.diametroMM,
		NumeroHilos:           c.numeroHilos,
		ResistenciaPVCPorKm:   c.resistenciaPVCPorKm,
		ResistenciaAlPorKm:    c.resistenciaAlPorKm,
		ResistenciaAceroPorKm: c.resistenciaAceroPorKm,
		ReactanciaPorKm:       c.reactanciaPorKm,
	})
}
