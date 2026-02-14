# Domain — Entity

Entidades y tipos del dominio. Sin dependencias externas.

> **Skills Reference**:
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — constructores, error handling, interfaces idiomaticas

### Auto-invoke

| Accion | Skill |
|--------|-------|
| Crear nueva entidad o tipo | `golang-patterns` |
| Agregar nuevo TipoEquipo o TipoCanalizacion | `golang-patterns` |
| Implementar nueva interfaz (`CalculadorCorriente`) | `golang-patterns` |

## Tipos y Constantes

**TipoEquipo (4):** `FILTRO_ACTIVO`, `FILTRO_RECHAZO`, `TRANSFORMADOR`, `CARGA`

**TipoCanalizacion (6):** `TUBERIA_PVC`, `TUBERIA_ALUMINIO`, `TUBERIA_ACERO_PG`, `TUBERIA_ACERO_PD`, `CHAROLA_CABLE_ESPACIADO`, `CHAROLA_CABLE_TRIANGULAR`

**SistemaElectrico (4):** `DELTA`, `ESTRELLA`, `BIFASICO`, `MONOFASICO`
- Determina cantidad de conductores: Delta=3, Estrella=4, Bifasico=3, Monofasico=2

Los 4 tipos de tuberia comparten tabla de ampacidad. Cada tipo mapea a columna R diferente en Tabla 9.

## Formulas por Tipo de Equipo

| Tipo | Parametro | Formula In |
|------|-----------|-----------|
| FiltroActivo | AmperajeNominal | In = directo |
| FiltroRechazo | KVAR | In = KVAR / (KV x sqrt3) |
| Transformador | KVA | In = KVA / (KV x sqrt3) |
| Carga | KW + FP + Fases | In = KW / (KV x factor x FP) |

Todos implementan `CalculadorCorriente`. Carga y Transformador tambien implementan `CalculadorPotencia`.

## MemoriaCalculo y ResultadoCaidaTension

`MemoriaCalculo` agrupa todos los pasos: `CorrienteNominal`, `CorrienteAjustada`, `FactoresAjuste`, `Potencias`, `ConductorAlimentacion`, `ConductorTierra`, `TipoCanalizacion`, `Canalizacion`, `TemperaturaUsada`, `CaidaTension`, `CumpleNormativa`.

`ResultadoCaidaTension` vive en `entity/` (no en `service/`) para evitar ciclo de dependencias — `service` importa `entity`. Contiene: `Porcentaje`, `CaidaVolts`, `Cumple`, `Impedancia` (término efectivo R·cosθ + X·senθ), `Resistencia` (R_ef), `Reactancia` (X_ef).

## Patron para nuevo equipo (Fase 2+)

1. Crear `entity/xxx.go` implementando `CalculadorCorriente`
2. Agregar constante en `TipoEquipo`
3. Agregar test table-driven en `entity/xxx_test.go`
4. Los servicios existentes no requieren cambio (reciben interfaces)

---

## CRITICAL RULES

- ALWAYS: Toda entidad implementa una interfaz (`CalculadorCorriente` o `CalculadorPotencia`)
- ALWAYS: Errores de entidades viven en `errors.go` con el patron `ErrXxx = errors.New(...)`
- ALWAYS: Constantes de tipo (`TipoEquipo`, `TipoCanalizacion`) como `string` tipado — nunca `int` o `iota`
- ALWAYS: Interfaces pequeñas — un solo metodo si es posible; definir donde se consumen, no donde se implementan
- NEVER: Dependencias externas (sin Gin, sin pgx, sin encoding/csv)
- NEVER: Campos mutables expuestos sin constructor

---

## NAMING CONVENTIONS

| Entidad | Patron | Ejemplo |
|---------|--------|---------|
| Tipo constante | `SCREAMING_SNAKE_CASE` | `FILTRO_ACTIVO`, `TUBERIA_PVC` |
| Struct entidad | `PascalCase` | `FiltroActivo`, `Transformador` |
| Interfaz | `PascalCase` + verbo | `CalculadorCorriente` |
| Error sentinel | `ErrPascalCase` | `ErrTipoEquipoInvalido` |
| Archivo | `snake_case.go` | `filtro_activo.go` |

---

## QA CHECKLIST

- [ ] `go test ./internal/domain/entity/...` pasa
- [ ] Nueva entidad implementa `CalculadorCorriente`
- [ ] Test table-driven cubre todos los casos de la formula
- [ ] Sin imports externos al stdlib de Go
