<script lang="ts">
	import type { ModoCalculo } from '$lib/features/calculos/domain/types/memoria.types';
	import type { EquipoFiltro } from '$lib/features/equipos/domain/types/equipo.types';
	import type { FormularioManualData } from '$lib/components/calculos/FormularioManual.svelte';
	import type { CamposInstalacionData } from '$lib/components/calculos/CamposInstalacion.svelte';
	import type { DatosEquipo } from '$lib/features/calculos/domain/types/memoria.types';

	import SelectorModo from '$lib/components/calculos/SelectorModo.svelte';
	import FormularioManual from '$lib/components/calculos/FormularioManual.svelte';
	import FormularioListado from '$lib/components/calculos/FormularioListado.svelte';
	import CamposInstalacion from '$lib/components/calculos/CamposInstalacion.svelte';

	import { memoriaStore } from '$lib/features/calculos/application/stores/memoria.store.svelte';
	import { equiposStore } from '$lib/features/equipos/application/stores/equipos.store.svelte';
	import {
		mapConexionToSistemaElectrico,
		mapTipoVoltajeToTipoVoltaje
	} from '$lib/features/calculos/infrastructure/mappers/memoria.mapper';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	// ── Estado Local (solo UI) ───────────────────────────────────────────────────
	// Modo - se sincroniza con el store
	let modo = $state<ModoCalculo>(memoriaStore.input.modo);

	// Datos del formulario manual (UI state)
	let datosManual = $state<FormularioManualData>({
		tipo_equipo: '',
		amperaje_nominal: undefined,
		potencia_nominal: undefined,
		potencia_unidad: 'KW',
		factor_potencia: undefined,
		itm: undefined,
		tension: 220,
		tipo_voltaje: 'FASE_NEUTRO',
		sistema_electrico: ''
	});

	// Datos de instalación (UI state)
	let instalacion = $state<CamposInstalacionData>({
		tension: undefined,
		tension_unidad: 'V',
		sistema_electrico: '',
		estado: '',
		tipo_canalizacion: '',
		num_tuberias: undefined,
		longitud_circuito: undefined,
		tipo_voltaje: '',
		material: 'CU',
		hilos_por_fase: 1,
		porcentaje_caida_maximo: 3.0,
		temperatura_override: undefined,
		diametro_control_mm: undefined
	});

	// Equipo seleccionado en modo LISTADO
	let equipoSeleccionado = $state<EquipoFiltro | undefined>(undefined);

	// ── Derived Values ───────────────────────────────────────────────────────────
	let esModoManual = $derived(modo === 'MANUAL_AMPERAJE' || modo === 'MANUAL_POTENCIA');
	let esModoListado = $derived(modo === 'LISTADO');

	// Derived from store
	let loading = $derived(memoriaStore.loading);
	let error = $derived(memoriaStore.error);

	// ── Handlers ─────────────────────────────────────────────────────────────────
	function handleModoChange(newModo: ModoCalculo) {
		modo = newModo;
		memoriaStore.actualizarInput({ modo: newModo });

		// Reset incompatible data when switching modes
		if (newModo === 'LISTADO') {
			datosManual = {
				tipo_equipo: '',
				amperaje_nominal: undefined,
				potencia_nominal: undefined,
				potencia_unidad: 'KW',
				factor_potencia: undefined,
				itm: undefined,
				tension: 220,
				tipo_voltaje: 'FASE_NEUTRO',
				sistema_electrico: ''
			};
		} else {
			equipoSeleccionado = undefined;
			// Clear installation fields that were auto-filled from equipment
			instalacion = {
				tension: undefined,
				tension_unidad: 'V',
				sistema_electrico: '',
				estado: instalacion.estado,
				tipo_canalizacion: '',
				num_tuberias: undefined,
				longitud_circuito: undefined,
				tipo_voltaje: '',
				material: 'CU',
				hilos_por_fase: 1,
				porcentaje_caida_maximo: 3.0,
				temperatura_override: undefined,
				diametro_control_mm: undefined
			};
		}
	}

	function handleDatosManualChange(newDatos: FormularioManualData) {
		datosManual = newDatos;

		// Update store with manual-specific fields
		const update: Record<string, number | string | undefined> = {
			tension: newDatos.tension,
			tipo_voltaje: newDatos.tipo_voltaje,
			sistema_electrico: newDatos.sistema_electrico
		};

		if (newDatos.tipo_equipo) update['tipo_equipo'] = newDatos.tipo_equipo;
		if (newDatos.amperaje_nominal !== undefined)
			update['amperaje_nominal'] = newDatos.amperaje_nominal;
		if (newDatos.potencia_nominal !== undefined)
			update['potencia_nominal'] = newDatos.potencia_nominal;
		if (newDatos.potencia_unidad) update['potencia_unidad'] = newDatos.potencia_unidad;
		if (newDatos.factor_potencia !== undefined)
			update['factor_potencia'] = newDatos.factor_potencia;
		if (newDatos.itm !== undefined) update['itm'] = newDatos.itm;

		memoriaStore.actualizarInput(update as Parameters<typeof memoriaStore.actualizarInput>[0]);
	}

	function handleEquipoChange(equipo: EquipoFiltro | undefined) {
		equipoSeleccionado = equipo;

		if (equipo) {
			// Auto-fill installation fields from selected equipment
			const sistemaElectrico = mapConexionToSistemaElectrico(equipo.conexion);
			const tipoVoltaje = mapTipoVoltajeToTipoVoltaje(equipo.tipo_voltaje);

			instalacion = {
				...instalacion,
				tension: equipo.voltaje,
				tension_unidad: 'V',
				sistema_electrico: sistemaElectrico ?? '',
				tipo_voltaje: tipoVoltaje ?? ''
			};

			// Build equipment data for store
			const equipoData: DatosEquipo = {
				clave: equipo.clave,
				tipo: equipo.tipo,
				voltaje: equipo.voltaje,
				amperaje: equipo.amperaje,
				itm: equipo.itm
			};
			if (equipo.bornes !== undefined && equipo.bornes !== null) {
				equipoData.bornes = equipo.bornes;
			}

			// Update store with equipment data (only include valid values)
			const update: Parameters<typeof memoriaStore.actualizarInput>[0] = {
				equipo: equipoData,
				tension: equipo.voltaje
			};

			// Only add sistema_electrico and tipo_voltaje if they have valid values
			if (sistemaElectrico) {
				update['sistema_electrico'] = sistemaElectrico;
			}
			if (tipoVoltaje) {
				update['tipo_voltaje'] = tipoVoltaje;
			}

			memoriaStore.actualizarInput(update);
		} else {
			// Clear equipment - use a empty object trick or just don't pass it
			// Since we can't set to undefined with exactOptionalPropertyTypes,
			// we need to reset the entire input for this mode
			const currentInput = memoriaStore.input;
			// eslint-disable-next-line @typescript-eslint/no-unused-vars
			const { equipo: _, ...inputWithoutEquipo } = currentInput;
			memoriaStore.input = inputWithoutEquipo;
		}
	}

	function handleInstalacionChange(newDatos: CamposInstalacionData) {
		instalacion = newDatos;

		// Update store with installation fields
		const update: Record<string, number | string | undefined> = {
			// Solo incluir tensión si tiene un valor válido (en modo MANUAL viene del formulario manual)
			...(newDatos.tension !== undefined && newDatos.tension > 0 && { tension: newDatos.tension }),
			tension_unidad: newDatos.tension_unidad,
			sistema_electrico: newDatos.sistema_electrico || undefined,
			estado: newDatos.estado,
			tipo_canalizacion: newDatos.tipo_canalizacion,
			longitud_circuito: newDatos.longitud_circuito ?? 0,
			tipo_voltaje: newDatos.tipo_voltaje || undefined,
			material: newDatos.material,
			hilos_por_fase: newDatos.hilos_por_fase,
			porcentaje_caida_maximo: newDatos.porcentaje_caida_maximo
		};

		if (newDatos.temperatura_override !== undefined) {
			update['temperatura_override'] = newDatos.temperatura_override;
		}
		if (newDatos.diametro_control_mm !== undefined) {
			update['diametro_control_mm'] = newDatos.diametro_control_mm;
		}
		if (newDatos.num_tuberias !== undefined && newDatos.num_tuberias > 0) {
			update['num_tuberias'] = newDatos.num_tuberias;
		}

		memoriaStore.actualizarInput(update as Parameters<typeof memoriaStore.actualizarInput>[0]);
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();

		await memoriaStore.calcular();

		if (memoriaStore.output) {
			// Serializar output ANTES de resetear
			const jsonStr = JSON.stringify(memoriaStore.output);
			const encodedData = btoa(unescape(encodeURIComponent(jsonStr)));
			// Resetear store para que al volver la página esté limpia
			memoriaStore.resetear();
			goto(`/calculos/resultado?data=${encodedData}`);
		}
	}

	// ── Equipos Store Integration ─────────────────────────────────────────────────
	// Handler for search from FormularioListado
	async function handleBusquedaChange(query: string) {
		await equiposStore.buscar(query);
	}

	// Handler for page change from FormularioListado
	async function handlePaginaChange(pagina: number) {
		await equiposStore.cambiarPagina(pagina);
	}

	// Resetear store y estados locales al montar la página
	// Garantiza formulario limpio si el usuario vuelve desde resultados
	onMount(() => {
		memoriaStore.resetear();
		datosManual = {
			tipo_equipo: '',
			amperaje_nominal: undefined,
			potencia_nominal: undefined,
			potencia_unidad: 'KW',
			factor_potencia: undefined,
			itm: undefined,
			tension: 220,
			tipo_voltaje: 'FASE_NEUTRO',
			sistema_electrico: ''
		};
		instalacion = {
			tension: undefined,
			tension_unidad: 'V',
			sistema_electrico: '',
			estado: '',
			tipo_canalizacion: '',
			num_tuberias: undefined,
			longitud_circuito: undefined,
			tipo_voltaje: '',
			material: 'CU',
			hilos_por_fase: 1,
			porcentaje_caida_maximo: 3.0,
			temperatura_override: undefined,
			diametro_control_mm: undefined
		};
		equipoSeleccionado = undefined;
	});

	// Load equipos when mode changes to LISTADO
	$effect(() => {
		if (modo === 'LISTADO') {
			// Only load if we haven't loaded yet or if we need fresh data
			if (equiposStore.equipos.length === 0 && !equiposStore.loading) {
				void equiposStore.cargar();
			}
		}
	});
