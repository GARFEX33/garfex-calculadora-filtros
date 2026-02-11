// internal/domain/service/ajuste_corriente_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAjustarCorriente_SingleFactor(t *testing.T) {
	cn, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	factores := map[string]float64{
		"proteccion": 1.25,
	}

	result, err := service.AjustarCorriente(cn, factores)
	require.NoError(t, err)
	assert.InDelta(t, 125.0, result.Valor(), 0.001)
}

func TestAjustarCorriente_MultipleFactors(t *testing.T) {
	cn, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	factores := map[string]float64{
		"proteccion":   1.25,
		"temperatura":  0.88,
		"agrupamiento": 0.80,
	}

	// 100 × 1.25 × 0.88 × 0.80 = 88.0
	result, err := service.AjustarCorriente(cn, factores)
	require.NoError(t, err)
	assert.InDelta(t, 88.0, result.Valor(), 0.01)
}

func TestAjustarCorriente_EmptyFactors(t *testing.T) {
	cn, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	result, err := service.AjustarCorriente(cn, nil)
	require.NoError(t, err)
	assert.InDelta(t, 100.0, result.Valor(), 0.001)
}

func TestAjustarCorriente_ZeroFactorFails(t *testing.T) {
	cn, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	factores := map[string]float64{
		"proteccion": 0,
	}

	_, err = service.AjustarCorriente(cn, factores)
	assert.Error(t, err)
}
