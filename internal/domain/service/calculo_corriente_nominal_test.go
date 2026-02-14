// internal/domain/service/calculo_corriente_nominal_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func itmHelper(t *testing.T, amperaje, voltaje int) entity.ITM {
	t.Helper()
	itm, err := entity.NewITM(amperaje, 3, 3, voltaje)
	require.NoError(t, err)
	return itm
}

func TestCalcularCorrienteNominal_FiltroActivo(t *testing.T) {
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, itmHelper(t, 125, 480))
	require.NoError(t, err)

	corriente, err := service.CalcularCorrienteNominal(fa)
	require.NoError(t, err)
	assert.InDelta(t, 100.0, corriente.Valor(), 0.001)
}

func TestCalcularCorrienteNominal_FiltroRechazo(t *testing.T) {
	fr, err := entity.NewFiltroRechazo("FR-001", 480, 100, itmHelper(t, 125, 480))
	require.NoError(t, err)

	corriente, err := service.CalcularCorrienteNominal(fr)
	require.NoError(t, err)
	// I = 100 / (0.48 × √3) ≈ 120.28 A
	assert.InDelta(t, 120.28, corriente.Valor(), 0.01)
}

func TestCalcularCorrienteNominal_Transformador(t *testing.T) {
	tr, err := entity.NewTransformador("TR-001", 480, 500, itmHelper(t, 200, 480))
	require.NoError(t, err)

	corriente, err := service.CalcularCorrienteNominal(tr)
	require.NoError(t, err)
	// I = 500 / (0.48 × √3) ≈ 601.40 A
	assert.InDelta(t, 601.40, corriente.Valor(), 0.01)
}

func TestCalcularCorrienteNominal_Carga(t *testing.T) {
	c, err := entity.NewCarga("C-001", 480, 50, 0.85, 3, itmHelper(t, 100, 480))
	require.NoError(t, err)

	corriente, err := service.CalcularCorrienteNominal(c)
	require.NoError(t, err)
	// I = 50 / (0.48 × √3 × 0.85) ≈ 70.76 A
	assert.InDelta(t, 70.76, corriente.Valor(), 0.01)
}
