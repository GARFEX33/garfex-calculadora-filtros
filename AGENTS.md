# Garfex Calculadora Filtros

Backend API en Go para memorias de calculo de instalaciones electricas segun normativa NOM (Mexico).

## Como Usar Esta Guia

- Empieza aqui para normas globales del proyecto
- Cada feature y capa tiene su propio AGENTS.md con guias especificas
- El AGENTS.md del directorio tiene precedencia sobre este archivo cuando hay conflicto

## Regla de Skills (OBLIGATORIO)

**ANTES de cualquier accion, verificar si aplica un skill.** Si hay 1% de probabilidad de que aplique, invocar el skill con la herramienta `Skill`.

Orden de prioridad:

1. **Skills de proceso primero** (brainstorming, debugging) — determinan COMO abordar la tarea
2. **Skills de implementacion segundo** (golang-patterns, api-design) — guian la ejecucion

Si el skill tiene checklist, crear todos con TodoWrite antes de seguirlo.

## Workflow de Desarrollo (OBLIGATORIO)

Para cualquier feature o bugfix, seguir este flujo de skills en orden:

| Paso | Skill                    | Trigger                     | Que hace                                                                                                                          |
| ---- | ------------------------ | --------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| 1    | `brainstorming`          | Usuario pide feature/cambio | Refina ideas con preguntas, explora alternativas, presenta diseño por secciones para validar. Guarda documento de diseño.         |
| 2    | `writing-plans`          | Diseño aprobado             | Divide el trabajo en tareas pequeñas (2-5 min cada una). Cada tarea tiene: rutas exactas, código completo, pasos de verificación. |
| 3    | `executing-plans`        | Plan listo                  | Despacha subagente fresco por tarea con revisión de dos etapas (spec + calidad)                                                   |
| 4    | `requesting-code-review` | Entre tareas                | Revisa contra el plan, reporta issues por severidad. Issues críticos bloquean progreso.                                           |

**IMPORTANTE:** No saltear pasos. Si el usuario dice "agregá X", empezar con `brainstorming`, NO con código.

## Sistema de Agentes Especializados (OBLIGATORIO)

Cada capa tiene su propio agente especializado con su ciclo completo de trabajo. El coordinador delega a los agentes en este orden:

```
domain-agent → application-agent → infrastructure-agent
```

### Cuándo invocar cada agente

| Accion | Agente | Skills del agente |
| ------ | ------ | ----------------- |
| Crear/modificar entidades, value objects, servicios de dominio | `domain-agent` | `brainstorming-dominio` → `writing-plans-dominio` → `executing-plans-dominio` |
| Crear/modificar ports, use cases, DTOs | `application-agent` | `brainstorming-application` → `writing-plans-application` → `executing-plans-application` |
| Crear/modificar adapters, repositorios, handlers HTTP | `infrastructure-agent` | `brainstorming-infrastructure` → `writing-plans-infrastructure` → `executing-plans-infrastructure` |

### Reglas de delegación entre agentes

- **domain-agent** trabaja primero — no sabe que existen los otros agentes
- **application-agent** toma el output del domain-agent — no toca infraestructura
- **infrastructure-agent** toma el output del application-agent — implementa ports, no define reglas
- **Coordinador** es el único que conoce el orden y hace el wiring en `cmd/api/main.go`
- Cada agente crea sus propias tareas con TodoWrite antes de ejecutar
- Cada agente verifica con `go test` antes de entregar

## Estructura del Proyecto (Vertical Slices)

```
internal/
  shared/
    kernel/
      valueobject/      ← Corriente, Tension, Temperatura, MaterialConductor, Conductor, etc.
  calculos/             ← feature: memoria de cálculo eléctrico
    domain/
      entity/           ← TipoCanalizacion, SistemaElectrico, ITM, MemoriaCalculo, etc.
      service/          ← 7 servicios de cálculo eléctrico (IEEE-141, NOM)
    application/
      port/             ← TablaNOMRepository, EquipoRepository (interfaces)
      usecase/          ← OrquestadorMemoriaCalculo y micro use cases
      dto/              ← EquipoInput, MemoriaOutput
    infrastructure/
      adapter/
        driver/http/    ← CalculoHandler, formatters, middleware
        driven/csv/     ← CSVTablaNOMRepository
  equipos/              ← feature: catálogo de equipos (placeholder futuro)
    domain/
    application/
    infrastructure/
cmd/api/main.go         ← único lugar que conoce todas las features, hace wiring
data/tablas_nom/        ← tablas CSV NOM
tests/integration/
```

## Reglas de Aislamiento Entre Features (CRITICO)

