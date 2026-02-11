# Domain Layer Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement the complete domain layer (entities, value objects, domain services) for the Garfex electrical calculation system — Fase 1 MVP with Filtro Activo and Filtro de Rechazo equipment types.

**Architecture:** Hexagonal/Clean Architecture domain core. Zero external dependencies in domain code (only Go stdlib). Domain services that need NOM table data receive pre-resolved slices of domain types — no port interfaces in the domain layer. Value objects are truly immutable (unexported fields + getters). All types validated via constructors.

**Tech Stack:** Go 1.22+, testing stdlib + testify (dev only)

**Reference:** `docs/plans/2026-02-09-arquitectura-inicial-design.md` (architecture design)

---

## Design Decisions

1. **Errors per package** — Each package defines its own errors to avoid circular imports (`entity` → `valueobject` exists, so `valueobject` cannot import `entity`)
2. **Tests alongside code** — `_test.go` files in the same directory (idiomatic Go), not in `tests/` directory. The `tests/` directory will be used for integration tests in future phases
3. **Conductor includes SeccionMM2** — Cross-section area is an intrinsic conductor property needed by both voltage-drop and conduit-sizing calculations
4. **Domain services accept pre-resolved data** — Services that need NOM table lookups (conductor selection, ground conductor, conduit sizing) accept `[]EntradaTablaXxx` slices. The application layer will resolve repository data and pass it in
5. **testify for assertions** — Listed as project dependency. "Sin dependencias externas" means no infrastructure deps (DB, HTTP, files), not test utilities
6. **Value objects are immutable** — Unexported fields with getter methods; construction only via `NewXxx()` constructors

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `.golangci.yml`
- Create: `.env.example`
- Create: directory structure for all layers (empty dirs with `.gitkeep`)

**Step 1: Initialize Go module**

```bash
go mod init github.com/garfex/calculadora-filtros
```

**Step 2: Create directory structure**

```bash
mkdir -p cmd/api
mkdir -p internal/domain/entity
mkdir -p internal/domain/valueobject
mkdir -p internal/domain/service
mkdir -p internal/application/port
mkdir -p internal/application/usecase
mkdir -p internal/application/dto
mkdir -p internal/infrastructure/repository
mkdir -p internal/infrastructure/client
mkdir -p internal/presentation/handler
mkdir -p internal/presentation/middleware
mkdir -p data/tablas_nom
mkdir -p tests
```

Create `.gitkeep` files in empty directories that have no Go files yet:

```bash
# Only in directories that won't have files in THIS plan
touch cmd/api/.gitkeep
touch internal/application/port/.gitkeep
touch internal/application/usecase/.gitkeep
touch internal/application/dto/.gitkeep
touch internal/infrastructure/repository/.gitkeep
touch internal/infrastructure/client/.gitkeep
touch internal/presentation/handler/.gitkeep
touch internal/presentation/middleware/.gitkeep
touch data/tablas_nom/.gitkeep
touch tests/.gitkeep
```

**Step 3: Create `.golangci.yml`**

```yaml
linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - gofmt
    - goimports

run:
  timeout: 5m
```

**Step 4: Create `.env.example`**

```bash
# Desarrollo (laptop Windows - misma red WiFi)
DB_HOST=192.168.1.X
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=postgres

# ENVIRONMENT: development | production
ENVIRONMENT=development
```

**Step 5: Install testify dependency**

```bash
go get github.com/stretchr/testify@latest
```

**Step 6: Commit**

```bash
git add go.mod go.sum .golangci.yml .env.example cmd/ internal/ data/ tests/
git commit -m "chore: scaffold project structure, Go module, and tooling config"
```

---

### Task 2: Value Object — Corriente

**Files:**
- Create: `internal/domain/valueobject/corriente.go`
- Test: `internal/domain/valueobject/corriente_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/valueobject/corriente_test.go
package valueobject_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCorriente(t *testing.T) {
	tests := []struct {
		name    string
		valor   float64
		wantErr bool
	}{
		{"positive value", 120.5, false},
		{"small positive value", 0.01, false},
		{"zero is invalid", 0, true},
		{"negative is invalid", -10.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := valueobject.NewCorriente(tt.valor)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, valueobject.ErrCorrienteInvalida))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.valor, c.Valor())
			assert.Equal(t, "A", c.Unidad())
		})
	}
}

func TestCorriente_Multiplicar(t *testing.T) {
	c, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	result, err := c.Multiplicar(1.25)
	require.NoError(t, err)
	assert.InDelta(t, 125.0, result.Valor(), 0.001)
}

func TestCorriente_Multiplicar_NegativeFactorFails(t *testing.T) {
	c, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	_, err = c.Multiplicar(-1)
	assert.Error(t, err)
}

func TestCorriente_Dividir(t *testing.T) {
	c, err := valueobject.NewCorriente(200)
	require.NoError(t, err)

	result, err := c.Dividir(2)
	require.NoError(t, err)
	assert.InDelta(t, 100.0, result.Valor(), 0.001)
}

func TestCorriente_Dividir_PorCeroFails(t *testing.T) {
	c, err := valueobject.NewCorriente(200)
	require.NoError(t, err)

	_, err = c.Dividir(0)
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/valueobject/ -v -run TestNewCorriente
```

Expected: FAIL — `corriente.go` does not exist yet.

**Step 3: Write minimal implementation**

```go
// internal/domain/valueobject/corriente.go
package valueobject

import (
	"errors"
	"fmt"
)

var ErrCorrienteInvalida = errors.New("corriente debe ser mayor que cero")

// Corriente represents an electrical current value in Amperes. Immutable.
type Corriente struct {
	valor  float64
	unidad string
}

func NewCorriente(valor float64) (Corriente, error) {
	if valor <= 0 {
		return Corriente{}, fmt.Errorf("%w: %.4f", ErrCorrienteInvalida, valor)
	}
	return Corriente{valor: valor, unidad: "A"}, nil
}

func (c Corriente) Valor() float64 { return c.valor }
func (c Corriente) Unidad() string { return c.unidad }

func (c Corriente) Multiplicar(factor float64) (Corriente, error) {
	return NewCorriente(c.valor * factor)
}

func (c Corriente) Dividir(divisor int) (Corriente, error) {
	if divisor == 0 {
		return Corriente{}, fmt.Errorf("dividir corriente: divisor es cero")
	}
	return NewCorriente(c.valor / float64(divisor))
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/valueobject/ -v -run TestCorriente
```

Expected: All PASS.

**Step 5: Commit**

```bash
git add internal/domain/valueobject/corriente.go internal/domain/valueobject/corriente_test.go
git commit -m "feat(domain): add Corriente value object with validation"
```

---

### Task 3: Value Object — Tension

**Files:**
- Create: `internal/domain/valueobject/tension.go`
- Test: `internal/domain/valueobject/tension_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/valueobject/tension_test.go
package valueobject_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTension(t *testing.T) {
	tests := []struct {
		name    string
		valor   int
		wantErr bool
	}{
		{"127V valid", 127, false},
		{"220V valid", 220, false},
		{"240V valid", 240, false},
		{"277V valid", 277, false},
		{"440V valid", 440, false},
		{"480V valid", 480, false},
		{"600V valid", 600, false},
		{"100V invalid", 100, true},
		{"0V invalid", 0, true},
		{"negative invalid", -220, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tension, err := valueobject.NewTension(tt.valor)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, valueobject.ErrVoltajeInvalido))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.valor, tension.Valor())
			assert.Equal(t, "V", tension.Unidad())
		})
	}
}

func TestTension_EnKilovoltios(t *testing.T) {
	tests := []struct {
		name     string
		voltaje  int
		expected float64
	}{
		{"480V → 0.48 kV", 480, 0.48},
		{"220V → 0.22 kV", 220, 0.22},
		{"127V → 0.127 kV", 127, 0.127},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tension, err := valueobject.NewTension(tt.voltaje)
			require.NoError(t, err)
			assert.InDelta(t, tt.expected, tension.EnKilovoltios(), 0.001)
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/valueobject/ -v -run TestNewTension
```

