/**
 * Tipos legacy de cálculos eléctricos.
 * Este archivo re-exporta los tipos del dominio para compatibilidad hacia atrás.
 *
 * IMPORTANTE: Los tipos principales están en:
 *   src/lib/features/calculos/domain/types/memoria.types.ts
 *
 * Este archivo existe para mantener compatibilidad con imports existentes.
 * Preferir imports directos desde $lib/features/calculos/ cuando sea posible.
 */

// Re-exportar tipos del dominio
export type {
	DatosEquipo,
	MemoriaOutput,
	MemoriaRequest,
	CalcularMemoriaRequest,
	ResultadoConductor,
	ResultadoCanalizacion,
	DetalleCharola,
	DetalleTuberia,
	ResultadoCaidaTension,
	CalcularMemoriaResponse,
	// Nuevos tipos agrupados
	DatosInstalacion,
	DatosCorrientes,
	DatosCanalizacion,
	DatosProteccion
} from '$lib/features/calculos/domain/types/memoria.types';

// Re-exportar enums
export type {
	ModoCalculo,
	TipoFiltro,
	TipoEquipo,
	SistemaElectrico,
	TipoVoltaje,
	UnidadPotencia,
	UnidadTension
} from '$lib/features/calculos/domain/types/calculo.enums';

export type { TipoCanalizacion } from '$lib/features/calculos/domain/types/tipo-canalizacion';
export type { MaterialConductor } from '$lib/features/calculos/domain/types/material-conductor';

// Re-exportar labels
export {
	TIPO_CANALIZACION_LABELS,
	MATERIAL_CONDUCTOR_LABELS
} from '$lib/features/calculos/domain/types/memoria.types';
