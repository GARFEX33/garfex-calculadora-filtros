// internal/calculos/domain/service/calculo_caida_tension_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalcularCaidaTension_FormulaIEEE141 verifica la fórmula NOM con FP.
//
// Fórmula: %e = 173 × (In/CF) × L_km × (R·cosθ + X·senθ) / E_FF
//
// Caso referencia FP=1 (FA/FR/TR): 2 AWG Cu, tubería PVC, 120A, 30m, 480V
//
//	R = 0.62 Ω/km, X = 0.051 Ω/km (Tabla 9), FP = 1.0
//	term = 0.62×1.0 + 0.051×0.0 = 0.62
//	%e   = 173 × 120 × 0.030 × 0.62 / 480 = 387.07 / 480 / 100 × 100 = 0.807%
//	VD   = 480 × 0.00807 = 3.874 V
//
// Caso FP=0.85 (Carga): mismo conductor
//
//	cosθ = 0.85, senθ = √(1-0.7225) = 0.5268
//	term = 0.62×0.85 + 0.051×0.5268 = 0.527 + 0.02687 = 0.5539
//	%e   = 173 × 120 × 0.030 × 0.5539 / 480 = 344.83 / 480 = 0.719%
//	VD   = 480 × 0.00719 = 3.451 V
func TestCalcularCaidaTension_FormulaIEEE141(t *testing.T) {
	tests := []struct {
		name               string
		entrada            service.EntradaCalculoCaidaTension
		corrienteA         float64
		distanciaM         float64
		voltaje            int
		limitePorc         float64
		expectedPorcentaje float64
		expectedVD         float64
		expectedCumple     bool
		skipNumeric        bool
	}{
		{
			name: "FP=1.0 (FiltroActivo/FiltroRechazo/Transformador) - solo resistencia importa",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    entity.SistemaElectricoEstrella,
				HilosPorFase:        1,
				FactorPotencia:      1.0,
			},
			corrienteA:         120,
			distanciaM:         30,
			voltaje:            480,
			limitePorc:         3.0,
			expectedPorcentaje: 0.80445,
			expectedVD:         3.861,
			expectedCumple:     true,
		},
		{
			name: "FP=0.85 (Carga) - reactancia contribuye",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    entity.SistemaElectricoEstrella,
				HilosPorFase:        1,
				FactorPotencia:      0.85,
			},
			corrienteA:         120,
			distanciaM:         30,
			voltaje:            480,
			limitePorc:         3.0,
			expectedPorcentaje: 0.719,
			expectedVD:         3.451,
			expectedCumple:     true,
		},
		{
			name: "2 hilos por fase reduce R y X a la mitad",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    entity.SistemaElectricoEstrella,
				HilosPorFase:        2,
				FactorPotencia:      1.0,
			},
			corrienteA: 120,
			distanciaM: 30,
			voltaje:    480,
			limitePorc: 3.0,
			// R_ef = 0.62/2 = 0.31, X_ef = 0.051/2 = 0.0255
			// term = 0.31×1.0 + 0.0255×0.0 = 0.31
			// %e   = 173 × 120 × 0.030 × 0.31 / 480 = 0.4035%
			expectedPorcentaje: 0.4035,
			expectedVD:         1.937,
			expectedCumple:     true,
		},
		{
			name: "charola espaciado - misma formula, X de reactancia_al",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionCharolaCableEspaciado,
				SistemaElectrico:    entity.SistemaElectricoEstrella,
				HilosPorFase:        1,
				FactorPotencia:      1.0,
			},
			corrienteA:         120,
			distanciaM:         30,
			voltaje:            480,
			limitePorc:         3.0,
			expectedPorcentaje: 0.80445,
			expectedVD:         3.861,
			expectedCumple:     true,
		},
		{
			name: "excede limite NOM 3%",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 5.21,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    entity.SistemaElectricoEstrella,
				HilosPorFase:        1,
				FactorPotencia:      1.0,
			},
			corrienteA:     25,
			distanciaM:     100,
			voltaje:        220,
			limitePorc:     3.0,
			expectedCumple: false,
			skipNumeric:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corriente, err := valueobject.NewCorriente(tt.corrienteA)
			require.NoError(t, err)

			tension, err := valueobject.NewTension(tt.voltaje)
			require.NoError(t, err)

			resultado, err := service.CalcularCaidaTension(
				tt.entrada, corriente, tt.distanciaM, tension, tt.limitePorc,
			)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCumple, resultado.Cumple)

			if !tt.skipNumeric {
				assert.InDelta(t, tt.expectedPorcentaje, resultado.Porcentaje, 0.01, "porcentaje")
				assert.InDelta(t, tt.expectedVD, resultado.CaidaVolts, 0.01, "caida volts")
			}
		})
	}
}

