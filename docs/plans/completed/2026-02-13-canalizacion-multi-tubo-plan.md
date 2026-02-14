# Canalizacion Multi-Tubo Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Agregar soporte para N tubos en paralelo en `CalcularCanalizacion`, pasando `numeroDeTubos int` como parámetro explícito.

**Architecture:** Modificar `entity.Canalizacion` para incluir `NumeroDeTubos`, actualizar la firma y lógica de `CalcularCanalizacion`, actualizar el único caller en `application/usecase/calcular_memoria.go`, y actualizar todos los tests existentes más agregar tests nuevos.

**Tech Stack:** Go 1.22+, testify

---

### Task 1: Agregar `NumeroDeTubos` a `entity.Canalizacion`

**Files:**
- Modify: `internal/domain/entity/canalizacion.go`

**Step 1: Leer el archivo actual**

```
internal/domain/entity/canalizacion.go (ya leído — 10 líneas)
```

**Step 2: Agregar el campo**

Reemplazar el contenido de `canalizacion.go` con:

```go
// internal/domain/entity/canalizacion.go
package entity

// Canalizacion represents the conduit or cable tray selected for the installation.
type Canalizacion struct {
	Tipo           string  // "TUBERIA" | "CHAROLA"
	Tamano         string  // e.g., "1 1/2" (inches for tubería)
	AnchoRequerido float64 // for charola: required width in mm; for tubería: total conductor area in mm²
	NumeroDeTubos  int     // number of parallel conduits; 1 = single conduit installation
}
```

**Step 3: Verificar que compila**

```bash
go build ./internal/domain/entity/...
```
Esperado: sin errores.

**Step 4: Commit**

```bash
git add internal/domain/entity/canalizacion.go
git commit -m "feat(entity): add NumeroDeTubos field to Canalizacion"
```

---

### Task 2: Actualizar `CalcularCanalizacion` — firma y lógica

**Files:**
- Modify: `internal/domain/service/calculo_canalizacion.go`

**Step 1: Reemplazar el archivo completo**

```go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ErrCanalizacionNoDisponible is returned when no conduit size fits the required area.
var ErrCanalizacionNoDisponible = errors.New("no se encontró canalización con área suficiente")

func determinarFillFactor(cantidad int) float64 {
	switch cantidad {
	case 1:
		return 0.53
	case 2:
		return 0.31
	default:
		return 0.40
	}
}

// ConductorParaCanalizacion holds the quantity and cross-section area
// of a group of identical conductors for conduit sizing calculations.
type ConductorParaCanalizacion struct {
	Cantidad   int
	SeccionMM2 float64
}

// CalcularCanalizacion selects the smallest conduit whose usable area
// (interior area × fill factor) accommodates all conductors.
// tipo should be a TipoCanalizacion string value (e.g., "TUBERIA_CONDUIT").
// numeroDeTubos indicates how many parallel conduits to use (must be >= 1).
// When numeroDeTubos > 1, the total conductor area and count are divided evenly
// among the tubes; fill factor is determined per-tube conductor count.
func CalcularCanalizacion(
	conductores []ConductorParaCanalizacion,
	tipo string,
	tabla []valueobject.EntradaTablaCanalizacion,
	numeroDeTubos int,
) (entity.Canalizacion, error) {
	if numeroDeTubos < 1 {
		return entity.Canalizacion{}, fmt.Errorf("numeroDeTubos debe ser mayor a cero")
	}
	if len(conductores) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("lista de conductores vacía")
	}
	if len(tabla) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("%w: tabla vacía", ErrCanalizacionNoDisponible)
	}

	var areaTotal float64
	var cantidadTotal int
	for _, c := range conductores {
		areaTotal += float64(c.Cantidad) * c.SeccionMM2
		cantidadTotal += c.Cantidad
	}

	conductoresPorTubo := cantidadTotal / numeroDeTubos
	factorRelleno := determinarFillFactor(conductoresPorTubo)
	areaPorTubo := areaTotal / float64(numeroDeTubos)
	areaRequerida := areaPorTubo / factorRelleno

	for _, entrada := range tabla {
		if entrada.AreaInteriorMM2 >= areaRequerida {
			return entity.Canalizacion{
				Tipo:           tipo,
				Tamano:         entrada.Tamano,
				AnchoRequerido: areaTotal,
				NumeroDeTubos:  numeroDeTubos,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf(
		"%w: área requerida %.2f mm² excede máxima disponible %.2f mm²",
		ErrCanalizacionNoDisponible, areaRequerida, tabla[len(tabla)-1].AreaInteriorMM2,
	)
}
```

