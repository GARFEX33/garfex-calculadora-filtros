// internal/calculos/domain/service/calculo_caida_tension_sistemas_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalcularCaidaTension_SistemasElectricos verifica la fórmula para diferentes sistemas eléctricos.
//
// Monofásico y Bifásico usan factor 2:
//
//	%e = 2 × I × L_km × (R·cosθ + X·senθ) / V × 100
//
// Trifásico (Estrella/Delta) usa factor √3:
//
//	%e = √3 × I × L_km × (R·cosθ + X·senθ) / V × 100
func TestCalcularCaidaTension_SistemasElectricos(t *testing.T) {
	// Datos base: 2 AWG Cu, tubería PVC, 50A, 20m, 220V, FP=1.0
	// R = 0.62 Ω/km, X = 0.051 Ω/km
	// L_km = 0.020 km
	// término = 0.62×1.0 + 0.051×0 = 0.62 Ω/km

	tests := []struct {
		name               string
		sistema            entity.SistemaElectrico
		expectedPorcentaje float64 // Diferente por el factor (2 vs √3)
		expectedVD         float64
	}{
		{
			name:    "Monofásico usa factor 2",
			sistema: entity.SistemaElectricoMonofasico,
			// %e = 2 × 50 × 0.020 × 0.62 / 220 × 100 = 1.24 / 220 × 100 = 0.5636%
			expectedPorcentaje: 0.5636,
			expectedVD:         1.24,
		},
		{
			name:    "Bifásico usa factor 2",
			sistema: entity.SistemaElectricoBifasico,
			// %e = 2 × 50 × 0.020 × 0.62 / 220 × 100 = 0.5636%
			expectedPorcentaje: 0.5636,
			expectedVD:         1.24,
		},
		{
			name:    "Trifásico Estrella usa factor √3",
			sistema: entity.SistemaElectricoEstrella,
			// %e = 1.732 × 50 × 0.020 × 0.62 / 220 × 100 = 1.074 / 220 × 100 = 0.4881%
			expectedPorcentaje: 0.4881,
			expectedVD:         1.074,
		},
		{
			name:    "Trifásico Delta usa factor √3",
			sistema: entity.SistemaElectricoDelta,
			// %e = 1.732 × 50 × 0.020 × 0.62 / 220 × 100 = 0.4881%
			expectedPorcentaje: 0.4881,
			expectedVD:         1.074,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entrada := service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				HilosPorFase:        1,
				FactorPotencia:      1.0,
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
			assert.InDelta(t, tt.expectedVD, resultado.CaidaVolts, 0.01, "caida volts")
			assert.True(t, resultado.Cumple, "debe cumplir límite 3%")
		})
	}
}

// TestCalcularCaidaTension_MonofasicoVsTrifasico compara directamente monofásico vs trifásico.
func TestCalcularCaidaTension_MonofasicoVsTrifasico(t *testing.T) {
	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.051,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		HilosPorFase:        1,
		FactorPotencia:      1.0,
	}

	corriente, _ := valueobject.NewCorriente(100)
	tension, _ := valueobject.NewTension(220)

	// Monofásico
	entrada.SistemaElectrico = entity.SistemaElectricoMonofasico
	resultadoMono, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
	require.NoError(t, err)

	// Trifásico
	entrada.SistemaElectrico = entity.SistemaElectricoEstrella
	resultadoTri, err := service.CalcularCaidaTension(entrada, corriente, 30.0, tension, 3.0)
	require.NoError(t, err)

	// Monofásico debe tener mayor caída (factor 2 vs √3)
	// Relación esperada: mono / tri = 2 / √3 = 1.1547
	ratio := resultadoMono.Porcentaje / resultadoTri.Porcentaje
	assert.InDelta(t, 1.1547, ratio, 0.01, "relación monofásico/trifásico debe ser 2/√3")
	assert.Greater(t, resultadoMono.Porcentaje, resultadoTri.Porcentaje, "monofásico debe tener mayor caída")
}
