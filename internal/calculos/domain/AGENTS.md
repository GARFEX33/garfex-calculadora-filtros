# Calculos — Domain Layer

Capa de negocio pura para la feature de cálculos eléctricos. Sin dependencias externas (sin Gin, pgx, CSV).

## Trabajar en esta Capa

Esta capa es responsabilidad del **`domain-agent`**. El agente ejecuta su ciclo completo:

```
brainstorming-dominio → writing-plans-dominio → executing-plans-dominio
```

**NO modificar directamente** — usar el sistema de orquestación.

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

### Opción A: Nueva feature (recomendado)

Si necesitás agregar/modificar dominio, crear una **nueva feature**:

```bash
# Orquestador (este chat):
orchestrate-agents --agent domain --feature nueva-feature
```

El `domain-agent` hará:
1. Brainstorming del dominio
2. Plan de implementación
3. Ejecución con tests

### Opción B: Cambio pequeño en calculos existente

Para cambios menores (ej: nueva validación, fix de bug):

```bash
# Orquestador:
# "domain-agent: agregar validación X a entidad Y"
```

El agente ejecuta solo la fase de executing-plans-dominio (si no hay diseño nuevo).

## Guías

| Subdirectorio | Contenido |
|---------------|-----------|
| `entity/` | Entidades: TipoEquipo, TipoCanalizacion, SistemaElectrico, MemoriaCalculo, etc. |
| `service/` | 12 servicios de cálculo NOM + IEEE-141 |

## Referencias

- Agente: `.opencode/agents/domain-agent.md`
- Comando: `.opencode/commands/orchestrate-agents.md`
- Skill: `.agents/skills/brainstorming-dominio/SKILL.md`

## Reglas de Oro

1. Domain nunca depende de Application ni Infrastructure
2. Sin I/O (no leer archivos, no HTTP, no DB)
3. Puro Go + lógica de negocio
4. Todo cambio pasa por `domain-agent`
