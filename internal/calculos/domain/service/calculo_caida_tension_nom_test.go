// internal/calculos/domain/service/calculo_caida_tension_nom_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalcularCaidaTension_VoltajesNOMReales verifica la fórmula NOM con voltajes reales.
//
// Voltajes NOM comunes en México:
//   - Vfn = 127V (fase-neutro)
//   - Vff = 220V (fase-fase)
//   - Relación: Vff = √3 × Vfn → 220 = 1.732 × 127
//
// Fórmula NOM (IEEE-141 — impedancia efectiva):
//
//	Zef = R·cosθ + X·senθ,   senθ = √(1 - cos²θ)
//	e   = factor × I × Zef × L
//	%e  = (e / V_referencia) × 100
//
// Datos base: 2 AWG Cu, tubería PVC, 70A, 30m, FP=0.9
// R = 0.62 Ω/km, X = 0.148 Ω/km (Tabla 9 NOM)
// cosθ = 0.9, senθ = √(1-0.81) = 0.43589
// Zef = 0.62×0.9 + 0.148×0.43589 = 0.558 + 0.06451 = 0.62251 Ω/km
// L   = 0.030 km
func TestCalcularCaidaTension_VoltajesNOMReales(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(70.0)

	tests := []struct {
		name                      string
		sistema                   entity.SistemaElectrico
		tipoVoltaje               entity.TipoVoltaje
		voltajeIngresado          int
		voltajeReferenciaEsperado float64
		factor                    string
		expectedPorcentaje        float64
		expectedCaida             float64
	}{
		{
			name:                      "MONOFASICO - usuario ingresa Vfn 127V",
			sistema:                   entity.SistemaElectricoMonofasico,
			tipoVoltaje:               entity.TipoVoltajeFaseNeutro,
			voltajeIngresado:          127,
			voltajeReferenciaEsperado: 127.0,
			factor:                    "2",
			// e = 2 × 70 × 0.62251 × 0.030 = 2.6145 V
			// %e = 2.6145 / 127 × 100 = 2.059%
			expectedPorcentaje: 2.059,
			expectedCaida:      2.6145,
		},
		{
			name:                      "MONOFASICO - usuario ingresa Vff 220V (se convierte a Vfn 127V)",
			sistema:                   entity.SistemaElectricoMonofasico,
			tipoVoltaje:               entity.TipoVoltajeFaseFase,
			voltajeIngresado:          220,
			voltajeReferenciaEsperado: 127.0, // 220 / √3
			factor:                    "2",
			// e = 2 × 70 × 0.62251 × 0.030 = 2.6145 V
			// %e = 2.6145 / 127 × 100 = 2.059%
			expectedPorcentaje: 2.059,
			expectedCaida:      2.6145,
		},
		{
			name:                      "BIFASICO - usuario ingresa Vfn 127V",
			sistema:                   entity.SistemaElectricoBifasico,
			tipoVoltaje:               entity.TipoVoltajeFaseNeutro,
			voltajeIngresado:          127,
			voltajeReferenciaEsperado: 127.0,
			factor:                    "1",
			// e = 1 × 70 × 0.62251 × 0.030 = 1.3073 V
			// %e = 1.3073 / 127 × 100 = 1.029%
			expectedPorcentaje: 1.029,
			expectedCaida:      1.3073,
		},
		{
			name:                      "DELTA - usuario ingresa Vff 220V",
			sistema:                   entity.SistemaElectricoDelta,
			tipoVoltaje:               entity.TipoVoltajeFaseFase,
			voltajeIngresado:          220,
			voltajeReferenciaEsperado: 220.0,
			factor:                    "√3",
			// e = √3 × 70 × 0.62251 × 0.030 = 2.2641 V
			// %e = 2.2641 / 220 × 100 = 1.029%
			expectedPorcentaje: 1.029,
			expectedCaida:      2.2641,
		},
		{
			name:                      "DELTA - usuario ingresa Vfn 127V (se convierte a Vff 220V)",
			sistema:                   entity.SistemaElectricoDelta,
			tipoVoltaje:               entity.TipoVoltajeFaseNeutro,
			voltajeIngresado:          127,
			voltajeReferenciaEsperado: 220.0, // 127 × √3
			factor:                    "√3",
			// e = √3 × 70 × 0.62251 × 0.030 = 2.2641 V
			// %e = 2.2641 / 220 × 100 = 1.029%
			expectedPorcentaje: 1.029,
			expectedCaida:      2.2641,
		},
		{
			name:                      "ESTRELLA - usuario ingresa Vfn 127V",
			sistema:                   entity.SistemaElectricoEstrella,
			tipoVoltaje:               entity.TipoVoltajeFaseNeutro,
			voltajeIngresado:          127,
			voltajeReferenciaEsperado: 220.0, // se convierte a Vff
			factor:                    "√3",
			// e = √3 × 70 × 0.62251 × 0.030 = 2.264 V
			// %e = 2.264 / 220 × 100 = 1.029%
			expectedPorcentaje: 1.029,
			expectedCaida:      2.264,
		},
		{
			name:                      "ESTRELLA - usuario ingresa Vff 220V",
			sistema:                   entity.SistemaElectricoEstrella,
			tipoVoltaje:               entity.TipoVoltajeFaseFase,
			voltajeIngresado:          220,
			voltajeReferenciaEsperado: 220.0,
			factor:                    "√3",
			// e = √3 × 70 × 0.62251 × 0.030 = 2.264 V
			// %e = 2.264 / 220 × 100 = 1.029%
			expectedPorcentaje: 1.029,
			expectedCaida:      2.264,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entrada := service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.148,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    tt.sistema,
				TipoVoltaje:         tt.tipoVoltaje,
				HilosPorFase:        1,
				FactorPotencia:      0.9,
			}

			tension, err := valueobject.NewTension(float64(tt.voltajeIngresado), "V")
			require.NoError(t, err)

			resultado, err := service.CalcularCaidaTension(
				entrada, corriente, 30.0, tension, 3.0,
			)
			require.NoError(t, err)

			assert.InDelta(t, tt.expectedPorcentaje, resultado.Porcentaje, 0.01, "porcentaje")
			assert.InDelta(t, tt.expectedCaida, resultado.CaidaVolts, 0.01, "caida volts")
			assert.True(t, resultado.Cumple, "debe cumplir límite 3%")
		})
	}
}

