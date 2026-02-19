// internal/calculos/domain/service/calcular_tamanio_tuberia_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tabla de ocupación PVC 40% (para tests)
var tablaOcupacionTest = []valueobject.EntradaTablaOcupacion{
	{Tamano: "1/2", AreaOcupacionMM2: 74, DesignacionMetrica: "16"},
	{Tamano: "3/4", AreaOcupacionMM2: 131, DesignacionMetrica: "21"},
	{Tamano: "1", AreaOcupacionMM2: 214, DesignacionMetrica: "27"},
	{Tamano: "1 1/4", AreaOcupacionMM2: 374, DesignacionMetrica: "35"},
	{Tamano: "1 1/2", AreaOcupacionMM2: 513, DesignacionMetrica: "41"},
	{Tamano: "2", AreaOcupacionMM2: 849, DesignacionMetrica: "53"},
	{Tamano: "2 1/2", AreaOcupacionMM2: 1212, DesignacionMetrica: "63"},
	{Tamano: "3", AreaOcupacionMM2: 1877, DesignacionMetrica: "78"},
	{Tamano: "4", AreaOcupacionMM2: 3237, DesignacionMetrica: "103"},
}

func TestCalcularAreaPorTubo_SingleTube(t *testing.T) {
	// 3 fases × 86 mm² + 1 neutro × 86 mm² + 1 tierra × 17.09 mm² = 361.09 mm²
	area := service.CalcularAreaPorTubo(
		3,     // fases
		1,     // neutros
		1,     // tierras
		86.0,  // área fase
		86.0,  // área neutro
		17.09, // área tierra
		1,     // 1 tubería
	)

	// fases por tubo: 3/1 = 3 → 3 × 86 = 258
	// neutros por tubo: 1/1 = 1 → 1 × 86 = 86
	// tierras por tubo: 1 × 17.09 × 1 = 17.09
	// total: 258 + 86 + 17.09 = 361.09
	assert.InDelta(t, 361.09, area, 0.01)
}

func TestCalcularAreaPorTubo_MultipleTubes(t *testing.T) {
	// 6 fases, 2 neutros, 1 tierra en 2 tubos
	// fases por tubo: 6/2 = 3 → 3 × 86 = 258
	// neutros por tubo: 2/2 = 1 → 1 × 86 = 86
	// tierras por tubo: 1 × 17.09 × 2 = 34.18 (NO se divide)
	// total por tubo: 258 + 86 + 34.18 = 378.18
	area := service.CalcularAreaPorTubo(
		6, // fases
		2, // neutros
		1, // tierras
		86.0,
		86.0,
		17.09,
		2, // 2 tuberías
	)

	expected := 378.18
	assert.InDelta(t, expected, area, 0.01)
}

func TestCalcularAreaPorTubo_ZeroNeutros(t *testing.T) {
	// Sin neutros - solo fases y tierra
	area := service.CalcularAreaPorTubo(
		3,     // fases
		0,     // neutros
		1,     // tierras
		33.6,  // área fase 2 AWG
		0,     // área neutro (no hay)
		17.09, // área tierra 6 AWG
		1,     // 1 tubería
	)

	// 3 × 33.6 + 0 + 1 × 17.09 = 117.89
	assert.InDelta(t, 117.89, area, 0.01)
}

func TestCalcularAreaRequerida(t *testing.T) {
	// Con factor 40%: área_requerida = área / 0.40
	area := service.CalcularAreaRequerida(100)
	assert.InDelta(t, 250.0, area, 0.01)
}

func TestBuscarTamanioTuberia_FindsCorrectSize(t *testing.T) {
	// Área requerida: 361.09 mm²
	// Tabla PVC 40%: 1 1/4 tiene 374 mm² (primera que cabe)
	result, err := service.BuscarTamanioTuberia(
		361.09,
		entity.TipoCanalizacionTuberiaPVC,
		tablaOcupacionTest,
	)

	require.NoError(t, err)
	assert.Equal(t, "1 1/4", result.TuberiaRecomendada())
	assert.Equal(t, "35", result.DesignacionMetrica())
}

func TestBuscarTamanioTuberia_SmallArea(t *testing.T) {
	// Área requerida: 50 mm² → cabe en 1/2 (74 mm²)
	result, err := service.BuscarTamanioTuberia(
		50.0,
		entity.TipoCanalizacionTuberiaPVC,
		tablaOcupacionTest,
	)

	require.NoError(t, err)
	assert.Equal(t, "1/2", result.TuberiaRecomendada())
}

func TestBuscarTamanioTuberia_NoFit(t *testing.T) {
	// Área requerida muy grande excede la máxima
	_, err := service.BuscarTamanioTuberia(
		5000.0, // mayor que 4" (3237 mm²)
		entity.TipoCanalizacionTuberiaPVC,
		tablaOcupacionTest,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no se encontró tamaño de tubería")
}

func TestBuscarTamanioTuberia_EmptyTable(t *testing.T) {
	_, err := service.BuscarTamanioTuberia(
		100.0,
		entity.TipoCanalizacionTuberiaPVC,
		[]valueobject.EntradaTablaOcupacion{},
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "tabla de ocupación vacía")
}

func TestBuscarTamanioTuberia_InvalidArea(t *testing.T) {
	_, err := service.BuscarTamanioTuberia(
		0.0,
		entity.TipoCanalizacionTuberiaPVC,
		tablaOcupacionTest,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "área requerida debe ser mayor que cero")
}

func TestCalcularTamanioTuberiaWithMultiplePipes(t *testing.T) {
	// Caso del usuario: 3 fases, 1 neutro, 1 tierra, 1 tubería
	// fases: 3 × 86 = 258
	// neutro: 1 × 86 = 86
	// tierra: 1 × 17.09 = 17.09
	// total: 361.09 mm² → 1 1/4 (374 mm²)
	result, err := service.CalcularTamanioTuberiaWithMultiplePipes(
		3,     // fases
		1,     // neutros
		1,     // tierras
		86.0,  // área fase
		86.0,  // área neutro
		17.09, // área tierra
		1,     // 1 tubería
		entity.TipoCanalizacionTuberiaPVC,
		tablaOcupacionTest,
	)

	require.NoError(t, err)
	assert.Equal(t, "1 1/4", result.TuberiaRecomendada())
	assert.Equal(t, "35", result.DesignacionMetrica())
	assert.Equal(t, 1, result.NumTuberias())
}

func TestCalcularTamanioTuberiaWithMultiplePipes_InvalidNumTuberias(t *testing.T) {
	_, err := service.CalcularTamanioTuberiaWithMultiplePipes(
		3, 1, 1,
		86.0, 86.0, 17.09,
		0, // inválido
		entity.TipoCanalizacionTuberiaPVC,
		tablaOcupacionTest,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "numeroDeTubos")
}

func TestDiseñoMetricoConduit(t *testing.T) {
	tests := []struct {
		tamano    string
		esperado  string
		errExpect bool
	}{
		{"1/2", "13mm", false},
		{"3/4", "19mm", false},
		{"1", "25mm", false},
		{"1 1/4", "32mm", false},
		{"2", "51mm", false},
		{"invalid", "", true},
	}

	for _, tt := range tests {
		result, err := service.DiseñoMetricoConduit(tt.tamano)
		if tt.errExpect {
			assert.Error(t, err, tt.tamano)
		} else {
			assert.NoError(t, err, tt.tamano)
			assert.Equal(t, tt.esperado, result, tt.tamano)
		}
	}
}
