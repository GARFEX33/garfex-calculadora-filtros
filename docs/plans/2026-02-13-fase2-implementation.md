# Fase 2: Memoria de Cálculo Completa - Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implementar cálculo automático de factores (temperatura y agrupamiento) y fórmulas NOM correctas para canalización (tubería y charola).

**Architecture:** Domain services puros para cálculos, application layer orquesta con datos de estado y sistema eléctrico, infrastructure lee tablas CSV NOM. Hexagonal architecture mantenida.

**Tech Stack:** Go 1.22+, testify, CSV parsing

**Reference Design:** `docs/plans/2026-02-13-fase2-memoria-calculo-design.md`

---

## Task 1: Crear Tabla de Estados y Temperaturas

**Files:**

- Create: `data/tablas_nom/estados_temperatura.csv`
- Create: `data/tablas_nom/estados_temperatura_test.go` (validación de formato)

**Step 1: Crear archivo CSV con temperaturas Maxima**

```csv
estado,temperatura_max_c
Aguascalientes,18
Baja California,22
Baja California Sur,25
Campeche,26
Coahuila,21
Colima,26
Chiapas,24
Chihuahua,19
Ciudad de Mexico,16
Durango,17
Guanajuato,19
Guerrero,26
Hidalgo,18
Jalisco,21
Mexico,16
Michoacan,21
Morelos,21
Nayarit,25
Nuevo Leon,21
Oaxaca,24
Puebla,18
Queretaro,19
Quintana Roo,26
San Luis Potosi,20
Sinaloa,25
Sonora,24
Tabasco,27
Tamaulipas,24
Tlaxcala,16
Veracruz,24
Yucatan,26
Zacatecas,17
```

**Step 2: Validar formato del CSV**

```bash
cat data/tablas_nom/estados_temperatura.csv | head -5
```

Expected: Muestra header + 4 primeras líneas con formato correcto.

**Step 3: Commit**

```bash
git add data/tablas_nom/estados_temperatura.csv
git commit -m "feat: add estados_temperatura.csv with 32 Mexican states"
```

---

## Task 2: Crear Tabla NOM 310-15(b)(2)(a) - Factores por Temperatura

**Files:**

- Create: `data/tablas_nom/310-15-b-2-a.csv`

**Step 1: Crear archivo CSV con factores de corrección**

```csv
rango_temp_c,factor_60c,factor_75c,factor_90c
10-15,1.20,1.15,1.12
16-20,1.15,1.11,1.09
21-25,1.10,1.07,1.05
26-30,1.05,1.03,1.02
31-35,1.00,1.00,1.00
36-40,0.94,0.95,0.96
41-45,0.88,0.90,0.91
46-50,0.82,0.85,0.87
51-55,0.75,0.80,0.82
56-60,0.67,0.74,0.77
61-70,0.58,0.67,0.71
71-80,0.47,0.58,0.63
```

**Step 2: Validar formato**

```bash
cat data/tablas_nom/310-15-b-2-a.csv | head -5
```

Expected: Header correcto y 4 primeros rangos.

**Step 3: Commit**

```bash
git add data/tablas_nom/310-15-b-2-a.csv
git commit -m "feat: add NOM 310-15(b)(2)(a) temperature correction factors"
```

---

## Task 3: Crear Tabla NOM 310-15(b)(3)(a) - Factores por Agrupamiento

**Files:**

- Create: `data/tablas_nom/310-15-b-3-a.csv`

**Step 1: Crear archivo CSV con factores de ajuste**

```csv
cantidad_conductores,factor
1,1.00
2,0.80
3,0.70
4,0.65
5-6,0.60
7-9,0.50
10-20,0.45
21-30,0.40
31-40,0.35
41+,0.30
```

**Step 2: Commit**

```bash
git add data/tablas_nom/310-15-b-3-a.csv
git commit -m "feat: add NOM 310-15(b)(3)(a) grouping adjustment factors"
```

---

## Task 4: Crear Tabla de Dimensiones de Charola

**Files:**

- Create: `data/tablas_nom/charola_dimensiones.csv`

**Step 1: Crear archivo CSV con tamaños comerciales**

```csv
tamano_pulgadas,ancho_mm
6,152.4
9,228.6
12,304.8
16,406.4
18,457.2
20,508.0
24,609.6
30,762.0
36,914.4
```

**Step 2: Commit**

```bash
git add data/tablas_nom/charola_dimensiones.csv
git commit -m "feat: add charola cable tray standard sizes table"
```

---

## Task 5: Agregar SistemaElectrico a Domain (Entity)

**Files:**

