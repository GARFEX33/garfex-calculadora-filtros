// internal/domain/entity/filtro_rechazo_test.go
package entity_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// itmFR returns a standard 3-phase ITM for FiltroRechazo tests.
func itmFR(t *testing.T, voltaje int) entity.ITM {
	t.Helper()
	itm, err := entity.NewITM(125, 3, 3, voltaje)
	require.NoError(t, err)
	return itm
}

func TestNewFiltroRechazo(t *testing.T) {
	fr, err := entity.NewFiltroRechazo("FR-001", 480, 100, itmFR(t, 480))
	require.NoError(t, err)

	assert.Equal(t, "FR-001", fr.Clave)
	assert.Equal(t, entity.TipoFiltroRechazo, fr.Tipo)
	assert.Equal(t, 480, fr.Voltaje)
	assert.Equal(t, 100, fr.KVAR)
	assert.Equal(t, 125, fr.ITM.Amperaje)
	assert.Equal(t, 3, fr.ITM.Bornes)
}

func TestNewFiltroRechazo_KVARInvalido(t *testing.T) {
	_, err := entity.NewFiltroRechazo("FR-001", 480, 0, itmFR(t, 480))
	assert.Error(t, err)

	_, err = entity.NewFiltroRechazo("FR-001", 480, -50, itmFR(t, 480))
	assert.Error(t, err)
}

func TestFiltroRechazo_CalcularCorrienteNominal(t *testing.T) {
	tests := []struct {
		name     string
		voltaje  int
		kvar     int
		expected float64
	}{
		{
			name:    "100 KVAR at 480V",
			voltaje: 480,
			kvar:    100,
			// I = 100 / (0.48 × √3) = 100 / 0.83138... = 120.28 A
			expected: 120.28,
		},
		{
			name:    "50 KVAR at 220V",
			voltaje: 220,
			kvar:    50,
			// I = 50 / (0.22 × √3) = 50 / 0.38105... = 131.22 A
			expected: 131.22,
		},
		{
			name:    "200 KVAR at 440V",
			voltaje: 440,
			kvar:    200,
			// I = 200 / (0.44 × √3) = 200 / 0.76210... = 262.43 A
			expected: 262.43,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr, err := entity.NewFiltroRechazo("FR-TEST", tt.voltaje, tt.kvar, itmFR(t, tt.voltaje))
			require.NoError(t, err)

			corriente, err := fr.CalcularCorrienteNominal()
			require.NoError(t, err)
			assert.InDelta(t, tt.expected, corriente.Valor(), 0.01)
		})
	}
}

func TestFiltroRechazo_CalcularCorrienteNominal_VoltajeCero(t *testing.T) {
	// Voltage 0 would cause division by zero.
	// NewFiltroRechazo should reject voltaje <= 0.
	itm, _ := entity.NewITM(125, 3, 3, 480)
	_, err := entity.NewFiltroRechazo("FR-BAD", 0, 100, itm)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, entity.ErrDivisionPorCero))
}

func TestFiltroRechazo_ImplementsCalculadorCorriente(t *testing.T) {
	fr, _ := entity.NewFiltroRechazo("FR-001", 480, 100, itmFR(t, 480))
	var _ entity.CalculadorCorriente = fr
}

func TestFiltroRechazo_ImplementsCalculadorPotencia(t *testing.T) {
	fr, _ := entity.NewFiltroRechazo("FR-001", 480, 100, itmFR(t, 480))
	var _ entity.CalculadorPotencia = fr
}

func TestFiltroRechazo_Potencias(t *testing.T) {
	// 100 kVAR purely reactive: kVAR=100, kVA=100, kW=0
	fr, err := entity.NewFiltroRechazo("FR-001", 480, 100, itmFR(t, 480))
	require.NoError(t, err)

	assert.InDelta(t, 100.0, fr.PotenciaKVAR(), 0.001)
	assert.InDelta(t, 100.0, fr.PotenciaKVA(), 0.001) // purely reactive
	assert.InDelta(t, 0.0, fr.PotenciaKW(), 0.001)    // purely reactive
}
