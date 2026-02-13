# Caída de Tensión IEEE-141 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Reemplazar la fórmula de impedancia pura por la fórmula IEEE-141/NOM con factor de potencia en `CalcularCaidaTension`.

**Architecture:** Servicio de dominio puro. `EntradaCalculoCaidaTension` se simplifica (elimina campos geométricos, agrega `ReactanciaOhmPorKm` y `FactorPotencia`). La fórmula pasa de 7 pasos a 5. El struct `ResultadoCaidaTension` mantiene los mismos campos pero `Impedancia` cambia de semántica a "término efectivo R·cosθ + X·senθ". Todo en `internal/domain/service/` y `internal/domain/entity/`.

**Tech Stack:** Go 1.22+, testify

**Reference:** `docs/plans/2026-02-12-caida-tension-ieee141-design.md`

---

## Contexto para el implementador

La fórmula **actual** (método impedancia Z):
```
%VD = (√3 × I × √(R² + X²) × L) / V × 100
```

La fórmula **nueva** (IEEE-141 / NOM con FP):
```
%e = 173 × (In/CF) × L_km × (R·cosθ + X·senθ) / E_FF
VD = E_FF × (%e / 100)
```

`173 = √3 × 100`, `cosθ = FactorPotencia`, `senθ = √(1 - FP²)`

El struct `EntradaCalculoCaidaTension` vive en `internal/domain/service/calculo_caida_tension.go`.
El struct `ResultadoCaidaTension` vive en `internal/domain/entity/memoria_calculo.go`.

---

### Task 1: Actualizar `EntradaCalculoCaidaTension` — tests RED primero

**Files:**
- Modify: `internal/domain/service/calculo_caida_tension_test.go`

El test existente usa los campos viejos (`DiametroExteriorMM`, `DiametroConductorMM`, `NumeroHilos`).
Hay que reescribirlo con el struct nuevo antes de tocar el código.

**Step 1: Reemplazar el contenido completo del archivo de test**

Sobreescribir `internal/domain/service/calculo_caida_tension_test.go` con:

