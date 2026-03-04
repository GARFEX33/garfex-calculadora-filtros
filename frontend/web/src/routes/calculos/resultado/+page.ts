import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { MemoriaOutput } from '$lib/types/calculos.types';

export const load: PageLoad = ({ url }) => {
	const id = url.searchParams.get('id');

	if (!id) {
		error(400, 'Faltan datos del cálculo');
	}

	const data = sessionStorage.getItem(`memoria-${id}`);

	if (!data) {
		error(400, 'Los datos del cálculo han expirado');
	}

	try {
		const parsed = JSON.parse(data) as MemoriaOutput;
		// Limpiar sessionStorage después de leer
		sessionStorage.removeItem(`memoria-${id}`);
		return { memoria: parsed };
	} catch (err) {
		console.error('Error parsing memoria data:', err);
		error(400, 'Datos del cálculo inválidos');
	}
};
