package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularCharolaTriangular(t *testing.T) {
	tablaCharola := []valueobject.EntradaTablaCanalizacion{
		{Tamano: "6", AreaInteriorMM2: 152.4},
		{Tamano: "9", AreaInteriorMM2: 228.6},
		{Tamano: "12", AreaInteriorMM2: 304.8},
	}

	t.Run("2 hilos por fase - conductor 500 KCM (25.48mm) + tierra 2 AWG (7.42mm)", func(t *testing.T) {
		// Formula: [(2 - 1) * 2.15 * 25.48] + 7.42 = 54.78 + 7.42 = 62.2mm
		// Charola 6" (152.4mm) es suficiente

		conductorFase := service.ConductorConDiametro{DiametroMM: 25.48}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 7.42}

		result, err := service.CalcularCharolaTriangular(
			2,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		assert.Equal(t, "6", result.Tamano)
	})

	t.Run("1 hilo por fase - conductor pequeño", func(t *testing.T) {
		// Formula: [(1 - 1) * 2.15 * 10] + 5 = 0 + 5 = 5mm
		// Requiere charola de 6"

		conductorFase := service.ConductorConDiametro{DiametroMM: 10.0}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 5.0}

		result, err := service.CalcularCharolaTriangular(
			1,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		assert.Equal(t, "6", result.Tamano)
	})
}
