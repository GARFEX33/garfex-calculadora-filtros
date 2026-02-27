/**
 * @deprecated Usar desde '$lib/features/equipos/domain/types'
 * Re-export para compatibilidad con imports existentes.
 */
export type { TipoFiltroEquipo, Conexion } from '$lib/features/equipos/domain/types/equipo.enums';
export type {
	EquipoFiltro,
	EquiposPagination,
	ListarEquiposData,
	ListarEquiposResponse,
	ListarEquiposParams
} from '$lib/features/equipos/domain/types/equipo.types';

// Tipos legacy (deprecated) — mantener por compatibilidad
// Los nuevos desarrollos deben usar las exportaciones desde '$lib/features/equipos'
