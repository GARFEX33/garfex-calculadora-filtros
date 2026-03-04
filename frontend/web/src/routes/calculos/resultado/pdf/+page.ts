/**
 * Página de configuración de PDF.
 * Los datos del cálculo llegan via estado de navegación de SvelteKit ($page.state),
 * accesible desde el componente +page.svelte.
 *
 * Este archivo no exporta un load function porque el estado de navegación
 * no está disponible en load functions — solo en componentes mediante $page.state.
 *
 * Navegación: goto('/calculos/resultado/pdf', { state: { memoria: memoriaOutput } })
 */
export const ssr = false;
