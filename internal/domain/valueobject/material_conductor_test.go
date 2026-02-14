// internal/domain/valueobject/material_conductor_test.go
package valueobject

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestMaterialConductor_UnmarshalJSON(t *testing.T) {
	type Input struct {
		Material MaterialConductor `json:"material"`
	}

	tests := []struct {
		name     string
		json     string
		expected MaterialConductor
	}{
		{"empty", `{}`, MaterialCobre},
		{"Cu uppercase", `{"material":"Cu"}`, MaterialCobre},
		{"Cu lower", `{"material":"cu"}`, MaterialCobre},
		{"Al uppercase", `{"material":"Al"}`, MaterialAluminio},
		{"Al lower", `{"material":"al"}`, MaterialAluminio},
		{"number 0", `{"material":"0"}`, MaterialCobre},
		{"number 1", `{"material":"1"}`, MaterialAluminio},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input Input
			err := json.Unmarshal([]byte(tt.json), &input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, input.Material)
		})
	}
}
