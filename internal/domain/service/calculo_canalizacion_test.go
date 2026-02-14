package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Simplified conduit sizing table (tubería EMT)
var tablaCanalizacionTest = []valueobject.EntradaTablaCanalizacion{
	{Tamano: "1/2", AreaInteriorMM2: 78.0},
	{Tamano: "3/4", AreaInteriorMM2: 122.0},
	{Tamano: "1", AreaInteriorMM2: 198.0},
	{Tamano: "1 1/4", AreaInteriorMM2: 277.0},
	{Tamano: "1 1/2", AreaInteriorMM2: 360.0},
	{Tamano: "2", AreaInteriorMM2: 572.0},
	{Tamano: "2 1/2", AreaInteriorMM2: 885.0},
	{Tamano: "3", AreaInteriorMM2: 1327.0},
	{Tamano: "4", AreaInteriorMM2: 2165.0},
}

func TestCalcularCanalizacion_Tuberia(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 33.62}, // 3 phases × 2 AWG
		{Cantidad: 1, SeccionMM2: 13.30}, // 1 ground × 6 AWG
	}
	// Total area = 3×33.62 + 1×13.30 = 114.16 mm²
	// Required conduit area at 40% fill = 114.16 / 0.40 = 285.4 mm²
	// Smallest conduit ≥ 285.4 → "1 1/2" (360 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 1)
	require.NoError(t, err)
	assert.Equal(t, "TUBERIA_CONDUIT", result.Tipo)
	assert.Equal(t, "1 1/2", result.Tamano)
	assert.InDelta(t, 114.16, result.AnchoRequerido, 0.01)
}

func TestCalcularCanalizacion_SmallConductors(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 3.31}, // 3 × 12 AWG
		{Cantidad: 1, SeccionMM2: 2.08}, // 1 × 14 AWG ground
	}
	// Total = 3×3.31 + 1×2.08 = 12.01 mm²
	// Required = 12.01 / 0.40 = 30.025 mm²
	// Smallest ≥ 30.025 → "1/2" (78 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 1)
	require.NoError(t, err)
	assert.Equal(t, "1/2", result.Tamano)
}

func TestCalcularCanalizacion_NoFit(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 20, SeccionMM2: 253.4},
	}
	// Total = 20 × 253.4 = 5068 mm² → required = 12670 mm² → exceeds all conduits

	_, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 1)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrCanalizacionNoDisponible))
}

func TestCalcularCanalizacion_EmptyConductors(t *testing.T) {
	_, err := service.CalcularCanalizacion(nil, "TUBERIA_CONDUIT", tablaCanalizacionTest, 1)
	assert.Error(t, err)
}

func TestCalcularCanalizacion_FillFactor1Conductor(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 1, SeccionMM2: 8.37}, // 1 × 8 AWG
	}
	// Total area = 1×8.37 = 8.37 mm²
	// Fill factor for 1 conductor = 53%
	// Required conduit area at 53% fill = 8.37 / 0.53 = 15.79 mm²
	// Smallest conduit ≥ 15.79 → "1/2" (78 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 1)
	require.NoError(t, err)
	assert.Equal(t, "1/2", result.Tamano)
}

func TestCalcularCanalizacion_FillFactor2Conductores(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 2, SeccionMM2: 3.31}, // 2 × 12 AWG (phase + neutral)
	}
	// Total area = 2×3.31 = 6.62 mm²
	// Fill factor for 2 conductors = 31%
	// Required conduit area at 31% fill = 6.62 / 0.31 = 21.35 mm²
	// Smallest conduit ≥ 21.35 → "1/2" (78 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 1)
	require.NoError(t, err)
	assert.Equal(t, "1/2", result.Tamano)
}

func TestCalcularCanalizacion_FillFactor3Conductores(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 3.31}, // 3 × 12 AWG (3-phase)
	}
	// Total area = 3×3.31 = 9.93 mm²
	// Fill factor for 3+ conductors = 40%
	// Required conduit area at 40% fill = 9.93 / 0.40 = 24.83 mm²
	// Smallest conduit ≥ 24.83 → "1/2" (78 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 1)
	require.NoError(t, err)
	assert.Equal(t, "1/2", result.Tamano)
}

func TestCalcularCanalizacion_DosTubos(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 33.62}, // 3 phases × 2 AWG
		{Cantidad: 1, SeccionMM2: 13.30}, // 1 ground × 6 AWG
	}
	// Total area = 114.16 mm², cantidadTotal = 4
	// Con 2 tubos: conductoresPorTubo = 4/2 = 2 → fillFactor = 0.31
	// areaPorTubo = 114.16/2 = 57.08 mm²
	// areaRequerida = 57.08/0.31 = 184.13 mm²
	// Smallest ≥ 184.13 → "1" (198 mm²)
	// Con 1 tubo habría dado "1 1/2" (360 mm²) — el tubo individual es más chico

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 2)
	require.NoError(t, err)
	assert.Equal(t, "TUBERIA_CONDUIT", result.Tipo)
	assert.Equal(t, "1", result.Tamano)
	assert.Equal(t, 2, result.NumeroDeTubos)
	assert.InDelta(t, 114.16, result.AnchoRequerido, 0.01)
}

func TestCalcularCanalizacion_DosTubosSmall(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 4, SeccionMM2: 3.31}, // 4 × 12 AWG
	}
	// Total area = 13.24 mm², cantidadTotal = 4
	// Con 2 tubos: conductoresPorTubo = 4/2 = 2 → fillFactor = 0.31
	// areaPorTubo = 13.24/2 = 6.62 mm²
	// areaRequerida = 6.62/0.31 = 21.35 mm²
	// Smallest ≥ 21.35 → "1/2" (78 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 2)
	require.NoError(t, err)
	assert.Equal(t, "1/2", result.Tamano)
	assert.Equal(t, 2, result.NumeroDeTubos)
}

func TestCalcularCanalizacion_NumeroDeTubosCero(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 3.31},
	}
	_, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "numeroDeTubos debe ser mayor a cero")
}

func TestCalcularCanalizacion_NumeroDeTubosNegativo(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 3.31},
	}
	_, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "numeroDeTubos debe ser mayor a cero")
}

func TestCalcularCanalizacion_NumeroDeTubosMayorQueConductores(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 2, SeccionMM2: 3.31},
	}
	// 2 conductores, 5 tubos — permitido, el servicio no juzga
	// cantidadTotal=2, conductoresPorTubo=2/5=0 → fillFactor=0.40 (default)
	// areaPorTubo = 6.62/5 = 1.324 mm²
	// areaRequerida = 1.324/0.40 = 3.31 mm²
	// Smallest ≥ 3.31 → "1/2" (78 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 5)
	require.NoError(t, err)
	assert.Equal(t, "1/2", result.Tamano)
	assert.Equal(t, 5, result.NumeroDeTubos)
}
