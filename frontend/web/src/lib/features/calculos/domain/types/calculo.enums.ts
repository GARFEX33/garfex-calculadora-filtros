/**
 * Tipos de modo de cálculo.
 */
export type ModoCalculo = 'LISTADO' | 'MANUAL_AMPERAJE' | 'MANUAL_POTENCIA';

/**
 * Tipos de filtro según potencia.
 */
export type TipoFiltro = 'A' | 'KVA' | 'KVAR';

/**
 * Tipos de equipo eléctrico.
 */
export type TipoEquipo = 'FILTRO_ACTIVO' | 'TRANSFORMADOR' | 'FILTRO_RECHAZO' | 'CARGA';

/**
 * Sistemas eléctricos soportados.
 */
export type SistemaElectrico = 'DELTA' | 'ESTRELLA' | 'BIFASICO' | 'MONOFASICO';

/**
 * Tipos de voltaje.
 */
export type TipoVoltaje = 'FASE_NEUTRO' | 'FASE_FASE';

/**
 * Unidades de potencia.
 */
export type UnidadPotencia = 'W' | 'KW' | 'KVA' | 'KVAR';

/**
 * Unidades de tensión.
 */
export type UnidadTension = 'V' | 'kV';
