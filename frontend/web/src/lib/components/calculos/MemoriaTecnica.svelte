<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';
	import { cn } from '$lib/utils';

	import SeccionEncabezado from './secciones/SeccionEncabezado.svelte';
	import SeccionCorriente from './secciones/SeccionCorriente.svelte';
	import SeccionAlimentador from './secciones/SeccionAlimentador.svelte';
	import SeccionTierra from './secciones/SeccionTierra.svelte';
	import SeccionCaidaTension from './secciones/SeccionCaidaTension.svelte';
	import SeccionCanalizacion from './secciones/SeccionCanalizacion.svelte';
	import SeccionConclusion from './secciones/SeccionConclusion.svelte';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// Derived values
	let cumpleNormativa = $derived(memoria.cumple_normativa ?? false);
</script>

<div class="flex flex-col gap-6">
	<!-- Status Banner -->
	<div
		class={cn(
			'rounded-lg border p-4 text-center',
			cumpleNormativa
				? 'border-success/50 bg-success/10 text-success'
				: 'border-destructive/50 bg-destructive/10 text-destructive'
		)}
	>
		<span class="text-lg font-semibold">
			{cumpleNormativa ? '✓ CUMPLE CON NOM-001-SEDE-2012' : '⚠ REQUIERE REVISIÓN'}
		</span>
	</div>

	<!-- 1. Encabezado -->
	<SeccionEncabezado {memoria} />

	<!-- 2. Cálculo de Corriente Nominal -->
	<SeccionCorriente {memoria} />

	<!-- 3. Dimensionamiento del Alimentador -->
	<SeccionAlimentador {memoria} />

	<!-- 4. Conductor de Puesta a Tierra -->
	<SeccionTierra {memoria} />

	<!-- 5. Cálculo de Caída de Tensión -->
	<SeccionCaidaTension {memoria} />

	<!-- 6. Cálculo de Canalización -->
	<SeccionCanalizacion {memoria} />

	<!-- 7. Conclusión Técnica -->
	<SeccionConclusion {memoria} />
</div>
