# Garfex Calculadora Filtros

Backend API en Go para memorias de calculo de instalaciones electricas segun normativa NOM (Mexico).

## Como Usar Esta Guia

- Empieza aqui para normas globales del proyecto
- Cada feature y capa tiene su propio AGENTS.md con guias especificas
- El AGENTS.md del directorio tiene precedencia sobre este archivo cuando hay conflicto

## Regla de Skills (OBLIGATORIO)

**ANTES de cualquier accion, verificar si aplica un skill.** Si hay 1% de probabilidad de que aplique, invocar el skill con la herramienta `Skill`.

Si el skill tiene checklist, crear todos con TodoWrite antes de seguirlo.

## Documentación

| Tema                    | Archivo                                                        |
| ----------------------- | -------------------------------------------------------------- |
| Estructura del proyecto | [docs/reference/structure.md](docs/reference/structure.md)     |
| Skills disponibles      | [docs/reference/skills.md](docs/reference/skills.md)           |
| Auto-invocación         | [docs/reference/auto-invoke.md](docs/reference/auto-invoke.md) |
| Comandos y endpoints    | [docs/reference/commands.md](docs/reference/commands.md)       |

## Documentacion Implementada

| Tema                      | Archivo                                                                                  |
| ------------------------- | ---------------------------------------------------------------------------------------- |
| Feature Cálculos          | [internal/calculos/AGENTS.md](internal/calculos/AGENTS.md)                               |
| Feature Equipos           | [internal/equipos/AGENTS.md](internal/equipos/AGENTS.md)                                 |
| Kernel Compartido         | [internal/shared/kernel/AGENTS.md](internal/shared/kernel/AGENTS.md)                     |
| Tablas NOM (datos)        | [data/tablas_nom/AGENTS.md](data/tablas_nom/AGENTS.md)                                   |
| Cálculos — Domain         | [internal/calculos/domain/AGENTS.md](internal/calculos/domain/AGENTS.md)                 |
| Cálculos — Application    | [internal/calculos/application/AGENTS.md](internal/calculos/application/AGENTS.md)       |
| Cálculos — Infrastructure | [internal/calculos/infrastructure/AGENTS.md](internal/calculos/infrastructure/AGENTS.md) |
| Equipos — Domain          | [internal/equipos/domain/AGENTS.md](internal/equipos/domain/AGENTS.md)                   |
| Equipos — Application     | [internal/equipos/application/AGENTS.md](internal/equipos/application/AGENTS.md)         |
| Equipos — Infrastructure  | [internal/equipos/infrastructure/AGENTS.md](internal/equipos/infrastructure/AGENTS.md)   |
