// internal/calculos/application/usecase/orquestador_memoria_test.go
package usecase

import (
	"context"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepoForIntegration es un mock completo para tests de integración
type mockRepoForIntegration struct {
	tablaTierra       []valueobject.EntradaTablaTierra
	tablaAmpacidad    []valueobject.EntradaTablaConductor
	tablaCanalizacion []valueobject.EntradaTablaCanalizacion
	impedancia        valueobject.ResistenciaReactancia
	estadosTemp       map[string]int
}

func newMockRepoForIntegration() *mockRepoForIntegration {
	return &mockRepoForIntegration{
		tablaTierra: []valueobject.EntradaTablaTierra{
			{ITMHasta: 100, ConductorCu: valueobject.ConductorParams{Calibre: "8 AWG", Material: valueobject.MaterialCobre, SeccionMM2: 8.37}},
			{ITMHasta: 200, ConductorCu: valueobject.ConductorParams{Calibre: "6 AWG", Material: valueobject.MaterialCobre, SeccionMM2: 13.3}},
		},
		tablaAmpacidad: []valueobject.EntradaTablaConductor{
			{Capacidad: 30, Conductor: valueobject.ConductorParams{Calibre: "10 AWG", Material: valueobject.MaterialCobre, SeccionMM2: 5.26}},
			{Capacidad: 55, Conductor: valueobject.ConductorParams{Calibre: "8 AWG", Material: valueobject.MaterialCobre, SeccionMM2: 8.37}},
			{Capacidad: 75, Conductor: valueobject.ConductorParams{Calibre: "6 AWG", Material: valueobject.MaterialCobre, SeccionMM2: 13.3}},
		},
		tablaCanalizacion: []valueobject.EntradaTablaCanalizacion{
			{Tamano: "1/2", AreaInteriorMM2: 78},
			{Tamano: "3/4", AreaInteriorMM2: 122},
			{Tamano: "1", AreaInteriorMM2: 188},
		},
		impedancia: func() valueobject.ResistenciaReactancia {
			rr, _ := valueobject.NewResistenciaReactancia(3.9, 0.164)
			return rr
		}(),
		estadosTemp: map[string]int{"Sonora": 25},
	}
}

func (m *mockRepoForIntegration) ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error) {
	return m.tablaTierra, nil
}

func (m *mockRepoForIntegration) ObtenerTablaAmpacidad(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura) ([]valueobject.EntradaTablaConductor, error) {
	return m.tablaAmpacidad, nil
}

func (m *mockRepoForIntegration) ObtenerCapacidadConductor(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura, calibre string) (float64, error) {
	for _, e := range m.tablaAmpacidad {
		if e.Conductor.Calibre == calibre {
			return e.Capacidad, nil
		}
	}
	return 0, nil
}

func (m *mockRepoForIntegration) ObtenerImpedancia(ctx context.Context, calibre string, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor) (valueobject.ResistenciaReactancia, error) {
	return m.impedancia, nil
}

func (m *mockRepoForIntegration) ObtenerTablaCanalizacion(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]valueobject.EntradaTablaCanalizacion, error) {
	return m.tablaCanalizacion, nil
}

func (m *mockRepoForIntegration) ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error) {
	if temp, ok := m.estadosTemp[estado]; ok {
		return temp, nil
	}
	return 25, nil
}

func (m *mockRepoForIntegration) ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error) {
	return 1.0, nil
}

func (m *mockRepoForIntegration) ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error) {
	return 1.0, nil
}

// mockSeleccionarTemperatura implements SeleccionarTemperaturaPort for testing
type mockSeleccionarTemperatura struct{}

func (m *mockSeleccionarTemperatura) SeleccionarTemperatura(
	corriente valueobject.Corriente,
	tipoCanalizacion entity.TipoCanalizacion,
	override *valueobject.Temperatura,
) valueobject.Temperatura {
	// Default to 60°C for testing
	if override != nil {
		return *override
	}
	if corriente.Valor() <= 100 {
		return valueobject.Temp60
	}
	return valueobject.Temp75
}

func (m *mockRepoForIntegration) ObtenerDiametroConductor(ctx context.Context, calibre string, material string, conAislamiento bool) (float64, error) {
	return 3.5, nil
}

func (m *mockRepoForIntegration) ObtenerCharolaPorAncho(ctx context.Context, anchoRequeridoMM float64) (valueobject.EntradaTablaCanalizacion, error) {
	return valueobject.EntradaTablaCanalizacion{Tamano: "100mm", AreaInteriorMM2: 5000}, nil
}

func TestOrquestadorMemoriaCalculo_Execute(t *testing.T) {
	// Setup
	repo := newMockRepoForIntegration()

	// Create micro use cases
	calcularCorrienteUC := NewCalcularCorrienteUseCase(nil)
	ajustarCorrienteUC := NewAjustarCorrienteUseCase(repo)
	seleccionarConductorUC := NewSeleccionarConductorUseCase(repo)
	dimensionarCanalizacionUC := NewDimensionarCanalizacionUseCase(repo)
	calcularCaidaTensionUC := NewCalcularCaidaTensionUseCase(repo)

	// Create orquestador
	orquestador := NewOrquestadorMemoriaCalculo(
		calcularCorrienteUC,
		ajustarCorrienteUC,
		seleccionarConductorUC,
		dimensionarCanalizacionUC,
		calcularCaidaTensionUC,
		repo,
	)

	// Input
	tension, err := valueobject.NewTension(220)
	require.NoError(t, err)

	input := dto.EquipoInput{
		Modo:             dto.ModoManualAmperaje,
		AmperajeNominal:  50,
		TipoEquipo:       string(entity.TipoEquipoFiltroActivo),
		Tension:          tension,
		FactorPotencia:   0.9,
		ITM:              100,
		TipoCanalizacion: "TUBERIA_PVC",
		HilosPorFase:     1,
		NumTuberias:      1,
		Material:         valueobject.MaterialCobre,
		LongitudCircuito: 10,
		Estado:           "Sonora",
		SistemaElectrico: dto.SistemaElectricoDelta,
	}

	// Execute
	output, err := orquestador.Execute(context.Background(), input)
	require.NoError(t, err)

	// Verify basic fields
	assert.Equal(t, 50.0, output.CorrienteNominal)
	// Corriente ajustada = nominal × factor_temp × factor_agrupamiento × factor_uso
	// 50 × 1.0 × 1.0 × 1.35 (filtro activo) = 67.5
	assert.InDelta(t, 67.5, output.CorrienteAjustada, 0.1)
	assert.Equal(t, 220, output.Tension)
	assert.Equal(t, 100, output.ITM)
	assert.Equal(t, "Sonora", output.Estado)
	assert.Equal(t, dto.SistemaElectricoDelta, output.SistemaElectrico)
	assert.Equal(t, 3, output.CantidadConductores) // Delta = 3

	// Verify conductor selected
	assert.NotEmpty(t, output.ConductorAlimentacion.Calibre)
	assert.Equal(t, "CU", output.ConductorAlimentacion.Material)
	assert.NotZero(t, output.ConductorAlimentacion.SeccionMM2)

	// Verify tierra
	assert.NotEmpty(t, output.ConductorTierra.Calibre)

	// Verify canalización
	assert.NotEmpty(t, output.Canalizacion.Tamano)
	assert.Greater(t, output.Canalizacion.AreaTotalMM2, 0.0)

	// Verify caída de tensión
	assert.True(t, output.CaidaTension.Cumple) // With short length, should pass
	assert.Greater(t, output.CaidaTension.LimitePorcentaje, 0.0)

	// Verify normativa
	assert.True(t, output.CumpleNormativa)
}