**Step 2: Verificar que compila**

```bash
go build ./internal/domain/service/...
```
Esperado: error de compilación en el test (firma cambió) — eso es correcto, lo arreglamos en Task 3.

**Step 3: Commit**

```bash
git add internal/domain/service/calculo_canalizacion.go
git commit -m "feat(service): add numeroDeTubos param to CalcularCanalizacion"
```

---

### Task 3: Actualizar tests existentes en `calculo_canalizacion_test.go`

**Files:**
- Modify: `internal/domain/service/calculo_canalizacion_test.go`

**Step 1: Agregar `1` como último argumento en todas las llamadas existentes**

Hay 6 llamadas a `service.CalcularCanalizacion` en los tests actuales. Todas pasan de:

```go
service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest)
```

a:

```go
service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 1)
```

Los tests afectados son:
- `TestCalcularCanalizacion_Tuberia` (línea 35)
- `TestCalcularCanalizacion_SmallConductors` (línea 51)
- `TestCalcularCanalizacion_NoFit` (línea 62)
- `TestCalcularCanalizacion_EmptyConductors` (línea 68)
- `TestCalcularCanalizacion_FillFactor1Conductor` (línea 81)
- `TestCalcularCanalizacion_FillFactor2Conductores` (línea 95)
- `TestCalcularCanalizacion_FillFactor3Conductores` (línea 109)

**Step 2: Correr los tests existentes para verificar que pasan**

```bash
go test ./internal/domain/service/... -run TestCalcularCanalizacion -v
```
Esperado: todos los tests existentes PASS.

**Step 3: Commit**

```bash
git add internal/domain/service/calculo_canalizacion_test.go
git commit -m "test(service): update existing CalcularCanalizacion tests to new signature"
```

---

### Task 4: Agregar tests nuevos para multi-tubo

**Files:**
- Modify: `internal/domain/service/calculo_canalizacion_test.go`

**Step 1: Escribir los tests nuevos (agregar al final del archivo)**

```go
func TestCalcularCanalizacion_DosTubos(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 33.62}, // 3 phases × 2 AWG
		{Cantidad: 1, SeccionMM2: 13.30}, // 1 ground × 6 AWG
	}
	// Total area = 114.16 mm², cantidadTotal = 4
	// Con 2 tubos: conductoresPorTubo = 4/2 = 2 → fillFactor = 0.31
	// areaPorTubo = 114.16/2 = 57.08 mm²
	// areaRequerida = 57.08/0.31 = 184.13 mm²
	// Smallest ≥ 184.13 → "1" (198 mm²)
	// Con 1 tubo habría dado "1 1/2" (360 mm²) — el tubo individual es más chico

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 2)
	require.NoError(t, err)
	assert.Equal(t, "TUBERIA_CONDUIT", result.Tipo)
	assert.Equal(t, "1", result.Tamano)
	assert.Equal(t, 2, result.NumeroDeTubos)
	assert.InDelta(t, 114.16, result.AnchoRequerido, 0.01)
}

func TestCalcularCanalizacion_DosTubosSmall(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 4, SeccionMM2: 3.31}, // 4 × 12 AWG
	}
	// Total area = 13.24 mm², cantidadTotal = 4
	// Con 2 tubos: conductoresPorTubo = 4/2 = 2 → fillFactor = 0.31
	// areaPorTubo = 13.24/2 = 6.62 mm²
	// areaRequerida = 6.62/0.31 = 21.35 mm²
	// Smallest ≥ 21.35 → "1/2" (78 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 2)
	require.NoError(t, err)
	assert.Equal(t, "1/2", result.Tamano)
	assert.Equal(t, 2, result.NumeroDeTubos)
}

func TestCalcularCanalizacion_NumeroDeTubosCero(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 3.31},
	}
	_, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "numeroDeTubos debe ser mayor a cero")
}

func TestCalcularCanalizacion_NumeroDeTubosNegativo(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 3.31},
	}
	_, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "numeroDeTubos debe ser mayor a cero")
}

func TestCalcularCanalizacion_NumeroDeTubosMayorQueConductores(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 2, SeccionMM2: 3.31},
	}
	// 2 conductores, 5 tubos — permitido, el servicio no juzga
	// cantidadTotal=2, conductoresPorTubo=2/5=0 → fillFactor=0.40 (default)
	// areaPorTubo = 6.62/5 = 1.324 mm²
	// areaRequerida = 1.324/0.40 = 3.31 mm²
	// Smallest ≥ 3.31 → "1/2" (78 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA_CONDUIT", tablaCanalizacionTest, 5)
	require.NoError(t, err)
	assert.Equal(t, "1/2", result.Tamano)
	assert.Equal(t, 5, result.NumeroDeTubos)
}
```

