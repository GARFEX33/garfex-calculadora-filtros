// internal/calculos/application/usecase/seleccionar_conductor_alimentacion_test.go
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
)

// Mock implementation of TablaNOMRepository for testing
type mockConductorAlimentacionRepo struct {
	tablaAmpacidad     []valueobject.EntradaTablaConductor
	tablaAmpacidadErr error
	capacidadConductor float64
	capacidadErr       error
}

func (m *mockConductorAlimentacionRepo) ObtenerTablaAmpacidad(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
) ([]valueobject.EntradaTablaConductor, error) {
	if m.tablaAmpacidadErr != nil {
		return nil, m.tablaAmpacidadErr
	}
	return m.tablaAmpacidad, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerCapacidadConductor(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
	calibre string,
) (float64, error) {
	if m.capacidadErr != nil {
		return 0, m.capacidadErr
	}
	return m.capacidadConductor, nil
}

// Unused mock methods - required to implement interface
func (m *mockConductorAlimentacionRepo) ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error) {
	return nil, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerImpedancia(
	ctx context.Context,
	calibre string,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
) (valueobject.ResistenciaReactancia, error) {
	return valueobject.ResistenciaReactancia{}, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerTablaCanalizacion(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
) ([]valueobject.EntradaTablaCanalizacion, error) {
	return nil, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error) {
	return 30, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerFactorTemperatura(
	ctx context.Context,
	tempAmbiente int,
	tempConductor valueobject.Temperatura,
) (float64, error) {
	return 1.0, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error) {
	return 1.0, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerDiametroConductor(
	ctx context.Context,
	calibre string,
	material string,
	conAislamiento bool,
) (float64, error) {
	return 0, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerCharolaPorAncho(
	ctx context.Context,
	anchoRequeridoMM float64,
) (valueobject.EntradaTablaCanalizacion, error) {
	return valueobject.EntradaTablaCanalizacion{}, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerTablaCharola(
	ctx context.Context,
	tipo entity.TipoCanalizacion,
) ([]valueobject.EntradaTablaCanalizacion, error) {
	return nil, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerAreaConductor(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerAreaConductorDesnudo(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerTablaOcupacionTuberia(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
) ([]valueobject.EntradaTablaOcupacion, error) {
	return nil, nil
}

func (m *mockConductorAlimentacionRepo) ObtenerSeccionConductor(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

// Test table entries - using proper ConductorParams structure
// entradaConductorTest builds an EntradaTablaConductor with the fields relevant
// for conductor selection.
func entradaConductorTest(calibre string, capacidad, seccionMM2 float64) valueobject.EntradaTablaConductor {
	return valueobject.EntradaTablaConductor{
		Capacidad: capacidad,
		Conductor: valueobject.ConductorParams{
			Calibre:         calibre,
			Material:        valueobject.MaterialCobre,
			TipoAislamiento: "THHN",
			SeccionMM2:      seccionMM2,
		},
	}
}

// Simplified NOM table 310-15(b)(16) excerpt for testing
// Note: Only valid NOM calibres are used (1 AWG and 3 AWG are NOT valid)
var tablaConductorTestData = []valueobject.EntradaTablaConductor{
	entradaConductorTest("14 AWG", 15, 2.08),
	entradaConductorTest("12 AWG", 20, 3.31),
	entradaConductorTest("10 AWG", 30, 5.26),
	entradaConductorTest("8 AWG", 40, 8.37),
	entradaConductorTest("6 AWG", 55, 13.30),
	entradaConductorTest("4 AWG", 70, 21.15),
	entradaConductorTest("2 AWG", 95, 33.62),
	entradaConductorTest("1/0 AWG", 130, 53.49),
	entradaConductorTest("2/0 AWG", 150, 67.43),
	entradaConductorTest("3/0 AWG", 175, 85.01),
	entradaConductorTest("4/0 AWG", 200, 107.2),
	entradaConductorTest("250 MCM", 230, 126.64),
	entradaConductorTest("300 MCM", 255, 152.0),
	entradaConductorTest("350 MCM", 285, 177.3),
	entradaConductorTest("400 MCM", 310, 202.7),
	entradaConductorTest("500 MCM", 340, 253.4),
	entradaConductorTest("600 MCM", 375, 304.0),
	entradaConductorTest("750 MCM", 400, 380.0),
	entradaConductorTest("1000 MCM", 430, 506.7),
}

// TestSeleccionarConductorAlimentacion_EscenarioA tests that current < 100A selects 60°C
func TestSeleccionarConductorAlimentacion_EscenarioA_CorrienteMenos100A(t *testing.T) {
	// Scenario A: CorrienteAjustada < 100A → temperatura = 60°C
	// GIVEN CorrienteAjustada = 75A y tipoCanalizacion = "TUBERIA_PVC"
	// THEN temperatura = 60°C

	// Setup
	mockRepo := &mockConductorAlimentacionRepo{
		tablaAmpacidad:     tablaConductorTestData,
		capacidadConductor: 85, // 2 AWG capacity at selected temp
	}
	uc := NewSeleccionarConductorAlimentacionUseCase(mockRepo)

	ctx := context.Background()
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: 75.0,
		TipoCanalizacion:  "TUBERIA_PVC",
		Material:          "Cu",
		Temperatura:       nil, // No override
		HilosPorFase:      1,
	}

	// Execute
	result, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Calibre)
	// Verify that the table used contains "60°C" (temperatura = 60°C for current < 100A)
	assert.Contains(t, result.TablaUsada, "60°C", "Expected 60°C table for current < 100A")
}

// TestSeleccionarConductorAlimentacion_EscenarioB tests that current >= 100A selects 75°C
func TestSeleccionarConductorAlimentacion_EscenarioB_CorrienteMayorOIgual100A(t *testing.T) {
	// Scenario B: CorrienteAjustada >= 100A → temperatura = 75°C
	// GIVEN CorrienteAjustada = 150A y tipoCanalizacion = "TUBERIA_PVC"
	// THEN temperatura = 75°C

	// Setup
	mockRepo := &mockConductorAlimentacionRepo{
		tablaAmpacidad:     tablaConductorTestData,
		capacidadConductor: 150, // capacity for selected conductor
	}
	uc := NewSeleccionarConductorAlimentacionUseCase(mockRepo)

	ctx := context.Background()
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: 150.0,
		TipoCanalizacion:  "TUBERIA_PVC",
		Material:          "Cu",
		Temperatura:       nil, // No override
		HilosPorFase:      1,
	}

	// Execute
	result, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Calibre)
	// Verify that the table used contains "75°C" (temperatura = 75°C for current >= 100A)
	assert.Contains(t, result.TablaUsada, "75°C", "Expected 75°C table for current >= 100A")
}

// TestSeleccionarConductorAlimentacion_EscenarioC tests boundary at exactly 100A
func TestSeleccionarConductorAlimentacion_EscenarioC_CorrienteExactamente100A(t *testing.T) {
	// Scenario C: CorrienteAjustada exactamente 100A → temperatura = 75°C (boundary)
	// GIVEN CorrienteAjustada = 100A y tipoCanalizacion = "TUBERIA_PVC"
	// THEN temperatura = 75°C

	// Setup
	mockRepo := &mockConductorAlimentacionRepo{
		tablaAmpacidad:     tablaConductorTestData,
		capacidadConductor: 130, // capacity for 1/0 AWG
	}
	uc := NewSeleccionarConductorAlimentacionUseCase(mockRepo)

	ctx := context.Background()
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: 100.0,
		TipoCanalizacion:  "TUBERIA_PVC",
		Material:          "Cu",
		Temperatura:       nil, // No override
		HilosPorFase:      1,
	}

	// Execute
	result, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Calibre)
	// Verify that the table used contains "75°C" (boundary: 100A uses 75°C, not 60°C)
	assert.Contains(t, result.TablaUsada, "75°C", "Expected 75°C table for current = 100A (boundary)")
}

// TestSeleccionarConductorAlimentacion_EscenarioD tests that temperature override is ignored
func TestSeleccionarConductorAlimentacion_EscenarioD_OverrideIgnorado(t *testing.T) {
	// Scenario D: Override de temperatura ignorado cuando contradice regla NOM
	// GIVEN CorrienteAjustada = 120A Y input.Temperatura = 60 (override del caller)
	// WHEN se ejecuta SeleccionarConductorAlimentacion
	// THEN temperatura = 75°C (override ignorado, se usa regla NOM)

	// Setup
	mockRepo := &mockConductorAlimentacionRepo{
		tablaAmpacidad:     tablaConductorTestData,
		capacidadConductor: 150,
	}
	uc := NewSeleccionarConductorAlimentacionUseCase(mockRepo)

	ctx := context.Background()

	// Input with temperature override that contradicts NOM rule
	temp60 := 60
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: 120.0,
		TipoCanalizacion:  "TUBERIA_PVC",
		Material:          "Cu",
		Temperatura:       &temp60, // Override: caller wants 60°C
		HilosPorFase:      1,
	}

	// Execute
	result, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Calibre)
	// Verify that the table used contains "75°C" (override ignored, NOM rule applied)
	// Current >= 100A must use 75°C regardless of override
	assert.Contains(t, result.TablaUsada, "75°C", "Expected 75°C table even with 60°C override (NOM rule)")
}

// TestSeleccionarConductorAlimentacion_ValidationError tests input validation
func TestSeleccionarConductorAlimentacion_ValidationError(t *testing.T) {
	// Setup
	mockRepo := &mockConductorAlimentacionRepo{}
	uc := NewSeleccionarConductorAlimentacionUseCase(mockRepo)

	ctx := context.Background()

	// Test with invalid input: zero current
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: 0,
		TipoCanalizacion:  "TUBERIA_PVC",
	}

	// Execute
	_, err := uc.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "corriente_ajustada")
}

