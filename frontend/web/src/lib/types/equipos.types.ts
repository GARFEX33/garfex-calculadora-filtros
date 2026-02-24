export type TipoFiltroEquipo = 'A' | 'KVA' | 'KVAR';

export interface EquipoFiltro {
	id: string;
	clave: string;
	tipo: TipoFiltroEquipo;
	voltaje: number;
	amperaje: number; // Qn/In — amperaje nominal (A), potencia (KVA o KVAR) según tipo
	itm: number;
	bornes?: number | null;
	conexion?: string | null; // DELTA, ESTRELLA, MONOFASICO, BIFASICO
	tipo_voltaje?: string | null; // FF (Fase-Fase), FN (Fase-Neutro)
	created_at: string;
}

// Paginación real del backend
export interface EquiposPagination {
	page: number;
	page_size: number;
	total: number;
	total_pages: number;
	has_next: boolean;
	has_prev: boolean;
}

// Respuesta real del backend: { success: true, data: { equipos: [], pagination: {} } }
export interface ListarEquiposData {
	equipos: EquipoFiltro[];
	pagination: EquiposPagination;
}

export interface ListarEquiposResponse {
	success: boolean;
	data: ListarEquiposData;
}

export interface ListarEquiposParams {
	page?: number;
	per_page?: number;
	tipo?: TipoFiltroEquipo;
	voltaje?: number;
	buscar?: string;
}
