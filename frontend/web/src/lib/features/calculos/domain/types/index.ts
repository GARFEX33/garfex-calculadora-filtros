/**
 * Domain types exports for calculos feature.
 */

// Enums and union types
export type {
	ModoCalculo,
	TipoFiltro,
	TipoEquipo,
	SistemaElectrico,
	TipoVoltaje,
	UnidadPotencia,
	UnidadTension
} from './calculo.enums.js';
export type { TipoCanalizacion, isTipoCanalizacion } from './tipo-canalizacion.js';
export type { MaterialConductor, isMaterialConductor } from './material-conductor.js';

// Labels
export { TIPO_CANALIZACION_LABELS } from './tipo-canalizacion.js';
export { MATERIAL_CONDUCTOR_LABELS } from './material-conductor.js';

// Data types
export type {
	DatosEquipo,
	CalcularMemoriaRequest,
	MemoriaRequest,
	ResultadoConductor,
	ResultadoCanalizacion,
	DetalleCharola,
	DetalleTuberia,
	ResultadoCaidaTension,
	MemoriaOutput,
	CalcularMemoriaResponse
} from './memoria.types.js';
