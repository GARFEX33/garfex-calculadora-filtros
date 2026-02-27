/**
 * Validators for calculos domain.
 * Pure functions that validate data before sending to the API.
 * NO dependencies on Svelte, fetch, or application layer.
 */

import type { CalcularMemoriaRequest } from '../types/memoria.types.js';
import type { ModoCalculo, TipoEquipo, SistemaElectrico, TipoVoltaje, TipoCanalizacion, MaterialConductor } from '../types/index.js';

// Valid enum values
const MODO_CALCULO_VALUES: readonly ModoCalculo[] = ['LISTADO', 'MANUAL_AMPERAJE', 'MANUAL_POTENCIA'];
const TIPO_EQUIPO_VALUES: readonly TipoEquipo[] = ['FILTRO_ACTIVO', 'TRANSFORMADOR', 'FILTRO_RECHAZO', 'CARGA'];
const SISTEMA_ELECTRICO_VALUES: readonly SistemaElectrico[] = ['DELTA', 'ESTRELLA', 'BIFASICO', 'MONOFASICO'];
const TIPO_VOLTAJE_VALUES: readonly TipoVoltaje[] = ['FASE_NEUTRO', 'FASE_FASE'];
const TIPO_CANALIZACION_VALUES: readonly TipoCanalizacion[] = [
	'TUBERIA_PVC',
	'TUBERIA_ALUMINIO',
	'TUBERIA_ACERO_PG',
	'TUBERIA_ACERO_PD',
	'CHAROLA_CABLE_ESPACIADO',
	'CHAROLA_CABLE_TRIANGULAR'
];
const MATERIAL_CONDUCTOR_VALUES: readonly MaterialConductor[] = ['CU', 'AL'];

/**
 * Validation error with field and message.
 */
export interface ValidationError {
	field: string;
	message: string;
}

/**
 * Validation result.
 */
export interface ValidationResult {
	valid: boolean;
	errors: ValidationError[];
}

/**
 * Validates that a value is in a list of allowed values.
 */
function isValidEnum<T>(value: unknown, allowed: readonly T[]): value is T {
	return allowed.includes(value as T);
}

/**
 * Validates a CalcularMemoriaRequest before sending to the API.
 *
 * Rules:
 * - modo is a valid ModoCalculo value
 * - If modo === 'MANUAL_AMPERAJE': amperaje_nominal > 0 and tipo_equipo is valid
 * - If modo === 'MANUAL_POTENCIA': potencia_nominal > 0 and tipo_equipo is valid
 * - If modo === 'LISTADO': equipo object exists with required fields
 * - tension > 0
 * - tipo_voltaje is 'FASE_NEUTRO' or 'FASE_FASE'
 * - sistema_electrico is valid
 * - tipo_canalizacion is valid (including TUBERIA_ALUMINIO)
 * - longitud_circuito > 0
 * - hilos_por_fase >= 1
 */
