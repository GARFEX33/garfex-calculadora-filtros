# Diseño: Soporte de Material (Cu/Al) en Conductor de Tierra y Limpieza de Calibres

**Fecha:** 2026-02-13
**Estado:** Aprobado

## Contexto

El sistema tenía tres problemas relacionados:

1. `calibresValidos` en `conductor.go` incluía calibres fuera del rango de uso (18 AWG, 16 AWG, 3 AWG, 1 AWG, 700-2000 MCM) y faltaban 750 MCM y 1000 MCM.
2. Las tablas CSV de ampacidad (b-16, b-17, b-20) tenían filas para esos calibres inválidos.
3. `250-122.csv` solo tenía columnas para cobre — la NOM define también aluminio — y el material estaba hardcodeado como Cu en el use case.

## Objetivo

- Limpiar calibres a la lista oficial de uso.
- Actualizar `250-122.csv` con estructura Cu+Al según la NOM, cortado en ITM 4000.
- Permitir que el usuario elija material (Cu o Al) en el input.
- Cuando Al no está disponible para un ITM (≤ 100A), hacer fallback silencioso a Cu.

## Cambios por Capa

---

### 1. `data/tablas_nom/250-122.csv`

Nueva estructura con 5 columnas. Al vacío = fallback a Cu.

```csv
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

### 2. Tablas de ampacidad (b-16, b-17, b-20)

Eliminar filas de calibres fuera de la lista:
- `18 AWG`, `16 AWG`, `3 AWG`, `1 AWG`
- `700 MCM`, `800 MCM`, `900 MCM`, `1250 MCM`, `1500 MCM`, `1750 MCM`, `2000 MCM`

---

### 3. `internal/domain/valueobject/conductor.go`

`calibresValidos` queda exactamente:

```go
var calibresValidos = map[string]bool{
    "14 AWG": true, "12 AWG": true, "10 AWG": true, "8 AWG": true,
    "6 AWG":  true, "4 AWG":  true, "2 AWG":  true,
    "1/0 AWG": true, "2/0 AWG": true, "3/0 AWG": true, "4/0 AWG": true,
    "250 MCM": true, "300 MCM": true, "350 MCM": true, "400 MCM": true,
    "500 MCM": true, "600 MCM": true, "750 MCM": true, "1000 MCM": true,
}
```

---

### 4. `internal/domain/valueobject/tabla_entrada.go`

`EntradaTablaTierra` pasa de un solo conductor a Cu + Al opcional:

```go
// Antes:
type EntradaTablaTierra struct {
    ITMHasta  int
    Conductor ConductorParams
}

// Después:
type EntradaTablaTierra struct {
    ITMHasta    int
    ConductorCu ConductorParams  // siempre presente
    ConductorAl *ConductorParams // nil = no disponible para este ITM, usar fallback Cu
}
```

---

### 5. `internal/domain/service/calculo_tierra.go`

Nueva firma:

```go
func SeleccionarConductorTierra(
    itm int,
    material valueobject.MaterialConductor,
    tabla []valueobject.EntradaTablaTierra,
) (valueobject.Conductor, error)
```

Lógica:
```
para cada entrada donde itm <= entrada.ITMHasta:
    si material == Al Y entrada.ConductorAl != nil → NewConductor(*entrada.ConductorAl)
    si material == Al Y entrada.ConductorAl == nil → fallback silencioso → NewConductor(entrada.ConductorCu)
    si material == Cu → NewConductor(entrada.ConductorCu)
```

---

### 6. `internal/application/dto/equipo_input.go`

Agregar campo:

```go
Material valueobject.MaterialConductor // "Cu" o "Al"; default Cu si vacío
```

Sin `binding:"required"`.

---

### 7. `internal/application/usecase/calcular_memoria.go`

**Cambio 1** — reemplazar hardcoded:
```go
// Antes:
material := valueobject.MaterialCobre

// Después:
material := input.Material
if material == "" {
    material = valueobject.MaterialCobre
}
```

**Cambio 2** — pasar material a SeleccionarConductorTierra:
```go
conductorTierra, err := service.SeleccionarConductorTierra(input.ITM, material, tablaTierra)
```

---

### 8. `internal/application/dto/memoria_output.go`

Agregar `Material` al DTO de salida para que el reporte lo refleje:

```go
Material string `json:"material"` // "Cu" o "Al"
```

Y en el mapeo del use case:
```go
output.Material = string(material)
```

---

### 9. Infrastructure — CSV reader de `250-122`

El reader de `ObtenerTablaTierra` debe parsear las 5 columnas nuevas y construir `EntradaTablaTierra` con `ConductorAl = nil` cuando `al_calibre` está vacío.

---

## Tests

### Actualizar
- `calculo_tierra_test.go` — agregar `material` como 2do argumento en todos los casos existentes (pasar `MaterialCobre`)
- Tests del CSV reader de `250-122` en infrastructure

### Tests nuevos en `calculo_tierra_test.go`

| Test | Input | Esperado |
|------|-------|----------|
| `TestSeleccionarConductorTierra_AluminioDisponible` | ITM 200, Al | 4 AWG Al |
| `TestSeleccionarConductorTierra_AluminioFallbackCu` | ITM 60, Al | 10 AWG Cu (fallback) |
| `TestSeleccionarConductorTierra_CobreExplicito` | ITM 200, Cu | 6 AWG Cu |
| `TestSeleccionarConductorTierra_ITMMaximo` | ITM 4000, Al | 750 MCM Al |
| `TestSeleccionarConductorTierra_ITMExcede` | ITM 5000 | `ErrConductorNoEncontrado` |

### Tests nuevos en `conductor_test.go`

- `3 AWG` y `1 AWG` → `ErrConductorInvalido`
- `750 MCM` y `1000 MCM` → válidos

## Calibres Válidos (lista definitiva)

`14, 12, 10, 8, 6, 4, 2, 1/0, 2/0, 3/0, 4/0, 250, 300, 350, 400, 500, 600, 750, 1000`

Solo AWG y MCM. Sin 18, 16, 3, 1 AWG. Sin 700, 800, 900, 1250, 1500, 1750, 2000 MCM.
