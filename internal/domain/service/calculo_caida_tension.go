package service

import (
	"errors"
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ErrDistanciaInvalida is returned when the distance is zero or negative.
var ErrDistanciaInvalida = errors.New("distancia debe ser mayor que cero")

// resistividad maps conductor material to its resistivity in Ω·mm²/m at ~75°C.
var resistividad = map[string]float64{
	"Cu": 0.01724,
	"Al": 0.02826,
}

// CalcularCaidaTension calculates the voltage drop percentage for a three-phase system.
//
// Formula: VD% = (√3 × ρ × L × I) / (S × V) × 100
//
// Where:
//   - ρ = resistivity (Cu: 0.01724, Al: 0.02826 Ω·mm²/m at 75°C)
//   - L = one-way distance in meters
//   - I = current in amperes
//   - S = conductor cross-section in mm²
//   - V = system voltage in volts
//
// Returns the voltage drop percentage, whether it meets the NOM limit, and any error.
func CalcularCaidaTension(
	conductor valueobject.Conductor,
	corriente valueobject.Corriente,
	distancia float64,
	tension valueobject.Tension,
	limiteNOM float64,
) (porcentaje float64, cumple bool, err error) {
	if distancia <= 0 {
		return 0, false, fmt.Errorf("%w: %.2f", ErrDistanciaInvalida, distancia)
	}

	rho, ok := resistividad[conductor.Material()]
	if !ok {
		return 0, false, fmt.Errorf("material desconocido para resistividad: %s", conductor.Material())
	}

	vd := (math.Sqrt(3) * rho * distancia * corriente.Valor()) / conductor.SeccionMM2()
	porcentaje = (vd / float64(tension.Valor())) * 100

	return porcentaje, porcentaje <= limiteNOM, nil
}