```go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalcularCaidaTension_FormulaIEEE141 verifica la fórmula NOM con FP.
//
// Fórmula: %e = 173 × (In/CF) × L_km × (R·cosθ + X·senθ) / E_FF
//
// Caso referencia FP=1 (FA/FR/TR): 2 AWG Cu, tubería PVC, 120A, 30m, 480V
//   R = 0.62 Ω/km, X = 0.051 Ω/km (Tabla 9), FP = 1.0
//   term = 0.62×1.0 + 0.051×0.0 = 0.62
//   %e   = 173 × 120 × 0.030 × 0.62 / 480 = 387.07 / 480 / 100 × 100 = 0.807%
//   VD   = 480 × 0.00807 = 3.874 V
//
// Caso FP=0.85 (Carga): mismo conductor
//   cosθ = 0.85, senθ = √(1-0.7225) = 0.5268
//   term = 0.62×0.85 + 0.051×0.5268 = 0.527 + 0.02687 = 0.5539
//   %e   = 173 × 120 × 0.030 × 0.5539 / 480 = 344.83 / 480 = 0.719%
//   VD   = 480 × 0.00719 = 3.451 V
func TestCalcularCaidaTension_FormulaIEEE141(t *testing.T) {
	tests := []struct {
		name               string
		entrada            service.EntradaCalculoCaidaTension
		corrienteA         float64
		distanciaM         float64
		voltaje            int
		limitePorc         float64
		expectedPorcentaje float64
		expectedVD         float64
		expectedCumple     bool
		skipNumeric        bool
	}{
		{
			name: "FP=1.0 (FiltroActivo/FiltroRechazo/Transformador) - solo resistencia importa",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				HilosPorFase:        1,
				FactorPotencia:      1.0,
			},
			corrienteA:         120,
			distanciaM:         30,
			voltaje:            480,
			limitePorc:         3.0,
			expectedPorcentaje: 0.807,
			expectedVD:         3.874,
			expectedCumple:     true,
		},
		{
			name: "FP=0.85 (Carga) - reactancia contribuye",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				HilosPorFase:        1,
				FactorPotencia:      0.85,
			},
			corrienteA:         120,
			distanciaM:         30,
			voltaje:            480,
			limitePorc:         3.0,
			expectedPorcentaje: 0.719,
			expectedVD:         3.451,
			expectedCumple:     true,
		},
		{
			name: "2 hilos por fase reduce R y X a la mitad",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				HilosPorFase:        2,
				FactorPotencia:      1.0,
			},
			corrienteA: 120,
			distanciaM: 30,
			voltaje:    480,
			limitePorc: 3.0,
			// R_ef = 0.62/2 = 0.31, X_ef = 0.051/2 = 0.0255
			// term = 0.31×1.0 + 0.0255×0.0 = 0.31
			// %e   = 173 × 120 × 0.030 × 0.31 / 480 = 0.4035%
			expectedPorcentaje: 0.4035,
			expectedVD:         1.937,
			expectedCumple:     true,
		},
		{
			name: "charola espaciado - misma formula, X de reactancia_al",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionCharolaCableEspaciado,
				HilosPorFase:        1,
				FactorPotencia:      1.0,
			},
			corrienteA:         120,
			distanciaM:         30,
			voltaje:            480,
			limitePorc:         3.0,
			expectedPorcentaje: 0.807,
			expectedVD:         3.874,
			expectedCumple:     true,
		},
		{
			name: "excede limite NOM 3%",
			entrada: service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 5.21,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				HilosPorFase:        1,
				FactorPotencia:      1.0,
			},
			corrienteA:     25,
			distanciaM:     100,
			voltaje:        220,
			limitePorc:     3.0,
			expectedCumple: false,
			skipNumeric:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corriente, err := valueobject.NewCorriente(tt.corrienteA)
			require.NoError(t, err)

			tension, err := valueobject.NewTension(tt.voltaje)
			require.NoError(t, err)

			resultado, err := service.CalcularCaidaTension(
				tt.entrada, corriente, tt.distanciaM, tension, tt.limitePorc,
			)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCumple, resultado.Cumple)

			if !tt.skipNumeric {
				assert.InDelta(t, tt.expectedPorcentaje, resultado.Porcentaje, 0.01, "porcentaje")
				assert.InDelta(t, tt.expectedVD, resultado.CaidaVolts, 0.01, "caida volts")
			}
		})
	}
}

func TestCalcularCaidaTension_ResultadoContieneRXTerminoEfectivo(t *testing.T) {
	// Verifica que el struct expone R_ef, X_ef y el término efectivo para el reporte
	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.051,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		HilosPorFase:        1,
		FactorPotencia:      0.85,
	}
	corriente, _ := valueobject.NewCorriente(120)
	tension, _ := valueobject.NewTension(480)

	resultado, err := service.CalcularCaidaTension(entrada, corriente, 30, tension, 3.0)
	require.NoError(t, err)

	// R_ef = 0.62 / 1 = 0.62
	assert.InDelta(t, 0.62, resultado.Resistencia, 0.001, "R_ef")
	// X_ef = 0.051 / 1 = 0.051
	assert.InDelta(t, 0.051, resultado.Reactancia, 0.001, "X_ef")
	// término efectivo = R·cosθ + X·senθ = 0.62×0.85 + 0.051×0.5268 = 0.5539
	assert.InDelta(t, 0.5539, resultado.Impedancia, 0.002, "término efectivo")
}

func TestCalcularCaidaTension_ErrorDistanciaInvalida(t *testing.T) {
	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.051,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		HilosPorFase:        1,
		FactorPotencia:      1.0,
	}
	corriente, _ := valueobject.NewCorriente(120)
	tension, _ := valueobject.NewTension(480)

	t.Run("distancia cero", func(t *testing.T) {
		_, err := service.CalcularCaidaTension(entrada, corriente, 0, tension, 3.0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrDistanciaInvalida)
	})

	t.Run("distancia negativa", func(t *testing.T) {
		_, err := service.CalcularCaidaTension(entrada, corriente, -10, tension, 3.0)
		require.Error(t, err)
		assert.ErrorIs(t, err, service.ErrDistanciaInvalida)
	})
}

func TestCalcularCaidaTension_ErrorHilosPorFaseInvalido(t *testing.T) {
	entrada := service.EntradaCalculoCaidaTension{
		ResistenciaOhmPorKm: 0.62,
		ReactanciaOhmPorKm:  0.051,
		TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
		HilosPorFase:        0,
		FactorPotencia:      1.0,
	}
	corriente, _ := valueobject.NewCorriente(120)
	tension, _ := valueobject.NewTension(480)

	_, err := service.CalcularCaidaTension(entrada, corriente, 30, tension, 3.0)
	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrHilosPorFaseInvalido)
}

func TestCalcularCaidaTension_ErrorFactorPotenciaInvalido(t *testing.T) {
	corriente, _ := valueobject.NewCorriente(120)
	tension, _ := valueobject.NewTension(480)

	casos := []struct {
		nombre string
		fp     float64
	}{
		{"FP cero", 0.0},
		{"FP negativo", -0.5},
		{"FP mayor que 1", 1.1},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			entrada := service.EntradaCalculoCaidaTension{
				ResistenciaOhmPorKm: 0.62,
				ReactanciaOhmPorKm:  0.051,
				TipoCanalizacion:    entity.TipoCanalizacionTuberiaPVC,
				HilosPorFase:        1,
				FactorPotencia:      c.fp,
			}
			_, err := service.CalcularCaidaTension(entrada, corriente, 30, tension, 3.0)
			require.Error(t, err)
			assert.ErrorIs(t, err, service.ErrFactorPotenciaInvalido)
		})
	}
}
```

