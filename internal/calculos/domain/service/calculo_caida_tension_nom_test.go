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
// Fórmula NOM: e = factor × I × Z × L
//
//	Z = √(R² + X²)
//	%e = (e / V_referencia) × 100
//
// Datos base: 2 AWG Cu, tubería PVC, 70A, 30m
// R = 0.62 Ω/km, X = 0.148 Ω/km (de tabla 9)
// Z = √(0.62² + 0.148²) = 0.6374 Ω/km
// L = 0.030 km
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
			// e = 2 × 70 × 0.6374 × 0.030 = 2.677 V
			// %e = 2.677 / 127 × 100 = 2.108%
			expectedPorcentaje: 2.108,
			expectedCaida:      2.677,
		},
		{
			name:                      "MONOFASICO - usuario ingresa Vff 220V (se convierte a Vfn 127V)",
			sistema:                   entity.SistemaElectricoMonofasico,
			tipoVoltaje:               entity.TipoVoltajeFaseFase,
			voltajeIngresado:          220,
			voltajeReferenciaEsperado: 127.0, // 220 / √3
			factor:                    "2",
			// e = 2 × 70 × 0.6374 × 0.030 = 2.677 V
			// %e = 2.677 / 127 × 100 = 2.108%
			expectedPorcentaje: 2.108,
			expectedCaida:      2.677,
		},
		{
			name:                      "BIFASICO - usuario ingresa Vfn 127V",
			sistema:                   entity.SistemaElectricoBifasico,
			tipoVoltaje:               entity.TipoVoltajeFaseNeutro,
			voltajeIngresado:          127,
			voltajeReferenciaEsperado: 127.0,
			factor:                    "1",
			// e = 1 × 70 × 0.6374 × 0.030 = 1.338 V
			// %e = 1.338 / 127 × 100 = 1.054%
			expectedPorcentaje: 1.054,
			expectedCaida:      1.338,
		},
		{
			name:                      "DELTA - usuario ingresa Vff 220V",
			sistema:                   entity.SistemaElectricoDelta,
			tipoVoltaje:               entity.TipoVoltajeFaseFase,
			voltajeIngresado:          220,
			voltajeReferenciaEsperado: 220.0,
			factor:                    "√3",
			// e = √3 × 70 × 0.6374 × 0.030 = 2.318 V
			// %e = 2.318 / 220 × 100 = 1.054%
			expectedPorcentaje: 1.054,
			expectedCaida:      2.318,
		},
		{
			name:                      "DELTA - usuario ingresa Vfn 127V (se convierte a Vff 220V)",
			sistema:                   entity.SistemaElectricoDelta,
			tipoVoltaje:               entity.TipoVoltajeFaseNeutro,
			voltajeIngresado:          127,
			voltajeReferenciaEsperado: 220.0, // 127 × √3
			factor:                    "√3",
			// e = √3 × 70 × 0.6374 × 0.030 = 2.318 V
			// %e = 2.318 / 220 × 100 = 1.054%
			expectedPorcentaje: 1.054,
			expectedCaida:      2.318,
		},
		{
			name:                      "ESTRELLA - usuario ingresa Vfn 127V",
			sistema:                   entity.SistemaElectricoEstrella,
			tipoVoltaje:               entity.TipoVoltajeFaseNeutro,
			voltajeIngresado:          127,
			voltajeReferenciaEsperado: 127.0,
			factor:                    "1",
			// e = 1 × 70 × 0.6374 × 0.030 = 1.338 V
			// %e = 1.338 / 127 × 100 = 1.054%
			expectedPorcentaje: 1.054,
			expectedCaida:      1.338,
		},
		{
			name:                      "ESTRELLA - usuario ingresa Vff 220V (se convierte a Vfn 127V)",
			sistema:                   entity.SistemaElectricoEstrella,
			tipoVoltaje:               entity.TipoVoltajeFaseFase,
			voltajeIngresado:          220,
			voltajeReferenciaEsperado: 127.0, // 220 / √3
			factor:                    "1",
			// e = 1 × 70 × 0.6374 × 0.030 = 1.338 V
			// %e = 1.338 / 127 × 100 = 1.054%
			expectedPorcentaje: 1.054,
			expectedCaida:      1.338,
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
			}

			tension, err := valueobject.NewTension(tt.voltajeIngresado)
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
// entre sistemas cuando todos usan el mismo voltaje de referencia.
//
// IMPORTANTE: Estas relaciones solo son válidas cuando se compara con el mismo
// voltaje de REFERENCIA (Vfn para MONOFASICO/BIFASICO/ESTRELLA, Vff para DELTA).
func TestCalcularCaidaTension_RelacionesSistemas(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(70.0)

	// Caso 1: Sistemas que usan Vfn (127V)
	sistemasVfn := []struct {
		nombre  string
		sistema entity.SistemaElectrico
		factor  float64
	}{
		{"MONOFASICO", entity.SistemaElectricoMonofasico, 2.0},
		{"BIFASICO", entity.SistemaElectricoBifasico, 1.0},
		{"ESTRELLA", entity.SistemaElectricoEstrella, 1.0},
	}

	resultadosVfn := make(map[string]float64)
	for _, s := range sistemasVfn {
		entrada := service.EntradaCalculoCaidaTension{
			ResistenciaOhmPorKm: 0.62,
			ReactanciaOhmPorKm:  0.148,
			TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
			SistemaElectrico:    s.sistema,
			TipoVoltaje:         entity.TipoVoltajeFaseNeutro,
			HilosPorFase:        1,
		}

		tension, _ := valueobject.NewTension(127) // Vfn
		resultado, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
		require.NoError(t, err)
		resultadosVfn[s.nombre] = resultado.Porcentaje
	}

	// Verificar relaciones
	assert.InDelta(t, resultadosVfn["MONOFASICO"], 2*resultadosVfn["BIFASICO"], 0.01, "MONOFASICO = 2 × BIFASICO")
	assert.InDelta(t, resultadosVfn["ESTRELLA"], resultadosVfn["BIFASICO"], 0.01, "ESTRELLA = BIFASICO")

	// Caso 2: Sistema DELTA usa Vff (220V)
	entradaDelta := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		SistemaElectrico:    entity.SistemaElectricoDelta,
		TipoVoltaje:         entity.TipoVoltajeFaseFase,
		HilosPorFase:        1,
	}

	tensionDelta, _ := valueobject.NewTension(220) // Vff
	resultadoDelta, err := service.CalcularCaidaTension(entradaDelta, corriente, 30.0, tensionDelta, 3.0)
	require.NoError(t, err)

	// DELTA con Vff debe dar el mismo porcentaje que BIFASICO/ESTRELLA con Vfn
	// porque ambos usan factor 1 (DELTA=√3) con sus respectivas referencias
	assert.InDelta(t, resultadosVfn["BIFASICO"], resultadoDelta.Porcentaje, 0.01, "DELTA (con Vff) = BIFASICO (con Vfn)")
}

