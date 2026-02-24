<script lang="ts">
	import type { PageData } from './$types';
	import { goto } from '$app/navigation';
	import MemoriaTecnica from '$lib/components/calculos/MemoriaTecnica.svelte';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Memoria de Cálculo - Resultados</title>
</svelte:head>

<div class="min-h-screen bg-background px-4 py-8">
	<div class="mx-auto max-w-4xl">
		<!-- Header with back button -->
		<header class="mb-8 flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-bold text-foreground">Memoria de Cálculo</h1>
				<p class="text-sm text-muted-foreground">Resultados del cálculo eléctrico</p>
			</div>
			<button
				onclick={() => goto('/')}
				class="rounded-md border border-border bg-card px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent"
			>
				Nuevo Cálculo
			</button>
		</header>

		<!-- Technical Memory Content -->
		<MemoriaTecnica memoria={data.memoria} />
	</div>
</div>
