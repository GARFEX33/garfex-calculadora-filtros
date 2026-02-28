/**
 * Mappers between domain types and API types.
 * This is the critical boundary between domain and infrastructure.
 *
 * Domain types use camelCase, TypeScript-friendly naming.
 * API types use snake_case (backend format) and exact field names.
 *
 * IMPORTANT: Uses exactOptionalPropertyTypes - optional fields must be omitted,
 * not explicitly set to undefined.
 */

import type {
	DatosEquipo,
	MemoriaOutput,
	MemoriaRequest,
	DetalleCharola,
	DetalleTuberia,
	ResultadoConductor,
	ResultadoCanalizacion,
	ResultadoCaidaTension,
	DatosInstalacion,
	DatosCorrientes,
	DatosCanalizacion,
	DatosProteccion
} from '../../domain/types/memoria.types';
import type {
	ApiMemoriaRequest,
	ApiMemoriaOutput,
	ApiResultadoConductor,
	ApiResultadoCanalizacion,
	ApiResultadoCaidaTension,
	ApiDetalleCharola,
	ApiDetalleTuberia,
	ApiDatosInstalacion,
	ApiDatosCorrientes,
	ApiDatosCanalizacion,
	ApiDatosProteccion,
	ApiDatosEquipo
} from '../api/memoria.api';
import type { SistemaElectrico, TipoVoltaje } from '../../domain/types/memoria.types';

/**
 * Maps equipment connection type to SistemaElectrico.
 * Returns undefined if connection is not valid.
 */
export function mapConexionToSistemaElectrico(
	conexion: string | null | undefined
): SistemaElectrico | undefined {
	if (!conexion) return undefined;

	const mapa: Record<string, SistemaElectrico> = {
		DELTA: 'DELTA',
		ESTRELLA: 'ESTRELLA',
		MONOFASICO: 'MONOFASICO',
		BIFASICO: 'BIFASICO'
	};
	return mapa[conexion];
}

/**
 * Maps equipment voltage type to TipoVoltaje.
 * Returns undefined if voltage type is not valid.
 */
export function mapTipoVoltajeToTipoVoltaje(
	tipoVoltaje: string | null | undefined
): TipoVoltaje | undefined {
	if (!tipoVoltaje) return undefined;

	const mapa: Record<string, TipoVoltaje> = {
		FF: 'FASE_FASE',
		FN: 'FASE_NEUTRO'
	};
	return mapa[tipoVoltaje];
}

/**
 * Maps domain MemoriaRequest to API request format.
 *
 * Transformations:
 * - factor_potencia: converts from percentage (e.g., 98) to decimal (0.98)
 * - MaterialConductor: domain uses 'CU'/'AL', API uses same format (no change)
 * - Handles optional field inclusion (only includes fields with actual values)
 */