- `calculos/` NUNCA importa `equipos/` y viceversa
- `shared/kernel/` NO importa ninguna feature
- `cmd/api/main.go` es el ÚNICO archivo que puede importar múltiples features
- Comunicación entre features en el futuro: solo vía interfaces en `shared/kernel/`

## Guias por Capa

| Capa                    | Ubicacion                                                   | AGENTS.md                                          |
| ----------------------- | ----------------------------------------------------------- | -------------------------------------------------- |
| Shared Kernel           | `internal/shared/kernel/`                                   | `internal/shared/kernel/AGENTS.md`                 |
| Feature Calculos        | `internal/calculos/`                                        | ver subcapas abajo                                 |
| Domain — Entity         | `internal/calculos/domain/entity/`                          | `internal/calculos/domain/AGENTS.md`               |
| Domain — Services       | `internal/calculos/domain/service/`                         | `internal/calculos/domain/AGENTS.md`               |
| Application             | `internal/calculos/application/`                            | `internal/calculos/application/AGENTS.md`          |
| Infrastructure          | `internal/calculos/infrastructure/`                         | `internal/calculos/infrastructure/AGENTS.md`       |
| Feature Equipos         | `internal/equipos/`                                         | `internal/equipos/AGENTS.md`                       |
| Datos NOM               | `data/tablas_nom/`                                          | `data/tablas_nom/AGENTS.md`                        |

## Skills Disponibles

### Skills Genericos

| Skill                   | Descripcion                                             | Ruta                                                      |
| ----------------------- | ------------------------------------------------------- | --------------------------------------------------------- |
| `golang-patterns`       | Patrones Go idiomaticos, error handling, interfaces     | [SKILL.md](.agents/skills/golang-patterns/SKILL.md)       |
| `golang-pro`            | Go avanzado: concurrencia, microservicios, performance  | [SKILL.md](.agents/skills/golang-pro/SKILL.md)            |
| `api-design-principles` | Diseno REST/GraphQL, convenciones API                   | [SKILL.md](.agents/skills/api-design-principles/SKILL.md) |
| `skill-creator`         | Crear nuevos skills siguiendo el spec de Agent Skills   | [SKILL.md](.agents/skills/skill-creator/SKILL.md)         |
| `skill-sync`            | Sincronizar metadata de skills a tablas Auto-invocacion | [SKILL.md](.agents/skills/skill-sync/SKILL.md)            |
| `commit-work`           | Commits de calidad: staging, split logico, mensajes     | [SKILL.md](.agents/skills/commit-work/SKILL.md)           |

### Skills de Proceso

| Skill                            | Descripcion                                        | Ruta                                                                   |
| -------------------------------- | -------------------------------------------------- | ---------------------------------------------------------------------- |
| `brainstorming`                  | Explorar ideas antes de implementar                | [SKILL.md](.agents/skills/brainstorming/SKILL.md)                      |
| `brainstorming-dominio`          | Diseñar dominio: entidades, VOs, servicios         | [SKILL.md](.agents/skills/brainstorming-dominio/SKILL.md)              |
| `brainstorming-application`      | Diseñar application: ports, use cases, DTOs        | [SKILL.md](.agents/skills/brainstorming-application/SKILL.md)          |
| `brainstorming-infrastructure`   | Diseñar infrastructure: adapters, repos            | [SKILL.md](.agents/skills/brainstorming-infrastructure/SKILL.md)       |
| `writing-plans-dominio`          | Plan de implementacion de dominio                  | [SKILL.md](.agents/skills/writing-plans-dominio/SKILL.md)              |
| `writing-plans-application`      | Plan de implementacion de application              | [SKILL.md](.agents/skills/writing-plans-application/SKILL.md)          |
| `writing-plans-infrastructure`   | Plan de implementacion de infrastructure           | [SKILL.md](.agents/skills/writing-plans-infrastructure/SKILL.md)       |
| `executing-plans-dominio`        | Ejecutar plan de dominio                           | [SKILL.md](.agents/skills/executing-plans-dominio/SKILL.md)            |
| `executing-plans-application`    | Ejecutar plan de application                       | [SKILL.md](.agents/skills/executing-plans-application/SKILL.md)        |
| `executing-plans-infrastructure` | Ejecutar plan de infrastructure                    | [SKILL.md](.agents/skills/executing-plans-infrastructure/SKILL.md)     |

### Skills de Proyecto

| Skill                                       | Descripcion                                                            | Ruta                                                                                          |
| ------------------------------------------- | ---------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| `agents-md-manager`                         | Crear y auditar jerarquia AGENTS.md                                    | [SKILL.md](.agents/skills/agents-md-manager/SKILL.md)                                        |
| `clean-ddd-hexagonal-vertical-go-enterprise`| Arquitectura Enterprise: Clean + DDD + Hexagonal + Vertical Slices     | [SKILL.md](.agents/skills/clean-ddd-hexagonal-vertical-go-enterprise/SKILL.md)               |

