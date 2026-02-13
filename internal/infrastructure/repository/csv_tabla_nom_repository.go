// internal/infrastructure/repository/csv_tabla_nom_repository.go
package repository

import (
	"fmt"
	"os"
)

// CSVTablaNOMRepository reads NOM tables from CSV files with in-memory caching.
type CSVTablaNOMRepository struct {
	basePath string
}

// NewCSVTablaNOMRepository creates a new repository and validates the base path.
// Tables are loaded into memory on first use (lazy loading).
func NewCSVTablaNOMRepository(basePath string) (*CSVTablaNOMRepository, error) {
	// Verify directory exists
	info, err := os.Stat(basePath)
	if err != nil {
		return nil, fmt.Errorf("cannot access base path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("base path is not a directory: %s", basePath)
	}

	return &CSVTablaNOMRepository{
		basePath: basePath,
	}, nil
}
