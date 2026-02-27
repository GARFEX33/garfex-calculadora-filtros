/**
 * @deprecated Usar desde '$lib/features/equipos'
 * Re-export para compatibilidad con imports existentes.
 */
export { listarEquipos, obtenerEquipo } from '$lib/features/equipos/infrastructure/api/equipos.api';

// Re-export tipos para compatibilidad
export type {
	ListarEquiposParams,
	ListarEquiposResponse,
	EquipoFiltro
} from '$lib/features/equipos/domain/types/equipo.types';