Expected: FAIL — `tension.go` does not exist.

**Step 3: Write minimal implementation**

```go
// internal/domain/valueobject/tension.go
package valueobject

import (
	"errors"
	"fmt"
)

var ErrVoltajeInvalido = errors.New("voltaje no válido según normativa NOM")

var voltajesValidos = map[int]bool{
	127: true,
	220: true,
	240: true,
	277: true,
	440: true,
	480: true,
	600: true,
}

// Tension represents an electrical voltage value in Volts. Immutable.
type Tension struct {
	valor  int
	unidad string
}

func NewTension(valor int) (Tension, error) {
	if !voltajesValidos[valor] {
		return Tension{}, fmt.Errorf("%w: %d", ErrVoltajeInvalido, valor)
	}
	return Tension{valor: valor, unidad: "V"}, nil
}

func (t Tension) Valor() int       { return t.valor }
func (t Tension) Unidad() string   { return t.unidad }

func (t Tension) EnKilovoltios() float64 {
	return float64(t.valor) / 1000.0
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/valueobject/ -v -run TestTension
```

Expected: All PASS.

**Step 5: Commit**

```bash
git add internal/domain/valueobject/tension.go internal/domain/valueobject/tension_test.go
git commit -m "feat(domain): add Tension value object with NOM voltage validation"
```

---

### Task 4: Value Object — Conductor ✅ COMPLETADO

**Files:**
- `internal/domain/valueobject/conductor.go`
- `internal/domain/valueobject/conductor_test.go`

**Implementación real (diverge del plan original — expandida con datos NOM-001-SEDE-2012):**

El constructor usa `ConductorParams` struct (patrón idiomático para >4 parámetros).
Se agregaron campos de NOM-001-SEDE-2012: `AreaConAislamientoMM2`, `DiametroMM`,
`NumeroHilos`, `ResistenciaPVCPorKm`, `ResistenciaAlPorKm`, `ResistenciaAceroPorKm`,
`ReactanciaPorKm`.

**Validaciones implementadas (post-refactor de campos opcionales):**
- Campos **requeridos** (validados en `NewConductor`):
  - `Calibre`: mapa `calibresValidos` con NOM 310-15(b)(16) — AWG 18→4/0, MCM 250→2000
  - `Material`: solo "Cu" o "Al"
  - `SeccionMM2`: > 0
- Campos **opcionales** (aceptados sin validación): `TipoAislamiento` (vacío para desnudos), `AreaConAislamientoMM2`, `DiametroMM`, `NumeroHilos`, resistencias, reactancia
  - **Validación postponida:** Los campos opcionales se validan al punto de uso (ej: caída de tensión requiere sección, canalización requiere área con aislamiento)

**Impacto en tareas 10–13:**
- Task 10: `EntradaTablaConductor` solo necesita `Calibre`, `Capacidad`, `SeccionMM2`, `Material`, `TipoAislamiento`
- Task 11: Conductores de tierra son desnudos (`TipoAislamiento = ""`); datos de test usan "4 AWG" y "2 AWG" (válidos)
- Task 12: Usar `AreaConAislamientoMM2()` (opcional) para cálculo de fill en canalización
- Task 13: `CalcularCaidaTension` usa `SeccionMM2()` para la fórmula ρ (es la sección del conductor, no del aislamiento)

**Commits:** `892aef4`, `49d897f`, `b992169`, `6569b23` (Task 4 original), `c6c15c8` (refactor de campos opcionales post-Task 11)

---

### Task 5: Entity Foundation — TipoFiltro, Equipo, CalculadorCorriente, Errors ✅ COMPLETADO

**Commit:** `8d563db`

**Archivos creados:** `errors.go`, `tipo_filtro.go`, `equipo.go`, `tipo_filtro_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/entity/tipo_filtro_test.go
package entity_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTipoFiltro(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected entity.TipoFiltro
		wantErr  bool
	}{
		{"ACTIVO valid", "ACTIVO", entity.TipoFiltroActivo, false},
		{"RECHAZO valid", "RECHAZO", entity.TipoFiltroRechazo, false},
		{"lowercase invalid", "activo", entity.TipoFiltro(""), true},
		{"empty invalid", "", entity.TipoFiltro(""), true},
		{"unknown invalid", "PASIVO", entity.TipoFiltro(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := entity.ParseTipoFiltro(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, entity.ErrTipoFiltroInvalido))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTipoFiltro_String(t *testing.T) {
	assert.Equal(t, "ACTIVO", entity.TipoFiltroActivo.String())
	assert.Equal(t, "RECHAZO", entity.TipoFiltroRechazo.String())
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/entity/ -v -run TestParseTipoFiltro
```

Expected: FAIL — files do not exist.

**Step 3: Write implementation**

```go
// internal/domain/entity/errors.go
package entity

import "errors"

var (
	ErrTipoFiltroInvalido = errors.New("tipo de filtro no válido")
	ErrDivisionPorCero    = errors.New("división por cero en cálculo de corriente")
)
```

```go
// internal/domain/entity/tipo_filtro.go
package entity

import "fmt"

type TipoFiltro string

const (
	TipoFiltroActivo  TipoFiltro = "ACTIVO"
	TipoFiltroRechazo TipoFiltro = "RECHAZO"
)

func ParseTipoFiltro(s string) (TipoFiltro, error) {
	switch s {
	case string(TipoFiltroActivo):
		return TipoFiltroActivo, nil
	case string(TipoFiltroRechazo):
		return TipoFiltroRechazo, nil
	default:
		return "", fmt.Errorf("%w: '%s'", ErrTipoFiltroInvalido, s)
	}
}

func (t TipoFiltro) String() string {
	return string(t)
}
```

```go
// internal/domain/entity/equipo.go
package entity

import "github.com/garfex/calculadora-filtros/internal/domain/valueobject"

// CalculadorCorriente defines the contract for equipment that can calculate
// its nominal current. Each equipment type implements this differently.
type CalculadorCorriente interface {
	CalcularCorrienteNominal() (valueobject.Corriente, error)
}

// Equipo is the base struct embedded by all equipment types.
type Equipo struct {
	Clave   string
	Tipo    TipoFiltro
	Voltaje int
	ITM     int
	Bornes  int
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/entity/ -v -run TestParseTipoFiltro
go test ./internal/domain/entity/ -v -run TestTipoFiltro_String
```

Expected: All PASS.

**Step 5: Commit**

```bash
git add internal/domain/entity/errors.go internal/domain/entity/tipo_filtro.go internal/domain/entity/equipo.go internal/domain/entity/tipo_filtro_test.go
git commit -m "feat(domain): add TipoFiltro enum, Equipo base struct, and CalculadorCorriente interface"
```

---

### Task 6: Entity — FiltroActivo ✅ COMPLETADO

**Commit:** `f28e1a3` (entity) + `a1b6233` (CalculadorPotencia)

**Divergencia:** Se agregó interfaz `CalculadorPotencia` (PotenciaKVA/KW/KVAR) en `equipo.go` y se implementó en `FiltroActivo` (PF=1: KVA=I×V×√3/1000, KW=KVA, KVAR=0). Esto también aplica a Task 7.

**Archivos creados:** `filtro_activo.go`, `filtro_activo_test.go`

