// internal/calculos/domain/entity/filtro_activo_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// itmFA returns a standard 3-phase ITM for FiltroActivo tests.
func itmFA(t *testing.T) entity.ITM {
	t.Helper()
	itm, err := entity.NewITM(125, 3, 3, 480)
	require.NoError(t, err)
	return itm
}

func TestNewFiltroActivo(t *testing.T) {
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, itmFA(t))
	require.NoError(t, err)

	assert.Equal(t, "FA-001", fa.Clave)
	assert.Equal(t, entity.TipoEquipoFiltroActivo, fa.Tipo)
	assert.Equal(t, 480, fa.Voltaje)
	assert.Equal(t, 100, fa.AmperajeNominal)
	assert.Equal(t, 125, fa.ITM.Amperaje)
	assert.Equal(t, 3, fa.ITM.Bornes)
	assert.Equal(t, 3, fa.ITM.Polos)
}

func TestNewFiltroActivo_AmperajeInvalido(t *testing.T) {
	_, err := entity.NewFiltroActivo("FA-001", 480, 0, itmFA(t))
	assert.Error(t, err)

	_, err = entity.NewFiltroActivo("FA-001", 480, -10, itmFA(t))
	assert.Error(t, err)
}

func TestFiltroActivo_CalcularCorrienteNominal(t *testing.T) {
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, itmFA(t))
	require.NoError(t, err)

	corriente, err := fa.CalcularCorrienteNominal()
	require.NoError(t, err)
	assert.InDelta(t, 100.0, corriente.Valor(), 0.001)
	assert.Equal(t, "A", corriente.Unidad())
}

func TestFiltroActivo_ImplementsCalculadorCorriente(t *testing.T) {
	fa, _ := entity.NewFiltroActivo("FA-001", 480, 100, itmFA(t))
	var _ entity.CalculadorCorriente = fa
}

func TestFiltroActivo_ImplementsCalculadorPotencia(t *testing.T) {
	fa, _ := entity.NewFiltroActivo("FA-001", 480, 100, itmFA(t))
	var _ entity.CalculadorPotencia = fa
}

func TestFiltroActivo_Potencias(t *testing.T) {
	// 100 A @ 480 V → kVA = 100 × 480 × √3 / 1000 = 83.138...
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, itmFA(t))
	require.NoError(t, err)

	assert.InDelta(t, 83.138, fa.PotenciaKVA(), 0.01)
	assert.InDelta(t, 83.138, fa.PotenciaKW(), 0.01) // PF=1
	assert.InDelta(t, 0.0, fa.PotenciaKVAR(), 0.001) // PF=1
}
