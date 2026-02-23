# Frontend Web — SvelteKit + TypeScript

App web responsive (PWA) para la calculadora de filtros eléctricos. Consume la API Go en `internal/`.

## Skills OBLIGATORIOS (cargar PRIMERO)

| Acción | Skill |
| ------ | ----- |
| Crear o editar cualquier archivo `.svelte`, `.svelte.ts`, `.svelte.js` | `svelte-code-writer` |
| Usar runes (`$state`, `$derived`, `$effect`, `$props`, `$bindable`) o snippets | `svelte5-best-practices` |
| Definir rutas, layouts, `+page.svelte`, `+layout.svelte`, `+server.ts`, SSR | `sveltekit-structure` |
| Crear componentes UI, design tokens, responsive layout, dark mode | `tailwind-design-system` |
| Tipos genéricos, tipos condicionales, type-safe API clients, utilidades TS | `typescript-advanced-types` |
| Detectar o corregir bugs en componentes | `systematic-debugging` |
| Verificar antes de declarar trabajo listo | `verification-before-completion` |

> **Nota Tailwind**: el SKILL.md usa ejemplos React — adaptar CVA/componentes al patrón Svelte 5.
> En Svelte se usa `class:` y props en lugar de `VariantProps`. Ver `svelte5-best-practices` para integración.

## Estructura

```
frontend/web/
├── src/
│   ├── lib/
│   │   ├── components/     # Componentes reutilizables (PascalCase.svelte)
│   │   ├── api/            # Clientes HTTP hacia la API Go
│   │   └── types/          # Tipos TypeScript (*.types.ts)
│   └── routes/
│       ├── +layout.svelte  # Layout raíz (nav + responsive shell)
│       ├── +page.svelte    # Home
│       ├── calculos/       # Rutas de cálculo
│       └── equipos/        # Rutas de catálogo
├── static/
├── package.json
├── svelte.config.js
├── tsconfig.json
└── vite.config.ts
```

## Reglas de Arquitectura

- **Svelte 5 siempre** — runes, snippets, callback props. Sin `$store` legacy ni `on:event`
- **Responsive-first** — mobile → tablet → desktop con Tailwind. Sin breakpoints hardcodeados
- **TypeScript estricto** — todos los props y retornos tipados. Usar `typescript-advanced-types` para API clients
- **Sin lógica de negocio en componentes** — delegar a `src/lib/api/`
- **SSR-safe** — no acceder a `window`/`document` sin guard `browser`
- **Design tokens via `@theme`** — colores, radios y animaciones en CSS, nunca valores arbitrarios

## Convenciones

| Tipo | Naming | Ubicación |
| ---- | ------ | --------- |
| Componente | `PascalCase.svelte` | `src/lib/components/` |
| Página | `+page.svelte` | `src/routes/{ruta}/` |
| Layout | `+layout.svelte` | `src/routes/{ruta}/` |
| Cliente API | `camelCase.ts` | `src/lib/api/` |
| Tipos | `*.types.ts` | `src/lib/types/` |

## API Base URL

Variable de entorno: `PUBLIC_API_URL=http://localhost:8080`

## Comandos

```bash
npm run dev        # Desarrollo
npm run build      # Build producción
npm run check      # Check de tipos
npm run lint       # Lint
```

## QA Checklist

- [ ] Todos los componentes usan `$state`/`$derived` (no `let` reactivo)
- [ ] Props tipados con TypeScript
- [ ] Sin acceso a `window`/`document` sin guard `browser`
- [ ] `svelte-autofixer` ejecutado sobre nuevos componentes
- [ ] Layout responsive verificado en mobile, tablet y desktop
- [ ] Design tokens usados desde `@theme` (sin valores arbitrarios de color)