## Auto-invocacion

Cuando realices estas acciones, LEE el AGENTS.md o skill correspondiente PRIMERO:

| Accion                                        | Agente / Referencia                                          |
| --------------------------------------------- | ------------------------------------------------------------ |
| Crear/modificar entidad o value object        | `domain-agent` → `internal/calculos/domain/AGENTS.md`       |
| Crear/modificar servicio de calculo           | `domain-agent` → `internal/calculos/domain/AGENTS.md`       |
| Agregar value object al kernel compartido     | `domain-agent` → `internal/shared/kernel/AGENTS.md`         |
| Trabajar con ports o use cases                | `application-agent` → `internal/calculos/application/AGENTS.md` |
| Trabajar con DTOs o flujo de orquestacion     | `application-agent` → `internal/calculos/application/AGENTS.md` |
| Modificar repositorios o CSV reader           | `infrastructure-agent` → `internal/calculos/infrastructure/AGENTS.md` |
| Crear/modificar endpoints API o handlers      | `infrastructure-agent` → `internal/calculos/infrastructure/AGENTS.md` |
| Configurar BD o variables de entorno          | `infrastructure-agent` → `internal/calculos/infrastructure/AGENTS.md` |
| Trabajar con tablas NOM CSV                   | `data/tablas_nom/AGENTS.md`                                  |
| Agregar nueva tabla NOM                       | `data/tablas_nom/AGENTS.md`                                  |
| Aplicar patrones Go idiomaticos               | skill `golang-patterns`                                      |
| Crear/auditar AGENTS.md                       | skill `agents-md-manager`                                    |
| Disenar API endpoints                         | skill `api-design-principles`                                |
| Crear nuevo skill                             | skill `skill-creator`                                        |
| Sincronizar skills a AGENTS.md                | skill `skill-sync`                                           |
| Hacer commits o pull requests                 | skill `commit-work`                                          |

## Stack

Go 1.22+, Gin, PostgreSQL (pgx/v5), testify, golangci-lint

## Comandos

### Desarrollo

```bash
go test ./...           # Tests
go test -race ./...     # Tests con race detector
go build ./...          # Compilacion
go vet ./...            # Analisis estatico
golangci-lint run       # Linting completo
```

### Iniciar Servidor

**IMPORTANTE:** Asegurarse de que el puerto 8080 esté libre antes de iniciar:

```bash
# Opción 1: Compilar y ejecutar (recomendado)
go build -o server.exe ./cmd/api/main.go
./server.exe

# Opción 2: Ejecutar directamente (sin compilar)
go run cmd/api/main.go

# Opción 3: Puerto personalizado (si 8080 está ocupado)
set PORT=8090
go run cmd/api/main.go
```

**Verificar que el servidor está corriendo:**

```bash
curl http://localhost:8080/health
# Respuesta esperada: {"status":"ok"}
```

**Endpoint principal:**

```bash
curl -X POST http://localhost:8080/api/v1/calculos/memoria \
  -H "Content-Type: application/json" \
  -d '{"modo":"MANUAL_AMPERAJE","amperaje_nominal":50,"tension":220,"tipo_canalizacion":"TUBERIA_PVC","hilos_por_fase":1,"longitud_circuito":10,"itm":100,"factor_potencia":0.9,"estado":"Sonora","sistema_electrico":"DELTA","material":"Cu"}'
```

**Campos obligatorios:** `modo`, `tension`, `tipo_canalizacion`, `itm`, `longitud_circuito`, `estado`, `sistema_electrico`

**Campo `material`:** Opcional, valores: `"Cu"` (default) o `"Al"`

## Convenciones Globales

## Actualizacion de Documentacion

Al terminar cada tarea, actualizar: plan si diverge, AGENTS.md si cambia una regla, MEMORY.md si debe persistir entre sesiones.

## Documentacion

### Implementados (en `docs/plans/completed/`)

- Arquitectura inicial: `completed/2026-02-09-arquitectura-inicial-design.md`
- Domain layer: `completed/2026-02-10-domain-layer.md`
- Tablas NOM canalizacion: `completed/2026-02-11-tablas-nom-canalizacion-design.md`
- Caída de tension IEEE-141: `completed/2026-02-12-caida-tension-ieee141-design.md`
- Ports CSV infrastructure: `completed/2026-02-12-ports-csv-infrastructure-design.md`
- Material Cu/Al conductor tierra: `completed/2026-02-13-material-conductor-tierra-design.md`

### Refactorizacion Vertical Slices (en `docs/plans/`)

- Diseño: `2026-02-15-vertical-slices-refactor-design.md`
- Plan: `2026-02-15-vertical-slices-refactor-plan.md`
