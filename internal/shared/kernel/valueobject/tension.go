// internal/shared/kernel/valueobject/tension.go
package valueobject

import (
	"encoding/json"
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

// MarshalJSON serializa Tension a JSON.
func (t Tension) MarshalJSON() ([]byte, error) {
	type alias Tension
	return json.Marshal(&struct {
		Valor  int    `json:"valor"`
		Unidad string `json:"unidad"`
	}{
		Valor:  t.valor,
		Unidad: t.unidad,
	})
}

func (t Tension) EnKilovoltios() float64 {
	return float64(t.valor) / 1000.0
}
