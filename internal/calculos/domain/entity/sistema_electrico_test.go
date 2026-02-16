// internal/calculos/domain/entity/sistema_electrico_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSistemaElectrico_CantidadConductores(t *testing.T) {
	tests := []struct {
		sistema  entity.SistemaElectrico
		expected int
	}{
		{entity.SistemaElectricoDelta, 3},
		{entity.SistemaElectricoEstrella, 4},
		{entity.SistemaElectricoBifasico, 3},
		{entity.SistemaElectricoMonofasico, 2},
	}

	for _, tt := range tests {
		t.Run(string(tt.sistema), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.sistema.CantidadConductores())
		})
	}
}

func TestValidarSistemaElectrico(t *testing.T) {
	validos := []entity.SistemaElectrico{
		entity.SistemaElectricoDelta,
		entity.SistemaElectricoEstrella,
		entity.SistemaElectricoBifasico,
		entity.SistemaElectricoMonofasico,
	}
	for _, se := range validos {
		t.Run(string(se), func(t *testing.T) {
			err := entity.ValidarSistemaElectrico(se)
			assert.NoError(t, err)
		})
	}

	t.Run("invalido", func(t *testing.T) {
		err := entity.ValidarSistemaElectrico("TRIFASICO")
		require.Error(t, err)
		assert.ErrorIs(t, err, entity.ErrSistemaElectricoInvalido)
	})
}

func TestParseSistemaElectrico(t *testing.T) {
	tests := []struct {
		input    string
		expected entity.SistemaElectrico
		wantErr  bool
	}{
		{"DELTA", entity.SistemaElectricoDelta, false},
		{"ESTRELLA", entity.SistemaElectricoEstrella, false},
		{"BIFASICO", entity.SistemaElectricoBifasico, false},
		{"MONOFASICO", entity.SistemaElectricoMonofasico, false},
		{"INVALIDO", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := entity.ParseSistemaElectrico(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
