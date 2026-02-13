# Fase 2: Memoria de Cálculo Completa - Diseño

> **Fecha:** 2026-02-13  
> **Feature:** Fase 2 - Cálculo automático de factores y canalización NOM  
> **Status:** Diseño aprobado, listo para implementación

---

## Resumen Ejecutivo

Fase 2 mejora el cálculo de memoria de instalaciones eléctricas implementando:

1. **Factores de ajuste automáticos** (temperatura y agrupamiento) basados en tablas NOM
2. **Cálculo correcto de canalización** según tipo:
   - Tubería: fill factor según cantidad de conductores (53%/31%/40%)
   - Charola espaciada: fórmula NOM con espacios
   - Charola triangular: fórmula NOM para arreglos

---

## Problema Actual (Fase 1)

- Factores de temperatura y agrupamiento son **inputs manuales** del usuario
- Cálculo de tubería usa **fill factor fijo 40%** sin considerar cantidad de conductores
- Charolas usan **simplificación** sin fórmulas NOM específicas

---

## Solución Propuesta

### Cambios en Inputs

| Campo Actual                    | Nuevo Campo                | Descripción                                            |
| ------------------------------- | -------------------------- | ------------------------------------------------------ |
| `factor_agrupamiento` (float64) | `estado` (string)          | Estados de México para temperatura promedio            |
| `factor_temperatura` (float64)  | `sistema_electrico` (enum) | Tipo de sistema: DELTA, ESTRELLA, BIFASICO, MONOFASICO |

Los factores se calculan **automáticamente** a partir de estos inputs.

---

## Nuevas Tablas CSV

### 1. estados_temperatura.csv

Temperatura promedio anual por estado de México.

```csv
estado,temperatura_promedio_c
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

### 2. 310-15-b-2-a.csv

Factores de corrección por temperatura ambiente (NOM-001-SEDE Tabla 310-15(b)(2)(a)).

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

### 3. 310-15-b-3-a.csv

Factores de ajuste por cantidad de conductores (NOM-001-SEDE Tabla 310-15(b)(3)(a)).

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

### 4. charola_dimensiones.csv

Tamaños comerciales de charola portacables.

```csv
tamano_pulgadas,ancho_mm,tipo
6,152.4,comercial
9,228.6,comercial
12,304.8,comercial
16,406.4,comercial
18,457.2,comercial
20,508.0,comercial
24,609.6,comercial
30,762.0,comercial
36,914.4,comercial
```

### 5. Tablas existentes a usar

- `tabla-5-dimensiones-aislamiento.csv` → Columna `area_total_mm2`
- `tabla-8-conductor-desnudo.csv` → Columna `area_total_mm2`
- `tabla-conduit-dimensiones.csv` → Columna `area_interior_mm2`

---

## Nuevos Domain Services

### 1. CalcularFactorTemperatura

```go
// Input: temperatura ambiente, temperatura del conductor (60C/75C/90C)
// Output: factor de corrección
// Tabla: 310-15-b-2-a.csv
func CalcularFactorTemperatura(tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error)
```

**Ejemplo:**

- Input: tempAmbiente=21 (Nuevo León), tempConductor=75C
- Busca rango 21-25 → factor 1.07
- Output: 1.07

### 2. CalcularFactorAgrupamiento

```go
// Input: cantidad de conductores portadores de corriente
// Output: factor de ajuste
// Tabla: 310-15-b-3-a.csv
func CalcularFactorAgrupamiento(cantidadConductores int) (float64, error)
```

**Mapeo Sistema Eléctrico → Cantidad de Conductores:**

| Sistema            | Fases | Neutro | Total Conductores Portadores |
| ------------------ | ----- | ------ | ---------------------------- |
| Trifásico Delta    | 3     | 0      | 3                            |
| Trifásico Estrella | 3     | 1      | 4                            |
| Bifásico           | 2     | 1      | 3                            |
| Monofásico         | 1     | 1      | 2                            |

**Ejemplo:**

- Input: Sistema=ESTRELLA → 4 conductores
- Tabla: 4 conductores → factor 0.65
- Output: 0.65

### 3. CalcularTuberia (mejorado)

```go
// Input: lista de conductores con sus áreas
// Output: canalización seleccionada
// Fill factor: 53% (1 conductor), 31% (2 conductores), 40% (3+ conductores)
func CalcularTuberia(
    conductores []ConductorParaTuberia,
    tabla []valueobject.EntradaTablaCanalizacion,
) (entity.Canalizacion, error)