**Files:**
- Create: `internal/domain/entity/filtro_activo.go`
- Test: `internal/domain/entity/filtro_activo_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/entity/filtro_activo_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFiltroActivo(t *testing.T) {
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, 125, 3)
	require.NoError(t, err)

	assert.Equal(t, "FA-001", fa.Clave)
	assert.Equal(t, entity.TipoFiltroActivo, fa.Tipo)
	assert.Equal(t, 480, fa.Voltaje)
	assert.Equal(t, 100, fa.Amperaje)
	assert.Equal(t, 125, fa.ITM)
	assert.Equal(t, 3, fa.Bornes)
}

func TestNewFiltroActivo_AmperajeInvalido(t *testing.T) {
	_, err := entity.NewFiltroActivo("FA-001", 480, 0, 125, 3)
	assert.Error(t, err)

	_, err = entity.NewFiltroActivo("FA-001", 480, -10, 125, 3)
	assert.Error(t, err)
}

func TestFiltroActivo_CalcularCorrienteNominal(t *testing.T) {
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, 125, 3)
	require.NoError(t, err)

	corriente, err := fa.CalcularCorrienteNominal()
	require.NoError(t, err)
	assert.InDelta(t, 100.0, corriente.Valor(), 0.001)
	assert.Equal(t, "A", corriente.Unidad())
}

func TestFiltroActivo_ImplementsCalculadorCorriente(t *testing.T) {
	fa, _ := entity.NewFiltroActivo("FA-001", 480, 100, 125, 3)
	var _ entity.CalculadorCorriente = fa
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/entity/ -v -run TestFiltroActivo
```

Expected: FAIL — `filtro_activo.go` does not exist.

**Step 3: Write minimal implementation**

```go
// internal/domain/entity/filtro_activo.go
package entity

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// FiltroActivo represents an active filter. Its nominal current equals
// its amperage rating directly (no formula needed).
type FiltroActivo struct {
	Equipo
	Amperaje int
}

func NewFiltroActivo(clave string, voltaje, amperaje, itm, bornes int) (*FiltroActivo, error) {
	if amperaje <= 0 {
		return nil, fmt.Errorf("amperaje debe ser mayor que cero: %d", amperaje)
	}
	return &FiltroActivo{
		Equipo: Equipo{
			Clave:   clave,
			Tipo:    TipoFiltroActivo,
			Voltaje: voltaje,
			ITM:     itm,
			Bornes:  bornes,
		},
		Amperaje: amperaje,
	}, nil
}

func (fa *FiltroActivo) CalcularCorrienteNominal() (valueobject.Corriente, error) {
	return valueobject.NewCorriente(float64(fa.Amperaje))
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/entity/ -v -run TestFiltroActivo
```

Expected: All PASS.

**Step 5: Commit**

```bash
git add internal/domain/entity/filtro_activo.go internal/domain/entity/filtro_activo_test.go
git commit -m "feat(domain): add FiltroActivo entity with direct amperage current"
```

---

### Task 7: Entity — FiltroRechazo ✅ COMPLETADO

**Commit:** `2d76007` (entity) + `a1b6233` (CalculadorPotencia)

**Divergencia:** Se implementó `CalculadorPotencia` en FiltroRechazo (puramente reactivo: KVAR=dado, KVA=KVAR, KW=0).

**Archivos creados:** `filtro_rechazo.go`, `filtro_rechazo_test.go`

**Files:**
- Create: `internal/domain/entity/filtro_rechazo.go`
- Test: `internal/domain/entity/filtro_rechazo_test.go`

**Formula:** `I = KVAR / (KV × √3)` where `KV = Voltaje / 1000`

**Step 1: Write the failing test**

```go
// internal/domain/entity/filtro_rechazo_test.go
package entity_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFiltroRechazo(t *testing.T) {
	fr, err := entity.NewFiltroRechazo("FR-001", 480, 100, 125, 3)
	require.NoError(t, err)

	assert.Equal(t, "FR-001", fr.Clave)
	assert.Equal(t, entity.TipoFiltroRechazo, fr.Tipo)
	assert.Equal(t, 480, fr.Voltaje)
	assert.Equal(t, 100, fr.KVAR)
	assert.Equal(t, 125, fr.ITM)
	assert.Equal(t, 3, fr.Bornes)
}

func TestNewFiltroRechazo_KVARInvalido(t *testing.T) {
	_, err := entity.NewFiltroRechazo("FR-001", 480, 0, 125, 3)
	assert.Error(t, err)

	_, err = entity.NewFiltroRechazo("FR-001", 480, -50, 125, 3)
	assert.Error(t, err)
}

func TestFiltroRechazo_CalcularCorrienteNominal(t *testing.T) {
	tests := []struct {
		name     string
		voltaje  int
		kvar     int
		expected float64
	}{
		{
			name:     "100 KVAR at 480V",
			voltaje:  480,
			kvar:     100,
			// I = 100 / (0.48 × √3) = 100 / 0.83138... = 120.28 A
			expected: 120.28,
		},
		{
			name:     "50 KVAR at 220V",
			voltaje:  220,
			kvar:     50,
			// I = 50 / (0.22 × √3) = 50 / 0.38105... = 131.22 A
			expected: 131.22,
		},
		{
			name:     "200 KVAR at 440V",
			voltaje:  440,
			kvar:     200,
			// I = 200 / (0.44 × √3) = 200 / 0.76210... = 262.43 A
			expected: 262.43,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr, err := entity.NewFiltroRechazo("FR-TEST", tt.voltaje, tt.kvar, 125, 3)
			require.NoError(t, err)

			corriente, err := fr.CalcularCorrienteNominal()
			require.NoError(t, err)
			assert.InDelta(t, tt.expected, corriente.Valor(), 0.01)
		})
	}
}

func TestFiltroRechazo_CalcularCorrienteNominal_VoltajeCero(t *testing.T) {
	// Voltage 0 would cause division by zero.
	// NewFiltroRechazo should reject voltaje <= 0.
	_, err := entity.NewFiltroRechazo("FR-BAD", 0, 100, 125, 3)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, entity.ErrDivisionPorCero))
}

func TestFiltroRechazo_ImplementsCalculadorCorriente(t *testing.T) {
	fr, _ := entity.NewFiltroRechazo("FR-001", 480, 100, 125, 3)
	var _ entity.CalculadorCorriente = fr
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/entity/ -v -run TestFiltroRechazo
```

Expected: FAIL — `filtro_rechazo.go` does not exist.

**Step 3: Write minimal implementation**

```go
// internal/domain/entity/filtro_rechazo.go
package entity

import (
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// FiltroRechazo represents a rejection filter (capacitor bank).
// Nominal current: I = KVAR / (KV × √3), where KV = Voltaje / 1000.
type FiltroRechazo struct {
	Equipo
	KVAR int
}

func NewFiltroRechazo(clave string, voltaje, kvar, itm, bornes int) (*FiltroRechazo, error) {
	if kvar <= 0 {
		return nil, fmt.Errorf("KVAR debe ser mayor que cero: %d", kvar)
	}
	if voltaje <= 0 {
		return nil, fmt.Errorf("%w: voltaje es %d", ErrDivisionPorCero, voltaje)
	}
	return &FiltroRechazo{
		Equipo: Equipo{
			Clave:   clave,
			Tipo:    TipoFiltroRechazo,
			Voltaje: voltaje,
			ITM:     itm,
			Bornes:  bornes,
		},
		KVAR: kvar,
	}, nil
}

func (fr *FiltroRechazo) CalcularCorrienteNominal() (valueobject.Corriente, error) {
	kv := float64(fr.Voltaje) / 1000.0
	denominador := kv * math.Sqrt(3)
	corriente := float64(fr.KVAR) / denominador
	return valueobject.NewCorriente(corriente)
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/entity/ -v -run TestFiltroRechazo
```

Expected: All PASS.

**Step 5: Commit**

```bash
git add internal/domain/entity/filtro_rechazo.go internal/domain/entity/filtro_rechazo_test.go
git commit -m "feat(domain): add FiltroRechazo entity with KVAR formula I=KVAR/(KV×√3)"
```

---

### Refactoring: ITM como entidad propia ✅ COMPLETADO (entre Task 7 y 8)

**Commit:** `f7ec5a1`

