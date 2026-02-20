// internal/calculos/application/usecase/calcular_charola_espaciado_test.go
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

// mockCharolaRepo es un mock para TablaNOMRepository específico para tests de charola.
type mockCharolaRepo struct {
	tablaCharola []valueobject.EntradaTablaCanalizacion
	tablaErr     error
}

func (m *mockCharolaRepo) ObtenerTablaAmpacidad(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura) ([]valueobject.EntradaTablaConductor, error) {
	return nil, nil
}

func (m *mockCharolaRepo) ObtenerCapacidadConductor(ctx context.Context, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor, temperatura valueobject.Temperatura, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockCharolaRepo) ObtenerTablaTierra(ctx context.Context) ([]valueobject.EntradaTablaTierra, error) {
	return nil, nil
}

func (m *mockCharolaRepo) ObtenerImpedancia(ctx context.Context, calibre string, canalizacion entity.TipoCanalizacion, material valueobject.MaterialConductor) (valueobject.ResistenciaReactancia, error) {
	return valueobject.ResistenciaReactancia{}, nil
}

func (m *mockCharolaRepo) ObtenerTablaCanalizacion(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]valueobject.EntradaTablaCanalizacion, error) {
	return nil, nil
}

func (m *mockCharolaRepo) ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error) {
	return 30, nil
}

func (m *mockCharolaRepo) ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error) {
	return 1.0, nil
}

func (m *mockCharolaRepo) ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error) {
	return 1.0, nil
}

func (m *mockCharolaRepo) ObtenerDiametroConductor(ctx context.Context, calibre string, material string, conAislamiento bool) (float64, error) {
	return 0, nil
}

func (m *mockCharolaRepo) ObtenerCharolaPorAncho(ctx context.Context, anchoRequeridoMM float64) (valueobject.EntradaTablaCanalizacion, error) {
	return valueobject.EntradaTablaCanalizacion{}, nil
}

func (m *mockCharolaRepo) ObtenerTablaCharola(ctx context.Context, tipo entity.TipoCanalizacion) ([]valueobject.EntradaTablaCanalizacion, error) {
	return m.tablaCharola, m.tablaErr
}

func (m *mockCharolaRepo) ObtenerAreaConductor(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockCharolaRepo) ObtenerAreaConductorDesnudo(ctx context.Context, calibre string) (float64, error) {
	return 0, nil
}

func (m *mockCharolaRepo) ObtenerTablaOcupacionTuberia(ctx context.Context, canalizacion entity.TipoCanalizacion) ([]valueobject.EntradaTablaOcupacion, error) {
	return nil, nil
}

// pointerToFloat64 returns a pointer to a float64 value (helper for tests)
func pointerToFloat64(v float64) *float64 {
	return &v
}

