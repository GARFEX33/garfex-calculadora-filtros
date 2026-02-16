# {Nombre del Proyecto}

{Descripción en una línea del proyecto.}

## Cómo Usar Esta Guía

- Empezá aquí para normas globales del proyecto
- Cada feature y capa tiene su propio AGENTS.md con guías específicas
- El AGENTS.md del directorio tiene precedencia sobre este archivo cuando hay conflicto

## Regla Anti-Duplicación (OBLIGATORIO) — RESPONSABILIDAD DEL ORQUESTADOR

⚠️ **Los agentes especializados NO se conocen entre sí.** El orquestador es el único con visión global.

**Flujo del Orquestador (antes de despachar agentes):**

**Paso 1: Investigar**
```bash
ls internal/{feature}/domain/service/*.go 2>/dev/null
rg "TODO|FIXME|XXX" internal/{feature}/application/usecase --type go
rg -i "func.*[Cc]alcular" internal/{feature} --type go
```

**Paso 2: Decidir**
| Situación | Decisión |
|-----------|----------|
| Existe servicio similar | Extender, no crear nuevo |
| Use case tiene TODO | Implementar TODO primero |
| Nada similar | Crear nuevo |

**Paso 3: Comunicar (en el prompt al agente)**

❌ Mal: "Creá un servicio para X"

✅ Bien: "Implementá el método Y() que tiene un TODO en ZUseCase. 
          Usá el servicio W que ya existe. NO crees un use case nuevo."

**Checklist (orquestador):**
- [ ] ¿Investigué qué ya existe?
- [ ] ¿Tomé la decisión de extender vs crear?
- [ ] ¿Comuniqué claramente qué hacer y qué NO hacer?

## Sistema de Agentes Especializados

Cada capa tiene su agente. El orquestador delega en orden:

```
domain-agent → application-agent → infrastructure-agent
```

### Cuándo invocar cada agente

| Acción | Agente | Skills |
|--------|--------|--------|
| Entidades, VOs, servicios de dominio | `domain-agent` | `brainstorming-dominio` → `writing-plans-dominio` → `executing-plans-dominio` |
| Ports, use cases, DTOs | `application-agent` | `brainstorming-application` → `writing-plans-application` → `executing-plans-application` |
| Adapters, repos, handlers HTTP | `infrastructure-agent` | `brainstorming-infrastructure` → `writing-plans-infrastructure` → `executing-plans-infrastructure` |
| Actualizar/auditar AGENTS.md | `agents-md-curator` | `agents-md-manager` |

### Reglas de delegación

- **domain-agent** trabaja primero — no sabe que existen los otros
- **application-agent** toma el output del domain — no toca infra
- **infrastructure-agent** toma el output del application — implementa ports
- **Orquestador** es el único que conoce el orden y hace wiring
- Cada agente verifica con `go test` antes de entregar

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

### Skills Genéricos

| Skill | Descripción | Ruta |
|-------|-------------|------|
| `golang-patterns` | Patrones Go idiomáticos | [SKILL.md](.agents/skills/golang-patterns/SKILL.md) |
| `api-design-principles` | Diseño REST/GraphQL | [SKILL.md](.agents/skills/api-design-principles/SKILL.md) |
| `commit-work` | Commits de calidad | [SKILL.md](.agents/skills/commit-work/SKILL.md) |

### Skills de Proceso

| Skill | Descripción | Ruta |
|-------|-------------|------|
| `brainstorming` | Explorar ideas antes de implementar | [SKILL.md](.agents/skills/brainstorming/SKILL.md) |
| `brainstorming-dominio` | Diseñar dominio | [SKILL.md](.agents/skills/brainstorming-dominio/SKILL.md) |
| `brainstorming-application` | Diseñar application | [SKILL.md](.agents/skills/brainstorming-application/SKILL.md) |
| `brainstorming-infrastructure` | Diseñar infrastructure | [SKILL.md](.agents/skills/brainstorming-infrastructure/SKILL.md) |

### Skills de Proyecto

| Skill | Descripción | Ruta |
|-------|-------------|------|
| `agents-md-manager` | Gestionar AGENTS.md | [SKILL.md](.agents/skills/agents-md-manager/SKILL.md) |

## Auto-invocación

| Acción | Agente / Skill |
|--------|----------------|
| Crear/modificar entidad o VO | `domain-agent` |
| Crear/modificar servicio de cálculo | `domain-agent` |
| Trabajar con ports o use cases | `application-agent` |
| Crear/modificar endpoints HTTP | `infrastructure-agent` |
| Crear/auditar AGENTS.md | `agents-md-curator` → `agents-md-manager` |
| Aplicar patrones Go | skill `golang-patterns` |
| Hacer commits | skill `commit-work` |

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