- Create: `internal/domain/entity/sistema_electrico.go`
- Create: `internal/domain/entity/sistema_electrico_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/entity/sistema_electrico_test.go
package entity_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestSistemaElectrico_CantidadConductores(t *testing.T) {
	tests := []struct {
		sistema  entity.SistemaElectrico
		expected int
	}{
		{entity.SistemaElectricoDelta, 3},
		{entity.SistemaElectricoEstrella, 4},
		{entity.SistemaElectricoBifasico, 3},
		{entity.SistemaElectricoMonofasico, 2},
	}

	for _, tt := range tests {
		t.Run(string(tt.sistema), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.sistema.CantidadConductores())
		})
	}
}

func TestParseSistemaElectrico(t *testing.T) {
	tests := []struct {
		input    string
		expected entity.SistemaElectrico
		wantErr  bool
	}{
		{"DELTA", entity.SistemaElectricoDelta, false},
		{"ESTRELLA", entity.SistemaElectricoEstrella, false},
		{"BIFASICO", entity.SistemaElectricoBifasico, false},
		{"MONOFASICO", entity.SistemaElectricoMonofasico, false},
		{"INVALIDO", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := entity.ParseSistemaElectrico(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/entity/ -v -run TestSistemaElectrico
```

Expected: FAIL - sistema_electrico.go does not exist

**Step 3: Write minimal implementation**

```go
// internal/domain/entity/sistema_electrico.go
package entity

import "fmt"

type SistemaElectrico string

const (
	SistemaElectricoDelta       SistemaElectrico = "DELTA"
	SistemaElectricoEstrella    SistemaElectrico = "ESTRELLA"
	SistemaElectricoBifasico    SistemaElectrico = "BIFASICO"
	SistemaElectricoMonofasico  SistemaElectrico = "MONOFASICO"
)

var ErrSistemaElectricoInvalido = fmt.Errorf("sistema eléctrico no válido")

func ParseSistemaElectrico(s string) (SistemaElectrico, error) {
	switch s {
	case string(SistemaElectricoDelta):
		return SistemaElectricoDelta, nil
	case string(SistemaElectricoEstrella):
		return SistemaElectricoEstrella, nil
	case string(SistemaElectricoBifasico):
		return SistemaElectricoBifasico, nil
	case string(SistemaElectricoMonofasico):
		return SistemaElectricoMonofasico, nil
	default:
		return "", fmt.Errorf("%w: '%s'", ErrSistemaElectricoInvalido, s)
	}
}

func (s SistemaElectrico) CantidadConductores() int {
	switch s {
	case SistemaElectricoDelta:
		return 3
	case SistemaElectricoEstrella:
		return 4
	case SistemaElectricoBifasico:
		return 3
	case SistemaElectricoMonofasico:
		return 2
	default:
		return 0
	}
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/entity/ -v -run TestSistemaElectrico
```

Expected: All PASS

**Step 5: Commit**

```bash
git add internal/domain/entity/sistema_electrico.go internal/domain/entity/sistema_electrico_test.go
git commit -m "feat(domain): add SistemaElectrico enum with conductor count mapping"
```

---

## Task 6: Service - CalcularFactorTemperatura

**Files:**

- Create: `internal/domain/service/calcular_factor_temperatura.go`
- Create: `internal/domain/service/calcular_factor_temperatura_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/service/calcular_factor_temperatura_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestCalcularFactorTemperatura(t *testing.T) {
	tests := []struct {
		name            string
		tempAmbiente    int
		tempConductor   valueobject.Temperatura
		expectedFactor  float64
		wantErr         bool
	}{
		{"21°C + 75C conductor", 21, valueobject.Temp75, 1.07, false},
		{"26°C + 75C conductor", 26, valueobject.Temp75, 1.03, false},
		{"36°C + 75C conductor", 36, valueobject.Temp75, 0.95, false},
		{"21°C + 60C conductor", 21, valueobject.Temp60, 1.10, false},
		{"31°C exacto", 31, valueobject.Temp75, 1.00, false},
		{"Temperatura negativa", -5, valueobject.Temp75, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factor, err := service.CalcularFactorTemperatura(tt.tempAmbiente, tt.tempConductor)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.InDelta(t, tt.expectedFactor, factor, 0.001)
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/service/ -v -run TestCalcularFactorTemperatura
```

Expected: FAIL - function does not exist

**Step 3: Write implementation**

