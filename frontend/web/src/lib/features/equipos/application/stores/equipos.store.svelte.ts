/**
 * Equipos store using Svelte 5 runes.
 *
 * Manages:
 * - Equipos list state
 * - Pagination state
 * - Loading and error states
 * - Search/filter state
 *
 * Uses Svelte 5 $state and $derived for reactivity.
 */

import type { EquipoFiltro, ListarEquiposParams } from '../../domain/types/equipo.types';
import { listarEquipos } from '../../infrastructure/api/equipos.api';
import { mapApiToListarEquiposData } from '../../infrastructure/mappers/equipo.mapper';

/**
 * EquiposStore class using Svelte 5 runes.
 *
 * Usage:
 * ```typescript
 * import { equiposStore } from '$lib/features/equipos';
 *
 * // Load equipos on mount
 * equiposStore.cargar();
 *
 * // Search
 * await equiposStore.buscar('400');
 *
 * // Change page
 * await equiposStore.cambiarPagina(2);
 * ```
 */
class EquiposStore {
	// ── State ─────────────────────────────────────────────────────────────────
	equipos = $state<EquipoFiltro[]>([]);

	total = $state(0);

	totalPages = $state(0);

	loading = $state(false);

	error = $state<string | null>(null);

	query = $state('');

	pagina = $state(1);

	limite = $state(20);

	// ── Derived Values ────────────────────────────────────────────────────────
	/**
	 * Whether there are equipos to display.
	 */
	hasEquipos = $derived(this.equipos.length > 0);

	/**
	 * Whether there's a next page.
	 */
	hasNextPage = $derived(this.pagina < this.totalPages);

	/**
	 * Whether there's a previous page.
	 */
	hasPrevPage = $derived(this.pagina > 1);

	// ── Methods ───────────────────────────────────────────────────────────────

	/**
	 * Loads equipos with current params.
	 * Uses current query, page, and limit.
	 */
	async cargar(): Promise<void> {
		this.loading = true;
		this.error = null;

		try {
			const params: ListarEquiposParams = {
				page: this.pagina,
				per_page: this.limite
			};

			if (this.query) {
				params.buscar = this.query;
			}

			const result = await listarEquipos(params);

			if (result.ok) {
				const data = mapApiToListarEquiposData(result.data.data);
				this.equipos = data.equipos;
				this.total = data.pagination.total;
				this.totalPages = data.pagination.total_pages;
			} else {
				this.error = result.error.error;
				this.equipos = [];
				this.total = 0;
				this.totalPages = 0;
			}
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Error desconocido';
			this.equipos = [];
			this.total = 0;
			this.totalPages = 0;
		} finally {
			this.loading = false;
		}
	}

	/**
	 * Searches for equipos with a new query.
	 * Resets page to 1.
	 */
	async buscar(query: string): Promise<void> {
		this.query = query;
		this.pagina = 1;
		await this.cargar();
	}

	/**
	 * Changes to a specific page.
	 */
	async cambiarPagina(pagina: number): Promise<void> {
		if (pagina < 1 || pagina > this.totalPages) {
			return;
		}
		this.pagina = pagina;
		await this.cargar();
	}

	/**
	 * Goes to the next page.
	 */
	async siguientePagina(): Promise<void> {
		if (this.hasNextPage) {
			await this.cambiarPagina(this.pagina + 1);
		}
	}

	/**
	 * Goes to the previous page.
	 */
	async paginaAnterior(): Promise<void> {
		if (this.hasPrevPage) {
			await this.cambiarPagina(this.pagina - 1);
		}
	}

	/**
	 * Clears the current search and reloads.
	 */
	async limpiarBusqueda(): Promise<void> {
		this.query = '';
		this.pagina = 1;
		await this.cargar();
	}

	/**
	 * Resets the store to initial state.
	 */
	resetear(): void {
		this.equipos = [];
		this.total = 0;
		this.totalPages = 0;
		this.loading = false;
		this.error = null;
		this.query = '';
		this.pagina = 1;
		this.limite = 20;
	}
}

/**
 * Singleton instance of the EquiposStore.
 * Use this in components: `equiposStore.cargar()`
 */
export const equiposStore = new EquiposStore();
