# Diseño: Refactorización a Vertical Slices + Sistema de Agentes Especializados

**Fecha:** 2026-02-15  
**Estado:** Aprobado  
**Objetivo:** Refactorizar la estructura de carpetas a Vertical Slices con un solo `go.mod`, e implementar un sistema de agentes especializados por capa con sus propios ciclos de trabajo.

---

## Contexto

El proyecto actual tiene una estructura por capas horizontales (`domain/`, `application/`, `infrastructure/`, `presentation/`) dentro de `internal/`. Esta estructura funciona bien para el scope actual pero no escala bien cuando se agreguen nuevos módulos (equipos, proyectos, reportes, etc.).

**Driver principal:** Escalabilidad. Que nada dependa de nada entre features.

---

## Decisiones de Diseño

| Decisión | Elección | Razón |
|----------|----------|-------|
| Estrategia de módulos | Un solo `go.mod` con vertical slices | `go.work` tiene overhead administrativo alto sin beneficio real a este escala |
| Shared code | `shared/kernel/` solo value objects estables | Los VOs eléctricos son usados por múltiples features (equipos también tendrá tensión/corriente) |
| Features | `calculos/` completo, `equipos/` placeholder vacío | `equipos/` no tiene implementación real hoy, pero se prepara la estructura |
| Handlers HTTP | Dentro de `infrastructure/adapter/driver/http/` | Alineado con Hexagonal estricta: los adapters de entrada son infraestructura |

---

## Estructura de Carpetas Objetivo

```
garfex-calculadora-filtros/
├── go.mod                                      ← un solo módulo Go
├── go.sum
├── cmd/
│   └── api/
│       └── main.go                             ← único lugar que conoce todo, hace wiring
├── internal/
│   ├── shared/
│   │   └── kernel/
│   │       ├── valueobject/
│   │       │   ├── corriente.go
│   │       │   ├── tension.go
│   │       │   ├── temperatura.go
│   │       │   ├── material_conductor.go
│   │       │   ├── conductor.go
│   │       │   ├── resistencia_reactancia.go
│   │       │   └── charola.go
│   │       └── errors.go
│   │
│   ├── calculos/                               ← feature: memoria de cálculo eléctrico
│   │   ├── domain/
│   │   │   ├── entity/
│   │   │   │   ├── tipo_canalizacion.go
│   │   │   │   ├── sistema_electrico.go
│   │   │   │   ├── itm.go
│   │   │   │   ├── tipo_equipo.go
│   │   │   │   ├── memoria_calculo.go
│   │   │   │   ├── transformador.go
│   │   │   │   ├── filtro_activo.go
│   │   │   │   ├── filtro_rechazo.go
│   │   │   │   ├── canalizacion.go
│   │   │   │   ├── carga.go
│   │   │   │   └── equipo.go
│   │   │   ├── service/
│   │   │   │   ├── calculo_corriente_nominal.go
│   │   │   │   ├── ajuste_corriente.go
│   │   │   │   ├── calculo_conductor.go
│   │   │   │   ├── calculo_tierra.go
│   │   │   │   ├── calculo_canalizacion.go
│   │   │   │   ├── calculo_caida_tension.go
│   │   │   │   ├── calcular_factor_temperatura.go
│   │   │   │   ├── calcular_factor_agrupamiento.go
│   │   │   │   ├── calcular_charola_espaciado.go
│   │   │   │   ├── calcular_charola_triangular.go
│   │   │   │   └── seleccionar_temperatura.go
│   │   │   └── errors.go
│   │   ├── application/
│   │   │   ├── port/
│   │   │   │   ├── tabla_nom_repository.go
│   │   │   │   ├── equipo_repository.go
│   │   │   │   └── seleccionar_temperatura.go
│   │   │   ├── usecase/
│   │   │   │   ├── calcular_corriente.go
│   │   │   │   ├── ajustar_corriente.go
│   │   │   │   ├── seleccionar_conductor.go
│   │   │   │   ├── dimensionar_canalizacion.go
│   │   │   │   ├── calcular_caida_tension.go
│   │   │   │   ├── calcular_memoria.go
│   │   │   │   ├── orquestador_memoria.go
│   │   │   │   └── helpers/
│   │   │   │       └── nombre_tabla.go
│   │   │   └── dto/
│   │   │       ├── equipo_input.go
│   │   │       ├── memoria_output.go
│   │   │       └── errors.go
│   │   └── infrastructure/
│   │       ├── adapter/
│   │       │   ├── driver/
│   │       │   │   └── http/
│   │       │   │       ├── calculo_handler.go
│   │       │   │       ├── formatters/
│   │       │   │       │   ├── nombre_tabla.go
│   │       │   │       │   └── observaciones.go
│   │       │   │       └── middleware/
│   │       │   │           └── .gitkeep
│   │       │   └── driven/
│   │       │       └── csv/
│   │       │           ├── csv_tabla_nom_repository.go
│   │       │           ├── seleccionar_temperatura.go
│   │       │           └── testdata/
│   │       └── router.go
│   │
│   └── equipos/                                ← feature: catálogo de equipos (placeholder)
│       ├── domain/
│       │   └── .gitkeep
│       ├── application/
│       │   ├── port/
│       │   │   └── .gitkeep
│       │   └── dto/
│       │       └── .gitkeep
│       └── infrastructure/
│           └── adapter/
│               └── driven/
│                   └── postgres/
│                       └── .gitkeep
│
├── data/
│   └── tablas_nom/                             ← CSVs NOM, no cambian
├── tests/
│   └── integration/
└── docs/
    └── plans/
```

