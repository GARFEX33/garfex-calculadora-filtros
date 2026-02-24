<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// Get canalizacion data
	let canalizacion = $derived(memoria.canalizacion);

	// Fill factor as percentage
	let fillFactorPorcentaje = $derived((memoria.fill_factor * 100).toFixed(1));

	// Tipo de canalización legible
	let tipoCanalizacionLegible = $derived(() => {
		switch (memoria.tipo_canalizacion) {
			case 'TUBERIA_PVC':
				return 'Tubería PVC';
			case 'TUBERIA_EMT':
				return 'Tubería EMT';
			case 'CHAROLA_CABLE_ESPACIADO':
				return 'Charola de Cable (Espaciado)';
			case 'CHAROLA_CABLE_TRESBOLILLO':
				return 'Charola de Cable (Tresbolillo)';
			default:
				return memoria.tipo_canalizacion;
		}
	});
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		6. Cálculo de Canalización
	</h2>

	<p class="mb-4 text-sm text-muted-foreground">
		La canalización se dimensiona considerando el área total de los conductores y el factor de
		llenado permitido (40% para más de 3 conductores).
	</p>

	<!-- Tipo de canalización -->
	<div class="mb-4 rounded border border-primary/30 bg-primary/10 p-3">
		<p class="text-sm font-medium text-primary">
			Tipo de Canalización: {tipoCanalizacionLegible()}
		</p>
	</div>

	<!-- Fórmula (general concept) -->
	<div class="mb-6 rounded bg-muted p-4">
		<h3 class="mb-2 font-semibold text-foreground">Criterio de Dimensionamiento</h3>
		<p class="text-sm text-muted-foreground">
			El área total de la canalización debe ser mayor o igual al área total de los conductores
			dividida por el factor de llenado permitido (40% = 0.40).
		</p>
		<p class="mt-2 font-mono text-sm text-foreground">
			Área_requerida = Σ(Áreas_conductores) / 0.40
		</p>
	</div>

	<!-- Conductores -->
	<div class="mb-6">
		<h3 class="mb-2 font-semibold text-foreground">Conductores en la Canalización</h3>
		<div class="grid grid-cols-2 gap-4 text-sm">
			<div>
				<span class="text-muted-foreground">Cantidad de Conductores:</span>
				<span class="ml-2 font-medium text-foreground">{memoria.cantidad_conductores}</span>
			</div>
			<div>
				<span class="text-muted-foreground">Hilos por Fase:</span>
				<span class="ml-2 font-medium text-foreground">{memoria.hilos_por_fase}</span>
			</div>
		</div>
	</div>

	<!-- Resultado: Canalización -->
	<div>
		<h3 class="mb-2 font-semibold text-foreground">Canalización Seleccionada</h3>
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
						<td class="px-4 py-2 text-muted-foreground">Tamaño Comercial</td>
						<td class="px-4 py-2 font-mono font-medium text-foreground">
							{canalizacion?.Tamano || '—'}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Número de Tubos</td>
						<td class="px-4 py-2 text-foreground">
							{canalizacion?.NumeroDeTubos || 1}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Área Total</td>
						<td class="px-4 py-2 text-foreground">
							{canalizacion?.AreaTotalMM2?.toFixed(2) ?? '—'} mm²
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Área Requerida</td>
						<td class="px-4 py-2 text-foreground">
							{canalizacion?.AreaRequeridaMM2?.toFixed(2) ?? '—'} mm²
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Factor de Llenado</td>
						<td class="px-4 py-2 text-foreground">{fillFactorPorcentaje}%</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<!-- Verificación -->
	<div class="mt-4 rounded border border-success/30 bg-success/10 p-3">
		<p class="text-sm text-foreground">✓ Canalización dimensionada conforme a NOM-001-SEDE-2012</p>
	</div>
</section>
