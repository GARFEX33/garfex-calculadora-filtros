# Calculos — Domain Layer

Capa de negocio pura para la feature de cálculos eléctricos. Sin dependencias externas (sin Gin, pgx, CSV).

> **Workflow:** Ver [docs/architecture/agents.md](../../../docs/architecture/agents.md)

## Estructura

| Subdirectorio | Contenido                                            |
| ------------- | ---------------------------------------------------- |
| `entity/`     | Entidades, tipos, interfaces del dominio de cálculos |
| `service/`    | Servicios de cálculo puros (sin I/O)                 |

## Dependencias permitidas

- `internal/shared/kernel/valueobject` — value objects compartidos
- stdlib de Go

## Dependencias prohibidas

> Ver reglas consolidadas en [docs/reference/structure.md](../../../docs/reference/structure.md)

- `internal/calculos/application/`
- `internal/calculos/infrastructure/`
- Gin, pgx, encoding/csv, cualquier framework externo

## Cómo modificar esta capa

> Ver flujo completo en [docs/architecture/workflow.md](../../../docs/architecture/workflow.md)

## Guías

| Subdirectorio | Contenido                                                                       |
| ------------- | ------------------------------------------------------------------------------- |
| `entity/`     | Entidades: TipoEquipo, TipoCanalizacion, SistemaElectrico, MemoriaCalculo, etc. |
| `service/`    | Servicios de cálculo NOM (ls internal/calculos/domain/service/*.go)            |

> **Nota:** Las subcarpetas `entity/` y `service/` heredan las reglas de este AGENTS.md. No necesitan AGENTS.md propio.

## Referencias

- Agente: `domain-agent`
- Skill: [.agents/skills/orchestrating-agents/SKILL.md](../../.agents/skills/orchestrating-agents/SKILL.md)

## Reglas de Oro

1. Domain nunca depende de Application ni Infrastructure
2. Sin I/O (no leer archivos, no HTTP, no DB)
3. Puro Go + lógica de negocio
4. Todo cambio pasa por `domain-agent`
