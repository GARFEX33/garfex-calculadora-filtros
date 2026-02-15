// internal/calculos/domain/service/calculo_tierra_test.go
package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func entradaTierraCu(itmHasta int, calibre string, seccionMM2 float64) valueobject.EntradaTablaTierra {
	return valueobject.EntradaTablaTierra{
		ITMHasta: itmHasta,
		ConductorCu: valueobject.ConductorParams{
			Calibre:    calibre,
			Material:   valueobject.MaterialCobre,
			SeccionMM2: seccionMM2,
		},
		ConductorAl: nil,
	}
}

func entradaTierraCuAl(itmHasta int, cuCalibre string, cuSeccion float64, alCalibre string, alSeccion float64) valueobject.EntradaTablaTierra {
	al := valueobject.ConductorParams{
		Calibre:    alCalibre,
		Material:   valueobject.MaterialAluminio,
		SeccionMM2: alSeccion,
	}
	return valueobject.EntradaTablaTierra{
		ITMHasta: itmHasta,
		ConductorCu: valueobject.ConductorParams{
			Calibre:    cuCalibre,
			Material:   valueobject.MaterialCobre,
			SeccionMM2: cuSeccion,
		},
		ConductorAl: &al,
	}
}

var tablaTierraTest = []valueobject.EntradaTablaTierra{
	entradaTierraCu(15, "14 AWG", 2.08),
	entradaTierraCu(20, "12 AWG", 3.31),
	entradaTierraCu(60, "10 AWG", 5.26),
	entradaTierraCu(100, "8 AWG", 8.37),
	entradaTierraCuAl(200, "6 AWG", 13.3, "4 AWG", 21.2),
	entradaTierraCuAl(400, "2 AWG", 33.6, "1/0 AWG", 42.4),
	entradaTierraCuAl(800, "1/0 AWG", 53.5, "3/0 AWG", 85.0),
	entradaTierraCuAl(1000, "2/0 AWG", 67.4, "4/0 AWG", 107.2),
	entradaTierraCuAl(4000, "500 MCM", 253.0, "750 MCM", 380.0),
}

func TestSeleccionarConductorTierra_CuExplicito(t *testing.T) {
	tests := []struct {
		name            string
		itm             int
		expectedCalibre string
	}{
		{"ITM 15 → 14 AWG Cu", 15, "14 AWG"},
		{"ITM 20 → 12 AWG Cu", 20, "12 AWG"},
		{"ITM 30 → 10 AWG Cu (≤60)", 30, "10 AWG"},
		{"ITM 100 → 8 AWG Cu", 100, "8 AWG"},
		{"ITM 125 → 6 AWG Cu (≤200)", 125, "6 AWG"},
		{"ITM 400 → 2 AWG Cu", 400, "2 AWG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conductor, err := service.SeleccionarConductorTierra(tt.itm, valueobject.MaterialCobre, tablaTierraTest)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCalibre, conductor.Calibre())
			assert.Equal(t, valueobject.MaterialCobre, conductor.Material())
		})
	}
}

func TestSeleccionarConductorTierra_AluminioDisponible(t *testing.T) {
	conductor, err := service.SeleccionarConductorTierra(200, valueobject.MaterialAluminio, tablaTierraTest)
	require.NoError(t, err)
	assert.Equal(t, "4 AWG", conductor.Calibre())
	assert.Equal(t, valueobject.MaterialAluminio, conductor.Material())
}

func TestSeleccionarConductorTierra_AluminioFallbackCu(t *testing.T) {
	conductor, err := service.SeleccionarConductorTierra(60, valueobject.MaterialAluminio, tablaTierraTest)
	require.NoError(t, err)
	assert.Equal(t, "10 AWG", conductor.Calibre())
	assert.Equal(t, valueobject.MaterialCobre, conductor.Material())
}

func TestSeleccionarConductorTierra_AluminioITMMaximo(t *testing.T) {
	conductor, err := service.SeleccionarConductorTierra(4000, valueobject.MaterialAluminio, tablaTierraTest)
	require.NoError(t, err)
	assert.Equal(t, "750 MCM", conductor.Calibre())
	assert.Equal(t, valueobject.MaterialAluminio, conductor.Material())
}

func TestSeleccionarConductorTierra_ITMExceedsTable(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(5000, valueobject.MaterialCobre, tablaTierraTest)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrConductorNoEncontrado))
}

func TestSeleccionarConductorTierra_InvalidITM(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(0, valueobject.MaterialCobre, tablaTierraTest)
	assert.Error(t, err)
}

func TestSeleccionarConductorTierra_EmptyTable(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(100, valueobject.MaterialCobre, nil)
	assert.Error(t, err)
}
