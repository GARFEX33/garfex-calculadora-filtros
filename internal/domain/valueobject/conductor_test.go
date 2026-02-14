// internal/domain/valueobject/conductor_test.go
package valueobject_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// conductor12AWGCu returns a valid ConductorParams for 12 AWG Cu THHN
// based on NOM-001-SEDE-2012 tables.
func conductor12AWGCu() valueobject.ConductorParams {
	return valueobject.ConductorParams{
		Calibre:               "12 AWG",
		Material:              valueobject.MaterialCobre,
		TipoAislamiento:       "THHN",
		SeccionMM2:            3.31,
		AreaConAislamientoMM2: 11.68,
		DiametroMM:            3.861,
		NumeroHilos:           7,
		ResistenciaPVCPorKm:   6.6,
		ResistenciaAlPorKm:    6.6,
		ResistenciaAceroPorKm: 6.6,
		ReactanciaPorKm:       0.177,
	}
}

func TestNewConductor_Valid(t *testing.T) {
	c, err := valueobject.NewConductor(conductor12AWGCu())
	require.NoError(t, err)

	assert.Equal(t, "12 AWG", c.Calibre())
	assert.Equal(t, valueobject.MaterialCobre, c.Material())
	assert.Equal(t, "THHN", c.TipoAislamiento())
	assert.Equal(t, 3.31, c.SeccionMM2())
	assert.Equal(t, 11.68, c.AreaConAislamientoMM2())
	assert.Equal(t, 3.861, c.DiametroMM())
	assert.Equal(t, 7, c.NumeroHilos())
	assert.Equal(t, 6.6, c.ResistenciaPVCPorKm())
	assert.Equal(t, 6.6, c.ResistenciaAlPorKm())
	assert.Equal(t, 6.6, c.ResistenciaAceroPorKm())
	assert.Equal(t, 0.177, c.ReactanciaPorKm())
}

func TestNewConductor_CalibreInvalido(t *testing.T) {
	tests := []struct {
		name    string
		calibre string
	}{
		{"empty", ""},
		{"not in NOM", "3 AWG"},
		{"without suffix", "12"},
		{"wrong format", "12AWG"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := conductor12AWGCu()
			p.Calibre = tt.calibre
			_, err := valueobject.NewConductor(p)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, valueobject.ErrConductorInvalido))
		})
	}
}

func TestNewConductor_MaterialInvalido(t *testing.T) {
	p := conductor12AWGCu()
	p.Material = valueobject.MaterialConductor(999) // Invalid material
	_, err := valueobject.NewConductor(p)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, valueobject.ErrConductorInvalido))
}

func TestNewConductor_SeccionCero(t *testing.T) {
	p := conductor12AWGCu()
	p.SeccionMM2 = 0
	_, err := valueobject.NewConductor(p)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, valueobject.ErrConductorInvalido))
}

func TestNewConductor_SeccionNegativa(t *testing.T) {
	p := conductor12AWGCu()
	p.SeccionMM2 = -1
	_, err := valueobject.NewConductor(p)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, valueobject.ErrConductorInvalido))
}

func TestNewConductor_Minimal(t *testing.T) {
	// Bare conductor with only required fields (Calibre, Material, SeccionMM2).
	// All other fields are zero-value — valid for ground conductors, etc.
	c, err := valueobject.NewConductor(valueobject.ConductorParams{
		Calibre:    "8 AWG",
		Material:   valueobject.MaterialCobre,
		SeccionMM2: 8.37,
	})
	require.NoError(t, err)

	assert.Equal(t, "8 AWG", c.Calibre())
	assert.Equal(t, valueobject.MaterialCobre, c.Material())
	assert.Equal(t, "", c.TipoAislamiento())
	assert.Equal(t, 8.37, c.SeccionMM2())
	assert.Equal(t, 0.0, c.AreaConAislamientoMM2())
	assert.Equal(t, 0.0, c.DiametroMM())
	assert.Equal(t, 0, c.NumeroHilos())
}

func TestNewConductor_ExtremosCalibre(t *testing.T) {
	base := conductor12AWGCu()

	// Extremo inferior: 14 AWG
	base.Calibre = "14 AWG"
	_, err := valueobject.NewConductor(base)
	assert.NoError(t, err)

	// Extremo superior: 1000 MCM
	base.Calibre = "1000 MCM"
	_, err = valueobject.NewConductor(base)
	assert.NoError(t, err)

	// Calibres eliminados: ahora inválidos
	base.Calibre = "18 AWG"
	_, err = valueobject.NewConductor(base)
	assert.Error(t, err)

	base.Calibre = "16 AWG"
	_, err = valueobject.NewConductor(base)
	assert.Error(t, err)

	base.Calibre = "3 AWG"
	_, err = valueobject.NewConductor(base)
	assert.Error(t, err)

	base.Calibre = "1 AWG"
	_, err = valueobject.NewConductor(base)
	assert.Error(t, err)

	base.Calibre = "2000 MCM"
	_, err = valueobject.NewConductor(base)
	assert.Error(t, err)

	base.Calibre = "750 MCM"
	_, err = valueobject.NewConductor(base)
	assert.NoError(t, err)
}