// TestSeleccionarConductorAlimentacion_EmptyChannelType tests empty channelization type
func TestSeleccionarConductorAlimentacion_EmptyChannelType(t *testing.T) {
	// Setup
	mockRepo := &mockConductorAlimentacionRepo{}
	uc := NewSeleccionarConductorAlimentacionUseCase(mockRepo)

	ctx := context.Background()

	// Test with invalid input: empty channelization type
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: 50.0,
		TipoCanalizacion:  "", // Empty - should fail validation
	}

	// Execute
	_, err := uc.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tipo_canalizacion")
}

// TestSeleccionarConductorAlimentacion_RepositoryError tests repository error handling
func TestSeleccionarConductorAlimentacion_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := &mockConductorAlimentacionRepo{
		tablaAmpacidadErr: errors.New("tabla no encontrada"),
	}
	uc := NewSeleccionarConductorAlimentacionUseCase(mockRepo)

	ctx := context.Background()
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: 50.0,
		TipoCanalizacion:  "TUBERIA_PVC",
		Material:          "Cu",
	}

	// Execute
	_, err := uc.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tabla ampacidad")
}

// TestSeleccionarConductorAlimentacion_HilosPorFaseDefault tests default HilosPorFase
func TestSeleccionarConductorAlimentacion_HilosPorFaseDefault(t *testing.T) {
	// Test that HilosPorFase defaults to 1 when set to 0

	// Setup
	mockRepo := &mockConductorAlimentacionRepo{
		tablaAmpacidad:     tablaConductorTestData,
		capacidadConductor: 85,
	}
	uc := NewSeleccionarConductorAlimentacionUseCase(mockRepo)

	ctx := context.Background()
	input := dto.ConductorAlimentacionInput{
		CorrienteAjustada: 75.0,
		TipoCanalizacion:  "TUBERIA_PVC",
		Material:          "Cu",
		HilosPorFase:      0, // Should default to 1
	}

	// Execute
	result, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Calibre)
}

