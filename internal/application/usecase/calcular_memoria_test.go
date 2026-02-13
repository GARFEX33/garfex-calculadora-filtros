// internal/application/usecase/calcular_memoria_test.go
package usecase

import (
	"context"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTablaNOMRepository es un mock simple para tests
type mockTablaNOMRepository struct {
	tablaTierra       []valueobject.EntradaTablaTierra
	tablaAmpacidad    []valueobject.EntradaTablaConductor
	tablaCanalizacion []valueobject.EntradaTablaCanalizacion
	impedancia        valueobject.ResistenciaReactancia
}

func (m *mockTablaNOMRepository) ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error) {
	return m.tablaTierra, nil
}

func (m *mockTablaNOMRepository) ObtenerTablaAmpacidad(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
) ([]valueobject.EntradaTablaConductor, error) {
	return m.tablaAmpacidad, nil
}

func (m *mockTablaNOMRepository) ObtenerImpedancia(
	ctx context.Context,
	calibre string,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
) (valueobject.ResistenciaReactancia, error) {
	return m.impedancia, nil
}

func (m *mockTablaNOMRepository) ObtenerTablaCanalizacion(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
) ([]valueobject.EntradaTablaCanalizacion, error) {
	return m.tablaCanalizacion, nil
}

// mockEquipoRepository es un mock simple para tests
type mockEquipoRepository struct {
	equipo entity.CalculadorCorriente
	err    error
}

func (m *mockEquipoRepository) BuscarPorClave(ctx context.Context, clave string) (entity.CalculadorCorriente, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.equipo, nil
}

func TestCalcularMemoriaUseCase_Execute_ManualAmperaje(t *testing.T) {
	// Setup mocks
	tablaTierra := []valueobject.EntradaTablaTierra{
		{ITMHasta: 100, Conductor: valueobject.ConductorParams{Calibre: "8 AWG", Material: "Cu", SeccionMM2: 8.37}},
		{ITMHasta: 200, Conductor: valueobject.ConductorParams{Calibre: "6 AWG", Material: "Cu", SeccionMM2: 13.3}},
	}

	tablaAmpacidad := []valueobject.EntradaTablaConductor{
		{Capacidad: 30, Conductor: valueobject.ConductorParams{Calibre: "10 AWG", Material: "Cu", SeccionMM2: 5.26}},
		{Capacidad: 55, Conductor: valueobject.ConductorParams{Calibre: "8 AWG", Material: "Cu", SeccionMM2: 8.37}},
		{Capacidad: 75, Conductor: valueobject.ConductorParams{Calibre: "6 AWG", Material: "Cu", SeccionMM2: 13.3}},
	}

	tablaCanalizacion := []valueobject.EntradaTablaCanalizacion{
		{Tamano: "1/2", AreaInteriorMM2: 78},
		{Tamano: "3/4", AreaInteriorMM2: 122},
	}

	impedancia := valueobject.ResistenciaReactancia{R: 3.9, X: 0.164}

	mockTabla := &mockTablaNOMRepository{
		tablaTierra:       tablaTierra,
		tablaAmpacidad:    tablaAmpacidad,
		tablaCanalizacion: tablaCanalizacion,
		impedancia:        impedancia,
	}

	mockEquipo := &mockEquipoRepository{}

	uc := NewCalcularMemoriaUseCase(mockTabla, mockEquipo)

	// Crear input
	ctx := context.Background()
	tension, _ := valueobject.NewTension(220)
	input := dto.EquipoInput{
		Modo:                  dto.ModoManualAmperaje,
		AmperajeNominal:       50,
		Tension:               tension,
		FactorPotencia:        1.0,
		ITM:                   100,
		TipoCanalizacion:      entity.TipoCanalizacionTuberiaPVC,
		LongitudCircuito:      10,
		PorcentajeCaidaMaximo: 3.0,
	}

	// Ejecutar
	output, err := uc.Execute(ctx, input)

	// Verificar
	require.NoError(t, err)
	assert.Equal(t, 50.0, output.CorrienteNominal.Valor())
	assert.Equal(t, 50.0, output.CorrienteAjustada.Valor())        // Sin factores
	assert.Equal(t, "8 AWG", output.ConductorAlimentacion.Calibre) // 50A < 55A
	assert.Equal(t, "8 AWG", output.ConductorTierra.Calibre)       // ITM 100
	assert.True(t, output.CumpleNormativa)
}

func TestCalcularMemoriaUseCase_Execute_ValidationError(t *testing.T) {
	uc := NewCalcularMemoriaUseCase(nil, nil)

	ctx := context.Background()
	input := dto.EquipoInput{
		Modo: dto.ModoCalculo("INVALIDO"),
	}

	_, err := uc.Execute(ctx, input)
	assert.Error(t, err)
}

func TestCalcularMemoriaUseCase_seleccionarTemperatura(t *testing.T) {
	uc := &CalcularMemoriaUseCase{}

	tests := []struct {
		name      string
		corriente float64
		override  *valueobject.Temperatura
		canaliz   entity.TipoCanalizacion
		want      valueobject.Temperatura
	}{
		{"<= 100A default", 50, nil, entity.TipoCanalizacionTuberiaPVC, valueobject.Temp60},
		{"> 100A default", 150, nil, entity.TipoCanalizacionTuberiaPVC, valueobject.Temp75},
		{"Charola triangular", 50, nil, entity.TipoCanalizacionCharolaCableTriangular, valueobject.Temp75},
		{"Override 90C", 50, &[]valueobject.Temperatura{valueobject.Temp90}[0], entity.TipoCanalizacionTuberiaPVC, valueobject.Temp90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corriente, _ := valueobject.NewCorriente(tt.corriente)
			input := dto.EquipoInput{
				TipoCanalizacion:    tt.canaliz,
				TemperaturaOverride: tt.override,
			}
			got := uc.seleccionarTemperatura(corriente, input)
			assert.Equal(t, tt.want, got)
		})
	}
}
