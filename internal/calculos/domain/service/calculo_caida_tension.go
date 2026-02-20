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

// EntradaCalculoCaidaTension contains the pre-resolved NOM table data needed
// to calculate voltage drop using the NOM simplified formula.
// The application layer is responsible for resolving R, X from Tabla 9.
type EntradaCalculoCaidaTension struct {
	ResistenciaOhmPorKm float64                 // Tabla 9 → res_{material}_{conduit}
	ReactanciaOhmPorKm  float64                 // Tabla 9 → reactancia_al o reactancia_acero
	TipoCanalizacion    entity.TipoCanalizacion // Documented in memoria de cálculo report
	HilosPorFase        int                     // N ≥ 1 (parallel conductors per phase)
	SistemaElectrico    entity.SistemaElectrico // For determining voltage drop factor
}

// CalcularCaidaTension calculates the voltage drop using the NOM simplified formula:
//
// Sistema Monofásico 1F-2H:
//
//	e = 2·I·Z·L        →  %e = (e / Vfn) × 100
//
// Sistema Bifásico 2F-3H:
//
//	e = I·Z·L          →  %e = (e / Vfn) × 100
//
// Sistema Trifásico Delta 3F-3H:
//
//	e = √3·I·Z·L       →  %e = (e / Vff) × 100
//
// Sistema Trifásico Estrella 3F-4H:
//
//	e = I·Z·L          →  %e = (e / Vfn) × 100
//
// Where:
//
//	I = Corriente nominal en Amperes (sin factor de corrección)
//	Z = Impedancia de Tabla 9 = √(R² + X²) en Ω/km
//	L = Longitud del alimentador en km
//	N = HilosPorFase (conductores en paralelo)
//	Vfn = Voltaje entre fase y neutro
//	Vff = Voltaje entre fases
//
// Note: La corriente I NO lleva factor de potencia en la fórmula NOM.
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

	n := float64(entrada.HilosPorFase)

	// Step 1: Calculate effective R and X per parallel conductor
	rEf := entrada.ResistenciaOhmPorKm / n
	xEf := entrada.ReactanciaOhmPorKm / n

	// Step 2: Calculate impedance Z = √(R² + X²) from Tabla 9
	impedancia := math.Sqrt(rEf*rEf + xEf*xEf)

	// Step 3: Determine voltage drop factor based on electrical system per NOM
	var factorSistema float64
	switch entrada.SistemaElectrico {
	case entity.SistemaElectricoMonofasico:
		factorSistema = 2.0 // Monofásico 1F-2H
	case entity.SistemaElectricoBifasico:
		factorSistema = 1.0 // Bifásico 2F-3H
	case entity.SistemaElectricoDelta:
		factorSistema = math.Sqrt(3) // Trifásico 3F-3H (Delta)
	case entity.SistemaElectricoEstrella:
		factorSistema = 1.0 // Trifásico 3F-4H (Estrella)
	default:
		return entity.ResultadoCaidaTension{}, fmt.Errorf("sistema eléctrico inválido: %v", entrada.SistemaElectrico)
	}

	// Step 4: Calculate voltage drop: e = factor × I × Z × L
	lKm := distancia / 1000.0
	caida := factorSistema * corriente.Valor() * impedancia * lKm

	// Step 5: Calculate percentage: %e = (e / V) × 100
	porcentaje := (caida / float64(tension.Valor())) * 100

	return entity.ResultadoCaidaTension{
		Porcentaje:  porcentaje,
		CaidaVolts:  caida,
		Cumple:      porcentaje <= limiteNOM,
		Impedancia:  impedancia, // Z = √(R² + X²) from Tabla 9
		Resistencia: rEf,
		Reactancia:  xEf,
	}, nil
}
