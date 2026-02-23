// internal/equipos/domain/entity/tipo_filtro_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/equipos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTipoFiltro(t *testing.T) {
	t.Run("valores válidos", func(t *testing.T) {
		cases := []struct {
			input    string
			expected entity.TipoFiltro
		}{
			{"A", entity.TipoFiltroA},
			{"KVA", entity.TipoFiltroKVA},
			{"KVAR", entity.TipoFiltroKVAR},
		}
		for _, tc := range cases {
			got, err := entity.ParseTipoFiltro(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		}
	})

	t.Run("valor inválido devuelve error", func(t *testing.T) {
		_, err := entity.ParseTipoFiltro("FILTRO_ACTIVO")
		assert.ErrorIs(t, err, entity.ErrTipoFiltroInvalido)
	})

	t.Run("string vacío devuelve error", func(t *testing.T) {
		_, err := entity.ParseTipoFiltro("")
		assert.ErrorIs(t, err, entity.ErrTipoFiltroInvalido)
	})
}

func TestTipoFiltro_String(t *testing.T) {
	assert.Equal(t, "A", entity.TipoFiltroA.String())
	assert.Equal(t, "KVA", entity.TipoFiltroKVA.String())
	assert.Equal(t, "KVAR", entity.TipoFiltroKVAR.String())
}
