import { apiClient } from '$lib/api/client';
import type { ApiResult } from '$lib/types/index.js';
import type {
	EquipoFiltro,
	ListarEquiposParams,
	ListarEquiposResponse
} from '$lib/types/equipos.types.js';

export async function listarEquipos(
	params: ListarEquiposParams = {}
): Promise<ApiResult<ListarEquiposResponse>> {
	const searchParams = new URLSearchParams();
	if (params.page !== undefined) searchParams.set('page', String(params.page));
	if (params.per_page !== undefined) searchParams.set('per_page', String(params.per_page));
	if (params.tipo !== undefined) searchParams.set('tipo', params.tipo);
	if (params.voltaje !== undefined) searchParams.set('voltaje', String(params.voltaje));
	if (params.buscar !== undefined) searchParams.set('buscar', params.buscar);

	const query = searchParams.toString();
	const path = query ? `/api/v1/equipos?${query}` : '/api/v1/equipos';
	return apiClient.get<ListarEquiposResponse>(path);
}

export async function obtenerEquipo(id: string): Promise<ApiResult<EquipoFiltro>> {
	return apiClient.get<EquipoFiltro>(`/api/v1/equipos/${id}`);
}
