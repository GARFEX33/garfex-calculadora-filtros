# {Nombre del Proyecto}

{Descripción en una línea del proyecto.}

## Cómo Usar Esta Guía

- Empezá aquí para normas globales del proyecto
- Cada feature y capa tiene su propio AGENTS.md con guías específicas
- El AGENTS.md del directorio tiene precedencia sobre este archivo cuando hay conflicto

## Regla Anti-Duplicación (OBLIGATORIO)

Antes de crear algo nuevo, verificar si ya existe:

1. **Buscar servicios similares:**
   ```bash
   ls internal/{feature}/domain/service/*.go 2>/dev/null
   rg "TODO|FIXME|XXX" internal/{feature}/application/usecase --type go
   ```

2. **Si existe similar → extender, no crear nuevo**

3. **Si existe TODO → implementar el TODO primero**

## Estructura del Proyecto

```
internal/
  shared/kernel/valueobject/   ← VOs compartidos entre features
  {feature}/                   ← Una feature = un vertical slice
    domain/
      entity/
      service/
    application/
      port/
      usecase/
      dto/
    infrastructure/
      adapter/
        driver/{protocolo}/    ← HTTP, gRPC, CLI
        driven/{tecnología}/   ← Postgres, CSV, Redis
  cmd/api/main.go                ← Wiring de todas las features
data/                          ← Datos estáticos (CSV, etc.)
```

## Guías por Capa

| Capa | Ubicación | AGENTS.md |
|------|-----------|-----------|
| Shared Kernel | `internal/shared/kernel/` | `internal/shared/kernel/AGENTS.md` |
| Feature {name} | `internal/{feature}/` | Ver subcapas |
| Domain | `internal/{feature}/domain/` | `internal/{feature}/domain/AGENTS.md` |
| Application | `internal/{feature}/application/` | `internal/{feature}/application/AGENTS.md` |
| Infrastructure | `internal/{feature}/infrastructure/` | `internal/{feature}/infrastructure/AGENTS.md` |

## Skills Disponibles

| Skill | Descripción | Ruta |
|-------|-------------|------|
| `golang-patterns` | Patrones Go idiomáticos | [SKILL.md](.agents/skills/golang-patterns/SKILL.md) |
| `api-design-principles` | Diseño REST/GraphQL | [SKILL.md](.agents/skills/api-design-principles/SKILL.md) |
| `commit-work` | Commits de calidad | [SKILL.md](.agents/skills/commit-work/SKILL.md) |
| `agents-md-manager` | Gestionar AGENTS.md | [SKILL.md](.agents/skills/agents-md-manager/SKILL.md) |

## Stack

{Lenguaje} {versión}, {Framework}, {DB}, {Testing}, {Linting}

## Comandos

```bash
# Tests
go test ./...
go test -race ./...

# Build
go build ./...

# Lint
go vet ./...
golangci-lint run
```

## Documentación

### Planes completados
- `docs/plans/completed/` — Diseños e implementaciones pasadas

### Planes en progreso
- `docs/plans/` — Trabajo actual

## Convenciones Globales

- **Nombres de negocio en español** — términos del dominio
- **Código en inglés idiomático** — variables, funciones, paquetes
- **Errores:** wrap con contexto, nunca silenciar
- **Tests:** table-driven con subtests
- **No panic** en código de dominio/librería
- **YAGNI** — solo fase actual, no especular
