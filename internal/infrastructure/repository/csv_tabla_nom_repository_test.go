// internal/infrastructure/repository/csv_tabla_nom_repository_test.go
package repository

import (
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