```go
// internal/domain/service/calcular_factor_temperatura.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// EntradaTablaFactorTemperatura representa una fila de la tabla NOM 310-15(b)(2)(a)
type EntradaTablaFactorTemperatura struct {
	RangoTempC string
	Factor60C  float64
	Factor75C  float64
	Factor90C  float64
}

// CalcularFactorTemperatura retorna el factor de corrección según temperatura ambiente y del conductor
func CalcularFactorTemperatura(
	tempAmbiente int,
	tempConductor valueobject.Temperatura,
	tabla []EntradaTablaFactorTemperatura,
) (float64, error) {
	if tempAmbiente < -10 {
		return 0, fmt.Errorf("temperatura ambiente inválida: %d°C", tempAmbiente)
	}

	for _, entrada := range tabla {
		if rangoContiene(entrada.RangoTempC, tempAmbiente) {
			switch tempConductor {
			case valueobject.Temp60:
				return entrada.Factor60C, nil
			case valueobject.Temp75:
				return entrada.Factor75C, nil
			case valueobject.Temp90:
				return entrada.Factor90C, nil
			default:
				return 0, fmt.Errorf("temperatura de conductor no soportada: %v", tempConductor)
			}
		}
	}

	return 0, fmt.Errorf("no se encontró factor para temperatura ambiente %d°C", tempAmbiente)
}

func rangoContiene(rango string, temp int) bool {
	var min, max int
	if _, err := fmt.Sscanf(rango, "%d-%d", &min, &max); err == nil {
		return temp >= min && temp <= max
	}
	if _, err := fmt.Sscanf(rango, "%d+", &min); err == nil {
		return temp >= min
	}
	return false
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/service/ -v -run TestCalcularFactorTemperatura
```

Expected: PASS (con tabla de test)

**Step 5: Commit**

```bash
git add internal/domain/service/calcular_factor_temperatura.go internal/domain/service/calcular_factor_temperatura_test.go
git commit -m "feat(domain): add CalcularFactorTemperatura service"
```

---

## Task 7: Service - CalcularFactorAgrupamiento

**Files:**

- Create: `internal/domain/service/calcular_factor_agrupamiento.go`
- Create: `internal/domain/service/calcular_factor_agrupamiento_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/service/calcular_factor_agrupamiento_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/stretchr/testify/assert"
)

func TestCalcularFactorAgrupamiento(t *testing.T) {
	tests := []struct {
		cantidad       int
		expectedFactor float64
		wantErr        bool
	}{
		{1, 1.00, false},
		{2, 0.80, false},
		{3, 0.70, false},
		{4, 0.65, false},
		{5, 0.60, false},
		{6, 0.60, false},
		{7, 0.50, false},
		{10, 0.45, false},
		{21, 0.40, false},
		{31, 0.35, false},
		{41, 0.30, false},
		{50, 0.30, false},
		{0, 0, true},
		{-1, 0, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_conductores", tt.cantidad), func(t *testing.T) {
			factor, err := service.CalcularFactorAgrupamiento(tt.cantidad)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.InDelta(t, tt.expectedFactor, factor, 0.001)
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./internal/domain/service/ -v -run TestCalcularFactorAgrupamiento
```

Expected: FAIL

**Step 3: Write implementation**

```go
// internal/domain/service/calcular_factor_agrupamiento.go
package service

import (
	"errors"
	"fmt"
)

var ErrCantidadConductoresInvalida = errors.New("cantidad de conductores debe ser mayor que cero")

// EntradaTablaFactorAgrupamiento representa una fila de la tabla NOM 310-15(b)(3)(a)
type EntradaTablaFactorAgrupamiento struct {
	CantidadMin int
	CantidadMax int // -1 significa "o más"
	Factor      float64
}

// CalcularFactorAgrupamiento retorna el factor de ajuste por cantidad de conductores
func CalcularFactorAgrupamiento(
	cantidad int,
	tabla []EntradaTablaFactorAgrupamiento,
) (float64, error) {
	if cantidad <= 0 {
		return 0, fmt.Errorf("%w: %d", ErrCantidadConductoresInvalida, cantidad)
	}

	for _, entrada := range tabla {
		if entrada.CantidadMax == -1 {
			if cantidad >= entrada.CantidadMin {
				return entrada.Factor, nil
			}
		} else {
			if cantidad >= entrada.CantidadMin && cantidad <= entrada.CantidadMax {
				return entrada.Factor, nil
			}
		}
	}

	// Si no encontramos, usar el factor más bajo (41+)
	return 0.30, nil
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./internal/domain/service/ -v -run TestCalcularFactorAgrupamiento
```

Expected: PASS

**Step 5: Commit**

```bash
git add internal/domain/service/calcular_factor_agrupamiento.go internal/domain/service/calcular_factor_agrupamiento_test.go
git commit -m "feat(domain): add CalcularFactorAgrupamiento service"
```

---

## Task 8: Service - CalcularCharolaEspaciado

**Files:**

