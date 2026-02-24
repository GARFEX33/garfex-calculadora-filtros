<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';
	import { cn } from '$lib/utils';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// Determine if three-phase or single-phase
	let esTrifasico = $derived(
		memoria.sistema_electrico === 'ESTRELLA' || memoria.sistema_electrico === 'DELTA'
	);

	// Get voltage drop data
	let caida = $derived(memoria.caida_tension);
	let cumple = $derived(caida?.cumple ?? false);
	let porcentaje = $derived(caida?.porcentaje ?? 0);
	let limite = $derived(caida?.limite_porcentaje ?? 3.0);
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		6. Cálculo de Caída de Tensión
	</h2>

	<p class="mb-4 text-sm text-muted-foreground">
		La caída de tensión se calcula considerando la impedancia del conductor, la longitud del
		circuito y el tipo de sistema eléctrico.
	</p>

	<!-- Fórmula -->
	<div class="mb-6 rounded bg-muted p-4">
		<h3 class="mb-2 font-semibold text-foreground">Fórmula</h3>
		{#if esTrifasico}
			<p class="font-mono text-foreground">%V = (√3 × I × Z × L / V) × 100</p>
			<p class="mt-1 text-sm text-muted-foreground">
				Sistema trifásico: factor √3, Z = impedancia por conductor (Ω/km)
			</p>
		{:else}
			<p class="font-mono text-foreground">%V = (2 × I × Z × L / V) × 100</p>
			<p class="mt-1 text-sm text-muted-foreground">
				Sistema monofásico: factor 2 (ida y retorno), Z = impedancia (Ω/km)
			</p>
		{/if}
	</div>

	<!-- Parámetros -->
	<div class="mb-6">
		<h3 class="mb-2 font-semibold text-foreground">Parámetros</h3>
		<div class="grid grid-cols-2 gap-4 text-sm">
			<div>
				<span class="text-muted-foreground">Corriente:</span>
				<span class="ml-2 font-mono text-foreground">{memoria.corriente_nominal.toFixed(2)} A</span>
			</div>
			<div>
				<span class="text-muted-foreground">Longitud:</span>
				<span class="ml-2 font-mono text-foreground">{memoria.longitud_circuito.toFixed(2)} m</span>
			</div>
			<div>
				<span class="text-muted-foreground">Voltaje:</span>
				<span class="ml-2 font-mono text-foreground">{memoria.tension} V</span>
			</div>
			<div>
				<span class="text-muted-foreground">Impedancia (Zef):</span>
				<span class="ml-2 font-mono text-foreground">
					{caida?.impedancia?.toFixed(4) ?? '—'} Ω/km
				</span>
			</div>
			<div>
				<span class="text-muted-foreground">Factor de Potencia:</span>
				<span class="ml-2 font-mono text-foreground">
					{memoria.factor_potencia?.toFixed(2) ?? '—'}
				</span>
			</div>
		</div>
	</div>

	<!-- Desarrollo -->
	<div class="mb-6">
		<h3 class="mb-2 font-semibold text-foreground">Desarrollo</h3>
		<div class="rounded bg-muted p-3 font-mono text-sm">
			{#if esTrifasico}
				<p class="text-foreground">
					%V = (√3 × {memoria.corriente_nominal.toFixed(2)} × {caida?.impedancia?.toFixed(4) ?? '—'} ×
					{memoria.longitud_circuito.toFixed(2)} / {memoria.tension}) × 100
				</p>
			{:else}
				<p class="text-foreground">
					%V = (2 × {memoria.corriente_nominal.toFixed(2)} × {caida?.impedancia?.toFixed(4) ?? '—'} ×
					{memoria.longitud_circuito.toFixed(2)} / {memoria.tension}) × 100
				</p>
			{/if}
		</div>
	</div>

	<!-- Resultado -->
	<div
		class={cn(
			'rounded border p-4',
			cumple ? 'border-success/30 bg-success/10' : 'border-destructive/30 bg-destructive/10'
		)}
	>
		<h3 class="mb-2 font-semibold text-foreground">Resultado</h3>
		<div class="grid grid-cols-2 gap-4">
			<div>
				<p class="text-sm text-muted-foreground">Caída de Tensión</p>
				<p class={cn('text-2xl font-bold', cumple ? 'text-success' : 'text-destructive')}>
					{porcentaje.toFixed(2)}%
				</p>
			</div>
			<div>
				<p class="text-sm text-muted-foreground">Caída en Volts</p>
				<p class="text-2xl font-bold text-foreground">
					{caida?.caida_volts?.toFixed(2) ?? '—'} V
				</p>
			</div>
		</div>
	</div>

	<!-- Verificación -->
	<div class="mt-4">
		<h3 class="mb-2 font-semibold text-foreground">Verificación</h3>
		<div
			class={cn(
				'rounded border p-3',
				cumple ? 'border-success/30 bg-success/5' : 'border-destructive/30 bg-destructive/5'
			)}
		>
			<p class="text-sm">
				{#if cumple}
					<span class="font-medium text-success">✓ CUMPLE</span>
					<span class="text-foreground">
						- La caída de tensión ({porcentaje.toFixed(2)}%) está dentro del límite permitido ({limite.toFixed(
							1
						)}%)
					</span>
				{:else}
					<span class="font-medium text-destructive">✗ NO CUMPLE</span>
					<span class="text-foreground">
						- La caída de tensión ({porcentaje.toFixed(2)}%) excede el límite permitido ({limite.toFixed(
							1
						)}%)
					</span>
				{/if}
			</p>
		</div>
	</div>
</section>
