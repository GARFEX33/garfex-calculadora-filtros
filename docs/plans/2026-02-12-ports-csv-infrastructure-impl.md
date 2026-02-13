# Ports y CSV Infrastructure Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implementar los ports (interfaces) en application/port/ y el CSV reader en infrastructure/repository/ con caché en memoria.

**Architecture:** Hexagonal/Clean Architecture — ports definen contratos en application, CSV repository implementa en infrastructure. Las tablas NOM se cargan una vez al iniciar y se mantienen en memoria (~50-100KB).

**Tech Stack:** Go 1.22+, encoding/csv, testify, table-driven tests

**Design Doc:** `docs/plans/2026-02-12-ports-csv-infrastructure-design.md`

---

## Task 1: Crear enum MaterialConductor en valueobject

**Files:**
- Create: `internal/domain/valueobject/material_conductor.go`
- Test: `internal/domain/valueobject/material_conductor_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/valueobject/material_conductor_test.go
package valueobject

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestMaterialConductor_String(t *testing.T) {
    tests := []struct {
        name     string
        material MaterialConductor
        want     string
    }{
        {"cobre", MaterialCobre, "CU"},
        {"aluminio", MaterialAluminio, "AL"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.want, tt.material.String())
        })
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/domain/valueobject/... -v -run TestMaterialConductor`
Expected: FAIL — `MaterialConductor` y `MaterialCobre` no definidos

**Step 3: Write minimal implementation**

```go
// internal/domain/valueobject/material_conductor.go
package valueobject

// MaterialConductor represents the conductor material (Cu or Al).
type MaterialConductor int

const (
    MaterialCobre MaterialConductor = iota
    MaterialAluminio
)

// String returns the NOM standard abbreviation for the material.
func (m MaterialConductor) String() string {
    switch m {
    case MaterialCobre:
        return "CU"
    case MaterialAluminio:
        return "AL"
    default:
        return "UNKNOWN"
    }
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/domain/valueobject/... -v -run TestMaterialConductor`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/domain/valueobject/material_conductor.go internal/domain/valueobject/material_conductor_test.go
git commit -m "feat: add MaterialConductor enum"
```

---

## Task 2: Crear enum Temperatura en valueobject

**Files:**
- Create: `internal/domain/valueobject/temperatura.go`
- Test: `internal/domain/valueobject/temperatura_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/valueobject/temperatura_test.go
package valueobject

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTemperatura_Valor(t *testing.T) {
    tests := []struct {
        name string
        temp Temperatura
        want int
    }{
        {"60C", Temp60, 60},
        {"75C", Temp75, 75},
        {"90C", Temp90, 90},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.want, tt.temp.Valor())
        })
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/domain/valueobject/... -v -run TestTemperatura`
Expected: FAIL — `Temperatura` y `Temp60` no definidos

**Step 3: Write minimal implementation**

```go
// internal/domain/valueobject/temperatura.go
package valueobject

// Temperatura represents the temperature rating in Celsius (60, 75, or 90).
type Temperatura int

const (
    Temp60 Temperatura = 60
    Temp75 Temperatura = 75
    Temp90 Temperatura = 90
)

