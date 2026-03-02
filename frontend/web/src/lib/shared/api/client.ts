import { PUBLIC_API_URL } from '$env/static/public';
import type { ApiError, ApiResult } from '$lib/types/index.js';

/**
 * Cliente HTTP base para la API Go.
 * Centraliza base URL, headers y manejo de errores.
 *
 * NO usar fetch directamente en componentes — siempre a través de este cliente
 * o los módulos específicos en src/lib/api/.
 */

// Usar URL relativa ('') para que el proxy de Vite intercepte las requests en dev
// En prod, usar PUBLIC_API_URL si está definida
const BASE_URL = PUBLIC_API_URL ?? '';

interface RequestOptions extends Omit<RequestInit, 'body'> {
	body?: unknown;
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<ApiResult<T>> {
	const { body, headers, ...rest } = options;

	const init: RequestInit = {
		...rest,
		headers: {
			'Content-Type': 'application/json',
			Accept: 'application/json',
			...headers
		}
	};

	if (body !== undefined) {
		init.body = JSON.stringify(body);
	}

	try {
		const response = await fetch(`${BASE_URL}${path}`, init);

		if (!response.ok) {
			let errorMessage = `Error ${response.status}`;
			try {
				const errorBody = (await response.json()) as { error?: string };
				errorMessage = errorBody?.error ?? errorMessage;
			} catch {
				// no-op: respuesta no tiene JSON
			}

			return {
				ok: false,
				error: { error: errorMessage, status: response.status } satisfies ApiError
			};
		}

		// 204 No Content — sin cuerpo
		if (response.status === 204) {
			return { ok: true, data: null as T };
		}

		const data = (await response.json()) as T;
		return { ok: true, data };
	} catch (err) {
		const message = err instanceof Error ? err.message : 'Error de conexión';
		return {
			ok: false,
			error: { error: message, status: 0 } satisfies ApiError
		};
	}
}

export const apiClient = {
	get: <T>(path: string, options?: Omit<RequestOptions, 'body'>) =>
		request<T>(path, { method: 'GET', ...options }),

	post: <T>(path: string, body?: unknown, options?: Omit<RequestOptions, 'body'>) =>
		request<T>(path, { method: 'POST', body, ...options }),

	put: <T>(path: string, body?: unknown, options?: Omit<RequestOptions, 'body'>) =>
		request<T>(path, { method: 'PUT', body, ...options }),

	patch: <T>(path: string, body?: unknown, options?: Omit<RequestOptions, 'body'>) =>
		request<T>(path, { method: 'PATCH', body, ...options }),

	delete: <T>(path: string, options?: Omit<RequestOptions, 'body'>) =>
		request<T>(path, { method: 'DELETE', ...options })
};