type ConductorParaTuberia struct {
    Cantidad   int
    AreaMM2    float64  // De Tabla 5 o Tabla 8
    EsTierra   bool
}
```

**Fórmula:**

```
areaTotal = Σ (cantidad × areaMM2) para todos los conductores
fillFactor = 0.53 (si totalConductores == 1)
           = 0.31 (si totalConductores == 2)
           = 0.40 (si totalConductores >= 3)
areaRequerida = areaTotal / fillFactor
seleccionar conduit donde areaInterior >= areaRequerida
```

### 4. CalcularCharolaEspaciado

```go
// Fórmula NOM: [(hilos_fase + hilos_neutro) - 1] × diametro_mayor + diametro_tierra
func CalcularCharolaEspaciado(
    hilosPorFase int,
    sistemaElectrico SistemaElectrico,
    conductorFase valueobject.Conductor,
    conductorTierra valueobject.Conductor,
    tablaCharola []valueobject.EntradaTablaCharola,
) (entity.Canalizacion, error)
```

**Fórmula:**

```
diamentroFase = conductorFase.DiametroMM()  // De Tabla 5
diametroTierra = conductorTierra.DiametroMM()  // De Tabla 8

hilosNeutro = 0 (si Delta) o 1 (si Estrella/Bifasico/Monofasico)
totalHilos = (hilosPorFase * 3) + hilosNeutro  // Para trifásico

anchoRequerido = ((totalHilos) - 1) * diametroFase + diametroTierra

seleccionar charola donde ancho >= anchoRequerido
```

### 5. CalcularCharolaTriangular

```go
// Fórmula NOM: [(hilos_por_fase - 1) × 2.15 × diametro_mayor] + diametro_tierra
func CalcularCharolaTriangular(
    hilosPorFase int,
    conductorFase valueobject.Conductor,
    conductorTierra valueobject.Conductor,
    tablaCharola []valueobject.EntradaTablaCharola,
) (entity.Canalizacion, error)
```

**Fórmula:**

```
diamentroFase = conductorFase.DiametroMM()  // De Tabla 5
diametroTierra = conductorTierra.DiametroMM()  // De Tabla 8

// Arreglo triangular: los cables se apilan en forma triangular
anchoRequerido = ((hilosPorFase - 1) * 2.15 * diametroFase) + diametroTierra

seleccionar charola donde ancho >= anchoRequerido
```

---

## Cambios en Value Objects

### Conductor

Agregar método para obtener diámetro:

```go
func (c Conductor) DiametroMM() float64 {
    // Calcular diámetro a partir del área: d = 2 * sqrt(area / π)
    return 2 * math.Sqrt(c.areaConAislamientoMM2 / math.Pi)
}
```

Nota: Para conductores desnudos (tierra), usar el área de la sección sin aislamiento.

---

## Cambios en DTOs

### EquipoInput

```go
type EquipoInput struct {
    // ... campos existentes ...

    // NUEVO: Reemplaza factor_agrupamiento
    Estado string `json:"estado" binding:"required"`

    // NUEVO: Reemplaza factor_temperatura
    SistemaElectrico SistemaElectrico `json:"sistema_electrico" binding:"required"`

    // Eliminados (ahora se calculan):
    // FactorAgrupamiento float64
    // FactorTemperatura float64
}

type SistemaElectrico string

