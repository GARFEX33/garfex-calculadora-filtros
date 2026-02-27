/**
 * Tipos base para la comunicación con la API Go.
 * Reflejan la estructura de respuestas del backend.
 */

// ── Respuesta base ────────────────────────────────────────────────────────────

/** Respuesta exitosa envuelta en data */
export interface ApiResponse<T> {
	data: T;
}

/** Respuesta paginada */
export interface ApiPaginatedResponse<T> {
	data: T[];
	meta: PaginationMeta;
}

export interface PaginationMeta {
	total: number;
	page: number;
	per_page: number;
	total_pages: number;
}

// ── Error de API ──────────────────────────────────────────────────────────────

export interface ApiError {
	/** Mensaje legible para el usuario */
	error: string;
	/** Código HTTP */
	status: number;
}

/** Resultado tipado — éxito o error */
export type ApiResult<T> = { ok: true; data: T } | { ok: false; error: ApiError };

// ── Parámetros de consulta comunes ────────────────────────────────────────────

export interface PaginationParams {
	page?: number;
	per_page?: number;
}