export function mapMemoriaInputToApi(domain: MemoriaRequest): ApiMemoriaRequest {
	// Start with required fields
	const base: ApiMemoriaRequest = {
		modo: domain.modo,
		tension: domain.tension,
		tipo_voltaje: domain.tipo_voltaje,
		sistema_electrico: domain.sistema_electrico,
		estado: domain.estado,
		tipo_canalizacion: domain.tipo_canalizacion,
		longitud_circuito: domain.longitud_circuito
	};

	// Optional: tension_unidad
	const optionalFields: Partial<ApiMemoriaRequest> = {};
	if (domain.tension_unidad !== undefined) {
		optionalFields.tension_unidad = domain.tension_unidad;
	}

	// LISTADO mode: include equipo
	if (domain.modo === 'LISTADO' && domain.equipo) {
		const equipo: ApiMemoriaRequest['equipo'] = {
			clave: domain.equipo.clave,
			tipo: domain.equipo.tipo,
			voltaje: domain.equipo.voltaje,
			amperaje: domain.equipo.amperaje,
			itm: domain.equipo.itm
		};
		// Only include bornes if it has a real value
		if (domain.equipo.bornes !== undefined) {
			equipo.bornes = domain.equipo.bornes;
		}
		optionalFields.equipo = equipo;
	}

	// MANUAL modes: include tipo_equipo and mode-specific fields
	if (domain.modo === 'MANUAL_AMPERAJE' || domain.modo === 'MANUAL_POTENCIA') {
		if (domain.tipo_equipo) {
			optionalFields.tipo_equipo = domain.tipo_equipo;
		}

		if (domain.modo === 'MANUAL_AMPERAJE' && domain.amperaje_nominal !== undefined) {
			optionalFields.amperaje_nominal = domain.amperaje_nominal;
		}

		if (domain.modo === 'MANUAL_POTENCIA') {
			if (domain.potencia_nominal !== undefined) {
				optionalFields.potencia_nominal = domain.potencia_nominal;
			}
			if (domain.potencia_unidad) {
				optionalFields.potencia_unidad = domain.potencia_unidad;
			}
			// Convert factor_potencia from percentage (e.g., 98) to decimal (0.98)
			if (domain.factor_potencia !== undefined) {
				optionalFields.factor_potencia = domain.factor_potencia / 100;
			}
		}

		// itm is required in MANUAL modes
		if (domain.itm !== undefined) {
			optionalFields.itm = domain.itm;
		}
	}

	// Optional fields
	if (domain.hilos_por_fase !== undefined) {
		optionalFields.hilos_por_fase = domain.hilos_por_fase;
	}

	if (domain.num_tuberias !== undefined && domain.num_tuberias > 0) {
		optionalFields.num_tuberias = domain.num_tuberias;
	}

	if (domain.material !== undefined) {
		optionalFields.material = domain.material;
	}

	if (domain.porcentaje_caida_maximo !== undefined) {
		optionalFields.porcentaje_caida_maximo = domain.porcentaje_caida_maximo;
	}

	if (domain.temperatura_override !== undefined) {
		optionalFields.temperatura_override = domain.temperatura_override;
	}

	if (domain.diametro_control_mm !== undefined) {
		optionalFields.diametro_control_mm = domain.diametro_control_mm;
	}

	return { ...base, ...optionalFields };
}

/**
 * Maps API response to domain MemoriaOutput.
 * Uses type assertion pattern for optional fields to satisfy exactOptionalPropertyTypes.
 *
 * La nueva estructura agrupada del backend coincide casi 1:1 con el domain,
 * así que el mapper es principalmente un pass-through con validaciones.
 */
export function mapApiToMemoriaOutput(api: ApiMemoriaOutput): MemoriaOutput {
	// Verificar campos obligatorios
	if (!api.equipo) {
		throw new Error('API response missing required field: equipo');
	}
	if (!api.instalacion) {
		throw new Error('API response missing required field: instalacion');
	}
	if (!api.corrientes) {
		throw new Error('API response missing required field: corrientes');
	}
	if (!api.cable_fase) {
		throw new Error('API response missing required field: cable_fase');
	}
	if (!api.cable_tierra) {
		throw new Error('API response missing required field: cable_tierra');
	}
	if (!api.canalizacion) {
		throw new Error('API response missing required field: canalizacion');
	}
	if (!api.canalizacion.resultado) {
		throw new Error('API response missing required field: canalizacion.resultado');
	}
	if (!api.proteccion) {
		throw new Error('API response missing required field: proteccion');
	}
	if (!api.caida_tension) {
		throw new Error('API response missing required field: caida_tension');
	}

	const result: MemoriaOutput = {
		// Datos del equipo
		equipo: mapApiToDatosEquipo(api.equipo),
		tipo_equipo: api.tipo_equipo,
		factor_potencia: api.factor_potencia,
		estado: api.estado,

		// Parámetros de instalación
		instalacion: mapApiToDatosInstalacion(api.instalacion),

		// Cálculos de corriente
		corrientes: mapApiToDatosCorrientes(api.corrientes),

		// Conductores
		cable_fase: mapApiToResultadoConductor(api.cable_fase),
		cable_tierra: mapApiToResultadoConductor(api.cable_tierra),

		// Cable neutro es opcional (nil para sistemas DELTA) - omitir si no existe
		...(api.cable_neutro && { cable_neutro: mapApiToResultadoConductor(api.cable_neutro) }),

		// Canalización
		canalizacion: mapApiToDatosCanalizacion(api.canalizacion),

		// Protección
		proteccion: mapApiToDatosProteccion(api.proteccion),

		// Caída de tensión
		caida_tension: mapApiToResultadoCaidaTension(api.caida_tension),

		// Resumen y metadatos
		cumple_normativa: api.cumple_normativa,
		observaciones: api.observaciones ?? [],
		pasos: api.pasos ?? []
	};

	return result;
}

