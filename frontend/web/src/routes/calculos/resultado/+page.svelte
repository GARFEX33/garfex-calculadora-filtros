<script lang="ts">
	import type { PageData } from './$types';
	import { goto } from '$app/navigation';
	import MemoriaTecnica from '$lib/components/calculos/MemoriaTecnica.svelte';

	let { data }: { data: PageData } = $props();

	function irAGenerarPdf() {
		goto('/calculos/resultado/pdf', { state: { memoria: data.memoria } });
	}
</script>

<svelte:head>
	<title>Memoria de Cálculo - Resultados</title>
</svelte:head>

<div class="min-h-screen bg-background px-4 py-8">
	<div class="mx-auto max-w-4xl">
		<!-- Header with action buttons -->
		<header class="mb-8 flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-bold text-foreground">Memoria de Cálculo</h1>
				<p class="text-sm text-muted-foreground">Resultados del cálculo eléctrico</p>
			</div>
			<div class="flex items-center gap-2">
				<!-- Generar PDF button -->
				<button
					onclick={irAGenerarPdf}
					class="inline-flex items-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary-hover focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
					title="Generar memoria de cálculo en PDF"
				>
					<!-- PDF icon -->
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
					<!-- Text: visible on desktop, hidden on mobile -->
					<span class="hidden sm:inline">Generar PDF</span>
				</button>

				<!-- Nuevo Cálculo button -->
				<button
					onclick={() => goto('/')}
					class="rounded-md border border-border bg-card px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent focus-visible:ring-2 focus-visible:ring-ring focus-visible:outline-none"
				>
					Nuevo Cálculo
				</button>
			</div>
		</header>

		<!-- Technical Memory Content -->
		<MemoriaTecnica memoria={data.memoria} />
	</div>
</div>
