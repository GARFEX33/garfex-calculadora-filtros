/**
 * Tipos del dominio de cálculos eléctricos.
 * Corresponden a la estructura de datos para el cálculo de memorias técnicas.
 */

import type { TipoCanalizacion } from './tipo-canalizacion.js';
import type { MaterialConductor } from './material-conductor.js';

// Enums / union types
export type { ModoCalculo, TipoFiltro, TipoEquipo, SistemaElectrico, TipoVoltaje, UnidadPotencia, UnidadTension } from './calculo.enums.js';
export type { TipoCanalizacion } from './tipo-canalizacion.js';
export type { MaterialConductor } from './material-conductor.js';

// Labels
export { TIPO_CANALIZACION_LABELS } from './tipo-canalizacion.js';
export { MATERIAL_CONDUCTOR_LABELS } from './material-conductor.js';

// Equipment data (sent in LISTADO mode)
export interface DatosEquipo {
	clave: string;
	tipo: import('./calculo.enums.js').TipoFiltro;
	voltaje: number;
	amperaje: number;
	itm: number;
	bornes?: number;
}

// Request body for POST /api/v1/calculos/memoria
export interface CalcularMemoriaRequest {
	modo: import('./calculo.enums.js').ModoCalculo;
	// LISTADO mode
	equipo?: DatosEquipo;
	// MANUAL modes
	tipo_equipo?: import('./calculo.enums.js').TipoEquipo;
	amperaje_nominal?: number; // MANUAL_AMPERAJE
	potencia_nominal?: number; // MANUAL_POTENCIA
	potencia_unidad?: import('./calculo.enums.js').UnidadPotencia; // MANUAL_POTENCIA
	factor_potencia?: number; // MANUAL_POTENCIA (0-1, required for CARGA)
	// Common installation fields
	tension: number;
	tension_unidad?: import('./calculo.enums.js').UnidadTension;
	itm?: number; // Required in MANUAL modes
	sistema_electrico: import('./calculo.enums.js').SistemaElectrico;
	estado: string;
	tipo_canalizacion: TipoCanalizacion;
	longitud_circuito: number;
	tipo_voltaje: import('./calculo.enums.js').TipoVoltaje;
	// Optional
	hilos_por_fase?: number;
	num_tuberias?: number;
	material?: MaterialConductor;
	porcentaje_caida_maximo?: number;
	temperatura_override?: number;
	diametro_control_mm?: number;
}

/**
 * Alias for CalcularMemoriaRequest — used in application layer.
 * Represents the user input for memoria calculation.
 */
export type MemoriaRequest = CalcularMemoriaRequest;

// ═══════════════════════════════════════════════════════════════════════════
// NUEVA ESTRUCTURA AGRUPADA (Phase 1) — Coincide con backend reorganizado
// ═══════════════════════════════════════════════════════════════════════════

// Datos de instalación — agrupados bajo 'instalacion'
export interface DatosInstalacion {
	tension: number;
	sistema_electrico: import('./calculo.enums.js').SistemaElectrico;
	tipo_canalizacion: string;
	material: string;
	longitud_circuito: number;
	hilos_por_fase: number;
	porcentaje_caida_maximo: number;
}

// Datos de corrientes — agrupados bajo 'corrientes'
export interface DatosCorrientes {
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

// Conductor result — backend serializes in snake_case
export interface ResultadoConductor {
	calibre: string;
	material: string;
	seccion_mm2: number;
	tipo_aislamiento: string;
	capacidad: number;
	num_hilos?: number;
	// Campos para recalculo por caída de tensión
	seleccion_por_caida_tension?: boolean;
	calibre_original_ampacidad?: string;
	nota_seleccion?: string;
}

// Conduit/raceway result — backend serializes in snake_case
export interface ResultadoCanalizacion {
	tamano: string;
	area_total_mm2: number;
	area_requerida_mm2: number;
	numero_de_tubos: number;
	ancho_comercial_mm?: number;
}

// Detalle de charola — valores intermedios del cálculo para el desarrollo en memoria
export interface DetalleCharola {
	// Diámetros
	diametro_fase_mm: number;
	diametro_tierra_mm: number;
	diametro_control_mm?: number;
	// Espaciado
	num_hilos_total?: number;
	espacio_fuerza_mm: number;
	ancho_fuerza_mm?: number;
	espacio_control_mm?: number;
	ancho_control_mm?: number;
	ancho_tierra_mm: number;
	// Triangular
	ancho_potencia_mm?: number;
	factor_triangular?: number;
}

// Detalle de tubería — valores intermedios del cálculo para el desarrollo en memoria
export interface DetalleTuberia {
	// Áreas físicas de conductores (mm²)
	// Fase/Neutro: Tabla 5 NOM (área con aislamiento THW)
	// Tierra: Tabla 8 NOM (conductor desnudo)
	area_fase_mm2: number;
	area_neutro_mm2?: number; // undefined si DELTA (sin neutro)
	area_tierra_mm2: number;
	// Distribución por tubo
	num_fases_por_tubo: number;
	num_neutros_por_tubo: number; // 0 si DELTA
	num_tierras: number;
	// Tubo seleccionado de la tabla NOM
	area_ocupacion_tubo_mm2: number; // área de ocupación del CSV (40% del interior ya aplicado)
	designacion_metrica: string; // ej: "63" → mostrar como "63 mm"
	fill_factor: number;
}

// Canalización agrupada — resultado + detalles
export interface DatosCanalizacion {
	resultado: ResultadoCanalizacion;
	fill_factor: number;
	detalle_charola?: DetalleCharola;
	detalle_tuberia?: DetalleTuberia;
}

// Datos de protección
export interface DatosProteccion {
	itm: number;
}

// Voltage drop result
export interface ResultadoCaidaTension {
	porcentaje: number;
	caida_volts: number;
	cumple: boolean;
	limite_porcentaje: number;
	impedancia: number;
	resistencia: number;
	reactancia: number;
}

// Full calculation result (from POST response data field)
// Nueva estructura agrupada por entidad — coincide con backend reorganizado
export interface MemoriaOutput {
	// Datos del equipo
	equipo: DatosEquipo;
	tipo_equipo: string;
	factor_potencia: number;
	estado: string;

	// Parámetros de instalación
	instalacion: DatosInstalacion;

	// Cálculos de corriente
	corrientes: DatosCorrientes;

	// Conductores
	cable_fase: ResultadoConductor;
	cable_neutro?: ResultadoConductor;
	cable_tierra: ResultadoConductor;

	// Canalización
	canalizacion: DatosCanalizacion;

	// Protección
	proteccion: DatosProteccion;

	// Caída de tensión
	caida_tension: ResultadoCaidaTension;

	// Resumen y metadatos
	cumple_normativa: boolean;
	observaciones: string[];
	pasos: unknown[];
}

// API response wrapper (backend wraps in { success, data })
export interface CalcularMemoriaResponse {
	success: boolean;
	data: MemoriaOutput;
}
