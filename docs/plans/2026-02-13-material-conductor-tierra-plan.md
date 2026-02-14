# Material Cu/Al Conductor Tierra Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Limpiar calibres a la lista oficial, actualizar `250-122.csv` con columnas Cu+Al según NOM, y permitir que el usuario elija material (Cu o Al) en toda la memoria de cálculo.

**Architecture:** Los cambios van de abajo hacia arriba: CSV → value object → service → application → DTO. Cada tarea es independiente excepto que las de service dependen del value object, y las de application dependen del service.

**Tech Stack:** Go 1.22+, encoding/csv, testify

---

### Task 1: Limpiar tablas CSV de ampacidad (b-16, b-17, b-20)

**Files:**
- Modify: `data/tablas_nom/310-15-b-16.csv`
- Modify: `data/tablas_nom/310-15-b-17.csv`
- Modify: `data/tablas_nom/310-15-b-20.csv`

**Eliminar estas filas en los 3 archivos** (calibres fuera de la lista oficial):
- `18 AWG`
- `16 AWG`
- `3 AWG`
- `1 AWG`
- `700 MCM`
- `800 MCM`
- `900 MCM`
- `1250 MCM`
- `1500 MCM`
- `1750 MCM`
- `2000 MCM`

**Nota:** `310-15-b-20.csv` solo tiene desde `8 AWG` — solo eliminar `3 AWG`, `1 AWG`, `700 MCM`, `800 MCM`, `900 MCM`.

**Resultado esperado de `310-15-b-16.csv`** (solo estas filas, header incluido):
```
seccion_mm2,calibre,cu_60c,cu_75c,cu_90c,al_60c,al_75c,al_90c
2.08,14 AWG,15,20,25,,,
3.31,12 AWG,20,25,30,,,
5.26,10 AWG,30,35,40,,,
8.37,8 AWG,40,50,55,,,
13.3,6 AWG,55,65,75,40,50,55
21.2,4 AWG,70,85,95,55,65,75
33.6,2 AWG,95,115,130,75,90,100
53.5,1/0 AWG,125,150,170,100,120,135
67.4,2/0 AWG,145,175,195,115,135,150
85.0,3/0 AWG,165,200,225,130,155,175
107.2,4/0 AWG,195,230,260,150,180,205
127,250 MCM,215,255,290,170,205,230
152,300 MCM,240,285,320,195,230,260
177,350 MCM,260,310,350,210,250,280
203,400 MCM,280,335,380,225,270,305
253,500 MCM,320,380,430,260,310,350
304,600 MCM,350,420,475,285,340,385
380,750 MCM,400,475,535,320,385,435
507,1000 MCM,455,545,615,375,445,500
```

**Verificar que el servidor sigue compilando:**
```bash
go build ./...
```
Esperado: sin errores.

**Step: Commit**
```bash
git add data/tablas_nom/310-15-b-16.csv data/tablas_nom/310-15-b-17.csv data/tablas_nom/310-15-b-20.csv
git commit -m "data: remove out-of-spec calibres from ampacity tables"
```

---

### Task 2: Reemplazar `250-122.csv` con estructura Cu+Al

**Files:**
- Modify: `data/tablas_nom/250-122.csv`

**Reemplazar el contenido completo con:**
```
itm_hasta,cu_calibre,cu_seccion_mm2,al_calibre,al_seccion_mm2
1,14 AWG,2.08,,
15,14 AWG,2.08,,
20,12 AWG,3.31,,
60,10 AWG,5.26,,
100,8 AWG,8.37,,
200,6 AWG,13.3,4 AWG,21.2
300,4 AWG,21.2,2 AWG,33.6
400,2 AWG,33.6,1/0 AWG,42.4
500,2 AWG,33.6,1/0 AWG,53.5
800,1/0 AWG,53.5,3/0 AWG,85.0
1000,2/0 AWG,67.4,4/0 AWG,107.2
1200,3/0 AWG,85.0,250 MCM,127.0
1600,4/0 AWG,107.2,350 MCM,177.0
2000,250 MCM,127.0,400 MCM,203.0
2500,350 MCM,177.0,600 MCM,304.0
3000,400 MCM,203.0,600 MCM,304.0
4000,500 MCM,253.0,750 MCM,380.0
```