/**
 * Maps API DatosEquipo to domain DatosEquipo.
 * Handles optional 'bornes' field.
 */
function mapApiToDatosEquipo(api: ApiDatosEquipo): DatosEquipo {
	const result: DatosEquipo = {
		clave: api.clave,
		tipo: api.tipo,
		voltaje: api.voltaje,
		amperaje: api.amperaje,
		itm: api.itm
	};

	if (api.bornes !== undefined) {
		result.bornes = api.bornes;
	}

	return result;
}

/**
 * Maps API DatosInstalacion to domain DatosInstalacion.
 */
function mapApiToDatosInstalacion(api: ApiDatosInstalacion): DatosInstalacion {
	return {
		tension: api.tension,
		sistema_electrico: api.sistema_electrico,
		tipo_canalizacion: api.tipo_canalizacion,
		material: api.material,
		longitud_circuito: api.longitud_circuito,
		hilos_por_fase: api.hilos_por_fase,
		porcentaje_caida_maximo: api.porcentaje_caida_maximo
	};
}

/**
 * Maps API DatosCorrientes to domain DatosCorrientes.
 */
function mapApiToDatosCorrientes(api: ApiDatosCorrientes): DatosCorrientes {
	return {
		corriente_nominal: api.corriente_nominal,
		corriente_ajustada: api.corriente_ajustada,
		corriente_por_hilo: api.corriente_por_hilo,
		factor_temperatura: api.factor_temperatura,
		factor_agrupamiento: api.factor_agrupamiento,
		factor_total_ajuste: api.factor_total_ajuste,
		temperatura_ambiente: api.temperatura_ambiente,
		temperatura_referencia: api.temperatura_referencia,
		conductores_por_tubo: api.conductores_por_tubo,
		cantidad_conductores: api.cantidad_conductores,
		tabla_ampacidad_usada: api.tabla_ampacidad_usada
	};
}

/**
 * Maps API DatosCanalizacion to domain DatosCanalizacion.
 */
function mapApiToDatosCanalizacion(api: ApiDatosCanalizacion): DatosCanalizacion {
	const result: DatosCanalizacion = {
		resultado: mapApiToResultadoCanalizacion(api.resultado),
		fill_factor: api.fill_factor
	};

	// Optional: detalle_charola
	if (api.detalle_charola) {
		result.detalle_charola = mapApiToDetalleCharola(api.detalle_charola);
	}

	// Optional: detalle_tuberia
	if (api.detalle_tuberia) {
		result.detalle_tuberia = mapApiToDetalleTuberia(api.detalle_tuberia);
	}

	return result;
}

/**
 * Maps API DatosProteccion to domain DatosProteccion.
 */
function mapApiToDatosProteccion(api: ApiDatosProteccion): DatosProteccion {
	return {
		itm: api.itm
	};
}

/**
 * Maps API ResultadoConductor to domain ResultadoConductor.
 * Handles optional fields.
 */