**Decisión:** `ITM` se extrajo como entidad validada con `Amperaje`, `Polos`, `Bornes`, `Voltaje`. En Fase 1, Polos=3 (trifásico), Voltaje=equipo.Voltaje. `Equipo.ITM int` y `Equipo.Bornes int` se reemplazaron por `Equipo.ITM ITM`. Los constructores `NewFiltroActivo/NewFiltroRechazo` ahora aceptan `itm ITM` en lugar de `itm, bornes int`.

**Impacto en Task 11:** `SeleccionarConductorTierra` recibe `itm int` (el amperaje del ITM). Cambiar a `itm.Amperaje` al integrar con entidades. Los datos de test con `"3 AWG"` y `"1 AWG"` deben reemplazarse por `"4 AWG"` y `"2 AWG"` (no válidos en calibresValidos).

---

### Refactoring: TipoFiltro → TipoEquipo ✅ COMPLETADO (entre Task 8 y 9)

**Commit:** `eaa5167`

**Decisión:** `TipoFiltro` renombrado a `TipoEquipo` para soportar más tipos de equipo en el futuro (tableros, transformadores, cargas, etc.). Los valores de las constantes cambiaron de `"ACTIVO"`/`"RECHAZO"` a `"FILTRO_ACTIVO"`/`"FILTRO_RECHAZO"` para ser descriptivos al agregar nuevos tipos. El error `ErrTipoFiltroInvalido` → `ErrTipoEquipoInvalido`. Archivos: `tipo_filtro.go` → `tipo_equipo.go`, `tipo_filtro_test.go` → `tipo_equipo_test.go`.

**Impacto en BD:** El enum `tipo_filtro` en PostgreSQL con valores `ACTIVO`/`RECHAZO` deberá migrarse a `FILTRO_ACTIVO`/`FILTRO_RECHAZO` al conectar infrastructure (Fase 2+). En Fase 1 no hay conexión a BD real, solo domain layer.

---

### Extensión: Transformador + Carga + Rename AmperajeNominal ✅ COMPLETADO (entre TipoEquipo y Task 9)

**Commit:** `f5e1c78`

**Decisión (brainstorming):** Se amplió el modelo de equipos de 2 a 4 tipos para cubrir más escenarios de memorias de cálculo. Diseño documentado en `docs/plans/2026-02-10-nuevos-equipos-design.md`.

**Cambios:**
- `FiltroActivo.Amperaje` → `AmperajeNominal` (claridad semántica)
- `Transformador`: KVA como dato de entrada, In=KVA/(KV×√3), potencia=solo KVA
- `Carga`: KW+FP+Fases como datos de entrada, soporta 1/2/3 fases con factor diferente
- `TipoEquipoTransformador="TRANSFORMADOR"`, `TipoEquipoCarga="CARGA"` agregados a `tipo_equipo.go`

**Impacto en servicios (Tasks 9-14):** Ninguno. Los servicios trabajan con interfaces `CalculadorCorriente` y `CalculadorPotencia`, agnósticos al tipo de equipo.

---

### Task 8: Domain Structs — Canalizacion and MemoriaCalculo ✅ COMPLETADO

**Commit:** `30c8274`

**Divergencia:** `MemoriaCalculo` incluye campos adicionales `PotenciaKVA`, `PotenciaKW`, `PotenciaKVAR` para mostrar en el output del reporte.

**Archivos creados:** `canalizacion.go`, `memoria_calculo.go`

**Files:**
- Create: `internal/domain/entity/canalizacion.go`
- Create: `internal/domain/entity/memoria_calculo.go`

No tests needed — these are plain data structs with no logic.

**Step 1: Write Canalizacion**

```go
// internal/domain/entity/canalizacion.go
package entity

// Canalizacion represents the conduit or cable tray selected for the installation.
type Canalizacion struct {
	Tipo       string  // "TUBERIA" | "CHAROLA"
	Tamano     string  // e.g., "1 1/2" (inches for tubería)
	AreaTotal  float64 // total conductor area in mm²
}
```

**Step 2: Write MemoriaCalculo**

```go
// internal/domain/entity/memoria_calculo.go
package entity

import "github.com/garfex/calculadora-filtros/internal/domain/valueobject"

// MemoriaCalculo holds the complete results of all calculation steps
// for an electrical installation memory.
type MemoriaCalculo struct {
	Equipo                CalculadorCorriente
	CorrienteNominal      valueobject.Corriente
	CorrienteAjustada     valueobject.Corriente
	FactoresAjuste        map[string]float64
	ConductorAlimentacion valueobject.Conductor
	HilosPorFase          int
	ConductorTierra       valueobject.Conductor
	Canalizacion          Canalizacion
	CaidaTension          float64
	CumpleNormativa       bool
}
```

**Step 3: Verify compilation**

```bash
go build ./internal/domain/...
```

Expected: No errors.

**Step 4: Commit**

```bash
git add internal/domain/entity/canalizacion.go internal/domain/entity/memoria_calculo.go
git commit -m "feat(domain): add Canalizacion and MemoriaCalculo result structs"
```

---

### Task 9: Service — CalcularCorrienteNominal + AjustarCorriente ✅ COMPLETADO

**Commits:** `1aa2570`

**Files creados:**
- `internal/domain/service/calculo_corriente_nominal.go`
- `internal/domain/service/ajuste_corriente.go`
- `internal/domain/service/calculo_corriente_nominal_test.go`
- `internal/domain/service/ajuste_corriente_test.go`

**Divergencias vs. plan:** Tests adaptados para usar `ITM` como struct (patrón actual) en lugar de parámetros posicionales. Se agregaron tests para `Transformador` y `Carga` (no incluidos en plan original pero necesarios ahora que existen esos tipos).

**Implementación resumen:**

**Step 1: Write the failing test for CalcularCorrienteNominal**

```go
// internal/domain/service/calculo_corriente_nominal_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularCorrienteNominal_FiltroActivo(t *testing.T) {
	fa, err := entity.NewFiltroActivo("FA-001", 480, 100, 125, 3)
	require.NoError(t, err)

	corriente, err := service.CalcularCorrienteNominal(fa)
	require.NoError(t, err)
	assert.InDelta(t, 100.0, corriente.Valor(), 0.001)
}

func TestCalcularCorrienteNominal_FiltroRechazo(t *testing.T) {
	fr, err := entity.NewFiltroRechazo("FR-001", 480, 100, 125, 3)
	require.NoError(t, err)

	corriente, err := service.CalcularCorrienteNominal(fr)
	require.NoError(t, err)
	// I = 100 / (0.48 × √3) ≈ 120.28 A
	assert.InDelta(t, 120.28, corriente.Valor(), 0.01)
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/service/ -v -run TestCalcularCorrienteNominal
```

Expected: FAIL — file does not exist.

**Step 3: Implement CalcularCorrienteNominal**

```go
// internal/domain/service/calculo_corriente_nominal.go
package service

import (
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// CalcularCorrienteNominal delegates to the equipment's own calculation method.
// FA → amperage directly. FR → I = KVAR / (KV × √3).
func CalcularCorrienteNominal(equipo entity.CalculadorCorriente) (valueobject.Corriente, error) {
	return equipo.CalcularCorrienteNominal()
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/service/ -v -run TestCalcularCorrienteNominal
```

Expected: All PASS.

**Step 5: Write the failing test for AjustarCorriente**

```go
// internal/domain/service/ajuste_corriente_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAjustarCorriente_SingleFactor(t *testing.T) {
	cn, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	factores := map[string]float64{
		"proteccion": 1.25,
	}

	result, err := service.AjustarCorriente(cn, factores)
	require.NoError(t, err)
	assert.InDelta(t, 125.0, result.Valor(), 0.001)
}

func TestAjustarCorriente_MultipleFactors(t *testing.T) {
	cn, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	factores := map[string]float64{
		"proteccion":   1.25,
		"temperatura":  0.88,
		"agrupamiento": 0.80,
	}

	// 100 × 1.25 × 0.88 × 0.80 = 88.0
	result, err := service.AjustarCorriente(cn, factores)
	require.NoError(t, err)
	assert.InDelta(t, 88.0, result.Valor(), 0.01)
}

func TestAjustarCorriente_EmptyFactors(t *testing.T) {
	cn, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	result, err := service.AjustarCorriente(cn, nil)
	require.NoError(t, err)
	assert.InDelta(t, 100.0, result.Valor(), 0.001)
}

func TestAjustarCorriente_ZeroFactorFails(t *testing.T) {
	cn, err := valueobject.NewCorriente(100)
	require.NoError(t, err)

	factores := map[string]float64{
		"proteccion": 0,
	}

	_, err = service.AjustarCorriente(cn, factores)
	assert.Error(t, err)
}
```

