// internal/calculos/application/usecase/calcular_charola_triangular_test.go
package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestCalcularCharolaTriangularUseCase_Execute(t *testing.T) {
	// Tabla de charolas triangular de prueba
	// Nota: el servicio de triangular compara AreaInteriorMM2 directamente con anchoRequerido
	tablaCharola := []valueobject.EntradaTablaCanalizacion{
		{Tamano: "100mm", AreaInteriorMM2: 5000},
		{Tamano: "150mm", AreaInteriorMM2: 7500},
		{Tamano: "200mm", AreaInteriorMM2: 10000},
		{Tamano: "300mm", AreaInteriorMM2: 15000},
	}

	tests := []struct {
		name        string
		input       dto.CharolaTriangularInput
		mockTabla   []valueobject.EntradaTablaCanalizacion
		mockErr     error
		wantOutput  dto.CharolaTriangularOutput
		wantErr     bool
		errContains string
	}{
		{
			name: "happy path - caso básico",
			// triangular con 1 hilo por fase:
			// anchoPotencia = 2 * 25.4 = 50.8
			// espacioFuerza = (1-1) * 2.15 * 25.4 = 0
			// anchoRequerido = 50.8 + 0 + 0 + 8.5 = 59.3 mm
			// 5000 (100mm) > 59.3 -> 100mm
			input: dto.CharolaTriangularInput{
				HilosPorFase:     1,
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla: tablaCharola,
			wantOutput: dto.CharolaTriangularOutput{
				Tipo:           "CHAROLA_CABLE_TRIANGULAR",
				Tamano:         "100mm",
				TamanoPulgadas: "3.94\"",
			},
			wantErr: false,
		},
		{
			name: "happy path - con cable de control",
			// triangular con 2 hilos por fase y cable de control:
			// anchoPotencia = 2 * 25.4 * 2 = 101.6
			// espacioFuerza = (2-1) * 2.15 * 25.4 = 54.61
			// espacioControl = 2.15 * 10 = 21.5
			// anchoControl = 10
			// anchoRequerido = 101.6 + 54.61 + 21.5 + 10 + 8.5 = 196.21 mm
			// 10000 (200mm) > 196.21 -> 200mm
			input: dto.CharolaTriangularInput{
				HilosPorFase:      2,
				DiametroFaseMM:    25.4,
				DiametroTierraMM:  8.5,
				DiametroControlMM: func() *float64 { v := 10.0; return &v }(),
			},
			mockTabla: tablaCharola,
			wantOutput: dto.CharolaTriangularOutput{
				Tipo:           "CHAROLA_CABLE_TRIANGULAR",
				Tamano:         "200mm",
				TamanoPulgadas: "7.87\"",
			},
			wantErr: false,
		},
		{
			name: "error - hilos por fase cero",
			input: dto.CharolaTriangularInput{
				HilosPorFase:     0,
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla:   tablaCharola,
			wantErr:     true,
			errContains: "hilos_por_fase",
		},
		{
			name: "error - diámetro fase cero",
			input: dto.CharolaTriangularInput{
				HilosPorFase:     1,
				DiametroFaseMM:   0,
				DiametroTierraMM: 8.5,
			},
			mockTabla:   tablaCharola,
			wantErr:     true,
			errContains: "diametro_fase_mm",
		},
		{
			name: "error - tabla vacía",
			input: dto.CharolaTriangularInput{
				HilosPorFase:     1,
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla:   []valueobject.EntradaTablaCanalizacion{},
			wantErr:     true,
			errContains: "charola",
		},
		{
			name: "error - repositorio falla",
			input: dto.CharolaTriangularInput{
				HilosPorFase:     1,
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
			uc := NewCalcularCharolaTriangularUseCase(mockRepo)

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