**Step 2: Correr los tests nuevos**

```bash
go test ./internal/domain/service/... -run TestCalcularCanalizacion -v
```
Esperado: todos PASS incluyendo los 5 tests nuevos.

**Step 3: Correr con race detector**

```bash
go test -race ./internal/domain/service/...
```
Esperado: PASS sin data races.

**Step 4: Commit**

```bash
git add internal/domain/service/calculo_canalizacion_test.go
git commit -m "test(service): add multi-tubo tests for CalcularCanalizacion"
```

---

### Task 5: Actualizar caller en `application/usecase/calcular_memoria.go`

**Files:**
- Modify: `internal/application/usecase/calcular_memoria.go:193`

**Step 1: Leer el contexto alrededor de la línea 193**

El caller actual (línea 193):
```go
canalizacion, err := service.CalcularCanalizacion(conductores, string(input.TipoCanalizacion), tablaCanalizacion)
```

`hilosPorFase` ya está disponible en el contexto (se usa en línea 189). Actualizar a:

```go
canalizacion, err := service.CalcularCanalizacion(conductores, string(input.TipoCanalizacion), tablaCanalizacion, hilosPorFase)
```

**Step 2: Verificar que el DTO de output expone NumeroDeTubos**

Revisar `internal/application/dto/` para ver si `ResultadoCanalizacion` tiene campo para número de tubos. Si no existe, agregar:

```go
// En el struct ResultadoCanalizacion del dto correspondiente:
NumeroDeTubos int `json:"numero_de_tubos"`
```

Y en el mapeo (línea ~198):
```go
output.Canalizacion = dto.ResultadoCanalizacion{
    Tipo:             input.TipoCanalizacion,
    Tamano:           canalizacion.Tamano,
    NumeroDeTubos:    canalizacion.NumeroDeTubos,
    AreaTotalMM2:     canalizacion.AnchoRequerido,
    AreaRequeridaMM2: canalizacion.AnchoRequerido / 0.40,
}
```

**Step 3: Verificar que compila todo**

```bash
go build ./...
```
Esperado: sin errores.

**Step 4: Correr todos los tests**

```bash
go test ./...
```
Esperado: PASS.

**Step 5: Commit**

```bash
git add internal/application/usecase/calcular_memoria.go
git add internal/application/dto/  # si se modificó
git commit -m "feat(application): pass hilosPorFase as numeroDeTubos to CalcularCanalizacion"
```

---

### Task 6: Verificación final

**Step 1: Tests completos con race detector**

```bash
go test -race ./...
```
Esperado: PASS sin data races.

**Step 2: Build limpio**

```bash
go build ./...
```

**Step 3: Vet**

```bash
go vet ./...
```
Esperado: sin warnings.