// TestCalcularCaidaTension_RelacionesSistemas verifica las relaciones matemáticas
// entre sistemas eléctricos.
//
// IMPORTANTE: Las relaciones dependen del factor y voltaje de referencia:
//   - MONOFASICO: factor=2, Vref=Vfn
//   - BIFASICO: factor=1, Vref=Vfn
//   - ESTRELLA: factor=√3, Vref=Vff
//   - DELTA: factor=√3, Vref=Vff
//
// La fórmula usa impedancia efectiva Zef = R·cosθ + X·senθ.
func TestCalcularCaidaTension_RelacionesSistemas(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(70.0)

	// Caso 1: MONOFASICO usa Vfn (127V)
	entradaMonofasico := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		SistemaElectrico:    entity.SistemaElectricoMonofasico,
		TipoVoltaje:         entity.TipoVoltajeFaseNeutro,
		HilosPorFase:        1,
		FactorPotencia:      0.9,
	}

	tensionMonofasico, _ := valueobject.NewTension(127, "V") // Vfn
	resultadoMonofasico, err := service.CalcularCaidaTension(entradaMonofasico, corriente, 30.0, tensionMonofasico, 3.0)
	require.NoError(t, err)

	// MONOFASICO: factor=2, Vref=127 → %e = 2.059%
	assert.InDelta(t, 2.059, resultadoMonofasico.Porcentaje, 0.01, "MONOFASICO")

	// Caso 2: BIFASICO usa Vfn (127V)
	entradaBifasico := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		SistemaElectrico:    entity.SistemaElectricoBifasico,
		TipoVoltaje:         entity.TipoVoltajeFaseNeutro,
		HilosPorFase:        1,
		FactorPotencia:      0.9,
	}

	tensionBifasico, _ := valueobject.NewTension(127, "V") // Vfn
	resultadoBifasico, err := service.CalcularCaidaTension(entradaBifasico, corriente, 30.0, tensionBifasico, 3.0)
	require.NoError(t, err)

	// BIFASICO: factor=1, Vref=127 → %e = 1.029%
	assert.InDelta(t, 1.029, resultadoBifasico.Porcentaje, 0.01, "BIFASICO")

	// Caso 3: ESTRELLA usa Vff (220V)
	entradaEstrella := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		SistemaElectrico:    entity.SistemaElectricoEstrella,
		TipoVoltaje:         entity.TipoVoltajeFaseFase,
		HilosPorFase:        1,
		FactorPotencia:      0.9,
	}

	tensionEstrella, _ := valueobject.NewTension(220, "V") // Vff
	resultadoEstrella, err := service.CalcularCaidaTension(entradaEstrella, corriente, 30.0, tensionEstrella, 3.0)
	require.NoError(t, err)

	// ESTRELLA: factor=√3, Vref=220 → %e = 1.029%
	assert.InDelta(t, 1.029, resultadoEstrella.Porcentaje, 0.01, "ESTRELLA")

	// Caso 4: DELTA usa Vff (220V)
	entradaDelta := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		SistemaElectrico:    entity.SistemaElectricoDelta,
		TipoVoltaje:         entity.TipoVoltajeFaseFase,
		HilosPorFase:        1,
		FactorPotencia:      0.9,
	}

	tensionDelta, _ := valueobject.NewTension(220, "V") // Vff
	resultadoDelta, err := service.CalcularCaidaTension(entradaDelta, corriente, 30.0, tensionDelta, 3.0)
	require.NoError(t, err)

	// DELTA: factor=√3, Vref=220 → %e = 1.029%
	assert.InDelta(t, 1.029, resultadoDelta.Porcentaje, 0.01, "DELTA")

	// Verificar relaciones:
	// - BIFASICO = MONOFASICO / 2 (mismo Vref, factor 1 vs 2)
	assert.InDelta(t, resultadoMonofasico.Porcentaje/2, resultadoBifasico.Porcentaje, 0.01, "BIFASICO = MONOFASICO/2")

	// - ESTRELLA = DELTA (mismo factor √3 y mismo Vref=220)
	assert.InDelta(t, resultadoDelta.Porcentaje, resultadoEstrella.Porcentaje, 0.01, "ESTRELLA = DELTA")
}

