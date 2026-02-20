# Feature: calculos

Memoria de cálculo eléctrico según normativa NOM (México).

Esta feature implementa el cálculo completo de una instalación eléctrica:
corriente nominal → ajuste por temperatura/agrupamiento → selección de conductor →
conductor de tierra → dimensionamiento de canalización → caída de tensión (NOM).

## Estructura

```
internal/calculos/
├── domain/          ← entidades y servicios de cálculo puro
│   ├── entity/      ← Proyecto, TipoCanalizacion, SistemaElectrico, etc.
│   └── service/     ← Servicios de cálculo NOM
├── application/     ← ports, use cases, DTOs
│   ├── port/        ← TablaNOMRepository, EquipoRepository
│   ├── usecase/     ← OrquestadorMemoriaCalculo y micro use cases
│   │   └── helpers/ ← Funciones auxiliares
│   └── dto/         ← EquipoInput, MemoriaOutput
└── infrastructure/  ← adapters HTTP (driver) y CSV (driven)
    └── adapter/
        ├── driver/http/     ← CalculoHandler, formatters
        └── driven/csv/      ← CSVTablaNOMRepository
```

## Cómo modificar esta feature

**NUNCA modificar directamente.** Usar el sistema de orquestación.

> Ver flujo completo en [docs/architecture/workflow.md](../../../docs/architecture/workflow.md)

## Referencias

- Workflow: [docs/architecture/workflow.md](../../../docs/architecture/workflow.md)
- Sistema de agentes: [docs/architecture/agents.md](../../../docs/architecture/agents.md)
- Estructura y reglas: [docs/reference/structure.md](../../../docs/reference/structure.md)
- Skill: [orchestrating-agents](../../.agents/skills/orchestrating-agents/SKILL.md)
