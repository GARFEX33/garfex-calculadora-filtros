// internal/shared/kernel/valueobject/potencia.go
package valueobject

import (
	"errors"
	"fmt"
)

// ErrPotenciaInvalida is returned when potencia value is invalid.
var ErrPotenciaInvalida = errors.New("potencia debe ser mayor que cero")

// ErrUnidadPotenciaInvalida is returned when potencia unit is not recognized.
var ErrUnidadPotenciaInvalida = errors.New("unidad de potencia inválida: use W, KW, KVA, o KVAR")

// UnidadPotencia represents the unit of electrical power.
type UnidadPotencia string

const (
	UnidadPotenciaW    UnidadPotencia = "W"
	UnidadPotenciaKW   UnidadPotencia = "KW"
	UnidadPotenciaKVA  UnidadPotencia = "KVA"
	UnidadPotenciaKVAR UnidadPotencia = "KVAR"
)

// ParseUnidadPotencia converts a string to UnidadPotencia.
func ParseUnidadPotencia(s string) (UnidadPotencia, error) {
	switch s {
	case "W", "w":
		return UnidadPotenciaW, nil
	case "KW", "kw", "kW":
		return UnidadPotenciaKW, nil
	case "KVA", "kva", "Kva":
		return UnidadPotenciaKVA, nil
	case "KVAR", "kvar", "KVAr":
		return UnidadPotenciaKVAR, nil
	default:
		return "", ErrUnidadPotenciaInvalida
	}
}

// Potencia represents an electrical power value. Immutable.
// Internally stores as watts (W).
type Potencia struct {
	valor  float64
	unidad UnidadPotencia
}

// NewPotencia creates a Potencia value object from a value and unit.
// The value is normalized to watts internally.
func NewPotencia(valor float64, unidad string) (Potencia, error) {
	if valor <= 0 {
		return Potencia{}, ErrPotenciaInvalida
	}

	unidadParsed, err := ParseUnidadPotencia(unidad)
	if err != nil {
		return Potencia{}, err
	}

	// Normalize to watts
	watts := normalizarAWatts(valor, unidadParsed)

	return Potencia{valor: watts, unidad: unidadParsed}, nil
}

// normalizarAWatts converts a value from any unit to watts.
func normalizarAWatts(valor float64, unidad UnidadPotencia) float64 {
	switch unidad {
	case UnidadPotenciaW:
		return valor
	case UnidadPotenciaKW:
		return valor * 1000.0
	case UnidadPotenciaKVA:
		// KVA = KW (assuming PF=1 for simplicity, domain handles the conversion)
		return valor * 1000.0
	case UnidadPotenciaKVAR:
		// KVAR is reactive power - stored as VAR for consistency
		return valor * 1000.0
	default:
		return valor
	}
}

// Valor returns the power value in watts.
func (p Potencia) Valor() float64 {
	return p.valor
}

// Unidad returns the original unit.
func (p Potencia) Unidad() UnidadPotencia {
	return p.unidad
}

// KW returns the power in kilowatts.
func (p Potencia) KW() float64 {
	return p.valor / 1000.0
}

// KVA returns the power in kilovolt-amperes.
func (p Potencia) KVA() float64 {
	return p.valor / 1000.0
}

// KVAR returns the power in kilovolt-amperes reactive.
func (p Potencia) KVAR() float64 {
	return p.valor / 1000.0
}

// String returns a string representation.
func (p Potencia) String() string {
	return fmt.Sprintf("%.2f W", p.valor)
}