// TestCalcularCaidaTension_ConHilosPorFase verifica que múltiples hilos dividen R y X.
//
// Con FP=0.9:
// R = 0.62, X = 0.148 (valores por conductor, NO dividos por N)
// Zef = 0.62×0.9 + 0.148×0.43589 = 0.558 + 0.0645 = 0.6225 Ω/km
func TestCalcularCaidaTension_ConHilosPorFase(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(70.0)
	tension, _ := valueobject.NewTension(127, "V") // Vfn

	entrada2Hilos := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		SistemaElectrico:    entity.SistemaElectricoMonofasico,
		TipoVoltaje:         entity.TipoVoltajeFaseNeutro,
		HilosPorFase:        2,
		FactorPotencia:      0.9,
	}

	resultado, err := service.CalcularCaidaTension(entrada2Hilos, corriente, 30.0, tension, 3.0)
	require.NoError(t, err)

	// R y X se mantienen como valores por conductor (no se dividen por N)
	assert.InDelta(t, 0.62, resultado.Resistencia, 0.001, "R original por conductor")
	assert.InDelta(t, 0.148, resultado.Reactancia, 0.001, "X original por conductor")

	// Zef = R×cosθ + X×senθ = 0.62×0.9 + 0.148×0.43589 = 0.6225 Ω/km
	assert.InDelta(t, 0.6225, resultado.Impedancia, 0.001, "Zef por conductor")

	// e = 2 × (70/2) × 0.6225 × 0.030 = 1.3073 V
	// (factor 2, corriente I/N = 35A, Zef por conductor, L=0.03km)
	assert.InDelta(t, 1.3073, resultado.CaidaVolts, 0.01, "caída con 2 hilos")

	// %e = 1.3073 / 127 × 100 = 1.029%
	assert.InDelta(t, 1.029, resultado.Porcentaje, 0.01, "porcentaje con 2 hilos")
}

