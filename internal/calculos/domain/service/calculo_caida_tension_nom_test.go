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

// TestCalcularCaidaTension_FormulaNOM verifica la fórmula NOM simplificada.
//
// Fórmula NOM: e = factor × I × Z × L
//
//	Z = √(R² + X²)
//	%e = (e / V) × 100
//
// Datos base: 2 AWG Cu, tubería PVC, 50A, 20m, 220V
// R = 0.62 Ω/km, X = 0.148 Ω/km (de tabla 9)
// Z = √(0.62² + 0.148²) = √(0.3844 + 0.0219) = √0.4063 = 0.6374 Ω/km
// L = 0.020 km
func TestCalcularCaidaTension_FormulaNOM(t *testing.T) {
	// Z = √(0.62² + 0.148²) = 0.6374 Ω/km
	// L = 0.020 km
	// I = 50 A
	// V = 220 V

	tests := []struct {
		name               string
		sistema            entity.SistemaElectrico
		factor             string
		expectedPorcentaje float64
		expectedCaida      float64
	}{
		{
			name:    "Monofásico 1F-2H usa factor 2",
			sistema: entity.SistemaElectricoMonofasico,
			factor:  "2",
			// e = 2 × 50 × 0.6374 × 0.020 = 1.2748 V
			// %e = 1.2748 / 220 × 100 = 0.5794%
			expectedPorcentaje: 0.5794,
			expectedCaida:      1.2748,
		},
		{
			name:    "Bifásico 2F-3H usa factor 1",
			sistema: entity.SistemaElectricoBifasico,
			factor:  "1",
			// e = 1 × 50 × 0.6374 × 0.020 = 0.6374 V
			// %e = 0.6374 / 220 × 100 = 0.2897%
			expectedPorcentaje: 0.2897,
			expectedCaida:      0.6374,
		},
		{
			name:    "Trifásico Delta 3F-3H usa factor √3",
			sistema: entity.SistemaElectricoDelta,
			factor:  "√3",
			// e = 1.732 × 50 × 0.6374 × 0.020 = 1.1040 V
			// %e = 1.1040 / 220 × 100 = 0.5018%
			expectedPorcentaje: 0.5018,
			expectedCaida:      1.1040,
		},
		{
			name:    "Trifásico Estrella 3F-4H usa factor 1",
			sistema: entity.SistemaElectricoEstrella,
			factor:  "1",
			// e = 1 × 50 × 0.6374 × 0.020 = 0.6374 V
			// %e = 0.6374 / 220 × 100 = 0.2897%
			expectedPorcentaje: 0.2897,
			expectedCaida:      0.6374,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entrada := service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,  // Tabla 9: 2 AWG Cu PVC
				ReactanciaOhmPorKm:  0.148, // Tabla 9: 2 AWG reactancia_al
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				HilosPorFase:        1,
				SistemaElectrico:    tt.sistema,
			}

			corriente, err := valueobject.NewCorriente(50)
			require.NoError(t, err)

			tension, err := valueobject.NewTension(220)
			require.NoError(t, err)

			resultado, err := service.CalcularCaidaTension(
				entrada, corriente, 20.0, tension, 3.0,
			)
			require.NoError(t, err)

			assert.InDelta(t, tt.expectedPorcentaje, resultado.Porcentaje, 0.01, "porcentaje")
			assert.InDelta(t, tt.expectedCaida, resultado.CaidaVolts, 0.01, "caida volts")
			assert.True(t, resultado.Cumple, "debe cumplir límite 3%")

			// Verificar que Z se calculó correctamente: √(R² + X²)
			expectedZ := 0.6374
			assert.InDelta(t, expectedZ, resultado.Impedancia, 0.001, "impedancia Z")
		})
	}
}

