# Application Layer

Orquesta domain services. Define contratos (ports), no implementaciones.

> **Skills Reference**:
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — interfaces pequeñas, error wrapping, convenciones de ports

### Auto-invoke

| Accion | Skill |
|--------|-------|
| Definir nuevo port (interface) | `golang-patterns` |
| Crear o modificar use case | `golang-patterns` |
| Agregar o modificar DTOs | `golang-patterns` |

## Estructura

- `port/` — Interfaces que infrastructure implementa
- `usecase/` — CalcularMemoriaUseCase (orquesta los 7 pasos)
- `dto/` — EquipoInput, MemoriaOutput (entrada/salida de la API)

## Ports (interfaces)

- **EquipoRepository** — buscar equipos en BD (PostgreSQL)
- **TablaNOMRepository** — leer tablas CSV de ampacidad, tierra, impedancia

Las interfaces se definen aqui, se implementan en `infrastructure/`.
Pequenas y enfocadas (pocos metodos por interface).

## Flujo del UseCase (orden obligatorio)

1. Corriente Nominal (segun TipoEquipo)
2. Ajuste de Corriente (factores)
3. Seleccionar TipoCanalizacion — determina tabla NOM
4. Resolver tabla ampacidad + columna temperatura — llamar SeleccionarConductorAlimentacion
5. Conductor de Tierra (ITM -> tabla 250-122)
6. Dimensionar Canalizacion (40% fill)
7. Resolver datos Tablas 9/5/8 — llamar CalcularCaidaTension

## Seleccion de Temperatura (logica aqui, no en domain)

- <= 100A -> 60C (o 75C si charola triangular sin columna 60C)
- > 100A -> 75C
- 90C solo con `temperatura_override: 90` explicito del usuario

## DTOs

- **EquipoInput:** modo (LISTADO/MANUAL_AMPERAJE/MANUAL_POTENCIA), datos del equipo, parametros de instalacion, TipoCanalizacion, TemperaturaOverride
- **MemoriaOutput:** resultado completo de todos los pasos para el reporte

## Convenciones

- `context.Context` como primer parametro en operaciones I/O
- Errores de flujo: `ErrEquipoNoEncontrado`, `ErrModoInvalido`
- DTOs son structs planos sin logica de negocio
- Nunca importar infrastructure — solo domain

---

## CRITICAL RULES

### Ports
- ALWAYS: Interfaces pequeñas y enfocadas — pocos metodos por port
- ALWAYS: Definir ports en `application/port/`, implementar en `infrastructure/`
- NEVER: Logica de negocio en ports — solo contratos
- NEVER: Importar tipos de infrastructure en application

### Use Cases
- ALWAYS: `context.Context` como primer parametro
- ALWAYS: Seguir el orden obligatorio de los 7 pasos del flujo
- ALWAYS: Logica de seleccion de temperatura aqui, no en domain ni infrastructure
- NEVER: Llamar directamente a infrastructure — solo via ports

### DTOs
- ALWAYS: Structs planos sin metodos de negocio
- ALWAYS: Validacion de input en el use case, no en el DTO
- NEVER: Entidades de domain expuestas directamente en DTOs — mapear siempre

---

## NAMING CONVENTIONS

| Entidad | Patron | Ejemplo |
|---------|--------|---------|
| Port (interface) | `PascalCaseRepository` | `EquipoRepository`, `TablaNOMRepository` |
| Use case struct | `PascalCaseUseCase` | `CalcularMemoriaUseCase` |
| DTO entrada | `PascalCaseInput` | `EquipoInput` |
| DTO salida | `PascalCaseOutput` | `MemoriaOutput` |
| Error sentinel | `ErrPascalCase` | `ErrEquipoNoEncontrado`, `ErrModoInvalido` |

---

## QA CHECKLIST

- [ ] `go test ./internal/application/...` pasa
- [ ] Flujo del UseCase respeta el orden obligatorio de 7 pasos
- [ ] Seleccion de temperatura implementada en application (no en domain)
- [ ] Sin imports de infrastructure
- [ ] Errores de flujo usan sentinels definidos en esta capa
