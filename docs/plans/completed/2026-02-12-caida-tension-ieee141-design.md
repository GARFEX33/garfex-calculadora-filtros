# Diseño: Caída de Tensión — Fórmula IEEE-141 / NOM con Factor de Potencia

**Fecha:** 2026-02-12
**Estado:** Aprobado
**Contexto:** Fase 1 — Reemplazo del método impedancia Z por fórmula IEEE-141 con FP

## Problema

El servicio `CalcularCaidaTension` implementado en `2026-02-11-caida-tension-impedancia-design.md`
usa el método de impedancia puro:

```
%VD = (√3 × I × Z × L_km / V) × 100   donde Z = √(R² + X²)
```

Este método no considera el ángulo de la carga (factor de potencia). Para cargas con
FP < 1 el resultado es incorrecto porque usa el módulo completo de la impedancia en vez
de su proyección sobre el ángulo de la corriente.

## Fórmula Correcta: IEEE-141 / NOM

```
%e = 173 × (In / CF) × L_km × (R·cosθ + X·senθ) / E_FF
VD = E_FF × (%e / 100)
```

Donde:
- `173` = √3 × 100
- `In` = corriente nominal [A]
- `CF` = número de conductores por fase (HilosPorFase)
- `L_km` = longitud en kilómetros
- `R` = resistencia AC [Ω/km] de Tabla 9
- `X` = reactancia inductiva [Ω/km] de Tabla 9
- `cosθ` = factor de potencia de la carga
- `senθ` = √(1 - cosθ²)
- `E_FF` = tensión entre fases [V]

## Factor de Potencia por Tipo de Equipo

Los filtros activos y de rechazo, y los transformadores, son equipos de corrección o control
del FP, por lo que en el cálculo de caída de tensión se usa FP = 1.0:

| Equipo | FP usado | cosθ | senθ | Razón |
|--------|----------|------|------|-------|
| `FiltroActivo` | 1.0 | 1.0 | 0.0 | Equipo de corrección FP |
| `FiltroRechazo` | 1.0 | 1.0 | 0.0 | Equipo de corrección FP |
| `Transformador` | 1.0 | 1.0 | 0.0 | FP asumido en cálculo In |
| `Carga` | FP explícito | FP | √(1-FP²) | Campo `FactorPotencia` |

Con FP = 1 la fórmula se reduce a:
```
%e = 173 × (In/CF) × L_km × R / E_FF
```
Solo la resistencia impacta — la reactancia desaparece (senθ = 0).

## Reactancia X — Fuente: Tabla 9

Se elimina el cálculo geométrico DMG/RMG. La reactancia se lee directamente de
`tabla-9-resistencia-reactancia.csv`:

| TipoCanalizacion | Columna CSV |
|-----------------|-------------|
| `TUBERIA_PVC` | `reactancia_al` |
| `TUBERIA_ALUMINIO` | `reactancia_al` |
| `TUBERIA_ACERO_PG` | `reactancia_acero` |
| `TUBERIA_ACERO_PD` | `reactancia_acero` |
| `CHAROLA_CABLE_ESPACIADO` | `reactancia_al` |
| `CHAROLA_CABLE_TRIANGULAR` | `reactancia_al` |

**Razón charola → `reactancia_al`:** La charola no tiene conduit metálico, comportamiento
equivalente a aluminio.

## Struct de Entrada — Simplificado

```go
type EntradaCalculoCaidaTension struct {
    ResistenciaOhmPorKm float64                 // Tabla 9 → res_{material}_{conduit}
    ReactanciaOhmPorKm  float64                 // Tabla 9 → reactancia_al o reactancia_acero
    TipoCanalizacion    entity.TipoCanalizacion  // Para documentar en reporte
    HilosPorFase        int                     // CF ≥ 1
    FactorPotencia      float64                 // FA/FR/TR = 1.0 | Carga = FP explícito
}
```

### Campos eliminados respecto al diseño anterior

Los siguientes campos desaparecen porque ya no se necesita el cálculo geométrico:
- `DiametroExteriorMM` ❌
- `DiametroConductorMM` ❌
- `NumeroHilos` ❌

### Mapas eliminados del servicio

- `factorHilos` ❌ (factores RMG por número de hilos)
- `factorDMG` ❌ (factores DMG por tipo de canalización)

## Flujo Interno del Servicio (5 pasos)

```
1. cosθ  = entrada.FactorPotencia
2. senθ  = √(1 - cosθ²)
3. R_ef  = entrada.ResistenciaOhmPorKm / HilosPorFase
4. X_ef  = entrada.ReactanciaOhmPorKm  / HilosPorFase
5. %e    = 173 × corriente × (distancia/1000) × (R_ef·cosθ + X_ef·senθ) / tension
6. VD    = tension × (%e / 100)
```

