# Frontend Web — SvelteKit + TypeScript

App web responsive (PWA) para la calculadora de filtros eléctricos. Consume la API Go en `internal/`.

## Skills OBLIGATORIOS (cargar PRIMERO)

| Acción                                                                         | Skill                            |
| ------------------------------------------------------------------------------ | -------------------------------- |
| Crear o editar cualquier archivo `.svelte`, `.svelte.ts`, `.svelte.js`         | `svelte-code-writer`             |
| Usar runes (`$state`, `$derived`, `$effect`, `$props`, `$bindable`) o snippets | `svelte5-best-practices`         |
| Definir rutas, layouts, `+page.svelte`, `+layout.svelte`, `+server.ts`, SSR    | `sveltekit-structure`            |
| Crear componentes UI, design tokens, responsive layout, dark mode              | `tailwind-design-system`         |
| Tipos genéricos, tipos condicionales, type-safe API clients, utilidades TS     | `typescript-advanced-types`      |
| Detectar o corregir bugs en componentes                                        | `systematic-debugging`           |
| Verificar antes de declarar trabajo listo                                      | `verification-before-completion` |

> **Nota Tailwind**: el skill usa ejemplos React — adaptar CVA/componentes al patrón Svelte 5.
> En Svelte se usa `class:` y props en lugar de `VariantProps`. Ver `svelte5-best-practices` para integración.

## Stack Tecnológico

| Tecnología                | Versión | Propósito                                   |
| ------------------------- | ------- | ------------------------------------------- |
| SvelteKit                 | 2.x     | Framework full-stack SSR/SPA                |
| Svelte                    | 5.x     | UI reactivo con runes                       |
| TypeScript                | 5.x     | Tipado estricto (`strict: true` + extras)   |
| Tailwind CSS              | 4.x     | Estilos CSS-first con `@theme`              |
| `@tailwindcss/vite`       | 4.x     | Plugin Vite para Tailwind v4                |
| `clsx` + `tailwind-merge` | latest  | Helper `cn()` para clases dinámicas         |
| ESLint                    | 10.x    | Linting — flat config (`eslint.config.js`)  |
| Prettier                  | 3.x     | Formateo — con plugins svelte + tailwind    |
| Husky                     | 9.x     | Git hooks (pre-commit)                      |
| lint-staged               | 16.x    | Ejecuta lint+format solo en archivos staged |
| Vite                      | 7.x     | Bundler                                     |

## Estructura

```
frontend/web/
├── src/
│   ├── app.css                    # CSS global + design tokens (@theme)
│   ├── app.html                   # HTML shell (lang="es", meta PWA)
│   ├── app.d.ts                   # Tipos globales SvelteKit
│   ├── lib/
│   │   ├── index.ts               # Re-exporta utils y tipos principales
│   │   ├── api/
│   │   │   ├── client.ts          # Cliente HTTP base (fetch wrapper)
│   │   │   ├── calculos.ts        # Endpoints de cálculos (memoria, amperaje, etc.)
│   │   │   └── equipos.ts         # Endpoints de equipos (CRUD catálogo)
│   │   ├── components/
│   │   │   ├── ui/                # Componentes primitivos reutilizables
│   │   │   └── calculos/          # Componentes de la calculadora eléctrica
│   │   │       ├── CamposInstalacion.svelte
│   │   │       ├── FormularioListado.svelte
│   │   │       ├── FormularioManual.svelte
│   │   │       ├── MemoriaTecnica.svelte
│   │   │       ├── ResultadosMemoria.svelte
│   │   │       ├── SelectorModo.svelte
│   │   │       └── secciones/     # Secciones de la memoria de cálculo
│   │   │           ├── SeccionAlimentador.svelte
│   │   │           ├── SeccionCaidaTension.svelte
│   │   │           ├── SeccionCanalizacion.svelte
│   │   │           ├── SeccionConclusion.svelte
│   │   │           ├── SeccionCorriente.svelte
│   │   │           ├── SeccionEncabezado.svelte
│   │   │           └── SeccionTierra.svelte
│   │   ├── types/
│   │   │   ├── index.ts           # Re-exporta todos los tipos
│   │   │   └── api.types.ts       # ApiResponse, ApiError, ApiResult, etc.
│   │   └── utils/
│   │       ├── index.ts           # Re-exporta utilidades
│   │       └── cn.ts              # Helper cn() para combinar clases
│   └── routes/
│       ├── +layout.svelte         # Layout raíz (importa app.css)
│       ├── +page.svelte           # Home /
│       ├── calculos/              # Rutas de cálculo eléctrico
│       │   └── resultado/         # /calculos/resultado — memoria de cálculo completa
│       └── equipos/               # Rutas de catálogo de filtros
├── static/                        # Assets estáticos públicos
├── .env                           # Variables de entorno locales (no commit)
├── .env.example                   # Template de variables (sí commit)
├── .prettierrc                    # Prettier: tabs, singleQuote, plugins svelte+tailwind
├── .prettierignore
├── eslint.config.js               # ESLint v9 flat config: JS + TS + Svelte
├── package.json                   # Scripts: dev, build, check, lint, format, qa
├── svelte.config.js
├── tsconfig.json                  # TypeScript strict + noUncheckedIndexedAccess + extras
└── vite.config.ts                 # Tailwind v4 via @tailwindcss/vite

# En el root del repo (../../):
# .husky/pre-commit               # Hook: cd frontend/web && npx lint-staged
```

