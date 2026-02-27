// internal/equipos/domain/entity/equipo_filtro_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEquipoFiltro(t *testing.T) {
	clave := "FA-480-100"
	bornes := 6
	conexion := entity.ConexionEstrella
	tipoVoltaje := entity.TipoVoltajeFaseFase

	t.Run("crea equipo válido con conexion estrella", func(t *testing.T) {
		eq, err := entity.NewEquipoFiltro(&clave, entity.TipoFiltroA, 480, 100, 125, &bornes, &conexion, &tipoVoltaje)
		require.NoError(t, err)
		assert.Equal(t, &clave, eq.Clave)
		assert.Equal(t, entity.TipoFiltroA, eq.Tipo)
		assert.Equal(t, 480, eq.Voltaje)
		assert.Equal(t, 100, eq.Amperaje)
		assert.Equal(t, 125, eq.ITM)
		assert.Equal(t, &bornes, eq.Bornes)
		assert.Equal(t, &conexion, eq.Conexion)
		assert.Equal(t, &tipoVoltaje, eq.TipoVoltaje)
	})

	t.Run("crea equipo sin clave, bornes ni conexion (nullables)", func(t *testing.T) {
		eq, err := entity.NewEquipoFiltro(nil, entity.TipoFiltroKVAR, 220, 50, 60, nil, nil, nil)
		require.NoError(t, err)
		assert.Nil(t, eq.Clave)
		assert.Nil(t, eq.Bornes)
		assert.Nil(t, eq.Conexion)
		assert.Nil(t, eq.TipoVoltaje)
	})

	t.Run("rechaza voltaje cero", func(t *testing.T) {
		_, err := entity.NewEquipoFiltro(nil, entity.TipoFiltroA, 0, 100, 125, nil, nil, nil)
		assert.ErrorIs(t, err, entity.ErrVoltajeInvalido)
	})

	t.Run("rechaza voltaje negativo", func(t *testing.T) {
		_, err := entity.NewEquipoFiltro(nil, entity.TipoFiltroA, -1, 100, 125, nil, nil, nil)
		assert.ErrorIs(t, err, entity.ErrVoltajeInvalido)
	})

	t.Run("rechaza amperaje cero", func(t *testing.T) {
		_, err := entity.NewEquipoFiltro(nil, entity.TipoFiltroA, 480, 0, 125, nil, nil, nil)
		assert.ErrorIs(t, err, entity.ErrAmperajeInvalido)
	})

	t.Run("rechaza ITM cero", func(t *testing.T) {
		_, err := entity.NewEquipoFiltro(nil, entity.TipoFiltroA, 480, 100, 0, nil, nil, nil)
		assert.ErrorIs(t, err, entity.ErrITMInvalido)
	})
}
