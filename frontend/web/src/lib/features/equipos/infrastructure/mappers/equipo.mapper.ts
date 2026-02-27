/**
 * Mappers between domain types and API types.
 * This is the critical boundary between domain and infrastructure.
 *
 * Domain types use camelCase, TypeScript-friendly naming.
 * API types use snake_case (backend format) and exact field names.
 *
 * For equipos, the API response structure closely matches domain types,
 * but we still validate and normalize the data.
 */

import type { EquipoFiltro, ListarEquiposData, EquiposPagination } from '../../domain/types/equipo.types';

/**
 * Raw API response for a single equipo.
 * Backend uses snake_case.
 * Note: Uses string type for conexion as API may return any string value
 */
export interface ApiEquipoFiltro {
	id: string;
	clave: string;
	tipo: 'A' | 'KVA' | 'KVAR';
	voltaje: number;
	amperaje: number;
	itm: number;
	bornes?: number | null;
	conexion?: string | null;
	tipo_voltaje?: string | null;
	created_at: string;
}

/**
 * Raw API response pagination.
 */
export interface ApiEquiposPagination {
	page: number;
	page_size: number;
	total: number;
	total_pages: number;
	has_next: boolean;
	has_prev: boolean;
}

/**
 * Raw API response for list endpoint.
 */
export interface ApiListarEquiposResponse {
	success: boolean;
	data: {
		equipos: ApiEquipoFiltro[];
		pagination: ApiEquiposPagination;
	};
}

/**
 * Raw API response for single equipo.
 */
export interface ApiObtenerEquipoResponse {
	success: boolean;
	data: ApiEquipoFiltro;
}

/**
 * Maps raw API equipo to domain EquipoFiltro.
 * Handles optional field conversion using exactOptionalPropertyTypes pattern:
 * - Only include property if it has a meaningful value
 * - Omit the property entirely if undefined (satisfies exactOptionalPropertyTypes)
 */
export function mapApiToEquipo(api: ApiEquipoFiltro): EquipoFiltro {
	// Start with required fields only
	const result: EquipoFiltro = {
		id: api.id,
		clave: api.clave,
		tipo: api.tipo,
		voltaje: api.voltaje,
		amperaje: api.amperaje,
		itm: api.itm,
		created_at: api.created_at
	};

	// Optional fields - add only if they have meaningful values (not null/undefined)
	// Using type assertion to satisfy exactOptionalPropertyTypes
	if (api.bornes !== undefined && api.bornes !== null) {
		// Using explicit assignment to satisfy exactOptionalPropertyTypes
		(result as { bornes?: number | null }).bornes = api.bornes;
	}

	if (api.conexion !== undefined && api.conexion !== null) {
		// Cast string to Conexion type (backend sends valid values)
		(result as { conexion?: string | null }).conexion = api.conexion;
	}

	if (api.tipo_voltaje !== undefined && api.tipo_voltaje !== null) {
		(result as { tipo_voltaje?: string | null }).tipo_voltaje = api.tipo_voltaje;
	}

	return result;
}

/**
 * Maps raw API pagination to domain EquiposPagination.
 */
function mapApiToPagination(api: ApiEquiposPagination): EquiposPagination {
	return {
		page: api.page,
		page_size: api.page_size,
		total: api.total,
		total_pages: api.total_pages,
		has_next: api.has_next,
		has_prev: api.has_prev
	};
}

/**
 * Maps raw API list response to domain ListarEquiposData.
 */
export function mapApiToListarEquiposData(api: ApiListarEquiposResponse['data']): ListarEquiposData {
	return {
		equipos: api.equipos.map(mapApiToEquipo),
		pagination: mapApiToPagination(api.pagination)
	};
}
