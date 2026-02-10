// internal/domain/valueobject/tension.go
package valueobject

import (
	"errors"
	"fmt"
)

var ErrVoltajeInvalido = errors.New("voltaje no válido según normativa NOM")

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
type Tension struct {
	valor  int
	unidad string
}

func NewTension(valor int) (Tension, error) {
	if !voltajesValidos[valor] {
		return Tension{}, fmt.Errorf("%w: %d", ErrVoltajeInvalido, valor)
	}
	return Tension{valor: valor, unidad: "V"}, nil
}

func (t Tension) Valor() int     { return t.valor }
func (t Tension) Unidad() string { return t.unidad }

func (t Tension) EnKilovoltios() float64 {
	return float64(t.valor) / 1000.0
}
