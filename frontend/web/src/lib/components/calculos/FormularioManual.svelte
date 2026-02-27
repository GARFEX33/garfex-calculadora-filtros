<script lang="ts">
	import type {
		TipoEquipo,
		UnidadPotencia,
		TipoVoltaje,
		SistemaElectrico
	} from '$lib/types/calculos.types';

	export interface FormularioManualData {
		tipo_equipo: TipoEquipo | '';
		amperaje_nominal: number | undefined;
		potencia_nominal: number | undefined;
		potencia_unidad: UnidadPotencia;
		factor_potencia: number | undefined;
		itm: number | undefined;
		tension: number;
		tipo_voltaje: TipoVoltaje;
		sistema_electrico: SistemaElectrico | '';
	}

	interface Props {
		modo: 'MANUAL_AMPERAJE' | 'MANUAL_POTENCIA';
		datos: FormularioManualData;
		onDatosChange: (datos: FormularioManualData) => void;
	}

	// Ajustar el tipo para que acepte los valores correctos
	let { modo, datos = $bindable(), onDatosChange }: Props = $props();

	const tipoEquipoOptions: { value: TipoEquipo; label: string }[] = [
		{ value: 'FILTRO_ACTIVO', label: 'Filtro Activo' },
		{ value: 'TRANSFORMADOR', label: 'Transformador' },
		{ value: 'FILTRO_RECHAZO', label: 'Filtro de Rechazo' },
		{ value: 'CARGA', label: 'Carga General' }
	];

	const potenciaUnidadOptions: { value: UnidadPotencia; label: string }[] = [
		{ value: 'W', label: 'W' },
		{ value: 'KW', label: 'kW' },
		{ value: 'KVA', label: 'kVA' },
		{ value: 'KVAR', label: 'kVAR' }
	];

	const tipoVoltajeOptions: { value: TipoVoltaje; label: string }[] = [
		{ value: 'FASE_NEUTRO', label: 'Fase-Neutro (FN)' },
		{ value: 'FASE_FASE', label: 'Fase-Fase (FF)' }
	];

	const sistemaElectricoOptions: { value: SistemaElectrico; label: string }[] = [
		{ value: 'ESTRELLA', label: 'Trifásico Estrella (4 hilos)' },
		{ value: 'DELTA', label: 'Trifásico Delta (3 hilos)' },
		{ value: 'MONOFASICO', label: 'Monofásico (2 hilos)' },
		{ value: 'BIFASICO', label: 'Bifásico (3 hilos)' }
	];

	// Using derived to conditionally show fields
	let mostrarAmperaje = $derived(modo === 'MANUAL_AMPERAJE');
	let mostrarPotencia = $derived(modo === 'MANUAL_POTENCIA');
	let mostrarFactorPotencia = $derived(modo === 'MANUAL_POTENCIA' && datos.tipo_equipo === 'CARGA');

	function updateDatos<K extends keyof FormularioManualData>(
		key: K,
		value: FormularioManualData[K]
	) {
		onDatosChange({ ...datos, [key]: value });
	}
</script>

