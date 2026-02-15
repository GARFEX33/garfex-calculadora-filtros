// internal/shared/kernel/valueobject/charola_test.go
package valueobject_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- ConductorCharola ---

func TestNewConductorCharola_Valido(t *testing.T) {
	c, err := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 25.48})
	require.NoError(t, err)
	assert.InDelta(t, 25.48, c.DiametroMM(), 0.001)
}

func TestNewConductorCharola_DiametroCero(t *testing.T) {
	_, err := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: 0})
	assert.ErrorIs(t, err, valueobject.ErrConductorInvalido)
}

func TestNewConductorCharola_DiametroNegativo(t *testing.T) {
	_, err := valueobject.NewConductorCharola(valueobject.ConductorCharolaParams{DiametroMM: -1.5})
	assert.ErrorIs(t, err, valueobject.ErrConductorInvalido)
}

// --- CableControl ---

func TestNewCableControl_Valido(t *testing.T) {
	c, err := valueobject.NewCableControl(valueobject.CableControlParams{Cantidad: 3, DiametroMM: 8.5})
	require.NoError(t, err)
	assert.Equal(t, 3, c.Cantidad())
	assert.InDelta(t, 8.5, c.DiametroMM(), 0.001)
}

func TestNewCableControl_CantidadCero_Permitida(t *testing.T) {
	// cantidad 0 es válida: representa "sin cables de control"
	_, err := valueobject.NewCableControl(valueobject.CableControlParams{Cantidad: 0, DiametroMM: 8.5})
	require.NoError(t, err)
}

func TestNewCableControl_CantidadNegativa(t *testing.T) {
	_, err := valueobject.NewCableControl(valueobject.CableControlParams{Cantidad: -1, DiametroMM: 8.5})
	assert.ErrorIs(t, err, valueobject.ErrCableControlInvalido)
}

func TestNewCableControl_DiametroCero(t *testing.T) {
	_, err := valueobject.NewCableControl(valueobject.CableControlParams{Cantidad: 2, DiametroMM: 0})
	assert.ErrorIs(t, err, valueobject.ErrCableControlInvalido)
}

func TestNewCableControl_DiametroNegativo(t *testing.T) {
	_, err := valueobject.NewCableControl(valueobject.CableControlParams{Cantidad: 2, DiametroMM: -3.0})
	assert.ErrorIs(t, err, valueobject.ErrCableControlInvalido)
}