const (
    SistemaElectricoDelta    SistemaElectrico = "DELTA"
    SistemaElectricoEstrella SistemaElectrico = "ESTRELLA"
    SistemaElectricoBifasico SistemaElectrico = "BIFASICO"
    SistemaElectricoMonofasico SistemaElectrico = "MONOFASICO"
)
```

### MemoriaOutput

```go
type MemoriaOutput struct {
    // ... campos existentes ...

    // NUEVOS: Información de cálculo
    Estado string `json:"estado"`
    TemperaturaAmbiente int `json:"temperatura_ambiente"`
    SistemaElectrico SistemaElectrico `json:"sistema_electrico"`
    CantidadConductores int `json:"cantidad_conductores"`

    // Calculados automáticamente
    FactorTemperatura float64 `json:"factor_temperatura"`
    FactorAgrupamiento float64 `json:"factor_agrupamiento"`
}
```

---

## Cambios en Ports (Repository)

Nuevos métodos en `TablaNOMRepository`:

```go
type TablaNOMRepository interface {
    // ... métodos existentes ...

    // NUEVO: Temperatura promedio por estado
    ObtenerTemperaturaPorEstado(ctx context.Context, estado string) (int, error)

    // NUEVO: Factor de corrección por temperatura
    ObtenerFactorTemperatura(ctx context.Context, tempAmbiente int, tempConductor valueobject.Temperatura) (float64, error)

    // NUEVO: Factor de ajuste por agrupamiento
    ObtenerFactorAgrupamiento(ctx context.Context, cantidadConductores int) (float64, error)

    // NUEVO: Área del conductor con aislamiento (Tabla 5)
    ObtenerAreaConductorConAislamiento(ctx context.Context, calibre string, aislamiento string) (float64, error)

    // NUEVO: Área del conductor desnudo (Tabla 8)
    ObtenerAreaConductorDesnudo(ctx context.Context, calibre string) (float64, error)

    // NUEVO: Diámetro del conductor con aislamiento
    ObtenerDiametroConductorConAislamiento(ctx context.Context, calibre string, aislamiento string) (float64, error)

    // NUEVO: Diámetro del conductor desnudo
    ObtenerDiametroConductorDesnudo(ctx context.Context, calibre string) (float64, error)

    // NUEVO: Seleccionar charola por ancho requerido
    ObtenerCharolaPorAncho(ctx context.Context, anchoRequeridoMM float64) (valueobject.EntradaTablaCharola, error)
}
```

---

## Flujo del UseCase Actualizado

```
PASO 1: Calcular Corriente Nominal
└── Igual que Fase 1

PASO 2: Calcular Factores de Ajuste (ACTUALIZADO)
├── 2.1: Obtener temperatura promedio del estado (estados_temperatura.csv)
├── 2.2: Calcular factor temperatura (310-15-b-2-a.csv)
├── 2.3: Determinar cantidad de conductores desde sistema_electrico
├── 2.4: Calcular factor agrupamiento (310-15-b-3-a.csv)
└── 2.5: Ajustar corriente nominal

PASO 3: Seleccionar Canalización
└── Igual que Fase 1 (input del usuario)

PASO 4: Seleccionar Conductor Alimentación
└── Igual que Fase 1

PASO 5: Seleccionar Conductor Tierra
└── Igual que Fase 1

PASO 6: Dimensionar Canalización (ACTUALIZADO)
├── 6.1: Determinar tipo de canalización
├── 6.2: Si TUBERIA → CalcularTuberia
│   ├── Obtener áreas de Tabla 5 (fases) y Tabla 8 (tierra)
│   ├── Calcular fill factor según cantidad de conductores
│   └── Seleccionar conduit
├── 6.3: Si CHAROLA ESPACIADO → CalcularCharolaEspaciado
│   ├── Obtener diámetros de Tabla 5 y Tabla 8
│   └── Calcular ancho con fórmula de espacios
└── 6.4: Si CHAROLA TRIANGULAR → CalcularCharolaTriangular
    ├── Obtener diámetros de Tabla 5 y Tabla 8
    └── Calcular ancho con fórmula triangular

