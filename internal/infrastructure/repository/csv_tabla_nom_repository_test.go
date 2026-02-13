// internal/infrastructure/repository/csv_tabla_nom_repository_test.go
package repository

import (
	"context"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCSVTablaNOMRepository(t *testing.T) {
	t.Run("successfully loads all tables", func(t *testing.T) {
		repo, err := NewCSVTablaNOMRepository("testdata")
		require.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("fails with invalid path", func(t *testing.T) {
		repo, err := NewCSVTablaNOMRepository("nonexistent")
		assert.Error(t, err)
		assert.Nil(t, repo)
	})
}

func TestCSVTablaNOMRepository_ObtenerTablaTierra(t *testing.T) {
	repo, err := NewCSVTablaNOMRepository("testdata")
	require.NoError(t, err)

	ctx := context.Background()
	tabla, err := repo.ObtenerTablaTierra(ctx)

	require.NoError(t, err)
	assert.Greater(t, len(tabla), 0)

	// Check first entry
	assert.Equal(t, 15, tabla[0].ITMHasta)
	assert.Equal(t, "14 AWG", tabla[0].Conductor.Calibre)
}

func TestCSVTablaNOMRepository_ObtenerTablaAmpacidad(t *testing.T) {
	repo, err := NewCSVTablaNOMRepository("testdata")
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name         string
		canalizacion entity.TipoCanalizacion
		material     valueobject.MaterialConductor
		temp         valueobject.Temperatura
	}{
		{"PVC Copper 75C", entity.TipoCanalizacionTuberiaPVC, valueobject.MaterialCobre, valueobject.Temp75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabla, err := repo.ObtenerTablaAmpacidad(ctx, tt.canalizacion, tt.material, tt.temp)
			require.NoError(t, err)
			assert.Greater(t, len(tabla), 0)

			// Check that entries have capacity values
			assert.Greater(t, tabla[0].Capacidad, 0.0)
		})
	}
}

func TestCSVTablaNOMRepository_ObtenerImpedancia(t *testing.T) {
	repo, err := NewCSVTablaNOMRepository("testdata")
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name         string
		calibre      string
		canalizacion entity.TipoCanalizacion
		material     valueobject.MaterialConductor
		wantR        float64
		wantX        float64
	}{
		{"14 AWG PVC Copper", "14 AWG", entity.TipoCanalizacionTuberiaPVC, valueobject.MaterialCobre, 10.2, 0.19},
		{"12 AWG PVC Copper", "12 AWG", entity.TipoCanalizacionTuberiaPVC, valueobject.MaterialCobre, 6.6, 0.177},
		{"6 AWG PVC Aluminium", "6 AWG", entity.TipoCanalizacionTuberiaPVC, valueobject.MaterialAluminio, 2.66, 0.167},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imp, err := repo.ObtenerImpedancia(ctx, tt.calibre, tt.canalizacion, tt.material)
			require.NoError(t, err)
			assert.InDelta(t, tt.wantR, imp.R, 0.01, "R mismatch")
			assert.InDelta(t, tt.wantX, imp.X, 0.01, "X mismatch")
		})
	}
}

func TestCSVTablaNOMRepository_ObtenerTablaCanalizacion(t *testing.T) {
	repo, err := NewCSVTablaNOMRepository("testdata")
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name         string
		canalizacion entity.TipoCanalizacion
		minEntries   int
	}{
		{"PVC", entity.TipoCanalizacionTuberiaPVC, 10},
		{"Aluminio", entity.TipoCanalizacionTuberiaAluminio, 10},
		{"Acero PG", entity.TipoCanalizacionTuberiaAceroPG, 10},
		{"Acero PD", entity.TipoCanalizacionTuberiaAceroPD, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabla, err := repo.ObtenerTablaCanalizacion(ctx, tt.canalizacion)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(tabla), tt.minEntries)

			// Check first entry has valid data
			assert.NotEmpty(t, tabla[0].Tamano)
			assert.Greater(t, tabla[0].AreaInteriorMM2, 0.0)
		})
	}
}
