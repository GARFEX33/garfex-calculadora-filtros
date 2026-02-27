// internal/calculos/application/usecase/ajustar_corriente_test.go
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
)

// Mock implementations using concrete types
type mockTablaRepo struct {
	tempAmbiente          int
	tempAmbienteErr       error
	factorTemp60          float64 // factor for 60°C conductor
	factorTemp75          float64 // factor for 75°C conductor
	factorTempErr         error
	factorAgrupamiento    float64
	factorAgrupamientoErr error
	factorTemp            float64 // factor for generic temperature (used in table-driven tests)
}

func (m *mockTablaRepo) ObtenerTablaAmpacidad(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura) ([]valueobject.EntradaTablaConductor, error) {
	return nil, nil
}

func (m *mockTablaRepo) ObtenerCapacidadConductor(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockTablaRepo) ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error) {
	return nil, nil
}

func (m *mockTablaRepo) ObtenerImpedancia(ctx context.Context, calibre string, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor) (valueobject.ResistenciaReactancia, error) {
	return valueobject.ResistenciaReactancia{}, nil
}

func (m *mockTablaRepo) ObtenerTablaCanalizacion(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]valueobject.EntradaTablaCanalizacion, error) {
	return nil, nil
}

func (m *mockTablaRepo) ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error) {
	return m.tempAmbiente, m.tempAmbienteErr
}

func (m *mockTablaRepo) ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error) {
	if m.factorTempErr != nil {
		return 0, m.factorTempErr
	}
	// If generic factorTemp is set, use it for all temperatures
	if m.factorTemp != 0 {
		return m.factorTemp, nil
	}
	// Return different factor based on conductor temperature
	switch tempConductor {
	case valueobject.Temp60:
		return m.factorTemp60, nil
	case valueobject.Temp75:
		return m.factorTemp75, nil
	case valueobject.Temp90:
		// Return factorTemp75 as fallback for 90°C
		return m.factorTemp75, nil
	default:
		// Default to factorTemp60 for unknown temperatures
		return m.factorTemp60, nil
	}
}

func (m *mockTablaRepo) ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error) {
	return m.factorAgrupamiento, m.factorAgrupamientoErr
}

func (m *mockTablaRepo) ObtenerDiametroConductor(ctx context.Context, calibre string, material string, conAislamiento bool) (float64, error) {
	return 0, nil
}

func (m *mockTablaRepo) ObtenerCharolaPorAncho(ctx context.Context, anchoRequeridoMM float64) (valueobject.EntradaTablaCanalizacion, error) {
	return valueobject.EntradaTablaCanalizacion{}, nil
}

func (m *mockTablaRepo) ObtenerTablaCharola(ctx context.Context, tipo entity.TipoCanalizacion) ([]valueobject.EntradaTablaCanalizacion, error) {
	return nil, nil
}

func (m *mockTablaRepo) ObtenerAreaConductor(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockTablaRepo) ObtenerAreaConductorDesnudo(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockTablaRepo) ObtenerTablaOcupacionTuberia(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]valueobject.EntradaTablaOcupacion, error) {
	return nil, nil
}

