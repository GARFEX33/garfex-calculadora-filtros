// internal/calculos/application/usecase/calcular_amperaje_nominal_test.go
package usecase

import (
	"context"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/application/dto"
	"github.com/garfex/calculadora-filtros/internal/calculos/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularAmperajeNominalUseCase_Execute(t *testing.T) {
	// Arrange
	uc := NewCalcularAmperajeNominalUseCase()
	ctx := context.Background()

	tests := []struct {
		name           string
		input          dto.AmperajeNominalInput
		expectedAmp    float64
		expectedUnidad string
		expectedErr    bool
	}{
		{
			name: "caso exitoso - carga monofásica",
			input: dto.AmperajeNominalInput{
				PotenciaWatts:    10000, // 10 kW
				Tension:          220,
				TipoCarga:        dto.TipoCargaDTOMonofasica,
				SistemaElectrico: dto.SistemaElectricoDTOEstrella,
				FactorPotencia:   0.9,
			},
			expectedAmp:    50.5, // I = 10000 / (220 * 0.9) ≈ 50.5 A
			expectedUnidad: "A",
			expectedErr:    false,
		},
		{
			name: "caso exitoso - carga trifásica",
			input: dto.AmperajeNominalInput{
				PotenciaWatts:    30000, // 30 kW
				Tension:          220,
				TipoCarga:        dto.TipoCargaDTOTrifasica,
				SistemaElectrico: dto.SistemaElectricoDTODelta,
				FactorPotencia:   0.9,
			},
			expectedAmp:    87.4, // I = 30000 / (220 * 1.732 * 0.9) ≈ 87.4 A
			expectedUnidad: "A",
			expectedErr:    false,
		},
		{
			name: "error - potencia cero",
			input: dto.AmperajeNominalInput{
				PotenciaWatts:    0,
				Tension:          220,
				TipoCarga:        dto.TipoCargaDTOMonofasica,
				SistemaElectrico: dto.SistemaElectricoDTOEstrella,
				FactorPotencia:   0.9,
			},
			expectedErr: true,
		},
		{
			name: "error - tensión inválida",
			input: dto.AmperajeNominalInput{
				PotenciaWatts:    10000,
				Tension:          100, // No válido según NOM
				TipoCarga:        dto.TipoCargaDTOMonofasica,
				SistemaElectrico: dto.SistemaElectricoDTOEstrella,
				FactorPotencia:   0.9,
			},
			expectedErr: true,
		},
		{
			name: "error - factor de potencia cero",
			input: dto.AmperajeNominalInput{
				PotenciaWatts:    10000,
				Tension:          220,
				TipoCarga:        dto.TipoCargaDTOMonofasica,
				SistemaElectrico: dto.SistemaElectricoDTOEstrella,
				FactorPotencia:   0,
			},
			expectedErr: true,
		},
		{
			name: "error - factor de potencia mayor a 1",
			input: dto.AmperajeNominalInput{
				PotenciaWatts:    10000,
				Tension:          220,
				TipoCarga:        dto.TipoCargaDTOMonofasica,
				SistemaElectrico: dto.SistemaElectricoDTOEstrella,
				FactorPotencia:   1.5,
			},
			expectedErr: true,
		},
		{
			name: "caso exitoso - tensión 127V monofásico",
			input: dto.AmperajeNominalInput{
				PotenciaWatts:    1500, // 1.5 kW
				Tension:          127,
				TipoCarga:        dto.TipoCargaDTOMonofasica,
				SistemaElectrico: dto.SistemaElectricoDTOEstrella,
				FactorPotencia:   0.85,
			},
			expectedAmp:    13.8, // I = 1500 / (127 * 0.85) ≈ 13.8 A
			expectedUnidad: "A",
			expectedErr:    false,
		},
		{
			name: "caso exitoso - tensión 480V trifásico",
			input: dto.AmperajeNominalInput{
				PotenciaWatts:    100000, // 100 kW
				Tension:          480,
				TipoCarga:        dto.TipoCargaDTOTrifasica,
				SistemaElectrico: dto.SistemaElectricoDTOEstrella,
				FactorPotencia:   0.95,
			},
			expectedAmp:    126.6, // I = 100000 / (480 * 1.732 * 0.95) ≈ 126.6 A
			expectedUnidad: "A",
			expectedErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := uc.Execute(ctx, tt.input)

			// Assert
			if tt.expectedErr {
				require.Error(t, err, "se esperaba error pero no ocurrió")
				return
			}

			require.NoError(t, err, "no se esperaba error pero ocurrió: %v", err)
			assert.InDelta(t, tt.expectedAmp, result.Amperaje, 0.5, "amperaje fuera de rango")
			assert.Equal(t, tt.expectedUnidad, result.Unidad, "unidad incorrecta")
		})
	}
}

func TestTipoCargaDTO_ToEntity(t *testing.T) {
	tests := []struct {
		dto      dto.TipoCargaDTO
		expected service.TipoCarga
	}{
		{dto.TipoCargaDTOMonofasica, service.TipoCargaMonofasica},
		{dto.TipoCargaDTOTrifasica, service.TipoCargaTrifasica},
		{dto.TipoCargaDTOGenerica, service.TipoCargaTrifasica}, // default
		{"", service.TipoCargaTrifasica},                       // default
	}

	for _, tt := range tests {
		result := tt.dto.ToEntity()
		assert.Equal(t, tt.expected, result)
	}
}

func TestSistemaElectricoDTO_ToEntity(t *testing.T) {
	tests := []struct {
		dto      dto.SistemaElectricoDTO
		expected service.SistemaElectrico
	}{
		{dto.SistemaElectricoDTOEstrella, service.SistemaElectricoEstrella},
		{dto.SistemaElectricoDTODelta, service.SistemaElectricoDelta},
		{dto.SistemaElectricoDTO("OTRO"), service.SistemaElectricoEstrella}, // default
		{"", service.SistemaElectricoEstrella},                              // default
	}

	for _, tt := range tests {
		result := tt.dto.ToEntity()
		assert.Equal(t, tt.expected, result)
	}
}