---

## Flujo de Dependencias

```
shared/kernel/              ← no depende de nada
      ↑
calculos/domain/            ← importa solo kernel
      ↑
calculos/application/       ← importa domain + kernel
      ↑
calculos/infrastructure/    ← importa application + kernel

cmd/api/main.go             ← único que conoce todo, hace wiring completo
```

**Regla absoluta:** `calculos/` nunca importa `equipos/` y viceversa. Solo `main.go` los conoce a los dos.

---

## Qué se Mueve y Adónde

| Código actual | Destino |
|---|---|
| `internal/domain/entity/` | `internal/calculos/domain/entity/` |
| `internal/domain/service/` | `internal/calculos/domain/service/` |
| `internal/domain/valueobject/` | `internal/shared/kernel/valueobject/` |
| `internal/application/port/` | `internal/calculos/application/port/` |
| `internal/application/usecase/` | `internal/calculos/application/usecase/` |
| `internal/application/dto/` | `internal/calculos/application/dto/` |
| `internal/infrastructure/repository/` | `internal/calculos/infrastructure/adapter/driven/csv/` |
| `internal/infrastructure/client/` | `internal/calculos/infrastructure/adapter/driven/postgres/` |
| `internal/presentation/handler/` | `internal/calculos/infrastructure/adapter/driver/http/` |
| `internal/presentation/router.go` | `internal/calculos/infrastructure/router.go` |
| `internal/presentation/formatters/` | `internal/calculos/infrastructure/adapter/driver/http/formatters/` |
| `internal/presentation/middleware/` | `internal/calculos/infrastructure/adapter/driver/http/middleware/` |

---

## Estrategia de Migración (sin romper tests)

Migrar de adentro hacia afuera, verificando `go test ./...` en cada paso:

| Paso | Acción | Verificación |
|------|--------|--------------|
| 1 | Crear `shared/kernel/` y mover value objects | `go test ./internal/shared/...` |
| 2 | Crear `calculos/domain/` y mover entities + services | `go test ./internal/calculos/domain/...` |
| 3 | Crear `calculos/application/` y mover ports + use cases + DTOs | `go test ./internal/calculos/application/...` |
| 4 | Crear `calculos/infrastructure/` y mover repos + adapters HTTP | `go test ./internal/calculos/infrastructure/...` |
| 5 | Actualizar `cmd/api/main.go` con nuevo wiring | `go build ./...` + `go test ./...` |
| 6 | Crear `equipos/` placeholder vacío | `go build ./...` |
| 7 | Eliminar carpetas viejas (`internal/domain/`, `internal/application/`, etc.) | `go test ./...` verde completo |
| 8 | Actualizar AGENTS.md raíz + cada capa | revisión manual |

