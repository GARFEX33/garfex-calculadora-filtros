// internal/domain/entity/filtro_activo_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFiltroActivo(t *testing.T) {
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, 125, 3)
	require.NoError(t, err)

	assert.Equal(t, "FA-001", fa.Clave)
	assert.Equal(t, entity.TipoFiltroActivo, fa.Tipo)
	assert.Equal(t, 480, fa.Voltaje)
	assert.Equal(t, 100, fa.Amperaje)
	assert.Equal(t, 125, fa.ITM)
	assert.Equal(t, 3, fa.Bornes)
}

func TestNewFiltroActivo_AmperajeInvalido(t *testing.T) {
	_, err := entity.NewFiltroActivo("FA-001", 480, 0, 125, 3)
	assert.Error(t, err)

	_, err = entity.NewFiltroActivo("FA-001", 480, -10, 125, 3)
	assert.Error(t, err)
}

func TestFiltroActivo_CalcularCorrienteNominal(t *testing.T) {
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, 125, 3)
	require.NoError(t, err)

	corriente, err := fa.CalcularCorrienteNominal()
	require.NoError(t, err)
	assert.InDelta(t, 100.0, corriente.Valor(), 0.001)
	assert.Equal(t, "A", corriente.Unidad())
}

func TestFiltroActivo_ImplementsCalculadorCorriente(t *testing.T) {
	fa, _ := entity.NewFiltroActivo("FA-001", 480, 100, 125, 3)
	var _ entity.CalculadorCorriente = fa
}

func TestFiltroActivo_ImplementsCalculadorPotencia(t *testing.T) {
	fa, _ := entity.NewFiltroActivo("FA-001", 480, 100, 125, 3)
	var _ entity.CalculadorPotencia = fa
}

func TestFiltroActivo_Potencias(t *testing.T) {
	// 100 A @ 480 V → kVA = 100 × 480 × √3 / 1000 = 83.138...
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, 125, 3)
	require.NoError(t, err)

	assert.InDelta(t, 83.138, fa.PotenciaKVA(), 0.01)
	assert.InDelta(t, 83.138, fa.PotenciaKW(), 0.01)  // PF=1
	assert.InDelta(t, 0.0, fa.PotenciaKVAR(), 0.001)  // PF=1
}
