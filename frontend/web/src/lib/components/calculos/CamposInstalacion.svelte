<script lang="ts">
	import { cn } from '$lib/utils';
	import type { TipoCanalizacion, MaterialConductor } from '$lib/types/calculos.types';

	export interface CamposInstalacionData {
		tension: number | undefined;
		tension_unidad: 'V' | 'kV';
		sistema_electrico: string;
		estado: string;
		tipo_canalizacion: TipoCanalizacion | '';
		longitud_circuito: number | undefined;
		tipo_voltaje: string;
		material: MaterialConductor;
		hilos_por_fase: number;
		porcentaje_caida_maximo: number;
		temperatura_override: number | undefined;
	}

	interface Props {
		datos: CamposInstalacionData;
		onDatosChange: (datos: CamposInstalacionData) => void;
	}

	let { datos = $bindable(), onDatosChange }: Props = $props();

	let mostrarAvanzadas = $state(false);

	const estadosMexico = [
		'Aguascalientes',
		'Baja California',
		'Baja California Sur',
		'Campeche',
		'Chiapas',
		'Chihuahua',
		'Ciudad de Mexico',
		'Coahuila',
		'Colima',
		'Durango',
		'Guanajuato',
		'Guerrero',
		'Hidalgo',
		'Jalisco',
		'Estado de Mexico',
		'Michoacan',
		'Morelos',
		'Nayarit',
		'Nuevo Leon',
		'Oaxaca',
		'Puebla',
		'Queretaro',
		'Quintana Roo',
		'San Luis Potosi',
		'Sinaloa',
		'Sonora',
		'Tabasco',
		'Tamaulipas',
		'Tlaxcala',
		'Veracruz',
		'Yucatan',
		'Zacatecas'
	];

	const canalizacionOptions: { value: TipoCanalizacion; label: string }[] = [
		{ value: 'TUBERIA_PVC', label: 'Tubería PVC' },
		{ value: 'TUBERIA_EMT', label: 'Tubería EMT' },
		{ value: 'CHAROLA_CABLE_ESPACIADO', label: 'Charola (Espaciado)' },
		{ value: 'CHAROLA_CABLE_TRESBOLILLO', label: 'Charola (Tresbolillo)' }
	];

	function updateDatos<K extends keyof CamposInstalacionData>(
		key: K,
		value: CamposInstalacionData[K]
	) {
		onDatosChange({ ...datos, [key]: value });
	}
</script>

<div class="flex flex-col gap-6">
	<!-- Estado - siempre visible -->
	<div class="space-y-4">
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<!-- Estado -->
			<div class="flex flex-col gap-1.5">
				<label for="estado" class="text-sm text-muted-foreground">Estado</label>
				<select
					id="estado"
					value={datos.estado}
					onchange={(e) => updateDatos('estado', e.currentTarget.value)}
					class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground focus:ring-2 focus:ring-ring focus:outline-none"
				>
					<option value="">Seleccionar...</option>
					{#each estadosMexico as estado}
						<option value={estado}>{estado}</option>
					{/each}
				</select>
			</div>
		</div>
	</div>

	<!-- Sección: Equipo Eléctrico - ELIMINADA PERMANENTEMENTE
		 Los campos de tensión, sistema eléctrico y tipo de voltaje ya no se muestran
		 en la página principal según requerimiento del usuario.
		 Los valores se manejan internamente desde el equipo seleccionado.
	 -->

	<!-- Sección: Canalización -->
	<div class="space-y-4">
		<h3 class="text-sm font-semibold text-foreground">Canalización</h3>
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<!-- Tipo Canalización -->
			<div class="flex flex-col gap-1.5">
				<label for="tipo_canalizacion" class="text-sm text-muted-foreground"
					>Tipo de Canalización</label
				>
				<select
					id="tipo_canalizacion"
					value={datos.tipo_canalizacion}
					onchange={(e) =>
						updateDatos('tipo_canalizacion', e.currentTarget.value as TipoCanalizacion | '')}
					class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground focus:ring-2 focus:ring-ring focus:outline-none"
				>
					<option value="">Seleccionar...</option>
					{#each canalizacionOptions as opt}
						<option value={opt.value}>{opt.label}</option>
					{/each}
				</select>
			</div>

			<!-- Longitud del Circuito -->
			<div class="flex flex-col gap-1.5">
				<label for="longitud_circuito" class="text-sm text-muted-foreground"
					>Longitud del Circuito (m)</label
				>
				<input
					type="number"
					id="longitud_circuito"
					placeholder="30"
					value={datos.longitud_circuito ?? ''}
					oninput={(e) =>
						updateDatos(
							'longitud_circuito',
							e.currentTarget.value ? Number(e.currentTarget.value) : undefined
						)}
					class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground focus:ring-2 focus:ring-ring focus:outline-none"
				/>
			</div>
		</div>
	</div>

	<!-- Sección: Opciones Avanzadas (collapsible) -->
	<div class="space-y-4">
		<button
			type="button"
			class="flex items-center gap-2 text-sm font-semibold text-foreground transition-colors hover:text-primary"
			onclick={() => (mostrarAvanzadas = !mostrarAvanzadas)}
		>
			<svg
				class={cn('h-4 w-4 transition-transform', mostrarAvanzadas && 'rotate-90')}
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
			>
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
			</svg>
			Opciones Avanzadas
		</button>

		{#if mostrarAvanzadas}
			<div class="grid grid-cols-1 gap-4 border-l-2 border-muted pl-6 md:grid-cols-3">
				<!-- Hilos por Fase -->
				<div class="flex flex-col gap-1.5">
					<label for="hilos_por_fase" class="text-sm text-muted-foreground">Hilos por Fase</label>
					<input
						type="number"
						id="hilos_por_fase"
						value={datos.hilos_por_fase}
						oninput={(e) => updateDatos('hilos_por_fase', Number(e.currentTarget.value))}
						class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground focus:ring-2 focus:ring-ring focus:outline-none"
					/>
				</div>

				<!-- % Caída Máximo -->
				<div class="flex flex-col gap-1.5">
					<label for="porcentaje_caida_maximo" class="text-sm text-muted-foreground"
						>% Caída Máximo</label
					>
					<input
						type="number"
						id="porcentaje_caida_maximo"
						value={datos.porcentaje_caida_maximo}
						oninput={(e) => updateDatos('porcentaje_caida_maximo', Number(e.currentTarget.value))}
						class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground focus:ring-2 focus:ring-ring focus:outline-none"
					/>
				</div>

				<!-- Temperatura Override -->
				<div class="flex flex-col gap-1.5">
					<label for="temperatura_override" class="text-sm text-muted-foreground"
						>Temperatura (°C) - Opcional</label
					>
					<input
						type="number"
						id="temperatura_override"
						placeholder="Auto"
						value={datos.temperatura_override ?? ''}
						oninput={(e) =>
							updateDatos(
								'temperatura_override',
								e.currentTarget.value ? Number(e.currentTarget.value) : undefined
							)}
						class="w-full rounded-md border border-input-border bg-input px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:ring-2 focus:ring-ring focus:outline-none"
					/>
				</div>
			</div>
		{/if}
	</div>
</div>
