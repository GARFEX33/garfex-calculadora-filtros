// internal/calculos/domain/service/calculo_caida_tension.go
package service

import (
	"errors"
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
)

// ErrDistanciaInvalida is returned when the distance is zero or negative.
var ErrDistanciaInvalida = errors.New("distancia debe ser mayor que cero")

// ErrHilosPorFaseInvalido is returned when HilosPorFase is zero or negative.
var ErrHilosPorFaseInvalido = errors.New("hilos por fase debe ser mayor que cero")

// ErrFactorPotenciaInvalido is returned when FactorPotencia is not in range (0, 1].
var ErrFactorPotenciaInvalido = errors.New("factor de potencia debe estar entre 0 (exclusivo) y 1")

// EntradaCalculoCaidaTension contains the pre-resolved NOM table data needed
// to calculate voltage drop using the IEEE-141 / NOM formula with power factor.
// The application layer is responsible for resolving R, X from Tabla 9 and
// the power factor from the equipment entity.
type EntradaCalculoCaidaTension struct {
	ResistenciaOhmPorKm float64                 // Tabla 9 → res_{material}_{conduit}
	ReactanciaOhmPorKm  float64                 // Tabla 9 → reactancia_al o reactancia_acero
	TipoCanalizacion    entity.TipoCanalizacion // Documented in memoria de cálculo report
	HilosPorFase        int                     // CF ≥ 1 (parallel conductors per phase)
	FactorPotencia      float64                 // cosθ: FA/FR/TR = 1.0 | Carga = explicit FP
}

// CalcularCaidaTension calculates the voltage drop for a three-phase system
// using the IEEE-141 / NOM formula with power factor:
//
//	%Vd = (√3 × Ib × L × (R·cosθ + X·senθ) / (V × N)) × 100
//	VD  = V × (%Vd / 100)
//
// Where cosθ = FactorPotencia, senθ = √(1 - FP²), N = HilosPorFase.
//
// For FP = 1.0 (FiltroActivo, FiltroRechazo, Transformador) the formula
// reduces to: %Vd = (√3 × Ib × L × R / (V × N)) × 100  (reactance has no effect).
func CalcularCaidaTension(
	entrada EntradaCalculoCaidaTension,
	corriente valueobject.Corriente,
	distancia float64,
	tension valueobject.Tension,
	limiteNOM float64,
) (entity.ResultadoCaidaTension, error) {
	if distancia <= 0 {
		return entity.ResultadoCaidaTension{}, fmt.Errorf("CalcularCaidaTension: %w: %.2f", ErrDistanciaInvalida, distancia)
	}
	if entrada.HilosPorFase <= 0 {
		return entity.ResultadoCaidaTension{}, fmt.Errorf("CalcularCaidaTension: %w: %d", ErrHilosPorFaseInvalido, entrada.HilosPorFase)
	}
	if entrada.FactorPotencia <= 0 || entrada.FactorPotencia > 1 {
		return entity.ResultadoCaidaTension{}, fmt.Errorf("CalcularCaidaTension: %w: %.4f", ErrFactorPotenciaInvalido, entrada.FactorPotencia)
	}

	n := float64(entrada.HilosPorFase)

	// Step 1-2: angle components
	cosTheta := entrada.FactorPotencia
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)

	// Step 3-4: effective R and X per parallel conductor
	rEf := entrada.ResistenciaOhmPorKm / n
	xEf := entrada.ReactanciaOhmPorKm / n

	// Step 5: effective impedance term (Ω/km)
	terminoEfectivo := rEf*cosTheta + xEf*sinTheta

	// Step 6: %Vd = (√3 × Ib × L × terminoEfectivo / (V × N)) × 100
	// Note: N (HilosPorFase) is already applied to R and X above via rEf and xEf.
	lKm := distancia / 1000.0
	porcentaje := math.Sqrt(3) * corriente.Valor() * lKm * terminoEfectivo / float64(tension.Valor()) * 100

	// Step 7: VD in volts
	vd := float64(tension.Valor()) * (porcentaje / 100)

	return entity.ResultadoCaidaTension{
		Porcentaje:  porcentaje,
		CaidaVolts:  vd,
		Cumple:      porcentaje <= limiteNOM,
		Impedancia:  terminoEfectivo, // R·cosθ + X·senθ — "effective impedance term"
		Resistencia: rEf,
		Reactancia:  xEf,
	}, nil
}