**Regla:** Si algo se rompe en un paso, se arregla ahí. No se acumula deuda.

---

## Sistema de Agentes Especializados

### Arquitectura del Sistema

```
┌─────────────────────────────────────────────────────────┐
│              Agente Coordinador (este chat)             │
│    brainstorming → writing-plans → ejecuta agentes      │
└─────────────────────────────────────────────────────────┘
         │               │                │
         ▼               ▼                ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────────────┐
│ domain-agent │ │application-  │ │ infrastructure-agent  │
│              │ │agent         │ │                       │
└──────────────┘ └──────────────┘ └──────────────────────┘
```

### Skills por Agente

| Agente | Scope | brainstorming | writing-plans | executing-plans |
|--------|-------|--------------|---------------|-----------------|
| `domain-agent` | `calculos/domain/` + `shared/kernel/` | `brainstorming-dominio` | `writing-plans-dominio` | `executing-plans-dominio` |
| `application-agent` | `calculos/application/` | `brainstorming-application` | `writing-plans-application` | `executing-plans-application` |
| `infrastructure-agent` | `calculos/infrastructure/` | `brainstorming-infrastructure` | `writing-plans-infrastructure` | `executing-plans-infrastructure` |

### Reglas de Delegación

- `domain-agent` no sabe que existen `application-agent` ni `infrastructure-agent`
- `application-agent` puede consultar output de `domain-agent` pero no lo invoca directamente
- `infrastructure-agent` depende del output de `application-agent` (ports definidos)
- El coordinador es el único que conoce el orden y la secuencia entre agentes
- Cada agente trabaja en su capa y entrega output antes de que el siguiente empiece

### Flujo de Trabajo Entre Agentes

```
Coordinador diseña
      ↓
domain-agent: brainstorming → writing-plans → executing-plans
      ↓  (domain completo y testeado)
application-agent: brainstorming → writing-plans → executing-plans
      ↓  (application completo y testeado)
infrastructure-agent: brainstorming → writing-plans → executing-plans
      ↓  (infrastructure completo y testeado)
Coordinador: wiring en main.go + go test ./... verde
```

---

## Impacto en AGENTS.md

### AGENTS.md Raíz (actualizar)

- Nueva tabla de estructura de carpetas (vertical slices)
- Nueva sección "Sistema de Agentes": cuándo invocar cada uno, orden de delegación
- Actualizar tabla de auto-invocación con los nuevos paths

### AGENTS.md por Capa (crear/actualizar)

| Archivo | Acción | Contenido clave |
|---------|--------|-----------------|
| `internal/shared/kernel/AGENTS.md` | Crear | Qué va acá, qué NO va acá, regla de estabilidad |
| `internal/calculos/AGENTS.md` | Crear | Overview de la feature, links a subcapas |
| `internal/calculos/domain/AGENTS.md` | Actualizar path | Mismo contenido, nuevo path |
| `internal/calculos/application/AGENTS.md` | Actualizar path | Mismo contenido, nuevo path |
| `internal/calculos/infrastructure/AGENTS.md` | Crear | Reglas de adapters, driver vs driven |
| `internal/equipos/AGENTS.md` | Crear | Placeholder, scope futuro |

---

## Criterios de Éxito

- [ ] `go test ./...` pasa completamente después de la migración
- [ ] `go build ./...` sin errores
- [ ] Ningún cross-import entre `calculos/` y `equipos/`
- [ ] `shared/kernel/` no importa ninguna feature
- [ ] Todos los AGENTS.md actualizados con nuevos paths
- [ ] `cmd/api/main.go` es el único archivo que importa múltiples features