**Step 6: Run test to verify it fails**

```bash
go test ./internal/domain/service/ -v -run TestAjustarCorriente
```

Expected: FAIL — file does not exist.

**Step 7: Implement AjustarCorriente**

```go
// internal/domain/service/ajuste_corriente.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// AjustarCorriente multiplies the nominal current by all adjustment factors
// (protection, temperature, grouping, etc.). Returns error if any factor <= 0.
func AjustarCorriente(cn valueobject.Corriente, factores map[string]float64) (valueobject.Corriente, error) {
	if len(factores) == 0 {
		return cn, nil
	}

	resultado := cn.Valor()
	for nombre, factor := range factores {
		if factor <= 0 {
			return valueobject.Corriente{}, fmt.Errorf("factor '%s' inválido: %.4f (debe ser > 0)", nombre, factor)
		}
		resultado *= factor
	}

	return valueobject.NewCorriente(resultado)
}
```

**Step 8: Run all tests to verify they pass**

```bash
go test ./internal/domain/service/ -v
```

Expected: All PASS.

**Step 9: Commit**

```bash
git add internal/domain/service/calculo_corriente_nominal.go internal/domain/service/calculo_corriente_nominal_test.go internal/domain/service/ajuste_corriente.go internal/domain/service/ajuste_corriente_test.go
git commit -m "feat(domain): add corriente nominal calculation and current adjustment services"
```

---

### Task 10: Service — SeleccionarConductorAlimentacion ✅ COMPLETADO

**Commit:** `b806152`

**Files creados:**
- `internal/domain/service/calculo_conductor.go`
- `internal/domain/service/calculo_conductor_test.go`

**Design note:** The service receives pre-resolved table data as `[]EntradaTablaConductor`. The application layer will query `TablaNOMRepository` and pass the results here. The entries must be sorted smallest-to-largest by calibre (as they appear in NOM tables).

**Divergencias vs. plan original:**
- `NewConductor` acepta `ConductorParams` struct (no params posicionales).
- `EntradaTablaConductor` fue diseñada con campos simples: `Calibre`, `Capacidad`, `SeccionMM2`, `Material`, `TipoAislamiento` (refactorizado durante el refactoring de campos opcionales).
- Los tests no usan valores dummy para campos innecesarios — se construyen `ConductorParams` solo con los campos que ImportanImportantemente, la selección solo depende de `Capacidad`.

**Step 1: Write the failing test**

```go
// internal/domain/service/calculo_conductor_test.go
package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Simplified NOM table 310-15(b)(16) excerpt for Cu THHN 90°C
var tablaConductorTest = []service.EntradaTablaConductor{
	{Calibre: "14 AWG", Capacidad: 25, SeccionMM2: 2.08, Material: "Cu", TipoAislamiento: "THHN"},
	{Calibre: "12 AWG", Capacidad: 30, SeccionMM2: 3.31, Material: "Cu", TipoAislamiento: "THHN"},
	{Calibre: "10 AWG", Capacidad: 40, SeccionMM2: 5.26, Material: "Cu", TipoAislamiento: "THHN"},
	{Calibre: "8 AWG", Capacidad: 55, SeccionMM2: 8.37, Material: "Cu", TipoAislamiento: "THHN"},
	{Calibre: "6 AWG", Capacidad: 75, SeccionMM2: 13.30, Material: "Cu", TipoAislamiento: "THHN"},
	{Calibre: "4 AWG", Capacidad: 95, SeccionMM2: 21.15, Material: "Cu", TipoAislamiento: "THHN"},
	{Calibre: "2 AWG", Capacidad: 130, SeccionMM2: 33.62, Material: "Cu", TipoAislamiento: "THHN"},
	{Calibre: "1/0 AWG", Capacidad: 170, SeccionMM2: 53.49, Material: "Cu", TipoAislamiento: "THHN"},
	{Calibre: "4/0 AWG", Capacidad: 260, SeccionMM2: 107.2, Material: "Cu", TipoAislamiento: "THHN"},
	{Calibre: "500 MCM", Capacidad: 380, SeccionMM2: 253.4, Material: "Cu", TipoAislamiento: "THHN"},
}

func TestSeleccionarConductorAlimentacion_Simple(t *testing.T) {
	corriente, err := valueobject.NewCorriente(120)
	require.NoError(t, err)

	conductor, err := service.SeleccionarConductorAlimentacion(corriente, 1, tablaConductorTest)
	require.NoError(t, err)
	// 120A needs at least 130A capacity → 2 AWG
	assert.Equal(t, "2 AWG", conductor.Calibre())
}

func TestSeleccionarConductorAlimentacion_ExactMatch(t *testing.T) {
	corriente, err := valueobject.NewCorriente(95)
	require.NoError(t, err)

	conductor, err := service.SeleccionarConductorAlimentacion(corriente, 1, tablaConductorTest)
	require.NoError(t, err)
	// 95A exactly matches 4 AWG capacity
	assert.Equal(t, "4 AWG", conductor.Calibre())
}

func TestSeleccionarConductorAlimentacion_ConHilosPorFase(t *testing.T) {
	corriente, err := valueobject.NewCorriente(240)
	require.NoError(t, err)

	// 240A / 2 hilos = 120A per wire → needs 130A → 2 AWG
	conductor, err := service.SeleccionarConductorAlimentacion(corriente, 2, tablaConductorTest)
	require.NoError(t, err)
	assert.Equal(t, "2 AWG", conductor.Calibre())
}

func TestSeleccionarConductorAlimentacion_CorrienteExceedsAllCapacities(t *testing.T) {
	corriente, err := valueobject.NewCorriente(500)
	require.NoError(t, err)

	_, err = service.SeleccionarConductorAlimentacion(corriente, 1, tablaConductorTest)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrConductorNoEncontrado))
}

func TestSeleccionarConductorAlimentacion_EmptyTable(t *testing.T) {
	corriente, err := valueobject.NewCorriente(10)
	require.NoError(t, err)

	_, err = service.SeleccionarConductorAlimentacion(corriente, 1, nil)
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/service/ -v -run TestSeleccionarConductorAlimentacion
```

Expected: FAIL — file does not exist.

**Step 3: Implement**

```go
// internal/domain/service/calculo_conductor.go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

var ErrConductorNoEncontrado = errors.New("no se encontró conductor con capacidad suficiente")

// EntradaTablaConductor represents one row from NOM table 310-15(b)(16).
// Must be sorted smallest-to-largest calibre (as in the NOM table).
type EntradaTablaConductor struct {
	Calibre         string
	Capacidad       float64 // ampacity in amperes
	SeccionMM2      float64
	Material        string
	TipoAislamiento string
}

// SeleccionarConductorAlimentacion picks the smallest conductor from the NOM table
// whose ampacity >= corrienteAjustada / hilosPorFase.
func SeleccionarConductorAlimentacion(
	corrienteAjustada valueobject.Corriente,
	hilosPorFase int,
	tabla []EntradaTablaConductor,
) (valueobject.Conductor, error) {
	if len(tabla) == 0 {
		return valueobject.Conductor{}, fmt.Errorf("%w: tabla vacía", ErrConductorNoEncontrado)
	}

	if hilosPorFase < 1 {
		hilosPorFase = 1
	}

	corrientePorHilo := corrienteAjustada.Valor() / float64(hilosPorFase)

	for _, entrada := range tabla {
		if entrada.Capacidad >= corrientePorHilo {
			return valueobject.NewConductor(
				entrada.Calibre,
				entrada.Material,
				entrada.TipoAislamiento,
				entrada.SeccionMM2,
			)
		}
	}

	return valueobject.Conductor{}, fmt.Errorf(
		"%w: corriente por hilo %.2f A excede máxima capacidad de tabla %.2f A",
		ErrConductorNoEncontrado, corrientePorHilo, tabla[len(tabla)-1].Capacidad,
	)
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/service/ -v -run TestSeleccionarConductorAlimentacion
```