## Arquitectura: Vertical Slicing

Este proyecto usa **Vertical Slicing** para organizar la lógica de negocio en lugar del patrón tradicional de capas horizontales.

### Por qué Vertical Slicing

- **Acoplamiento débil**: Cada feature es independiente y puede evolucionar sin afectar otras
- **Testabilidad**: Cada capa de una feature se puede testear de forma aislada
- **SSR-safe**: Los stores usan `$state` de Svelte 5, seguros para SSR
- **Preparación para OpenAPI**: Los tipos/domain están separados de la infraestructura, facilitando la generación automática将来

### Estructura de directorios features/

```
src/lib/features/
├── calculos/                    # Feature: cálculos eléctricos
│   ├── domain/                   # Tipos, validadores, reglas de negocio
│   │   ├── types/               # Tipos del dominio (sin dependencias externas)
│   │   │   ├── memoria.types.ts
│   │   │   ├── calculo.enums.ts
│   │   │   └── index.ts
│   │   └── validators/           # Validadores puros (sin side effects)
│   │       └── validar-memoria-request.ts
│   ├── infrastructure/            # Implementación técnica
│   │   ├── api/                 # Cliente HTTP (usa apiClient)
│   │   └── mappers/             # Transformación API ↔ Domain
│   └── application/              # Orquestación (stores, servicios)
│       └── stores/               # Svelte 5 stores con $state
│           └── memoria.store.svelte.ts
└── equipos/                      # Feature: catálogo de filtros
    ├── domain/
    ├── infrastructure/
    └── application/
```

### Reglas de dependencia

```
domain/    → SIN dependencias (solo tipos TS puros)
infrastructure/ → domain/ (usa los tipos del dominio)
application/   → domain/ + infrastructure/ (combina lógica + datos)
components/    → application/ (usa stores, NO lógica de negocio)
```

**Regla de oro**: Los componentes UI NUNCA deben tener lógica de negocio. Delegar siempre a los stores.

### Capas explicadas

#### domain/
- **Qué va aquí**: Tipos, enums, validadores puros
- **Qué NO va aquí**: `fetch`, `console.log`, referencias a `window`, estados reactivos
- **Ejemplo**:
  ```typescript
  // ✓ Correcto: tipo puro
  export interface EquipoFiltro { id: string; clave: string; }
  
  // ✗ Incorrecto:依赖 externa
  export interface EquipoConFetch extends EquipoFiltro { fetch: () => Promise }
  ```

#### infrastructure/
- **Qué va aquí**: Llamadas API, transformaciones de datos, configuraciones externas
- **Qué NO va aquí**: Lógica de negocio, estados reactivos
- **Ejemplo**:
  ```typescript
  // Mapea API → Domain
  export function mapApiToEquipo(api: ApiEquipo): EquipoFiltro { ... }
  ```

#### application/
- **Qué va aquí**: Stores con `$state`, servicios que orquestan múltiples llamadas
- **Qué NO va aquí**: Tipos definidos (usar los del domain)
- **Ejemplo**:
  ```typescript
  class EquiposStore {
    equipos = $state<EquipoFiltro[]>([]);
    loading = $state(false);
    async cargar() { ... }
  }
  ```

### Cómo agregar una nueva feature

