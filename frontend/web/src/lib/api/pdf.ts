import { PUBLIC_API_URL } from '$env/static/public';
import type { MemoriaOutput } from '$lib/types/calculos.types';

/** Tiempo de espera máximo para generación de PDF (ms) */
const PDF_TIMEOUT_MS = 30_000;

const BASE_URL = PUBLIC_API_URL ?? '';

/**
 * Datos de presentación del PDF ingresados por el usuario.
 * Corresponde a dto.PresentacionInput en el backend.
 */
export interface PresentacionInput {
	empresa_id: string;
	nombre_proyecto: string;
	direccion_proyecto: string;
	responsable: string;
	nombre_equipo_override?: string | undefined;
}

/**
 * Body del POST /api/v1/pdf/memoria.
 * Corresponde a dto.PdfMemoriaRequest en el backend.
 */
export interface PdfMemoriaRequest {
	memoria: MemoriaOutput;
	presentacion: PresentacionInput;
}

/**
 * Genera la memoria de cálculo en PDF.
 *
 * @param request - Datos del cálculo y de presentación
 * @returns Blob con el PDF generado (application/pdf)
 * @throws Error si la petición falla, supera los 30s o el servidor responde con error
 */
export async function generarMemoriaPdf(request: PdfMemoriaRequest): Promise<Blob> {
	const controller = new AbortController();
	const timeoutId = setTimeout(() => controller.abort(), PDF_TIMEOUT_MS);

	try {
		const response = await fetch(`${BASE_URL}/api/v1/pdf/memoria`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(request),
			signal: controller.signal
		});

		if (!response.ok) {
			let errorMessage = `Error ${response.status} al generar el PDF`;
			try {
				const errorBody = (await response.json()) as { error?: string };
				errorMessage = errorBody?.error ?? errorMessage;
			} catch {
				// no-op: respuesta sin body JSON
			}
			throw new Error(errorMessage);
		}

		const blob = await response.blob();
		return blob;
	} catch (err) {
		if (err instanceof Error && err.name === 'AbortError') {
			throw new Error('La generación del PDF superó el tiempo máximo de 30 segundos', {
				cause: err
			});
		}
		throw err;
	} finally {
		clearTimeout(timeoutId);
	}
}

/**
 * Dispara la descarga de un Blob como archivo en el navegador.
 *
 * @param blob - El contenido del archivo
 * @param filename - Nombre del archivo a descargar
 */
export function descargarBlob(blob: Blob, filename: string): void {
	const url = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = url;
	link.download = filename;
	document.body.appendChild(link);
	link.click();
	document.body.removeChild(link);
	URL.revokeObjectURL(url);
}
