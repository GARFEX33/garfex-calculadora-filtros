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
	factorTemp            float64
	factorTempErr         error
	factorAgrupamiento    float64
	factorAgrupamientoErr error
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
	return m.factorTemp, m.factorTempErr
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

func (m *mockTablaRepo) ObtenerAreaConductor(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockTablaRepo) ObtenerAreaConductorDesnudo(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockTablaRepo) ObtenerTablaOcupacionTuberia(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]valueobject.EntradaTablaOcupacion, error) {
	return nil, nil
}

type mockSeleccionarTempPort struct {
	temperatura valueobject.Temperatura
}

func (m *mockSeleccionarTempPort) SeleccionarTemperatura(corriente valueobject.Corriente, tipoCanalizacion entity.TipoCanalizacion, temperaturaOverride *valueobject.Temperatura) valueobject.Temperatura {
	return m.temperatura
}

func TestAjustarCorrienteUseCase_Execute_FiltroActivo(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:       40,
		factorTemp:         0.88,
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
	// Fórmula: corrienteAjustada = corrienteNominal * factorUso / (factorTemp * factorAgr)
	// 50 * 1.35 / (0.88 * 0.70) = 50 * 1.35 / 0.616 = 109.63
	expectedFactorTotal := 1.35 / (0.88 * 0.70) // factorUso / (factorTemp * factorAgr)
	expectedCorrienteAjustada := 50.0 * 1.35 / (0.88 * 0.70)
	assert.InDelta(t, expectedCorrienteAjustada, result.CorrienteAjustada, 0.001)
	assert.Equal(t, 0.88, result.FactorTemperatura)
	assert.Equal(t, 0.70, result.FactorAgrupamiento)
	assert.Equal(t, 1.35, result.FactorUso) // Filtro activo = 1.35
	assert.InDelta(t, expectedFactorTotal, result.FactorTotal, 0.001)
	assert.Equal(t, 3, result.ConductoresPorTubo)
	assert.Equal(t, 3, result.CantidadConductoresTotal)
	assert.Equal(t, 40, result.TemperaturaAmbiente)
	// Domain service returns 60°C for current <= 100A (non-triangular charola)
	assert.Equal(t, 60, result.Temperatura)
}

func TestAjustarCorrienteUseCase_Execute_Transformador(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:       35,
		factorTemp:         0.94,
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
	assert.Equal(t, 0.94, result.FactorTemperatura)
	assert.Equal(t, 0.65, result.FactorAgrupamiento)
	assert.Equal(t, 1.25, result.FactorUso) // Transformador = 1.25
	// Nueva fórmula: factorTotal = factorUso / (factorTemp * factorAgr)
	assert.InDelta(t, 1.25/(0.94*0.65), result.FactorTotal, 0.001)
	assert.Equal(t, 4, result.ConductoresPorTubo)       // (4 conductores × 2 hilos) / 2 tuberías = 4
	assert.Equal(t, 8, result.CantidadConductoresTotal) // 4 conductores × 2 hilos
	assert.Equal(t, 35, result.TemperaturaAmbiente)
}

func TestAjustarCorrienteUseCase_Execute_Carga(t *testing.T) {
	// Setup
	mockRepo := &mockTablaRepo{
		tempAmbiente:       30,
		factorTemp:         1.0,
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
		factorTemp:         0.91,
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
		factorTemp:         0.88,
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
		factorTemp:   0.88,
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
		factorTemp:         0.88,
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
		factorTemp:            0.88,
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
