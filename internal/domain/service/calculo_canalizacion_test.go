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