- Create: `internal/domain/service/calcular_charola_espaciado.go`
- Create: `internal/domain/service/calcular_charola_espaciado_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/service/calcular_charola_espaciado_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularCharolaEspaciado(t *testing.T) {
	// Tabla de test con tamaños de charola
	tablaCharola := []valueobject.EntradaTablaCanalizacion{
		{Tamano: "6", AreaInteriorMM2: 152.4},
		{Tamano: "9", AreaInteriorMM2: 228.6},
		{Tamano: "12", AreaInteriorMM2: 304.8},
		{Tamano: "16", AreaInteriorMM2: 406.4},
		{Tamano: "18", AreaInteriorMM2: 457.2},
		{Tamano: "20", AreaInteriorMM2: 508.0},
	}

	t.Run("3 hilos Delta - conductor 4 AWG (25.48mm) + tierra 8 AWG (8.5mm)", func(t *testing.T) {
		// Formula: [(3 + 0) - 1] * 25.48 + 8.5 = 2 * 25.48 + 8.5 = 59.46mm
		// Requiere charola de 9" (228.6mm)

		conductorFase := valueobject.Conductor{ /* diametro 25.48mm */ }
		conductorTierra := valueobject.Conductor{ /* diametro 8.5mm */ }

		result, err := service.CalcularCharolaEspaciado(
			1, // hilos por fase
			entity.SistemaElectricoDelta,
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		assert.Equal(t, "9", result.Tamano)
	})
}
```

**Step 2: Run test to verify it fails**

**Step 3: Write implementation**

```go
// internal/domain/service/calcular_charola_espaciado.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// ConductorConDiametro extiende conductor con diámetro para cálculo de charola
type ConductorConDiametro struct {
	Conductor  valueobject.Conductor
	DiametroMM float64
}

// CalcularCharolaEspaciado calcula el ancho requerido para charola con cables espaciados
// Formula NOM: [(hilos_fase_total + hilos_neutro) - 1] * diametro_mayor + diametro_tierra
func CalcularCharolaEspaciado(
	hilosPorFase int,
	sistema entity.SistemaElectrico,
	conductorFase ConductorConDiametro,
	conductorTierra ConductorConDiametro,
	tablaCharola []valueobject.EntradaTablaCanalizacion,
) (entity.Canalizacion, error) {
	// Calcular total de hilos de fase
	hilosFaseTotal := hilosPorFase * 3 // Siempre 3 fases

	// Calcular hilos de neutro según sistema
	hilosNeutro := 0
	if sistema == entity.SistemaElectricoEstrella ||
	   sistema == entity.SistemaElectricoBifasico ||
	   sistema == entity.SistemaElectricoMonofasico {
		hilosNeutro = 1
	}

	// Fórmula: [(total_hilos) - 1] * diametro_fase + diametro_tierra
	totalHilos := hilosFaseTotal + hilosNeutro
	anchoRequerido := float64(totalHilos-1) * conductorFase.DiametroMM + conductorTierra.DiametroMM

	// Seleccionar charola
	for _, entrada := range tablaCharola {
		// Usamos AreaInteriorMM2 como ancho para simplificar
		if entrada.AreaInteriorMM2 >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:      "CHAROLA",
				Tamano:    entrada.Tamano,
				AreaTotal: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf(
		"no se encontró charola suficiente: ancho requerido %.2f mm",
		anchoRequerido,
	)
}
```

**Step 4: Run test to verify it passes**

**Step 5: Commit**

```bash
git add internal/domain/service/calcular_charola_espaciado.go internal/domain/service/calcular_charola_espaciado_test.go
git commit -m "feat(domain): add CalcularCharolaEspaciado service"
```

---

## Task 9: Service - CalcularCharolaTriangular

**Files:**

- Create: `internal/domain/service/calcular_charola_triangular.go`
- Create: `internal/domain/service/calcular_charola_triangular_test.go`

**Step 1: Write the failing test**

```go
// internal/domain/service/calcular_charola_triangular_test.go
package service_test

import (
	"testing"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/service"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcularCharolaTriangular(t *testing.T) {
	tablaCharola := []valueobject.EntradaTablaCanalizacion{
		{Tamano: "6", AreaInteriorMM2: 152.4},
		{Tamano: "9", AreaInteriorMM2: 228.6},
		{Tamano: "12", AreaInteriorMM2: 304.8},
	}

	t.Run("2 hilos por fase - conductor 500 KCM (25.48mm) + tierra 2 AWG (7.42mm)", func(t *testing.T) {
		// Formula: [(2 - 1) * 2.15 * 25.48] + 7.42 = 54.78 + 7.42 = 62.2mm

		conductorFase := service.ConductorConDiametro{
			DiametroMM: 25.48,
		}
		conductorTierra := service.ConductorConDiametro{
			DiametroMM: 7.42,
		}

		result, err := service.CalcularCharolaTriangular(
			2, // hilos por fase
			conductorFase,
			conductorTierra,
			tablaCharola,
		)

		require.NoError(t, err)
		assert.Equal(t, "9", result.Tamano)
	})
}
```

