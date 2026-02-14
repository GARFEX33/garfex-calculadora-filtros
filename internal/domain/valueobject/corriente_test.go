// internal/domain/valueobject/corriente_test.go
package valueobject_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCorriente(t *testing.T) {
	tests := []struct {
		name    string
		valor   float64
		wantErr bool
	}{
		{"positive value", 120.5, false},
		{"small positive value", 0.01, false},
		{"zero is invalid", 0, true},
		{"negative is invalid", -10.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := valueobject.NewCorriente(tt.valor)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, valueobject.ErrCorrienteInvalida))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.valor, c.Valor())
			assert.Equal(t, "A", c.Unidad())
		})
	}
}

func TestCorriente_Multiplicar(t *testing.T) {
	c, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	result, err := c.Multiplicar(1.25)
	require.NoError(t, err)
	assert.InDelta(t, 125.0, result.Valor(), 0.001)
}

func TestCorriente_Multiplicar_NegativeFactorFails(t *testing.T) {
	c, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	_, err = c.Multiplicar(-1)
	assert.Error(t, err)
}

func TestCorriente_Dividir(t *testing.T) {
	c, err := valueobject.NewCorriente(200)
	require.NoError(t, err)

	result, err := c.Dividir(2)
	require.NoError(t, err)
	assert.InDelta(t, 100.0, result.Valor(), 0.001)
}

func TestCorriente_Dividir_PorCeroFails(t *testing.T) {
	c, err := valueobject.NewCorriente(200)
	require.NoError(t, err)

	_, err = c.Dividir(0)
	assert.Error(t, err)
}
