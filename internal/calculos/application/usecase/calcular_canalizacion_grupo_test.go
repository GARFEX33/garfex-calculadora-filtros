// internal/calculos/application/usecase/calcular_canalizacion_grupo_test.go
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTablaNOMRepository es un mock de TablaNOMRepository.
type MockTablaNOMRepository struct {
	mock.Mock
}

func (m *MockTablaNOMRepository) ObtenerTablaAmpacidad(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
) ([]valueobject.EntradaTablaConductor, error) {
	args := m.Called(ctx, canalizacion, material, temperatura)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]valueobject.EntradaTablaConductor), args.Error(1)
}

func (m *MockTablaNOMRepository) ObtenerCapacidadConductor(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
	temperatura valueobject.Temperatura,
	calibre string,
) (float64, error) {
	args := m.Called(ctx, canalizacion, material, temperatura, calibre)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTablaNOMRepository) ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]valueobject.EntradaTablaTierra), args.Error(1)
}

func (m *MockTablaNOMRepository) ObtenerImpedancia(
	ctx context.Context,
	calibre string,
	canalizacion entity.TipoCanalizacion,
	material valueobject.MaterialConductor,
) (valueobject.ResistenciaReactancia, error) {
	args := m.Called(ctx, calibre, canalizacion, material)
	return args.Get(0).(valueobject.ResistenciaReactancia), args.Error(1)
}

func (m *MockTablaNOMRepository) ObtenerTablaCanalizacion(
	ctx context.Context,
	canalizacion entity.TipoCanalizacion,
) ([]valueobject.EntradaTablaCanalizacion, error) {
	args := m.Called(ctx, canalizacion)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]valueobject.EntradaTablaCanalizacion), args.Error(1)
}

func (m *MockTablaNOMRepository) ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error) {
	args := m.Called(ctx, estado)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockTablaNOMRepository) ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error) {
	args := m.Called(ctx, tempAmbiente, tempConductor)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTablaNOMRepository) ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error) {
	args := m.Called(ctx, cantidadConductores)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTablaNOMRepository) ObtenerDiametroConductor(ctx context.Context, calibre string, material string, conAislamiento bool) (float64, error) {
	args := m.Called(ctx, calibre, material, conAislamiento)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTablaNOMRepository) ObtenerCharolaPorAncho(ctx context.Context, anchoRequeridoMM float64) (valueobject.EntradaTablaCanalizacion, error) {
	args := m.Called(ctx, anchoRequeridoMM)
	return args.Get(0).(valueobject.EntradaTablaCanalizacion), args.Error(1)
}

func TestCalcularCanalizacionGrupoUseCase_Execute_Exitoso1Tubo(t *testing.T) {
	// Arrange
	mockRepo := new(MockTablaNOMRepository)
	uc := NewCalcularCanalizacionGrupoUseCase(mockRepo)

	// Input: 3 fases de 21.2 mm² (calibre 12) + 1 tierra de 8.37 mm² (calibre 8)
	input := dto.CanalizacionGrupoInput{
		Conductores: []dto.ConductorGrupoInput{
			{Cantidad: 3, SeccionMM2: 21.2},
		},
		SeccionTierraMM2: 8.37,
		TipoCanalizacion: "TUBERIA_PVC",
		NumeroDeTubos:    1,
	}

	// Mock: retornar tabla con entradas
	tabla := []valueobject.EntradaTablaCanalizacion{
		{Tamano: "27 mm (1\")", AreaInteriorMM2: 489},
		{Tamano: "35 mm (1 1/4\")", AreaInteriorMM2: 823},
		{Tamano: "41 mm (1 1/2\")", AreaInteriorMM2: 1117},
		{Tamano: "53 mm (2\")", AreaInteriorMM2: 1816},
	}
	mockRepo.On("ObtenerTablaCanalizacion", mock.Anything, entity.TipoCanalizacionTuberiaPVC).Return(tabla, nil)

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result.Tamano)
	assert.Greater(t, result.AreaTotalMM2, 0.0)
	assert.Equal(t, 1, result.NumeroDeTubos)
	// Verificar que se agregaron tierras: 3 fases + 1 tierra = 4 conductores
	// Area: 3*21.2 + 1*8.37 = 72.97 mm²
	// FactorRelleno para 4 conductores = 0.40
	assert.InDelta(t, 0.40, result.FactorRelleno, 0.01)
	mockRepo.AssertExpectations(t)
}

