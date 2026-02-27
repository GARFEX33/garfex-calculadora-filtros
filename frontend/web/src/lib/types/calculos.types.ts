// Enums / union types
export type ModoCalculo = 'LISTADO' | 'MANUAL_AMPERAJE' | 'MANUAL_POTENCIA';
export type TipoFiltro = 'A' | 'KVA' | 'KVAR';
export type TipoEquipo = 'FILTRO_ACTIVO' | 'TRANSFORMADOR' | 'FILTRO_RECHAZO' | 'CARGA';
export type SistemaElectrico = 'DELTA' | 'ESTRELLA' | 'BIFASICO' | 'MONOFASICO';
export type TipoVoltaje = 'FASE_NEUTRO' | 'FASE_FASE';
export type TipoCanalizacion =
	| 'TUBERIA_PVC'
	| 'TUBERIA_ACERO_PG'
	| 'TUBERIA_ACERO_PD'
	| 'CHAROLA_CABLE_ESPACIADO'
	| 'CHAROLA_CABLE_TRIANGULAR';
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
export interface MemoriaOutput {
	equipo: DatosEquipo;
	tipo_equipo: string;
	tension: number;
	factor_potencia: number;
	estado: string;
	temperatura_ambiente: number;
	sistema_electrico: SistemaElectrico;
	cantidad_conductores: number;
	conductores_por_tubo?: number;
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
	detalle_charola?: DetalleCharola;
	detalle_tuberia?: DetalleTuberia;
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