// TestCalcularCaidaTension_ConHilosPorFase verifica que múltiples hilos dividen R y X.
func TestCalcularCaidaTension_ConHilosPorFase(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(70.0)
	tension, _ := valueobject.NewTension(127) // Vfn

	entrada2Hilos := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		SistemaElectrico:    entity.SistemaElectricoMonofasico,
		TipoVoltaje:         entity.TipoVoltajeFaseNeutro,
		HilosPorFase:        2,
	}

	resultado, err := service.CalcularCaidaTension(entrada2Hilos, corriente, 30.0, tension, 3.0)
	require.NoError(t, err)

	// R_ef = 0.62 / 2 = 0.31
	assert.InDelta(t, 0.31, resultado.Resistencia, 0.001, "R dividido por 2")

	// X_ef = 0.148 / 2 = 0.074
	assert.InDelta(t, 0.074, resultado.Reactancia, 0.001, "X dividido por 2")

	// Z = √(0.31² + 0.074²) = 0.319
	assert.InDelta(t, 0.319, resultado.Impedancia, 0.001, "Z con 2 hilos")

	// Caída debe ser la mitad que con 1 hilo
	// e = 2 × 70 × 0.319 × 0.030 = 1.338 V
	assert.InDelta(t, 1.338, resultado.CaidaVolts, 0.01, "caída con 2 hilos")

	// %e = 1.338 / 127 × 100 = 1.054%
	assert.InDelta(t, 1.054, resultado.Porcentaje, 0.01, "porcentaje con 2 hilos")
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
			name:             "ESTRELLA 480V-3F (usuario ingresa Vff, se convierte a Vfn 277V)",
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
			}

			tension, err := valueobject.NewTension(tt.voltajeIngresado)
			require.NoError(t, err)

			resultado, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCumple, resultado.Cumple)
		})
	}
}

// TestCalcularCaidaTension_ErroresValidacion verifica errores de validación.
func TestCalcularCaidaTension_ErroresValidacion(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(70.0)
	tension, _ := valueobject.NewTension(127)

	t.Run("distancia cero", func(t *testing.T) {
		entrada := service.EntradaCalculoCaidaTension{
			ResistenciaOhmPorKm: 0.62,
			ReactanciaOhmPorKm:  0.148,
			TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
			SistemaElectrico:    entity.SistemaElectricoMonofasico,
			TipoVoltaje:         entity.TipoVoltajeFaseNeutro,
			HilosPorFase:        1,
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
		}
		_, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrHilosPorFaseInvalido)
	})
}