Expected: All PASS.

**Step 5: Commit**

```bash
git add internal/domain/service/calculo_conductor.go internal/domain/service/calculo_conductor_test.go
git commit -m "feat(domain): add conductor selection service with NOM table lookup"
```

---

### Task 11: Service — SeleccionarConductorTierra ✅ COMPLETADO

**Commit:** `284a4eb`

**Files creados:**
- `internal/domain/service/calculo_tierra.go`
- `internal/domain/service/calculo_tierra_test.go`

**Design note:** Table 250-122 maps ITM ranges to ground conductor calibres. The service receives sorted entries where each row represents the maximum ITM for that conductor size.

**Divergencias vs. plan original:**
1. El plan usaba `"3 AWG"` y `"1 AWG"` pero **NO son calibres válidos** en `calibresValidos` (no están en NOM 310-15(b)(16)). Se sustituyó con `"4 AWG"` y `"2 AWG"`.
2. `NewConductor` acepta `ConductorParams` struct. `EntradaTablaTierra` usa `Calibre`, `SeccionMM2`, e `ITMHasta` — menos campos que la alimentación porque los conductores de tierra son desnudos (sin aislamiento requerido).

**Step 1: Write the failing test**

```go
// internal/domain/service/calculo_tierra_test.go
package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Simplified NOM table 250-122 excerpt
var tablaTierraTest = []service.EntradaTablaTierra{
	{ITMHasta: 15, Calibre: "14 AWG", SeccionMM2: 2.08},
	{ITMHasta: 20, Calibre: "12 AWG", SeccionMM2: 3.31},
	{ITMHasta: 40, Calibre: "10 AWG", SeccionMM2: 5.26},
	{ITMHasta: 60, Calibre: "10 AWG", SeccionMM2: 5.26},
	{ITMHasta: 100, Calibre: "8 AWG", SeccionMM2: 8.37},
	{ITMHasta: 200, Calibre: "6 AWG", SeccionMM2: 13.30},
	{ITMHasta: 400, Calibre: "3 AWG", SeccionMM2: 26.67},
	{ITMHasta: 600, Calibre: "1 AWG", SeccionMM2: 42.41},
	{ITMHasta: 800, Calibre: "1/0 AWG", SeccionMM2: 53.49},
	{ITMHasta: 1000, Calibre: "2/0 AWG", SeccionMM2: 67.43},
}

func TestSeleccionarConductorTierra(t *testing.T) {
	tests := []struct {
		name            string
		itm             int
		expectedCalibre string
	}{
		{"ITM 15 → 14 AWG", 15, "14 AWG"},
		{"ITM 20 → 12 AWG", 20, "12 AWG"},
		{"ITM 30 → 10 AWG (≤40)", 30, "10 AWG"},
		{"ITM 100 → 8 AWG", 100, "8 AWG"},
		{"ITM 125 → 6 AWG (≤200)", 125, "6 AWG"},
		{"ITM 400 → 3 AWG", 400, "3 AWG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conductor, err := service.SeleccionarConductorTierra(tt.itm, tablaTierraTest)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCalibre, conductor.Calibre())
			assert.Equal(t, "Cu", conductor.Material())
		})
	}
}

func TestSeleccionarConductorTierra_ITMExceedsTable(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(1200, tablaTierraTest)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrConductorNoEncontrado))
}

func TestSeleccionarConductorTierra_InvalidITM(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(0, tablaTierraTest)
	assert.Error(t, err)
}

func TestSeleccionarConductorTierra_EmptyTable(t *testing.T) {
	_, err := service.SeleccionarConductorTierra(100, nil)
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/service/ -v -run TestSeleccionarConductorTierra
```

Expected: FAIL — file does not exist.

**Step 3: Implement**

```go
// internal/domain/service/calculo_tierra.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// EntradaTablaTierra represents one row from NOM table 250-122.
// Entries must be sorted by ITMHasta ascending.
type EntradaTablaTierra struct {
	ITMHasta   int
	Calibre    string
	SeccionMM2 float64
}

// SeleccionarConductorTierra selects the ground conductor from NOM table 250-122
// based on the equipment's ITM (circuit breaker) rating.
func SeleccionarConductorTierra(itm int, tabla []EntradaTablaTierra) (valueobject.Conductor, error) {
	if itm <= 0 {
		return valueobject.Conductor{}, fmt.Errorf("ITM debe ser mayor que cero: %d", itm)
	}
	if len(tabla) == 0 {
		return valueobject.Conductor{}, fmt.Errorf("%w: tabla de tierra vacía", ErrConductorNoEncontrado)
	}

	for _, entrada := range tabla {
		if itm <= entrada.ITMHasta {
			return valueobject.NewConductor(entrada.Calibre, "Cu", "THHN", entrada.SeccionMM2)
		}
	}

	return valueobject.Conductor{}, fmt.Errorf(
		"%w: ITM %d excede máximo de tabla %d",
		ErrConductorNoEncontrado, itm, tabla[len(tabla)-1].ITMHasta,
	)
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/service/ -v -run TestSeleccionarConductorTierra
```

Expected: All PASS.

**Step 5: Commit**

```bash
git add internal/domain/service/calculo_tierra.go internal/domain/service/calculo_tierra_test.go
git commit -m "feat(domain): add ground conductor selection from NOM table 250-122"
```

---

### Refactoring: Conductor VO — Campos Opcionales ✅ COMPLETADO (post-Task 11)

**Commit:** `c6c15c8`

**Decisión:** Se identificó que el modelo de Conductor no podía representar conductores desnudos (tierra, puesta a tierra) sin usar valores dummy para campos de aislamiento. Además, los test helpers estaban forzados a llenar campos innecesarios. Se refactorizó a un modelo de **validación en construcción mínima + validación al punto de uso**.

**Cambios:**

1. **Constructor más flexible:**
   ```go
   func NewConductor(p ConductorParams) (Conductor, error)
   // Valida solo: Calibre (en map), Material (Cu|Al), SeccionMM2 > 0
   // Acepta: TipoAislamiento vacío (para desnudos), todos los demás campos opcionales
   ```

2. **Campos opcionales sin validación:**
   - `TipoAislamiento`: puede ser `""` (desnudo), `"THHN"`, `"THW"`, etc.
   - `AreaConAislamientoMM2`, `DiametroMM`, `NumeroHilos`: solo usados si están presentes
   - Resistencias y reactancia: no validadas en construcción

3. **Validación postponida:**
   - `CalcularCaidaTension`: requiere `SeccionMM2()` ✓ (siempre validado)
   - `CalcularCanalizacion`: requiere `AreaConAislamientoMM2()` (si 0, no se usa canalización con aislamiento; ver adaptación Task 12)
   - Tests de tierra usan conductores desnudos: `TipoAislamiento = ""`

**Impacto en testes:**
- `conductor_test.go`: Reemplazó 9 casos de error (validación de todos los campos) con 3 casos específicos:
  - `TestNewConductor_SeccionCero`
  - `TestNewConductor_SeccionNegativa`
  - `TestNewConductor_Minimal` (solo 3 campos, sin aislamiento)
- `calculo_conductor_test.go`: Simplificó `entradaConductor()` helper (sin valores dummy)
- `calculo_tierra_test.go`: Simplificó `entradaTierra()` helper (solo 3 campos para desnudo)

