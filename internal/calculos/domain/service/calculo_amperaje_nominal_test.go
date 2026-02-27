// internal/calculos/domain/service/calculo_amperaje_nominal_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularAmperajeNominalCircuito_Monofasico(t *testing.T) {
	// 220V, 1000W, FP=0.9 → I = 1000 / (220 × 0.9) ≈ 5.05 A
	tension, err := valueobject.NewTension(220, "V")
	require.NoError(t, err)

	potencia, err := valueobject.NewPotencia(1000, "W")
	require.NoError(t, err)

	corriente, err := service.CalcularAmperajeNominalCircuito(
		potencia,
		tension,
		entity.SistemaElectricoMonofasico,
		0.9,
	)
	require.NoError(t, err)

	assert.InDelta(t, 5.0505, corriente.Valor(), 0.01)
}

func TestCalcularAmperajeNominalCircuito_Trifasico(t *testing.T) {
	// 480V, 50000W, FP=0.85 → I = 50000 / (480 × √3 × 0.85) ≈ 70.75 A
	tension, err := valueobject.NewTension(480, "V")
	require.NoError(t, err)

	potencia, err := valueobject.NewPotencia(50000, "W")
	require.NoError(t, err)

	corriente, err := service.CalcularAmperajeNominalCircuito(
		potencia,
		tension,
		entity.SistemaElectricoDelta,
		0.85,
	)
	require.NoError(t, err)

	assert.InDelta(t, 70.75, corriente.Valor(), 0.1)
}

func TestCalcularAmperajeNominalCircuito_127V_Monofasico(t *testing.T) {
	// Caso: circuito monofásico 127V, 2000W, FP=0.95
	// I = 2000 / (127 × 0.95) ≈ 16.54 A
	tension, err := valueobject.NewTension(127, "V")
	require.NoError(t, err)

	potencia, err := valueobject.NewPotencia(2000, "W")
	require.NoError(t, err)

	corriente, err := service.CalcularAmperajeNominalCircuito(
		potencia,
		tension,
		entity.SistemaElectricoMonofasico,
		0.95,
	)
	require.NoError(t, err)

	assert.InDelta(t, 16.54, corriente.Valor(), 0.1)
}

func TestCalcularAmperajeNominalCircuito_Errores(t *testing.T) {
	tension, err := valueobject.NewTension(220, "V")
	require.NoError(t, err)

	t.Run("potencia_negativa", func(t *testing.T) {
		// NewPotencia already validates, so we test the value object creation
		_, err := valueobject.NewPotencia(-100, "W")
		require.Error(t, err)
		assert.ErrorIs(t, err, valueobject.ErrPotenciaInvalida)
	})

	t.Run("potencia_valida", func(t *testing.T) {
		// Verify valid potencia works with the service
		potencia, err := valueobject.NewPotencia(1000, "W")
		require.NoError(t, err)
		_, err = service.CalcularAmperajeNominalCircuito(
			potencia,
			tension,
			entity.SistemaElectricoMonofasico,
			0.9,
		)
		require.NoError(t, err)
	})

	t.Run("factor_potencia_cero", func(t *testing.T) {
		potencia, err := valueobject.NewPotencia(1000, "W")
		require.NoError(t, err)
		_, err = service.CalcularAmperajeNominalCircuito(
			potencia,
			tension,
			entity.SistemaElectricoMonofasico,
			0,
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrFactorPotenciaInvalido)
	})

	t.Run("factor_potencia_mayor_1", func(t *testing.T) {
		potencia, err := valueobject.NewPotencia(1000, "W")
		require.NoError(t, err)
		_, err = service.CalcularAmperajeNominalCircuito(
			potencia,
			tension,
			entity.SistemaElectricoMonofasico,
			1.5,
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrFactorPotenciaInvalido)
	})
}
