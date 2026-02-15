// internal/calculos/domain/service/calcular_charola_triangular_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
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
		// Formula: anchoPotencia + espacioFuerza + tierra
		// anchoPotencia = 2 * 25.48 = 50.96mm
		// espacioFuerza = (2-1) * 2.15 * 25.48 = 54.78mm
		// Total = 50.96 + 54.78 + 7.42 = 113.16mm
		// Charola 6" (152.4mm) es suficiente

		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 25.48})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 7.42})

		result, err := service.CalcularCharolaTriangular(
			2,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.NoError(t, err)
		assert.Equal(t, "6", result.Tamano)
	})

	t.Run("1 hilo por fase - conductor pequeño", func(t *testing.T) {
		// Formula: anchoPotencia + espacioFuerza + tierra
		// anchoPotencia = 2 * 1 * 2.15 * 10 = 43mm
		// espacioFuerza = (1-1) * 2.15 * 10 = 0mm
		// Total = 43 + 0 + 5 = 48mm
		// Requiere charola de 6"

		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		result, err := service.CalcularCharolaTriangular(
			1,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.NoError(t, err)
		assert.Equal(t, "6", result.Tamano)
	})

	t.Run("error: hilosPorFase menor que 1", func(t *testing.T) {
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		_, err := service.CalcularCharolaTriangular(
			0,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "hilos por fase debe ser >= 1")
	})

	t.Run("error: tabla vacía", func(t *testing.T) {
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		_, err := service.CalcularCharolaTriangular(
			2,
			conductorFase,
			conductorTierra,
			[]valueobject.EntradaTablaCanalizacion{},
			nil,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "tabla vacía")
	})

	t.Run("error: ninguna charola suficiente", func(t *testing.T) {
		// Conductor muy grande que no cabe en ninguna charola de la tabla
		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 100.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 50.0})

		_, err := service.CalcularCharolaTriangular(
			5,
			conductorFase,
			conductorTierra,
			tablaCharola,
			nil,
		)

		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrCharolaTriangularNoEncontrada)
	})

	t.Run("Con cables de control", func(t *testing.T) {
		// Formula: anchoPotencia + espacioFuerza + espacioControl + anchoControl + tierra
		// anchoPotencia = 2 * 10 = 20mm
		// espacioFuerza = (1-1) * 2.15 * 10 = 0mm
		// espacioControl = 2.15 * 4 = 8.6mm
		// anchoControl = 4mm
		// Total = 20 + 0 + 8.6 + 4 + 5 = 37.6mm
		// Charola 6" (152.4mm) es suficiente

		conductorFase, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 10.0})
		conductorTierra, _ := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 5.0})

		cableControl, _ := valueobject.NewCableControl(valueobject.CableControlParams{
			Cantidad:   1,
			DiametroMM: 4.0,
		})
		cablesControl := []valueobject.CableControl{cableControl}

		result, err := service.CalcularCharolaTriangular(
			1,
			conductorFase,
			conductorTierra,
			tablaCharola,
			cablesControl,
		)

		require.NoError(t, err)
		anchoRequerido := 37.6
		assert.Equal(t, anchoRequerido, result.AnchoRequerido)
		assert.Equal(t, "6", result.Tamano)
	})
}