func TestCalcularCharolaEspaciadoUseCase_Execute(t *testing.T) {
	// Tabla de charolas del CSV: charola_dimensiones.csv
	tablaCharola := []valueobject.EntradaTablaCanalizacion{
		{Tamano: "6", AreaInteriorMM2: 152.4},
		{Tamano: "9", AreaInteriorMM2: 228.6},
		{Tamano: "12", AreaInteriorMM2: 304.8},
		{Tamano: "16", AreaInteriorMM2: 406.4},
		{Tamano: "18", AreaInteriorMM2: 457.2},
		{Tamano: "20", AreaInteriorMM2: 508.0},
		{Tamano: "24", AreaInteriorMM2: 609.6},
		{Tamano: "30", AreaInteriorMM2: 762.0},
		{Tamano: "36", AreaInteriorMM2: 914.4},
	}

	tests := []struct {
		name        string
		input       dto.CharolaEspaciadoInput
		mockTabla   []valueobject.EntradaTablaCanalizacion
		mockErr     error
		wantOutput  dto.CharolaEspaciadoOutput
		wantErr     bool
		errContains string
	}{
		{
			name: "happy path - caso básico DELTA 1 hilo",
			// DELTA: 3 fases sin neutro, 1 hilo por fase = 3 hilos totales
			// anchoRequerido = 3*25.4 + 3*25.4 + 8.5 = 160.9mm -> 9" (228.6mm)
			input: dto.CharolaEspaciadoInput{
				HilosPorFase:     1,
				SistemaElectrico: "DELTA",
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla: tablaCharola,
			wantOutput: dto.CharolaEspaciadoOutput{
				Tipo:           "CHAROLA_CABLE_ESPACIADO",
				Tamano:         "9",
				TamanoPulgadas: "9\"",
			},
			wantErr: false,
		},
		{
			name: "happy path - caso MONOFASICO 1 hilo",
			// MONOFASICO: 1 fase con neutro, 1 hilo por fase = 2 hilos totales
			// anchoRequerido = 2*25.4 + 2*25.4 + 8.5 = 109.3mm -> 6" (152.4mm)
			input: dto.CharolaEspaciadoInput{
				HilosPorFase:     1,
				SistemaElectrico: "MONOFASICO",
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla: tablaCharola,
			wantOutput: dto.CharolaEspaciadoOutput{
				Tipo:           "CHAROLA_CABLE_ESPACIADO",
				Tamano:         "6",
				TamanoPulgadas: "6\"",
			},
			wantErr: false,
		},
		{
			name: "happy path - con cable de control",
			// MONOFASICO: 1 fase con neutro, 1 hilo por fase = 2 hilos totales
			// espacioFuerza = 2*10 = 20
			// anchoFuerza = 2*10 = 20
			// espacioControl = 2*6 = 12
			// anchoControl = 6
			// anchoTierra = 5
			// anchoRequerido = 20 + 20 + 12 + 6 + 5 = 63mm -> 6" (152.4mm)
			input: dto.CharolaEspaciadoInput{
				HilosPorFase:      1,
				SistemaElectrico:  "MONOFASICO",
				DiametroFaseMM:    10.0,
				DiametroTierraMM:  5.0,
				DiametroControlMM: pointerToFloat64(6.0),
			},
			mockTabla: tablaCharola,
			wantOutput: dto.CharolaEspaciadoOutput{
				Tipo:           "CHAROLA_CABLE_ESPACIADO",
				Tamano:         "6",
				TamanoPulgadas: "6\"",
			},
			wantErr: false,
		},
		{
			name: "error - hilos por fase cero",
			input: dto.CharolaEspaciadoInput{
				HilosPorFase:     0,
				SistemaElectrico: "DELTA",
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla:   tablaCharola,
			wantErr:     true,
			errContains: "hilos_por_fase",
		},
		{
			name: "error - diámetro fase cero",
			input: dto.CharolaEspaciadoInput{
				HilosPorFase:     1,
				SistemaElectrico: "DELTA",
				DiametroFaseMM:   0,
				DiametroTierraMM: 8.5,
			},
			mockTabla:   tablaCharola,
			wantErr:     true,
			errContains: "diametro_fase_mm",
		},
		{
			name: "error - sistema eléctrico inválido",
			input: dto.CharolaEspaciadoInput{
				HilosPorFase:     1,
				SistemaElectrico: "INVALIDO",
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla:   tablaCharola,
			wantErr:     true,
			errContains: "sistema eléctrico",
		},
		{
			name: "error - tabla vacía",
			input: dto.CharolaEspaciadoInput{
				HilosPorFase:     1,
				SistemaElectrico: "DELTA",
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla:   []valueobject.EntradaTablaCanalizacion{},
			wantErr:     true,
			errContains: "charola",
		},
		{
			name: "error - repositorio falla",
			input: dto.CharolaEspaciadoInput{
				HilosPorFase:     1,
				SistemaElectrico: "DELTA",
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla:   nil,
			mockErr:     errors.New("error de repositorio"),
			wantErr:     true,
			errContains: "repositorio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &mockCharolaRepo{
				tablaCharola: tt.mockTabla,
				tablaErr:     tt.mockErr,
			}
			uc := NewCalcularCharolaEspaciadoUseCase(mockRepo)

			// Act
			output, err := uc.Execute(context.Background(), tt.input)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantOutput.Tipo, output.Tipo)
			assert.Equal(t, tt.wantOutput.Tamano, output.Tamano)
			assert.Equal(t, tt.wantOutput.TamanoPulgadas, output.TamanoPulgadas)
		})
	}
}