**Step: Commit**
```bash
git add data/tablas_nom/250-122.csv
git commit -m "data: update 250-122 with Cu+Al columns per NOM, cut at ITM 4000"
```

---

### Task 3: Limpiar `calibresValidos` en `conductor.go`

**Files:**
- Modify: `internal/domain/valueobject/conductor.go`

**Reemplazar el mapa `calibresValidos` (líneas 20-30) con:**
```go
var calibresValidos = map[string]bool{
	// AWG
	"14 AWG": true, "12 AWG": true, "10 AWG": true, "8 AWG": true,
	"6 AWG":  true, "4 AWG":  true, "2 AWG":  true,
	"1/0 AWG": true, "2/0 AWG": true, "3/0 AWG": true, "4/0 AWG": true,
	// MCM
	"250 MCM": true, "300 MCM": true, "350 MCM": true, "400 MCM": true,
	"500 MCM": true, "600 MCM": true, "750 MCM": true, "1000 MCM": true,
}
```

**Eliminar: 18 AWG, 16 AWG, 700 MCM, 800 MCM, 900 MCM, 1250 MCM, 1500 MCM, 1750 MCM, 2000 MCM**
**Nota: 3 AWG y 1 AWG ya no estaban — confirmar que siguen sin estar.**

**Step: Actualizar `conductor_test.go`**

El test `TestNewConductor_ExtremosCalibre` (líneas 112-122) actualmente verifica que `18 AWG` y `2000 MCM` son válidos. Hay que actualizarlo para verificar los nuevos extremos:

```go
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

	// Calibres eliminados: inválidos
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
```

**Nota:** El test `TestNewConductor_CalibreInvalido` en línea 48 ya incluye `"3 AWG"` como inválido — ese test debe seguir pasando sin cambios.

**Verificar tests:**
```bash
go test ./internal/domain/valueobject/... -v
```
Esperado: todos PASS.

**Step: Commit**
```bash
git add internal/domain/valueobject/conductor.go internal/domain/valueobject/conductor_test.go
git commit -m "feat(valueobject): restrict calibresValidos to official NOM list"
```

---

### Task 4: Actualizar `EntradaTablaTierra` para Cu+Al

**Files:**
- Modify: `internal/domain/valueobject/tabla_entrada.go`

**Reemplazar `EntradaTablaTierra` (líneas 13-20) con:**
```go
// EntradaTablaTierra represents one row from NOM table 250-122.
// Entries must be sorted by ITMHasta ascending.
// ConductorCu is always present. ConductorAl is nil when aluminium is not
// permitted for this ITM range (per NOM) — callers fall back to ConductorCu.
type EntradaTablaTierra struct {
	ITMHasta    int
	ConductorCu ConductorParams  // always present
	ConductorAl *ConductorParams // nil = not available for this ITM, use Cu fallback
}
```

**Verificar que compila:**
```bash
go build ./internal/domain/valueobject/...
```
Esperado: errores de compilación en `calculo_tierra.go` y en el CSV reader — esperado, se arreglan en tasks 5 y 7.

**Step: Commit**
```bash
git add internal/domain/valueobject/tabla_entrada.go
git commit -m "feat(valueobject): add ConductorAl to EntradaTablaTierra for NOM 250-122"
```

---

### Task 5: Actualizar `SeleccionarConductorTierra` — nueva firma y lógica

**Files:**
- Modify: `internal/domain/service/calculo_tierra.go`
- Modify: `internal/domain/service/calculo_tierra_test.go`

**Step 1: Escribir los tests nuevos PRIMERO (TDD)**

Reemplazar el archivo `calculo_tierra_test.go` completo con:

```go
// internal/domain/service/calculo_tierra_test.go
package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helpers para construir EntradaTablaTierra

func entradaTierraCu(itmHasta int, calibre string, seccionMM2 float64) valueobject.EntradaTablaTierra {
	return valueobject.EntradaTablaTierra{
		ITMHasta: itmHasta,
		ConductorCu: valueobject.ConductorParams{
			Calibre:    calibre,
			Material:   "Cu",
			SeccionMM2: seccionMM2,
		},
		ConductorAl: nil,
	}
}

func entradaTierraCuAl(itmHasta int, cuCalibre string, cuSeccion float64, alCalibre string, alSeccion float64) valueobject.EntradaTablaTierra {
	al := valueobject.ConductorParams{
		Calibre:    alCalibre,
		Material:   "Al",
		SeccionMM2: alSeccion,
	}
	return valueobject.EntradaTablaTierra{
		ITMHasta: itmHasta,
		ConductorCu: valueobject.ConductorParams{
			Calibre:    cuCalibre,
			Material:   "Cu",
			SeccionMM2: cuSeccion,
		},
		ConductorAl: &al,
	}
}

// Tabla de prueba basada en NOM 250-122
var tablaTierraTest = []valueobject.EntradaTablaTierra{
	entradaTierraCu(15, "14 AWG", 2.08),
	entradaTierraCu(20, "12 AWG", 3.31),
	entradaTierraCu(60, "10 AWG", 5.26),
	entradaTierraCu(100, "8 AWG", 8.37),
	entradaTierraCuAl(200, "6 AWG", 13.3, "4 AWG", 21.2),
	entradaTierraCuAl(400, "2 AWG", 33.6, "1/0 AWG", 42.4),
	entradaTierraCuAl(800, "1/0 AWG", 53.5, "3/0 AWG", 85.0),
	entradaTierraCuAl(1000, "2/0 AWG", 67.4, "4/0 AWG", 107.2),
	entradaTierraCuAl(4000, "500 MCM", 253.0, "750 MCM", 380.0),
}

func TestSeleccionarConductorTierra_CuExplicito(t *testing.T) {
	tests := []struct {
		name            string
		itm             int
		expectedCalibre string
	}{
		{"ITM 15 → 14 AWG Cu", 15, "14 AWG"},
		{"ITM 20 → 12 AWG Cu", 20, "12 AWG"},
		{"ITM 30 → 10 AWG Cu (≤60)", 30, "10 AWG"},
		{"ITM 100 → 8 AWG Cu", 100, "8 AWG"},
		{"ITM 125 → 6 AWG Cu (≤200)", 125, "6 AWG"},
		{"ITM 400 → 2 AWG Cu", 400, "2 AWG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conductor, err := service.SeleccionarConductorTierra(tt.itm, valueobject.MaterialCobre, tablaTierraTest)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCalibre, conductor.Calibre())
			assert.Equal(t, "Cu", conductor.Material())
		})
	}
}

func TestSeleccionarConductorTierra_AluminioDisponible(t *testing.T) {
	// ITM 200, Al disponible → 4 AWG Al
	conductor, err := service.SeleccionarConductorTierra(200, valueobject.MaterialAluminio, tablaTierraTest)
	require.NoError(t, err)
	assert.Equal(t, "4 AWG", conductor.Calibre())
	assert.Equal(t, "Al", conductor.Material())
}

func TestSeleccionarConductorTierra_AluminioFallbackCu(t *testing.T) {
	// ITM 60, Al NO disponible → fallback silencioso a Cu (10 AWG)
	conductor, err := service.SeleccionarConductorTierra(60, valueobject.MaterialAluminio, tablaTierraTest)
	require.NoError(t, err)
	assert.Equal(t, "10 AWG", conductor.Calibre())
	assert.Equal(t, "Cu", conductor.Material())
}

func TestSeleccionarConductorTierra_AluminioITMMaximo(t *testing.T) {
	// ITM 4000, Al disponible → 750 MCM Al
	conductor, err := service.SeleccionarConductorTierra(4000, valueobject.MaterialAluminio, tablaTierraTest)
	require.NoError(t, err)
	assert.Equal(t, "750 MCM", conductor.Calibre())
	assert.Equal(t, "Al", conductor.Material())
}

func TestSeleccionarConductorTierra_ITMExceedsTable(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(5000, valueobject.MaterialCobre, tablaTierraTest)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrConductorNoEncontrado))
}

func TestSeleccionarConductorTierra_InvalidITM(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(0, valueobject.MaterialCobre, tablaTierraTest)
	assert.Error(t, err)
}

func TestSeleccionarConductorTierra_EmptyTable(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(100, valueobject.MaterialCobre, nil)
	assert.Error(t, err)
}
```

**Step 2: Verificar que los tests FALLAN (compilación falla por firma incorrecta)**
```bash
go test ./internal/domain/service/... -run TestSeleccionarConductorTierra -v
```
Esperado: error de compilación — firma todavía tiene 2 args.

