// internal/infrastructure/repository/csv_tabla_nom_repository.go
package repository

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// CSVTablaNOMRepository reads NOM tables from CSV files with in-memory caching.
type CSVTablaNOMRepository struct {
	basePath    string
	tablaTierra []service.EntradaTablaTierra
}

// NewCSVTablaNOMRepository creates a new repository and loads all tables into memory.
func NewCSVTablaNOMRepository(basePath string) (*CSVTablaNOMRepository, error) {
	// Verify directory exists
	info, err := os.Stat(basePath)
	if err != nil {
		return nil, fmt.Errorf("cannot access base path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("base path is not a directory: %s", basePath)
	}

	repo := &CSVTablaNOMRepository{
		basePath: basePath,
	}

	// Load ground conductor table
	tablaTierra, err := repo.loadTablaTierra()
	if err != nil {
		return nil, fmt.Errorf("failed to load ground table: %w", err)
	}
	repo.tablaTierra = tablaTierra

	return repo, nil
}

// ObtenerTablaTierra returns the ground conductor table (250-122).
func (r *CSVTablaNOMRepository) ObtenerTablaTierra(ctx context.Context) ([]service.EntradaTablaTierra, error) {
	return r.tablaTierra, nil
}

func (r *CSVTablaNOMRepository) loadTablaTierra() ([]service.EntradaTablaTierra, error) {
	filePath := filepath.Join(r.basePath, "250-122.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open 250-122.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read 250-122.csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("250-122.csv is empty or missing header")
	}

	// Validate header
	header := records[0]
	expectedHeader := []string{"itm_hasta", "calibre", "seccion_mm2", "material"}
	for i, col := range expectedHeader {
		if i >= len(header) || header[i] != col {
			return nil, fmt.Errorf("250-122.csv: invalid header, expected %v at position %d", expectedHeader, i)
		}
	}

	var result []service.EntradaTablaTierra
	for i, record := range records[1:] {
		itm, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("250-122.csv line %d: invalid ITM value: %w", i+2, err)
		}

		seccion, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("250-122.csv line %d: invalid seccion_mm2: %w", i+2, err)
		}

		result = append(result, service.EntradaTablaTierra{
			ITMHasta: itm,
			Conductor: valueobject.ConductorParams{
				Calibre:    record[1],
				SeccionMM2: seccion,
				Material:   record[3],
			},
		})
	}

	return result, nil
}
