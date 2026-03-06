/**
 * Catálogo estático de empresas para la memoria de cálculo PDF.
 * Espeja el catálogo del backend (internal/pdf/domain/empresa.go).
 * En MVP los datos son placeholder — actualizar con datos reales antes de producción.
 */

export interface EmpresaPdf {
	id: string;
	nombre: string;
	direccion: string;
	telefono: string;
	email: string;
}

export const EMPRESAS_PDF: EmpresaPdf[] = [
	{
		id: 'garfex',
		nombre: 'Garfex Ingeniería Eléctrica S.A. de C.V.',
		direccion: 'Av. Insurgentes Sur 1234, Col. Del Valle, CDMX, C.P. 03100',
		telefono: '+52 55 1193-0515',
		email: 'jcgarcia@garfex.mx'
	},
	{
		id: 'summaa',
		nombre: 'Summa Ingeniería Eléctrica S.A. de C.V.',
		direccion: 'Blvd. Manuel Ávila Camacho 800, Lomas de Chapultepec, CDMX, C.P. 11000',
		telefono: '+52 55 9876-5432',
		email: 'contacto@summa.mx'
	},
	{
		id: 'siemens',
		nombre: 'Siemens S.A. de C.V.',
		direccion: 'Lago Alberto 319, Anáhuac I Secc, Miguel Hidalgo, CDMX, C.P. 11320',
		telefono: '+52 55 5229-3600',
		email: 'contacto@siemens.com.mx'
	}
];

/** ID de la empresa seleccionada por defecto */
export const EMPRESA_DEFAULT_ID = 'garfex';
