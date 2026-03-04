<script lang="ts">
	import { EMPRESAS_PDF, EMPRESA_DEFAULT_ID } from '$lib/config/empresas-pdf';
	import type { EmpresaPdf } from '$lib/config/empresas-pdf';
	import { cn } from '$lib/utils';

	interface Props {
		empresaId?: string;
		onchange?: (id: string) => void;
	}

	let { empresaId = $bindable(EMPRESA_DEFAULT_ID), onchange }: Props = $props();

	let empresaSeleccionada = $derived<EmpresaPdf | undefined>(
		EMPRESAS_PDF.find((e) => e.id === empresaId)
	);

	function seleccionar(id: string) {
		empresaId = id;
		onchange?.(id);
	}
</script>

<fieldset class="space-y-3">
	<legend class="mb-2 text-sm font-medium text-foreground">Empresa presentadora</legend>

	<div class="space-y-2">
		{#each EMPRESAS_PDF as empresa (empresa.id)}
			<label
				class={cn(
					'flex cursor-pointer items-center gap-3 rounded-md border px-4 py-3 transition-colors',
					empresaId === empresa.id
						? 'border-primary bg-primary/5 text-foreground'
						: 'border-border bg-card text-foreground hover:border-primary/50 hover:bg-muted/50'
				)}
			>
				<input
					type="radio"
					name="empresa_id"
					value={empresa.id}
					checked={empresaId === empresa.id}
					onchange={() => seleccionar(empresa.id)}
					class="size-4 accent-primary"
				/>
				<span class="text-sm font-medium">{empresa.nombre}</span>
			</label>
		{/each}
	</div>

	<!-- Vista previa de la empresa seleccionada -->
	{#if empresaSeleccionada}
		<div
			class="mt-3 rounded-md border border-border bg-muted/30 px-4 py-3 text-sm text-muted-foreground"
		>
			<p class="font-medium text-foreground">{empresaSeleccionada.nombre}</p>
			<p class="mt-1">{empresaSeleccionada.direccion}</p>
			<p class="mt-0.5">{empresaSeleccionada.telefono} · {empresaSeleccionada.email}</p>
		</div>
	{/if}
</fieldset>