1. Crear directorio en `src/lib/features/{nombre-feature}/`
2. Crear las 3 capas: `domain/`, `infrastructure/`, `application/`
3. Definir tipos en `domain/types/`
4. Implementar API en `infrastructure/api/`
5. Crear store en `application/stores/`
6. Crear `index.ts` que re-exporte todo
7. Consumir desde componentes en `src/lib/components/`

### Auto-invocación de Skills

| Contexto                          | Skill                                |
| --------------------------------- | ------------------------------------ |
| Crear/modificar .svelte           | `svelte5-best-practices`             |
| Crear store con $state            | `svelte5-best-practices`             |
| Tipos domain, validadores puros   | `typescript-advanced-types`          |
| Client API, mapeos                | `api-design-principles`              |
| Componente con diseño responsivo  | `tailwind-design-system`             |

### Compatibilidad con re-exports

Los archivos legacy en `lib/api/` y `lib/types/` re-exportan desde `features/` para no romper imports existentes:

```typescript
// Legacy import (deprecated pero funciona)
import { EquipoFiltro } from '$lib/types/equipos.types';

// Nuevo import (recomendado)
import { EquipoFiltro } from '$lib/features/equipos';
```

### Preparación para OpenAPI/codegen

La separación stricta entre domain e infrastructure permite:

1. **Fase actual**: Escribir tipos manualmente en `domain/types/`
2. **Fase futura**: Generar tipos desde OpenAPI y colocarlos en `infrastructure/api/generated/`
3. **Mappers**: Mantienen la transformación sin cambiar el domain

Esto permite migrar gradualmente sin reescribir toda la aplicación.

## Design System — Tokens Tailwind v4

El archivo `src/app.css` define TODOS los tokens en `@theme`. **Nunca usar valores arbitrarios**.

### Categorías de tokens

| Categoría          | Prefijo       | Ejemplo                               |
| ------------------ | ------------- | ------------------------------------- |
| Colores semánticos | `--color-*`   | `bg-primary`, `text-muted-foreground` |
| Border radius      | `--radius-*`  | `rounded-md`, `rounded-xl`            |
| Tipografía         | `--font-*`    | `font-sans`, `font-mono`              |
| Animaciones        | `--animate-*` | `animate-fade-in`, `animate-slide-up` |

### Colores clave

| Token                            | Uso                               |
| -------------------------------- | --------------------------------- |
| `background` / `foreground`      | Fondo y texto de la app           |
| `primary` / `primary-foreground` | Acciones principales, botones CTA |
| `muted` / `muted-foreground`     | Textos secundarios, placeholders  |
| `card` / `card-foreground`       | Tarjetas y paneles                |
| `sidebar` / `sidebar-foreground` | Navegación lateral                |
| `destructive`                    | Errores, eliminar                 |
| `success`                        | Éxito, validaciones OK            |
| `warning`                        | Advertencias                      |
| `border` / `ring`                | Bordes e indicadores de foco      |

### Dark mode

Se activa agregando clase `.dark` al `<html>`. Los tokens se redefinen automáticamente.

## Calidad de Código

### TypeScript — `tsconfig.json`

Opciones activas más allá de `strict: true`:

| Opción                               | Efecto                                                                   |
| ------------------------------------ | ------------------------------------------------------------------------ |
| `strict: true`                       | Activa: `strictNullChecks`, `noImplicitAny`, `strictFunctionTypes`, etc. |
| `noUncheckedIndexedAccess`           | `arr[0]` es `T \| undefined`, no `T`                                     |
| `noImplicitOverride`                 | Requires `override` keyword on overriding methods                        |
| `noPropertyAccessFromIndexSignature` | Fuerza bracket notation en index signatures                              |
| `exactOptionalPropertyTypes`         | `prop?: T` ≠ `prop: T \| undefined`                                      |

### ESLint — `eslint.config.js` (flat config v9)

- `@typescript-eslint/no-explicit-any`: **error** (prohibido `any`)
- `@typescript-eslint/consistent-type-imports`: **error** (usar `import type`)
- `@typescript-eslint/no-floating-promises`: **error** (promesas sin `await` o `.catch`)
- `@typescript-eslint/no-misused-promises`: **error**
- `no-console`: **warn** (permitido `console.warn` y `console.error`)
- `prefer-const`: **error** en `.ts`, off en `.svelte` (por `$props()` destructuring)

### Prettier — `.prettierrc`

