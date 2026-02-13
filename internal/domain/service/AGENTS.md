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
| CalcularCaidaTension | `calculo_caida_tension.go` | Metodo impedancia Z = sqrt(R²+X²) |

## Caida de Tension (metodo impedancia)

Formula: `VD = sqrt3 x I x Z x L_km`

- **R:** recibida de Tabla 9 (pre-resuelta por infrastructure)
- **X:** calculada geometricamente con DMG/RMG
- **RMG** = (diametro_desnudo/2) x factorHilos[numHilos]
- **DMG** = diametro_exterior_thw x factorDMG[tipoCanalizacion]

| Parametro | Valores |
|-----------|---------|
| Factores hilos | 1=0.7788, 7=0.726, 19=0.758, 37=0.768, 61=0.772 |
| Factores DMG | tuberia/triangular=1.0, espaciado=2.0 |

Diseno completo: `docs/plans/2026-02-11-caida-tension-impedancia-design.md`

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
