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
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 25.48})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 8.5})

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.NoError(t, err)
		// Delta: 3 hilos, 2*3*25.48 + 8.5 = 161.38mm -> charola 200mm
		assert.Equal(t, "200mm", result.Tamano)
	})

	t.Run("hilosPorFase greater than 1", func(t *testing.T) {
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		result, err := service.CalcularCharolaEspaciado(
			2,
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.NoError(t, err)
		// Delta: 3 fases * 2 hilos = 6 hilos, 2*6*10 + 5 = 125mm -> charola 150mm
		assert.Equal(t, "150mm", result.Tamano)
	})

	t.Run("Delta system requires 3 conductors", func(t *testing.T) {
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.NoError(t, err)
		// Delta: 3 hilos, sin control
		// Total = 2 * 3 * 10 + 5 = 65mm
		anchoRequerido := 2.0*3.0*10.0 + 5.0
		assert.Equal(t, anchoRequerido, result.AnchoRequerido)
	})

	t.Run("Estrella system requires 4 conductors (with neutro)", func(t *testing.T) {
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoEstrella,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.NoError(t, err)
		// Estrella: 4 hilos, sin control
		// Total = 2 * 4 * 10 + 5 = 85mm
		anchoRequerido := 2.0*4.0*10.0 + 5.0
		assert.Equal(t, anchoRequerido, result.AnchoRequerido)
	})

	t.Run("Bifasico system requires 3 conductors (with neutro)", func(t *testing.T) {
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoBifasico,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.NoError(t, err)
		// Bifasico: 3 hilos, sin control
		// Total = 2 * 3 * 10 + 5 = 65mm
		anchoRequerido := 2.0*3.0*10.0 + 5.0
		assert.Equal(t, anchoRequerido, result.AnchoRequerido)
	})

	t.Run("Monofasico system requires 2 conductors (with neutro)", func(t *testing.T) {
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoMonofasico,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.NoError(t, err)
		// Monofasico: 2 hilos, sin control
		// Total = 2 * 2 * 10 + 5 = 45mm
		anchoRequerido := 2.0*2.0*10.0 + 5.0
		assert.Equal(t, anchoRequerido, result.AnchoRequerido)
	})

	t.Run("hilosPorFase less than 1 returns error", func(t *testing.T) {
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		_, err := service.CalcularCharolaEspaciado(
			0,
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.Error(t, err)
	})

	t.Run("Con cables de control", func(t *testing.T) {
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		cableControl, _ := valueobject.NewCableControl(valueobject.CableControlParams{
			Cantidad:   1,
			DiametroMM: 4.0,
		})
		cablesControl := []valueobject.CableControl{cableControl}

		result, err := service.CalcularCharolaEspaciado(
			1,
			entity.SistemaElectricoMonofasico,
			conductorFase,
			conductorTierra,
			tablaCharola,
			cablesControl,
		)

		require.NoError(t, err)
		// Monofasico: 2 hilos (1 fase + 1 neutro), 1 cable control de 4mm
		// Formula: 2*hilos*Ø_fase + 3*Ø_control + Ø_tierra
		// Total = 2*2*10 + 3*4 + 5 = 40 + 12 + 5 = 57mm
		anchoRequerido := 2.0*2.0*10.0 + 3.0*4.0 + 5.0
		assert.Equal(t, anchoRequerido, result.AnchoRequerido)
	})
}
