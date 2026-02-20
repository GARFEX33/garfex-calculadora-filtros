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
			// anchoPotencia = 2 * 25.4 * 1 = 50.8
			// espacioFuerza = (1-1) * 2.15 * 25.4 = 0
			// anchoRequerido = 50.8 + 0 + 8.5 = 59.3 mm
			// 6" (152.4mm) > 59.3 -> 6"
			input: dto.CharolaTriangularInput{
				HilosPorFase:     1,
				DiametroFaseMM:   25.4,
				DiametroTierraMM: 8.5,
			},
			mockTabla: tablaCharola,
			wantOutput: dto.CharolaTriangularOutput{
				Tipo:           "CHAROLA_CABLE_TRIANGULAR",
				Tamano:         "6",
				TamanoPulgadas: "6\"",
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
			// 9" (228.6mm) > 196.21 -> 9"
			input: dto.CharolaTriangularInput{
				HilosPorFase:      2,
				DiametroFaseMM:    25.4,
				DiametroTierraMM:  8.5,
				DiametroControlMM: func() *float64 { v := 10.0; return &v }(),
			},
			mockTabla: tablaCharola,
			wantOutput: dto.CharolaTriangularOutput{
				Tipo:           "CHAROLA_CABLE_TRIANGULAR",
				Tamano:         "9",
				TamanoPulgadas: "9\"",
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