// TestCalcularCaidaTension_VoltajesUSA verifica con voltajes 480V/277V (USA industrial).
func TestCalcularCaidaTension_VoltajesUSA(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(120.0)

	tests := []struct {
		name             string
		sistema          entity.SistemaElectrico
		tipoVoltaje      entity.TipoVoltaje
		voltajeIngresado int
		expectedCumple   bool
	}{
		{
			name:             "ESTRELLA 480V-3F (usuario ingresa Vff, usa Vff directo)",
			sistema:          entity.SistemaElectricoEstrella,
			tipoVoltaje:      entity.TipoVoltajeFaseFase,
			voltajeIngresado: 480,
			expectedCumple:   true, // caída pequeña con 480V sistema
		},
		{
			name:             "DELTA 480V-3F (usuario ingresa Vff 480V)",
			sistema:          entity.SistemaElectricoDelta,
			tipoVoltaje:      entity.TipoVoltajeFaseFase,
			voltajeIngresado: 480,
			expectedCumple:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entrada := service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.148,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    tt.sistema,
				TipoVoltaje:         tt.tipoVoltaje,
				HilosPorFase:        1,
				FactorPotencia:      0.9,
			}

			tension, err := valueobject.NewTension(float64(tt.voltajeIngresado), "V")
			require.NoError(t, err)

			resultado, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCumple, resultado.Cumple)
		})
	}
}

// TestCalcularCaidaTension_FactoresPotencia verifica el comportamiento con distintos FP.
//
// Con R=0.62, X=0.148:
//   - FP=1.0 (resistiva pura): Zef = 0.62×1 + 0.148×0 = 0.620 Ω/km
//   - FP=0.9:                  Zef = 0.62×0.9 + 0.148×0.43589 = 0.62251 Ω/km
//   - FP=0.8:                  Zef = 0.62×0.8 + 0.148×0.6 = 0.496 + 0.08880 = 0.58480 Ω/km
func TestCalcularCaidaTension_FactoresPotencia(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(70.0)
	tension, _ := valueobject.NewTension(127, "V")

	tests := []struct {
		name           string
		factorPotencia float64
		expectedZef    float64
	}{
		{
			name:           "FP=1.0 — carga resistiva pura, Zef = R",
			factorPotencia: 1.0,
			// Zef = 0.62×1 + 0.148×0 = 0.620
			expectedZef: 0.620,
		},
		{
			name:           "FP=0.9 — carga mixta típica",
			factorPotencia: 0.9,
			// Zef = 0.62×0.9 + 0.148×0.43589 = 0.558 + 0.06451 = 0.62251
			expectedZef: 0.62251,
		},
		{
			name:           "FP=0.8 — carga predominantemente inductiva",
			factorPotencia: 0.8,
			// senθ = √(1-0.64) = √0.36 = 0.6
			// Zef = 0.62×0.8 + 0.148×0.6 = 0.496 + 0.08880 = 0.58480
			expectedZef: 0.58480,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entrada := service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.148,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    entity.SistemaElectricoMonofasico,
				TipoVoltaje:         entity.TipoVoltajeFaseNeutro,
				HilosPorFase:        1,
				FactorPotencia:      tt.factorPotencia,
			}

			resultado, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
			require.NoError(t, err)
			assert.InDelta(t, tt.expectedZef, resultado.Impedancia, 0.001, "Zef efectiva")
		})
	}
}

// TestCalcularCaidaTension_ErroresValidacion verifica errores de validación.
func TestCalcularCaidaTension_ErroresValidacion(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(70.0)
	tension, _ := valueobject.NewTension(127, "V")

	t.Run("distancia cero", func(t *testing.T) {
		entrada := service.EntradaCalculoCaidaTension{
			ResistenciaOhmPorKm: 0.62,
			ReactanciaOhmPorKm:  0.148,
			TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
			SistemaElectrico:    entity.SistemaElectricoMonofasico,
			TipoVoltaje:         entity.TipoVoltajeFaseNeutro,
			HilosPorFase:        1,
			FactorPotencia:      0.9,
		}
		_, err := service.CalcularCaidaTension(entrada, corriente, 0, tension, 3.0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrDistanciaInvalida)
	})

	t.Run("hilos por fase cero", func(t *testing.T) {
		entrada := service.EntradaCalculoCaidaTension{
			ResistenciaOhmPorKm: 0.62,
			ReactanciaOhmPorKm:  0.148,
			TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
			SistemaElectrico:    entity.SistemaElectricoMonofasico,
			TipoVoltaje:         entity.TipoVoltajeFaseNeutro,
			HilosPorFase:        0,
			FactorPotencia:      0.9,
		}
		_, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrHilosPorFaseInvalido)
	})
}
