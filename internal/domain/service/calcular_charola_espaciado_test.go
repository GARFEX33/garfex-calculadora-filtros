package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularCharolaEspaciado(t *testing.T) {
	tablaCharola := []valueobject.EntradaTablaCanalizacion{
		{Tamano: "50mm", AreaInteriorMM2: 2500},
		{Tamano: "100mm", AreaInteriorMM2: 5000},
		{Tamano: "150mm", AreaInteriorMM2: 7500},
		{Tamano: "200mm", AreaInteriorMM2: 10000},
		{Tamano: "300mm", AreaInteriorMM2: 15000},
		{Tamano: "450mm", AreaInteriorMM2: 22500},
		{Tamano: "600mm", AreaInteriorMM2: 30000},
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
		assert.Equal(t, "100mm", result.Tamano)
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
		assert.Equal(t, "50mm", result.Tamano)
	})

	t.Run("Empty table returns error", func(t *testing.T) {
		emptyTable := []valueobject.EntradaTablaCanalizacion{}
		conductorFase := service.ConductorConDiametro{DiametroMM: 10.0}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 5.0}

		_, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			emptyTable,
		)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrCharolaNoEncontrada)
	})

	t.Run("hilosPorFase greater than 1", func(t *testing.T) {
		conductorFase := service.ConductorConDiametro{DiametroMM: 10.0}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 5.0}

		result, err := service.CalcularCharolaEspaciado(
			2,
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		assert.Equal(t, "100mm", result.Tamano)
	})

	t.Run("Delta system requires 3 conductors", func(t *testing.T) {
		conductorFase := service.ConductorConDiametro{DiametroMM: 10.0}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 5.0}

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		anchoRequerido := float64(3-1)*10.0 + 5.0
		assert.Equal(t, anchoRequerido, result.AreaTotal)
	})

	t.Run("Estrella system requires 4 conductors (with neutro)", func(t *testing.T) {
		conductorFase := service.ConductorConDiametro{DiametroMM: 10.0}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 5.0}

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoEstrella,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		anchoRequerido := float64(4-1)*10.0 + 5.0
		assert.Equal(t, anchoRequerido, result.AreaTotal)
	})

	t.Run("Bifasico system requires 3 conductors (no neutro)", func(t *testing.T) {
		conductorFase := service.ConductorConDiametro{DiametroMM: 10.0}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 5.0}

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoBifasico,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		anchoRequerido := float64(3-1)*10.0 + 5.0
		assert.Equal(t, anchoRequerido, result.AreaTotal)
	})

	t.Run("Monofasico system requires 2 conductors (with neutro)", func(t *testing.T) {
		conductorFase := service.ConductorConDiametro{DiametroMM: 10.0}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 5.0}

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoMonofasico,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		anchoRequerido := float64(3+1-1)*10.0 + 5.0
		assert.Equal(t, anchoRequerido, result.AreaTotal)
	})

	t.Run("hilosPorFase less than 1 returns error", func(t *testing.T) {
		conductorFase := service.ConductorConDiametro{DiametroMM: 10.0}
		conductorTierra := service.ConductorConDiametro{DiametroMM: 5.0}

		_, err := service.CalcularCharolaEspaciado(
			0,
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.Error(t, err)
	})
}
