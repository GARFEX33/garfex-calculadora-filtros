// internal/calculos/domain/service/seleccionar_temperatura_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeleccionarTemperatura_CorrienteBajaEnTuberia(t *testing.T) {
	// Corriente <= 100A en tubería → 60°C
	corriente, err := valueobject.NewCorriente(50.0)
	require.NoError(t, err)

	temp := service.SeleccionarTemperatura(corriente, entity.TipoCanalizacionTuberiaPVC, nil)

	assert.Equal(t, valueobject.Temp60, temp)
}

func TestSeleccionarTemperatura_CorrienteAltaEnTuberia(t *testing.T) {
	// Corriente > 100A en tubería → 75°C
	corriente, err := valueobject.NewCorriente(150.0)
	require.NoError(t, err)

	temp := service.SeleccionarTemperatura(corriente, entity.TipoCanalizacionTuberiaPVC, nil)

	assert.Equal(t, valueobject.Temp75, temp)
}

func TestSeleccionarTemperatura_CorrienteBajaEnCharolaTriangular(t *testing.T) {
	// Corriente <= 100A en charola triangular → 75°C (no tiene columna 60°C)
	corriente, err := valueobject.NewCorriente(80.0)
	require.NoError(t, err)

	temp := service.SeleccionarTemperatura(corriente, entity.TipoCanalizacionCharolaCableTriangular, nil)

	assert.Equal(t, valueobject.Temp75, temp)
}

func TestSeleccionarTemperatura_CorrienteBajaEnCharolaEspaciado(t *testing.T) {
	// Corriente <= 100A en charola espaciado → 60°C (tiene columna 60°C)
	corriente, err := valueobject.NewCorriente(90.0)
	require.NoError(t, err)

	temp := service.SeleccionarTemperatura(corriente, entity.TipoCanalizacionCharolaCableEspaciado, nil)

	assert.Equal(t, valueobject.Temp60, temp)
}

func TestSeleccionarTemperatura_OverrideTomaPrecedencia(t *testing.T) {
	// Override explícito → usa el override sin importar corriente ni canalización
	corriente, err := valueobject.NewCorriente(50.0)
	require.NoError(t, err)

	override := valueobject.Temp90
	temp := service.SeleccionarTemperatura(corriente, entity.TipoCanalizacionTuberiaPVC, &override)

	assert.Equal(t, valueobject.Temp90, temp)
}

func TestSeleccionarTemperatura_CorrienteExacta100A(t *testing.T) {
	// Corriente exactamente 100A → <=100A → 60°C en tubería
	corriente, err := valueobject.NewCorriente(100.0)
	require.NoError(t, err)

	temp := service.SeleccionarTemperatura(corriente, entity.TipoCanalizacionTuberiaPVC, nil)

	assert.Equal(t, valueobject.Temp60, temp)
}

func TestSeleccionarTemperatura_CorrienteExacta100AEnCharolaTriangular(t *testing.T) {
	// Corriente exactamente 100A en charola triangular → 75°C
	corriente, err := valueobject.NewCorriente(100.0)
	require.NoError(t, err)

	temp := service.SeleccionarTemperatura(corriente, entity.TipoCanalizacionCharolaCableTriangular, nil)

	assert.Equal(t, valueobject.Temp75, temp)
}