function mapApiToResultadoConductor(api: ApiResultadoConductor): ResultadoConductor {
	const result: ResultadoConductor = {
		calibre: api.calibre,
		material: api.material,
		seccion_mm2: api.seccion_mm2,
		tipo_aislamiento: api.tipo_aislamiento,
		capacidad: api.capacidad,
		num_hilos: api.num_hilos // Obligatorio según DTO backend
	};

	if (api.seleccion_por_caida_tension !== undefined)
		result.seleccion_por_caida_tension = api.seleccion_por_caida_tension;
	if (api.calibre_original_ampacidad !== undefined)
		result.calibre_original_ampacidad = api.calibre_original_ampacidad;
	if (api.nota_seleccion !== undefined) result.nota_seleccion = api.nota_seleccion;

	return result;
}

/**
 * Maps API ResultadoCanalizacion to domain ResultadoCanalizacion.
 * Handles optional fields.
 */
function mapApiToResultadoCanalizacion(api: ApiResultadoCanalizacion): ResultadoCanalizacion {
	const result: ResultadoCanalizacion = {
		tamano: api.tamano,
		area_total_mm2: api.area_total_mm2,
		area_requerida_mm2: api.area_requerida_mm2,
		numero_de_tubos: api.numero_de_tubos
	};

	if (api.ancho_comercial_mm !== undefined) {
		result.ancho_comercial_mm = api.ancho_comercial_mm;
	}

	return result;
}

/**
 * Maps API ResultadoCaidaTension to domain ResultadoCaidaTension.
 */
function mapApiToResultadoCaidaTension(api: ApiResultadoCaidaTension): ResultadoCaidaTension {
	return {
		porcentaje: api.porcentaje,
		caida_volts: api.caida_volts,
		cumple: api.cumple,
		limite_porcentaje: api.limite_porcentaje,
		impedancia: api.impedancia,
		resistencia: api.resistencia,
		reactancia: api.reactancia
	};
}

/**
 * Maps API DetalleCharola to domain DetalleCharola.
 * Handles optional fields.
 */
function mapApiToDetalleCharola(api: ApiDetalleCharola): DetalleCharola {
	const result: DetalleCharola = {
		diametro_fase_mm: api.diametro_fase_mm,
		diametro_tierra_mm: api.diametro_tierra_mm,
		espacio_fuerza_mm: api.espacio_fuerza_mm,
		ancho_tierra_mm: api.ancho_tierra_mm
	};

	// Optional fields - add only if present
	if (api.diametro_control_mm !== undefined) result.diametro_control_mm = api.diametro_control_mm;
	if (api.num_hilos_total !== undefined) result.num_hilos_total = api.num_hilos_total;
	if (api.ancho_fuerza_mm !== undefined) result.ancho_fuerza_mm = api.ancho_fuerza_mm;
	if (api.espacio_control_mm !== undefined) result.espacio_control_mm = api.espacio_control_mm;
	if (api.ancho_control_mm !== undefined) result.ancho_control_mm = api.ancho_control_mm;
	if (api.ancho_potencia_mm !== undefined) result.ancho_potencia_mm = api.ancho_potencia_mm;
	if (api.factor_triangular !== undefined) result.factor_triangular = api.factor_triangular;
	if (api.factor_control !== undefined) result.factor_control = api.factor_control;

	return result;
}

/**
 * Maps API DetalleTuberia to domain DetalleTuberia.
 * Handles optional fields.
 */
function mapApiToDetalleTuberia(api: ApiDetalleTuberia): DetalleTuberia {
	const result: DetalleTuberia = {
		area_fase_mm2: api.area_fase_mm2,
		area_tierra_mm2: api.area_tierra_mm2,
		num_fases_por_tubo: api.num_fases_por_tubo,
		num_neutros_por_tubo: api.num_neutros_por_tubo,
		num_tierras: api.num_tierras,
		area_ocupacion_tubo_mm2: api.area_ocupacion_tubo_mm2,
		designacion_metrica: api.designacion_metrica,
		fill_factor: api.fill_factor
	};

	// Optional fields - add only if present
	if (api.area_neutro_mm2 !== undefined) result.area_neutro_mm2 = api.area_neutro_mm2;

	return result;
}