// Valor returns the temperature value in Celsius.
func (t Temperatura) Valor() int {
    return int(t)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/domain/valueobject/... -v -run TestTemperatura`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/domain/valueobject/temperatura.go internal/domain/valueobject/temperatura_test.go
git commit -m "feat: add Temperatura enum"
```

---

## Task 3: Crear ResistenciaReactancia en application/port/

**Files:**
- Create: `internal/application/port/resistencia_reactancia.go`

**Step 1: Write the implementation (simple struct, no tests needed)**

```go
// internal/application/port/resistencia_reactancia.go
package port

// ResistenciaReactancia holds the impedance values for voltage drop calculation.
type ResistenciaReactancia struct {
    R float64 // Ohms per km
    X float64 // Ohms per km
}
```

**Step 2: Commit**

```bash
git add internal/application/port/resistencia_reactancia.go
git commit -m "feat: add ResistenciaReactancia struct in port"
```

---

## Task 4: Crear TablaNOMRepository interface

**Files:**
- Create: `internal/application/port/tabla_nom_repository.go`

**Step 1: Write the interface**

```go
// internal/application/port/tabla_nom_repository.go
package port

import (
    "context"
    "github.com/garfex/calculadora-filtros/internal/domain/entity"
    "github.com/garfex/calculadora-filtros/internal/domain/service"
    "github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// TablaNOMRepository defines the contract for reading NOM tables.
type TablaNOMRepository interface {
    // ObtenerTablaAmpacidad returns ampacity table entries for the given conduit type, material, and temperature.
    ObtenerTablaAmpacidad(
        ctx context.Context,
        canalizacion entity.TipoCanalizacion,
        material valueobject.MaterialConductor,
        temperatura valueobject.Temperatura,
    ) ([]service.EntradaTablaConductor, error)

    // ObtenerTablaTierra returns the ground conductor table (250-122).
    ObtenerTablaTierra(ctx context.Context) ([]service.EntradaTablaTierra, error)

    // ObtenerImpedancia returns R and X values for the given calibre and conduit type.
    ObtenerImpedancia(
        ctx context.Context,
        calibre string,
        canalizacion entity.TipoCanalizacion,
        material valueobject.MaterialConductor,
    ) (ResistenciaReactancia, error)

    // ObtenerTablaCanalizacion returns conduit sizing table entries.
    ObtenerTablaCanalizacion(
        ctx context.Context,
        canalizacion entity.TipoCanalizacion,
    ) ([]service.EntradaTablaCanalizacion, error)
}
```

**Step 2: Verify it compiles**

Run: `go build ./internal/application/port/...`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add internal/application/port/tabla_nom_repository.go
git commit -m "feat: add TablaNOMRepository interface"
```

---

## Task 5: Crear EquipoRepository interface

**Files:**
- Create: `internal/application/port/equipo_repository.go`

**Step 1: Write the interface**

```go
// internal/application/port/equipo_repository.go
package port

import (
    "context"
    "github.com/garfex/calculadora-filtros/internal/domain/entity"
)

// EquipoRepository defines the contract for equipment persistence.
type EquipoRepository interface {
    // BuscarPorClave finds an equipment by its unique key.
    BuscarPorClave(ctx context.Context, clave string) (entity.Equipo, error)
}
```

**Step 2: Verify it compiles**

Run: `go build ./internal/application/port/...`
Expected: SUCCESS

**Step 3: Commit**

```bash
git add internal/application/port/equipo_repository.go
git commit -m "feat: add EquipoRepository interface"
```

---

## Task 6: Crear tabla 250-122.csv

**Files:**
- Create: `data/tablas_nom/250-122.csv`

**Step 1: Create CSV with NOM 250-122 ground conductor table**

```csv
itm_hasta,calibre,seccion_mm2,material
15,14 AWG,2.08,CU
20,12 AWG,3.31,CU
30,10 AWG,5.26,CU
40,10 AWG,5.26,CU
60,10 AWG,5.26,CU
100,8 AWG,8.37,CU
200,6 AWG,13.3,CU
300,4 AWG,21.2,CU
400,3 AWG,26.7,CU
500,2 AWG,33.6,CU
600,1 AWG,42.4,CU
800,1/0 AWG,53.5,CU
1000,2/0 AWG,67.4,CU
1200,3/0 AWG,85.0,CU
1600,4/0 AWG,107.2,CU
2000,250 MCM,126.7,CU
2500,300 MCM,152.0,CU
3000,350 MCM,177.3,CU
4000,400 MCM,202.7,CU
5000,500 MCM,253.4,CU
6000,600 MCM,304.0,CU
```

**Step 2: Commit**

```bash
git add data/tablas_nom/250-122.csv
git commit -m "feat: add NOM 250-122 ground conductor table"
```

---

## Task 7: Crear CSVTablaNOMRepository (constructor y carga)

**Files:**
- Create: `internal/infrastructure/repository/csv_tabla_nom_repository.go`
- Create: `internal/infrastructure/repository/testdata/250-122.csv` (copia para tests)

**Step 1: Write the failing test**

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/infrastructure/repository/... -v -run TestNewCSVTablaNOMRepository`
Expected: FAIL — `NewCSVTablaNOMRepository` no definido

**Step 3: Copy test data**

```bash
mkdir -p internal/infrastructure/repository/testdata
cp data/tablas_nom/250-122.csv internal/infrastructure/repository/testdata/
cp data/tablas_nom/310-15-b-16.csv internal/infrastructure/repository/testdata/
cp data/tablas_nom/tabla-9-resistencia-reactancia.csv internal/infrastructure/repository/testdata/
```

**Step 4: Write minimal implementation (constructor only)**

```go
// internal/infrastructure/repository/csv_tabla_nom_repository.go
package repository

import (
    "fmt"
    "os"
    "path/filepath"
)

// CSVTablaNOMRepository reads NOM tables from CSV files with in-memory caching.
type CSVTablaNOMRepository struct {
    basePath string
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

    // TODO: Load tables into memory

    return &CSVTablaNOMRepository{
        basePath: basePath,
    }, nil
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./internal/infrastructure/repository/... -v -run TestNewCSVTablaNOMRepository`
Expected: PASS

**Step 6: Commit**

```bash
git add internal/infrastructure/repository/csv_tabla_nom_repository.go
git add internal/infrastructure/repository/csv_tabla_nom_repository_test.go
git add internal/infrastructure/repository/testdata/
git commit -m "feat: add CSVTablaNOMRepository constructor"
```

---

## Task 8: Implementar ObtenerTablaTierra

**Files:**
- Modify: `internal/infrastructure/repository/csv_tabla_nom_repository.go`

**Step 1: Write the failing test**

Add to `csv_tabla_nom_repository_test.go`:

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/infrastructure/repository/... -v -run TestCSVTablaNOMRepository_ObtenerTablaTierra`
Expected: FAIL — `ObtenerTablaTierra` no implementado

**Step 3: Implement the method**

Add to `CSVTablaNOMRepository` struct:

```go
// Add field to store loaded table
tablaTierra []service.EntradaTablaTierra
```

Add to constructor:

```go
// Load ground conductor table
tablaTierra, err := r.loadTablaTierra()
if err != nil {
    return nil, fmt.Errorf("failed to load ground table: %w", err)
}
```

Implement the method and loader:

```go
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
                Calibre:     record[1],
                SeccionMM2:  seccion,
                Material:    record[3],
            },
        })
    }

    return result, nil
}
```

**Step 4: Add imports**

```go
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
```

**Step 5: Run test to verify it passes**

Run: `go test ./internal/infrastructure/repository/... -v -run TestCSVTablaNOMRepository_ObtenerTablaTierra`
Expected: PASS

**Step 6: Commit**

```bash
git add internal/infrastructure/repository/csv_tabla_nom_repository.go
git commit -m "feat: implement ObtenerTablaTierra with CSV loading"
```

---

## Task 9: Implementar ObtenerTablaAmpacidad

**Files:**
- Modify: `internal/infrastructure/repository/csv_tabla_nom_repository.go`

**Step 1: Write the failing test**

Add to test file:

```go
func TestCSVTablaNOMRepository_ObtenerTablaAmpacidad(t *testing.T) {
    repo, err := NewCSVTablaNOMRepository("testdata")
    require.NoError(t, err)

    ctx := context.Background()
    
    tests := []struct {
        name        string
        canalizacion entity.TipoCanalizacion
        material    valueobject.MaterialConductor
        temp        valueobject.Temperatura
    }{
        {"PVC Copper 75C", entity.TUBERIA_PVC, valueobject.MaterialCobre, valueobject.Temp75},
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/infrastructure/repository/... -v -run TestCSVTablaNOMRepository_ObtenerTablaAmpacidad`
Expected: FAIL — método no implementado

**Step 3: Implement the method**

Add to struct:

```go
// Cache for ampacity tables: map[tipoCanalizacion]map[material]map[temperatura][]EntradaTablaConductor
tablasAmpacidad map[entity.TipoCanalizacion]map[valueobject.MaterialConductor]map[valueobject.Temperatura][]service.EntradaTablaConductor
```

Implement the method:

```go
func (r *CSVTablaNOMRepository) ObtenerTablaAmpacidad(
    ctx context.Context,
    canalizacion entity.TipoCanalizacion,
    material valueobject.MaterialConductor,
    temperatura valueobject.Temperatura,
) ([]service.EntradaTablaConductor, error) {
    byMaterial, ok := r.tablasAmpacidad[canalizacion]
    if !ok {
        return nil, fmt.Errorf("no ampacity table for conduit type: %s", canalizacion)
    }

    byTemp, ok := byMaterial[material]
    if !ok {
        return nil, fmt.Errorf("no ampacity table for material: %s", material)
    }

    tabla, ok := byTemp[temperatura]
    if !ok {
        return nil, fmt.Errorf("no ampacity table for temperature: %d°C", temperatura)
    }

    return tabla, nil
}
```

**Step 4: Load ampacity tables in constructor**

Add to constructor:

```go
repo.tablasAmpacidad = make(map[entity.TipoCanalizacion]map[valueobject.MaterialConductor]map[valueobject.Temperatura][]service.EntradaTablaConductor)

// Load 310-15-b-16 for all conduit types using it
for _, canalizacion := range []entity.TipoCanalizacion{
    entity.TUBERIA_PVC,
    entity.TUBERIA_ALUMINIO,
    entity.TUBERIA_ACERO_PG,
    entity.TUBERIA_ACERO_PD,
} {
    repo.tablasAmpacidad[canalizacion] = make(map[valueobject.MaterialConductor]map[valueobject.Temperatura][]service.EntradaTablaConductor)
    
    for _, material := range []valueobject.MaterialConductor{
        valueobject.MaterialCobre,
        valueobject.MaterialAluminio,
    } {
        repo.tablasAmpacidad[canalizacion][material] = make(map[valueobject.Temperatura][]service.EntradaTablaConductor)
        
        tabla, err := r.loadTablaAmpacidad("310-15-b-16.csv", material)
        if err != nil {
            return nil, fmt.Errorf("failed to load ampacity table for %s %s: %w", canalizacion, material, err)
        }
        
        // Split by temperature columns
        for _, temp := range []valueobject.Temperatura{valueobject.Temp60, valueobject.Temp75, valueobject.Temp90} {
            repo.tablasAmpacidad[canalizacion][material][temp] = r.extractByTemperature(tabla, material, temp)
        }
    }
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./internal/infrastructure/repository/... -v -run TestCSVTablaNOMRepository_ObtenerTablaAmpacidad`
Expected: PASS

**Step 6: Commit**

```bash
git add internal/infrastructure/repository/csv_tabla_nom_repository.go
git commit -m "feat: implement ObtenerTablaAmpacidad with CSV loading"
```

---

## Task 10: Run all tests and verify

**Step 1: Run all tests**

Run: `go test ./internal/... -v`
Expected: All tests PASS

**Step 2: Run with race detector**

Run: `go test -race ./internal/...`
Expected: No race conditions detected

**Step 3: Verify build**

Run: `go build ./...`
Expected: SUCCESS

---

## Summary

After completing these tasks:

- ✅ `MaterialConductor` and `Temperatura` enums in `valueobject/`
- ✅ `TablaNOMRepository` and `EquipoRepository` interfaces in `application/port/`
- ✅ `CSVTablaNOMRepository` implementation with in-memory caching
- ✅ `250-122.csv` ground conductor table
- ✅ All tests passing

Next iteration: PostgreSQL repository and use case orchestration.