## Struct de Resultado — Actualizado Semánticamente

```go
type ResultadoCaidaTension struct {
    Porcentaje  float64  // %e
    CaidaVolts  float64  // VD en volts
    Cumple      bool     // %e ≤ limiteNOM
    Impedancia  float64  // término efectivo: R_ef·cosθ + X_ef·senθ  [Ω/km]
    Resistencia float64  // R_ef = R / HilosPorFase  [Ω/km]
    Reactancia  float64  // X_ef = X / HilosPorFase  [Ω/km]
}
```

`Impedancia` pasa a representar el **término efectivo** `R·cosθ + X·senθ` en vez de `√(R²+X²)`.
El nombre se mantiene para no romper el struct; el reporte de memoria de cálculo lo mostrará
como "Impedancia efectiva".

## Errores

| Error | Estado |
|-------|--------|
| `ErrDistanciaInvalida` | Se mantiene |
| `ErrHilosPorFaseInvalido` | Se mantiene |
| `ErrNumeroHilosDesconocido` | **Eliminado** (ya no se usa NumeroHilos) |
| `ErrFactorPotenciaInvalido` | **Nuevo**: FP fuera de rango (0, 1] |

## Validaciones

```
distancia    > 0            → ErrDistanciaInvalida
HilosPorFase > 0            → ErrHilosPorFaseInvalido
FactorPotencia > 0 && ≤ 1  → ErrFactorPotenciaInvalido
```

## Ejemplo de Validación

Datos: 2 AWG Cu, tubería PVC, FP = 0.85, 120 A, 30 m, 480 V

Lookup Tabla 9:
- `res_cu_pvc` para 2 AWG = 0.62 Ω/km
- `reactancia_al` para 2 AWG = 0.051 Ω/km (valor típico NOM)

Cálculo:
```
cosθ  = 0.85
senθ  = √(1 - 0.85²) = √(1 - 0.7225) = √0.2775 = 0.5268
R_ef  = 0.62 / 1 = 0.62 Ω/km
X_ef  = 0.051 / 1 = 0.051 Ω/km
term  = 0.62×0.85 + 0.051×0.5268 = 0.527 + 0.0269 = 0.5539 Ω/km
%e    = 173 × 120 × 0.030 × 0.5539 / 480
      = 173 × 120 × 0.030 × 0.5539 / 480
      = 344.83 / 480 = 0.719%
VD    = 480 × 0.00719 = 3.45 V
```

Con FP = 1.0 (FA/FR/TR):
```
term  = 0.62×1.0 + 0.051×0.0 = 0.62 Ω/km
%e    = 173 × 120 × 0.030 × 0.62 / 480 = 387.07 / 480 = 0.807%
```

## Impacto en Código

### Cambios en domain/entity

| Archivo | Cambio |
|---------|--------|
| `memoria_calculo.go` | `ResultadoCaidaTension.Impedancia` — semántica actualizada en comentario |

### Cambios en domain/service

| Archivo | Cambio |
|---------|--------|
| `calculo_caida_tension.go` | Reescribir: nueva fórmula, nuevo struct entrada, eliminar campos geométricos |
| `calculo_caida_tension_test.go` | Reescribir: nuevos casos con FP variable, casos FA/FR/TR con FP=1 |

### Sin cambios

- `tipo_canalizacion.go` — 6 valores se mantienen
- `SeleccionarConductorAlimentacion` — independiente
- `SeleccionarConductorTierra` — independiente
- `CalcularCanalizacion` — independiente

### Cambios en infrastructure/application (Fase 2)

| Cambio | Descripción |
|--------|-------------|
| CSV reader Tabla 9 | Leer `reactancia_al` y `reactancia_acero` por calibre |
| Resolver `EntradaCalculoCaidaTension` | Mapear `TipoCanalizacion` → columna X correcta |
| Resolver `FactorPotencia` | Leer FP del equipo; usar 1.0 para FA/FR/TR |

## Decisiones Clave

1. **FP = 1.0 para FA/FR/TR** — son equipos de corrección, FP conocido y fijo
2. **X de Tabla 9** — elimina cálculo geométrico, simplifica el servicio
3. **Charola → `reactancia_al`** — sin conduit metálico, comportamiento equivalente
4. **FP en `EntradaCalculoCaidaTension`** — la capa application lo resuelve desde el equipo
5. **`Impedancia` = término efectivo** — mantiene el campo, cambia semántica para el reporte
6. **`173`** — constante explícita en código (no `math.Sqrt(3) * 100`) para legibilidad NOM
