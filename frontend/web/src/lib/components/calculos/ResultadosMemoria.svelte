<script lang="ts">
	import { cn } from '$lib/utils';
	import type { MemoriaOutput } from '$lib/types/calculos.types';

	interface Props {
		resultado: MemoriaOutput;
	}

	let { resultado }: Props = $props();

	// Derived values for conditional display
	let tieneObservaciones = $derived(resultado.observaciones && resultado.observaciones.length > 0);
	let cumpleCaidaTension = $derived(resultado.caida_tension?.cumple ?? false);
	let porcentajeCaida = $derived(resultado.caida_tension?.porcentaje ?? 0);

	// Nueva estructura agrupada - accesos directos para facilitar lectura
	let corrientes = $derived(resultado.corrientes);
	let instalacion = $derived(resultado.instalacion);
	let canalizacion = $derived(resultado.canalizacion);
</script>

<div class="flex flex-col gap-6">
	<!-- 1. Resumen -->
	<section class="rounded-lg border border-border bg-card p-4">
		<h3 class="mb-4 text-lg font-semibold text-card-foreground">Resumen</h3>
		<dl class="grid grid-cols-1 gap-x-6 gap-y-3 text-sm md:grid-cols-2">
			<div class="flex justify-between md:block">
				<dt class="text-muted-foreground">Cumple Norma</dt>
				<dd class="mt-1 md:mt-0">
					<span
						class={cn(
							'inline-flex items-center rounded-full px-2 py-1 text-xs font-medium',
							resultado.cumple_normativa
								? 'bg-success text-success-foreground'
								: 'bg-destructive text-destructive-foreground'
						)}
					>
						{resultado.cumple_normativa ? 'Cumple' : 'No Cumple'}
					</span>
				</dd>
			</div>
			<div class="flex justify-between md:block">
				<dt class="text-muted-foreground">Corriente Nominal</dt>
				<dd class="font-mono text-foreground">{corrientes.corriente_nominal.toFixed(2)} A</dd>
			</div>
			<div class="flex justify-between md:block">
				<dt class="text-muted-foreground">Corriente Ajustada</dt>
				<dd class="font-mono text-foreground">{corrientes.corriente_ajustada.toFixed(2)} A</dd>
			</div>
			<div class="flex justify-between md:block">
				<dt class="text-muted-foreground">Factor Total de Ajuste</dt>
				<dd class="font-mono text-foreground">{corrientes.factor_total_ajuste.toFixed(3)}</dd>
			</div>
		</dl>
	</section>

	<!-- 2. Sistema -->
	<section class="rounded-lg border border-border bg-card p-4">
		<h3 class="mb-4 text-lg font-semibold text-card-foreground">Sistema</h3>
		<dl class="grid grid-cols-2 gap-x-6 gap-y-3 text-sm md:grid-cols-3">
			<div>
				<dt class="text-muted-foreground">Tipo de Equipo</dt>
				<dd class="text-foreground">{resultado.tipo_equipo}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Tensión</dt>
				<dd class="text-foreground">{instalacion.tension} V</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Sistema Eléctrico</dt>
				<dd class="text-foreground">{instalacion.sistema_electrico}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Estado</dt>
				<dd class="text-foreground">{resultado.estado}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Temperatura Ambiente</dt>
				<dd class="text-foreground">{corrientes.temperatura_ambiente} °C</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Factor de Potencia</dt>
				<dd class="font-mono text-foreground">{resultado.factor_potencia.toFixed(2)}</dd>
			</div>
		</dl>
	</section>

	<!-- 3. Conductor de Alimentación -->
	<section class="rounded-lg border border-border bg-card p-4">
		<h3 class="mb-4 text-lg font-semibold text-card-foreground">Conductor de Alimentación</h3>
		<dl class="grid grid-cols-2 gap-x-6 gap-y-3 text-sm md:grid-cols-3">
			<div>
				<dt class="text-muted-foreground">Calibre</dt>
				<dd class="font-mono font-medium text-foreground">
					{resultado.cable_fase.calibre}
				</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Material</dt>
				<dd class="text-foreground">{resultado.cable_fase.material}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Tipo de Aislamiento</dt>
				<dd class="text-foreground">
					{resultado.cable_fase.tipo_aislamiento || '—'}
				</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Capacidad</dt>
				<dd class="text-foreground">{resultado.cable_fase.capacidad} A</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Sección</dt>
				<dd class="text-foreground">{resultado.cable_fase.seccion_mm2} mm²</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Tabla de Ampacidad</dt>
				<dd class="text-foreground">{corrientes.tabla_ampacidad_usada}</dd>
			</div>
		</dl>
	</section>

	<!-- 4. Conductor de Tierra -->
	<section class="rounded-lg border border-border bg-card p-4">
		<h3 class="mb-4 text-lg font-semibold text-card-foreground">Conductor de Tierra</h3>
		<dl class="grid grid-cols-2 gap-x-6 gap-y-3 text-sm md:grid-cols-3">
			<div>
				<dt class="text-muted-foreground">Calibre</dt>
				<dd class="font-mono font-medium text-foreground">{resultado.cable_tierra.calibre}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Material</dt>
				<dd class="text-foreground">{resultado.cable_tierra.material}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Tipo de Aislamiento</dt>
				<dd class="text-foreground">{resultado.cable_tierra.tipo_aislamiento || '—'}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Capacidad</dt>
				<dd class="text-foreground">{resultado.cable_tierra.capacidad} A</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">ITM</dt>
				<dd class="text-foreground">{resultado.proteccion.itm} A</dd>
			</div>
		</dl>
	</section>

	<!-- 5. Canalización -->
	<section class="rounded-lg border border-border bg-card p-4">
		<h3 class="mb-4 text-lg font-semibold text-card-foreground">Canalización</h3>
		<dl class="grid grid-cols-2 gap-x-6 gap-y-3 text-sm md:grid-cols-3">
			<div>
				<dt class="text-muted-foreground">Tamaño</dt>
				<dd class="font-mono font-medium text-foreground">{canalizacion.resultado.tamano}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Tipo de Canalización</dt>
				<dd class="text-foreground">{instalacion.tipo_canalizacion}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Número de Tubos</dt>
				<dd class="text-foreground">{canalizacion.resultado.numero_de_tubos}</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Fill Factor</dt>
				<dd class="text-foreground">{(canalizacion.fill_factor * 100).toFixed(1)}%</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Área Total</dt>
				<dd class="text-foreground">{canalizacion.resultado.area_total_mm2.toFixed(2)} mm²</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Área Requerida</dt>
				<dd class="text-foreground">{canalizacion.resultado.area_requerida_mm2.toFixed(2)} mm²</dd>
			</div>
		</dl>
	</section>

	<!-- 6. Caída de Tensión -->
	<section class="rounded-lg border border-border bg-card p-4">
		<h3 class="mb-4 text-lg font-semibold text-card-foreground">Caída de Tensión</h3>
		<dl class="grid grid-cols-2 gap-x-6 gap-y-3 text-sm md:grid-cols-3">
			<div>
				<dt class="text-muted-foreground">Porcentaje</dt>
				<dd
					class={cn(
						'font-mono font-medium',
						cumpleCaidaTension ? 'text-success' : 'text-destructive'
					)}
				>
					{porcentajeCaida.toFixed(2)}%
				</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Caída</dt>
				<dd class="font-mono text-foreground">
					{resultado.caida_tension.caida_volts.toFixed(2)} V
				</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Límite</dt>
				<dd class="text-foreground">{resultado.caida_tension.limite_porcentaje.toFixed(1)}%</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Impedancia</dt>
				<dd class="font-mono text-foreground">{resultado.caida_tension.impedancia.toFixed(4)} Ω</dd>
			</div>
			<div>
				<dt class="text-muted-foreground">Cumple</dt>
				<dd>
					<span
						class={cn(
							'inline-flex items-center rounded-full px-2 py-1 text-xs font-medium',
							cumpleCaidaTension
								? 'bg-success text-success-foreground'
								: 'bg-destructive text-destructive-foreground'
						)}
					>
						{cumpleCaidaTension ? 'Sí' : 'No'}
					</span>
				</dd>
			</div>
		</dl>
	</section>

	<!-- 7. Observaciones -->
	{#if tieneObservaciones}
		<section class="rounded-lg border border-border bg-card p-4">
			<h3 class="mb-4 text-lg font-semibold text-card-foreground">Observaciones</h3>
			<ul class="list-inside list-disc space-y-1 text-sm text-foreground">
				{#each resultado.observaciones as obs}
					<li>{obs}</li>
				{/each}
			</ul>
		</section>
	{/if}
</div>
