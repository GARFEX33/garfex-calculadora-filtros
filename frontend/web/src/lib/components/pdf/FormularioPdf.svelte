<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';
	import { generarMemoriaPdf, descargarBlob } from '$lib/api/pdf';
	import { EMPRESA_DEFAULT_ID } from '$lib/config/empresas-pdf';
	import { cn } from '$lib/utils';
	import SelectorEmpresa from './SelectorEmpresa.svelte';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	// ── Campos del formulario ────────────────────────────────────────────────
	let empresaId = $state(EMPRESA_DEFAULT_ID);
	let nombreProyecto = $state('');
	let direccionProyecto = $state('');
	let responsable = $state('');
	// svelte-ignore state_referenced_locally
	// Se usa valor inicial como prefill, no es necesario reaccivo
	let nombreEquipo = $state(memoria.equipo?.clave ?? '');

	// ── Estado de UI ─────────────────────────────────────────────────────────
	let cargando = $state(false);
	let errorMensaje = $state<string | null>(null);

	// ── Errores de validación inline ─────────────────────────────────────────
	let errores = $state<{ nombreProyecto?: string; responsable?: string }>({});

	function validar(): boolean {
		const nuevosErrores: typeof errores = {};

		if (!nombreProyecto.trim()) {
			nuevosErrores.nombreProyecto = 'El nombre del proyecto es obligatorio';
		} else if (nombreProyecto.length > 200) {
			nuevosErrores.nombreProyecto = 'Máximo 200 caracteres';
		}

		if (!responsable.trim()) {
			nuevosErrores.responsable = 'El nombre del responsable es obligatorio';
		} else if (responsable.length > 100) {
			nuevosErrores.responsable = 'Máximo 100 caracteres';
		}

		errores = nuevosErrores;
		return Object.keys(nuevosErrores).length === 0;
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		errorMensaje = null;

		if (!validar()) return;

		cargando = true;
		try {
			const blob = await generarMemoriaPdf({
				memoria,
				presentacion: {
					empresa_id: empresaId,
					nombre_proyecto: nombreProyecto.trim(),
					direccion_proyecto: direccionProyecto.trim(),
					responsable: responsable.trim(),
					...(nombreEquipo.trim() ? { nombre_equipo_override: nombreEquipo.trim() } : {})
				}
			});

			// Nombre descriptivo para el archivo descargado
			const fecha = new Date().toISOString().slice(0, 10).replace(/-/g, '');
			const proyectoSanitizado = nombreProyecto
				.trim()
				.replace(/\s+/g, '_')
				.replace(/[^a-zA-Z0-9_-]/g, '');
			const equipoSanitizado = (memoria.equipo?.clave ?? 'Equipo')
				.replace(/\s+/g, '_')
				.replace(/[^a-zA-Z0-9_-]/g, '');
			const filename = `MemoriaCalculo_${proyectoSanitizado}_${equipoSanitizado}_${fecha}.pdf`;

			descargarBlob(blob, filename);
		} catch (err) {
			errorMensaje = err instanceof Error ? err.message : 'Error desconocido al generar el PDF';
		} finally {
			cargando = false;
		}
	}
</script>

