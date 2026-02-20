// internal/calculos/domain/entity/tipo_voltaje_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTipoVoltaje_Validos(t *testing.T) {
	tests := []struct {
		input    string
		expected entity.TipoVoltaje
	}{
		{"FASE_NEUTRO", entity.TipoVoltajeFaseNeutro},
		{"fase_neutro", entity.TipoVoltajeFaseNeutro},
		{"FN", entity.TipoVoltajeFaseNeutro},
		{"fn", entity.TipoVoltajeFaseNeutro},
		{"FASE_FASE", entity.TipoVoltajeFaseFase},
		{"fase_fase", entity.TipoVoltajeFaseFase},
		{"FF", entity.TipoVoltajeFaseFase},
		{"ff", entity.TipoVoltajeFaseFase},
		{" FASE_NEUTRO ", entity.TipoVoltajeFaseNeutro}, // con espacios
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			resultado, err := entity.ParseTipoVoltaje(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, resultado)
		})
	}
}

func TestParseTipoVoltaje_Invalidos(t *testing.T) {
	tests := []string{
		"",
		"INVALIDO",
		"TRIFASICO",
		"127V",
		"NEUTRO",
		"FASE",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := entity.ParseTipoVoltaje(input)
			require.Error(t, err)
			assert.ErrorIs(t, err, entity.ErrTipoVoltajeInvalido)
		})
	}
}

func TestTipoVoltaje_String(t *testing.T) {
	assert.Equal(t, "FASE_NEUTRO", entity.TipoVoltajeFaseNeutro.String())
	assert.Equal(t, "FASE_FASE", entity.TipoVoltajeFaseFase.String())
}

func TestTipoVoltaje_EsFaseNeutro(t *testing.T) {
	assert.True(t, entity.TipoVoltajeFaseNeutro.EsFaseNeutro())
	assert.False(t, entity.TipoVoltajeFaseFase.EsFaseNeutro())
}

func TestTipoVoltaje_EsFaseFase(t *testing.T) {
	assert.True(t, entity.TipoVoltajeFaseFase.EsFaseFase())
	assert.False(t, entity.TipoVoltajeFaseNeutro.EsFaseFase())
}
