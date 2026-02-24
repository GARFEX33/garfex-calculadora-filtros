<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';
	import { cn } from '$lib/utils';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// Derived values
	let cumpleNormativa = $derived(memoria.cumple_normativa ?? false);
	let observaciones = $derived(memoria.observaciones ?? []);
	let tieneObservaciones = $derived(observaciones.length > 0);

	// Check individual criteria
	let cumpleAmpacidad = $derived(
		memoria.conductor_alimentacion.Capacidad >= memoria.corriente_nominal * 1.25
	);
	let cumpleCaida = $derived(memoria.caida_tension?.cumple ?? false);
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		7. Conclusión Técnica
	</h2>

	<!-- Summary status -->
	<div
		class={cn(
			'rounded-lg border p-6 text-center',
			cumpleNormativa
				? 'border-success/50 bg-success/10'
				: 'border-destructive/50 bg-destructive/10'
		)}
	>
		<h3 class={cn('text-lg font-bold', cumpleNormativa ? 'text-success' : 'text-destructive')}>
			{cumpleNormativa
				? 'EL DISEÑO CUMPLE CON LA NOM-001-SEDE-2012'
				: 'EL DISEÑO REQUIERE REVISIÓN'}
		</h3>
	</div>

	<!-- Criteria checklist -->
	<div class="mt-6">
		<h4 class="mb-3 font-semibold text-foreground">Criterios de Cumplimiento</h4>
		<div class="space-y-2">
			<div
				class={cn(
					'flex items-center rounded border p-3',
					cumpleAmpacidad
						? 'border-success/30 bg-success/5'
						: 'border-destructive/30 bg-destructive/5'
				)}
			>
				<span class="mr-2 {cumpleAmpacidad ? 'text-success' : 'text-destructive'}">
					{cumpleAmpacidad ? '✓' : '✗'}
				</span>
				<span class="text-sm text-foreground"> Ampacidad del conductor de alimentación </span>
			</div>
			<div
				class={cn(
					'flex items-center rounded border p-3',
					cumpleCaida ? 'border-success/30 bg-success/5' : 'border-destructive/30 bg-destructive/5'
				)}
			>
				<span class="mr-2 {cumpleCaida ? 'text-success' : 'text-destructive'}">
					{cumpleCaida ? '✓' : '✗'}
				</span>
				<span class="text-sm text-foreground">
					Caída de tensión dentro de límites ({memoria.caida_tension?.limite_porcentaje ?? 3}%)
				</span>
			</div>
		</div>
	</div>

	<!-- Observations -->
	{#if tieneObservaciones}
		<div class="mt-6">
			<h4 class="mb-3 font-semibold text-foreground">Observaciones</h4>
			<ul class="list-inside list-disc space-y-1 rounded border border-border bg-muted p-4 text-sm">
				{#each observaciones as obs}
					<li class="text-foreground">{obs}</li>
				{/each}
			</ul>
		</div>
	{/if}

	<!-- Final statement -->
	<div class="mt-6 rounded border border-primary/30 bg-primary/10 p-4">
		<p class="text-sm text-foreground">
			<strong>Resumen:</strong> Se aplicaron las fórmulas técnicas correctas para el ejemplo
			presentado. El sistema{' '}
			{cumpleAmpacidad ? 'cumple' : 'no cumple'} los criterios de ampacidad, y la{' '}
			{cumpleCaida ? 'caída de tensión está dentro' : 'caída de tensión excede'} los límites recomendados.
			El diseño {cumpleNormativa ? 'cumple' : 'no cumple'} con la NOM-001-SEDE-2012.
		</p>
	</div>

	<!-- Signature area -->
	<div class="mt-8 grid grid-cols-2 gap-8 pt-8">
		<div class="border-t border-border pt-2">
			<p class="text-xs text-muted-foreground">Elaboró</p>
			<div class="mt-4 h-8"></div>
		</div>
		<div class="border-t border-border pt-2">
			<p class="text-xs text-muted-foreground">Revisó</p>
			<div class="mt-4 h-8"></div>
		</div>
	</div>
</section>
