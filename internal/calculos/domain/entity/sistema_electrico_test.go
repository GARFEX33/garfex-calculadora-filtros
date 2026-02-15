// internal/calculos/domain/entity/sistema_electrico_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
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