// Table-driven tests for temperature selection
func TestSeleccionarConductorAlimentacion_TemperatureSelectionTableDriven(t *testing.T) {
	tests := []struct {
		name                string
		corrienteAjustada   float64
		overrideTemp        *int
		expectedTempInTable string
		description         string
	}{
		{
			name:                "corriente_75A_menor_100A",
			corrienteAjustada:   75.0,
			overrideTemp:        nil,
			expectedTempInTable: "60°C",
			description:         "Current < 100A → 60°C",
		},
		{
			name:                "corriente_99A_menor_100A",
			corrienteAjustada:   99.0,
			overrideTemp:        nil,
			expectedTempInTable: "60°C",
			description:         "Current 99A < 100A → 60°C",
		},
		{
			name:                "corriente_100A_exacto",
			corrienteAjustada:   100.0,
			overrideTemp:        nil,
			expectedTempInTable: "75°C",
			description:         "Current = 100A (boundary) → 75°C",
		},
		{
			name:                "corriente_150A_mayor_100A",
			corrienteAjustada:   150.0,
			overrideTemp:        nil,
			expectedTempInTable: "75°C",
			description:         "Current > 100A → 75°C",
		},
		{
			name:                "corriente_200A_mayor_100A",
			corrienteAjustada:   200.0,
			overrideTemp:        nil,
			expectedTempInTable: "75°C",
			description:         "Current >> 100A → 75°C",
		},
		{
			name:                "override_60C_ignorado_corriente_120A",
			corrienteAjustada:   120.0,
			overrideTemp:        intPtr(60),
			expectedTempInTable: "75°C",
			description:         "Override 60°C ignored when current >= 100A",
		},
		{
			name:                "override_75C_aceptado_pero_mismo_resultado",
			corrienteAjustada:   120.0,
			overrideTemp:        intPtr(75),
			expectedTempInTable: "75°C",
			description:         "Override 75°C matches NOM rule",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockConductorAlimentacionRepo{
				tablaAmpacidad:     tablaConductorTestData,
				capacidadConductor: 150,
			}
			uc := NewSeleccionarConductorAlimentacionUseCase(mockRepo)

			ctx := context.Background()
			input := dto.ConductorAlimentacionInput{
				CorrienteAjustada: tt.corrienteAjustada,
				TipoCanalizacion:  "TUBERIA_PVC",
				Material:           "Cu",
				Temperatura:        tt.overrideTemp,
				HilosPorFase:       1,
			}

			// Execute
			result, err := uc.Execute(ctx, input)

			// Assert
			assert.NoError(t, err, tt.description)
			assert.Contains(t, result.TablaUsada, tt.expectedTempInTable, tt.description)
		})
	}
}

// Helper function to create pointer to int
func intPtr(i int) *int {
	return &i
}