func TestCalcularCaidaTension_ResultadoContieneRXTerminoEfectivo(t *testing.T) {
	// Verifica que el struct expone R_ef, X_ef y el término efectivo para el reporte
	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.051,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    entity.SistemaElectricoEstrella,
		HilosPorFase:        1,
		FactorPotencia:      0.85,
	}
	corriente, _ := valueobject.NewCorriente(120)
	tension, _ := valueobject.NewTension(480)

	resultado, err := service.CalcularCaidaTension(entrada, corriente, 30, tension, 3.0)
	require.NoError(t, err)

	// R_ef = 0.62 / 1 = 0.62
	assert.InDelta(t, 0.62, resultado.Resistencia, 0.001, "R_ef")
	// X_ef = 0.051 / 1 = 0.051
	assert.InDelta(t, 0.051, resultado.Reactancia, 0.001, "X_ef")
	// término efectivo = R·cosθ + X·senθ = 0.62×0.85 + 0.051×0.5268 = 0.5539
	assert.InDelta(t, 0.5539, resultado.Impedancia, 0.002, "término efectivo")
}

func TestCalcularCaidaTension_ErrorDistanciaInvalida(t *testing.T) {
	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.051,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    entity.SistemaElectricoEstrella,
		HilosPorFase:        1,
		FactorPotencia:      1.0,
	}
	corriente, _ := valueobject.NewCorriente(120)
	tension, _ := valueobject.NewTension(480)

	t.Run("distancia cero", func(t *testing.T) {
		_, err := service.CalcularCaidaTension(entrada, corriente, 0, tension, 3.0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrDistanciaInvalida)
	})

	t.Run("distancia negativa", func(t *testing.T) {
		_, err := service.CalcularCaidaTension(entrada, corriente, -10, tension, 3.0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrDistanciaInvalida)
	})
}

func TestCalcularCaidaTension_ErrorHilosPorFaseInvalido(t *testing.T) {
	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.051,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    entity.SistemaElectricoEstrella,
		HilosPorFase:        0,
		FactorPotencia:      1.0,
	}
	corriente, _ := valueobject.NewCorriente(120)
	tension, _ := valueobject.NewTension(480)

	_, err := service.CalcularCaidaTension(entrada, corriente, 30, tension, 3.0)
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrHilosPorFaseInvalido)
}

func TestCalcularCaidaTension_ErrorFactorPotenciaInvalido(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(120)
	tension, _ := valueobject.NewTension(480)

	casos := []struct {
		nombre string
		fp     float64
	}{
		{"FP cero", 0.0},
		{"FP negativo", -0.5},
		{"FP mayor que 1", 1.1},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			entrada := service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				SistemaElectrico:    entity.SistemaElectricoEstrella,
				HilosPorFase:        1,
				FactorPotencia:      c.fp,
			}
			_, err := service.CalcularCaidaTension(entrada, corriente, 30, tension, 3.0)
			require.Error(t, err)
			assert.ErrorIs(t, err, service.ErrFactorPotenciaInvalido)
		})
	}
}
