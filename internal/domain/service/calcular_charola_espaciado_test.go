package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularCharolaEspaciado(t *testing.T) {
	tablaCharola := []struct {
		Tamano  string
		AnchoMM float64
	}{
		{"6", 152.4},
		{"9", 228.6},
		{"12", 304.8},
		{"16", 406.4},
		{"18", 457.2},
		{"20", 508.0},
	}

	t.Run("Delta 3 hilos - conductor 4 AWG (25.48mm) + tierra 8 AWG (8.5mm)", func(t *testing.T) {
		conductorFase := service.ConductorConDiametro{DiametroMM: 25.48}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 8.5}

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		assert.Equal(t, "6", result.Tamano)
	})

	t.Run("Estrella 3 hilos + neutro - conductor 2 AWG (7.42mm) + tierra 6 AWG (4.11mm)", func(t *testing.T) {
		conductorFase := service.ConductorConDiametro{DiametroMM: 7.42}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 4.11}

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoEstrella,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		assert.Equal(t, "6", result.Tamano)
	})
}
