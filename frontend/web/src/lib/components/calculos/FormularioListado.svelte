<script lang="ts">
	import { cn } from '$lib/utils';
	import type { EquipoFiltro } from '$lib/types/equipos.types';
	import { fade } from 'svelte/transition';

	interface Props {
		equipoSeleccionado: EquipoFiltro | undefined;
		onEquipoChange: (equipo: EquipoFiltro | undefined) => void;
		// Props para datos externos (en lugar de fetch interno)
		equipos?: EquipoFiltro[];
		totalEquipos?: number;
		loading?: boolean;
		error?: string | null;
		onBusquedaChange?: (query: string) => void;
		onPaginaChange?: (pagina: number) => void;
		// Indica si los datos se proporcionan externamente
		externalData?: boolean;
	}

	let {
		equipoSeleccionado = $bindable(),
		onEquipoChange,
		equipos = [],
		totalEquipos = 1,
		loading = false,
		error = null,
		onBusquedaChange,
		onPaginaChange,
		externalData = false
	}: Props = $props();

	// Estado local para búsqueda y paginación (si no hay datos externos)
	let searchTerm = $state('');
	let currentPage = $state(1);
	let totalPages = $state(1);

	// Sincronizar totalPages cuando cambian los datos
	$effect(() => {
		if (totalEquipos > 0) {
			totalPages = Math.ceil(totalEquipos / 20);
		}
	});

	// Manejar búsqueda local (sin datos externos)
	function handleSearchLocal() {
		if (onBusquedaChange) {
			onBusquedaChange(searchTerm);
		}
	}

	// Manejar cambio de página local
	function handlePreviousPage() {
		if (!externalData && onPaginaChange && currentPage > 1) {
			currentPage--;
			onPaginaChange(currentPage);
		}
	}

	function handleNextPage() {
		if (!externalData && onPaginaChange && currentPage < totalPages) {
			currentPage++;
			onPaginaChange(currentPage);
		}
	}

	// Debounce para búsqueda local
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;

	// Track previous searchTerm to detect actual changes
	let previousSearchTerm = $state('');

	// Search effect — only triggers when searchTerm actually changes from user input
	$effect(() => {
		// Use snapshot to track the actual value change
		const currentTerm = searchTerm;

		// Skip if this is the initial render (first run with empty string)
		// or if the value hasn't actually changed
		if (currentTerm === previousSearchTerm) {
			// On first run, just record the initial value
			if (previousSearchTerm === '') {
				previousSearchTerm = currentTerm;
			}
			return;
		}

		// Record the new value for next comparison
		previousSearchTerm = currentTerm;

		// Clear any pending search
		if (debounceTimer) {
			clearTimeout(debounceTimer);
		}

		// Debounce the actual search
		debounceTimer = setTimeout(() => {
			handleSearchLocal();
		}, 300);
	});

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
			onkeydown={(e) => e.key === 'Enter' && handleSearchLocal()}
			class="flex-1 rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:outline-none"
		/>
	</div>

	<!-- Estados: Loading, Error, Empty, Lista -->
	{#if loading}
		<div class="space-y-2">
			{#each Array(3) as _}
				<div class="animate-pulse rounded-md border border-border bg-card p-3">
					<div class="flex items-start gap-3">
						<div class="mt-1 h-4 w-4 rounded-full bg-muted"></div>
						<div class="flex-1 space-y-2">
							<div class="flex items-center gap-2">
								<div class="h-4 w-16 rounded bg-muted"></div>
								<div class="h-4 w-8 rounded bg-muted"></div>
							</div>
							<div class="flex gap-3">
								<div class="h-3 w-12 rounded bg-muted"></div>
								<div class="h-3 w-16 rounded bg-muted"></div>
								<div class="h-3 w-12 rounded bg-muted"></div>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{:else if error}
		<div class="flex flex-col items-center gap-2 py-8">
			<p class="text-sm text-destructive">{error}</p>
			<button
				type="button"
				onclick={() => onPaginaChange?.(currentPage)}
				class="rounded-md bg-secondary px-4 py-2 text-sm font-medium text-secondary-foreground transition-colors hover:bg-secondary/80"
			>
				Reintentar
			</button>
		</div>
	{:else if equipos.length === 0}
		<div class="flex flex-col items-center justify-center py-8 text-center">
			<svg
				class="mb-3 h-10 w-10 text-muted-foreground"
				fill="none"
				viewBox="0 0 24 24"
				stroke="currentColor"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="1.5"
					d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			<p class="text-sm font-medium text-foreground">No se encontraron equipos</p>
			<p class="text-xs text-muted-foreground">Try adjusting your search terms</p>
		</div>
	{:else}
		<!-- Lista de equipos -->
		<div
			class="max-h-64 overflow-y-auto rounded-md border border-border"
			transition:fade={{ duration: 150 }}
		>
			{#each equipos as equipo}
				<label
					class={cn(
						'flex cursor-pointer items-start gap-3 border-b border-border p-3 transition-all last:border-b-0 hover:bg-muted',
						equipoSeleccionado?.id === equipo.id ? 'border-l-2 border-l-primary bg-primary/5' : ''
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
