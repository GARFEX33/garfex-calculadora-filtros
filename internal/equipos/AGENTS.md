# Feature: equipos

Catálogo de equipos Garfex — PLACEHOLDER FUTURO.

Esta feature estará a cargo del catálogo de equipos eléctricos (filtros activos,
filtros de rechazo, transformadores, cargas) con su búsqueda y persistencia en PostgreSQL.

## Estado actual

**Estructura vacía.** Solo existe como placeholder para mantener los límites de la arquitectura.

## Implementar esta feature

Cuando se decida implementar, usar el sistema de orquestación:

```bash
# Paso 1: Domain
orchestrate-agents --agent domain --feature equipos

# Paso 2: Application
orchestrate-agents --agent application --feature equipos

# Paso 3: Infrastructure
orchestrate-agents --agent infrastructure --feature equipos

# Paso 4: Orquestador actualiza main.go
```

## Estructura esperada

```
internal/equipos/
├── domain/
│   ├── entity/      ← Equipo, FiltroActivo, FiltroRechazo, etc.
│   └── service/     ← Búsqueda, Validaciones
├── application/
│   ├── port/        ← EquipoRepository
│   ├── usecase/     ← BuscarEquipo, ListarEquipos
│   └── dto/         ← EquipoInput, EquipoOutput
└── infrastructure/
    └── adapter/
        ├── driver/http/      ← EquipoHandler
        └── driven/postgres/  ← PostgresEquipoRepository
```

## Referencias

- Workflow: [docs/architecture/workflow.md](../../../docs/architecture/workflow.md)
- Sistema de agentes: [docs/architecture/agents.md](../../../docs/architecture/agents.md)
- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)
- Orquestación: [orchestrating-agents](../../.agents/skills/orchestrating-agents/SKILL.md)
- Skill domain: [brainstorming-dominio](../../.agents/skills/brainstorming-dominio/SKILL.md)
- Skill application: [brainstorming-application](../../.agents/skills/brainstorming-application/SKILL.md)
- Skill infrastructure: [brainstorming-infrastructure](../../.agents/skills/brainstorming-infrastructure/SKILL.md)
