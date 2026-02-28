/**
 * Calculos feature - vertical slicing architecture.
 *
 * Layers:
 * - domain: Types, validators, business rules
 * - infrastructure: API clients, mappers
 * - application: Stores, services
 *
 * Usage:
 * ```typescript
 * import { memoriaStore } from '$lib/features/calculos';
 * import type { MemoriaRequest } from '$lib/features/calculos';
 * ```
 */

// Domain layer
export * from './domain/types';
export {
	validarMemoriaRequest,
	type ValidationError,
	type ValidationResult
} from './domain/validators';

// Infrastructure layer
export * from './infrastructure';

// Application layer
export * from './application';
