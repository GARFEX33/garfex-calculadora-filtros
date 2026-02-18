---
name: domain-agent
description: Especialista únicamente en la capa de dominio de calculos. Entidades, value objects y servicios de negocio.
model: opencode/minimax-m2.5-free
---

# Calculos — Domain Layer

Capa de negocio pura para la feature de cálculos eléctricos. Sin dependencias externas (sin Gin, pgx, CSV).

> **Workflow:** Ver [`AGENTS.md` raíz](../../../AGENTS.md) → "Sistema de Agentes Especializados"

## Estructura

| Subdirectorio | Contenido |
|---------------|-----------|
| `entity/` | Entidades, tipos, interfaces del dominio de cálculos |
| `service/` | Servicios de cálculo puros (sin I/O) |

## Dependencias permitidas

- `internal/shared/kernel/valueobject` — value objects compartidos
- stdlib de Go

## Dependencias prohibidas

- `internal/calculos/application/`
- `internal/calculos/infrastructure/`
- Gin, pgx, encoding/csv, cualquier framework externo

## Cómo modificar esta capa

> Ver flujo completo en [`AGENTS.md` raíz](../../../AGENTS.md)

## Guías

| Subdirectorio | Contenido |
|---------------|-----------|
| `entity/` | Entidades: TipoEquipo, TipoCanalizacion, SistemaElectrico, MemoriaCalculo, etc. |
| `service/` | 13 servicios de cálculo NOM + IEEE-141 |

## Referencias

- Agente: `domain-agent`
- Skill: `.agents/skills/orchestrating-agents/SKILL.md`

## Reglas de Oro

1. Domain nunca depende de Application ni Infrastructure
2. Sin I/O (no leer archivos, no HTTP, no DB)
3. Puro Go + lógica de negocio
4. Todo cambio pasa por `domain-agent`
