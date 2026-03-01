<script lang="ts">
	import type {
		DetalleCharola,
		DetalleTuberia,
		ResultadoCanalizacion
	} from '$lib/types/calculos.types';
	import type { SistemaElectrico } from '$lib/features/calculos/domain/types/calculo.enums.js';
	import DiagramaCharolaEspaciada from './diagramas/DiagramaCharolaEspaciada.svelte';
	import DiagramaCharolaTriangular from './diagramas/DiagramaCharolaTriangular.svelte';
	import DiagramaTuberia from './diagramas/DiagramaTuberia.svelte';

	interface Props {
		tipoCanalizacion: string;
		detalleCharola?: DetalleCharola | undefined;
		detalleTuberia?: DetalleTuberia | undefined;
		resultado: ResultadoCanalizacion;
		sistemaElectrico: SistemaElectrico;
		hilosPorFase: number;
		calibreFase: string;
		calibreTierra: string;
	}

	let {
		tipoCanalizacion,
		detalleCharola = undefined,
		detalleTuberia = undefined,
		resultado,
		sistemaElectrico,
		hilosPorFase,
		calibreFase,
		calibreTierra
	}: Props = $props();

	let esCharolaEspaciado = $derived(tipoCanalizacion === 'CHAROLA_CABLE_ESPACIADO');
	let esCharolaTriangular = $derived(tipoCanalizacion === 'CHAROLA_CABLE_TRIANGULAR');
	let esTuberia = $derived(
		tipoCanalizacion === 'TUBERIA_PVC' ||
			tipoCanalizacion === 'TUBERIA_ALUMINIO' ||
			tipoCanalizacion === 'TUBERIA_ACERO_PG' ||
			tipoCanalizacion === 'TUBERIA_ACERO_PD'
	);
</script>

<div class="rounded-lg border border-border bg-card">
	{#if esCharolaEspaciado && detalleCharola}
		<DiagramaCharolaEspaciada
			detalle={detalleCharola}
			{resultado}
			{sistemaElectrico}
			{hilosPorFase}
			{calibreFase}
			{calibreTierra}
		/>
	{:else if esCharolaTriangular && detalleCharola}
		<DiagramaCharolaTriangular
			detalle={detalleCharola}
			{resultado}
			{sistemaElectrico}
			{hilosPorFase}
			{calibreFase}
			{calibreTierra}
		/>
	{:else if esTuberia && detalleTuberia}
		<DiagramaTuberia
			detalle={detalleTuberia}
			{resultado}
			{sistemaElectrico}
			{calibreFase}
			{calibreTierra}
		/>
	{:else}
		<!-- Fallback para tipos no reconocidos o sin detalle -->
		<div class="p-4 text-center text-muted-foreground">
			<p class="text-sm">Diagrama no disponible para el tipo de canalización seleccionado.</p>
			<p class="mt-1 text-xs">Tipo: {tipoCanalizacion}</p>
		</div>
	{/if}
</div>
