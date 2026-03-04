<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { redirect } from '@sveltejs/kit';
	import type { MemoriaOutput } from '$lib/types/calculos.types';
	import FormularioPdf from '$lib/components/pdf/FormularioPdf.svelte';

	// Extraer memoria del estado de navegación
	let memoria = $derived($page.state.memoria as MemoriaOutput | undefined);

	// Redirect si no hay datos
	// svelte-ignore state_referenced_locally
	// Validación en tiempo de inicialización, no necesita reactividad
	if (!memoria) {
		throw redirect(302, '/calculos/resultado');
	}

	// Codificar datos para pasar como query param
	function volverAResultados() {
		try {
			// Generar ID único para almacenar en sessionStorage
			const resultId = crypto.randomUUID();
			sessionStorage.setItem(`memoria-${resultId}`, JSON.stringify(memoria));
			goto(`/calculos/resultado?id=${resultId}`);
		} catch (err) {
			console.error('Error al preparar datos para resultados:', err);
			// Fallback: ir sin datos (mostrará error al usuario)
			goto('/calculos/resultado');
		}
	}

	function nuevoCalculo() {
		goto('/');
	}
</script>

<svelte:head>
	<title>Generar PDF — Memoria de Cálculo</title>
</svelte:head>

<div class="min-h-screen bg-background px-4 py-8">
	<div class="mx-auto max-w-2xl">
		<!-- Header -->
		<header class="mb-8 flex items-center justify-between">
			<div>
				<button
					onclick={volverAResultados}
					class="inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
				>
					<svg
						class="size-4"
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 20 20"
						fill="currentColor"
						aria-hidden="true"
					>
						<path
							fill-rule="evenodd"
							d="M17 10a.75.75 0 01-.75.75H5.612l4.158 3.96a.75.75 0 11-1.04 1.08l-5.5-5.25a.75.75 0 010-1.08l5.5-5.25a.75.75 0 111.04 1.08L5.612 9.25H16.25A.75.75 0 0117 10z"
							clip-rule="evenodd"
						/>
					</svg>
					Volver a Resultados
				</button>
				<h1 class="mt-4 text-2xl font-bold text-foreground">Generar memoria de cálculo</h1>
				<p class="mt-1 text-sm text-muted-foreground">
					Complete los datos de presentación para generar el documento PDF.
				</p>
			</div>
			<button
				onclick={nuevoCalculo}
				class="rounded-md border border-border bg-card px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
			>
				Nuevo Cálculo
			</button>
		</header>

		<!-- Formulario de configuración PDF -->
		<FormularioPdf {memoria} />
	</div>
</div>