**Step 3: Implementar la nueva firma en `calculo_tierra.go`**

Reemplazar el archivo completo con:
```go
// internal/domain/service/calculo_tierra.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// SeleccionarConductorTierra selects the ground conductor from NOM table 250-122
// based on the equipment's ITM rating and the desired conductor material.
// If material is Al but Al is not available for the given ITM range (NOM restriction),
// it silently falls back to Cu.
func SeleccionarConductorTierra(
	itm int,
	material valueobject.MaterialConductor,
	tabla []valueobject.EntradaTablaTierra,
) (valueobject.Conductor, error) {
	if itm <= 0 {
		return valueobject.Conductor{}, fmt.Errorf("ITM debe ser mayor que cero: %d", itm)
	}
	if len(tabla) == 0 {
		return valueobject.Conductor{}, fmt.Errorf("%w: tabla de tierra vacía", ErrConductorNoEncontrado)
	}

	for _, entrada := range tabla {
		if itm <= entrada.ITMHasta {
			if material == valueobject.MaterialAluminio && entrada.ConductorAl != nil {
				return valueobject.NewConductor(*entrada.ConductorAl)
			}
			// Cu explícito, o Al sin disponibilidad → fallback a Cu
			return valueobject.NewConductor(entrada.ConductorCu)
		}
	}

	return valueobject.Conductor{}, fmt.Errorf(
		"%w: ITM %d excede máximo de tabla %d",
		ErrConductorNoEncontrado, itm, tabla[len(tabla)-1].ITMHasta,
	)
}
```

**Step 4: Verificar que los tests PASAN**
```bash
go test ./internal/domain/service/... -run TestSeleccionarConductorTierra -v
```
Esperado: todos PASS.

**Step 5: Correr todos los tests del service**
```bash
go test ./internal/domain/service/... -v
```
Esperado: todos PASS (algunos de calculo_tierra pueden fallar por el caller en application — se arregla en Task 6).

**Step 6: Commit**
```bash
git add internal/domain/service/calculo_tierra.go internal/domain/service/calculo_tierra_test.go
git commit -m "feat(service): add material param to SeleccionarConductorTierra with Al fallback"
```

---

### Task 6: Actualizar `EquipoInput` y `calcular_memoria.go`

**Files:**
- Modify: `internal/application/dto/equipo_input.go`
- Modify: `internal/application/usecase/calcular_memoria.go`
- Modify: `internal/application/dto/memoria_output.go`

**Step 1: Agregar `Material` a `EquipoInput`**

En `equipo_input.go`, agregar después de `HilosPorFase`:
```go
Material valueobject.MaterialConductor // "Cu" o "Al"; si vacío, default Cu
```

**Step 2: Agregar `Material` a `MemoriaOutput`**

En `memoria_output.go`, agregar en el struct `MemoriaOutput` (buscar el struct principal):
```go
Material string `json:"material"` // "Cu" o "Al"
```

**Step 3: Actualizar `calcular_memoria.go`**

**Cambio A** — reemplazar línea 135 (hardcoded `material := valueobject.MaterialCobre`):
```go
material := input.Material
if material == "" {
    material = valueobject.MaterialCobre
}
output.Material = string(material)
```

**Cambio B** — en la llamada a `SeleccionarConductorTierra` (~línea 169), agregar `material`:
```go
conductorTierra, err := service.SeleccionarConductorTierra(input.ITM, material, tablaTierra)
```

**Step 4: Verificar build completo**
```bash
go build ./...
```
Esperado: sin errores.

**Step 5: Commit**
```bash
git add internal/application/dto/equipo_input.go internal/application/dto/memoria_output.go internal/application/usecase/calcular_memoria.go
git commit -m "feat(application): add Material field to input/output, wire through use case"
```

---

### Task 7: Actualizar CSV reader de `250-122` en infrastructure

**Files:**
- Modify: `internal/infrastructure/repository/csv_tabla_nom_repository.go`

El método `loadTablaTierra` (~línea 349) actualmente espera el header `itm_hasta,calibre,seccion_mm2,material`. Debe cambiarse para parsear el nuevo header `itm_hasta,cu_calibre,cu_seccion_mm2,al_calibre,al_seccion_mm2`.

**Reemplazar el método `loadTablaTierra` completo** (líneas 349-408):

