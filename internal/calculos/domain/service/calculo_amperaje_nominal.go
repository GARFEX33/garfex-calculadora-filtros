// internal/calculos/domain/service/calculo_amperaje_nominal.go
package service

import (
	"errors"
	"math"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ErrPotenciaInvalida is returned when potenciaWatts is zero or negative.
var ErrPotenciaInvalida = errors.New("potencia debe ser mayor que cero")

// ErrTensionInvalida is returned when the tension value is zero or negative.
var ErrTensionInvalida = errors.New("tensión no válida")

// ErrFactorPotenciaInvalido is returned when factor de potencia is out of range (0, 1].
var ErrFactorPotenciaInvalido = errors.New("factor de potencia debe estar entre 0 y 1")

// CalcularAmperajeNominalCircuito calcula el amperaje nominal de un circuito
// eléctrico según fórmulas NOM-001-SEDE.
//
// Parámetros:
//   - potenciaWatts: potencia activa en Watts
//   - tension: tensión del circuito (V)
//   - sistema: sistema eléctrico canónico del dominio
//   - factorPotencia: factor de potencia (0 a 1]
//
// Fórmulas según NOM:
//   - Monofásico/Bifásico: I = W / (V × FP)
//   - Delta/Estrella (trifásico): I = W / (V × √3 × FP)
func CalcularAmperajeNominalCircuito(
	potenciaWatts float64,
	tension valueobject.Tension,
	sistema entity.SistemaElectrico,
	factorPotencia float64,
) (valueobject.Corriente, error) {
	if potenciaWatts <= 0 {
		return valueobject.Corriente{}, ErrPotenciaInvalida
	}
	if factorPotencia <= 0 || factorPotencia > 1 {
		return valueobject.Corriente{}, ErrFactorPotenciaInvalido
	}

	tensionV := float64(tension.Valor())
	if tensionV <= 0 {
		return valueobject.Corriente{}, ErrTensionInvalida
	}

	var amperaje float64
	switch sistema {
	case entity.SistemaElectricoMonofasico, entity.SistemaElectricoBifasico:
		// I = W / (V × FP)
		amperaje = potenciaWatts / (tensionV * factorPotencia)
	default:
		// Delta, Estrella → trifásico: I = W / (V × √3 × FP)
		amperaje = potenciaWatts / (tensionV * math.Sqrt(3) * factorPotencia)
	}

	return valueobject.NewCorriente(amperaje)
}
