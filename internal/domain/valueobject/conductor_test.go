// internal/domain/valueobject/conductor_test.go
package valueobject_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConductor(t *testing.T) {
	tests := []struct {
		name            string
		calibre         string
		material        string
		tipoAislamiento string
		seccionMM2      float64
		wantErr         bool
	}{
		{"valid Cu THHN 12AWG", "12 AWG", "Cu", "THHN", 3.31, false},
		{"valid Al THW 4/0AWG", "4/0 AWG", "Al", "THW", 107.2, false},
		{"valid Cu XHHW 500MCM", "500 MCM", "Cu", "XHHW", 253.4, false},
		{"empty calibre", "", "Cu", "THHN", 3.31, true},
		{"invalid material", "12 AWG", "Fe", "THHN", 3.31, true},
		{"empty material", "12 AWG", "", "THHN", 3.31, true},
		{"empty aislamiento", "12 AWG", "Cu", "", 3.31, true},
		{"zero seccion", "12 AWG", "Cu", "THHN", 0, true},
		{"negative seccion", "12 AWG", "Cu", "THHN", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := valueobject.NewConductor(tt.calibre, tt.material, tt.tipoAislamiento, tt.seccionMM2)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, valueobject.ErrConductorInvalido))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.calibre, c.Calibre())
			assert.Equal(t, tt.material, c.Material())
			assert.Equal(t, tt.tipoAislamiento, c.TipoAislamiento())
			assert.Equal(t, tt.seccionMM2, c.SeccionMM2())
		})
	}
}