**Step 2: Run test to verify it fails**

**Step 3: Write implementation**

```go
// internal/domain/service/calcular_charola_triangular.go
package service

import (
	"fmt"

	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/domain/valueobject"
)

// CalcularCharolaTriangular calcula el ancho requerido para charola con arreglo triangular
// Formula NOM: [(hilos_por_fase - 1) * 2.15 * diametro_fase] + diametro_tierra
func CalcularCharolaTriangular(
	hilosPorFase int,
	conductorFase ConductorConDiametro,
	conductorTierra ConductorConDiametro,
	tablaCharola []valueobject.EntradaTablaCanalizacion,
) (entity.Canalizacion, error) {
	if hilosPorFase < 1 {
		return entity.Canalizacion{}, fmt.Errorf("hilos por fase debe ser >= 1: %d", hilosPorFase)
	}

	const factorTriangular = 2.15

	// Fórmula: [(hilos_por_fase - 1) * 2.15 * diametro_fase] + diametro_tierra
	anchoRequerido := (float64(hilosPorFase-1) * factorTriangular * conductorFase.DiametroMM) +
		conductorTierra.DiametroMM

	// Seleccionar charola
	for _, entrada := range tablaCharola {
		if entrada.AreaInteriorMM2 >= anchoRequerido {
			return entity.Canalizacion{
				Tipo:      "CHAROLA_TRIANGULAR",
				Tamano:    entrada.Tamano,
				AreaTotal: anchoRequerido,
			}, nil
		}
	}

	return entity.Canalizacion{}, fmt.Errorf(
		"no se encontró charola triangular suficiente: ancho requerido %.2f mm",
		anchoRequerido,
	)
}
```

**Step 4: Run test to verify it passes**

**Step 5: Commit**

```bash
git add internal/domain/service/calcular_charola_triangular.go internal/domain/service/calcular_charola_triangular_test.go
git commit -m "feat(domain): add CalcularCharolaTriangular service"
```

---

## Task 10: Actualizar CalcularTuberia con Fill Factors Correctos

**Files:**

- Modify: `internal/domain/service/calcular_tuberia.go`
- Modify: `internal/domain/service/calcular_tuberia_test.go`

**Step 1: Review current implementation**

Actualmente usa fill factor fijo 0.40. Necesita usar:

- 0.53 (1 conductor)
- 0.31 (2 conductores)
- 0.40 (3+ conductores)

**Step 2: Update implementation**

```go
// CalcularTuberia selects conduit with correct NOM fill factors
func CalcularTuberia(
	conductores []ConductorParaTuberia,
	tipo string,
	tabla []valueobject.EntradaTablaCanalizacion,
) (entity.Canalizacion, error) {
	// ... validaciones existentes ...

	// Calcular área total
	var areaTotal float64
	totalConductores := 0
	for _, c := range conductores {
		areaTotal += float64(c.Cantidad) * c.AreaMM2
		totalConductores += c.Cantidad
	}

	// Determinar fill factor según NOM
	fillFactor := determinarFillFactor(totalConductores)

	areaRequerida := areaTotal / fillFactor

	// ... resto de la implementación ...
}

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
```

**Step 3: Update tests**

Agregar tests para los 3 casos de fill factor.

**Step 4: Commit**

```bash
git add internal/domain/service/calcular_tuberia.go internal/domain/service/calcular_tuberia_test.go
git commit -m "refactor(domain): update CalcularTuberia with correct NOM fill factors 53%/31%/40%"
```

---

## Task 11: Actualizar DTOs - EquipoInput

**Files:**

- Modify: `internal/application/dto/equipo_input.go`

**Step 1: Update struct definition**

```go
type EquipoInput struct {
	// ... campos existentes ...

	// NUEVO: Reemplaza factor_agrupamiento y factor_temperatura
	Estado           string                   `json:"estado" binding:"required"`
	SistemaElectrico entity.SistemaElectrico  `json:"sistema_electrico" binding:"required"`

	// ELIMINAR:
	// FactorAgrupamiento float64
	// FactorTemperatura float64
}
```

**Step 2: Update Validate() method**

```go
func (e *EquipoInput) Validate() error {
	// ... validaciones existentes ...

	// Validar estado no vacío
	if strings.TrimSpace(e.Estado) == "" {
		return fmt.Errorf("estado es requerido")
	}

	// Validar sistema eléctrico
	if _, err := entity.ParseSistemaElectrico(string(e.SistemaElectrico)); err != nil {
		return fmt.Errorf("sistema_electrico inválido: %w", err)
	}

	return nil
}
```

