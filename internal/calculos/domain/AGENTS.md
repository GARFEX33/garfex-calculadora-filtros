# Calculos — Domain Layer

Capa de negocio pura para la feature de cálculos eléctricos. Sin dependencias externas (sin Gin, pgx, CSV).

## Estructura

| Subdirectorio | Contenido                                            |
| ------------- | ---------------------------------------------------- |
| `entity/`     | Entidades, tipos, interfaces del dominio de cálculos |
| `service/`    | Servicios de cálculo puros (sin I/O)                 |

## Dependencias permitidas

- `internal/shared/kernel/valueobject` — value objects compartidos
- stdlib de Go

## Dependencias prohibidas

> Ver reglas consolidadas en [docs/reference/structure.md](../../../docs/reference/structure.md)

## Guías

| Subdirectorio | Contenido                                                                       |
| ------------- | ------------------------------------------------------------------------------- |
| `entity/`     | Entidades: TipoEquipo, TipoCanalizacion, SistemaElectrico, MemoriaCalculo, etc. |
| `service/`    | Servicios de cálculo NOM (ls internal/calculos/domain/service/*.go)            |

> **Nota:** Las subcarpetas `entity/` y `service/` heredan las reglas de este AGENTS.md. No necesitan AGENTS.md propio.

## Convenciones de Cálculo — Caída de Tensión

### Impedancia efectiva (obligatorio)

`EntradaCalculoCaidaTension` requiere `FactorPotencia float64` (cosθ, rango `(0, 1]`).

El service calcula la **impedancia efectiva** según NOM / IEEE-141:

```go
senTheta := math.Sqrt(1 - cosTheta*cosTheta)
Zef      := resistencia*cosTheta + reactancia*senTheta   // Ω/km por conductor
```

El campo `Impedancia` en `entity.ResultadoCaidaTension` representa **Zef por conductor**, no `√(R²+X²)`.

**Nota:** La división por N (número de hilos en paralelo) se aplica a la corriente I, NO a R ni X.

### Fórmula completa

```
e = factor × (I/N) × L × (R × cosθ + X × sinθ)
%e = (e / V_referencia) × 100
```

### Factores por sistema eléctrico

| Sistema         | Factor | Voltaje de referencia |
| --------------- | ------ | --------------------- |
| MONOFASICO 1F2H | 2.0    | Vfn                   |
| BIFASICO 2F3H   | 2.0    | Vfn                   |
| DELTA 3F3H      | √3     | Vff                   |
| ESTRELLA 3F4H   | √3     | Vfn                   |

## Reglas de Oro — Capa Domain

*Estas reglas son específicas para la capa Domain de cálculos. Ver [docs/reference/structure.md](../../../docs/reference/structure.md) para reglas globales.*

1. Domain nunca depende de Application ni Infrastructure
2. Sin I/O (no leer archivos, no HTTP, no DB)
3. Puro Go + lógica de negocio

## Referencias

- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)
