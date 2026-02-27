<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';
	import { cn } from '$lib/utils';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// Sistema eléctrico
	let esMonofasico = $derived(memoria.sistema_electrico === 'MONOFASICO');
	let esBifasico = $derived(memoria.sistema_electrico === 'BIFASICO');
	let esEstrella = $derived(memoria.sistema_electrico === 'ESTRELLA');
	let esDelta = $derived(memoria.sistema_electrico === 'DELTA');

	// Voltage reference: Mono/Bifasico → Vfn; Estrella/Delta → Vff
	let voltajeRefLabel = $derived(esMonofasico || esBifasico ? 'Vfn' : 'Vff');
	let voltajeRefDesc = $derived(esMonofasico || esBifasico ? 'fase-neutro' : 'fase-fase');

	// Zef components
	let cosTheta = $derived(memoria.factor_potencia);
	let senTheta = $derived(Math.sqrt(1 - cosTheta * cosTheta));

	// caida shorthand
	let caida = $derived(memoria.caida_tension);
	let cumple = $derived(caida?.cumple ?? false);
	let porcentaje = $derived(caida?.porcentaje ?? 0);
	let limite = $derived(caida?.limite_porcentaje ?? 3.0);

	// Recálculo por caída de tensión
	let recalculoPorCaida = $derived(
		memoria.conductor_alimentacion.seleccion_por_caida_tension === true
	);
	let calibreOriginal = $derived(memoria.conductor_alimentacion.calibre_original_ampacidad ?? '');
	let calibreFinal = $derived(memoria.conductor_alimentacion.calibre);
	let notaRecalculo = $derived(memoria.conductor_alimentacion.nota_seleccion ?? '');

	// Caso donde se agotaron todos los calibres disponibles y ninguno cumple
	let agotadoCalibre = $derived(!cumple && !recalculoPorCaida);

	// Valores para desarrollo
	let R = $derived(caida?.resistencia?.toFixed(4) ?? '—');
	let X = $derived(caida?.reactancia?.toFixed(4) ?? '—');
	let Zef = $derived(caida?.impedancia?.toFixed(4) ?? '—');
	let I = $derived(memoria.corriente_nominal.toFixed(2));
	let L = $derived(memoria.longitud_circuito.toFixed(2));
	let V = $derived(memoria.tension);

	// Derived values for Desarrollo section (Task 3.4)
	let L_km = $derived((memoria.longitud_circuito / 1000).toFixed(3));
	let caida_volts = $derived(caida?.caida_volts?.toFixed(4) ?? '—');
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		6. Cálculo de Caída de Tensión
	</h2>

	<!-- Descripción -->
	<p class="mb-4 text-sm text-muted-foreground">
		La caída de tensión se calcula considerando la impedancia efectiva del conductor, la longitud
		del circuito y el tipo de sistema eléctrico. Para sistemas trifásicos se utiliza el factor √3,
		mientras que para monofásicos se emplea el factor 2 (ida y retorno).
	</p>

	<!-- Fórmula -->
	<div class="mb-6 rounded bg-muted p-4">
		<h3 class="mb-2 font-semibold text-foreground">Fórmula</h3>
		{#if esMonofasico}
			<p class="font-mono text-foreground">e = 2 × I × Zef × L</p>
			<p class="font-mono text-sm text-muted-foreground">e% = (e / Vfn) × 100</p>
		{:else if esBifasico}
			<p class="font-mono text-foreground">e = 1 × I × Zef × L</p>
			<p class="font-mono text-sm text-muted-foreground">e% = (e / Vfn) × 100</p>
		{:else if esEstrella}
			<p class="font-mono text-foreground">e = √3 × I × Zef × L</p>
			<p class="font-mono text-sm text-muted-foreground">e% = (e / Vff) × 100</p>
		{:else if esDelta}
			<p class="font-mono text-foreground">e = √3 × I × Zef × L</p>
			<p class="font-mono text-sm text-muted-foreground">e% = (e / Vff) × 100</p>
		{/if}
	</div>

	<!-- Impedancia Efectiva (Zef) -->
	<div class="mb-6 rounded border border-border bg-card p-4">
		<h3 class="mb-2 font-semibold text-card-foreground">Impedancia Efectiva (Zef)</h3>
		<p class="font-mono text-foreground">Zef = R × cosθ + X × sinθ</p>
		<p class="mt-1 font-mono text-sm text-muted-foreground">
			Zef = {R} × {cosTheta.toFixed(3)} + {X} × {senTheta.toFixed(3)} = {Zef} Ω/km
		</p>
	</div>

	<!-- Parámetros -->
	<div class="mb-6">
		<h3 class="mb-2 font-semibold text-foreground">Parámetros</h3>
		<div class="grid grid-cols-2 gap-4 text-sm">
			<div>
				<span class="text-muted-foreground">Corriente:</span>
				<span class="ml-2 font-mono text-foreground">{I} A</span>
			</div>
			<div>
				<span class="text-muted-foreground">Longitud:</span>
				<span class="ml-2 font-mono text-foreground">{L} m</span>
			</div>
			<div>
				<span class="text-muted-foreground">Voltaje:</span>
				<span class="ml-2 font-mono text-foreground">{V} V</span>
			</div>
			<div>
				<span class="text-muted-foreground">Voltaje referencia:</span>
				<span class="ml-2 font-mono text-foreground">{voltajeRefLabel} ({voltajeRefDesc})</span>
			</div>
			<div>
				<span class="text-muted-foreground">Resistencia (R):</span>
				<span class="ml-2 font-mono text-foreground">{R} Ω/km</span>
			</div>
			<div>
				<span class="text-muted-foreground">Reactancia (X):</span>
				<span class="ml-2 font-mono text-foreground">{X} Ω/km</span>
			</div>
			<div>
				<span class="text-muted-foreground">Impedancia (Zef):</span>
				<span class="ml-2 font-mono text-foreground">{Zef} Ω/km</span>
			</div>
			<div>
				<span class="text-muted-foreground">Factor de Potencia:</span>
				<span class="ml-2 font-mono text-foreground">{cosTheta.toFixed(2)}</span>
			</div>
		</div>
	</div>

	<!-- Desarrollo -->
	<div class="mb-6">
		<h3 class="mb-2 font-semibold text-foreground">Desarrollo</h3>
		<div class="space-y-1 rounded bg-muted p-3 font-mono text-sm">
			{#if esMonofasico}
				<p class="text-foreground">e = 2 × {I} A × {Zef} Ω/km × {L_km} km = {caida_volts} V</p>
				<p class="text-muted-foreground">
					e% = ({caida_volts} / {V}) × 100 = {porcentaje.toFixed(2)}%
				</p>
			{:else if esBifasico}
				<p class="text-foreground">e = 1 × {I} A × {Zef} Ω/km × {L_km} km = {caida_volts} V</p>
				<p class="text-muted-foreground">
					e% = ({caida_volts} / {V}) × 100 = {porcentaje.toFixed(2)}%
				</p>
			{:else if esEstrella}
				<p class="text-foreground">e = √3 × {I} A × {Zef} Ω/km × {L_km} km = {caida_volts} V</p>
				<p class="text-muted-foreground">
					e% = ({caida_volts} / {V}) × 100 = {porcentaje.toFixed(2)}%
				</p>
			{:else if esDelta}
				<p class="text-foreground">e = √3 × {I} A × {Zef} Ω/km × {L_km} km = {caida_volts} V</p>
				<p class="text-muted-foreground">
					e% = ({caida_volts} / {V}) × 100 = {porcentaje.toFixed(2)}%
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

	<!-- Recálculo por caída de tensión — calibre aumentado exitosamente -->
	{#if recalculoPorCaida}
		<div class="mt-4 rounded border border-warning/40 bg-warning/10 p-4">
			<h3 class="mb-2 flex items-center gap-2 font-semibold text-foreground">
				<span class="rounded bg-warning px-2 py-0.5 text-xs font-medium text-warning-foreground">
					Recálculo NOM-001-SEDE
				</span>
				Ajuste automático de calibre por caída de tensión
			</h3>
			<p class="mb-2 text-sm text-muted-foreground">
				El calibre seleccionado por ampacidad no cumplió con el límite de caída de tensión. Se
				aumentó automáticamente al siguiente calibre superior que sí cumple.
			</p>
			{#if calibreOriginal}
				<div class="flex items-center gap-3 font-mono text-sm">
					<span class="rounded bg-destructive/10 px-2 py-1 text-destructive"
						>{calibreOriginal} AWG</span
					>
					<span class="text-muted-foreground">→</span>
					<span class="rounded bg-success/10 px-2 py-1 font-semibold text-success"
						>{calibreFinal} AWG</span
					>
				</div>
			{/if}
			{#if notaRecalculo}
				<p class="mt-2 text-xs text-muted-foreground">{notaRecalculo}</p>
			{/if}
		</div>
	{/if}

	<!-- Agotamiento de calibres — ninguno cumple hasta 1000 MCM -->
	{#if agotadoCalibre}
		<div class="mt-4 rounded border border-destructive/40 bg-destructive/10 p-4">
			<h3 class="mb-1 font-semibold text-destructive">⚠ Calibre máximo superado</h3>
			<p class="text-sm text-foreground">
				Se probaron todos los calibres disponibles hasta 1000 MCM y ninguno cumple con el límite de
				caída de tensión de {limite.toFixed(1)}% para esta instalación. Considere reducir la
				longitud del circuito, aumentar el voltaje del sistema o usar conductores en paralelo (hilos
				por fase).
			</p>
		</div>
	{/if}
</section>
