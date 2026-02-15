// internal/calculos/infrastructure/adapter/driven/csv/csv_tabla_nom_repository_test.go
package csv

import (
	"context"
	"fmt"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/calculos/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/shared/kernel/valueobject"
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

	// Check first entry (ITM hasta 1)
	assert.Equal(t, 1, tabla[0].ITMHasta)
	assert.Equal(t, "14 AWG", tabla[0].ConductorCu.Calibre)

	// Check second entry (ITM hasta 15)
	assert.Equal(t, 15, tabla[1].ITMHasta)
	assert.Equal(t, "14 AWG", tabla[1].ConductorCu.Calibre)
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

func TestCSVTablaNOMRepository_ObtenerTemperaturaPorEstado(t *testing.T) {
	repo, err := NewCSVTablaNOMRepository("testdata")
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		estado      string
		expectedMin int // temperatura mínima esperada (para verificar que no es el valor viejo)
		expectedMax int
		wantErr     bool
	}{
		{"Nuevo Leon", 36, 38, false},       // 36.8 → round → 37
		{"Sonora", 38, 39, false},           // 38.2 → round → 38
		{"Tlaxcala", 27, 28, false},         // 27.5 → round → 28
		{"Ciudad de Mexico", 29, 30, false}, // 29.1 → round → 29
		{"Estado Inexistente", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.estado, func(t *testing.T) {
			temp, err := repo.ObtenerTemperaturaPorEstado(ctx, tt.estado)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.GreaterOrEqual(t, temp, tt.expectedMin, "temperatura demasiado baja — posible dato incorrecto")
			assert.LessOrEqual(t, temp, tt.expectedMax, "temperatura demasiado alta — posible dato incorrecto")
		})
	}
}

func TestCSVTablaNOMRepository_ObtenerFactorAgrupamiento(t *testing.T) {
	repo, err := NewCSVTablaNOMRepository("testdata")
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		cantidad int
		expected float64
	}{
		{1, 1.00},
		{2, 0.80}, // bug anterior: retornaba 0.30
		{3, 0.70}, // bug anterior: retornaba 0.30
		{4, 0.65}, // bug anterior: retornaba 0.30
		{5, 0.60},
		{6, 0.60},
		{7, 0.50},
		{10, 0.45},
		{21, 0.40},
		{31, 0.35},
		{41, 0.30},
		{50, 0.30},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_conductores", tt.cantidad), func(t *testing.T) {
			factor, err := repo.ObtenerFactorAgrupamiento(ctx, tt.cantidad)
			require.NoError(t, err)
			assert.InDelta(t, tt.expected, factor, 0.001)
		})
	}
}

func TestCSVTablaNOMRepository_ObtenerFactorTemperatura(t *testing.T) {
	repo, err := NewCSVTablaNOMRepository("testdata")
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		tempAmbiente  int
		tempConductor valueobject.Temperatura
		expected      float64
	}{
		{37, valueobject.Temp60, 0.94}, // Nuevo Leon real: rango 36-40, 60C
		{37, valueobject.Temp75, 0.95}, // rango 36-40, 75C
		{31, valueobject.Temp75, 1.00}, // rango 31-35, base
		{21, valueobject.Temp75, 1.07}, // rango 21-25
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d°C_%dC", tt.tempAmbiente, tt.tempConductor), func(t *testing.T) {
			factor, err := repo.ObtenerFactorTemperatura(ctx, tt.tempAmbiente, tt.tempConductor)
			require.NoError(t, err)
			assert.InDelta(t, tt.expected, factor, 0.001)
		})
	}
}
