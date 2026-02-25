<script lang="ts">
	import type {
		ModoCalculo,
		CalcularMemoriaRequest,
		DatosEquipo,
		SistemaElectrico,
		TipoCanalizacion,
		MaterialConductor,
		TipoVoltaje
	} from '$lib/types/calculos.types';
	import type { EquipoFiltro } from '$lib/types/equipos.types';
	import type { FormularioManualData } from '$lib/components/calculos/FormularioManual.svelte';
	import type { CamposInstalacionData } from '$lib/components/calculos/CamposInstalacion.svelte';

	import SelectorModo from '$lib/components/calculos/SelectorModo.svelte';
	import FormularioManual from '$lib/components/calculos/FormularioManual.svelte';
	import FormularioListado from '$lib/components/calculos/FormularioListado.svelte';
	import CamposInstalacion from '$lib/components/calculos/CamposInstalacion.svelte';

	import { calcularMemoria } from '$lib/api/calculos';
	import { goto } from '$app/navigation';

	// ── Estado Principal ──────────────────────────────────────────────────────────
	let modo = $state<ModoCalculo>('MANUAL_AMPERAJE');

	// Manual form data
	let datosManual = $state<FormularioManualData>({
		tipo_equipo: '',
		amperaje_nominal: undefined,
		potencia_nominal: undefined,
		potencia_unidad: 'KW',
		factor_potencia: undefined,
		itm: undefined
	});

	// Listado form data
	let equipoSeleccionado = $state<EquipoFiltro | undefined>(undefined);

	// Installation fields
	let instalacion = $state<CamposInstalacionData>({
		tension: undefined,
		tension_unidad: 'V',
		sistema_electrico: '',
		estado: '',
		tipo_canalizacion: '',
		num_tuberias: undefined,
		longitud_circuito: undefined,
		tipo_voltaje: '',
		material: 'Cu',
		hilos_por_fase: 1,
		porcentaje_caida_maximo: 3.0,
		temperatura_override: undefined
	});

	// Submission state
	let loading = $state(false);
	let error = $state<string | undefined>(undefined);

	// ── Derived Values ────────────────────────────────────────────────────────────
	let esModoManual = $derived(modo === 'MANUAL_AMPERAJE' || modo === 'MANUAL_POTENCIA');
	let esModoListado = $derived(modo === 'LISTADO');

	// ── Mapping Functions ──────────────────────────────────────────────────────────
	// Maps equipment connection type to the expected format for installation
	function mapearConexionASistemaElectrico(
		conexion: string | null | undefined
	): SistemaElectrico | '' {
		if (!conexion) return '';

		const mapa: Record<string, SistemaElectrico> = {
			DELTA: 'DELTA',
			ESTRELLA: 'ESTRELLA',
			MONOFASICO: 'MONOFASICO',
			BIFASICO: 'BIFASICO'
		};
		return mapa[conexion] || '';
	}

	// Maps equipment voltage type to the expected format
	function mapearTipoVoltaje(tipoVoltaje: string | null | undefined): TipoVoltaje | '' {
		if (!tipoVoltaje) return '';

		const mapa: Record<string, TipoVoltaje> = {
			FF: 'FASE_FASE',
			FN: 'FASE_NEUTRO'
		};
		return mapa[tipoVoltaje] || '';
	}

	// ── Build Request ─────────────────────────────────────────────────────────────
	// Build CalcularMemoriaRequest from current state, only including optional fields with actual values
	let request = $derived.by((): CalcularMemoriaRequest => {
		const base: CalcularMemoriaRequest = {
			modo,
			tension: instalacion.tension ?? 0,
			tension_unidad: instalacion.tension_unidad,
			sistema_electrico: instalacion.sistema_electrico as SistemaElectrico,
			estado: instalacion.estado,
			tipo_canalizacion: instalacion.tipo_canalizacion as TipoCanalizacion,
			longitud_circuito: instalacion.longitud_circuito ?? 0,
			tipo_voltaje: instalacion.tipo_voltaje as TipoVoltaje,
			material: instalacion.material as MaterialConductor,
			hilos_por_fase: instalacion.hilos_por_fase,
			porcentaje_caida_maximo: instalacion.porcentaje_caida_maximo
		};

		// Add optional fields only if they have values
		if (equipoSeleccionado && modo === 'LISTADO') {
			const equipo: DatosEquipo = {
				clave: equipoSeleccionado.clave,
				tipo: equipoSeleccionado.tipo,
				voltaje: equipoSeleccionado.voltaje,
				amperaje: equipoSeleccionado.amperaje,
				itm: equipoSeleccionado.itm
			};
			// Only include bornes if it has a real value (null/undefined = omit)
			if (equipoSeleccionado.bornes !== null && equipoSeleccionado.bornes !== undefined) {
				equipo.bornes = equipoSeleccionado.bornes;
			}
			base.equipo = equipo;
		}

		if (modo === 'MANUAL_AMPERAJE' || modo === 'MANUAL_POTENCIA') {
			if (datosManual.tipo_equipo) {
				base.tipo_equipo = datosManual.tipo_equipo;
			}

			if (modo === 'MANUAL_AMPERAJE' && datosManual.amperaje_nominal !== undefined) {
				base.amperaje_nominal = datosManual.amperaje_nominal;
			}

			if (modo === 'MANUAL_POTENCIA') {
				if (datosManual.potencia_nominal !== undefined) {
					base.potencia_nominal = datosManual.potencia_nominal;
				}
				if (datosManual.potencia_unidad) {
					base.potencia_unidad = datosManual.potencia_unidad;
				}
				if (datosManual.factor_potencia !== undefined) {
					base.factor_potencia = datosManual.factor_potencia;
				}
			}

			if (datosManual.itm !== undefined) {
				base.itm = datosManual.itm;
			}
		}

		// Optional installation fields
		if (instalacion.temperatura_override !== undefined) {
			base.temperatura_override = instalacion.temperatura_override;
		}

		if (instalacion.diametro_control_mm !== undefined) {
			base.diametro_control_mm = instalacion.diametro_control_mm;
		}

		// num_tuberias: only include if it has a value and is greater than 0
		// (CHAROLA_* types don't show this field, so it will be undefined)
		if (instalacion.num_tuberias !== undefined && instalacion.num_tuberias > 0) {
			base.num_tuberias = instalacion.num_tuberias;
		}

		return base;
	});

	// ── Handlers ──────────────────────────────────────────────────────────────────
	function handleModoChange(newModo: ModoCalculo) {
		modo = newModo;
		// Reset incompatible data when switching modes
		if (newModo === 'LISTADO') {
			datosManual = {
				tipo_equipo: '',
				amperaje_nominal: undefined,
				potencia_nominal: undefined,
				potencia_unidad: 'KW',
				factor_potencia: undefined,
				itm: undefined
			};
		} else {
			equipoSeleccionado = undefined;
			// Clear installation fields that were auto-filled from equipment
			instalacion = {
				tension: undefined,
				tension_unidad: 'V',
				sistema_electrico: '',
				estado: instalacion.estado, // Keep estado as it's not equipment-specific
				tipo_canalizacion: '',
				num_tuberias: undefined,
				longitud_circuito: undefined,
				tipo_voltaje: '',
				material: 'Cu',
				hilos_por_fase: 1,
				porcentaje_caida_maximo: 3.0,
				temperatura_override: undefined
			};
		}
	}

	function handleDatosManualChange(newDatos: FormularioManualData) {
		datosManual = newDatos;
	}

	function handleEquipoChange(equipo: EquipoFiltro | undefined) {
		equipoSeleccionado = equipo;

		// Auto-fill installation fields from selected equipment
		if (equipo) {
			instalacion = {
				...instalacion,
				tension: equipo.voltaje,
				tension_unidad: 'V',
				sistema_electrico: mapearConexionASistemaElectrico(equipo.conexion),
				tipo_voltaje: mapearTipoVoltaje(equipo.tipo_voltaje)
			};
		}
	}

	function handleInstalacionChange(newDatos: CamposInstalacionData) {
		instalacion = newDatos;
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		loading = true;
		error = undefined;

		try {
			const result = await calcularMemoria(request);
			if (result.ok) {
				// Redirect to results page with data (handle Unicode characters)
				const jsonStr = JSON.stringify(result.data.data);
				const encodedData = btoa(unescape(encodeURIComponent(jsonStr)));
				goto(`/calculos/resultado?data=${encodedData}`);
			} else {
				error = result.error.error;
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Error desconocido';
		} finally {
			loading = false;
		}
	}
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
					<FormularioListado bind:equipoSeleccionado onEquipoChange={handleEquipoChange} />
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
	</div>
</div>
