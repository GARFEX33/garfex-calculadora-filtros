import { apiClient } from '$lib/api/client';
import type { ApiResult } from '$lib/types/index.js';
import type { CalcularMemoriaRequest, CalcularMemoriaResponse } from '$lib/types/calculos.types.js';

export async function calcularMemoria(
	request: CalcularMemoriaRequest
): Promise<ApiResult<CalcularMemoriaResponse>> {
	return apiClient.post<CalcularMemoriaResponse>('/api/v1/calculos/memoria', request);
}
