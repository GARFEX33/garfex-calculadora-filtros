/**
 * Material del conductor eléctrico.
 * Según convención NOM: mayúsculas (CU, AL).
 */

export type MaterialConductor = 'CU' | 'AL';

export const MATERIAL_CONDUCTOR_LABELS: Record<MaterialConductor, string> = {
	CU: 'Cobre (CU)',
	AL: 'Aluminio (AL)'
};

/**
 * Verifica si un string es un MaterialConductor válido
 */
export function isMaterialConductor(value: unknown): value is MaterialConductor {
	return typeof value === 'string' && (value === 'CU' || value === 'AL');
}
