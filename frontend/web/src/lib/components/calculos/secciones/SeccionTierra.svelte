<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// Número de hilos del conductor de tierra (del backend)
	let numHilosTierra = $derived(memoria.conductor_tierra.num_hilos ?? 1);

	// Canalización type and tube count
	let tipo = $derived(memoria.tipo_canalizacion);
	let canalizacion = $derived(memoria.canalizacion);
	let numTubos = $derived(canalizacion.numero_de_tubos ?? 1);

	// Clasificación del tipo de canalización
	let esTuberia = $derived(
		tipo === 'TUBERIA_PVC' ||
			tipo === 'TUBERIA_ALUMINIO' ||
			tipo === 'TUBERIA_ACERO_PG' ||
			tipo === 'TUBERIA_ACERO_PD'
	);

	// numHilosTierraMostrar: el backend ya envía el valor correcto (1 tierra por tubo = numTuberias total)
	let numHilosTierraMostrar = $derived(numHilosTierra);

	// Bug 2 fix: conductores por tubería para tubería multi-tubo
	let conductoresPorTubo = $derived(
		esTuberia && numTubos > 1
			? Math.ceil(memoria.cantidad_conductores / numTubos)
			: memoria.cantidad_conductores
	);

	// Bug 3 fix: multiplicador ×N solo aplica para charola (conductores paralelos por fase)
	let mostrarMultiplicadorTierra = $derived(!esTuberia && numHilosTierra > 1);
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		4. Conductor de Puesta a Tierra
	</h2>

	<p class="mb-4 text-sm text-muted-foreground">
		El conductor de puesta a tierra se dimensiona conforme a la Tabla 250-122 de la
		NOM-001-SEDE-2012, en función de la corriente del dispositivo de protección (ITM).
	</p>

	<!-- Norma de referencia -->
	<div class="mb-4 rounded border border-primary/30 bg-primary/10 p-3">
		<p class="text-sm font-medium text-primary">
			Referencia: Tabla 250-122 - Conductores de Puesta a Tierra de Equipos (AWG/kcmil)
		</p>
	</div>

	<!-- Fórmula y desarrollo -->
	<div class="mb-6 rounded bg-muted p-4">
		<h3 class="mb-2 font-semibold text-foreground">Criterio de Selección</h3>
		<p class="text-sm text-muted-foreground">
			El calibre del conductor de tierra se selecciona según el valor del Interruptor Termomagnético
			(ITM) que protege el circuito.
		</p>
	</div>

	<!-- ITM utilizado y contexto de hilos -->
	<div class="mb-4">
		<p class="text-sm text-muted-foreground">
			<span class="font-medium text-foreground">ITM:</span>
			{memoria.itm} A
		</p>
		{#if memoria.hilos_por_fase > 1}
			<p class="mt-1 text-sm text-muted-foreground">
				<span class="font-medium text-foreground"
					>Circuito con {memoria.hilos_por_fase} hilos por fase:</span
				>
				{#if esTuberia && numTubos > 1}
					{conductoresPorTubo} conductores por tubería/canalización
				{:else}
					{memoria.cantidad_conductores} conductores totales en la canalización
				{/if}
			</p>
		{/if}
	</div>

	<!-- Resultado: Conductor de tierra -->
	<div>
		<h3 class="mb-2 font-semibold text-foreground">Conductor Seleccionado</h3>
		<div class="overflow-hidden rounded border border-border">
			<table class="w-full text-sm">
				<thead class="bg-muted">
					<tr>
						<th class="px-4 py-2 text-left font-medium text-foreground">Parámetro</th>
						<th class="px-4 py-2 text-left font-medium text-foreground">Valor</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-border">
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Calibre</td>
						<td class="px-4 py-2 font-mono font-medium text-foreground">
							{memoria.conductor_tierra.calibre}
							{#if mostrarMultiplicadorTierra}
								<span class="ml-1 text-xs text-muted-foreground">× {numHilosTierra}</span>
							{/if}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Material</td>
						<td class="px-4 py-2 text-foreground">
							{memoria.conductor_tierra.material?.toUpperCase() === 'CU'
								? 'Cobre (Cu)'
								: 'Aluminio (Al)'}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Sección</td>
						<td class="px-4 py-2 text-foreground">
							{memoria.conductor_tierra.seccion_mm2.toFixed(2)} mm²
							{#if mostrarMultiplicadorTierra}
								<span class="ml-1 text-xs text-muted-foreground"
									>(× {numHilosTierra} = {(
										memoria.conductor_tierra.seccion_mm2 * numHilosTierra
									).toFixed(2)} mm² total)</span
								>
							{/if}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Tipo de Aislamiento</td>
						<td class="px-4 py-2 text-foreground">
							{memoria.conductor_tierra.tipo_aislamiento || 'Desnudo'}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Número de Hilos</td>
						<td class="px-4 py-2 text-foreground">
							{numHilosTierraMostrar}
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<!-- Justificación -->
	<div class="mt-4 rounded border border-success/30 bg-success/10 p-3">
		<p class="text-sm text-foreground">
			✓ Conductor de puesta a tierra seleccionado conforme a Tabla 250-122 para ITM de{' '}
			{memoria.itm} A
		</p>
	</div>

	<!-- Tabla utilizada -->
	<div class="mt-4 text-sm">
		<span class="text-muted-foreground">Tabla de Referencia Utilizada:</span>
		<span class="ml-2 font-medium text-foreground"
			>NOM-250-122 (Conductores de Puesta a Tierra)</span
		>
	</div>
</section>
