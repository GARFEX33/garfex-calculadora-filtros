/**
 * Equipos feature - vertical slicing architecture.
 *
 * Layers:
 * - domain: Types, validators, business rules
 * - infrastructure: API clients, mappers
 * - application: Stores, services
 *
 * Usage:
 * ```typescript
 * import { equiposStore } from '$lib/features/equipos';
 * import type { EquipoFiltro } from '$lib/features/equipos';
 * ```
 */

// Domain layer
export * from './domain/types';
export * from './domain/validators';

// Infrastructure layer
export * from './infrastructure';

// Application layer
export * from './application';
