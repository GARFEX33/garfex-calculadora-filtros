// internal/calculos/domain/service/calculo_conductor_test.go
package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// entradaConductor builds an EntradaTablaConductor with the fields relevant
// for conductor selection. Optional fields (resistance, reactance) omitted —
// they are validated at point of use (voltage drop, etc.), not at construction.
func entradaConductor(calibre string, capacidad, seccionMM2 float64) valueobject.EntradaTablaConductor {
	return valueobject.EntradaTablaConductor{
		Capacidad: capacidad,
		Conductor: valueobject.ConductorParams{
			Calibre:         calibre,
			Material:        valueobject.MaterialCobre,
			TipoAislamiento: "THHN",
			SeccionMM2:      seccionMM2,
		},
	}
}

// Simplified NOM table 310-15(b)(16) excerpt for Cu THHN 90°C
var tablaConductorTest = []valueobject.EntradaTablaConductor{
	entradaConductor("14 AWG", 25, 2.08),
	entradaConductor("12 AWG", 30, 3.31),
	entradaConductor("10 AWG", 40, 5.26),
	entradaConductor("8 AWG", 55, 8.37),
	entradaConductor("6 AWG", 75, 13.30),
	entradaConductor("4 AWG", 95, 21.15),
	entradaConductor("2 AWG", 130, 33.62),
	entradaConductor("1/0 AWG", 170, 53.49),
	entradaConductor("4/0 AWG", 260, 107.2),
	entradaConductor("500 MCM", 380, 253.4),
}

func TestSeleccionarConductorAlimentacion_Simple(t *testing.T) {
	corriente, err := valueobject.NewCorriente(120)
	require.NoError(t, err)

	conductor, err := service.SeleccionarConductorAlimentacion(corriente, 1, tablaConductorTest)
	require.NoError(t, err)
	// 120A needs at least 130A capacity → 2 AWG
	assert.Equal(t, "2 AWG", conductor.Calibre())
}

func TestSeleccionarConductorAlimentacion_ExactMatch(t *testing.T) {
	corriente, err := valueobject.NewCorriente(95)
	require.NoError(t, err)

	conductor, err := service.SeleccionarConductorAlimentacion(corriente, 1, tablaConductorTest)
	require.NoError(t, err)
	// 95A exactly matches 4 AWG capacity
	assert.Equal(t, "4 AWG", conductor.Calibre())
}

func TestSeleccionarConductorAlimentacion_ConHilosPorFase(t *testing.T) {
	corriente, err := valueobject.NewCorriente(240)
	require.NoError(t, err)

	// 240A / 2 hilos = 120A per wire → needs 130A → 2 AWG
	conductor, err := service.SeleccionarConductorAlimentacion(corriente, 2, tablaConductorTest)
	require.NoError(t, err)
	assert.Equal(t, "2 AWG", conductor.Calibre())
}

func TestSeleccionarConductorAlimentacion_CorrienteExceedsAllCapacities(t *testing.T) {
	corriente, err := valueobject.NewCorriente(500)
	require.NoError(t, err)

	_, err = service.SeleccionarConductorAlimentacion(corriente, 1, tablaConductorTest)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrConductorNoEncontrado))
}

func TestSeleccionarConductorAlimentacion_EmptyTable(t *testing.T) {
	corriente, err := valueobject.NewCorriente(10)
	require.NoError(t, err)

	_, err = service.SeleccionarConductorAlimentacion(corriente, 1, nil)
	assert.Error(t, err)
}