export function validarMemoriaRequest(input: unknown): ValidationResult {
	const errors: ValidationError[] = [];

	// Must be an object
	if (!input || typeof input !== 'object') {
		return {
			valid: false,
			errors: [{ field: 'root', message: 'La solicitud debe ser un objeto' }]
		};
	}

	const req = input as Partial<CalcularMemoriaRequest>;

	// Required: modo
	if (!req.modo) {
		errors.push({ field: 'modo', message: 'El modo de cálculo es requerido' });
	} else if (!isValidEnum(req.modo, MODO_CALCULO_VALUES)) {
		errors.push({ field: 'modo', message: `Modo inválido: ${req.modo}` });
	}

	// Validate mode-specific fields
	if (req.modo === 'LISTADO') {
		// LISTADO mode requires equipo
		if (!req.equipo) {
			errors.push({ field: 'equipo', message: 'El equipo es requerido en modo LISTADO' });
		} else {
			// Validate equipo fields
			if (!req.equipo.clave) {
				errors.push({ field: 'equipo.clave', message: 'La clave del equipo es requerida' });
			}
			if (typeof req.equipo.voltaje !== 'number' || req.equipo.voltaje <= 0) {
				errors.push({ field: 'equipo.voltaje', message: 'El voltaje del equipo debe ser mayor a 0' });
			}
			if (typeof req.equipo.amperaje !== 'number' || req.equipo.amperaje <= 0) {
				errors.push({ field: 'equipo.amperaje', message: 'El amperaje del equipo debe ser mayor a 0' });
			}
		}
	} else if (req.modo === 'MANUAL_AMPERAJE') {
		// MANUAL_AMPERAJE requires amperaje_nominal and tipo_equipo
		if (typeof req.amperaje_nominal !== 'number' || req.amperaje_nominal <= 0) {
			errors.push({ field: 'amperaje_nominal', message: 'El amperaje nominal debe ser mayor a 0' });
		}
		if (!req.tipo_equipo) {
			errors.push({ field: 'tipo_equipo', message: 'El tipo de equipo es requerido en modo MANUAL_AMPERAJE' });
		} else if (!isValidEnum(req.tipo_equipo, TIPO_EQUIPO_VALUES)) {
			errors.push({ field: 'tipo_equipo', message: `Tipo de equipo inválido: ${req.tipo_equipo}` });
		}
	} else if (req.modo === 'MANUAL_POTENCIA') {
		// MANUAL_POTENCIA requires potencia_nominal and tipo_equipo
		if (typeof req.potencia_nominal !== 'number' || req.potencia_nominal <= 0) {
			errors.push({ field: 'potencia_nominal', message: 'La potencia nominal debe ser mayor a 0' });
		}
		if (!req.tipo_equipo) {
			errors.push({ field: 'tipo_equipo', message: 'El tipo de equipo es requerido en modo MANUAL_POTENCIA' });
		} else if (!isValidEnum(req.tipo_equipo, TIPO_EQUIPO_VALUES)) {
			errors.push({ field: 'tipo_equipo', message: `Tipo de equipo inválido: ${req.tipo_equipo}` });
		}
	}

	// Required: tension
	if (typeof req.tension !== 'number' || req.tension <= 0) {
		errors.push({ field: 'tension', message: 'La tensión debe ser mayor a 0' });
	}

	// Required: tipo_voltaje
	if (!req.tipo_voltaje) {
		errors.push({ field: 'tipo_voltaje', message: 'El tipo de voltaje es requerido' });
	} else if (!isValidEnum(req.tipo_voltaje, TIPO_VOLTAJE_VALUES)) {
		errors.push({ field: 'tipo_voltaje', message: `Tipo de voltaje inválido: ${req.tipo_voltaje}` });
	}

	// Required: sistema_electrico
	if (!req.sistema_electrico) {
		errors.push({ field: 'sistema_electrico', message: 'El sistema eléctrico es requerido' });
	} else if (!isValidEnum(req.sistema_electrico, SISTEMA_ELECTRICO_VALUES)) {
		errors.push({ field: 'sistema_electrico', message: `Sistema eléctrico inválido: ${req.sistema_electrico}` });
	}

	// Required: estado (string, non-empty)
	if (!req.estado || typeof req.estado !== 'string' || req.estado.trim() === '') {
		errors.push({ field: 'estado', message: 'El estado es requerido' });
	}

	// Required: tipo_canalizacion
	if (!req.tipo_canalizacion) {
		errors.push({ field: 'tipo_canalizacion', message: 'El tipo de canalización es requerido' });
	} else if (!isValidEnum(req.tipo_canalizacion, TIPO_CANALIZACION_VALUES)) {
		errors.push({ field: 'tipo_canalizacion', message: `Tipo de canalización inválido: ${req.tipo_canalizacion}` });
	}

	// Required: longitud_circuito
	if (typeof req.longitud_circuito !== 'number' || req.longitud_circuito <= 0) {
		errors.push({ field: 'longitud_circuito', message: 'La longitud del circuito debe ser mayor a 0' });
	}

	// Optional: hilos_por_fase (must be >= 1 if provided)
	if (req.hilos_por_fase !== undefined && (typeof req.hilos_por_fase !== 'number' || req.hilos_por_fase < 1)) {
		errors.push({ field: 'hilos_por_fase', message: 'Los hilos por fase deben ser al menos 1' });
	}

	// Optional: material (must be valid if provided)
	if (req.material !== undefined && !isValidEnum(req.material, MATERIAL_CONDUCTOR_VALUES)) {
		errors.push({ field: 'material', message: `Material inválido: ${req.material}` });
	}

	// Optional: itm (must be > 0 if provided in MANUAL modes)
	if (req.itm !== undefined && (typeof req.itm !== 'number' || req.itm <= 0)) {
		errors.push({ field: 'itm', message: 'El ITM debe ser mayor a 0' });
	}

	return {
		valid: errors.length === 0,
		errors
	};
}
