<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';
	import { cn } from '$lib/utils';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// Helper para obtener label del tipo de equipo
	const labelsTipoEquipo: Record<string, string> = {
		FILTRO_ACTIVO: 'Filtro Activo',
		FILTRO_RECHAZO: 'Filtro de Rechazo',
		TRANSFORMADOR: 'Transformador',
		CARGA: 'Carga General'
	};

	// Factor de uso según tipo de equipo (Artículo 460-8)
	let factorUso = $derived(
		memoria.tipo_equipo === 'FILTRO_ACTIVO' || memoria.tipo_equipo === 'FILTRO_RECHAZO'
			? 1.35
			: 1.25
	);

	let justificacionFactorUso = $derived(
		memoria.tipo_equipo === 'FILTRO_ACTIVO' || memoria.tipo_equipo === 'FILTRO_RECHAZO'
			? 'Los conductores para capacitores deben tener al menos el 135% de la corriente nominal (Artículo 460-8)'
			: 'Factor de diseño estándar para equipos de carga general (125%)'
	);

	// Número de hilos del conductor de alimentación (del backend)
	let numHilosAlimentacion = $derived(
		(memoria.conductor_alimentacion.NumHilos ?? memoria.hilos_por_fase) || 1
	);

	// Verificación de capacidad
	let capacidadPorHilo = $derived(memoria.conductor_alimentacion.Capacidad);
	let capacidadTotal = $derived(
		numHilosAlimentacion > 1
			? memoria.conductor_alimentacion.Capacidad * numHilosAlimentacion
			: memoria.conductor_alimentacion.Capacidad
	);
	let cumpleCapacidad = $derived(capacidadTotal >= memoria.corriente_ajustada);

	// Verificar si es charola (no aplica agrupamiento)
	let esCharola = $derived(memoria.tipo_canalizacion.includes('CHAROLA'));
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		3. Dimensionamiento del Alimentador
	</h2>

	<p class="mb-4 text-sm text-muted-foreground">
		El conductor de alimentación se dimensiona aplicando los factores de corrección establecidos en
		la NOM-001-SEDE-2012, considerando el tipo de equipo, la temperatura ambiente y el agrupamiento
		de conductores.
	</p>

	<!-- Norma de referencia -->
	<div class="mb-4 rounded border border-primary/30 bg-primary/10 p-3">
		<p class="text-sm font-medium text-primary">
			Referencias: Artículo 460-8 (Conductores para Capacitores) • Artículo 310-15(b)(2)(A)
			(Temperatura) • Artículo 310-15(b)(3)(A) (Agrupamiento)
		</p>
	</div>

	<!-- Factor de Uso -->
	<div class="mb-4 rounded-lg border border-border bg-card p-4">
		<h3 class="mb-2 flex items-center gap-2 font-semibold text-foreground">
			<span class="text-primary">▪</span>
			Factor de Uso (Tipo de Equipo)
		</h3>
		<div class="space-y-1 text-sm">
			<p class="text-foreground">
				<strong>Tipo de Equipo:</strong>
				{labelsTipoEquipo[memoria.tipo_equipo] || memoria.tipo_equipo}
			</p>
			<p class="text-foreground">
				<strong>Factor de Uso:</strong>
				{(factorUso * 100).toFixed(0)}% ({factorUso})
			</p>
			<p class="text-muted-foreground">{justificacionFactorUso}</p>
		</div>
	</div>

	<!-- Factores de Temperatura y Agrupamiento en grid -->
	<div class="mb-4 grid gap-4 md:grid-cols-2">
		<!-- Factor de Temperatura -->
		<div class="rounded-lg border border-border bg-card p-4">
			<h3 class="mb-2 flex items-center gap-2 font-semibold text-foreground">
				<span class="text-primary">▪</span>
				Factor de Temperatura
			</h3>
			<div class="space-y-1 text-sm">
				<p class="text-foreground">
					<strong>Temperatura Ambiente:</strong>
					{memoria.temperatura_ambiente} °C
					<span class="text-muted-foreground">(Estado: {memoria.estado})</span>
				</p>
				<p class="text-foreground">
					<strong>Temperatura del Conductor:</strong>
					{memoria.temperatura_usada} °C
				</p>
				<p class="text-foreground">
					<strong>Factor de Temperatura:</strong>
					{memoria.factor_temperatura.toFixed(2)}
				</p>
				<p class="text-muted-foreground">Referencia: Tabla 310-15(b)(2)(A)</p>
			</div>
		</div>

		<!-- Factor de Agrupamiento -->
		<div class="rounded-lg border border-border bg-card p-4">
			<h3 class="mb-2 flex items-center gap-2 font-semibold text-foreground">
				<span class="text-primary">▪</span>
				Factor de Agrupamiento
			</h3>
			<div class="space-y-1 text-sm">
				<p class="text-foreground">
					<strong>Cantidad de Conductores:</strong>
					{memoria.cantidad_conductores}
				</p>
				<p class="text-foreground">
					<strong>Factor de Agrupamiento:</strong>
					{memoria.factor_agrupamiento.toFixed(2)}
					{#if esCharola}
						<span class="ml-1 text-success">(No aplica - Charola)</span>
					{/if}
				</p>
				<p class="text-muted-foreground">
					{#if esCharola}
						Los cables en charola van separados, no aplica factor de agrupamiento
					{:else}
						Referencia: Tabla 310-15(b)(3)(A)
					{/if}
				</p>
			</div>
		</div>
	</div>

	<!-- Fórmula -->
	<div class="mb-6 rounded bg-muted p-4">
		<h3 class="mb-2 font-semibold text-foreground">Fórmula de Dimensionamiento</h3>
		<p class="font-mono text-foreground">
			I<sub>ajustada</sub> = I<sub>nominal</sub> × F<sub>uso</sub> / (F<sub>temp</sub> × F<sub
				>agr</sub
			>)
		</p>
	</div>

	<!-- Desarrollo -->
	<div class="mb-6">
		<h3 class="mb-2 font-semibold text-foreground">Desarrollo</h3>
		<div class="space-y-2 font-mono text-sm">
			<p class="text-foreground">
				I<sub>ajustada</sub> = {memoria.corriente_nominal.toFixed(2)} A × {factorUso.toFixed(2)} / ({memoria.factor_temperatura.toFixed(
					2
				)} × {memoria.factor_agrupamiento.toFixed(2)})
			</p>
			<p class="text-foreground">
				I<sub>ajustada</sub> = {memoria.corriente_nominal.toFixed(2)} A × {factorUso.toFixed(2)} /
				{memoria.factor_total_ajuste.toFixed(3)}
			</p>
			<p class="text-lg font-bold text-primary">
				I<sub>ajustada</sub> = {memoria.corriente_ajustada.toFixed(2)} A
			</p>
			{#if numHilosAlimentacion > 1}
				<hr class="my-2 border-border" />
				<p class="text-foreground">
					<strong>{numHilosAlimentacion}</strong> hilos por fase en paralelo
				</p>
				<p class="text-foreground">
					I<sub>hilo</sub> = {memoria.corriente_ajustada.toFixed(2)} A / {numHilosAlimentacion} =
					<strong>{memoria.corriente_por_hilo.toFixed(2)} A</strong> por hilo
				</p>
			{/if}
		</div>
	</div>

	<!-- Resultado: Conductor seleccionado -->
	<div class="mb-6">
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
							{memoria.conductor_alimentacion.Calibre}
							{#if numHilosAlimentacion > 1}
								<span class="ml-1 text-xs text-muted-foreground">× {numHilosAlimentacion}</span>
							{/if}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Material</td>
						<td class="px-4 py-2 text-foreground">
							{memoria.conductor_alimentacion.Material?.toUpperCase() === 'CU'
								? 'Cobre (Cu)'
								: 'Aluminio (Al)'}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Sección</td>
						<td class="px-4 py-2 text-foreground">
							{memoria.conductor_alimentacion.SeccionMM2.toFixed(2)} mm²
							{#if numHilosAlimentacion > 1}
								<span class="ml-1 text-xs text-muted-foreground"
									>(× {numHilosAlimentacion} = {(
										memoria.conductor_alimentacion.SeccionMM2 * numHilosAlimentacion
									).toFixed(2)} mm² total)</span
								>
							{/if}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Tipo de Aislamiento</td>
						<td class="px-4 py-2 text-foreground">
							{memoria.conductor_alimentacion.TipoAislamiento || 'THWN'}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Número de Hilos por Fase</td>
						<td class="px-4 py-2 text-foreground">
							{numHilosAlimentacion}
						</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Temperatura de Referencia</td>
						<td class="px-4 py-2 text-foreground">{memoria.temperatura_usada} °C</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Ampacidad por Hilo</td>
						<td class="px-4 py-2 font-medium text-foreground">
							{memoria.conductor_alimentacion.Capacidad} A
						</td>
					</tr>
					{#if numHilosAlimentacion > 1}
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Capacidad Total</td>
							<td class="px-4 py-2 font-medium text-foreground">
								{memoria.conductor_alimentacion.Capacidad} A × {numHilosAlimentacion} = {(
									memoria.conductor_alimentacion.Capacidad * numHilosAlimentacion
								).toFixed(0)} A
							</td>
						</tr>
					{/if}
				</tbody>
			</table>
		</div>
	</div>

	<!-- Verificación -->
	<div
		class={cn(
			'rounded border p-4',
			cumpleCapacidad
				? 'border-success/30 bg-success/10'
				: 'border-destructive/30 bg-destructive/10'
		)}
	>
		<p class="text-sm">
			{#if cumpleCapacidad}
				<span class="font-medium text-success">✓</span>
				{#if numHilosAlimentacion > 1}
					<span class="text-foreground">
						La ampacidad total ({capacidadTotal} A = {memoria.conductor_alimentacion.Capacidad} A × {numHilosAlimentacion}
						hilos) es mayor o igual a la corriente ajustada ({memoria.corriente_ajustada.toFixed(2)} A).
						Cada hilo transporta {memoria.corriente_por_hilo.toFixed(2)} A (≤ {capacidadPorHilo} A).
					</span>
				{:else}
					<span class="text-foreground">
						La ampacidad ({capacidadTotal} A) es mayor o igual a la corriente ajustada ({memoria.corriente_ajustada.toFixed(
							2
						)} A). El conductor cumple.
					</span>
				{/if}
			{:else}
				<span class="font-medium text-destructive">✗</span>
				<span class="text-foreground">
					La ampacidad ({capacidadTotal} A) es menor a la corriente ajustada ({memoria.corriente_ajustada.toFixed(
						2
					)} A). El conductor NO cumple.
				</span>
			{/if}
		</p>
	</div>

	<!-- Tabla utilizada -->
	<div class="mt-4 text-sm">
		<span class="text-muted-foreground">Tabla de Ampacidad Utilizada:</span>
		<span class="ml-2 font-medium text-foreground">{memoria.tabla_ampacidad_usada}</span>
	</div>
</section>