- Tabs: **true** (indentación con tabs)
- `singleQuote`: **true**
- `trailingComma`: **none**
- `printWidth`: **100**
- `plugins`: `prettier-plugin-svelte` + `prettier-plugin-tailwindcss`
- `tailwindStylesheet`: apunta a `src/app.css` para ordenar clases con tokens propios

### Husky + lint-staged

El pre-commit hook en `../../.husky/pre-commit` ejecuta **lint-staged** automáticamente:

| Archivos staged           | Acciones                                       |
| ------------------------- | ---------------------------------------------- |
| `*.ts`, `*.js`            | `prettier --write` → `eslint --max-warnings=0` |
| `*.svelte`                | `prettier --write` → `eslint --max-warnings=0` |
| `*.css`, `*.json`, `*.md` | `prettier --write`                             |

## Reglas de Arquitectura

- **Svelte 5 siempre** — runes, snippets, callback props. Sin `$store` legacy ni `on:event`
- **Responsive-first** — mobile → tablet → desktop con Tailwind. Sin breakpoints hardcodeados
- **TypeScript estricto** — todos los props y retornos tipados. `strict: true` + extras en tsconfig
- **Sin lógica de negocio en componentes** — delegar a `src/lib/api/`
- **SSR-safe** — no acceder a `window`/`document` sin guard `browser`
- **Design tokens via `@theme`** — colores, radios y animaciones en CSS, nunca valores arbitrarios
- **`cn()` para clases dinámicas** — siempre usar el helper `cn()` de `$lib/utils`
- **`import type`** — obligatorio para imports de solo tipos (ESLint lo fuerza)

## Convenciones de Naming

| Tipo                    | Naming              | Ubicación                |
| ----------------------- | ------------------- | ------------------------ |
| Componente reutilizable | `PascalCase.svelte` | `src/lib/components/`    |
| Componente UI primitivo | `PascalCase.svelte` | `src/lib/components/ui/` |
| Página                  | `+page.svelte`      | `src/routes/{ruta}/`     |
| Layout                  | `+layout.svelte`    | `src/routes/{ruta}/`     |
| Server handler          | `+server.ts`        | `src/routes/{ruta}/`     |
| Cliente API             | `camelCase.ts`      | `src/lib/api/`           |
| Tipos                   | `*.types.ts`        | `src/lib/types/`         |
| Utilidades              | `camelCase.ts`      | `src/lib/utils/`         |

## Patrones de Imports

```typescript
// Utilidades
import { cn } from '$lib/utils';

// Tipos
import type { ApiResult, ApiError } from '$lib/types';

// Cliente API
import { apiClient } from '$lib/api/client';

// Componentes
import Button from '$lib/components/ui/Button.svelte';

// Evitar: import relativo largo
// ✗ import { cn } from '../../lib/utils/cn';
// ✓ import { cn } from '$lib/utils';
```

## Variables de Entorno

| Variable         | Descripción           | Default                 |
| ---------------- | --------------------- | ----------------------- |
| `PUBLIC_API_URL` | URL base de la API Go | `http://localhost:8080` |

Prefijo `PUBLIC_` = accesible en cliente. Sin prefijo = solo server-side.

## Comandos

```bash
npm run dev            # Servidor de desarrollo (http://localhost:5173)
npm run build          # Build producción
npm run preview        # Preview del build

# Calidad de código
npm run check          # svelte-check + TypeScript (0 errores requerido)
npm run lint           # ESLint sobre src/ (--max-warnings=0)
npm run lint:fix       # ESLint con auto-fix
npm run format         # Prettier --write src/
npm run format:check   # Prettier --check (solo verifica, no modifica)
npm run qa             # Pipeline completo: format:check + lint + check
```

> `npm run qa` es el comando a ejecutar antes de cada PR o commit manual.

## QA Checklist

- [ ] `npm run qa` pasa sin errores ni warnings
- [ ] Todos los componentes usan `$state`/`$derived` (no `let` reactivo sin rune)
- [ ] Props tipados con TypeScript (sin `any` — ESLint lo bloquea)
- [ ] `import type` para imports de solo tipos
- [ ] Sin acceso a `window`/`document` sin guard `browser`
- [ ] `svelte-autofixer` ejecutado sobre nuevos componentes
- [ ] Layout responsive verificado en mobile (375px), tablet (768px), desktop (1280px)
- [ ] Design tokens usados desde `@theme` (sin valores arbitrarios de color)
- [ ] `cn()` usado para clases dinámicas (no template strings)
- [ ] Variables de entorno públicas con prefijo `PUBLIC_`
