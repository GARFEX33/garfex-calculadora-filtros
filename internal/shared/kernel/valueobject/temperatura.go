// internal/shared/kernel/valueobject/temperatura.go
package valueobject

import (
	"errors"
	"fmt"
)

// ErrTemperaturaInvalida se retorna cuando la temperatura no corresponde a un
// rating normalizado por NOM (60 °C, 75 °C o 90 °C).
var ErrTemperaturaInvalida = errors.New("temperatura no válida; valores permitidos: 60, 75, 90 °C")

// Temperatura represents the temperature rating in Celsius (60, 75, or 90).
type Temperatura int

const (
	Temp60 Temperatura = 60
	Temp75 Temperatura = 75
	Temp90 Temperatura = 90
)

// ValidarTemperatura verifica que t sea un rating de temperatura reconocido
// por la normativa NOM (60 °C, 75 °C o 90 °C).
// Retorna nil si es válida, ErrTemperaturaInvalida en caso contrario.
func ValidarTemperatura(t Temperatura) error {
	switch t {
	case Temp60, Temp75, Temp90:
		return nil
	default:
		return fmt.Errorf("%w: %d °C", ErrTemperaturaInvalida, int(t))
	}
}

// Valor returns the temperature value in Celsius.
func (t Temperatura) Valor() int {
	return int(t)
}
