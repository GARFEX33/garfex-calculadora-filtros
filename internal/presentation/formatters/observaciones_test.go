// internal/presentation/formatters/observaciones_test.go
package formatters

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
)

func TestGenerarObservaciones(t *testing.T) {
	tests := []struct {
		name     string
		output   dto.MemoriaOutput
		expected []string
	}{
		{
			name: "caida tension excedida",
			output: dto.MemoriaOutput{
				CaidaTension: dto.ResultadoCaidaTension{
					Porcentaje:       4.5,
					LimitePorcentaje: 3.0,
					Cumple:           false,
				},
				HilosPorFase: 1,
			},
			expected: []string{"Caída de tensión 4.50% excede el límite de 3.00%"},
		},
		{
			name: "hilos en paralelo",
			output: dto.MemoriaOutput{
				CaidaTension: dto.ResultadoCaidaTension{
					Cumple: true,
				},
				HilosPorFase: 3,
			},
			expected: []string{"Se usan 3 hilos por fase en paralelo"},
		},
		{
			name: "sin observaciones",
			output: dto.MemoriaOutput{
				CaidaTension: dto.ResultadoCaidaTension{
					Cumple: true,
				},
				HilosPorFase: 1,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerarObservaciones(tt.output)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d observations, got %d", len(tt.expected), len(result))
			}
			for i, obs := range tt.expected {
				if i < len(result) && result[i] != obs {
					t.Errorf("expected %s, got %s", obs, result[i])
				}
			}
		})
	}
}
