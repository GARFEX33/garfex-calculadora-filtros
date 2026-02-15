// internal/shared/kernel/valueobject/temperatura.go
package valueobject

// Temperatura represents the temperature rating in Celsius (60, 75, or 90).
type Temperatura int

const (
	Temp60 Temperatura = 60
	Temp75 Temperatura = 75
	Temp90 Temperatura = 90
)

// Valor returns the temperature value in Celsius.
func (t Temperatura) Valor() int {
	return int(t)
}
