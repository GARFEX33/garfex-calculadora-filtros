package helpers

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestNombreTablaAmpacidad(t *testing.T) {
	tests := []struct {
		name         string
		canalizacion string
		material     valueobject.MaterialConductor
		temperatura  valueobject.Temperatura
		want         string
	}{
		{
			name:         "PVC con cobre a 75°C",
			canalizacion: "TUBERIA_PVC",
			material:     valueobject.MaterialCobre,
			temperatura:  valueobject.Temp75,
			want:         "NOM-310-15-B-16 (Cu, 75°C)",
		},
		{
			name:         "PVC con aluminio a 60°C",
			canalizacion: "TUBERIA_PVC",
			material:     valueobject.MaterialAluminio,
			temperatura:  valueobject.Temp60,
			want:         "NOM-310-15-B-16 (Al, 60°C)",
		},
		{
			name:         "TUBERIA_ALUMINIO con cobre a 90°C",
			canalizacion: "TUBERIA_ALUMINIO",
			material:     valueobject.MaterialCobre,
			temperatura:  valueobject.Temp90,
			want:         "NOM-310-15-B-16 (Cu, 90°C)",
		},
		{
			name:         "CHAROLA_CABLE_ESPACIADO con cobre a 75°C",
			canalizacion: "CHAROLA_CABLE_ESPACIADO",
			material:     valueobject.MaterialCobre,
			temperatura:  valueobject.Temp75,
			want:         "NOM-310-15-B-17 (Cu, 75°C)",
		},
		{
			name:         "CHAROLA_CABLE_TRIANGULAR con aluminio a 60°C",
			canalizacion: "CHAROLA_CABLE_TRIANGULAR",
			material:     valueobject.MaterialAluminio,
			temperatura:  valueobject.Temp60,
			want:         "NOM-310-15-B-20 (Al, 60°C)",
		},
		{
			name:         "canalización desconocida usa default PVC",
			canalizacion: "DESCONOCIDA",
			material:     valueobject.MaterialCobre,
			temperatura:  valueobject.Temp75,
			want:         "NOM-310-15-B-16 (Cu, 75°C)",
		},
		{
			name:         "TUBERIA_ACERO_PG con cobre a 60°C",
			canalizacion: "TUBERIA_ACERO_PG",
			material:     valueobject.MaterialCobre,
			temperatura:  valueobject.Temp60,
			want:         "NOM-310-15-B-16 (Cu, 60°C)",
		},
		{
			name:         "TUBERIA_ACERO_PD con aluminio a 90°C",
			canalizacion: "TUBERIA_ACERO_PD",
			material:     valueobject.MaterialAluminio,
			temperatura:  valueobject.Temp90,
			want:         "NOM-310-15-B-16 (Al, 90°C)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NombreTablaAmpacidad(tt.canalizacion, tt.material, tt.temperatura)
			assert.Equal(t, tt.want, got)
		})
	}
}
