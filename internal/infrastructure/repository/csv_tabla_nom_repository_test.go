// internal/infrastructure/repository/csv_tabla_nom_repository_test.go
package repository

import (
	"context"
	"testing"

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
