# Skill Registry — garfex-calculadora-filtros

Generated: 2026-05-06

## Project Context

**Stack**: Go 1.24 (Gin, pgx/v5, Gotenberg) + SvelteKit 2 / Svelte 5 / TypeScript / Tailwind 4  
**Architecture**: Clean + Hexagonal + Vertical Slicing (Screaming Architecture)  
**Test runner**: `go test ./...` (stretchr/testify) — standard mode  
**Artifact store**: engram only (no openspec/)

## Compact Rules

### `clean-ddd-hexagonal-vertical-go-enterprise`
- Hexagonal strict: domain never imports infra
- Feature slices in `internal/{feature}/domain|application|infrastructure`
- Ports (interfaces) in `application/port/`, adapters in `infrastructure/adapter/driven|driver`
- No business logic in infrastructure layer

### `golang-patterns`
- Idiomatic Go: errors as values, no panics in business logic
- Interfaces defined at consumption site (domain/application ports)
- Constructor functions `New*` return interface, not concrete type

### `svelte5-best-practices`
- Svelte 5 runes only: `$state`, `$derived`, `$effect`, `$props`
- No legacy `$store` syntax, no `on:event` (use callback props)
- SSR-safe: no `window`/`document` without `browser` guard

### `sveltekit-structure`
- Routes in `src/routes/`, layouts in `+layout.svelte`
- Server logic in `+server.ts`, data loading in `+page.ts`/`+layout.ts`

### `tailwind-design-system`
- Tailwind v4 with `@theme` tokens in `src/app.css`
- No arbitrary color values, use semantic tokens
- `cn()` helper for dynamic classes

### `commit-work`
- Conventional commits: `feat:`, `fix:`, `refactor:`, `chore:`
- No AI attribution in commits

### `verification-before-completion`
- Run `go test ./...` + `npm run qa` before claiming done

## User Skills Trigger Table

| Trigger | Skill |
|---------|-------|
| `*.go` in `internal/` | `clean-ddd-hexagonal-vertical-go-enterprise`, `golang-patterns` |
| Go concurrency/generics | `golang-pro` |
| `*.svelte`, `*.svelte.ts` | `svelte5-best-practices` |
| Routes, layouts, SSR | `sveltekit-structure` |
| UI components, Tailwind | `tailwind-design-system` |
| TypeScript types/generics | `typescript-advanced-types` |
| API design/review | `api-design-principles` |
| Commit | `commit-work` |
| Bug/test failure | `systematic-debugging` |
| Before declaring done | `verification-before-completion` |
| Finish branch | `finishing-a-development-branch` |