func TestCalcularCanalizacionGrupoUseCase_Execute_Exitoso2Tubos(t *testing.T) {
	// Arrange
	mockRepo := new(MockTablaNOMRepository)
	uc := NewCalcularCanalizacionGrupoUseCase(mockRepo)

	// Input: 3 fases de 21.2 mm² + 1 tierra de 8.37 mm², 2 tubos
	// Por norma: 2 tierras (1 por tubo)
	input := dto.CanalizacionGrupoInput{
		Conductores: []dto.ConductorGrupoInput{
			{Cantidad: 3, SeccionMM2: 21.2},
		},
		SeccionTierraMM2: 8.37,
		TipoCanalizacion: "TUBERIA_PVC",
		NumeroDeTubos:    2,
	}

	// Mock: tabla con tubo pequeño
	tabla := []valueobject.EntradaTablaCanalizacion{
		{Tamano: "27 mm (1\")", AreaInteriorMM2: 489},
		{Tamano: "35 mm (1 1/4\")", AreaInteriorMM2: 823},
		{Tamano: "41 mm (1 1/2\")", AreaInteriorMM2: 1117},
	}
	mockRepo.On("ObtenerTablaCanalizacion", mock.Anything, entity.TipoCanalizacionTuberiaPVC).Return(tabla, nil)

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, result.NumeroDeTubos)
	// Con 2 tubos: (3 fases + 2 tierras) / 2 = 2.5 → 2 o 3 por tubo
	// FactorRelleno para 2-3 conductores = 0.31
	assert.InDelta(t, 0.31, result.FactorRelleno, 0.01)
	mockRepo.AssertExpectations(t)
}

func TestCalcularCanalizacionGrupoUseCase_Execute_ErrorValidacionInput(t *testing.T) {
	// Arrange
	mockRepo := new(MockTablaNOMRepository)
	uc := NewCalcularCanalizacionGrupoUseCase(mockRepo)

	// Input inválido: sin conductores
	input := dto.CanalizacionGrupoInput{
		Conductores:      []dto.ConductorGrupoInput{},
		SeccionTierraMM2: 8.37,
		TipoCanalizacion: "TUBERIA_PVC",
		NumeroDeTubos:    1,
	}

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result.Tamano)
	assert.Contains(t, err.Error(), "validar input")
	mockRepo.AssertNotCalled(t, "ObtenerTablaCanalizacion")
}

func TestCalcularCanalizacionGrupoUseCase_Execute_ErrorTablaNoDisponible(t *testing.T) {
	// Arrange
	mockRepo := new(MockTablaNOMRepository)
	uc := NewCalcularCanalizacionGrupoUseCase(mockRepo)

	input := dto.CanalizacionGrupoInput{
		Conductores: []dto.ConductorGrupoInput{
			{Cantidad: 3, SeccionMM2: 21.2},
		},
		SeccionTierraMM2: 8.37,
		TipoCanalizacion: "TUBERIA_PVC",
		NumeroDeTubos:    1,
	}

	// Mock: tabla vacía
	mockRepo.On("ObtenerTablaCanalizacion", mock.Anything, entity.TipoCanalizacionTuberiaPVC).Return([]valueobject.EntradaTablaCanalizacion{}, nil)

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result.Tamano)
	assert.Contains(t, err.Error(), "calcular canalización")
	mockRepo.AssertExpectations(t)
}

func TestCalcularCanalizacionGrupoUseCase_Execute_ErrorTipoCanalizacionInvalido(t *testing.T) {
	// Arrange
	mockRepo := new(MockTablaNOMRepository)
	uc := NewCalcularCanalizacionGrupoUseCase(mockRepo)

	input := dto.CanalizacionGrupoInput{
		Conductores: []dto.ConductorGrupoInput{
			{Cantidad: 3, SeccionMM2: 21.2},
		},
		SeccionTierraMM2: 8.37,
		TipoCanalizacion: "TIPO_INVALIDO",
		NumeroDeTubos:    1,
	}

	// Act
	result, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result.Tamano)
	assert.Contains(t, err.Error(), "parsear tipo canalización")
	mockRepo.AssertNotCalled(t, "ObtenerTablaCanalizacion")
}

func TestConductorGrupoInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   dto.CanalizacionGrupoInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "válido con 1 grupo",
			input: dto.CanalizacionGrupoInput{
				Conductores:      []dto.ConductorGrupoInput{{Cantidad: 3, SeccionMM2: 21.2}},
				SeccionTierraMM2: 8.37,
				TipoCanalizacion: "TUBERIA_PVC",
				NumeroDeTubos:    1,
			},
			wantErr: false,
		},
		{
			name: "válido con 2 grupos",
			input: dto.CanalizacionGrupoInput{
				Conductores: []dto.ConductorGrupoInput{
					{Cantidad: 3, SeccionMM2: 21.2},
					{Cantidad: 2, SeccionMM2: 13.3},
				},
				SeccionTierraMM2: 8.37,
				TipoCanalizacion: "TUBERIA_CONDUIT",
				NumeroDeTubos:    2,
			},
			wantErr: false,
		},
		{
			name: "inválido - sin conductores",
			input: dto.CanalizacionGrupoInput{
				Conductores:      []dto.ConductorGrupoInput{},
				SeccionTierraMM2: 8.37,
				TipoCanalizacion: "TUBERIA_PVC",
				NumeroDeTubos:    1,
			},
			wantErr: true,
			errMsg:  "conductores",
		},
		{
			name: "inválido - tierra cero",
			input: dto.CanalizacionGrupoInput{
				Conductores:      []dto.ConductorGrupoInput{{Cantidad: 3, SeccionMM2: 21.2}},
				SeccionTierraMM2: 0,
				TipoCanalizacion: "TUBERIA_PVC",
				NumeroDeTubos:    1,
			},
			wantErr: true,
			errMsg:  "tierra",
		},
		{
			name: "inválido - sin tipo canalización",
			input: dto.CanalizacionGrupoInput{
				Conductores:      []dto.ConductorGrupoInput{{Cantidad: 3, SeccionMM2: 21.2}},
				SeccionTierraMM2: 8.37,
				TipoCanalizacion: "",
				NumeroDeTubos:    1,
			},
			wantErr: true,
			errMsg:  "tipo_canalizacion",
		},
		{
			name: "inválido - tubes menor a 1",
			input: dto.CanalizacionGrupoInput{
				Conductores:      []dto.ConductorGrupoInput{{Cantidad: 3, SeccionMM2: 21.2}},
				SeccionTierraMM2: 8.37,
				TipoCanalizacion: "TUBERIA_PVC",
				NumeroDeTubos:    0,
			},
			wantErr: true,
			errMsg:  "tubos",
		},
		{
			name: "inválido - conductor con cantidad cero",
			input: dto.CanalizacionGrupoInput{
				Conductores:      []dto.ConductorGrupoInput{{Cantidad: 0, SeccionMM2: 21.2}},
				SeccionTierraMM2: 8.37,
				TipoCanalizacion: "TUBERIA_PVC",
				NumeroDeTubos:    1,
			},
			wantErr: true,
			errMsg:  "cantidad",
		},
		{
			name: "inválido - conductor con sección cero",
			input: dto.CanalizacionGrupoInput{
				Conductores:      []dto.ConductorGrupoInput{{Cantidad: 3, SeccionMM2: 0}},
				SeccionTierraMM2: 8.37,
				TipoCanalizacion: "TUBERIA_PVC",
				NumeroDeTubos:    1,
			},
			wantErr: true,
			errMsg:  "seccion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCalcularCanalizacionGrupoUseCase_New(t *testing.T) {
	mockRepo := new(MockTablaNOMRepository)
	uc := NewCalcularCanalizacionGrupoUseCase(mockRepo)

	assert.NotNil(t, uc)
	assert.NotNil(t, uc.tablaRepo)
}

// Error de dominio simulado para el test de tabla no disponible
var errTablaVacia = errors.New("tabla vacía")
