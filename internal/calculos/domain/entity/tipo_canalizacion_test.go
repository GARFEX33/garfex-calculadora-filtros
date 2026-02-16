// internal/calculos/domain/entity/tipo_canalizacion_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTipoCanalizacion_validos(t *testing.T) {
	casos := []struct {
		input    string
		expected entity.TipoCanalizacion
	}{
		{"TUBERIA_PVC", entity.TipoCanalizacionTuberiaPVC},
		{"TUBERIA_ALUMINIO", entity.TipoCanalizacionTuberiaAluminio},
		{"TUBERIA_ACERO_PG", entity.TipoCanalizacionTuberiaAceroPG},
		{"TUBERIA_ACERO_PD", entity.TipoCanalizacionTuberiaAceroPD},
		{"CHAROLA_CABLE_ESPACIADO", entity.TipoCanalizacionCharolaCableEspaciado},
		{"CHAROLA_CABLE_TRIANGULAR", entity.TipoCanalizacionCharolaCableTriangular},
	}
	for _, c := range casos {
		t.Run(c.input, func(t *testing.T) {
			got, err := entity.ParseTipoCanalizacion(c.input)
			require.NoError(t, err)
			assert.Equal(t, c.expected, got)
		})
	}
}

func TestParseTipoCanalizacion_invalido(t *testing.T) {
	_, err := entity.ParseTipoCanalizacion("TUBERIA_MADERA")
	require.Error(t, err)
	assert.ErrorIs(t, err, entity.ErrTipoCanalizacionInvalido)
}

func TestTipoCanalizacion_EsCharola(t *testing.T) {
	charolas := []entity.TipoCanalizacion{
		entity.TipoCanalizacionCharolaCableEspaciado,
		entity.TipoCanalizacionCharolaCableTriangular,
	}
	for _, tc := range charolas {
		t.Run(string(tc), func(t *testing.T) {
			assert.True(t, tc.EsCharola())
		})
	}

	tuberias := []entity.TipoCanalizacion{
		entity.TipoCanalizacionTuberiaPVC,
		entity.TipoCanalizacionTuberiaAluminio,
		entity.TipoCanalizacionTuberiaAceroPG,
		entity.TipoCanalizacionTuberiaAceroPD,
	}
	for _, tc := range tuberias {
		t.Run(string(tc), func(t *testing.T) {
			assert.False(t, tc.EsCharola())
		})
	}
}

func TestValidarTipoCanalizacion_validos(t *testing.T) {
	casos := []struct {
		nombre string
		tc     entity.TipoCanalizacion
	}{
		{"tuberia PVC", entity.TipoCanalizacionTuberiaPVC},
		{"tuberia aluminio", entity.TipoCanalizacionTuberiaAluminio},
		{"tuberia acero pared gruesa", entity.TipoCanalizacionTuberiaAceroPG},
		{"tuberia acero pared delgada", entity.TipoCanalizacionTuberiaAceroPD},
		{"charola cable espaciado", entity.TipoCanalizacionCharolaCableEspaciado},
		{"charola cable triangular", entity.TipoCanalizacionCharolaCableTriangular},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			err := entity.ValidarTipoCanalizacion(c.tc)
			assert.NoError(t, err)
		})
	}
}

func TestValidarTipoCanalizacion_invalido(t *testing.T) {
	err := entity.ValidarTipoCanalizacion("INVALIDO")
	require.Error(t, err)
	assert.ErrorIs(t, err, entity.ErrTipoCanalizacionInvalido)
}

func TestTipoCanalizacion_valoresString(t *testing.T) {
	assert.Equal(t, entity.TipoCanalizacion("TUBERIA_PVC"), entity.TipoCanalizacionTuberiaPVC)
	assert.Equal(t, entity.TipoCanalizacion("TUBERIA_ALUMINIO"), entity.TipoCanalizacionTuberiaAluminio)
	assert.Equal(t, entity.TipoCanalizacion("TUBERIA_ACERO_PG"), entity.TipoCanalizacionTuberiaAceroPG)
	assert.Equal(t, entity.TipoCanalizacion("TUBERIA_ACERO_PD"), entity.TipoCanalizacionTuberiaAceroPD)
	assert.Equal(t, entity.TipoCanalizacion("CHAROLA_CABLE_ESPACIADO"), entity.TipoCanalizacionCharolaCableEspaciado)
	assert.Equal(t, entity.TipoCanalizacion("CHAROLA_CABLE_TRIANGULAR"), entity.TipoCanalizacionCharolaCableTriangular)
}
