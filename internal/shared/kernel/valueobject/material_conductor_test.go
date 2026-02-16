// internal/shared/kernel/valueobject/material_conductor_test.go
package valueobject_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaterialConductor_String(t *testing.T) {
	tests := []struct {
		name     string
		material valueobject.MaterialConductor
		want     string
	}{
		{"cobre", valueobject.MaterialCobre, "CU"},
		{"aluminio", valueobject.MaterialAluminio, "AL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.material.String())
		})
	}
}

func TestParseMaterialConductor_validos(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected valueobject.MaterialConductor
	}{
		{"Cu", "Cu", valueobject.MaterialCobre},
		{"CU", "CU", valueobject.MaterialCobre},
		{"cu", "cu", valueobject.MaterialCobre},
		{"cobre", "cobre", valueobject.MaterialCobre},
		{"COBRE", "COBRE", valueobject.MaterialCobre},
		{"Al", "Al", valueobject.MaterialAluminio},
		{"AL", "AL", valueobject.MaterialAluminio},
		{"al", "al", valueobject.MaterialAluminio},
		{"aluminio", "aluminio", valueobject.MaterialAluminio},
		{"ALUMINIO", "ALUMINIO", valueobject.MaterialAluminio},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valueobject.ParseMaterialConductor(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParseMaterialConductor_invalido(t *testing.T) {
	_, err := valueobject.ParseMaterialConductor("HIERRO")
	require.Error(t, err)
	assert.ErrorIs(t, err, valueobject.ErrMaterialConductorInvalido)
}
