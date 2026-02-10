// internal/domain/entity/tipo_filtro_test.go
package entity_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTipoFiltro(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected entity.TipoFiltro
		wantErr  bool
	}{
		{"ACTIVO valid", "ACTIVO", entity.TipoFiltroActivo, false},
		{"RECHAZO valid", "RECHAZO", entity.TipoFiltroRechazo, false},
		{"lowercase invalid", "activo", entity.TipoFiltro(""), true},
		{"empty invalid", "", entity.TipoFiltro(""), true},
		{"unknown invalid", "PASIVO", entity.TipoFiltro(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := entity.ParseTipoFiltro(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, entity.ErrTipoFiltroInvalido))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTipoFiltro_String(t *testing.T) {
	assert.Equal(t, "ACTIVO", entity.TipoFiltroActivo.String())
	assert.Equal(t, "RECHAZO", entity.TipoFiltroRechazo.String())
}
