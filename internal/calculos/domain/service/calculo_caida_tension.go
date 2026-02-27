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

// calcularVoltajeReferencia convierte el voltaje ingresado al voltaje de referencia
// requerido por el sistema eléctrico según NOM-001-SEDE-2012.
//
// Reglas de conversión:
//   - MONOFASICO, BIFASICO → requieren Vfn (fase-neutro)
//   - ESTRELLA, DELTA → requieren Vff (fase-fase)
//
// Si el voltaje ingresado no coincide con el requerido, se convierte usando:
//   - Vfn = Vff / √3
//   - Vff = Vfn × √3
func calcularVoltajeReferencia(
	voltajeIngresado float64,
	tipoVoltaje entity.TipoVoltaje,
	sistema entity.SistemaElectrico,
) float64 {
	// Determinar qué tipo de voltaje requiere el sistema
	var requiereVfn bool
	switch sistema {
	case entity.SistemaElectricoMonofasico, entity.SistemaElectricoBifasico:
		requiereVfn = true
	case entity.SistemaElectricoDelta, entity.SistemaElectricoEstrella:
		requiereVfn = false // requiere Vff
	default:
		// Sistema inválido, retornar voltaje sin conversión
		return voltajeIngresado
	}

	// Convertir si es necesario
	if requiereVfn && tipoVoltaje.EsFaseFase() {
		// Usuario ingresó Vff, pero necesitamos Vfn
		return voltajeIngresado / math.Sqrt(3)
	}

	if !requiereVfn && tipoVoltaje.EsFaseNeutro() {
		// Usuario ingresó Vfn, pero necesitamos Vff
		return voltajeIngresado * math.Sqrt(3)
	}

	// Ya es el tipo correcto, no convertir
	return voltajeIngresado
}

// EntradaCalculoCaidaTension contains the pre-resolved NOM table data needed
// to calculate voltage drop using the NOM simplified formula.
// The application layer is responsible for resolving R, X from Tabla 9.
type EntradaCalculoCaidaTension struct {
	ResistenciaOhmPorKm float64                 // Tabla 9 → res_{material}_{conduit}
	ReactanciaOhmPorKm  float64                 // Tabla 9 → reactancia_al o reactancia_acero
	TipoCanalizacion    entity.TipoCanalizacion // Documented in memoria de cálculo report
	HilosPorFase        int                     // N ≥ 1 (parallel conductors per phase)
	SistemaElectrico    entity.SistemaElectrico // For determining voltage drop factor
	TipoVoltaje         entity.TipoVoltaje      // FASE_NEUTRO or FASE_FASE (for voltage reference)
	FactorPotencia      float64                 // cosθ ∈ (0, 1] — factor de potencia del circuito
}

// CalcularCaidaTension calculates the voltage drop using the NOM / IEEE-141 formula
// with effective impedance:
//
// Formula general: e = factor × (I/N) × L × (R × cosθ + X × sinθ)
//
// Sistema Monofásico 1F-2H (Circuito Monofásico 2 hilos):
//
//	e = 2 × (I/N) × L × Zef
//	%e = (e / Vfn) × 100
//
// Sistema Bifásico 2F-3H (Circuito Monofásico 3 hilos):
//
//	e = (I/N) × L × Zef
//	%e = (e / Vfn) × 100
//
// Sistema Trifásico Delta 3F-3H:
//
//	e = √3 × (I/N) × L × Zef
//	%e = (e / Vff) × 100
//
// Sistema Trifásico Estrella 3F-4H:
//
//	e = √3 × (I/N) × L × Zef
//	%e = (e / Vff) × 100
//
// Where:
//
//	I   = Corriente nominal en Amperes
//	N   = HilosPorFase (conductores en paralelo por fase)
//	Zef = Impedancia efectiva = R·cosθ + X·senθ  en Ω/km  (IEEE-141 / NOM Tabla 9)
//	L   = Longitud del alimentador en km
//	Vfn = Voltaje fase-neutro
//	Vff = Voltaje fase-fase
//	cosθ = FactorPotencia,  senθ = √(1 - cos²θ)
//
// Note: Se usa impedancia EFECTIVA (R·cosθ + X·senθ), NO la magnitud √(R²+X²).
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

	// Step 1: Calculate effective impedance Zef = R·cosθ + X·senθ (per conductor)
	// Note: We use R and X directly, NOT divided by N. The division by N
	// is applied to the current I in Step 4.
	cosTheta := entrada.FactorPotencia
	senTheta := math.Sqrt(1 - cosTheta*cosTheta)
	impedancia := entrada.ResistenciaOhmPorKm*cosTheta + entrada.ReactanciaOhmPorKm*senTheta

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
		factorSistema = math.Sqrt(3) // Trifásico 3F-4H (Estrella)
	default:
		return entity.ResultadoCaidaTension{}, fmt.Errorf("sistema eléctrico inválido: %v", entrada.SistemaElectrico)
	}

	// Step 4: Calculate voltage drop: e = factor × (I/N) × Z × L
	lKm := distancia / 1000.0
	corrientePorHilo := corriente.Valor() / n
	caida := factorSistema * corrientePorHilo * impedancia * lKm

	// Step 5: Determine voltage reference based on electrical system and convert if needed
	// According to NOM-001-SEDE-2012:
	//   - MONOFASICO, BIFASICO → use Vfn; ESTRELLA, DELTA → use Vff (phase-to-neutral)
	//   - DELTA → use Vff (phase-to-phase)
	voltajeReferencia := calcularVoltajeReferencia(
		float64(tension.Valor()),
		entrada.TipoVoltaje,
		entrada.SistemaElectrico,
	)

	// Step 6: Calculate percentage: %e = (e / V_referencia) × 100
	porcentaje := (caida / voltajeReferencia) * 100

	return entity.ResultadoCaidaTension{
		Porcentaje:  porcentaje,
		CaidaVolts:  caida,
		Cumple:      porcentaje <= limiteNOM,
		Impedancia:  impedancia, // Zef = R·cosθ + X·senθ (IEEE-141 efectiva)
		Resistencia: entrada.ResistenciaOhmPorKm,
		Reactancia:  entrada.ReactanciaOhmPorKm,
	}, nil
}
