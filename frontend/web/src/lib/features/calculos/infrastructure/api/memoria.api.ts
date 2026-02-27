/**
 * API layer for memoria calculations.
 * Handles HTTP communication with the backend.
 * NO business logic — just transport.
 */

import { apiClient } from '$lib/shared/api/client';
import type { ApiResult } from '$lib/shared/types/api.types';

/**
 * Raw API request for POST /api/v1/calculos/memoria.
 * Matches the exact format expected by the backend.
 */
export interface ApiMemoriaRequest {
	modo: 'LISTADO' | 'MANUAL_AMPERAJE' | 'MANUAL_POTENCIA';
	// LISTADO mode
	equipo?: ApiDatosEquipo;
	// MANUAL modes
	tipo_equipo?: 'FILTRO_ACTIVO' | 'TRANSFORMADOR' | 'FILTRO_RECHAZO' | 'CARGA';
	amperaje_nominal?: number;
	potencia_nominal?: number;
	potencia_unidad?: 'W' | 'KW' | 'KVA' | 'KVAR';
	factor_potencia?: number;
	// Common
	tension: number;
	tension_unidad?: 'V' | 'kV';
	itm?: number;
	sistema_electrico: 'DELTA' | 'ESTRELLA' | 'BIFASICO' | 'MONOFASICO';
	estado: string;
	tipo_canalizacion:
		| 'TUBERIA_PVC'
		| 'TUBERIA_ALUMINIO'
		| 'TUBERIA_ACERO_PG'
		| 'TUBERIA_ACERO_PD'
		| 'CHAROLA_CABLE_ESPACIADO'
		| 'CHAROLA_CABLE_TRIANGULAR';
	longitud_circuito: number;
	tipo_voltaje: 'FASE_NEUTRO' | 'FASE_FASE';
	hilos_por_fase?: number;
	num_tuberias?: number;
	material?: 'CU' | 'AL';
	porcentaje_caida_maximo?: number;
	temperatura_override?: number;
	diametro_control_mm?: number;
}

/**
 * Equipment data in API format.
 */
export interface ApiDatosEquipo {
	clave: string;
	tipo: 'A' | 'KVA' | 'KVAR';
	voltaje: number;
	amperaje: number;
	itm: number;
	bornes?: number;
}

/**
 * Datos de instalación en formato API.
 * Agrupados bajo 'instalacion' en la respuesta.
 */
export interface ApiDatosInstalacion {
	tension: number;
	sistema_electrico: 'DELTA' | 'ESTRELLA' | 'BIFASICO' | 'MONOFASICO';
	tipo_canalizacion: string;
	material: string;
	longitud_circuito: number;
	hilos_por_fase: number;
	porcentaje_caida_maximo: number;
}

/**
 * Datos de corrientes en formato API.
 * Agrupados bajo 'corrientes' en la respuesta.
 */
export interface ApiDatosCorrientes {
	corriente_nominal: number;
	corriente_ajustada: number;
	corriente_por_hilo: number;
	factor_temperatura: number;
	factor_agrupamiento: number;
	factor_total_ajuste: number;
	temperatura_ambiente: number;
	temperatura_referencia: number;
	conductores_por_tubo: number;
	cantidad_conductores: number;
	tabla_ampacidad_usada: string;
}

/**
 * Conductor result from API — backend serializes in snake_case.
 */
export interface ApiResultadoConductor {
	calibre: string;
	material: string;
	seccion_mm2: number;
	tipo_aislamiento: string;
	capacidad: number;
	num_hilos?: number;
	seleccion_por_caida_tension?: boolean;
	calibre_original_ampacidad?: string;
	nota_seleccion?: string;
}

/**
 * Conduit result from API.
 */
export interface ApiResultadoCanalizacion {
	tamano: string;
	area_total_mm2: number;
	area_requerida_mm2: number;
	numero_de_tubos: number;
	ancho_comercial_mm?: number;
}

/**
 * Charola detail from API.
 */
