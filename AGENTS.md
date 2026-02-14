# Garfex Calculadora Filtros

Backend API en Go para memorias de calculo de instalaciones electricas segun normativa NOM (Mexico).

## Como Usar Esta Guia

- Empieza aqui para normas globales del proyecto
- Cada capa tiene su propio AGENTS.md con guias especificas
- El AGENTS.md del directorio tiene precedencia sobre este archivo cuando hay conflicto

## Regla de Skills (OBLIGATORIO)

**ANTES de cualquier accion, verificar si aplica un skill.** Si hay 1% de probabilidad de que aplique, invocar el skill con la herramienta `Skill`.

Orden de prioridad:
1. **Skills de proceso primero** (brainstorming, debugging) — determinan COMO abordar la tarea
2. **Skills de implementacion segundo** (golang-patterns, api-design) — guian la ejecucion

Si el skill tiene checklist, crear todos con TodoWrite antes de seguirlo.

## Workflow de Desarrollo (OBLIGATORIO)

Para cualquier feature o bugfix, seguir este flujo de skills en orden:

| Paso | Skill | Trigger | Que hace |
|------|-------|---------|----------|
| 1 | `brainstorming` | Usuario pide feature/cambio | Refina ideas con preguntas, explora alternativas, presenta diseño por secciones para validar. Guarda documento de diseño. |
| 2 | `writing-plans` | Diseño aprobado | Divide el trabajo en tareas pequeñas (2-5 min cada una). Cada tarea tiene: rutas exactas, código completo, pasos de verificación. |
| 3 | `subagent-driven-development` o `executing-plans` | Plan listo | Despacha subagente fresco por tarea con revisión de dos etapas (spec + calidad), o ejecuta en batches con checkpoints humanos. |
| 4 | `test-driven-development` | Durante implementación | RED-GREEN-REFACTOR: escribir test que falla → verlo fallar → código mínimo → verlo pasar → commit. Borra código escrito antes de tests. |
| 5 | `requesting-code-review` | Entre tareas | Revisa contra el plan, reporta issues por severidad. Issues críticos bloquean progreso. |
| 6 | `finishing-a-development-branch` | Tareas completas | Verifica tests, presenta opciones (merge/PR/keep/discard), limpia worktree. |

**IMPORTANTE:** No saltear pasos. Si el usuario dice "agregá X", empezar con `brainstorming`, NO con código.

## Guias por Capa

| Capa                    | Ubicacion                          | AGENTS.md contiene                           |
| ----------------------- | ---------------------------------- | -------------------------------------------- |
| Domain                  | `internal/domain/`                 | Orquestador — apunta a entity/, vo/, service/ |
| Domain — Entity         | `internal/domain/entity/`          | Entidades, TipoEquipo, TipoCanalizacion, MemoriaCalculo |
| Domain — Value Objects  | `internal/domain/valueobject/`     | Corriente, Tension, Conductor, MaterialConductor |
| Domain — Services       | `internal/domain/service/`         | 6 servicios de calculo, formula IEEE-141 caida tension |
| Application             | `internal/application/`            | Ports, use cases, DTOs, orquestacion         |
| Infrastructure          | `internal/infrastructure/`         | Repos, CSV, PostgreSQL, mapeos, entorno      |
| Presentation            | `internal/presentation/`           | API REST, handlers, errores HTTP, versionado |
| Datos NOM               | `data/tablas_nom/`                 | Tablas CSV, formatos, reglas de validacion   |

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

### Skills de Proyecto

| Skill               | Descripcion                         | Ruta                                                  |
| ------------------- | ----------------------------------- | ----------------------------------------------------- |
| `agents-md-manager` | Crear y auditar jerarquia AGENTS.md | [SKILL.md](.agents/skills/agents-md-manager/SKILL.md) |

## Auto-invocacion

Cuando realices estas acciones, LEE el AGENTS.md o skill correspondiente PRIMERO:

| Accion                                    | Referencia                          |
| ----------------------------------------- | ----------------------------------- |
| Crear/modificar entidad o value object    | `internal/domain/entity/AGENTS.md` o `internal/domain/valueobject/AGENTS.md` |
| Crear/modificar servicio de calculo       | `internal/domain/service/AGENTS.md`  |
| Trabajar con ports o use cases            | `internal/application/AGENTS.md`    |
| Trabajar con DTOs o flujo de orquestacion | `internal/application/AGENTS.md`    |
| Modificar repositorios o CSV reader       | `internal/infrastructure/AGENTS.md` |
| Configurar BD o variables de entorno      | `internal/infrastructure/AGENTS.md` |
| Crear/modificar endpoints API             | `internal/presentation/AGENTS.md`   |
| Trabajar con tablas NOM CSV               | `data/tablas_nom/AGENTS.md`         |
| Agregar nueva tabla NOM                   | `data/tablas_nom/AGENTS.md`         |
| Aplicar patrones Go idiomaticos           | skill `golang-patterns`             |
| Crear/auditar AGENTS.md                   | skill `agents-md-manager`           |
| Disenar API endpoints                     | skill `api-design-principles`       |
| Crear nuevo skill                         | skill `skill-creator`               |
| Sincronizar skills a AGENTS.md            | skill `skill-sync`                  |
| Hacer commits o pull requests             | skill `commit-work`                 |

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

## Fases

1. **Fase 1 (actual):** 4 equipos, 6 servicios, 7 tablas NOM, 6 canalizaciones
2. **Fase 2:** Mas equipos y tablas NOM
3. **Fase 3:** PDF + frontend (repo separado)

**IMPORTANTE:** No adelantarse. Solo implementar lo necesario para la fase actual.

## Convenciones Globales

- **Nombres de negocio en espanol** (`MemoriaCalculo`, `CorrienteNominal`)
- **Codigo Go en ingles idiomatico** (packages, variables internas)
- **Errores:** `ErrXxx = errors.New(...)`, wrap con `fmt.Errorf("%w: ...", ErrXxx)`
- **Tests:** table-driven con `t.Run`, testify, `_test.go` en mismo directorio
- **Sin panic**, sin context en structs, receptores consistentes
- **Domain sin dependencias externas** — sin Gin, sin pgx, sin CSV; ports en `application/port/`; DI manual en `cmd/api/main.go`
- **YAGNI** — Fase 1 unicamente: sin PDF, sin auth, sin frontend, sin cache avanzado

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

### Pendientes (en `docs/plans/`)

- Canalizacion multi-tubo: `2026-02-13-canalizacion-multi-tubo-plan.md`
- Fase 2 memoria calculo: `2026-02-13-fase2-memoria-calculo-design.md`
