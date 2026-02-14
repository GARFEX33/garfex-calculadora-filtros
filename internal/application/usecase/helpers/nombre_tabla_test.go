package helpers

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

func TestNombreTablaAmpacidad(t *testing.T) {
	tests := []struct {
		name         string
		canalizacion string
		material     valueobject.MaterialConductor
		temperatura  valueobject.Temperatura
		expected     string
	}{
		{
			name:         "PVC con cobre 75C",
			canalizacion: "TUBERIA_PVC",
			material:     valueobject.MaterialCobre,
			temperatura:  valueobject.Temp75,
			expected:     "NOM-310-15-B-16 (Cu, 75°C)",
		},
		{
			name:         "Charola triangular con aluminio 90C",
			canalizacion: "CHAROLA_CABLE_TRIANGULAR",
			material:     valueobject.MaterialAluminio,
			temperatura:  valueobject.Temp90,
			expected:     "NOM-310-15-B-20 (Al, 90°C)",
		},
		{
			name:         "Acero PG con cobre 60C",
			canalizacion: "TUBERIA_ACERO_PG",
			material:     valueobject.MaterialCobre,
			temperatura:  valueobject.Temp60,
			expected:     "NOM-310-15-B-16 (Cu, 60°C)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NombreTablaAmpacidad(tt.canalizacion, tt.material, tt.temperatura)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
