<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// Determine if three-phase or single-phase
	let esTrifasico = $derived(
		memoria.sistema_electrico === 'ESTRELLA' || memoria.sistema_electrico === 'DELTA'
	);

	// Get equipment type
	let tipoEquipo = $derived(memoria.tipo_equipo?.toUpperCase() ?? '');

	// Helper to get calculation info based on equipment type and system
	function getInfoCalculo(): {
		tipo: string;
		formula: string;
		desarrollo: string[];
		valores: Record<string, string>;
	} {
		const tensionV = memoria.tension;
		const fp = memoria.factor_potencia ?? 1;

		// FILTRO_ACTIVO - amperaje directo
		if (tipoEquipo === 'FILTRO_ACTIVO') {
			return {
				tipo: 'Amperaje directo',
				formula: 'I = Iₙominal',
				desarrollo: [`I = ${memoria.corriente_nominal.toFixed(2)} A (dato del equipo)`],
				valores: {
					Amperaje: `${memoria.equipo?.amperaje ?? memoria.corriente_nominal.toFixed(2)} A`,
					Tipo: 'Filtro Activo (FP = 1.0)'
				}
			};
		}

		// TRANSFORMADOR - I = KVA / (kV × √3)
		if (tipoEquipo === 'TRANSFORMADOR') {
			const kva = (memoria.corriente_nominal * tensionV * Math.sqrt(3)) / 1000;
			const kv = tensionV / 1000;
			return {
				tipo: 'Desde KVA (Transformador)',
				formula: 'I = KVA / (kV × √3)',
				desarrollo: [
					`I = ${kva.toFixed(2)} kVA / (${kv.toFixed(3)} kV × 1.732)`,
					`I = ${kva.toFixed(2)} / ${(kv * Math.sqrt(3)).toFixed(3)}`,
					`I = ${memoria.corriente_nominal.toFixed(2)} A`
				],
				valores: {
					KVA: `${kva.toFixed(2)} kVA`,
					Voltaje: `${tensionV} V (${kv.toFixed(3)} kV)`,
					Fórmula: 'I = KVA / (kV × √3)'
				}
			};
		}

		// FILTRO_RECHAZO - I = KVAR / (kV × √3)
		if (tipoEquipo === 'FILTRO_RECHAZO') {
			const kvar = (memoria.corriente_nominal * tensionV * Math.sqrt(3)) / 1000;
			const kv = tensionV / 1000;
			return {
				tipo: 'Desde KVAR (Filtro de Rechazo)',
				formula: 'I = KVAR / (kV × √3)',
				desarrollo: [
					`I = ${kvar.toFixed(2)} kVAR / (${kv.toFixed(3)} kV × 1.732)`,
					`I = ${kvar.toFixed(2)} / ${(kv * Math.sqrt(3)).toFixed(3)}`,
					`I = ${memoria.corriente_nominal.toFixed(2)} A`
				],
				valores: {
					KVAR: `${kvar.toFixed(2)} kVAR`,
					Voltaje: `${tensionV} V (${kv.toFixed(3)} kV)`,
					Fórmula: 'I = KVAR / (kV × √3)'
				}
			};
		}

		// CARGA o MANUAL_POTENCIA - depends on system type
		if (esTrifasico) {
			// Trifásico: I = W / (V × √3 × FP)
			const potenciaKW = (memoria.corriente_nominal * tensionV * Math.sqrt(3) * fp) / 1000;
			const potenciaW = potenciaKW * 1000;
			return {
				tipo: 'Desde Potencia (Sistema Trifásico)',
				formula: 'I = P / (V × √3 × cosθ)',
				desarrollo: [
					`P = ${potenciaKW.toFixed(2)} kW = ${potenciaW.toFixed(0)} W`,
					`I = ${potenciaW.toFixed(0)} / (${tensionV} × 1.732 × ${fp.toFixed(2)})`,
					`I = ${potenciaW.toFixed(0)} / ${(tensionV * Math.sqrt(3) * fp).toFixed(2)}`,
					`I = ${memoria.corriente_nominal.toFixed(2)} A`
				],
				valores: {
					Potencia: `${potenciaKW.toFixed(2)} kW`,
					Voltaje: `${tensionV} V`,
					'Factor de Potencia': fp.toFixed(2),
					Sistema: memoria.sistema_electrico,
					Fórmula: 'I = P / (V × √3 × cosθ)'
				}
			};
		} else {
			// Monofásico/Bifásico: I = W / (V × FP)
			const potenciaKW = (memoria.corriente_nominal * tensionV * fp) / 1000;
			const potenciaW = potenciaKW * 1000;
			return {
				tipo: 'Desde Potencia (Sistema Monofásico)',
				formula: 'I = P / (V × cosθ)',
				desarrollo: [
					`P = ${potenciaKW.toFixed(2)} kW = ${potenciaW.toFixed(0)} W`,
					`I = ${potenciaW.toFixed(0)} / (${tensionV} × ${fp.toFixed(2)})`,
					`I = ${potenciaW.toFixed(0)} / ${(tensionV * fp).toFixed(2)}`,
					`I = ${memoria.corriente_nominal.toFixed(2)} A`
				],
				valores: {
					Potencia: `${potenciaKW.toFixed(2)} kW`,
					Voltaje: `${tensionV} V`,
					'Factor de Potencia': fp.toFixed(2),
					Sistema: memoria.sistema_electrico,
					Fórmula: 'I = P / (V × cosθ)'
				}
			};
		}
	}

	let info = $derived(getInfoCalculo());
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		2. Cálculo de Corriente Nominal
	</h2>

	<p class="mb-4 text-sm text-muted-foreground">
		La corriente nominal se calcula en función del tipo de equipo y sistema eléctrico.
	</p>

	<!-- Tipo de cálculo -->
	<div class="mb-4 rounded bg-muted p-3">
		<span class="text-sm font-medium text-muted-foreground">Tipo de cálculo:</span>
		<span class="ml-2 font-semibold text-foreground">{info.tipo}</span>
	</div>

	<!-- Fórmula -->
	<div class="mb-6 rounded bg-muted p-4">
		<h3 class="mb-2 font-semibold text-foreground">Fórmula</h3>
		<p class="font-mono text-lg text-foreground">{info.formula}</p>
	</div>

	<!-- Desarrollo -->
	<div class="mb-6">
		<h3 class="mb-2 font-semibold text-foreground">Desarrollo</h3>
		<div class="space-y-1 rounded border border-border bg-muted/50 p-4">
			{#each info.desarrollo as paso, i}
				<p
					class="font-mono text-sm {i === info.desarrollo.length - 1
						? 'font-semibold text-foreground'
						: 'text-muted-foreground'}"
				>
					{#if i < info.desarrollo.length - 1}
						<span class="text-muted-foreground">→</span>
					{/if}
					{paso}
				</p>
			{/each}
		</div>
	</div>

	<!-- Resultado -->
	<div class="rounded border border-success/30 bg-success/10 p-4">
		<h3 class="mb-2 font-semibold text-success">Resultado</h3>
		<p class="text-2xl font-bold text-foreground">
			I<sub>n</sub> = {memoria.corriente_nominal.toFixed(2)} A
		</p>
	</div>

	<!-- Datos de referencia -->
	<div class="mt-4 grid grid-cols-2 gap-4 text-sm">
		<div>
			<span class="text-muted-foreground">Sistema:</span>
			<span class="ml-2 font-medium text-foreground">{memoria.sistema_electrico}</span>
		</div>
		<div>
			<span class="text-muted-foreground">Voltaje:</span>
			<span class="ml-2 font-medium text-foreground">{memoria.tension} V</span>
		</div>
		{#if info.valores['Factor de Potencia']}
			<div>
				<span class="text-muted-foreground">Factor de Potencia:</span>
				<span class="ml-2 font-medium text-foreground">{info.valores['Factor de Potencia']}</span>
			</div>
		{/if}
		{#if info.valores['Tipo']}
			<div>
				<span class="text-muted-foreground">Equipo:</span>
				<span class="ml-2 font-medium text-foreground">{info.valores['Tipo']}</span>
			</div>
		{/if}
	</div>
</section>
