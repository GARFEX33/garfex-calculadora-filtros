// internal/shared/kernel/valueobject/corriente.go
package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrCorrienteInvalida = errors.New("corriente debe ser mayor que cero")

// Corriente represents an electrical current value in Amperes. Immutable.
type Corriente struct {
	valor  float64
	unidad string
}

func NewCorriente(valor float64) (Corriente, error) {
	if valor <= 0 {
		return Corriente{}, fmt.Errorf("%w: %.4f", ErrCorrienteInvalida, valor)
	}
	return Corriente{valor: valor, unidad: "A"}, nil
}

func (c Corriente) Valor() float64 { return c.valor }
func (c Corriente) Unidad() string { return c.unidad }

// MarshalJSON serializa Corriente a JSON.
func (c Corriente) MarshalJSON() ([]byte, error) {
	type alias Corriente
	return json.Marshal(&struct {
		Valor  float64 `json:"valor"`
		Unidad string  `json:"unidad"`
	}{
		Valor:  c.valor,
		Unidad: c.unidad,
	})
}

func (c Corriente) Multiplicar(factor float64) (Corriente, error) {
	return NewCorriente(c.valor * factor)
}

func (c Corriente) Dividir(divisor int) (Corriente, error) {
	if divisor == 0 {
		return Corriente{}, fmt.Errorf("dividir corriente: divisor es cero")
	}
	return NewCorriente(c.valor / float64(divisor))
}