**Step 2: Correr los tests para verificar RED**

```bash
go test ./internal/domain/service/ -run TestCalcularCaidaTension -v
```

Esperado: errores de compilación por campos inexistentes en `EntradaCalculoCaidaTension`
(`ReactanciaOhmPorKm`, `FactorPotencia`) y error sentinel nuevo (`ErrFactorPotenciaInvalido`).

---

### Task 2: Reescribir `CalcularCaidaTension` — implementación GREEN

**Files:**
- Modify: `internal/domain/service/calculo_caida_tension.go`

**Step 1: Sobreescribir el archivo completo con la nueva implementación**

```go
package service

import (
	"errors"
	"fmt"
	"math"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ErrDistanciaInvalida is returned when the distance is zero or negative.
var ErrDistanciaInvalida = errors.New("distancia debe ser mayor que cero")

// ErrHilosPorFaseInvalido is returned when HilosPorFase is zero or negative.
var ErrHilosPorFaseInvalido = errors.New("hilos por fase debe ser mayor que cero")

// ErrFactorPotenciaInvalido is returned when FactorPotencia is not in range (0, 1].
var ErrFactorPotenciaInvalido = errors.New("factor de potencia debe estar entre 0 (exclusivo) y 1")

// EntradaCalculoCaidaTension contains the pre-resolved NOM table data needed
// to calculate voltage drop using the IEEE-141 / NOM formula with power factor.
// The application layer is responsible for resolving R, X from Tabla 9 and
// the power factor from the equipment entity.
type EntradaCalculoCaidaTension struct {
	ResistenciaOhmPorKm float64                 // Tabla 9 → res_{material}_{conduit}
	ReactanciaOhmPorKm  float64                 // Tabla 9 → reactancia_al o reactancia_acero
	TipoCanalizacion    entity.TipoCanalizacion  // Documented in memoria de cálculo report
	HilosPorFase        int                     // CF ≥ 1 (parallel conductors per phase)
	FactorPotencia      float64                 // cosθ: FA/FR/TR = 1.0 | Carga = explicit FP
}

// CalcularCaidaTension calculates the voltage drop for a three-phase system
// using the IEEE-141 / NOM formula with power factor:
//
//	%e = 173 × (In/CF) × L_km × (R·cosθ + X·senθ) / E_FF
//	VD = E_FF × (%e / 100)
//
// Where 173 = √3 × 100, cosθ = FactorPotencia, senθ = √(1 - FP²).
//
// For FP = 1.0 (FiltroActivo, FiltroRechazo, Transformador) the formula
// reduces to: %e = 173 × (In/CF) × L_km × R / E_FF  (reactance has no effect).
func CalcularCaidaTension(
	entrada EntradaCalculoCaidaTension,
	corriente valueobject.Corriente,
	distancia float64,
	tension valueobject.Tension,
	limiteNOM float64,
) (entity.ResultadoCaidaTension, error) {
	if distancia <= 0 {
		return entity.ResultadoCaidaTension{}, fmt.Errorf("CalcularCaidaTension: %w: %.2f", ErrDistanciaInvalida, distancia)
	}
	if entrada.HilosPorFase <= 0 {
		return entity.ResultadoCaidaTension{}, fmt.Errorf("CalcularCaidaTension: %w: %d", ErrHilosPorFaseInvalido, entrada.HilosPorFase)
	}
	if entrada.FactorPotencia <= 0 || entrada.FactorPotencia > 1 {
		return entity.ResultadoCaidaTension{}, fmt.Errorf("CalcularCaidaTension: %w: %.4f", ErrFactorPotenciaInvalido, entrada.FactorPotencia)
	}

	n := float64(entrada.HilosPorFase)

	// Step 1-2: angle components
	cosTheta := entrada.FactorPotencia
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)

	// Step 3-4: effective R and X per parallel conductor
	rEf := entrada.ResistenciaOhmPorKm / n
	xEf := entrada.ReactanciaOhmPorKm / n

	// Step 5: effective impedance term (Ω/km)
	terminoEfectivo := rEf*cosTheta + xEf*sinTheta

	// Step 6: %e = 173 × (In/CF) × L_km × terminoEfectivo / E_FF
	// Note: corriente already represents In (total, not per conductor).
	// CF (HilosPorFase) is already applied to R and X above.
	lKm := distancia / 1000.0
	porcentaje := 173 * corriente.Valor() * lKm * terminoEfectivo / float64(tension.Valor())

	// Step 7: VD in volts
	vd := float64(tension.Valor()) * (porcentaje / 100)

	return entity.ResultadoCaidaTension{
		Porcentaje:  porcentaje,
		CaidaVolts:  vd,
		Cumple:      porcentaje <= limiteNOM,
		Impedancia:  terminoEfectivo, // R·cosθ + X·senθ — "effective impedance term"
		Resistencia: rEf,
		Reactancia:  xEf,
	}, nil
}
```

