// internal/domain/valueobject/material_conductor_test.go
package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaterialConductor_String(t *testing.T) {
	tests := []struct {
		name     string
		material MaterialConductor
		want     string
	}{
		{"cobre", MaterialCobre, "CU"},
		{"aluminio", MaterialAluminio, "AL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.material.String())
		})
	}
}
