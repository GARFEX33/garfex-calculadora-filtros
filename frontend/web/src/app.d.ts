// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
import type { MemoriaOutput } from '$lib/types/calculos.types';

declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		interface PageState {
			/** MemoriaOutput para la página de configuración PDF */
			memoria?: MemoriaOutput;
		}
		// interface Platform {}
	}
}

export {};
