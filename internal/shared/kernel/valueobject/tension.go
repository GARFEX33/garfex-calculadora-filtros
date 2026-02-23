// internal/shared/kernel/valueobject/tension.go
package valueobject

import (
	"errors"
	"fmt"
	"math"
)

var ErrVoltajeInvalido = errors.New("voltaje no válido según normativa NOM")
var ErrUnidadTensionInvalida = errors.New("unidad de tensión inválida: use V o kV")

// UnidadTension representa la unidad de voltaje.
type UnidadTension string

const (
	UnidadTensionV  UnidadTension = "V"
	UnidadTensionkV UnidadTension = "kV"
)

// ParseUnidadTension convierte un string a UnidadTension.
func ParseUnidadTension(s string) (UnidadTension, error) {
	switch s {
	case "V", "v":
		return UnidadTensionV, nil
	case "kV", "KV", "kv", "Kv":
		return UnidadTensionkV, nil
	case "":
		// Default a V para compatibilidad hacia atrás
		return UnidadTensionV, nil
	default:
		return "", ErrUnidadTensionInvalida
	}
}

var voltajesValidos = map[int]bool{
	127: true,
	220: true,
	240: true,
	277: true,
	440: true,
	480: true,
	600: true,
}

// Tension represents an electrical voltage value in Volts. Immutable.
// Internally stores as volts (V).
type Tension struct {
	valor  int           // siempre en volts
	unidad UnidadTension // unidad original
}

// normalizarAVolts convierte un valor de cualquier unidad a volts.
func normalizarAVolts(valor float64, unidad UnidadTension) int {
	switch unidad {
	case UnidadTensionV:
		return int(math.Round(valor))
	case UnidadTensionkV:
		// kV a V: multiplicar por 1000
		return int(math.Round(valor * 1000.0))
	default:
		return int(math.Round(valor))
	}
}

// NewTension crea un value object Tension validando que el voltaje sea uno de los
// valores normalizados por NOM: 127, 220, 240, 277, 440, 480 o 600 V.
// Acepta valor como float64 para soportar decimales (ej: 0.48 kV).
// La unidad puede ser "V", "kV" o vacía (default: "V").
// Retorna ErrVoltajeInvalido si el valor no está en la lista.
// Retorna ErrUnidadTensionInvalida si la unidad no es reconocida.
func NewTension(valor float64, unidad string) (Tension, error) {
	// Parsear la unidad (default a V si está vacía)
	unidadParsed, err := ParseUnidadTension(unidad)
	if err != nil {
		return Tension{}, err
	}

	// Normalizar a volts
	volts := normalizarAVolts(valor, unidadParsed)

	// Validar contra valores NOM
	if !voltajesValidos[volts] {
		return Tension{}, fmt.Errorf("%w: %v %s (equivalente a %d V)", ErrVoltajeInvalido, valor, unidadParsed, volts)
	}

	return Tension{valor: volts, unidad: unidadParsed}, nil
}

func (t Tension) Valor() int            { return t.valor }
func (t Tension) Unidad() UnidadTension { return t.unidad }

func (t Tension) EnKilovoltios() float64 {
	return float64(t.valor) / 1000.0
}