<div class="flex flex-col gap-6">
	<!-- Tipo de Equipo -->
	<div class="flex flex-col gap-1.5">
		<label for="tipo_equipo" class="text-sm font-medium text-foreground">Tipo de Equipo</label>
		<select
			id="tipo_equipo"
			value={datos.tipo_equipo}
			onchange={(e) => updateDatos('tipo_equipo', e.currentTarget.value as TipoEquipo | '')}
			class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground focus:ring-2 focus:ring-ring focus:outline-none"
		>
			<option value="">Seleccionar...</option>
			{#each tipoEquipoOptions as opt}
				<option value={opt.value}>{opt.label}</option>
			{/each}
		</select>
	</div>

	<!-- Amperaje Nominal (solo modo MANUAL_AMPERAJE) -->
	{#if mostrarAmperaje}
		<div class="flex flex-col gap-1.5">
			<label for="amperaje_nominal" class="text-sm font-medium text-foreground"
				>Amperaje Nominal (A)</label
			>
			<input
				type="number"
				id="amperaje_nominal"
				placeholder="30"
				value={datos.amperaje_nominal ?? ''}
				oninput={(e) =>
					updateDatos(
						'amperaje_nominal',
						e.currentTarget.value ? Number(e.currentTarget.value) : undefined
					)}
				class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:outline-none"
			/>
		</div>
	{/if}

	<!-- Potencia Nominal (solo modo MANUAL_POTENCIA) -->
	{#if mostrarPotencia}
		<div class="flex flex-col gap-1.5">
			<label for="potencia_nominal" class="text-sm font-medium text-foreground"
				>Potencia Nominal</label
			>
			<div class="flex gap-2">
				<input
					type="number"
					id="potencia_nominal"
					placeholder="5"
					value={datos.potencia_nominal ?? ''}
					oninput={(e) =>
						updateDatos(
							'potencia_nominal',
							e.currentTarget.value ? Number(e.currentTarget.value) : undefined
						)}
					class="flex-1 rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:outline-none"
				/>
				<select
					value={datos.potencia_unidad}
					onchange={(e) => updateDatos('potencia_unidad', e.currentTarget.value as UnidadPotencia)}
					class="w-24 rounded-md border border-input-border bg-input px-2 py-2 text-sm text-foreground focus:ring-2 focus:ring-ring focus:outline-none"
				>
					{#each potenciaUnidadOptions as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
			</div>
		</div>
	{/if}

	<!-- Tensión ( siempre visible en modos manuales) -->
	<div class="flex flex-col gap-1.5">
		<label for="tension" class="text-sm font-medium text-foreground">Tensión (V)</label>
		<input
			type="number"
			id="tension"
			required
			placeholder="220"
			value={datos.tension ?? ''}
			oninput={(e) =>
				updateDatos('tension', e.currentTarget.value ? Number(e.currentTarget.value) : 0)}
			class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:outline-none"
		/>
	</div>

	<!-- Tipo de Voltaje (siempre visible en modos manuales) -->
	<div class="flex flex-col gap-1.5">
		<label for="tipo_voltaje" class="text-sm font-medium text-foreground">Tipo de Voltaje</label>
		<select
			id="tipo_voltaje"
			value={datos.tipo_voltaje}
			onchange={(e) => updateDatos('tipo_voltaje', e.currentTarget.value as TipoVoltaje)}
			class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground focus:ring-2 focus:ring-ring focus:outline-none"
		>
			{#each tipoVoltajeOptions as opt}
				<option value={opt.value}>{opt.label}</option>
			{/each}
		</select>
	</div>

	<!-- Sistema Eléctrico (siempre visible en modos manuales) -->
	<div class="flex flex-col gap-1.5">
		<label for="sistema_electrico" class="text-sm font-medium text-foreground"
			>Sistema Eléctrico</label
		>
		<select
			id="sistema_electrico"
			value={datos.sistema_electrico}
			onchange={(e) =>
				updateDatos('sistema_electrico', e.currentTarget.value as SistemaElectrico | '')}
			class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground focus:ring-2 focus:ring-ring focus:outline-none"
		>
			<option value="">Seleccionar...</option>
			{#each sistemaElectricoOptions as opt}
				<option value={opt.value}>{opt.label}</option>
			{/each}
		</select>
	</div>

	<!-- Factor de Potencia (solo modo MANUAL_POTENCIA + tipo CARGA) -->
	{#if mostrarFactorPotencia}
		<div class="flex flex-col gap-1.5">
			<label for="factor_potencia" class="text-sm font-medium text-foreground"
				>Factor de Potencia (%)</label
			>
			<input
				type="number"
				id="factor_potencia"
				min="1"
				max="100"
				placeholder="98"
				value={datos.factor_potencia ?? ''}
				oninput={(e) =>
					updateDatos(
						'factor_potencia',
						e.currentTarget.value ? Number(e.currentTarget.value) : undefined
					)}
				class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:outline-none"
			/>
			<p class="text-xs text-muted-foreground">Ingrese valor entre 1 y 100 (ej: 98 = 98%)</p>
		</div>
	{/if}

	<!-- ITM (Interruptor Termomagnético) -->
	<div class="flex flex-col gap-1.5">
		<label for="itm" class="text-sm font-medium text-foreground"
			>Interruptor Termomagnético (A)</label
		>
		<input
			type="number"
			id="itm"
			placeholder="30"
			value={datos.itm ?? ''}
			oninput={(e) =>
				updateDatos('itm', e.currentTarget.value ? Number(e.currentTarget.value) : undefined)}
			class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:outline-none"
		/>
	</div>
</div>