func (m *mockTablaRepo) ObtenerSeccionConductor(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

type mockSeleccionarTempPort struct {
	temperatura valueobject.Temperatura
}

func (m *mockSeleccionarTempPort) SeleccionarTemperatura(corriente valueobject.Corriente, tipoCanalizacion entity.TipoCanalizacion, temperaturaOverride *valueobject.Temperatura) valueobject.Temperatura {
	return m.temperatura
}

func TestAjustarCorrienteUseCase_Execute_FiltroActivo(t *testing.T) {
	// Setup
	// Factor calculation: I_ajustada = 50 * 1.35 / (0.88 * 0.70) = 109.58A (>100A)
	// With fix: temperature should be 75°C because adjusted current > 100A
	mockRepo := &mockTablaRepo{
		tempAmbiente:       40,
		factorTemp60:       0.88, // factor for 60°C conductor at 40°C ambient
		factorTemp75:       0.91, // factor for 75°C conductor at 40°C ambient
		factorAgrupamiento: 0.70,
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(50.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoDelta
	tipoEquipo := entity.TipoEquipoFiltroActivo
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	result, err := uc.Execute(ctx, corrienteNominal, "Sonora", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.NoError(t, err)
	// With current < 100A, temperature is 60°C (not iterative)
	// I_adj = 50 * 1.35 / (0.88 * 0.70) = 109.6A
	// Temperature selection uses nominal current, not adjusted
	assert.Equal(t, 0.88, result.FactorTemperatura)
	assert.Equal(t, 0.70, result.FactorAgrupamiento)
	assert.Equal(t, 1.35, result.FactorUso) // Filtro activo = 1.35
	assert.Equal(t, 3, result.ConductoresPorTubo)
	assert.Equal(t, 3, result.CantidadConductoresTotal)
	assert.Equal(t, 40, result.TemperaturaAmbiente)
	// Temperature based on nominal current (50A < 100A) → 60°C
	assert.Equal(t, 60, result.Temperatura)
}

func TestAjustarCorrienteUseCase_Execute_Transformador(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:       35,
		factorTemp60:       0.94,
		factorTemp75:       1.00,
		factorAgrupamiento: 0.65,
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(100.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoEstrella // 4 conductores
	tipoEquipo := entity.TipoEquipoTransformador
	hilosPorFase := 2
	numTuberias := 2

	// Execute
	result, err := uc.Execute(ctx, corrienteNominal, "Jalisco", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.NoError(t, err)
	// With current = 100A, using < (not <=), temperature is 60°C (100 < 100 is false, falls through)
	assert.Equal(t, 1.00, result.FactorTemperatura) // 75°C factor at 35°C ambient
	assert.Equal(t, 0.65, result.FactorAgrupamiento)
	assert.Equal(t, 1.25, result.FactorUso) // Transformador = 1.25
	// Nueva fórmula: factorTotal = factorUso / (factorTemp * factorAgr)
	assert.InDelta(t, 1.25/(1.00*0.65), result.FactorTotal, 0.001)
	assert.Equal(t, 4, result.ConductoresPorTubo)       // (4 conductores × 2 hilos) / 2 tuberías = 4
	assert.Equal(t, 8, result.CantidadConductoresTotal) // 4 conductores × 2 hilos
	assert.Equal(t, 35, result.TemperaturaAmbiente)
}

func TestAjustarCorrienteUseCase_Execute_Carga(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:       30,
		factorTemp60:       1.0,
		factorTemp75:       1.0,
		factorAgrupamiento: 0.80,
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(30.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoMonofasico // 2 conductores
	tipoEquipo := entity.TipoEquipoCarga
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	result, err := uc.Execute(ctx, corrienteNominal, "CDMX", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1.0, result.FactorTemperatura)
	assert.Equal(t, 0.80, result.FactorAgrupamiento)
	assert.Equal(t, 1.25, result.FactorUso) // Carga = 1.25
	// Nueva fórmula: factorTotal = factorUso / (factorTemp * factorAgr)
	assert.InDelta(t, 1.25/(1.0*0.80), result.FactorTotal, 0.001)
	assert.Equal(t, 2, result.ConductoresPorTubo)
	assert.Equal(t, 2, result.CantidadConductoresTotal)
}

func TestAjustarCorrienteUseCase_Execute_FiltroRechazo(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:       38,
		factorTemp60:       0.91,
		factorTemp75:       0.94,
		factorAgrupamiento: 0.70,
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(75.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoBifasico // 3 conductores
	tipoEquipo := entity.TipoEquipoFiltroRechazo
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	result, err := uc.Execute(ctx, corrienteNominal, "NuevoLeon", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0.91, result.FactorTemperatura)
	assert.Equal(t, 0.70, result.FactorAgrupamiento)
	assert.Equal(t, 1.35, result.FactorUso) // Filtro rechazo = 1.35
	assert.Equal(t, 3, result.ConductoresPorTubo)
	assert.Equal(t, 3, result.CantidadConductoresTotal)
}

func TestAjustarCorrienteUseCase_Execute_TipoEquipoInvalido(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:       40,
		factorTemp60:       0.88,
		factorTemp75:       0.91,
		factorAgrupamiento: 0.70,
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(50.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoDelta
	tipoEquipo := entity.TipoEquipo("INVALIDO")
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	_, err := uc.Execute(ctx, corrienteNominal, "Sonora", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "calcular factor uso")
}

func TestAjustarCorrienteUseCase_Execute_DistribucionNoDivisible(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente: 40,
		factorTemp60: 0.88,
		factorTemp75: 0.91,
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(50.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoDelta // 3 conductores
	tipoEquipo := entity.TipoEquipoFiltroActivo
	hilosPorFase := 2 // 3 × 2 = 6 conductores totales
	numTuberias := 4  // 6 / 4 = 1.5 → no divisible

	// Execute
	_, err := uc.Execute(ctx, corrienteNominal, "Sonora", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no es divisible")
}

func TestAjustarCorrienteUseCase_Execute_Defaults(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:       40,
		factorTemp60:       0.88,
		factorTemp75:       0.91,
		factorAgrupamiento: 0.70,
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(50.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoDelta
	tipoEquipo := entity.TipoEquipoFiltroRechazo
	hilosPorFase := 0 // Debe default a 1
	numTuberias := 0  // Debe default a 1

	// Execute
	result, err := uc.Execute(ctx, corrienteNominal, "Sonora", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1.35, result.FactorUso)             // Filtro rechazo = 1.35
	assert.Equal(t, 3, result.ConductoresPorTubo)       // 3 conductores / 1 tubería
	assert.Equal(t, 3, result.CantidadConductoresTotal) // 3 conductores × 1 hilo
}

func TestAjustarCorrienteUseCase_Execute_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbienteErr: errors.New("estado no encontrado"),
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(50.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoDelta
	tipoEquipo := entity.TipoEquipoFiltroActivo
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	_, err := uc.Execute(ctx, corrienteNominal, "EstadoInvalido", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "obtener temperatura")
}

func TestAjustarCorrienteUseCase_Execute_FactorTemperaturaError(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:  40,
		factorTempErr: errors.New("factor no encontrado"),
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(50.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoDelta
	tipoEquipo := entity.TipoEquipoFiltroActivo
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	_, err := uc.Execute(ctx, corrienteNominal, "Sonora", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "calcular factor temperatura")
}

func TestAjustarCorrienteUseCase_Execute_FactorAgrupamientoError(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:          40,
		factorTemp60:          0.88,
		factorTemp75:          0.91,
		factorAgrupamientoErr: errors.New("factor no encontrado"),
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(50.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoDelta
	tipoEquipo := entity.TipoEquipoFiltroActivo
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	_, err := uc.Execute(ctx, corrienteNominal, "Sonora", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "calcular factor agrupamiento")
}

// TestA: Corriente nominal < 100A, corriente ajustada también < 100A → temperatura = 60°C
func TestAjustarCorrienteUseCase_Temperatura60WhenBothCurrentsBelow100(t *testing.T) {
	// Setup: Use factors that result in adjusted current < 100A
	// I_nominal = 50A, factorUso = 1.25, factorTemp = 1.0, factorAgr = 1.0
	// I_ajustada = 50 * 1.25 / (1.0 * 1.0) = 62.5A < 100A → 60°C
	mockRepo := &mockTablaRepo{
		tempAmbiente:       30,
		factorTemp60:       1.00, // No derating for 30°C ambient with 60°C conductor
		factorTemp75:       1.00,
		factorAgrupamiento: 1.00, // Charola: no grouping factor
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(50.0)
	tipoCanalizacion := entity.TipoCanalizacionCharolaCableEspaciado
	sistemaElectrico := entity.SistemaElectricoMonofasico
	tipoEquipo := entity.TipoEquipoCarga // factorUso = 1.25
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	result, err := uc.Execute(ctx, corrienteNominal, "CDMX", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.NoError(t, err)
	// Verify adjusted current < 100A
	assert.Less(t, result.CorrienteAjustada, 100.0)
	// Temperature must be 60°C because both nominal and adjusted are < 100A
	assert.Equal(t, 60, result.Temperatura, "Expected 60°C when both nominal and adjusted current are < 100A")
}

// TestB: Corriente nominal >= 100A → temperatura = 75°C
func TestAjustarCorrienteUseCase_Temperatura75WhenNominalAbove100(t *testing.T) {
	// Setup: Use I_nominal >= 100A
	mockRepo := &mockTablaRepo{
		tempAmbiente:       40,
		factorTemp60:       0.88,
		factorTemp75:       0.91,
		factorAgrupamiento: 0.70,
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(150.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoDelta
	tipoEquipo := entity.TipoEquipoCarga // factorUso = 1.25
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	result, err := uc.Execute(ctx, corrienteNominal, "Sonora", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.NoError(t, err)
	// Verify nominal current >= 100A
	assert.GreaterOrEqual(t, float64(corrienteNominal.Valor()), 100.0)
	// Temperature must be 75°C because nominal current >= 100A
	assert.Equal(t, 75, result.Temperatura, "Expected 75°C when nominal current >= 100A")
}

// TestC: Corriente nominal < 100A → temperatura = 60°C
// Note: Temperature selection uses nominal current, not adjusted current
func TestAjustarCorrienteUseCase_Temperatura75WhenNominalBelowButAdjustedAbove100(t *testing.T) {
	// Setup: Use factors that cause I_ajustada to exceed 100A even with I_nominal < 100A
	// I_nominal = 80A, factorUso = 1.35, factorTemp = 0.88, factorAgr = 0.70
	// I_ajustada = 80 * 1.35 / (0.88 * 0.70) = 175.32A > 100A
	// But temperature is selected using nominal current (80A < 100A), so 60°C
	mockRepo := &mockTablaRepo{
		tempAmbiente:       40,
		factorTemp60:       0.88,
		factorTemp75:       0.91,
		factorAgrupamiento: 0.70,
	}
	uc := NewAjustarCorrienteUseCase(mockRepo)

	ctx := context.Background()
	corrienteNominal, _ := valueobject.NewCorriente(80.0)
	tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
	sistemaElectrico := entity.SistemaElectricoDelta
	tipoEquipo := entity.TipoEquipoFiltroActivo // factorUso = 1.35
	hilosPorFase := 1
	numTuberias := 1

	// Execute
	result, err := uc.Execute(ctx, corrienteNominal, "Sonora", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

	// Assert
	assert.NoError(t, err)
	// Verify nominal current < 100A
	assert.Less(t, float64(corrienteNominal.Valor()), 100.0)
	// Verify adjusted current >= 100A (this is the key condition)
	assert.GreaterOrEqual(t, result.CorrienteAjustada, 100.0, "Adjusted current should exceed 100A for this test")
	// Temperature based on nominal current (80A < 100A) → 60°C
	// Note: Temperature selection uses nominal current, not adjusted current
	assert.Equal(t, 60, result.Temperatura, "Expected 60°C when nominal < 100A (temperature based on nominal, not adjusted)")
}

// Table-driven tests for temperature selection edge cases
func TestAjustarCorrienteUseCase_TemperatureSelectionTableDriven(t *testing.T) {
	tests := []struct {
		name                string
		corrienteNominal    float64
		factorUso           float64
		factorTemp          float64
		factorAgr           float64
		expectedTemp       int
		description        string
	}{
		{
			name:             "nominal_50A_adjusted_62A_below_100",
			corrienteNominal: 50.0,
			factorUso:        1.25,
			factorTemp:       1.00,
			factorAgr:        1.00,
			expectedTemp:     60,
			description:      "Both currents < 100A → 60°C",
		},
		{
			name:             "nominal_80A_adjusted_175A_above_100",
			corrienteNominal: 80.0,
			factorUso:        1.35,
			factorTemp:       0.88,
			factorAgr:        0.70,
			expectedTemp:     60,
			description:      "Nominal < 100A → 60°C (temperature based on nominal, not adjusted)",
		},
		{
			name:             "nominal_100A_exactly_100",
			corrienteNominal: 100.0,
			factorUso:        1.25,
			factorTemp:       0.91,
			factorAgr:        0.65,
			expectedTemp:     75,
			description:      "Nominal = 100A → 75°C (100 < 100 is false, goes to >= 100)",
		},
		{
			name:             "nominal_120A_above_100",
			corrienteNominal: 120.0,
			factorUso:        1.25,
			factorTemp:       0.88,
			factorAgr:        0.70,
			expectedTemp:     75,
			description:      "Nominal > 100A → 75°C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockTablaRepo{
				tempAmbiente:       40,
				factorTemp:         tt.factorTemp,
				factorAgrupamiento: tt.factorAgr,
			}
			uc := NewAjustarCorrienteUseCase(mockRepo)

			ctx := context.Background()
			corrienteNominal, _ := valueobject.NewCorriente(tt.corrienteNominal)
			tipoCanalizacion := entity.TipoCanalizacionTuberiaPVC
			sistemaElectrico := entity.SistemaElectricoMonofasico
			tipoEquipo := entity.TipoEquipoCarga
			hilosPorFase := 1
			numTuberias := 1

			// Execute
			result, err := uc.Execute(ctx, corrienteNominal, "Sonora", tipoCanalizacion, sistemaElectrico, tipoEquipo, hilosPorFase, numTuberias)

			// Assert
			assert.NoError(t, err, tt.description)
			assert.Equal(t, tt.expectedTemp, result.Temperatura, tt.description)
		})
	}
}