PASO 7: Calcular Caída de Tensión
└── Igual que Fase 1
```

---

## Consideraciones de Equipos

Cada tipo de equipo define su **sistema eléctrico por defecto**:

| Equipo        | Sistema por Defecto | ¿Configurable? |
| ------------- | ------------------- | -------------- |
| FiltroActivo  | Trifásico Delta     | No             |
| FiltroRechazo | Trifásico Delta     | No             |
| Transformador | Trifásico Estrella  | No             |
| Carga         | Trifásico Estrella  | Sí             |

Para **Carga**, el usuario puede especificar `SistemaElectrico` en el input. Si no especifica, usa Estrella por defecto.

---

## Criterios de Aceptación

- [ ] Tabla `estados_temperatura.csv` con 32 estados de México
- [ ] Tabla `310-15-b-2-a.csv` con factores por temperatura
- [ ] Tabla `310-15-b-3-a.csv` con factores por agrupamiento
- [ ] Tabla `charola_dimensiones.csv` con tamaños comerciales
- [ ] Service `CalcularFactorTemperatura` con tests
- [ ] Service `CalcularFactorAgrupamiento` con tests
- [ ] Service `CalcularTuberia` mejorado con fill factors 53%/31%/40%
- [ ] Service `CalcularCharolaEspaciado` con fórmula NOM
- [ ] Service `CalcularCharolaTriangular` con fórmula NOM
- [ ] Repository implementa nuevos métodos de lectura de tablas
- [ ] UseCase actualizado para calcular factores automáticamente
- [ ] API acepta `estado` y `sistema_electrico` en lugar de factores manuales
- [ ] Tests de integración pasan con casos reales

---

## Notas de Implementación

1. **Temperatura del conductor:** Usar la misma lógica de Fase 1:
   - ≤100A → 60°C
   - > 100A → 75°C
   - Override a 90°C si el usuario lo especifica

2. **Mapeo de temperatura ambiente a rango:** Usar el rango que contenga la temperatura exacta. Ej: 21°C → rango 21-25.

3. **Cantidad de conductores:** Incluir neutro como conductor portador solo si existe (Estrella, Bifasico, Monofasico). Delta no tiene neutro.

4. **Diámetros:** Almacenar diámetros en las tablas CSV, no calcularlos del área. La NOM ya da los diámetros exactos.

5. **Error handling:** Si un estado no existe en la tabla, retornar error específico. Lo mismo para rangos de temperatura no encontrados.

---

## Archivos a Modificar/Crear

### Domain Services (Nuevos)

- `internal/domain/service/calcular_factor_temperatura.go`
- `internal/domain/service/calcular_factor_agrupamiento.go`
- `internal/domain/service/calcular_charola_espaciado.go`
- `internal/domain/service/calcular_charola_triangular.go`
- Actualizar: `internal/domain/service/calcular_tuberia.go`

### Domain Services (Tests)

- `internal/domain/service/calcular_factor_temperatura_test.go`
- `internal/domain/service/calcular_factor_agrupamiento_test.go`
- `internal/domain/service/calcular_charola_espaciado_test.go`
- `internal/domain/service/calcular_charola_triangular_test.go`
- Actualizar: `internal/domain/service/calcular_tuberia_test.go`

### Application (DTOs)

- Actualizar: `internal/application/dto/equipo_input.go`
- Actualizar: `internal/application/dto/memoria_output.go`

### Application (Ports)

- Actualizar: `internal/application/port/tabla_nom_repository.go`

### Application (UseCase)

- Actualizar: `internal/application/usecase/calcular_memoria.go`
- Actualizar: `internal/application/usecase/calcular_memoria_test.go`

### Infrastructure

- Actualizar: `internal/infrastructure/repository/csv_tabla_nom_repository.go`

### Data (Tablas CSV)

- Crear: `data/tablas_nom/estados_temperatura.csv`
- Crear: `data/tablas_nom/310-15-b-2-a.csv`
- Crear: `data/tablas_nom/310-15-b-3-a.csv`
- Crear: `data/tablas_nom/charola_dimensiones.csv`
- Actualizar: `data/tablas_nom/tabla-5-dimensiones-aislamiento.csv` (agregar columna diametro_mm si no existe)
- Actualizar: `data/tablas_nom/tabla-8-conductor-desnudo.csv` (agregar columna diametro_mm si no existe)

---

## Referencias

- NOM-001-SEDE-2012: Instalaciones Eléctricas (Utilización)
- Tabla 310-15(b)(2)(a): Factores de corrección por temperatura
- Tabla 310-15(b)(3)(a): Factores de ajuste por agrupamiento
- Tabla 5: Dimensiones de conductores aislados
- Tabla 8: Dimensiones de conductores desnudos
- IEEE-141: Cálculo de caída de tensión
