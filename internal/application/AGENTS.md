# Application Layer

Orquesta domain services. Define contratos (ports), no implementaciones.

> **Skills Reference**:
>
> - [`golang-patterns`](.agents/skills/golang-patterns/SKILL.md) — interfaces pequeñas, error wrapping, convenciones de ports

### Auto-invoke

| Accion                         | Skill             |
| ------------------------------ | ----------------- |
| Definir nuevo port (interface) | `golang-patterns` |
| Crear o modificar use case     | `golang-patterns` |
| Agregar o modificar DTOs       | `golang-patterns` |

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
7. Resolver R y X de Tabla 9 + FP del equipo — llamar CalcularCaidaTension

## Seleccion de Temperatura (logica aqui, no en domain)

- <= 100A -> 60C (o 75C si charola triangular sin columna 60C)
- > 100A -> 75C
- 90C solo con `temperatura_override: 90` explicito del usuario

## DTOs

- **EquipoInput:** modo (LISTADO/MANUAL_AMPERAJE/MANUAL_POTENCIA), datos del equipo, parametros de instalacion, TipoCanalizacion, TemperaturaOverride, **Material** (Cu/Al, default Cu)
- **MemoriaOutput:** resultado completo de todos los pasos para el reporte

### Campo Material (Cu/Al)

- Tipo: `valueobject.MaterialConductor` (internamente int, JSON como string "CU"/"AL")
- Valores aceptados: "Cu", "cu", "CU", "cobre", "Al", "al", "AL", "aluminio"
- Default: Cobre (`MaterialCobre`) si no se especifica o string vacío
- Afecta:
  - Paso 4: Selección de conductor de alimentación (tabla ampacidad Cu vs Al)
  - Paso 5: Selección de conductor de tierra según ITM + material
  - Paso 7: Cálculo de caída de tensión (R y X diferente por material)

### Conductor de Tierra - Lógica de Material (Paso 5)

1. Buscar entrada en tabla 250-122 donde `ITM <= ITMHasta`
2. Si `material == Al` Y entrada tiene `ConductorAl` → usar Al
3. Si `material == Al` PERO entrada NO tiene `ConductorAl` → **fallback a Cu** (regla NOM)
4. Si `material == Cu` → usar siempre Cu

Tabla 250-122 solo tiene Al para ITM > 100A (aprox). Para ITM ≤ 100A, Al no está definido → siempre fallback a Cu.

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
- ALWAYS: Seguir el orden obligatorio
- ALWAYS: Logica de seleccion de temperatura aqui, no en domain ni infrastructure
- NEVER: Llamar directamente a infrastructure — solo via ports

### DTOs

- ALWAYS: Structs planos sin metodos de negocio
- ALWAYS: Validacion de input en el use case, no en el DTO
- NEVER: Entidades de domain expuestas directamente en DTOs — mapear siempre

---

## NAMING CONVENTIONS

| Entidad          | Patron                 | Ejemplo                                    |
| ---------------- | ---------------------- | ------------------------------------------ |
| Port (interface) | `PascalCaseRepository` | `EquipoRepository`, `TablaNOMRepository`   |
| Use case struct  | `PascalCaseUseCase`    | `CalcularMemoriaUseCase`                   |
| DTO entrada      | `PascalCaseInput`      | `EquipoInput`                              |
| DTO salida       | `PascalCaseOutput`     | `MemoriaOutput`                            |
| Error sentinel   | `ErrPascalCase`        | `ErrEquipoNoEncontrado`, `ErrModoInvalido` |

---

## QA CHECKLIST

- [ ] `go test ./internal/application/...` pasa
- [ ] Flujo del UseCase respeta el orden obligatorio de 7 pasos
- [ ] Seleccion de temperatura implementada en application (no en domain)
- [ ] Sin imports de infrastructure
- [ ] Errores de flujo usan sentinels definidos en esta capa
