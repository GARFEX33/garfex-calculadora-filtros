# Garfex Calculadora Filtros

API Go (backend) + Frontend Svelte/TypeScript para memorias de calculo de instalaciones electricas segun normativa NOM (Mexico).

## Como Usar Esta Guia

- Empieza aqui para normas globales del proyecto
- Cada feature y capa tiene su propio AGENTS.md con guias especificas
- El AGENTS.md del directorio tiene precedencia sobre este archivo cuando hay conflicto

## Regla de Skills (OBLIGATORIO)

**ANTES de cualquier accion, verificar si aplica un skill.** Si hay 1% de probabilidad de que aplique, invocar el skill con la herramienta `Skill`.

Si el skill tiene checklist, crear todos con TodoWrite antes de seguirlo.

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

Cuando realices estas acciones, invoca el skill correspondiente PRIMERO:

| Acción | Skill |
| ------ | ----- |
| Crear/modificar archivos `.go` en `internal/` | `clean-ddd-hexagonal-vertical-go-enterprise` |
| Patrones idiomáticos de Go (interfaces, errores, concurrencia) | `golang-patterns` |
| Go avanzado (goroutines, generics, gRPC, microservices) | `golang-pro` |
| Diseñar o revisar endpoints REST / API contracts | `api-design-principles` |

### Frontend (Svelte / SvelteKit)

Cuando realices estas acciones, invoca el skill correspondiente PRIMERO:

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

## Skills Disponibles

### Skills de Backend (Go)

| Skill | Descripción |
| ----- | ----------- |
| `clean-ddd-hexagonal-vertical-go-enterprise` | Arquitectura hexagonal + DDD + vertical slices en Go |
| `golang-patterns` | Patrones idiomáticos de Go |
| `golang-pro` | Go avanzado: goroutines, generics, gRPC |
| `api-design-principles` | Principios de diseño REST y GraphQL |

### Skills de Frontend (Svelte)

| Skill | Descripción |
| ----- | ----------- |
| `svelte-code-writer` | CLI `@sveltejs/mcp` para docs y autofixer — OBLIGATORIO al tocar `.svelte` |
| `svelte5-best-practices` | Runes, snippets, eventos, TypeScript, migración Svelte 4→5 |
| `sveltekit-structure` | Routing, layouts, error handling, SSR, hidratación |
| `tailwind-design-system` | Design tokens, componentes UI, responsive, dark mode con Tailwind v4 |
| `typescript-advanced-types` | Generics, conditional/mapped types, type-safe API clients, utilidades TS |

### Skills Generales

| Skill | Descripción |
| ----- | ----------- |
| `finishing-a-development-branch` | Guía para finalizar e integrar ramas |
| `commit-work` | Commits convencionales de alta calidad |
| `systematic-debugging` | Debugging sistemático ante bugs o fallos |
| `skill-creator` | Crear nuevos skills para agentes |
| `skill-sync` | Sincronizar metadatos de skills a AGENTS.md |
| `verification-before-completion` | Verificar antes de declarar trabajo listo |
| `agents-md-manager` | Gestionar AGENTS.md y README.md jerárquicos |

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

| Tema           | Archivo                                              |
| -------------- | ---------------------------------------------------- |
| Frontend Web   | [frontend/web/AGENTS.md](frontend/web/AGENTS.md)     |
| Frontend Mobile | [frontend/mobile/AGENTS.md](frontend/mobile/AGENTS.md) — reservado, no activo |
