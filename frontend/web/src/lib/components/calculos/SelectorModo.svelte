<script lang="ts">
	import { cn } from '$lib/utils';
	import type { ModoCalculo } from '$lib/types/calculos.types';

	interface Props {
		modo: ModoCalculo;
		onModoChange: (modo: ModoCalculo) => void;
	}

	let { modo = $bindable(), onModoChange }: Props = $props();

	function handleModoTabClick(tab: 'LISTADO' | 'MANUAL_AMPERAJE' | 'MANUAL_POTENCIA') {
		onModoChange(tab);
	}

	function handleSubModoChange(sub: 'MANUAL_AMPERAJE' | 'MANUAL_POTENCIA') {
		onModoChange(sub);
	}
</script>

<div class="flex flex-col gap-4">
	<!-- Tabs principales -->
	<div class="flex rounded-lg bg-muted p-1">
		<button
			type="button"
			class={cn(
				'flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors',
				modo === 'LISTADO'
					? 'bg-primary text-primary-foreground shadow-sm'
					: 'text-muted-foreground hover:text-foreground'
			)}
			onclick={() => handleModoTabClick('LISTADO')}
		>
			Seleccionar desde Equipo
		</button>
		<button
			type="button"
			class={cn(
				'flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors',
				modo === 'MANUAL_AMPERAJE' || modo === 'MANUAL_POTENCIA'
					? 'bg-primary text-primary-foreground shadow-sm'
					: 'text-muted-foreground hover:text-foreground'
			)}
			onclick={() => handleModoTabClick('MANUAL_AMPERAJE')}
		>
			Carga Manual
		</button>
	</div>

	<!-- Sub-opciones cuando está en modo manual -->
	{#if modo === 'MANUAL_AMPERAJE' || modo === 'MANUAL_POTENCIA'}
		<div class="flex gap-4 pl-2">
			<label class="flex cursor-pointer items-center gap-2">
				<input
					type="radio"
					name="subModoManual"
					value="MANUAL_AMPERAJE"
					checked={modo === 'MANUAL_AMPERAJE'}
					class="h-4 w-4 border-border text-primary focus:ring-primary"
					onchange={() => handleSubModoChange('MANUAL_AMPERAJE')}
				/>
				<span class="text-sm text-foreground">Por Amperaje</span>
			</label>
			<label class="flex cursor-pointer items-center gap-2">
				<input
					type="radio"
					name="subModoManual"
					value="MANUAL_POTENCIA"
					checked={modo === 'MANUAL_POTENCIA'}
					class="h-4 w-4 border-border text-primary focus:ring-primary"
					onchange={() => handleSubModoChange('MANUAL_POTENCIA')}
				/>
				<span class="text-sm text-foreground">Por Potencia</span>
			</label>
		</div>
	{/if}
</div>