**Step 2: Correr los tests**

```bash
go test ./internal/domain/service/ -run TestCalcularCaidaTension -v
```

Esperado: todos PASS.

**Step 3: Correr toda la suite**

```bash
go test ./...
```

Esperado: `ok` en los 3 paquetes (entity, service, valueobject).

**Step 4: Commit**

```bash
git add internal/domain/service/calculo_caida_tension.go \
        internal/domain/service/calculo_caida_tension_test.go
git commit -m "feat(domain): replace impedance method with IEEE-141 voltage drop formula"
```

---

### Task 3: Actualizar comentario semántico en `ResultadoCaidaTension`

**Files:**
- Modify: `internal/domain/entity/memoria_calculo.go`

El campo `Impedancia` ahora representa el término efectivo `R·cosθ + X·senθ`, no `√(R²+X²)`.
Actualizar solo el comentario del campo.

**Step 1: Editar el comentario**

Cambiar la línea del campo `Impedancia` en `ResultadoCaidaTension`:

```go
// Antes:
Impedancia  float64 // Z (Ω/km) = √(R² + X²)

// Después:
Impedancia  float64 // Término efectivo (Ω/km) = R·cosθ + X·senθ  (IEEE-141)
```

**Step 2: Compilar para verificar**

```bash
go build ./...
```

Esperado: sin errores.

**Step 3: Commit**

```bash
git add internal/domain/entity/memoria_calculo.go
git commit -m "docs(domain): update Impedancia field comment to reflect IEEE-141 effective term"
```

---

### Task 4: Verificación final

**Step 1: Correr toda la suite completa**

```bash
go test ./...
```

Esperado:
```
ok  github.com/garfex/calculadora-filtros/internal/domain/entity
ok  github.com/garfex/calculadora-filtros/internal/domain/service
ok  github.com/garfex/calculadora-filtros/internal/domain/valueobject
```

**Step 2: Verificar que compila limpio**

```bash
go build ./...
go vet ./...
```

Esperado: sin output (sin errores ni warnings).

**Step 3: Commit final si hay cambios pendientes**

```bash
git status
# Si hay algo sin commitear:
git add -A
git commit -m "chore: final cleanup after IEEE-141 voltage drop implementation"
```