<form onsubmit={handleSubmit} class="space-y-6" novalidate>
	<!-- Sección 1: Empresa -->
	<section class="rounded-lg border border-border bg-card p-5">
		<h2 class="mb-4 text-base font-semibold text-foreground">Empresa presentadora</h2>
		<SelectorEmpresa bind:empresaId />
	</section>

	<!-- Sección 2: Datos del proyecto -->
	<section class="rounded-lg border border-border bg-card p-5">
		<h2 class="mb-4 text-base font-semibold text-foreground">Datos del proyecto</h2>

		<div class="space-y-4">
			<!-- Nombre del proyecto (requerido) -->
			<div class="space-y-1.5">
				<label for="nombre_proyecto" class="text-sm font-medium text-foreground">
					Nombre del proyecto
					<span class="ml-0.5 text-destructive" aria-hidden="true">*</span>
				</label>
				<input
					id="nombre_proyecto"
					type="text"
					bind:value={nombreProyecto}
					maxlength={200}
					placeholder="Ej. Planta de producción norte — tablero TG-01"
					aria-required="true"
					aria-invalid={!!errores.nombreProyecto}
					aria-describedby={errores.nombreProyecto ? 'error-nombre-proyecto' : undefined}
					class={cn(
						'flex h-10 w-full rounded-md border bg-input px-3 py-2 text-sm placeholder:text-muted-foreground',
						'focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none',
						'disabled:cursor-not-allowed disabled:opacity-50',
						errores.nombreProyecto ? 'border-destructive' : 'border-input-border'
					)}
				/>
				{#if errores.nombreProyecto}
					<p id="error-nombre-proyecto" class="text-sm text-destructive" role="alert">
						{errores.nombreProyecto}
					</p>
				{/if}
			</div>

			<!-- Dirección del proyecto (opcional) -->
			<div class="space-y-1.5">
				<label for="direccion_proyecto" class="text-sm font-medium text-foreground">
					Dirección de la instalación
					<span class="ml-1 text-xs text-muted-foreground">(opcional)</span>
				</label>
				<input
					id="direccion_proyecto"
					type="text"
					bind:value={direccionProyecto}
					maxlength={300}
					placeholder="Ej. Av. Industrial 500, Parque Industrial, Monterrey N.L."
					class="flex h-10 w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50"
				/>
			</div>
		</div>
	</section>

	<!-- Sección 3: Datos del reporte -->
	<section class="rounded-lg border border-border bg-card p-5">
		<h2 class="mb-4 text-base font-semibold text-foreground">Datos del reporte</h2>

		<div class="space-y-4">
			<!-- Responsable (requerido) -->
			<div class="space-y-1.5">
				<label for="responsable" class="text-sm font-medium text-foreground">
					Responsable del cálculo
					<span class="ml-0.5 text-destructive" aria-hidden="true">*</span>
				</label>
				<input
					id="responsable"
					type="text"
					bind:value={responsable}
					maxlength={100}
					placeholder="Nombre completo del ingeniero o técnico"
					aria-required="true"
					aria-invalid={!!errores.responsable}
					aria-describedby={errores.responsable ? 'error-responsable' : undefined}
					class={cn(
						'flex h-10 w-full rounded-md border bg-input px-3 py-2 text-sm placeholder:text-muted-foreground',
						'focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none',
						'disabled:cursor-not-allowed disabled:opacity-50',
						errores.responsable ? 'border-destructive' : 'border-input-border'
					)}
				/>
				{#if errores.responsable}
					<p id="error-responsable" class="text-sm text-destructive" role="alert">
						{errores.responsable}
					</p>
				{/if}
			</div>

			<!-- Nombre del equipo (pre-llenado, editable) -->
			<div class="space-y-1.5">
				<label for="nombre_equipo" class="text-sm font-medium text-foreground">
					Identificación del equipo
					<span class="ml-1 text-xs text-muted-foreground">(editable)</span>
				</label>
				<input
					id="nombre_equipo"
					type="text"
					bind:value={nombreEquipo}
					maxlength={100}
					placeholder="Clave o descripción del equipo"
					class="flex h-10 w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50"
				/>
				<p class="text-xs text-muted-foreground">
					Dejar en blanco para usar la clave del equipo calculado ({memoria.equipo?.clave ?? '—'})
				</p>
			</div>
		</div>
	</section>

	<!-- Error global -->
	{#if errorMensaje}
		<div
			class="flex items-start gap-3 rounded-md border border-destructive/30 bg-destructive/10 px-4 py-3 text-sm text-destructive"
			role="alert"
		>
			<svg
				class="mt-0.5 size-4 shrink-0"
				xmlns="http://www.w3.org/2000/svg"
				viewBox="0 0 20 20"
				fill="currentColor"
				aria-hidden="true"
			>
				<path
					fill-rule="evenodd"
					d="M10 18a8 8 0 100-16 8 8 0 000 16zm-.75-4.75a.75.75 0 001.5 0v-4.5a.75.75 0 00-1.5 0v4.5zm.75-7a.75.75 0 100 1.5.75.75 0 000-1.5z"
					clip-rule="evenodd"
				/>
			</svg>
			<span>{errorMensaje}</span>
		</div>
	{/if}

	<!-- Botón submit -->
	<button
		type="submit"
		disabled={cargando}
		class={cn(
			'inline-flex w-full items-center justify-center gap-2 rounded-md px-6 py-3',
			'bg-primary text-sm font-semibold text-primary-foreground',
			'transition-colors hover:bg-primary-hover',
			'focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none',
			'disabled:cursor-not-allowed disabled:opacity-60'
		)}
	>
		{#if cargando}
			<!-- Spinner -->
			<svg
				class="size-4 animate-spin"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				aria-hidden="true"
			>
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
				></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
				></path>
			</svg>
			<span>Generando PDF...</span>
		{:else}
			<!-- Icono PDF -->
			<svg
				class="size-4"
				xmlns="http://www.w3.org/2000/svg"
				viewBox="0 0 20 20"
				fill="currentColor"
				aria-hidden="true"
			>
				<path
					fill-rule="evenodd"
					d="M4.5 2A1.5 1.5 0 003 3.5v13A1.5 1.5 0 004.5 18h11a1.5 1.5 0 001.5-1.5V7.621a1.5 1.5 0 00-.44-1.06l-4.12-4.122A1.5 1.5 0 0011.378 2H4.5zm4.75 6.75a.75.75 0 011.5 0v2.546l.943-1.048a.75.75 0 111.114 1.004l-2.25 2.5a.75.75 0 01-1.114 0l-2.25-2.5a.75.75 0 111.114-1.004l.943 1.048V8.75z"
					clip-rule="evenodd"
				/>
			</svg>
			<span>Generar PDF</span>
		{/if}
	</button>

	<p class="text-center text-xs text-muted-foreground">
		Los campos marcados con <span class="text-destructive">*</span> son obligatorios
	</p>
</form>
