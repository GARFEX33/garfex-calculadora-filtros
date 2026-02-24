// Enums / union types
export type ModoCalculo = 'LISTADO' | 'MANUAL_AMPERAJE' | 'MANUAL_POTENCIA';
export type TipoFiltro = 'A' | 'KVA' | 'KVAR';
export type TipoEquipo = 'FILTRO_ACTIVO' | 'TRANSFORMADOR' | 'FILTRO_RECHAZO' | 'CARGA';
export type SistemaElectrico = 'DELTA' | 'ESTRELLA' | 'BIFASICO' | 'MONOFASICO';
export type TipoVoltaje = 'FASE_NEUTRO' | 'FASE_FASE';
export type TipoCanalizacion =
	| 'TUBERIA_PVC'
	| 'TUBERIA_EMT'
	| 'CHAROLA_CABLE_ESPACIADO'
	| 'CHAROLA_CABLE_TRESBOLILLO';
export type MaterialConductor = 'Cu' | 'Al';
export type UnidadPotencia = 'W' | 'KW' | 'KVA' | 'KVAR';
export type UnidadTension = 'V' | 'kV';

// Equipment data (sent in LISTADO mode)
export interface DatosEquipo {
	clave: string;
	tipo: TipoFiltro;
	voltaje: number;
	amperaje: number;
	itm: number;
	bornes?: number;
}

// Request body for POST /api/v1/calculos/memoria
export interface CalcularMemoriaRequest {
	modo: ModoCalculo;
	// LISTADO mode
	equipo?: DatosEquipo;
	// MANUAL modes
	tipo_equipo?: TipoEquipo;
	amperaje_nominal?: number; // MANUAL_AMPERAJE
	potencia_nominal?: number; // MANUAL_POTENCIA
	potencia_unidad?: UnidadPotencia; // MANUAL_POTENCIA
	factor_potencia?: number; // MANUAL_POTENCIA (0-1, required for CARGA)
	// Common installation fields
	tension: number;
	tension_unidad?: UnidadTension;
	itm?: number; // Required in MANUAL modes
	sistema_electrico: SistemaElectrico;
	estado: string;
	tipo_canalizacion: TipoCanalizacion;
	longitud_circuito: number;
	tipo_voltaje: TipoVoltaje;
	// Optional
	hilos_por_fase?: number;
	num_tuberias?: number;
	material?: MaterialConductor;
	porcentaje_caida_maximo?: number;
	temperatura_override?: number;
	diametro_control_mm?: number;
}

// Conductor result — backend serializes in PascalCase
export interface ResultadoConductor {
	Calibre: string;
	Material: string;
	SeccionMM2: number;
	TipoAislamiento: string;
	Capacidad: number;
	NumHilos?: number;
}

// Conduit/raceway result — backend serializes in PascalCase
export interface ResultadoCanalizacion {
	Tamano: string;
	AreaTotalMM2: number;
	AreaRequeridaMM2: number;
	NumeroDeTubos: number;
}

// Voltage drop result
export interface ResultadoCaidaTension {
	porcentaje: number;
	caida_volts: number;
	cumple: boolean;
	limite_porcentaje: number;
	impedancia: number;
}

// Full calculation result (from POST response data field)
export interface MemoriaOutput {
	equipo: DatosEquipo;
	tipo_equipo: string;
	tension: number;
	factor_potencia: number;
	estado: string;
	temperatura_ambiente: number;
	sistema_electrico: SistemaElectrico;
	cantidad_conductores: number;
	corriente_nominal: number;
	corriente_ajustada: number;
	factor_temperatura: number;
	factor_agrupamiento: number;
	factor_total_ajuste: number;
	hilos_por_fase: number;
	corriente_por_hilo: number;
	tipo_canalizacion: string;
	material: string;
	temperatura_usada: number;
	conductor_alimentacion: ResultadoConductor;
	tabla_ampacidad_usada: string;
	conductor_tierra: ResultadoConductor;
	itm: number;
	canalizacion: ResultadoCanalizacion;
	fill_factor: number;
	longitud_circuito: number;
	caida_tension: ResultadoCaidaTension;
	cumple_normativa: boolean;
	observaciones: string[];
}

// API response wrapper (backend wraps in { success, data })
export interface CalcularMemoriaResponse {
	success: boolean;
	data: MemoriaOutput;
}