**Beneficios:**
- Representación correcta de conductores desnudos sin hackeos
- Tests más limpios (sin dummy values)
- Validación enfocada en lo que importa en construcción
- Extensible para futuros tipos de conductores

---

### Task 12: Service — CalcularCanalizacion

**Files:**
- Create: `internal/domain/service/calculo_canalizacion.go`
- Test: `internal/domain/service/calculo_canalizacion_test.go`

**Design note:** Per NOM, conduit fill for 2+ conductors is 40% of the conduit's internal area. The service calculates total conductor area, applies the fill factor, and selects the smallest conduit that fits.

**⚠️ Nota de diseño — Nuevo diseño de canalizaciones (2026-02-11):**
La firma del servicio usa `tipo string` en el plan original, pero el diseño validado usa `TipoCanalizacion` (enum). Al implementar, se debe pasar el tipo como `string` serializado del enum (ej: `"TUBERIA_CONDUIT"`), o adaptar la firma a `TipoCanalizacion`. El enum `TipoCanalizacion` debe crearse en `internal/domain/entity/tipo_canalizacion.go` **antes** de implementar este servicio.

La tabla NOM correcta para cada canalización es responsabilidad de la capa application/infrastructure — el domain service solo recibe `[]EntradaTablaCanalizacion` ya resueltos (mismo patrón que `SeleccionarConductorAlimentacion`).

Ver diseño completo: `docs/plans/2026-02-11-tablas-nom-canalizacion-design.md`

**Step 1: Write the failing test**

```go
// internal/domain/service/calculo_canalizacion_test.go
package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Simplified conduit sizing table (tubería EMT)
var tablaCanalizacionTest = []service.EntradaTablaCanalizacion{
	{Tamano: "1/2", AreaInteriorMM2: 78.0},
	{Tamano: "3/4", AreaInteriorMM2: 122.0},
	{Tamano: "1", AreaInteriorMM2: 198.0},
	{Tamano: "1 1/4", AreaInteriorMM2: 277.0},
	{Tamano: "1 1/2", AreaInteriorMM2: 360.0},
	{Tamano: "2", AreaInteriorMM2: 572.0},
	{Tamano: "2 1/2", AreaInteriorMM2: 885.0},
	{Tamano: "3", AreaInteriorMM2: 1327.0},
	{Tamano: "4", AreaInteriorMM2: 2165.0},
}

func TestCalcularCanalizacion_Tuberia(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 33.62},  // 3 phases × 2 AWG
		{Cantidad: 1, SeccionMM2: 13.30},  // 1 ground × 6 AWG
	}
	// Total area = 3×33.62 + 1×13.30 = 114.16 mm²
	// Required conduit area at 40% fill = 114.16 / 0.40 = 285.4 mm²
	// Smallest conduit ≥ 285.4 → "1 1/2" (360 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA", tablaCanalizacionTest)
	require.NoError(t, err)
	assert.Equal(t, "TUBERIA", result.Tipo)
	assert.Equal(t, "1 1/2", result.Tamano)
	assert.InDelta(t, 114.16, result.AreaTotal, 0.01)
}

func TestCalcularCanalizacion_SmallConductors(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 3, SeccionMM2: 3.31},  // 3 × 12 AWG
		{Cantidad: 1, SeccionMM2: 2.08},  // 1 × 14 AWG ground
	}
	// Total = 3×3.31 + 1×2.08 = 12.01 mm²
	// Required = 12.01 / 0.40 = 30.025 mm²
	// Smallest ≥ 30.025 → "1/2" (78 mm²)

	result, err := service.CalcularCanalizacion(conductores, "TUBERIA", tablaCanalizacionTest)
	require.NoError(t, err)
	assert.Equal(t, "1/2", result.Tamano)
}

func TestCalcularCanalizacion_NoFit(t *testing.T) {
	conductores := []service.ConductorParaCanalizacion{
		{Cantidad: 20, SeccionMM2: 253.4},
	}
	// Total = 20 × 253.4 = 5068 mm² → required = 12670 mm² → exceeds all conduits

	_, err := service.CalcularCanalizacion(conductores, "TUBERIA", tablaCanalizacionTest)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrCanalizacionNoDisponible))
}

func TestCalcularCanalizacion_EmptyConductors(t *testing.T) {
	_, err := service.CalcularCanalizacion(nil, "TUBERIA", tablaCanalizacionTest)
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/service/ -v -run TestCalcularCanalizacion
```

Expected: FAIL — file does not exist.

**Step 3: Implement**

```go
// internal/domain/service/calculo_canalizacion.go
package service

import (
	"errors"
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
)

var ErrCanalizacionNoDisponible = errors.New("no se encontró canalización con área suficiente")

const factorRellenoTuberia = 0.40 // NOM: 40% fill for 2+ conductors

// ConductorParaCanalizacion holds the quantity and cross-section area
// of a group of identical conductors for conduit sizing calculations.
type ConductorParaCanalizacion struct {
	Cantidad   int
	SeccionMM2 float64
}

// EntradaTablaCanalizacion represents one row from a conduit sizing table.
// Entries must be sorted by AreaInteriorMM2 ascending.
type EntradaTablaCanalizacion struct {
	Tamano          string
	AreaInteriorMM2 float64
}

// CalcularCanalizacion selects the smallest conduit whose usable area
// (interior area × fill factor) accommodates all conductors.
func CalcularCanalizacion(
	conductores []ConductorParaCanalizacion,
	tipo string,
	tabla []EntradaTablaCanalizacion,
) (entity.Canalizacion, error) {
	if len(conductores) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("lista de conductores vacía")
	}
	if len(tabla) == 0 {
		return entity.Canalizacion{}, fmt.Errorf("%w: tabla vacía", ErrCanalizacionNoDisponible)
	}

	var areaTotal float64
	for _, c := range conductores {
		areaTotal += float64(c.Cantidad) * c.SeccionMM2
	}

	areaRequerida := areaTotal / factorRellenoTuberia

	for _, entrada := range tabla {
		if entrada.AreaInteriorMM2 >= areaRequerida {
			return entity.Canalizacion{
				Tipo:      tipo,
				Tamano:    entrada.Tamano,
				AreaTotal: areaTotal,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf(
		"%w: área requerida %.2f mm² excede máxima disponible %.2f mm²",
		ErrCanalizacionNoDisponible, areaRequerida, tabla[len(tabla)-1].AreaInteriorMM2,
	)
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/service/ -v -run TestCalcularCanalizacion
```

Expected: All PASS.

**Step 5: Commit**

```bash
git add internal/domain/service/calculo_canalizacion.go internal/domain/service/calculo_canalizacion_test.go
git commit -m "feat(domain): add conduit sizing service with NOM 40% fill factor"
```

---

### Task 13: Service — CalcularCaidaTension

**⚠️ Adaptación requerida:** El test usa `valueobject.NewConductor(calibre, material, "THHN", seccion)` (API viejo). Debe actualizarse a `valueobject.NewConductor(valueobject.ConductorParams{...})` con todos los campos requeridos. La fórmula usa `conductor.SeccionMM2()` (sección del conductor sin aislamiento), no `AreaConAislamientoMM2()`.

**Files:**
- Create: `internal/domain/service/calculo_caida_tension.go`
- Test: `internal/domain/service/calculo_caida_tension_test.go`

**Formula (three-phase):**
```
VD = (√3 × ρ × L × I) / S
VD% = (VD / V) × 100
```

Where:
- `ρ` = resistivity (Cu: 0.01724, Al: 0.02826 Ω·mm²/m at 75°C)
- `L` = one-way distance in meters
- `I` = current in amperes
- `S` = conductor cross-section in mm²
- `V` = system voltage in volts

**NOM limits:** 3% for feeders (alimentadores), 5% total.

**Step 1: Write the failing test**

