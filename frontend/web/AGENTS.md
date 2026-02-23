# Frontend Web — SvelteKit + TypeScript

Aplicación web para la calculadora de filtros eléctricos. Consume la API Go en `internal/`.

## Skills OBLIGATORIOS (cargar PRIMERO)

| Acción | Skill |
| ------ | ----- |
| Crear o editar cualquier archivo `.svelte`, `.svelte.ts`, `.svelte.js` | `svelte-code-writer` |
| Usar runes (`$state`, `$derived`, `$effect`, `$props`, `$bindable`) o snippets | `svelte5-best-practices` |
| Definir rutas, layouts, `+page.svelte`, `+layout.svelte`, `+server.ts`, SSR | `sveltekit-structure` |
| Detectar o corregir bugs en componentes | `systematic-debugging` |
| Verificar antes de declarar trabajo listo | `verification-before-completion` |

## Estructura Esperada

```
frontend/web/
├── src/
│   ├── lib/
│   │   ├── components/     # Componentes reutilizables
│   │   ├── api/            # Clientes HTTP hacia la API Go
│   │   └── types/          # Tipos TypeScript compartidos
│   └── routes/
│       ├── +layout.svelte  # Layout raíz
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

- **Svelte 5 siempre** — usar runes, nunca `$store` legacy ni `on:event`
- **TypeScript estricto** — todos los props y retornos tipados
- **Sin lógica de negocio en componentes** — delegar a `src/lib/api/`
- **Snippets en lugar de slots** — usar `{#snippet}` y `{@render}`
- **Eventos como callback props** — no `createEventDispatcher`
- **SSR-safe** — no acceder a `window`/`document` sin `browser` check

## Convenciones

| Tipo | Naming | Ubicación |
| ---- | ------ | --------- |
| Componente reutilizable | `PascalCase.svelte` | `src/lib/components/` |
| Página SvelteKit | `+page.svelte` | `src/routes/{ruta}/` |
| Layout | `+layout.svelte` | `src/routes/{ruta}/` |
| Cliente API | `camelCase.ts` | `src/lib/api/` |
| Tipos | `*.types.ts` | `src/lib/types/` |

## API Base URL

Configurar en variable de entorno: `PUBLIC_API_URL=http://localhost:8080`

## Comandos

```bash
# Desarrollo
npm run dev

# Build producción
npm run build

# Preview build
npm run preview

# Check de tipos
npm run check

# Lint
npm run lint
```

## QA Checklist

- [ ] Todos los componentes usan `$state`/`$derived` (no `let` reactivo)
- [ ] Props tipados con TypeScript
- [ ] Sin imports de `window`/`document` sin guard `browser`
- [ ] `svelte-autofixer` ejecutado sobre nuevos componentes
- [ ] Rutas con `+page.svelte` y datos cargados en `+page.ts` (no inline fetch)