// TestCalcularCaidaTension_ConHilosPorFase verifica que múltiples hilos dividen R y X.
func TestCalcularCaidaTension_ConHilosPorFase(t *testing.T) {
	// Con 2 hilos por fase:
	// R_ef = 0.62 / 2 = 0.31 Ω/km
	// X_ef = 0.148 / 2 = 0.074 Ω/km
	// Z = √(0.31² + 0.074²) = √(0.0961 + 0.0055) = √0.1016 = 0.3187 Ω/km
	// e = √3 × 50 × 0.3187 × 0.020 = 0.5520 V (Delta)
	// %e = 0.5520 / 220 × 100 = 0.2509%

	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		HilosPorFase:        2, // 2 hilos en paralelo
		SistemaElectrico:    entity.SistemaElectricoDelta,
	}

	corriente, _ := valueobject.NewCorriente(50)
	tension, _ := valueobject.NewTension(220)

	resultado, err := service.CalcularCaidaTension(entrada, corriente, 20.0, tension, 3.0)
	require.NoError(t, err)

	// Z efectiva = 0.3187 Ω/km (mitad de Z con 1 hilo)
	assert.InDelta(t, 0.3187, resultado.Impedancia, 0.001, "impedancia con 2 hilos")
	assert.InDelta(t, 0.2509, resultado.Porcentaje, 0.01, "porcentaje")
	assert.InDelta(t, 0.5520, resultado.CaidaVolts, 0.01, "caida volts")
}

// TestCalcularCaidaTension_ComparacionSistemas verifica las relaciones entre sistemas.
func TestCalcularCaidaTension_ComparacionSistemas(t *testing.T) {
	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		HilosPorFase:        1,
	}

	corriente, _ := valueobject.NewCorriente(100)
	tension, _ := valueobject.NewTension(220)

	// Calcular para cada sistema
	resultados := make(map[string]float64)
	sistemas := map[string]entity.SistemaElectrico{
		"Monofasico": entity.SistemaElectricoMonofasico,
		"Bifasico":   entity.SistemaElectricoBifasico,
		"Delta":      entity.SistemaElectricoDelta,
		"Estrella":   entity.SistemaElectricoEstrella,
	}

	for nombre, sistema := range sistemas {
		entrada.SistemaElectrico = sistema
		resultado, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
		require.NoError(t, err)
		resultados[nombre] = resultado.Porcentaje
	}

	// Verificar relaciones esperadas
	// Monofásico (factor 2) debe ser exactamente el doble que Bifásico (factor 1)
	assert.InDelta(t, resultados["Monofasico"], 2*resultados["Bifasico"], 0.01, "Mono = 2 × Bifásico")

	// Delta (factor √3) debe ser √3 veces Bifásico (factor 1)
	ratioDeltatoBi := resultados["Delta"] / resultados["Bifasico"]
	assert.InDelta(t, 1.732, ratioDeltatoBi, 0.01, "Delta = √3 × Bifásico")

	// Estrella y Bifásico deben ser iguales (ambos factor 1)
	assert.InDelta(t, resultados["Estrella"], resultados["Bifasico"], 0.01, "Estrella = Bifásico")

	// Monofásico debe tener mayor caída que todos
	assert.Greater(t, resultados["Monofasico"], resultados["Delta"], "Mono > Delta")
	assert.Greater(t, resultados["Monofasico"], resultados["Estrella"], "Mono > Estrella")
	assert.Greater(t, resultados["Monofasico"], resultados["Bifasico"], "Mono > Bifásico")
}

// TestCalcularCaidaTension_EjemploReal verifica con valores reales NOM.
func TestCalcularCaidaTension_EjemploReal(t *testing.T) {
	// Ejemplo: Circuito trifásico Delta, 120A, 30m, 480V, 2 AWG Cu, Tubería PVC
	// R = 0.62 Ω/km, X = 0.148 Ω/km
	// Z = √(0.62² + 0.148²) = 0.6374 Ω/km
	// L = 0.030 km
	// e = √3 × 120 × 0.6374 × 0.030 = 3.9677 V
	// %e = 3.9677 / 480 × 100 = 0.8266%

	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.148,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		HilosPorFase:        1,
		SistemaElectrico:    entity.SistemaElectricoDelta,
	}

	corriente, _ := valueobject.NewCorriente(120)
	tension, _ := valueobject.NewTension(480)

	resultado, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
	require.NoError(t, err)

	assert.InDelta(t, 0.8266, resultado.Porcentaje, 0.01, "porcentaje")
	assert.InDelta(t, 3.9677, resultado.CaidaVolts, 0.01, "caida volts")
	assert.True(t, resultado.Cumple, "cumple límite 3%")
	assert.InDelta(t, 0.6374, resultado.Impedancia, 0.001, "impedancia Z")
}
