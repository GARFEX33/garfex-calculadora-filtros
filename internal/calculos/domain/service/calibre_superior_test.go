// internal/calculos/domain/service/calibre_superior_test.go
package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test ObtenerCalibreSuperior - Escenarios de specs
func TestObtenerCalibreSuperior_EscenarioA_MiddleOfList(t *testing.T) {
	// Scenario A: "1/0" → "2/0" (middle of list)
	resultado, err := service.ObtenerCalibreSuperior("1/0")
	require.NoError(t, err)
	assert.Equal(t, "2/0", resultado)
}

func TestObtenerCalibreSuperior_EscenarioB_AWGtoMCMBoundary(t *testing.T) {
	// Scenario B: "4/0" → "250" (AWG to MCM boundary)
	resultado, err := service.ObtenerCalibreSuperior("4/0")
	require.NoError(t, err)
	assert.Equal(t, "250", resultado)
}

func TestObtenerCalibreSuperior_EscenarioC_SmallestCaliber(t *testing.T) {
	// Scenario C: "14" → "12" (smallest valid caliber)
	resultado, err := service.ObtenerCalibreSuperior("14")
	require.NoError(t, err)
	assert.Equal(t, "12", resultado)
}

func TestObtenerCalibreSuperior_EscenarioD_MaximumCaliber(t *testing.T) {
	// Scenario D: "1000" → error (maximum caliber)
	_, err := service.ObtenerCalibreSuperior("1000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no existe calibre superior a 1000 MCM")
}

func TestObtenerCalibreSuperior_EscenarioE_InvalidCaliber(t *testing.T) {
	// Scenario E: "999" → error (invalid caliber)
	_, err := service.ObtenerCalibreSuperior("999")
	require.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrCalibreNoReconocido))
	assert.Contains(t, err.Error(), "no reconocido en la lista de calibres NOM")
}

// Test table-driven: todos los pares consecutivos en la lista
func TestObtenerCalibreSuperior_TodosLosParesConsecutivos(t *testing.T) {
	paresConsecutivos := []struct {
		calibreActual string
		calibreEsperado string
	}{
		{"14", "12"},
		{"12", "10"},
		{"10", "8"},
		{"8", "6"},
		{"6", "4"},
		{"4", "2"},
		{"2", "1/0"},
		{"1/0", "2/0"},
		{"2/0", "3/0"},
		{"3/0", "4/0"},
		{"4/0", "250"},
		{"250", "300"},
		{"300", "350"},
		{"350", "400"},
		{"400", "500"},
		{"500", "600"},
		{"600", "750"},
		{"750", "1000"},
	}

	for _, tt := range paresConsecutivos {
		t.Run(tt.calibreActual+"_"+tt.calibreEsperado, func(t *testing.T) {
			resultado, err := service.ObtenerCalibreSuperior(tt.calibreActual)
			require.NoError(t, err, "Error obteniendo calibre superior a %s", tt.calibreActual)
			assert.Equal(t, tt.calibreEsperado, resultado, "Para calibre %s se esperaba %s", tt.calibreActual, tt.calibreEsperado)
		})
	}
}