```go
func (r *CSVTablaNOMRepository) loadTablaTierra() ([]valueobject.EntradaTablaTierra, error) {
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
	expectedHeader := []string{"itm_hasta", "cu_calibre", "cu_seccion_mm2", "al_calibre", "al_seccion_mm2"}
	for i, col := range expectedHeader {
		if i >= len(header) || header[i] != col {
			return nil, fmt.Errorf("250-122.csv: invalid header at position %d, expected %q got %q", i, col, header[i])
		}
	}

	var result []valueobject.EntradaTablaTierra
	for i, record := range records[1:] {
		if len(record) < 3 {
			continue
		}

		itm, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("250-122.csv line %d: invalid ITM value: %w", i+2, err)
		}

		cuSeccion, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("250-122.csv line %d: invalid cu_seccion_mm2: %w", i+2, err)
		}

		entrada := valueobject.EntradaTablaTierra{
			ITMHasta: itm,
			ConductorCu: valueobject.ConductorParams{
				Calibre:    record[1],
				Material:   "Cu",
				SeccionMM2: cuSeccion,
			},
			ConductorAl: nil,
		}

		// Parse Al columns if present and non-empty
		if len(record) >= 5 && record[3] != "" && record[4] != "" {
			alSeccion, err := strconv.ParseFloat(record[4], 64)
			if err != nil {
				return nil, fmt.Errorf("250-122.csv line %d: invalid al_seccion_mm2: %w", i+2, err)
			}
			alParams := valueobject.ConductorParams{
				Calibre:    record[3],
				Material:   "Al",
				SeccionMM2: alSeccion,
			}
			entrada.ConductorAl = &alParams
		}

		result = append(result, entrada)
	}

	return result, nil
}
```

**Step: Verificar build y tests**
```bash
go build ./...
go test ./internal/infrastructure/... -v
```
Esperado: PASS.

**Step: Commit**
```bash
git add internal/infrastructure/repository/csv_tabla_nom_repository.go
git commit -m "feat(infrastructure): update loadTablaTierra to parse Cu+Al columns from 250-122"
```

---

### Task 8: Verificación final

**Step 1: Build limpio**
```bash
go build ./...
```

**Step 2: Todos los tests**
```bash
go test ./...
```
Esperado: todos PASS (el fallo preexistente `TestFase2_CalculoCompleto` en integration es un bug separado — no relacionado).

**Step 3: Vet**
```bash
go vet ./...
```
Esperado: sin warnings.

**Step 4: Prueba manual — Cu (default)**
```bash
curl -s -X POST http://localhost:8080/api/v1/calculos/memoria \
  -H "Content-Type: application/json" \
  -d '{"modo":"MANUAL_AMPERAJE","amperaje_nominal":50,"tension":220,"tipo_canalizacion":"TUBERIA_PVC","hilos_por_fase":1,"longitud_circuito":30,"itm":200,"factor_potencia":0.9,"estado":"Jalisco","sistema_electrico":"DELTA"}'
```
Esperado: `"material":"Cu"`, conductor tierra `6 AWG Cu`.

**Step 5: Prueba manual — Al explícito**
```bash
curl -s -X POST http://localhost:8080/api/v1/calculos/memoria \
  -H "Content-Type: application/json" \
  -d '{"modo":"MANUAL_AMPERAJE","amperaje_nominal":50,"tension":220,"tipo_canalizacion":"TUBERIA_PVC","hilos_por_fase":1,"longitud_circuito":30,"itm":200,"factor_potencia":0.9,"estado":"Jalisco","sistema_electrico":"DELTA","material":"Al"}'
```
Esperado: `"material":"Al"`, conductor tierra `4 AWG Al`.

**Step 6: Prueba manual — Al con fallback (ITM ≤ 100)**
```bash
curl -s -X POST http://localhost:8080/api/v1/calculos/memoria \
  -H "Content-Type: application/json" \
  -d '{"modo":"MANUAL_AMPERAJE","amperaje_nominal":30,"tension":220,"tipo_canalizacion":"TUBERIA_PVC","hilos_por_fase":1,"longitud_circuito":20,"itm":60,"factor_potencia":0.9,"estado":"Jalisco","sistema_electrico":"DELTA","material":"Al"}'
```
Esperado: `"material":"Al"`, conductor tierra `10 AWG Cu` (fallback silencioso).
