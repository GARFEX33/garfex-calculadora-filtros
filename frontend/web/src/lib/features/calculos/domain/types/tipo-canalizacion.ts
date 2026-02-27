/**
 * Tipos de canalización para instalaciones eléctricas.
 * Corresponde a la tabla de capacidades de tubería según NOM-001-SEDE.
 */

export type TipoCanalizacion =
	| 'TUBERIA_PVC'
	| 'TUBERIA_ALUMINIO'
	| 'TUBERIA_ACERO_PG'
	| 'TUBERIA_ACERO_PD'
	| 'CHAROLA_CABLE_ESPACIADO'
	| 'CHAROLA_CABLE_TRIANGULAR';

export const TIPO_CANALIZACION_LABELS: Record<TipoCanalizacion, string> = {
	TUBERIA_PVC: 'Tubería PVC',
	TUBERIA_ALUMINIO: 'Tubería Aluminio',
	TUBERIA_ACERO_PG: 'Tubería Acero PG',
	TUBERIA_ACERO_PD: 'Tubería Acero PD',
	CHAROLA_CABLE_ESPACIADO: 'Charola Cable Espaciado',
	CHAROLA_CABLE_TRIANGULAR: 'Charola Cable Triangular'
};

/**
 * Verifica si un string es un TipoCanalizacion válido
 */
export function isTipoCanalizacion(value: unknown): value is TipoCanalizacion {
	return (
		typeof value === 'string' &&
		(value === 'TUBERIA_PVC' ||
			value === 'TUBERIA_ALUMINIO' ||
			value === 'TUBERIA_ACERO_PG' ||
			value === 'TUBERIA_ACERO_PD' ||
			value === 'CHAROLA_CABLE_ESPACIADO' ||
			value === 'CHAROLA_CABLE_TRIANGULAR')
	);
}