**Step 3: Commit**

```bash
git add internal/application/dto/equipo_input.go
git commit -m "refactor(dto): update EquipoInput with estado and sistema_electrico"
```

---

## Task 12: Actualizar DTOs - MemoriaOutput

**Files:**

- Modify: `internal/application/dto/memoria_output.go`

**Step 1: Add new fields**

```go
type MemoriaOutput struct {
	// ... campos existentes ...

	// NUEVOS: Información de cálculo
	Estado                  string                   `json:"estado"`
	TemperaturaAmbiente     int                      `json:"temperatura_ambiente"`
	SistemaElectrico        entity.SistemaElectrico  `json:"sistema_electrico"`
	CantidadConductores     int                      `json:"cantidad_conductores"`

	// Factores calculados
	FactorTemperaturaCalculado  float64 `json:"factor_temperatura_calculado"`
	FactorAgrupamientoCalculado float64 `json:"factor_agrupamiento_calculado"`
}
```

**Step 2: Commit**

```bash
git add internal/application/dto/memoria_output.go
git commit -m "refactor(dto): add calculated factor fields to MemoriaOutput"
```

---

## Task 13: Actualizar Port - TablaNOMRepository

**Files:**

- Modify: `internal/application/port/tabla_nom_repository.go`

**Step 1: Add new methods to interface**

```go
type TablaNOMRepository interface {
	// ... métodos existentes ...

	// NUEVOS: Tablas de factores
	ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error)
	ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error)
	ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error)

	// NUEVOS: Dimensiones para canalización
	ObtenerDiametroConductor(ctx context.Context, calibre string, material string, conAislamiento bool) (float64, error)
	ObtenerCharolaPorAncho(ctx context.Context, anchoRequeridoMM float64) (valueobject.EntradaTablaCanalizacion, error)
}
```

**Step 2: Commit**

```bash
git add internal/application/port/tabla_nom_repository.go
git commit -m "refactor(port): add new TablaNOMRepository methods for Fase 2"
```

---

## Task 14: Implementar Repository - CSVTablaNOMRepository

**Files:**

- Modify: `internal/infrastructure/repository/csv_tabla_nom_repository.go`
- Modify: `internal/infrastructure/repository/csv_tabla_nom_repository_test.go`

**Step 1: Implement new loader methods**

```go
// loadEstadosTemperatura carga la tabla de estados
func (r *CSVTablaNOMRepository) loadEstadosTemperatura() (map[string]int, error) {
	// Leer CSV y retornar map[estado]temperatura
}

// loadFactoresTemperatura carga la tabla 310-15-b-2-a
func (r *CSVTablaNOMRepository) loadFactoresTemperatura() ([]service.EntradaTablaFactorTemperatura, error) {
	// Leer CSV y retornar slice de entradas
}

// loadFactoresAgrupamiento carga la tabla 310-15-b-3-a
func (r *CSVTablaNOMRepository) loadFactoresAgrupamiento() ([]service.EntradaTablaFactorAgrupamiento, error) {
	// Leer CSV y retornar slice de entradas
}
```

**Step 2: Implement new query methods**

```go
func (r *CSVTablaNOMRepository) ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error) {
	temp, ok := r.estadosTemperatura[strings.ToLower(estado)]
	if !ok {
		return 0, fmt.Errorf("estado no encontrado: %s", estado)
	}
	return temp, nil
}

func (r *CSVTablaNOMRepository) ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error) {
	// Buscar en r.factoresTemperatura
}

func (r *CSVTablaNOMRepository) ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error) {
	// Buscar en r.factoresAgrupamiento
}
```

**Step 3: Update constructor to load new tables**

```go
func NewCSVTablaNOMRepository(basePath string) (*CSVTablaNOMRepository, error) {
	// ... carga existente ...

	// Cargar nuevas tablas
	estados, err := repo.loadEstadosTemperatura()
	if err != nil {
		return nil, fmt.Errorf("failed to load estados temperatura: %w", err)
	}
	repo.estadosTemperatura = estados

	factoresTemp, err := repo.loadFactoresTemperatura()
	if err != nil {
		return nil, fmt.Errorf("failed to load factores temperatura: %w", err)
	}
	repo.factoresTemperatura = factoresTemp

	factoresAgrup, err := repo.loadFactoresAgrupamiento()
	if err != nil {
		return nil, fmt.Errorf("failed to load factores agrupamiento: %w", err)
	}
	repo.factoresAgrupamiento = factoresAgrup

	return repo, nil
}
```

**Step 4: Commit**

