# Garfex Calculadora Filtros

API Go (backend) + Frontend Svelte/TypeScript para memorias de calculo de instalaciones electricas segun normativa NOM (Mexico).

## Como Usar Esta Guia

- Empieza aqui para normas globales del proyecto
- Cada feature y capa tiene su propio AGENTS.md con guias especificas
- El AGENTS.md del directorio tiene precedencia sobre este archivo cuando hay conflicto

## Regla de Skills (OBLIGATORIO)

**ANTES de cualquier accion, verificar si aplica un skill.** Si hay 1% de probabilidad de que aplique, invocar el skill con la herramienta `Skill`.

Si el skill tiene checklist, crear todos con TodoWrite antes de seguirlo.

> Catálogo completo de skills: [.agents/skills/INDEX.md](.agents/skills/INDEX.md)

## Memoria Persistente (Engram) — ÚNICA FUENTE DE VERDAD

> **REGLA ABSOLUTA**: Todo historial de cambios, decisiones y contexto del proyecto vive EXCLUSIVAMENTE en Engram.
> **NUNCA** crear ni usar carpetas `openspec/` ni archivos de specs en el repo. No existe openspec en este proyecto.

**SIEMPRE usar Engram** para guardar el historial de decisiones, implementaciones y aprendizajes.

### Por qué SOLO Engram

- Engram persiste entre sesiones sin contaminar el repositorio
- Los archivos openspec/ quedan en git y confunden a futuros agentes
- Búsqueda semántica — encuentra decisiones pasadas por contexto, no por nombre de archivo
- Un único lugar para buscar: no hay que revisar carpetas de specs dispersas

### Cuándo guardar en Engram (OBLIGATORIO)

| Cuándo | Qué guardar |
|--------|-------------|
| Feature implementada | Resumen de cambios, archivos modificados |
| Decisión de arquitectura | Por qué se eligió X sobre Y |
| Bug corregido | Qué era, por qué ocurría, cómo se solucionó |
| Sesión terminada | Resumen de lo hecho, siguiente paso |
| Change SDD completado | Proposal, specs, design, tasks, verify (todo en Engram) |
| Change SDD pendiente | Estado actual + tareas pendientes |

### Cómo guardar

```typescript
// Usar la herramienta engram_mem_save con:
title: "Breve descripción"
type: "architecture | bugfix | decision | pattern | config"
content: "**What**: ... **Where**: ... **Learned**: ..."
```

### Cómo recuperar contexto al iniciar sesión

```
1. mem_context() — ver sesiones recientes
2. mem_search("tema relevante") — buscar decisiones pasadas
3. mem_get_observation(id) — leer contenido completo si está truncado
```

## Estructura del Proyecto

```
garfex-calculadora-filtros/
├── cmd/                        # Entrypoints Go
├── internal/                   # Lógica de negocio (hexagonal)
│   ├── calculos/               # Feature: cálculos eléctricos
│   ├── equipos/                # Feature: catálogo de filtros
│   └── shared/kernel/          # Value objects compartidos
├── data/tablas_nom/            # Tablas NOM (datos estáticos)
├── frontend/
│   ├── web/                    # SvelteKit + TypeScript (PWA responsive)
│   └── mobile/                 # Reservado — no activo
└── docker-compose.yml
```

## Auto-invocación de Skills

### Backend (Go)

| Acción | Skill |
| ------ | ----- |
| Crear/modificar archivos `.go` en `internal/` | `clean-ddd-hexagonal-vertical-go-enterprise` |
| Patrones idiomáticos de Go (interfaces, errores, concurrencia) | `golang-patterns` |
| Go avanzado (goroutines, generics, gRPC, microservices) | `golang-pro` |
| Diseñar o revisar endpoints REST / API contracts | `api-design-principles` |

### Frontend (Svelte / SvelteKit)

| Acción | Skill |
| ------ | ----- |
| Crear o editar cualquier componente `.svelte` o módulo `.svelte.ts/.svelte.js` | `svelte-code-writer` |
| Usar runes (`$state`, `$derived`, `$effect`, `$props`), snippets, eventos o migrar de Svelte 4 | `svelte5-best-practices` |
| Definir rutas, layouts, error boundaries, SSR o hidratación en SvelteKit | `sveltekit-structure` |
| Crear componentes UI, design tokens, responsive layout o sistema de estilos con Tailwind | `tailwind-design-system` |
| Tipos genéricos, tipos condicionales, mapped types, type-safe API clients o utilidades TS | `typescript-advanced-types` |

### General

| Acción | Skill |
| ------ | ----- |
| Terminar una rama de desarrollo | `finishing-a-development-branch` |
| Crear un commit | `commit-work` |
| Encontrar y corregir un bug o fallo de tests | `systematic-debugging` |
| Crear/modificar un skill | `skill-creator` |
| Regenerar tablas auto-invoke en AGENTS.md | `skill-sync` |
| Verificar antes de declarar trabajo completo | `verification-before-completion` |
| Crear o auditar AGENTS.md / README.md | `agents-md-manager` |

### SDD (Spec-Driven Development)

> **IMPORTANTE**: En este proyecto los artefactos SDD (proposal, specs, design, tasks, verify) se guardan en **Engram**, NO en carpetas `openspec/`. El modo de artefactos es `engram`.

| Acción | Skill |
| ------ | ----- |
| Inicializar SDD en el proyecto ("sdd init", "iniciar sdd") | `sdd-init` |
| Escribir specs de un cambio | `sdd-spec` |
| Crear diseño técnico | `sdd-design` |
| Investigar viabilidad y contexto | `sdd-explore` |
| Crear propuesta de cambio | `sdd-propose` |
| Descomponer en tareas | `sdd-tasks` |
| Implementar según specs | `sdd-apply` |
| Validar implementación contra specs | `sdd-verify` |
| Archivar cambio completado | `sdd-archive` |

## Documentacion Implementada

### Backend

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

### Frontend

| Tema            | Archivo                                                                        |
| --------------- | ------------------------------------------------------------------------------ |
| Frontend Web    | [frontend/web/AGENTS.md](frontend/web/AGENTS.md)                               |
| Frontend Mobile | [frontend/mobile/AGENTS.md](frontend/mobile/AGENTS.md) — reservado, no activo |
