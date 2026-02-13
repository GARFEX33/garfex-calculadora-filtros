# Domain — Services

Logica de calculo pura. Reciben datos ya interpretados — sin I/O, sin CSV, sin BD.

> **Skills Reference**:
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — patrones idiomáticos, interfaces, error handling

### Auto-invoke

| Accion | Skill |
|--------|-------|
| Crear o modificar servicio de calculo | `golang-patterns` |
| Implementar nueva formula NOM | `golang-patterns` |

## Servicios (6)

| Servicio | Archivo | Responsabilidad |
|----------|---------|-----------------|
| CorrienteNominal | `calculo_corriente_nominal.go` | Calcula In segun TipoEquipo |
| AjusteCorriente | `ajuste_corriente.go` | Aplica factores temperatura y agrupamiento |
| SeleccionarConductorAlimentacion | `calculo_conductor.go` | Recibe `[]EntradaTablaConductor` pre-resueltos |
| SeleccionarConductorTierra | `calculo_tierra.go` | ITM + tabla 250-122 pre-resuelta |
| CalcularCanalizacion | `calculo_canalizacion.go` | 40% fill NOM para tuberia |
| CalcularCaidaTension | `calculo_caida_tension.go` | Formula IEEE-141/NOM con factor de potencia |

## Caida de Tension (formula IEEE-141 / NOM)

```
%Vd = (√3 × Ib × L × (R·cosθ + X·senθ) / (V × N)) × 100
VD  = V × (%Vd / 100)
```

- **R:** de Tabla 9, columna `res_{material}_{conduit}` (pre-resuelta por infrastructure)
- **X:** de Tabla 9, columna `reactancia_al` o `reactancia_acero` (pre-resuelta por infrastructure)
- **cosθ = FactorPotencia:** FA/FR/TR = 1.0 fijo | Carga = FP explícito del equipo
- **N = HilosPorFase:** conductores en paralelo por fase

### EntradaCalculoCaidaTension

| Campo | Tipo | Fuente |
|-------|------|--------|
| `ResistenciaOhmPorKm` | float64 | Tabla 9 → columna R según material + conduit |
| `ReactanciaOhmPorKm` | float64 | Tabla 9 → `reactancia_al` o `reactancia_acero` |
| `TipoCanalizacion` | TipoCanalizacion | Para documentar en reporte |
| `HilosPorFase` | int | CF ≥ 1 |
| `FactorPotencia` | float64 | FA/FR/TR = 1.0 | Carga = FP explícito |

### ResultadoCaidaTension (vive en entity/)

| Campo | Semántica |
|-------|-----------|
| `Porcentaje` | %Vd calculado |
| `CaidaVolts` | VD en volts |
| `Cumple` | %Vd ≤ limiteNOM |
| `Impedancia` | Término efectivo R·cosθ + X·senθ (Ω/km) |
| `Resistencia` | R_ef = R / N (Ω/km) |
| `Reactancia` | X_ef = X / N (Ω/km) |

### Mapeo X en Tabla 9

| TipoCanalizacion | Columna X |
|-----------------|-----------|
| PVC, Aluminio, Charola espaciado, Charola triangular | `reactancia_al` |
| Acero PG, Acero PD | `reactancia_acero` |

Diseno completo: `docs/plans/2026-02-12-caida-tension-ieee141-design.md`

## Regla clave

Los servicios reciben interfaces (`CalculadorCorriente`), no tipos concretos.
Nunca toman decisiones sobre archivos, BD o temperatura — eso es responsabilidad de application.

---

## CRITICAL RULES

### Pureza del Servicio
- ALWAYS: Recibir interfaces, nunca tipos concretos (`CalculadorCorriente`, no `FiltroActivo`)
- ALWAYS: Sin estado — funciones puras o structs sin campos mutables
- NEVER: I/O de ningún tipo (sin CSV, sin BD, sin HTTP)
- NEVER: Decisiones de temperatura o canalizacion — eso es de `application`
- NEVER: Importar `infrastructure` o `presentation`

### Error Handling
- ALWAYS: Retornar `error` envuelto con contexto: `fmt.Errorf("CalcularCaidaTension: %w", err)`
- NEVER: Panic en servicios de calculo — siempre retornar error

---

## NAMING CONVENTIONS

| Entidad | Patron | Ejemplo |
|---------|--------|---------|
| Archivo servicio | `calculo_snake_case.go` | `calculo_corriente_nominal.go` |
| Funcion/metodo principal | `VerboPascalCase` | `Calcular`, `Seleccionar` |
| Struct entrada tabla | `EntradaPascalCase` | `EntradaTablaConductor` |
| Error sentinel | `ErrPascalCase` | `ErrConductorNoEncontrado` |

---

## QA CHECKLIST

- [ ] `go test ./internal/domain/service/...` pasa
- [ ] `go test -race ./internal/domain/service/...` pasa
- [ ] Tests table-driven cubren casos limite (corriente 0, longitud 0, calibre no encontrado)
- [ ] Sin imports de infrastructure, presentation o librerias externas
- [ ] Errores envueltos con contexto del servicio
