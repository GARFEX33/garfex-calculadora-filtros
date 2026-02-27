/**
 * Memoria calculation store using Svelte 5 runes.
 *
 * Manages:
 * - Input state (MemoriaRequest)
 * - Output state (MemoriaOutput)
 * - Loading and error states
 * - Validation
 *
 * Uses Svelte 5 $state and $derived for reactivity.
 */

import type { MemoriaRequest, MemoriaOutput } from '../../domain/types/memoria.types';
import { validarMemoriaRequest } from '../../domain/validators/validar-memoria-request';
import type { ValidationError } from '../../domain/validators/validar-memoria-request';
import { calcularMemoria } from '../../infrastructure/api/memoria.api';
import { mapMemoriaInputToApi, mapApiToMemoriaOutput } from '../../infrastructure/mappers/memoria.mapper';

/**
 * Default values for MemoriaRequest.
 */
function getDefaultMemoriaRequest(): MemoriaRequest {
	return {
		modo: 'MANUAL_AMPERAJE',
		tension: 220,
		tension_unidad: 'V',
		sistema_electrico: 'ESTRELLA',
		estado: '',
		tipo_canalizacion: 'TUBERIA_PVC',
		longitud_circuito: 0,
		tipo_voltaje: 'FASE_NEUTRO',
		material: 'CU',
		hilos_por_fase: 1,
		porcentaje_caida_maximo: 3.0
	};
}

/**
 * MemoriaStore class using Svelte 5 runes.
 *
 * Usage:
 * ```typescript
 * import { memoriaStore } from '$lib/features/calculos/application/stores';
 *
 * // Reactive access to state
 * $: if (memoriaStore.esValido) { ... }
 *
 * // Update input
 * memoriaStore.actualizarInput({ tension: 220 });
 *
 * // Calculate
 * await memoriaStore.calcular();
 * ```
 */
class MemoriaStore {
	// ── State ─────────────────────────────────────────────────────────────────
	// Input request data
	input = $state<MemoriaRequest>(getDefaultMemoriaRequest());

	// Output result from API
	output = $state<MemoriaOutput | null>(null);

	// Loading state
	loading = $state(false);

	// Error state (network/API errors)
	error = $state<string | null>(null);

	// Validation errors (field-level)
	erroresValidacion = $state<ValidationError[]>([]);

	// ── Derived Values ────────────────────────────────────────────────────────
	/**
	 * Whether the current input is valid for submission.
	 * Derived from validation errors.
	 */
	esValido = $derived(this.erroresValidacion.length === 0);

	/**
	 * Whether there's any error (validation or API).
	 * Useful for displaying error UI.
	 */
	tieneError = $derived(this.error !== null || this.erroresValidacion.length > 0);

	// ── Methods ───────────────────────────────────────────────────────────────

	/**
	 * Updates the input with partial data.
	 * Filters out undefined values to satisfy exactOptionalPropertyTypes.
	 * Only includes properties that are explicitly provided.
	 */
	actualizarInput(partial: Partial<MemoriaRequest>): void {
		// Filter out undefined values - with exactOptionalPropertyTypes,
		// undefined means "property should be removed", not "set to undefined"
		const filtered = Object.fromEntries(
			Object.entries(partial).filter(([, value]) => value !== undefined)
		);
		this.input = { ...this.input, ...filtered };
		// Clear validation errors when input changes
		this.erroresValidacion = [];
	}

	/**
	 * Resets the store to initial state.
	 */
	resetear(): void {
		this.input = getDefaultMemoriaRequest();
		this.output = null;
		this.loading = false;
		this.error = null;
		this.erroresValidacion = [];
	}

	/**
	 * Performs the memoria calculation.
	 *
	 * Flow:
	 * 1. Validate input
	 * 2. If invalid, populate validation errors and return
	 * 3. If valid, map to API format
	 * 4. Call API
	 * 5. Map response to domain format
	 * 6. Update output state
	 *
	 * @throws Never — errors are captured in state
	 */
	async calcular(): Promise<void> {
		// Step 1: Validate input
		const validacion = validarMemoriaRequest(this.input);

		if (!validacion.valid) {
			this.erroresValidacion = validacion.errors;
			this.output = null;
			this.error = null;
			return;
		}

		// Clear previous errors
		this.erroresValidacion = [];
		this.error = null;
		this.loading = true;

		try {
			// Step 2: Map to API format
			const apiRequest = mapMemoriaInputToApi(this.input);

			// Step 3: Call API
			const result = await calcularMemoria(apiRequest);

			if (result.ok) {
				// Step 4: Map response to domain format
				const dominioOutput = mapApiToMemoriaOutput(result.data.data);
				this.output = dominioOutput;
			} else {
				// API error
				this.error = result.error.error;
				this.output = null;
			}
		} catch (err) {
			// Network error or unexpected exception
			this.error = err instanceof Error ? err.message : 'Error desconocido';
			this.output = null;
		} finally {
			this.loading = false;
		}
	}

	/**
	 * Gets a validation error for a specific field.
	 * Returns undefined if field is valid.
	 */
	getErrorForField(field: string): string | undefined {
		return this.erroresValidacion.find((e) => e.field === field)?.message;
	}
}

/**
 * Singleton instance of the MemoriaStore.
 * Use this in components: `memoriaStore.calcular()`
 */
export const memoriaStore = new MemoriaStore();
