// internal/calculos/domain/entity/carga_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func itmCarga(t *testing.T, voltaje int) entity.ITM {
	t.Helper()
	itm, err := entity.NewITM(100, 3, 3, voltaje)
	require.NoError(t, err)
	return itm
}

func TestNewCarga(t *testing.T) {
	c, err := entity.NewCarga("C-001", 480, 50, 0.85, 3, itmCarga(t, 480))
	require.NoError(t, err)

	assert.Equal(t, "C-001", c.Clave)
	assert.Equal(t, entity.TipoEquipoCarga, c.Tipo)
	assert.Equal(t, 480, c.Voltaje)
	assert.Equal(t, 50, c.KW)
	assert.InDelta(t, 0.85, c.FactorPotencia, 0.001)
	assert.Equal(t, 3, c.Fases)
}

func TestNewCarga_Invalido(t *testing.T) {
	tests := []struct {
		name    string
		voltaje int
		kw      int
		fp      float64
		fases   int
	}{
		{"kw cero", 480, 0, 0.85, 3},
		{"kw negativo", 480, -10, 0.85, 3},
		{"fp cero", 480, 50, 0, 3},
		{"fp negativo", 480, 50, -0.5, 3},
		{"fp mayor a 1", 480, 50, 1.1, 3},
		{"fases cero", 480, 50, 0.85, 0},
		{"fases 4", 480, 50, 0.85, 4},
		{"voltaje cero", 0, 50, 0.85, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			itm, _ := entity.NewITM(100, 3, 3, 480)
			_, err := entity.NewCarga("C-BAD", tt.voltaje, tt.kw, tt.fp, tt.fases, itm)
			assert.Error(t, err)
		})
	}
}

func TestCarga_CalcularCorrienteNominal_Trifasico(t *testing.T) {
	// I = 50 / (0.48 × √3 × 0.85) = 50 / 0.70668... = 70.76 A
	c, err := entity.NewCarga("C-001", 480, 50, 0.85, 3, itmCarga(t, 480))
	require.NoError(t, err)

	corriente, err := c.CalcularCorrienteNominal()
	require.NoError(t, err)
	assert.InDelta(t, 70.76, corriente.Valor(), 0.01)
}

func TestCarga_CalcularCorrienteNominal_Bifasico(t *testing.T) {
	// I = 50 / (0.48 × 2 × 0.85) = 50 / 0.816 = 61.27 A
	c, err := entity.NewCarga("C-002", 480, 50, 0.85, 2, itmCarga(t, 480))
	require.NoError(t, err)

	corriente, err := c.CalcularCorrienteNominal()
	require.NoError(t, err)
	assert.InDelta(t, 61.27, corriente.Valor(), 0.01)
}

func TestCarga_CalcularCorrienteNominal_Monofasico(t *testing.T) {
	// I = 50 / (0.48 × 1 × 0.85) = 50 / 0.408 = 122.55 A
	c, err := entity.NewCarga("C-003", 480, 50, 0.85, 1, itmCarga(t, 480))
	require.NoError(t, err)

	corriente, err := c.CalcularCorrienteNominal()
	require.NoError(t, err)
	assert.InDelta(t, 122.55, corriente.Valor(), 0.01)
}

func TestCarga_ImplementsCalculadorCorriente(t *testing.T) {
	c, _ := entity.NewCarga("C-001", 480, 50, 0.85, 3, itmCarga(t, 480))
	var _ entity.CalculadorCorriente = c
}

func TestCarga_ImplementsCalculadorPotencia(t *testing.T) {
	c, _ := entity.NewCarga("C-001", 480, 50, 0.85, 3, itmCarga(t, 480))
	var _ entity.CalculadorPotencia = c
}

func TestCarga_Potencias(t *testing.T) {
	// 50 kW, FP=0.85
	// KVA = 50 / 0.85 = 58.8235...
	// KVAR = √(58.8235² - 50²) = √(3460.21 - 2500) = √960.21 = 30.987...
	c, err := entity.NewCarga("C-001", 480, 50, 0.85, 3, itmCarga(t, 480))
	require.NoError(t, err)

	assert.InDelta(t, 58.82, c.PotenciaKVA(), 0.01)
	assert.InDelta(t, 50.0, c.PotenciaKW(), 0.001)
	assert.InDelta(t, 30.99, c.PotenciaKVAR(), 0.01)
}

func TestCarga_Potencias_FP1(t *testing.T) {
	// FP=1: KVA=KW=50, KVAR=0
	c, err := entity.NewCarga("C-002", 480, 50, 1.0, 3, itmCarga(t, 480))
	require.NoError(t, err)

	assert.InDelta(t, 50.0, c.PotenciaKVA(), 0.001)
	assert.InDelta(t, 50.0, c.PotenciaKW(), 0.001)
	assert.InDelta(t, 0.0, c.PotenciaKVAR(), 0.001)
}
