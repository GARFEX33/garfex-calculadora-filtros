# Garfex Calculadora Filtros

Backend API en Go para memorias de calculo de instalaciones electricas segun normativa NOM (Mexico).

## Como Usar Esta Guia

- Empieza aqui para normas globales del proyecto
- Cada capa tiene su propio CLAUDE.md con guias especificas
- El CLAUDE.md del directorio tiene precedencia sobre este archivo cuando hay conflicto

## Guias por Capa

| Capa | Ubicacion | CLAUDE.md contiene |
|------|-----------|-------------------|
| Domain | `internal/domain/` | Entidades, VOs, servicios, formulas NOM |
| Application | `internal/application/` | Ports, use cases, DTOs, orquestacion |
| Infrastructure | `internal/infrastructure/` | Repos, CSV, PostgreSQL, mapeos, entorno |
| Presentation | `internal/presentation/` | API REST, handlers, errores HTTP, versionado |
| Datos NOM | `data/tablas_nom/` | Tablas CSV, formatos, reglas de validacion |

## Skills Disponibles

### Skills Genericos

| Skill | Descripcion | Ruta |
|-------|-------------|------|
| `golang-patterns` | Patrones Go idiomaticos, error handling, interfaces | [SKILL.md](.agents/skills/golang-patterns/SKILL.md) |
| `golang-pro` | Go avanzado: concurrencia, microservicios, performance | [SKILL.md](.agents/skills/golang-pro/SKILL.md) |
| `api-design-principles` | Diseno REST/GraphQL, convenciones API | [SKILL.md](.agents/skills/api-design-principles/SKILL.md) |
| `skill-creator` | Crear nuevos skills siguiendo el spec de Agent Skills | [SKILL.md](.agents/skills/skill-creator/SKILL.md) |
| `skill-sync` | Sincronizar metadata de skills a tablas Auto-invocacion | [SKILL.md](.agents/skills/skill-sync/SKILL.md) |

### Skills de Proyecto

| Skill | Descripcion | Ruta |
|-------|-------------|------|
| `claude-md-manager` | Crear y auditar jerarquia CLAUDE.md | [SKILL.md](.agents/skills/claude-md-manager/SKILL.md) |

## Auto-invocacion

Cuando realices estas acciones, LEE el CLAUDE.md o skill correspondiente PRIMERO:

| Accion | Referencia |
|--------|-----------|
| Crear/modificar entidad o value object | `internal/domain/CLAUDE.md` |
| Crear/modificar servicio de calculo | `internal/domain/CLAUDE.md` |
| Trabajar con ports o use cases | `internal/application/CLAUDE.md` |
| Trabajar con DTOs o flujo de orquestacion | `internal/application/CLAUDE.md` |
| Modificar repositorios o CSV reader | `internal/infrastructure/CLAUDE.md` |
| Configurar BD o variables de entorno | `internal/infrastructure/CLAUDE.md` |
| Crear/modificar endpoints API | `internal/presentation/CLAUDE.md` |
| Trabajar con tablas NOM CSV | `data/tablas_nom/CLAUDE.md` |
| Agregar nueva tabla NOM | `data/tablas_nom/CLAUDE.md` |
| Aplicar patrones Go idiomaticos | skill `golang-patterns` |
| Crear/auditar CLAUDE.md | skill `claude-md-manager` |
| Disenar API endpoints | skill `api-design-principles` |
| Crear nuevo skill | skill `skill-creator` |
| Sincronizar skills a CLAUDE.md | skill `skill-sync` |

## Stack

Go 1.22+, Gin, PostgreSQL (pgx/v5), testify, golangci-lint

## Principios de Diseno

1. **Domain sin dependencias externas** — sin Gin, sin pgx, sin CSV
2. **Domain no conoce NOM como archivos** — recibe datos ya interpretados
3. **Accept interfaces, return structs**
4. **Interfaces definidas donde se consumen** — ports en `application/port/`
5. **Inyeccion de dependencias manual** — en `cmd/api/main.go`
6. **YAGNI** — solo lo necesario para la fase actual

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

## Non-Goals (Fase 1)

- No PDF, no autenticacion, no multi-tenant
- No cache, no pooling avanzado
- No frontend (repositorio separado)

## Convenciones Globales

- **Nombres de negocio en espanol** (`MemoriaCalculo`, `CorrienteNominal`)
- **Codigo Go en ingles idiomatico** (packages, variables internas)
- **Errores:** `ErrXxx = errors.New(...)`, wrap con `fmt.Errorf("%w: ...", ErrXxx)`
- **Tests:** table-driven con `t.Run`, testify, `_test.go` en mismo directorio
- **Sin panic**, sin context en structs, receptores consistentes

## Actualizacion de Documentacion

Al terminar cada tarea, actualizar: plan si diverge, CLAUDE.md si cambia una regla, MEMORY.md si debe persistir entre sesiones.

## Documentacion

- Arquitectura: `docs/plans/2026-02-09-arquitectura-inicial-design.md`
- Domain layer: `docs/plans/2026-02-10-domain-layer.md`
- Caida de tension: `docs/plans/2026-02-11-caida-tension-impedancia-design.md`
- Canalizacion: `docs/plans/2026-02-11-tablas-nom-canalizacion-design.md`
