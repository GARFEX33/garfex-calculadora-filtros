import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { MemoriaOutput } from '$lib/types/calculos.types';

export const load: PageLoad = ({ url }) => {
	const data = url.searchParams.get('data');

	if (!data) {
		error(400, 'Faltan datos del cálculo');
	}

	try {
		const decoded = atob(data);
		const jsonStr = decodeURIComponent(escape(decoded));
		const parsed = JSON.parse(jsonStr) as MemoriaOutput;
		return { memoria: parsed };
	} catch {
		error(400, 'Datos del cálculo inválidos');
	}
};
