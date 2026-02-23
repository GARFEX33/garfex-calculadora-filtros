// internal/shared/kernel/valueobject/tension_test.go
package valueobject_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTension(t *testing.T) {
	tests := []struct {
		name    string
		valor   float64
		unidad  string
		wantErr bool
	}{
		{"127V valid", 127, "V", false},
		{"220V valid", 220, "V", false},
		{"240V valid", 240, "V", false},
		{"277V valid", 277, "V", false},
		{"440V valid", 440, "V", false},
		{"480V valid", 480, "V", false},
		{"600V valid", 600, "V", false},
		{"100V invalid", 100, "V", true},
		{"0V invalid", 0, "V", true},
		{"negative invalid", -220, "V", true},
		// Tests with kV
		{"0.48kV valid (480V)", 0.48, "kV", false},
		{"0.22kV valid (220V)", 0.22, "kV", false},
		{"0.127kV valid (127V)", 0.127, "kV", false},
		{"0.5kV invalid (500V)", 0.5, "kV", true},
		// Default empty unit
		{"empty unit defaults to V", 480, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tension, err := valueobject.NewTension(tt.valor, tt.unidad)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, valueobject.ErrVoltajeInvalido))
				return
			}
			require.NoError(t, err)
			// Verify normalized value in volts
			expectedVolts := int(tt.valor)
			if tt.unidad == "kV" {
				expectedVolts = int(tt.valor * 1000)
			}
			assert.Equal(t, expectedVolts, tension.Valor())
		})
	}
}

func TestTension_EnKilovoltios(t *testing.T) {
	tests := []struct {
		name     string
		voltaje  float64
		unidad   string
		expected float64
	}{
		{"480V → 0.48 kV", 480, "V", 0.48},
		{"220V → 0.22 kV", 220, "V", 0.22},
		{"127V → 0.127 kV", 127, "V", 0.127},
		{"0.48kV → 0.48 kV", 0.48, "kV", 0.48},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tension, err := valueobject.NewTension(tt.voltaje, tt.unidad)
			require.NoError(t, err)
			assert.InDelta(t, tt.expected, tension.EnKilovoltios(), 0.001)
		})
	}
}
