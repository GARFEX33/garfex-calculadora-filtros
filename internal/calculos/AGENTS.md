---
name: feature-calculos
description: Feature completa de memoria de cálculo eléctrico según normativa NOM.
model: opencode/minimax-m2.5-free
---

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
│           └── service/     ← 14 servicios de cálculo NOM
├── application/     ← ports, use cases, DTOs
│   ├── port/        ← TablaNOMRepository, EquipoRepository
│   ├── usecase/     ← OrquestadorMemoriaCalculo y micro use cases
│   └── dto/         ← EquipoInput, MemoriaOutput
└── infrastructure/  ← adapters HTTP (driver) y CSV (driven)
    └── adapter/
        ├── driver/http/     ← CalculoHandler, formatters
        └── driven/csv/      ← CSVTablaNOMRepository
```

## Reglas de Aislamiento

- **NO importa** `equipos/` ni ninguna otra feature
- **Solo importa** `shared/kernel/` para value objects compartidos
- **`cmd/api/main.go`** es el único que instancia y conecta esta feature

## Cómo modificar esta feature

**NUNCA modificar directamente.** Usar el sistema de orquestación.

> Ver flujo completo en [`AGENTS.md` raíz](../../../AGENTS.md)

## Referencias

- Workflow: [`AGENTS.md` raíz](../../../AGENTS.md) → "Sistema de Agentes Especializados"
- Skill: [`orchestrating-agents`](../../.agents/skills/orchestrating-agents/SKILL.md)
