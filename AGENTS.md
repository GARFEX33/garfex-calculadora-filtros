# Garfex Calculadora Filtros

Backend API en Go para memorias de calculo de instalaciones electricas segun normativa NOM (Mexico).

## Como Usar Esta Guia

- Empieza aqui para normas globales del proyecto
- Cada capa tiene su propio AGENTS.md con guias especificas
- El AGENTS.md del directorio tiene precedencia sobre este archivo cuando hay conflicto

## Guias por Capa

| Capa                    | Ubicacion                          | AGENTS.md contiene                           |
| ----------------------- | ---------------------------------- | -------------------------------------------- |
| Domain                  | `internal/domain/`                 | Orquestador — apunta a entity/, vo/, service/ |
| Domain — Entity         | `internal/domain/entity/`          | Entidades, TipoEquipo, TipoCanalizacion, MemoriaCalculo |
| Domain — Value Objects  | `internal/domain/valueobject/`     | Corriente, Tension, Conductor                |
| Domain — Services       | `internal/domain/service/`         | 6 servicios de calculo, caida de tension     |
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

```bash
go test ./...           # Tests
go test -race ./...     # Tests con race detector
go build ./...          # Compilacion
go vet ./...            # Analisis estatico
golangci-lint run       # Linting completo
go run cmd/api/main.go  # Servidor dev
```

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

- Arquitectura: `docs/plans/2026-02-09-arquitectura-inicial-design.md`
- Domain layer: `docs/plans/2026-02-10-domain-layer.md`
- Caida de tension: `docs/plans/2026-02-11-caida-tension-impedancia-design.md`
- Canalizacion: `docs/plans/2026-02-11-tablas-nom-canalizacion-design.md`
