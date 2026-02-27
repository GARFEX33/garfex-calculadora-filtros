/**
 * Validators for equipos domain.
 * Pure functions that validate data before sending to the API.
 * NO dependencies on Svelte, fetch, or application layer.
 */

import type { ListarEquiposParams, TipoFiltroEquipo } from '../types/equipo.types.js';

// Valid enum values
const TIPO_FILTRO_VALUES: readonly TipoFiltroEquipo[] = ['A', 'KVA', 'KVAR'];

/**
 * Validation error with field and message.
 */
export interface ValidationError {
	field: string;
	message: string;
}

/**
 * Validation result.
 */
export interface ValidationResult {
	valid: boolean;
	errors: ValidationError[];
}

/**
 * Validates that a value is in a list of allowed values.
 */
function isValidEnum<T>(value: unknown, allowed: readonly T[]): value is T {
	return allowed.includes(value as T);
}

/**
 * Validates listarEquipos parameters.
 *
 * Rules:
 * - page >= 1 if provided
 * - per_page must be between 1 and 100 if provided
 * - tipo must be a valid TipoFiltroEquipo if provided
 * - voltaje > 0 if provided
 */
export function validarListarEquiposParams(input: unknown): ValidationResult {
	const errors: ValidationError[] = [];

	// Must be an object
	if (!input || typeof input !== 'object') {
		return {
			valid: false,
			errors: [{ field: 'root', message: 'Los parámetros deben ser un objeto' }]
		};
	}

	const params = input as Partial<ListarEquiposParams>;

	// Optional: page (must be >= 1)
	if (params.page !== undefined) {
		if (typeof params.page !== 'number' || params.page < 1) {
			errors.push({ field: 'page', message: 'La página debe ser mayor o igual a 1' });
		}
	}

	// Optional: per_page (must be between 1 and 100)
	if (params.per_page !== undefined) {
		if (typeof params.per_page !== 'number' || params.per_page < 1 || params.per_page > 100) {
			errors.push({ field: 'per_page', message: 'El límite debe estar entre 1 y 100' });
		}
	}

	// Optional: tipo (must be valid TipoFiltroEquipo)
	if (params.tipo !== undefined && !isValidEnum(params.tipo, TIPO_FILTRO_VALUES)) {
		errors.push({ field: 'tipo', message: `Tipo de filtro inválido: ${params.tipo}` });
	}

	// Optional: voltaje (must be > 0)
	if (params.voltaje !== undefined) {
		if (typeof params.voltaje !== 'number' || params.voltaje <= 0) {
			errors.push({ field: 'voltaje', message: 'El voltaje debe ser mayor a 0' });
		}
	}

	return {
		valid: errors.length === 0,
		errors
	};
}
