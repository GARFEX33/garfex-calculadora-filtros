<script lang="ts">
	import { onMount } from 'svelte';
	import { cn } from '$lib/utils';
	import { listarEquipos } from '$lib/api/equipos';
	import type { EquipoFiltro } from '$lib/types/equipos.types';

	interface Props {
		equipoSeleccionado: EquipoFiltro | undefined;
		onEquipoChange: (equipo: EquipoFiltro | undefined) => void;
	}

	let { equipoSeleccionado = $bindable(), onEquipoChange }: Props = $props();

	let equipos = $state<EquipoFiltro[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let searchTerm = $state('');
	let currentPage = $state(1);
	let totalPages = $state(1);
	let mounted = $state(false);

	// Debounce timer for real-time search
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;

	onMount(async () => {
		await cargarEquipos();
		mounted = true;
	});

	// Debounced search effect - triggers 300ms after searchTerm changes
	$effect(() => {
		// Skip effect during initial mount - equipos were already loaded
		if (!mounted) return;

		// Clear existing timer
		if (debounceTimer) {
			clearTimeout(debounceTimer);
		}

		// Set new timer for debounced search
		// The searchTerm dependency is implicitly tracked by accessing it in cargarEquipos
		debounceTimer = setTimeout(() => {
			// Access searchTerm to maintain reactivity tracking
			void searchTerm;
			cargarEquipos(1);
		}, 300);
	});

	// Cleanup timer on unmount
	onMount(() => {
		return () => {
			if (debounceTimer) {
				clearTimeout(debounceTimer);
			}
		};
	});

	async function cargarEquipos(page = 1) {
		loading = true;
		error = null;
		currentPage = page;

		const params: { page: number; per_page: number; buscar?: string } = {
			page,
			per_page: 20
		};
		if (searchTerm) {
			params.buscar = searchTerm;
		}

		const result = await listarEquipos(params);

		if (result.ok) {
			equipos = result.data.data.equipos;
			totalPages = result.data.data.pagination.total_pages;
		} else {
			error = result.error.error || 'Error al cargar equipos';
		}
		loading = false;
	}

	async function handleSearch() {
		await cargarEquipos(1);
	}

	async function handlePreviousPage() {
		if (currentPage > 1) {
			await cargarEquipos(currentPage - 1);
		}
	}

	async function handleNextPage() {
		if (currentPage < totalPages) {
			await cargarEquipos(currentPage + 1);
		}
	}

	function handleSelectEquipo(equipo: EquipoFiltro) {
		onEquipoChange(equipo);
	}

	function getTipoBadgeClass(tipo: string): string {
		switch (tipo) {
			case 'A':
				return 'bg-primary text-primary-foreground';
			case 'KVA':
				return 'bg-accent text-accent-foreground';
			case 'KVAR':
				return 'bg-warning text-warning-foreground';
			default:
				return 'bg-muted text-muted-foreground';
		}
	}
</script>

<div class="flex flex-col gap-4">
	<!-- Buscador -->
	<div class="flex gap-2">
		<input
			type="text"
			placeholder="Buscar por clave (ej: 400, 48D)..."
			bind:value={searchTerm}
			onkeydown={(e) => e.key === 'Enter' && handleSearch()}
			class="flex-1 rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:outline-none"
		/>
		<button
			type="button"
			onclick={handleSearch}
			class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary-hover"
		>
			Buscar
		</button>
	</div>

	<!-- Estados: Loading, Error, Empty, Lista -->
	{#if loading}
		<div class="flex items-center justify-center py-8">
			<p class="text-sm text-muted-foreground">Cargando equipos...</p>
		</div>
	{:else if error}
		<div class="flex flex-col items-center gap-2 py-8">
			<p class="text-sm text-destructive">{error}</p>
			<button
				type="button"
				onclick={() => cargarEquipos(currentPage)}
				class="rounded-md bg-secondary px-4 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80"
			>
				Reintentar
			</button>
		</div>
	{:else if equipos.length === 0}
		<div class="flex items-center justify-center py-8">
			<p class="text-sm text-muted-foreground">No se encontraron equipos</p>
		</div>
	{:else}
		<!-- Lista de equipos -->
		<div class="max-h-64 overflow-y-auto rounded-md border border-border">
			{#each equipos as equipo}
				<label
					class={cn(
						'flex cursor-pointer items-start gap-3 border-b border-border p-3 transition-colors last:border-b-0 hover:bg-muted',
						equipoSeleccionado?.id === equipo.id && 'bg-primary/10'
					)}
				>
					<input
						type="radio"
						name="equipo"
						checked={equipoSeleccionado?.id === equipo.id}
						onchange={() => handleSelectEquipo(equipo)}
						class="mt-1 h-4 w-4 border-border text-primary focus:ring-primary"
					/>
					<div class="min-w-0 flex-1">
						<div class="flex items-center gap-2">
							<span class="font-mono text-sm font-medium text-foreground">{equipo.clave}</span>
							<span
								class={cn(
									'rounded px-1.5 py-0.5 text-xs font-medium',
									getTipoBadgeClass(equipo.tipo)
								)}
							>
								{equipo.tipo}
							</span>
						</div>
						<div class="mt-1 flex flex-wrap gap-x-3 text-xs text-muted-foreground">
							<span>{equipo.voltaje}V</span>
							<span>Qn: {equipo.amperaje}</span>
							<span>ITM: {equipo.itm}A</span>
							{#if equipo.bornes}
								<span>Bornes: {equipo.bornes}</span>
							{/if}
						</div>
					</div>
				</label>
			{/each}
		</div>

		<!-- Paginación -->
		{#if totalPages > 1}
			<div class="flex items-center justify-between">
				<button
					type="button"
					disabled={currentPage === 1}
					onclick={handlePreviousPage}
					class="rounded-md border border-border bg-card px-3 py-1.5 text-sm text-foreground transition-colors hover:bg-muted disabled:cursor-not-allowed disabled:opacity-50"
				>
					Anterior
				</button>
				<span class="text-sm text-muted-foreground">
					Página {currentPage} de {totalPages}
				</span>
				<button
					type="button"
					disabled={currentPage === totalPages}
					onclick={handleNextPage}
					class="rounded-md border border-border bg-card px-3 py-1.5 text-sm text-foreground transition-colors hover:bg-muted disabled:cursor-not-allowed disabled:opacity-50"
				>
					Siguiente
				</button>
			</div>
		{/if}
	{/if}

	<!-- Equipo seleccionado - Card de detalles -->
	{#if equipoSeleccionado}
		<div class="rounded-lg border border-border bg-card p-4">
			<h4 class="mb-3 text-sm font-semibold text-card-foreground">Equipo Seleccionado</h4>
			<dl class="grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
				<div>
					<dt class="text-muted-foreground">Clave</dt>
					<dd class="font-mono font-medium text-foreground">{equipoSeleccionado.clave}</dd>
				</div>
				<div>
					<dt class="text-muted-foreground">Tipo</dt>
					<dd
						class={cn(
							'inline-flex rounded px-2 py-0.5 text-xs font-medium',
							getTipoBadgeClass(equipoSeleccionado.tipo)
						)}
					>
						{equipoSeleccionado.tipo}
					</dd>
				</div>
				<div>
					<dt class="text-muted-foreground">Voltaje</dt>
					<dd class="text-foreground">{equipoSeleccionado.voltaje} V</dd>
				</div>
				<div>
					<dt class="text-muted-foreground">Qn/In</dt>
					<dd class="text-foreground">{equipoSeleccionado.amperaje}</dd>
				</div>
				<div>
					<dt class="text-muted-foreground">ITM</dt>
					<dd class="text-foreground">{equipoSeleccionado.itm} A</dd>
				</div>
				{#if equipoSeleccionado.bornes}
					<div>
						<dt class="text-muted-foreground">Bornes</dt>
						<dd class="text-foreground">{equipoSeleccionado.bornes}</dd>
					</div>
				{/if}
				{#if equipoSeleccionado.conexion}
					<div>
						<dt class="text-muted-foreground">Conexión</dt>
						<dd class="text-foreground">{equipoSeleccionado.conexion}</dd>
					</div>
				{/if}
				{#if equipoSeleccionado.tipo_voltaje}
					<div>
						<dt class="text-muted-foreground">Tipo Voltaje</dt>
						<dd class="text-foreground">{equipoSeleccionado.tipo_voltaje}</dd>
					</div>
				{/if}
			</dl>
		</div>
	{/if}
</div>