```bash
git add internal/infrastructure/repository/csv_tabla_nom_repository.go internal/infrastructure/repository/csv_tabla_nom_repository_test.go
git commit -m "feat(infra): implement new TablaNOMRepository methods for Fase 2 tables"
```

---

## Task 15: Actualizar UseCase - CalcularMemoriaUseCase

**Files:**

- Modify: `internal/application/usecase/calcular_memoria.go`
- Modify: `internal/application/usecase/calcular_memoria_test.go`

**Step 1: Update Execute method - Paso 2**

```go
// PASO 2: Calcular Factores de Ajuste (ACTUALIZADO)
// Obtener temperatura ambiente del estado
tempAmbiente, err := uc.tablaRepo.ObtenerTemperaturaPorEstado(ctx, input.Estado)
if err != nil {
	return dto.MemoriaOutput{}, fmt.Errorf("obtener temperatura para estado %s: %w", input.Estado, err)
}
output.TemperaturaAmbiente = tempAmbiente

// Calcular factor de temperatura
factorTemp, err := uc.tablaRepo.ObtenerFactorTemperatura(ctx, tempAmbiente, temperatura)
if err != nil {
	return dto.MemoriaOutput{}, fmt.Errorf("calcular factor temperatura: %w", err)
}
output.FactorTemperaturaCalculado = factorTemp

// Calcular cantidad de conductores desde sistema eléctrico
cantidadConductores := input.SistemaElectrico.CantidadConductores()
output.CantidadConductores = cantidadConductores

// Calcular factor de agrupamiento
factorAgrup, err := uc.tablaRepo.ObtenerFactorAgrupamiento(ctx, cantidadConductores)
if err != nil {
	return dto.MemoriaOutput{}, fmt.Errorf("calcular factor agrupamiento: %w", err)
}
output.FactorAgrupamientoCalculado = factorAgrup

// Ajustar corriente con factores calculados
factores := map[string]float64{
	"agrupamiento": factorAgrup,
	"temperatura":  factorTemp,
}
```

**Step 2: Update Paso 6 - Canalización**

```go
// PASO 6: Dimensionar Canalización (ACTUALIZADO)
switch input.TipoCanalizacion {
case entity.TipoCanalizacionTuberiaPVC,
	 entity.TipoCanalizacionTuberiaAluminio,
	 entity.TipoCanalizacionTuberiaAceroPG,
	 entity.TipoCanalizacionTuberiaAceroPD:
	// Calcular tubería con áreas reales
	conductoresTuberia := []service.ConductorParaTuberia{
		{Cantidad: hilosPorFase * 3, AreaMM2: areaFase}, // 3 fases
		{Cantidad: hilosNeutro, AreaMM2: areaNeutro},    // Neutro si aplica
		{Cantidad: 1, AreaMM2: areaTierra},              // Tierra
	}
	canalizacion, err = service.CalcularTuberia(conductoresTuberia, tipoCanalizacionStr, tablaConduit)

 case entity.TipoCanalizacionCharolaCableEspaciado:
	// Calcular charola espaciada
	canalizacion, err = service.CalcularCharolaEspaciado(
		hilosPorFase,
		input.SistemaElectrico,
		service.ConductorConDiametro{Conductor: conductor, DiametroMM: diametroFase},
		service.ConductorConDiametro{Conductor: conductorTierra, DiametroMM: diametroTierra},
		tablaCharola,
	)

 case entity.TipoCanalizacionCharolaCableTriangular:
	// Calcular charola triangular
	canalizacion, err = service.CalcularCharolaTriangular(
		hilosPorFase,
		service.ConductorConDiametro{Conductor: conductor, DiametroMM: diametroFase},
		service.ConductorConDiametro{Conductor: conductorTierra, DiametroMM: diametroTierra},
		tablaCharola,
	)
}
```

**Step 3: Update tests**

Actualizar tests para usar `Estado` y `SistemaElectrico` en lugar de factores manuales.

**Step 4: Commit**

```bash
git add internal/application/usecase/calcular_memoria.go internal/application/usecase/calcular_memoria_test.go
git commit -m "refactor(usecase): update CalcularMemoria with automatic factor calculation"
```

---

## Task 16: Actualizar Handler - CalculoHandler

**Files:**

- Modify: `internal/presentation/handler/calculo_handler.go`
- Modify: `internal/presentation/handler/calculo_handler_test.go`

**Step 1: Update request binding**

```go
type CalcularMemoriaRequest struct {
	// ... campos existentes ...

	// NUEVOS (requeridos)
	Estado           string `json:"estado" binding:"required"`
	SistemaElectrico string `json:"sistema_electrico" binding:"required,oneof=DELTA ESTRELLA BIFASICO MONOFASICO"`

	// ELIMINADOS:
	// FactorAgrupamiento float64 `json:"factor_agrupamiento,omitempty"`
	// FactorTemperatura  float64 `json:"factor_temperatura,omitempty"`
}
```

