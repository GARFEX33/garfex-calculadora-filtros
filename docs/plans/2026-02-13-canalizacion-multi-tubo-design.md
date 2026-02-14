# Diseño: Soporte Multi-Tubo en CalcularCanalizacion

**Fecha:** 2026-02-13  
**Estado:** Aprobado  
**Archivo principal:** `internal/domain/service/calculo_canalizacion.go`

## Contexto

`CalcularCanalizacion` actualmente selecciona siempre 1 tubo. En instalaciones con conductores en paralelo (`HilosPorFase > 1`), la NOM permite —y a veces requiere— distribuir los conductores en múltiples tubos del mismo tamaño.

## Objetivo

Agregar soporte para N tubos iguales en paralelo, manteniendo compatibilidad con el comportamiento actual (1 tubo).

## Cambios de Entidad

### `entity.Canalizacion`

Agregar campo `NumeroDeTubos int`:

```go
type Canalizacion struct {
    Tipo           string
    Tamano         string
    AnchoRequerido float64
    NumeroDeTubos  int  // cantidad de tubos en paralelo; 1 = instalación normal
}
```

## Cambios de Servicio

### Firma

```go
func CalcularCanalizacion(
    conductores []ConductorParaCanalizacion,
    tipo string,
    tabla []valueobject.EntradaTablaCanalizacion,
    numeroDeTubos int,
) (entity.Canalizacion, error)
```

### Lógica de cálculo

```
cantidadTotal      = suma de todos los Cantidad en conductores
conductoresPorTubo = cantidadTotal / numeroDeTubos   // división entera
fillFactor         = determinarFillFactor(conductoresPorTubo)
areaPorTubo        = areaTotal / numeroDeTubos
areaRequerida      = areaPorTubo / fillFactor
→ buscar en tabla el primer Tamano cuyo AreaInteriorMM2 >= areaRequerida
```

**Clave:** el fill factor se determina con `conductoresPorTubo`, no con el total. El NEC/NOM define el fill según cuántos conductores van en ese tubo individual.

## Validaciones

| Condición | Error |
|-----------|-------|
| `numeroDeTubos < 1` | `"numeroDeTubos debe ser mayor a cero"` |
| `conductores` vacío | `"lista de conductores vacía"` (existente) |
| `tabla` vacía | `ErrCanalizacionNoDisponible: tabla vacía` (existente) |
| área por tubo excede tabla | `ErrCanalizacionNoDisponible` (existente) |

## Casos límite

| Caso | Comportamiento |
|------|---------------|
| `numeroDeTubos = 1` | Idéntico al comportamiento actual |
| `numeroDeTubos = 2`, área entra en tubo más pequeño | Retorna tubo chico con `NumeroDeTubos: 2` |
| `numeroDeTubos = 2`, área dividida no entra en ningún tubo | `ErrCanalizacionNoDisponible` |
| `numeroDeTubos > cantidadTotal` | Permitido — responsabilidad del caller |

## Tests

### Tests existentes
Se actualizan agregando `1` como último argumento. Sin cambio de lógica.

### Tests nuevos

| Test | Descripción |
|------|-------------|
| `TestCalcularCanalizacion_DosTubos` | 4 conductores, `numeroDeTubos=2` → tubo más chico que con 1 tubo, `NumeroDeTubos: 2` |
| `TestCalcularCanalizacion_DosTubosSmall` | conductores chicos, `numeroDeTubos=2` → fill factor calculado con `cantidadTotal/2` |
| `TestCalcularCanalizacion_NumeroDeTubosCero` | `numeroDeTubos=0` → error |
| `TestCalcularCanalizacion_NumeroDeTubosNegativo` | `numeroDeTubos=-1` → error |
| `TestCalcularCanalizacion_NumeroDeTubosMayorQueConductores` | 2 conductores, `numeroDeTubos=5` → no error |

## Impacto en callers

Cualquier sitio que llame a `CalcularCanalizacion` necesita pasar `1` como `numeroDeTubos` para mantener comportamiento actual. Buscar con `rg "CalcularCanalizacion"` para identificar todos los callers.
