// internal/shared/kernel/valueobject/resistencia_reactancia_test.go
package valueobject_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResistenciaReactancia_Valida(t *testing.T) {
	rr, err := valueobject.NewResistenciaReactancia(3.9, 0.164)
	require.NoError(t, err)
	assert.InDelta(t, 3.9, rr.R(), 0.0001)
	assert.InDelta(t, 0.164, rr.X(), 0.0001)
}

func TestNewResistenciaReactancia_CeroEsValido(t *testing.T) {
	// R=0 y X=0 son válidos (conductor ideal / línea sin pérdidas)
	rr, err := valueobject.NewResistenciaReactancia(0, 0)
	require.NoError(t, err)
	assert.Equal(t, 0.0, rr.R())
	assert.Equal(t, 0.0, rr.X())
}

func TestNewResistenciaReactancia_ResistenciaNegativa(t *testing.T) {
	_, err := valueobject.NewResistenciaReactancia(-1.0, 0.164)
	require.Error(t, err)
	assert.ErrorIs(t, err, valueobject.ErrImpedanciaInvalida)
}

func TestNewResistenciaReactancia_ReactanciaNegativa(t *testing.T) {
	_, err := valueobject.NewResistenciaReactancia(3.9, -0.5)
	require.Error(t, err)
	assert.ErrorIs(t, err, valueobject.ErrImpedanciaInvalida)
}

func TestNewResistenciaReactancia_AmbasNegativas(t *testing.T) {
	_, err := valueobject.NewResistenciaReactancia(-1.0, -0.5)
	require.Error(t, err)
	assert.ErrorIs(t, err, valueobject.ErrImpedanciaInvalida)
}

func TestResistenciaReactancia_AccesoViaGetters(t *testing.T) {
	// Los campos son privados; acceso via getters R() y X()
	rr, err := valueobject.NewResistenciaReactancia(0.62, 0.051)
	require.NoError(t, err)
	assert.Equal(t, 0.62, rr.R())
	assert.Equal(t, 0.051, rr.X())
}
