// internal/calculos/domain/entity/canalizacion_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCanalizacion_valida(t *testing.T) {
	c, err := entity.NewCanalizacion(entity.TipoCanalizacionTuberiaPVC, "1 1/2", 150.0, 1)
	require.NoError(t, err)
	assert.Equal(t, entity.TipoCanalizacionTuberiaPVC, c.Tipo)
	assert.Equal(t, "1 1/2", c.Tamano)
	assert.InDelta(t, 150.0, c.AnchoRequerido, 0.001)
	assert.Equal(t, 1, c.NumeroDeTubos)
}

func TestNewCanalizacion_tipoInvalido(t *testing.T) {
	_, err := entity.NewCanalizacion("INVALIDO", "1 1/2", 150.0, 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, entity.ErrTipoCanalizacionInvalido)
}

func TestNewCanalizacion_tamanoVacio(t *testing.T) {
	_, err := entity.NewCanalizacion(entity.TipoCanalizacionTuberiaPVC, "", 150.0, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tamaño no puede estar vacío")
}

func TestNewCanalizacion_numeroDeTubosMenorA1(t *testing.T) {
	_, err := entity.NewCanalizacion(entity.TipoCanalizacionTuberiaPVC, "1 1/2", 150.0, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "numeroDeTubos debe ser >= 1")
}

func TestNewCanalizacion_charola(t *testing.T) {
	c, err := entity.NewCanalizacion(entity.TipoCanalizacionCharolaCableEspaciado, "300mm", 250.0, 1)
	require.NoError(t, err)
	assert.Equal(t, entity.TipoCanalizacionCharolaCableEspaciado, c.Tipo)
	assert.Equal(t, "300mm", c.Tamano)
}
