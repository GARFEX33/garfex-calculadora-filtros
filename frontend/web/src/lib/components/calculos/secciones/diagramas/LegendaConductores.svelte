<script lang="ts">
	import { COLORES } from './geometry.js';

	interface Props {
		mostrarNeutro?: boolean;
		mostrarControl?: boolean;
	}

	let { mostrarNeutro = false, mostrarControl = false }: Props = $props();
</script>

<div class="flex flex-wrap gap-4 text-xs text-foreground">
	<!-- Fase - siempre mostrar -->
	<div class="flex items-center gap-1.5">
		<span
			class="h-3 w-3 rounded-full border border-gray-400"
			style="background-color: {COLORES.fase};"
		></span>
		<span class="font-medium">Fase</span>
	</div>

	<!-- Tierra - siempre mostrar -->
	<div class="flex items-center gap-1.5">
		<span
			class="h-3 w-3 rounded-full border border-gray-400"
			style="background-color: {COLORES.tierra};"
		></span>
		<span class="font-medium">Tierra</span>
	</div>

	<!-- Neutro - opcional -->
	{#if mostrarNeutro}
		<div class="flex items-center gap-1.5">
			<span
				class="h-3 w-3 rounded-full border border-gray-400"
				style="background-color: {COLORES.neutro};"
			></span>
			<span class="font-medium">Neutro</span>
		</div>
	{/if}

	<!-- Control - opcional -->
	{#if mostrarControl}
		<div class="flex items-center gap-1.5">
			<span
				class="h-3 w-3 rounded-full border border-gray-400"
				style="background-color: {COLORES.control};"
			></span>
			<span class="font-medium">Control</span>
		</div>
	{/if}
</div>
