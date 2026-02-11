// internal/domain/entity/transformador_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func itmTransf(t *testing.T, voltaje int) entity.ITM {
	t.Helper()
	itm, err := entity.NewITM(200, 3, 3, voltaje)
	require.NoError(t, err)
	return itm
}

func TestNewTransformador(t *testing.T) {
	tr, err := entity.NewTransformador("TR-001", 480, 500, itmTransf(t, 480))
	require.NoError(t, err)

	assert.Equal(t, "TR-001", tr.Clave)
	assert.Equal(t, entity.TipoEquipoTransformador, tr.Tipo)
	assert.Equal(t, 480, tr.Voltaje)
	assert.Equal(t, 500, tr.KVA)
	assert.Equal(t, 200, tr.ITM.Amperaje)
}

func TestNewTransformador_KVAInvalido(t *testing.T) {
	_, err := entity.NewTransformador("TR-001", 480, 0, itmTransf(t, 480))
	assert.Error(t, err)

	_, err = entity.NewTransformador("TR-001", 480, -100, itmTransf(t, 480))
	assert.Error(t, err)
}

func TestNewTransformador_VoltajeCero(t *testing.T) {
	itm, _ := entity.NewITM(200, 3, 3, 480)
	_, err := entity.NewTransformador("TR-BAD", 0, 500, itm)
	assert.Error(t, err)
}

func TestTransformador_CalcularCorrienteNominal(t *testing.T) {
	tests := []struct {
		name     string
		voltaje  int
		kva      int
		expected float64
	}{
		{
			name:    "500 KVA at 480V",
			voltaje: 480,
			kva:     500,
			// I = 500 / (0.48 × √3) = 500 / 0.83138... = 601.40 A
			expected: 601.40,
		},
		{
			name:    "150 KVA at 220V",
			voltaje: 220,
			kva:     150,
			// I = 150 / (0.22 × √3) = 150 / 0.38105... = 393.65 A
			expected: 393.65,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := entity.NewTransformador("TR-TEST", tt.voltaje, tt.kva, itmTransf(t, tt.voltaje))
			require.NoError(t, err)

			corriente, err := tr.CalcularCorrienteNominal()
			require.NoError(t, err)
			assert.InDelta(t, tt.expected, corriente.Valor(), 0.01)
		})
	}
}

func TestTransformador_ImplementsCalculadorCorriente(t *testing.T) {
	tr, _ := entity.NewTransformador("TR-001", 480, 500, itmTransf(t, 480))
	var _ entity.CalculadorCorriente = tr
}

func TestTransformador_ImplementsCalculadorPotencia(t *testing.T) {
	tr, _ := entity.NewTransformador("TR-001", 480, 500, itmTransf(t, 480))
	var _ entity.CalculadorPotencia = tr
}

func TestTransformador_Potencias(t *testing.T) {
	tr, err := entity.NewTransformador("TR-001", 480, 500, itmTransf(t, 480))
	require.NoError(t, err)

	assert.InDelta(t, 500.0, tr.PotenciaKVA(), 0.001)
	assert.InDelta(t, 0.0, tr.PotenciaKW(), 0.001)
	assert.InDelta(t, 0.0, tr.PotenciaKVAR(), 0.001)
}