**Step 2: Update mapping to DTO**

```go
equipoInput := dto.EquipoInput{
	// ... mapeo existente ...
	Estado:           req.Estado,
	SistemaElectrico: entity.SistemaElectrico(req.SistemaElectrico),
}
```

**Step 3: Commit**

```bash
git add internal/presentation/handler/calculo_handler.go internal/presentation/handler/calculo_handler_test.go
git commit -m "refactor(handler): update CalculoHandler with new estado and sistema_electrico fields"
```

---

## Task 17: Tests de Integración

**Files:**

- Create: `tests/integration/fase2_calculo_test.go`

**Step 1: Create integration test**

```go
// tests/integration/fase2_calculo_test.go
package integration

import (
	"context"
	"testing"

	"github.com/garfex/calculadora-filtros/internal/application/dto"
	"github.com/garfex/calculadora-filtros/internal/application/usecase"
	"github.com/garfex/calculadora-filtros/internal/domain/entity"
	"github.com/garfex/calculadora-filtros/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFase2_CalculoCompleto(t *testing.T) {
	// Setup
	tablaRepo, err := repository.NewCSVTablaNOMRepository("../../data/tablas_nom")
	require.NoError(t, err)

	// Use in-memory equipo repo for test
	equipoRepo := repository.NewInMemoryEquipoRepository()

	uc := usecase.NewCalcularMemoriaUseCase(tablaRepo, equipoRepo)

	ctx := context.Background()
	input := dto.EquipoInput{
		Modo:             dto.ModoManualAmperaje,
		TipoEquipo:       entity.TipoEquipoFiltroActivo,
		Clave:            "FA-TEST-001",
		AmperajeNominal:  100,
		Tension:          480,
		Estado:           "Nuevo Leon",
		SistemaElectrico: entity.SistemaElectricoDelta,
		TipoCanalizacion: entity.TipoCanalizacionTuberiaPVC,
		ITM:              125,
		LongitudCircuito: 50,
		HilosPorFase:     1,
	}

	// Execute
	output, err := uc.Execute(ctx, input)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Nuevo Leon", output.Estado)
	assert.Equal(t, 21, output.TemperaturaAmbiente) // Temperatura promedio de Nuevo León
	assert.Equal(t, entity.SistemaElectricoDelta, output.SistemaElectrico)
	assert.Equal(t, 3, output.CantidadConductores)
	assert.InDelta(t, 1.07, output.FactorTemperaturaCalculado, 0.01) // 21°C + 75C
	assert.InDelta(t, 0.70, output.FactorAgrupamientoCalculado, 0.01) // 3 conductores
	assert.NotEmpty(t, output.Canalizacion.Tamano)
}
```

**Step 2: Run integration test**

```bash
go test ./tests/integration/ -v -run TestFase2_CalculoCompleto
```

Expected: PASS

**Step 3: Commit**

```bash
git add tests/integration/fase2_calculo_test.go
git commit -m "test(integration): add Fase 2 integration test"
```

---

## Task 18: Verificación Final

**Files:**

- All modified files

**Step 1: Run all tests**

```bash
go test ./... -race
```

Expected: All PASS

**Step 2: Build**

```bash
go build ./...
```

Expected: No errors

**Step 3: Lint**

```bash
golangci-lint run
```

Expected: No critical issues

**Step 4: Final commit**

```bash
git commit --allow-empty -m "chore: Fase 2 implementation complete - automatic factors and NOM conduit calculations"
```

---

## QA Checklist

- [ ] All new CSV tables created and validated
- [ ] Domain services implemented with tests (TDD)
- [ ] DTOs updated with new fields
- [ ] Repository implements all new methods
- [ ] UseCase calculates factors automatically
- [ ] Handler accepts new API format
- [ ] Integration tests pass
- [ ] All unit tests pass
- [ ] No breaking changes to existing API (except removed manual factors)
- [ ] Documentation updated

---

## Breaking Changes (Documentar)

1. **API Input:** Se eliminan campos `factor_agrupamiento` y `factor_temperatura`
2. **API Input:** Se agregan campos requeridos `estado` y `sistema_electrico`
3. **API Output:** Se agregan campos con prefijo `_calculado` para los factores

**Migration Guide:**

- Clientes deben actualizar para enviar `estado` (ej: "Nuevo Leon") en lugar de `factor_temperatura`
- Clientes deben actualizar para enviar `sistema_electrico` (ej: "DELTA") en lugar de `factor_agrupamiento`
