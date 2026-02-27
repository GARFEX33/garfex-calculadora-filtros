<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// Calcular capacidad en kVA si es posible
	let capacidadKVA = $derived(() => {
		if (memoria.factor_potencia && memoria.factor_potencia > 0) {
			const potencia = (memoria.corrientes.corriente_nominal * memoria.instalacion.tension * Math.sqrt(3)) / 1000;
			return (potencia / memoria.factor_potencia).toFixed(2);
		}
		return null;
	});
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		1. Encabezado
	</h2>

	<!-- Título del documento -->
	<div class="mb-6 text-center">
		<h3 class="text-lg font-bold text-foreground">Memoria de Cálculo de Alimentador</h3>
		<p class="text-sm text-muted-foreground">Conforme a NOM-001-SEDE-2012</p>
	</div>

	<!-- Datos del equipo y proyecto -->
	<dl class="grid grid-cols-1 gap-x-6 gap-y-4 text-sm md:grid-cols-2">
		<div class="border-l-2 border-primary pl-4">
			<dt class="text-xs tracking-wide text-muted-foreground uppercase">Equipo / Carga</dt>
			<dd class="font-medium text-foreground">
				{memoria.equipo?.clave || memoria.tipo_equipo || '—'}
			</dd>
		</div>

		<div class="border-l-2 border-primary pl-4">
			<dt class="text-xs tracking-wide text-muted-foreground uppercase">Tipo de Equipo</dt>
			<dd class="font-medium text-foreground">{memoria.tipo_equipo || '—'}</dd>
		</div>

		<div class="border-l-2 border-primary pl-4">
			<dt class="text-xs tracking-wide text-muted-foreground uppercase">Capacidad</dt>
			<dd class="font-medium text-foreground">
				{#if capacidadKVA()}
					{capacidadKVA()} kVA
				{:else}
					{memoria.corrientes.corriente_nominal.toFixed(2)} A
				{/if}
			</dd>
		</div>

		<div class="border-l-2 border-primary pl-4">
			<dt class="text-xs tracking-wide text-muted-foreground uppercase">Voltaje de Operación</dt>
			<dd class="font-medium text-foreground">{memoria.instalacion.tension} V</dd>
		</div>

		<div class="border-l-2 border-primary pl-4">
			<dt class="text-xs tracking-wide text-muted-foreground uppercase">Sistema Eléctrico</dt>
			<dd class="font-medium text-foreground">{memoria.instalacion.sistema_electrico}</dd>
		</div>

		<div class="border-l-2 border-primary pl-4">
			<dt class="text-xs tracking-wide text-muted-foreground uppercase">
				Longitud del Alimentador
			</dt>
			<dd class="font-medium text-foreground">{memoria.instalacion.longitud_circuito.toFixed(2)} m</dd>
		</div>

		<div class="border-l-2 border-primary pl-4">
			<dt class="text-xs tracking-wide text-muted-foreground uppercase">Estado</dt>
			<dd class="font-medium text-foreground">{memoria.estado}</dd>
		</div>

		<div class="border-l-2 border-primary pl-4">
			<dt class="text-xs tracking-wide text-muted-foreground uppercase">Temperatura Ambiente</dt>
			<dd class="font-medium text-foreground">{memoria.corrientes.temperatura_ambiente} °C</dd>
		</div>
	</dl>
</section>