export interface ApiDetalleCharola {
	diametro_fase_mm: number;
	diametro_tierra_mm: number;
	diametro_control_mm?: number;
	num_hilos_total?: number;
	espacio_fuerza_mm: number;
	ancho_fuerza_mm?: number;
	espacio_control_mm?: number;
	ancho_control_mm?: number;
	ancho_tierra_mm: number;
	ancho_potencia_mm?: number;
	factor_triangular?: number;
}

/**
 * Tuberia detail from API.
 */
export interface ApiDetalleTuberia {
	area_fase_mm2: number;
	area_neutro_mm2?: number;
	area_tierra_mm2: number;
	num_fases_por_tubo: number;
	num_neutros_por_tubo: number;
	num_tierras: number;
	area_ocupacion_tubo_mm2: number;
	designacion_metrica: string;
	fill_factor: number;
}

/**
 * Datos de canalización agrupados (resultado + detalles).
 */
export interface ApiDatosCanalizacion {
	resultado: ApiResultadoCanalizacion;
	fill_factor: number;
	detalle_charola?: ApiDetalleCharola;
	detalle_tuberia?: ApiDetalleTuberia;
}

/**
 * Datos de protección en formato API.
 */
export interface ApiDatosProteccion {
	itm: number;
}

/**
 * Voltage drop result from API.
 */
export interface ApiResultadoCaidaTension {
	porcentaje: number;
	caida_volts: number;
	cumple: boolean;
	limite_porcentaje: number;
	impedancia: number;
	resistencia: number;
	reactancia: number;
}

/**
 * Paso de memoria - detalle del cálculo.
 */
export interface ApiPasoMemoria {
	numero: number;
	nombre: string;
	descripcion: string;
	resultado: unknown;
}

/**
 * Full calculation output from API response data field.
 * Uses snake_case as returned by backend.
 * Nueva estructura agrupada por entidad.
 */
export interface ApiMemoriaOutput {
	// ═══════════════════════════════════════════════════════════════════════
	// DATOS DEL EQUIPO
	// ═══════════════════════════════════════════════════════════════════════
	equipo: ApiDatosEquipo;
	tipo_equipo: string;
	factor_potencia: number;
	estado: string;

	// ═══════════════════════════════════════════════════════════════════════
	// PARÁMETROS DE INSTALACIÓN
	// ═══════════════════════════════════════════════════════════════════════
	instalacion: ApiDatosInstalacion;

	// ═══════════════════════════════════════════════════════════════════════
	// CÁLCULOS DE CORRIENTE
	// ═══════════════════════════════════════════════════════════════════════
	corrientes: ApiDatosCorrientes;

	// ═══════════════════════════════════════════════════════════════════════
	// CONDUCTORES
	// ═══════════════════════════════════════════════════════════════════════
	cable_fase: ApiResultadoConductor;
	cable_neutro?: ApiResultadoConductor; // nil para sistemas DELTA
	cable_tierra: ApiResultadoConductor;

	// ═══════════════════════════════════════════════════════════════════════
	// CANALIZACIÓN
	// ═══════════════════════════════════════════════════════════════════════
	canalizacion: ApiDatosCanalizacion;

	// ═══════════════════════════════════════════════════════════════════════
	// PROTECCIÓN
	// ═══════════════════════════════════════════════════════════════════════
	proteccion: ApiDatosProteccion;

	// ═══════════════════════════════════════════════════════════════════════
	// CAÍDA DE TENSIÓN
	// ═══════════════════════════════════════════════════════════════════════
	caida_tension: ApiResultadoCaidaTension;

	// ═══════════════════════════════════════════════════════════════════════
	// RESUMEN Y METADATOS
	// ═══════════════════════════════════════════════════════════════════════
	cumple_normativa: boolean;
	observaciones: string[];
	pasos: ApiPasoMemoria[];
}

/**
 * API response wrapper.
 */
export interface ApiMemoriaResponse {
	success: boolean;
	data: ApiMemoriaOutput;
}

/**
 * Calculates a memoria técnica.
 *
 * @param input - Request data in API format
 * @returns Result with response data or error
 */
export async function calcularMemoria(input: ApiMemoriaRequest): Promise<ApiResult<ApiMemoriaResponse>> {
	return apiClient.post<ApiMemoriaResponse>('/api/v1/calculos/memoria', input);
}
