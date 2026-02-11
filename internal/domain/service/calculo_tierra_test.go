// internal/domain/service/calculo_tierra_test.go
package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// entradaTierra builds an EntradaTablaTierra with only the fields relevant
// for ground conductor selection (calibre, material, section).
// Ground conductors can be bare — no insulation or resistance data needed.
func entradaTierra(itmHasta int, calibre string, seccionMM2 float64) service.EntradaTablaTierra {
	return service.EntradaTablaTierra{
		ITMHasta: itmHasta,
		Conductor: valueobject.ConductorParams{
			Calibre:    calibre,
			Material:   "Cu",
			SeccionMM2: seccionMM2,
		},
	}
}

// Simplified NOM table 250-122 excerpt.
// "3 AWG" and "1 AWG" from original plan replaced with "4 AWG" and "2 AWG"
// because 3 AWG and 1 AWG are not valid calibres per NOM 310-15(b)(16).
var tablaTierraTest = []service.EntradaTablaTierra{
	entradaTierra(15, "14 AWG", 2.08),
	entradaTierra(20, "12 AWG", 3.31),
	entradaTierra(40, "10 AWG", 5.26),
	entradaTierra(60, "10 AWG", 5.26),
	entradaTierra(100, "8 AWG", 8.37),
	entradaTierra(200, "6 AWG", 13.30),
	entradaTierra(400, "4 AWG", 21.15),
	entradaTierra(600, "2 AWG", 33.62),
	entradaTierra(800, "1/0 AWG", 53.49),
	entradaTierra(1000, "2/0 AWG", 67.43),
}

func TestSeleccionarConductorTierra(t *testing.T) {
	tests := []struct {
		name            string
		itm             int
		expectedCalibre string
	}{
		{"ITM 15 → 14 AWG", 15, "14 AWG"},
		{"ITM 20 → 12 AWG", 20, "12 AWG"},
		{"ITM 30 → 10 AWG (≤40)", 30, "10 AWG"},
		{"ITM 100 → 8 AWG", 100, "8 AWG"},
		{"ITM 125 → 6 AWG (≤200)", 125, "6 AWG"},
		{"ITM 400 → 4 AWG", 400, "4 AWG"},
		{"ITM 600 → 2 AWG", 600, "2 AWG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conductor, err := service.SeleccionarConductorTierra(tt.itm, tablaTierraTest)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCalibre, conductor.Calibre())
			assert.Equal(t, "Cu", conductor.Material())
		})
	}
}

func TestSeleccionarConductorTierra_ITMExceedsTable(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(1200, tablaTierraTest)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrConductorNoEncontrado))
}

func TestSeleccionarConductorTierra_InvalidITM(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(0, tablaTierraTest)
	assert.Error(t, err)
}

func TestSeleccionarConductorTierra_EmptyTable(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(100, nil)
	assert.Error(t, err)
}
