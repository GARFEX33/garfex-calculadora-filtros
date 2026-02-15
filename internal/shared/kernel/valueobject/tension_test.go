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
		valor   int
		wantErr bool
	}{
		{"127V valid", 127, false},
		{"220V valid", 220, false},
		{"240V valid", 240, false},
		{"277V valid", 277, false},
		{"440V valid", 440, false},
		{"480V valid", 480, false},
		{"600V valid", 600, false},
		{"100V invalid", 100, true},
		{"0V invalid", 0, true},
		{"negative invalid", -220, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tension, err := valueobject.NewTension(tt.valor)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, valueobject.ErrVoltajeInvalido))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.valor, tension.Valor())
			assert.Equal(t, "V", tension.Unidad())
		})
	}
}

func TestTension_EnKilovoltios(t *testing.T) {
	tests := []struct {
		name     string
		voltaje  int
		expected float64
	}{
		{"480V → 0.48 kV", 480, 0.48},
		{"220V → 0.22 kV", 220, 0.22},
		{"127V → 0.127 kV", 127, 0.127},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tension, err := valueobject.NewTension(tt.voltaje)
			require.NoError(t, err)
			assert.InDelta(t, tt.expected, tension.EnKilovoltios(), 0.001)
		})
	}
}