</script>

<svelte:head>
	<title>Calculadora de Filtros Eléctricos - Memoria de Cálculo NOM</title>
</svelte:head>

<div class="min-h-screen bg-background px-4 py-8">
	<div class="mx-auto max-w-4xl">
		<!-- Header -->
		<header class="mb-8 text-center">
			<h1 class="text-3xl font-bold text-foreground">Calculadora de Filtros Eléctricos</h1>
			<p class="mt-1 text-muted-foreground">Memoria de Cálculo NOM</p>
		</header>

		<!-- Main Form Card -->
		<form onsubmit={handleSubmit} class="rounded-xl border border-border bg-card p-6 shadow-sm">
			<!-- Mode Selector -->
			<div class="mb-6">
				<SelectorModo bind:modo onModoChange={handleModoChange} />
			</div>

			<!-- Dynamic Form Based on Mode -->
			<div class="mb-6 space-y-6">
				{#if esModoManual}
					<FormularioManual
						modo={modo as 'MANUAL_AMPERAJE' | 'MANUAL_POTENCIA'}
						bind:datos={datosManual}
						onDatosChange={handleDatosManualChange}
					/>
				{:else if esModoListado}
					<FormularioListado
						bind:equipoSeleccionado
						onEquipoChange={handleEquipoChange}
						equipos={equiposStore.equipos}
						totalEquipos={equiposStore.total}
						loading={equiposStore.loading}
						error={equiposStore.error}
						onBusquedaChange={handleBusquedaChange}
						onPaginaChange={handlePaginaChange}
						externalData={true}
					/>
				{/if}
			</div>

			<!-- Installation Fields (always visible) -->
			<div class="mb-6">
				<CamposInstalacion bind:datos={instalacion} onDatosChange={handleInstalacionChange} />
			</div>

			<!-- Submit Button -->
			<button
				type="submit"
				disabled={loading}
				class="flex w-full items-center justify-center gap-2 rounded-md bg-primary px-6 py-3 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary-hover disabled:cursor-not-allowed disabled:opacity-50 md:w-auto md:px-8"
			>
				{#if loading}
					<svg class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
						></circle>
						<path
							class="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
						></path>
					</svg>
					Calculando...
				{:else}
					Calcular Memoria
				{/if}
			</button>
		</form>

		<!-- Error Display -->
		{#if error}
			<div class="mt-6 rounded-lg border border-destructive/50 bg-destructive/10 p-4">
				<p class="text-sm font-medium text-destructive">{error}</p>
			</div>
		{/if}

		<!-- Validation Errors Display -->
		{#if memoriaStore.erroresValidacion.length > 0}
			<div class="mt-4 rounded-lg border border-red-500 bg-red-50 p-4">
				<p class="mb-2 text-sm font-medium text-red-800">
					Por favor complete los siguientes campos:
				</p>
				<ul class="list-inside list-disc text-sm text-red-700">
					{#each memoriaStore.erroresValidacion as err}
						<li>{err.field}: {err.message}</li>
					{/each}
				</ul>
			</div>
		{/if}
	</div>
</div>
