// internal/domain/entity/tipo_equipo_test.go
package entity_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTipoEquipo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected entity.TipoEquipo
		wantErr  bool
	}{
		{"FILTRO_ACTIVO valid", "FILTRO_ACTIVO", entity.TipoEquipoFiltroActivo, false},
		{"FILTRO_RECHAZO valid", "FILTRO_RECHAZO", entity.TipoEquipoFiltroRechazo, false},
		{"lowercase invalid", "filtro_activo", entity.TipoEquipo(""), true},
		{"empty invalid", "", entity.TipoEquipo(""), true},
		{"old ACTIVO invalid", "ACTIVO", entity.TipoEquipo(""), true},
		{"unknown invalid", "TABLERO", entity.TipoEquipo(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := entity.ParseTipoEquipo(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, entity.ErrTipoEquipoInvalido))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTipoEquipo_String(t *testing.T) {
	assert.Equal(t, "FILTRO_ACTIVO", entity.TipoEquipoFiltroActivo.String())
	assert.Equal(t, "FILTRO_RECHAZO", entity.TipoEquipoFiltroRechazo.String())
}
