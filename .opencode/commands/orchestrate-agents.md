# Orchestrate Agents Command

## Description

Orquesta el trabajo de agentes especializados (domain-agent, application-agent, infrastructure-agent) siguiendo la arquitectura hexagonal con vertical slices.

## Usage

```bash
# Despachar domain-agent para pasos 1-2
orchestrate-agents --agent domain --steps "1-2" --plan docs/plans/2026-02-15-mi-plan.md

# Despachar application-agent para paso 3
orchestrate-agents --agent application --steps "3" --plan docs/plans/2026-02-15-mi-plan.md

# Despachar infrastructure-agent para paso 4
orchestrate-agents --agent infrastructure --steps "4" --plan docs/plans/2026-02-15-mi-plan.md
```

## Parameters

| Parameter | Description | Required |
|-----------|-------------|----------|
| `--agent` | Tipo de agente: `domain`, `application`, `infrastructure` | Yes |
| `--steps` | Pasos del plan a ejecutar (ej: "1-2", "3", "4-5") | Yes |
| `--plan` | Ruta al archivo de plan (desde raíz del proyecto) | Yes |
| `--feature` | Nombre de la feature (ej: `calculos`, `equipos`) | No, default: "nueva-feature" |

## Prerequisites

- Diseño aprobado en `docs/plans/YYYY-MM-DD-*-design.md`
- Plan de implementación en `docs/plans/YYYY-MM-DD-*-plan.md`
- Rama de trabajo creada y activa

## Agent Dispatch Flow

```
Coordinador (este chat)
    │
    ├──► domain-agent (pasos 1-2)
    │      └──► domain/ completo + testeado
    │
    ├──► application-agent (paso 3)
    │      └──► application/ completo + testeado
    │
    ├──► infrastructure-agent (paso 4)
    │      └──► infrastructure/ completo + testeado
    │
    └──► Coordinador: wiring + commit
```

## What Each Agent Does

### domain-agent
- **Scope:** `internal/shared/kernel/`, `internal/{feature}/domain/`
- **Creates:** Entities, Value Objects, Domain Services
- **Verifies:** `go test ./internal/{feature}/domain/...`

### application-agent
- **Scope:** `internal/{feature}/application/`
- **Creates:** Ports, Use Cases, DTOs
- **Verifies:** `go test ./internal/{feature}/application/...`

### infrastructure-agent
- **Scope:** `internal/{feature}/infrastructure/`
- **Creates:** Adapters (driver HTTP, driven repositories)
- **Verifies:** `go test ./internal/{feature}/infrastructure/...`

## Prompt Template Used

When you run this command, the following prompt template is sent to the agent:

---

Sos el **{agent}-agent** de este proyecto. Tu trabajo es ejecutar **{steps}** del plan.

## Proyecto

- Repositorio: `{absolute_path}`
- Rama activa: `{current_branch}`
- Módulo Go: `{go_module}`

## Contexto — qué hicieron los agentes anteriores

{context_from_previous_agents}

## Tu scope

{scope_directories}

**NO toches:**
{forbidden_directories}

## Plan a ejecutar

`{plan_path}`

## Instrucciones

1. Leé el plan completo
2. Creá tus propias tareas con TodoWrite antes de empezar
3. Ejecutá cada tarea marcando `in_progress` → `completed`
4. Verificá con `go test` antes de terminar cada paso
5. Si algo falla, arreglalo antes de seguir

## Al terminar

Reportá:
- Lista exacta de archivos creados/modificados
- Output de `go test ./...`
- Issues encontrados (si hay)

---

## Example Session

```
# Usuario solicita nueva feature
> Quiero agregar soporte para transformadores trifásicos

# Coordinador ejecuta workflow
1. skill: brainstorming → diseño aprobado
2. skill: writing-plans → plan creado
3. git checkout -b feature/transformador-trifasico

# Coordinador despacha agentes
4. orchestrate-agents --agent domain --steps "1-2" --plan docs/plans/2026-02-20-transformador-plan.md --feature calculos
   → Esperar a que domain-agent termine

5. orchestrate-agents --agent application --steps "3" --plan docs/plans/2026-02-20-transformador-plan.md --feature calculos
   → Esperar a que application-agent termine

6. orchestrate-agents --agent infrastructure --steps "4" --plan docs/plans/2026-02-20-transformador-plan.md --feature calculos
   → Esperar a que infrastructure-agent termine

# Coordinador finaliza
7. Actualizar cmd/api/main.go
8. git add -A && git commit -m "feat: add transformador trifasico support"
```

## Rules

1. **Wait for each agent** — don't dispatch next until previous reports completion
2. **One agent at a time** — no parallel execution
3. **Verify tests** — each agent must report `go test` output
4. **Respect scope** — agents only touch their assigned directories
5. **Coordinator only** — only this chat handles wiring and commits

## See Also

- `brainstorming` skill — step 1: design
- `writing-plans` skill — step 2: planning
- `.agents/skills/orchestrating-agents/SKILL.md` — complete workflow documentation
- `AGENTS.md` root — project-specific agent system rules