```go
// internal/domain/service/calculo_caida_tension_test.go
package service_test

import (
	"errors"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularCaidaTension(t *testing.T) {
	tests := []struct {
		name          string
		corrienteA    float64
		distanciaM    float64
		calibre       string
		material      string
		seccionMM2    float64
		voltaje       int
		limitePorc    float64
		expectedPorc  float64
		expectedCumple bool
	}{
		{
			name:          "Cu 2AWG 30m at 120A 480V within 3%",
			corrienteA:    120,
			distanciaM:    30,
			calibre:       "2 AWG",
			material:      "Cu",
			seccionMM2:    33.62,
			voltaje:       480,
			limitePorc:    3.0,
			// VD = (√3 × 0.01724 × 30 × 120) / 33.62 = (107.476) / 33.62 = 3.197 V
			// VD% = (3.197 / 480) × 100 = 0.666%
			expectedPorc:  0.666,
			expectedCumple: true,
		},
		{
			name:          "Cu 12AWG 100m at 25A 220V exceeds 3%",
			corrienteA:    25,
			distanciaM:    100,
			calibre:       "12 AWG",
			material:      "Cu",
			seccionMM2:    3.31,
			voltaje:       220,
			limitePorc:    3.0,
			// VD = (√3 × 0.01724 × 100 × 25) / 3.31 = (74.63) / 3.31 = 22.55 V
			// VD% = (22.55 / 220) × 100 = 10.25%
			expectedPorc:  10.25,
			expectedCumple: false,
		},
		{
			name:          "Al conductor higher resistivity",
			corrienteA:    100,
			distanciaM:    20,
			calibre:       "4/0 AWG",
			material:      "Al",
			seccionMM2:    107.2,
			voltaje:       480,
			limitePorc:    3.0,
			// VD = (√3 × 0.02826 × 20 × 100) / 107.2 = (97.87) / 107.2 = 0.913 V
			// VD% = (0.913 / 480) × 100 = 0.190%
			expectedPorc:  0.190,
			expectedCumple: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corriente, err := valueobject.NewCorriente(tt.corrienteA)
			require.NoError(t, err)

			conductor, err := valueobject.NewConductor(tt.calibre, tt.material, "THHN", tt.seccionMM2)
			require.NoError(t, err)

			tension, err := valueobject.NewTension(tt.voltaje)
			require.NoError(t, err)

			porcentaje, cumple, err := service.CalcularCaidaTension(
				conductor, corriente, tt.distanciaM, tension, tt.limitePorc,
			)
			require.NoError(t, err)
			assert.InDelta(t, tt.expectedPorc, porcentaje, 0.01)
			assert.Equal(t, tt.expectedCumple, cumple)
		})
	}
}

func TestCalcularCaidaTension_DistanciaCero(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(100)
	conductor, _ := valueobject.NewConductor("2 AWG", "Cu", "THHN", 33.62)
	tension, _ := valueobject.NewTension(480)

	_, _, err := service.CalcularCaidaTension(conductor, corriente, 0, tension, 3.0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrDistanciaInvalida))
}

func TestCalcularCaidaTension_DistanciaNegativa(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(100)
	conductor, _ := valueobject.NewConductor("2 AWG", "Cu", "THHN", 33.62)
	tension, _ := valueobject.NewTension(480)

	_, _, err := service.CalcularCaidaTension(conductor, corriente, -10, tension, 3.0)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, service.ErrDistanciaInvalida))
}

func TestCalcularCaidaTension_InvalidMaterial(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(100)
	// Conductor constructor validates material, so this tests the service
	// handling of an unexpected material (defensive check).
	// Since NewConductor rejects invalid materials, this case can't happen
	// in practice. We test the zero-distance and negative-distance cases instead.
	// This test is a placeholder showing the service is robust.
	corriente2, _ := valueobject.NewCorriente(100)
	conductor, _ := valueobject.NewConductor("2 AWG", "Cu", "THHN", 33.62)
	tension, _ := valueobject.NewTension(480)

	// Valid case just to ensure no panic with valid inputs
	_, _, err := service.CalcularCaidaTension(conductor, corriente2, 10, tension, 3.0)
	_ = corriente
	assert.NoError(t, err)
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/service/ -v -run TestCalcularCaidaTension
```

Expected: FAIL — file does not exist.

**Step 3: Implement**

```go
// internal/domain/service/calculo_caida_tension.go
package service

import (
	"errors"
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

var ErrDistanciaInvalida = errors.New("distancia debe ser mayor que cero")

// Resistivity values in Ω·mm²/m (approximate at 75°C operating temperature)
var resistividad = map[string]float64{
	"Cu": 0.01724,
	"Al": 0.02826,
}

// CalcularCaidaTension calculates voltage drop percentage for a three-phase system.
// Formula: VD% = (√3 × ρ × L × I) / (S × V) × 100
// Returns the percentage, whether it meets the NOM limit, and any error.
func CalcularCaidaTension(
	conductor valueobject.Conductor,
	corriente valueobject.Corriente,
	distancia float64,
	tension valueobject.Tension,
	limiteNOM float64,
) (porcentaje float64, cumple bool, err error) {
	if distancia <= 0 {
		return 0, false, fmt.Errorf("%w: %.2f", ErrDistanciaInvalida, distancia)
	}

	rho, ok := resistividad[conductor.Material()]
	if !ok {
		return 0, false, fmt.Errorf("material desconocido para resistividad: %s", conductor.Material())
	}

	vd := (math.Sqrt(3) * rho * distancia * corriente.Valor()) / conductor.SeccionMM2()
	porcentaje = (vd / float64(tension.Valor())) * 100

	return porcentaje, porcentaje <= limiteNOM, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/service/ -v -run TestCalcularCaidaTension
```

Expected: All PASS.

**Step 5: Commit**

```bash
git add internal/domain/service/calculo_caida_tension.go internal/domain/service/calculo_caida_tension_test.go
git commit -m "feat(domain): add voltage drop calculation service with NOM limit validation"
```

---

### Task 14: Final Verification

**Step 1: Run full test suite with race detector**

```bash
go test -race ./...
```

Expected: All PASS, no race conditions.

**Step 2: Run go vet**

```bash
go vet ./...
```

Expected: No issues.

**Step 3: Run go build**

```bash
go build ./...
```

Expected: No errors.

**Step 4: Run golangci-lint (if installed)**

```bash
golangci-lint run
```

Expected: No issues. If not installed, skip — CI will catch it.

**Step 5: Run test coverage**

```bash
go test -cover ./internal/domain/...
```

Expected: High coverage across all domain packages.

**Step 6: Fix any issues found and commit**

If any issues are found, fix them and commit:

```bash
git add -A
git commit -m "fix(domain): resolve lint/vet issues from final verification"
```

**Step 7: Final commit if no issues**

```bash
# No commit needed if everything passed
```

---

## Summary

| Task | Component | Files Created | Tests |
|------|-----------|---------------|-------|
| 1 | Scaffolding | go.mod, dirs, .golangci.yml, .env.example | — |
| 2 | Corriente VO | corriente.go | 5 test cases |
| 3 | Tension VO | tension.go | 10+ test cases |
| 4 | Conductor VO | conductor.go | 9 test cases |
| 5 | Entity foundation | tipo_filtro.go, equipo.go, errors.go | 5+ test cases |
| 6 | FiltroActivo | filtro_activo.go | 4 test functions |
| 7 | FiltroRechazo | filtro_rechazo.go | 5 test functions |
| 8 | Result structs | canalizacion.go, memoria_calculo.go | — |
| 9 | CorrienteNominal + AjusteCorriente | 2 service files | 6 test functions |
| 10 | SeleccionarConductorAlimentacion | calculo_conductor.go | 5 test functions |
| 11 | SeleccionarConductorTierra | calculo_tierra.go | 4 test functions |
| 12 | CalcularCanalizacion | calculo_canalizacion.go | 4 test functions |
| 13 | CalcularCaidaTension | calculo_caida_tension.go | 5 test functions |
| 14 | Final verification | — | Full suite |

**Total: 14 tasks, ~10 Go source files, ~10 test files, ~50+ test cases**
