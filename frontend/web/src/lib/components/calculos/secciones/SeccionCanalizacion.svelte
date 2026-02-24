<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	let canalizacion = $derived(memoria.canalizacion);
	let tipo = $derived(memoria.tipo_canalizacion);

	// Clasificación del tipo de canalización
	let esTuberia = $derived(
		tipo === 'TUBERIA_PVC' ||
			tipo === 'TUBERIA_ALUMINIO' ||
			tipo === 'TUBERIA_ACERO_PG' ||
			tipo === 'TUBERIA_ACERO_PD'
	);
	let esCharolaEspaciado = $derived(tipo === 'CHAROLA_CABLE_ESPACIADO');
	let esCharolaTriangular = $derived(tipo === 'CHAROLA_CABLE_TRIANGULAR');

	// Etiqueta legible del tipo
	let tipoLabel = $derived.by(() => {
		switch (tipo) {
			case 'TUBERIA_PVC':
				return 'Tubería PVC';
			case 'TUBERIA_ALUMINIO':
				return 'Tubería Aluminio';
			case 'TUBERIA_ACERO_PG':
				return 'Tubería Acero Pared Gruesa (Rígida)';
			case 'TUBERIA_ACERO_PD':
				return 'Tubería Acero Pared Delgada (EMT)';
			case 'CHAROLA_CABLE_ESPACIADO':
				return 'Charola de Cable — Espaciado (1 diámetro de separación)';
			case 'CHAROLA_CABLE_TRIANGULAR':
				return 'Charola de Cable — Arreglo Triangular (Tresbolillo)';
			default:
				return tipo;
		}
	});

	// Fill factor como porcentaje (solo tubería)
	let fillFactorPorcentaje = $derived((memoria.fill_factor * 100).toFixed(0));

	// Control conductor — optional
	let tieneControl = $derived(!!memoria.diametro_control_mm && memoria.diametro_control_mm > 0);

	// Detalle charola — valores intermedios para el desarrollo con números reales
	let detalle = $derived(memoria.detalle_charola);

	// Labels legibles de material y sistema eléctrico
	let materialLabel = $derived(
		memoria.conductor_alimentacion.Material?.toUpperCase() === 'CU' ? 'Cobre (Cu)' : 'Aluminio (Al)'
	);
	let sistemaLabel = $derived.by(() => {
		switch (memoria.sistema_electrico) {
			case 'DELTA':
				return 'Trifásico Delta (3F-3H)';
			case 'ESTRELLA':
				return 'Trifásico Estrella (3F-4H)';
			case 'BIFASICO':
				return 'Bifásico (2F-3H)';
			case 'MONOFASICO':
				return 'Monofásico (1F-2H)';
			default:
				return memoria.sistema_electrico;
		}
	});

	// Referencia normativa según tipo
	let referencianom = $derived.by(() => {
		if (esTuberia) return 'NOM-001-SEDE — Cap. 9, Tabla 4';
		if (esCharolaEspaciado) return 'NOM-001-SEDE — 310-15(b)(17)';
		if (esCharolaTriangular) return 'NOM-001-SEDE — 310-15(b)(20)';
		return 'NOM-001-SEDE';
	});
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		5. Cálculo de Canalización
	</h2>

	<!-- Tipo de canalización -->
	<div class="mb-4 rounded border border-primary/30 bg-primary/10 p-3">
		<p class="text-sm font-medium text-primary">
			Tipo de Canalización: {tipoLabel}
		</p>
		<p class="mt-1 text-xs text-muted-foreground">Referencia: {referencianom}</p>
	</div>

	<!-- ═══════════════════════════════════════════════════════════
	     TUBERÍA (PVC / Aluminio / Acero PG / Acero PD)
	     ═══════════════════════════════════════════════════════════ -->
	{#if esTuberia}
		<!-- Criterio de dimensionamiento NOM Cap. 9 -->
		<div class="mb-6 rounded bg-muted p-4">
			<h3 class="mb-2 font-semibold text-foreground">Criterio de Dimensionamiento — Cap. 9 NOM</h3>
			<p class="mb-2 text-sm text-muted-foreground">
				El área interior de la tubería debe alojar el área total de conductores respetando el factor
				de llenado permitido según el número de conductores:
			</p>
			<div class="mb-3 grid grid-cols-3 gap-2 text-sm">
				<div
					class="rounded border border-border p-2 text-center"
					class:border-primary={memoria.fill_factor === 0.53}
					class:bg-primary={memoria.fill_factor === 0.53}
					class:text-primary-foreground={memoria.fill_factor === 0.53}
				>
					<p class="font-mono font-bold">53%</p>
					<p class="text-xs">1 conductor</p>
				</div>
				<div
					class="rounded border border-border p-2 text-center"
					class:border-primary={memoria.fill_factor === 0.31}
					class:bg-primary={memoria.fill_factor === 0.31}
					class:text-primary-foreground={memoria.fill_factor === 0.31}
				>
					<p class="font-mono font-bold">31%</p>
					<p class="text-xs">2 conductores</p>
				</div>
				<div
					class="rounded border border-border p-2 text-center"
					class:border-primary={memoria.fill_factor === 0.4}
					class:bg-primary={memoria.fill_factor === 0.4}
					class:text-primary-foreground={memoria.fill_factor === 0.4}
				>
					<p class="font-mono font-bold">40%</p>
					<p class="text-xs">3+ conductores</p>
				</div>
			</div>
			<p class="font-mono text-sm text-foreground">
				Área_requerida = (Σ Áreas_conductores / N_tubos) / {fillFactorPorcentaje}%
			</p>
		</div>

		<!-- Conductores en la instalación -->
		<div class="mb-6">
			<h3 class="mb-2 font-semibold text-foreground">Conductores en la Instalación</h3>
			<div class="grid grid-cols-2 gap-4 text-sm">
				<div>
					<span class="text-muted-foreground">Total de conductores:</span>
					<span class="ml-2 font-medium text-foreground">{memoria.cantidad_conductores}</span>
				</div>
				<div>
					<span class="text-muted-foreground">Hilos por fase:</span>
					<span class="ml-2 font-medium text-foreground">{memoria.hilos_por_fase}</span>
				</div>
				<div>
					<span class="text-muted-foreground">Tubos en paralelo:</span>
					<span class="ml-2 font-medium text-foreground">{canalizacion?.NumeroDeTubos ?? 1}</span>
				</div>
				<div>
					<span class="text-muted-foreground">Factor de llenado aplicado:</span>
					<span class="ml-2 font-mono font-medium text-foreground">{fillFactorPorcentaje}%</span>
				</div>
			</div>
		</div>

		<!-- Resultado tubería -->
		<div>
			<h3 class="mb-2 font-semibold text-foreground">Tubería Seleccionada</h3>
			<div class="overflow-hidden rounded border border-border">
				<table class="w-full text-sm">
					<thead class="bg-muted">
						<tr>
							<th class="px-4 py-2 text-left font-medium text-foreground">Parámetro</th>
							<th class="px-4 py-2 text-left font-medium text-foreground">Valor</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-border">
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Tamaño Comercial</td>
							<td class="px-4 py-2 font-mono font-bold text-foreground">
								{canalizacion?.Tamano || '—'}
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Número de Tubos</td>
							<td class="px-4 py-2 font-medium text-foreground">
								{canalizacion?.NumeroDeTubos ?? 1}
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Área por Tubo</td>
							<td class="px-4 py-2 font-mono text-foreground">
								{canalizacion?.AreaTotalMM2?.toFixed(2) ?? '—'} mm²
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Factor de Llenado</td>
							<td class="px-4 py-2 font-mono text-foreground">{fillFactorPorcentaje}%</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<!-- Verificación -->
		<div class="mt-4 rounded border border-success/30 bg-success/10 p-3">
			<p class="text-sm text-foreground">
				✓ Tubería dimensionada conforme a NOM-001-SEDE-2012 Cap. 9, Tabla 4
			</p>
		</div>

		<!-- ═══════════════════════════════════════════════════════════
	     CHAROLA ESPACIADO — 1 diámetro de separación entre cables
	     ═══════════════════════════════════════════════════════════ -->
	{:else if esCharolaEspaciado}
		<p class="mb-4 text-sm text-muted-foreground">
			Los cables se instalan con una separación mínima igual a 1 diámetro exterior entre sí,
			conforme a NOM-001-SEDE-2012 Art. 310-15(b)(17). El ancho de charola se determina sumando los
			espacios de fuerza, control (si aplica) y tierra.
		</p>

		<!-- Norma de referencia -->
		<div class="mb-4 rounded border border-primary/30 bg-primary/10 p-3">
			<p class="text-sm font-medium text-primary">
				Referencia: NOM-001-SEDE-2012 Art. 392 / 310-15(b)(17) — Cables en charola con espaciado de
				1 diámetro
			</p>
		</div>

		<!-- Fórmula general -->
		<div class="mb-4 rounded bg-muted p-4">
			<h3 class="mb-3 font-semibold text-foreground">Fórmula de Dimensionamiento</h3>
			<div class="space-y-1 text-sm text-foreground">
				<p>
					<strong>E<sub>f</sub></strong> = Espacio de fuerza = N<sub>hilos</sub> × Ø<sub>fase</sub>
				</p>
				<p>
					<strong>A<sub>f</sub></strong> = Ancho de fuerza = N<sub>hilos</sub> × Ø<sub>fase</sub>
				</p>
				{#if tieneControl}
					<p>
						<strong>E<sub>c</sub></strong> = Espacio de control = 2 × Ø<sub>control</sub>
						<span class="text-muted-foreground">(uno a cada lado)</span>
					</p>
					<p>
						<strong>A<sub>c</sub></strong> = Ancho de control = Ø<sub>control</sub>
					</p>
				{/if}
				<p class="mt-3 border-t border-border pt-3 font-semibold">
					{#if tieneControl}
						A<sub>req</sub> = E<sub>f</sub> + A<sub>f</sub> + E<sub>c</sub> + A<sub>c</sub> + Ø<sub
							>tierra</sub
						>
					{:else}
						A<sub>req</sub> = E<sub>f</sub> + A<sub>f</sub> + Ø<sub>tierra</sub>
					{/if}
				</p>
			</div>
			<p class="mt-3 text-xs text-muted-foreground">
				DMG = 2.0 | Sin factor de agrupamiento | Tabla de ampacidad: 310-15(b)(17)
			</p>
		</div>

		<!-- Conductores en la Instalación -->
		<div class="mb-6 overflow-hidden rounded border border-border">
			<table class="w-full text-sm">
				<thead class="bg-muted">
					<tr>
						<th class="px-4 py-2 text-left font-medium text-foreground">Conductor</th>
						<th class="px-4 py-2 text-left font-medium text-foreground">Especificación</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-border">
					<!-- Sistema -->
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Sistema eléctrico</td>
						<td class="px-4 py-2 font-medium text-foreground">{sistemaLabel}</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Hilos por fase</td>
						<td class="px-4 py-2 font-medium text-foreground">{memoria.hilos_por_fase}</td>
					</tr>
					<!-- Separador visual -->
					<tr class="bg-muted/40">
						<td
							class="px-4 py-1 text-xs font-semibold tracking-wide text-muted-foreground uppercase"
							>Conductores físicos</td
						>
						<td class="px-4 py-1"></td>
					</tr>
					<!-- Fase -->
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Fase</td>
						<td class="px-4 py-2 text-foreground">
							<span class="font-mono font-medium">{memoria.conductor_alimentacion.Calibre}</span>
							<span class="ml-1 text-muted-foreground">{materialLabel}</span>
							{#if detalle}
								<span class="ml-2 text-xs text-muted-foreground"
									>Ø {detalle.diametro_fase_mm.toFixed(2)} mm</span
								>
							{/if}
						</td>
					</tr>
					<!-- Tierra -->
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Tierra</td>
						<td class="px-4 py-2 text-foreground">
							<span class="font-mono font-medium">{memoria.conductor_tierra.Calibre}</span>
							<span class="ml-1 text-muted-foreground">{materialLabel}</span>
							{#if detalle}
								<span class="ml-2 text-xs text-muted-foreground"
									>Ø {detalle.diametro_tierra_mm.toFixed(2)} mm</span
								>
							{/if}
							<span
								class="ml-2 rounded bg-muted px-1 py-0.5 text-xs font-medium text-muted-foreground"
								>Desnudo</span
							>
						</td>
					</tr>
					<!-- Control (opcional) -->
					{#if tieneControl && detalle?.diametro_control_mm}
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Control</td>
							<td class="px-4 py-2 text-foreground">
								<span class="text-xs text-muted-foreground"
									>Ø {detalle.diametro_control_mm.toFixed(2)} mm</span
								>
								<span
									class="ml-2 rounded bg-muted px-1 py-0.5 text-xs font-medium text-muted-foreground"
									>Espaciado a cada lado</span
								>
							</td>
						</tr>
					{/if}
				</tbody>
			</table>
		</div>

		<!-- Desarrollo con números reales -->
		{#if detalle}
			<div class="mb-6">
				<h3 class="mb-2 font-semibold text-foreground">Desarrollo</h3>
				<div class="space-y-1.5 rounded bg-muted/30 p-3 font-mono text-sm">
					<p class="text-foreground">
						E<sub>f</sub> = {detalle.num_hilos_total} hilos × {detalle.diametro_fase_mm.toFixed(2)} mm
						= <strong>{detalle.espacio_fuerza_mm.toFixed(2)} mm</strong>
					</p>
					<p class="text-foreground">
						A<sub>f</sub> = {detalle.num_hilos_total} hilos × {detalle.diametro_fase_mm.toFixed(2)} mm
						= <strong>{detalle.ancho_fuerza_mm?.toFixed(2) ?? '—'} mm</strong>
					</p>
					{#if tieneControl && detalle.diametro_control_mm}
						<p class="text-foreground">
							E<sub>c</sub> = 2 × {detalle.diametro_control_mm.toFixed(2)} mm =
							<strong>{detalle.espacio_control_mm?.toFixed(2) ?? '—'} mm</strong>
						</p>
						<p class="text-foreground">
							A<sub>c</sub> = <strong>{detalle.ancho_control_mm?.toFixed(2) ?? '—'} mm</strong>
						</p>
					{/if}
					<p class="text-foreground">
						Ø<sub>tierra</sub> = <strong>{detalle.ancho_tierra_mm.toFixed(2)} mm</strong>
					</p>
					<hr class="my-2 border-border/60" />
					{#if tieneControl && detalle.diametro_control_mm}
						<p class="text-muted-foreground">
							A<sub>req</sub> = {detalle.espacio_fuerza_mm.toFixed(2)} + {detalle.ancho_fuerza_mm?.toFixed(
								2
							)} + {detalle.espacio_control_mm?.toFixed(2)} + {detalle.ancho_control_mm?.toFixed(2)} +
							{detalle.ancho_tierra_mm.toFixed(2)}
						</p>
					{:else}
						<p class="text-muted-foreground">
							A<sub>req</sub> = {detalle.espacio_fuerza_mm.toFixed(2)} + {detalle.ancho_fuerza_mm?.toFixed(
								2
							)} + {detalle.ancho_tierra_mm.toFixed(2)}
						</p>
					{/if}
					<p class="text-lg font-bold text-primary">
						A<sub>req</sub> = {canalizacion?.AreaRequeridaMM2?.toFixed(2) ?? '—'} mm
					</p>
				</div>
			</div>
		{/if}

		<!-- Resultado charola espaciado -->
		<div class="mb-6">
			<h3 class="mb-2 font-semibold text-foreground">Charola Seleccionada</h3>
			<div class="overflow-hidden rounded border border-border">
				<table class="w-full text-sm">
					<thead class="bg-muted">
						<tr>
							<th class="px-4 py-2 text-left font-medium text-foreground">Parámetro</th>
							<th class="px-4 py-2 text-left font-medium text-foreground">Valor</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-border">
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Ancho Requerido (A<sub>req</sub>)</td>
							<td class="px-4 py-2 font-mono text-foreground">
								{canalizacion?.AreaRequeridaMM2?.toFixed(1) ?? '—'} mm
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Ancho Comercial Seleccionado</td>
							<td class="px-4 py-2 font-mono font-bold text-primary">
								{canalizacion?.Tamano || '—'}
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<!-- Verificación -->
		<div class="rounded border border-success/30 bg-success/10 p-3">
			<p class="text-sm text-foreground">
				✓ Charola dimensionada conforme a NOM-001-SEDE-2012 Art. 392 / 310-15(b)(17)
			</p>
		</div>

		<!-- ═══════════════════════════════════════════════════════════
	     CHAROLA TRIANGULAR — cables en tresbolillo (tocándose)
	     ═══════════════════════════════════════════════════════════ -->
	{:else if esCharolaTriangular}
		<!-- Advertencia temperatura -->
		<div class="mb-4 rounded border border-warning/40 bg-warning/10 p-3">
			<p class="text-sm font-medium text-warning-foreground">
				⚠ Temperatura mínima: 75 °C — Esta configuración no tiene columna de 60 °C en la tabla de
				ampacidad (310-15(b)(20))
			</p>
		</div>

		<p class="mb-4 text-sm text-muted-foreground">
			Los cables se instalan en disposición triangular (tresbolillo), tocándose entre sí. Se aplica
			un factor de espaciado triangular de 2.15 conforme a NOM-001-SEDE-2012 Art. 310-15(b)(20).
		</p>

		<!-- Norma de referencia -->
		<div class="mb-4 rounded border border-primary/30 bg-primary/10 p-3">
			<p class="text-sm font-medium text-primary">
				Referencia: NOM-001-SEDE-2012 Art. 392 / 310-15(b)(20) — Cables en charola, arreglo
				triangular (tresbolillo)
			</p>
		</div>

		<!-- Fórmula general -->
		<div class="mb-4 rounded bg-muted p-4">
			<h3 class="mb-3 font-semibold text-foreground">Fórmula de Dimensionamiento</h3>
			<div class="space-y-1 text-sm text-foreground">
				<p>
					<strong>A<sub>p</sub></strong> = Ancho de potencia = 2 × Ø<sub>fase</sub> × N<sub
						>hilos</sub
					>
				</p>
				<p>
					<strong>E<sub>f</sub></strong> = Espacio de fuerza = (N<sub>hilos</sub> − 1) × 2.15 × Ø<sub
						>fase</sub
					>
				</p>
				{#if tieneControl}
					<p>
						<strong>E<sub>c</sub></strong> = Espacio de control = 2.15 × Ø<sub>control</sub>
						<span class="text-muted-foreground">(a cada lado)</span>
					</p>
					<p>
						<strong>A<sub>c</sub></strong> = Ancho de control = Ø<sub>control</sub>
					</p>
				{/if}
				<p class="mt-3 border-t border-border pt-3 font-semibold">
					{#if tieneControl}
						A<sub>req</sub> = A<sub>p</sub> + E<sub>f</sub> + E<sub>c</sub> + A<sub>c</sub> + Ø<sub
							>tierra</sub
						>
					{:else}
						A<sub>req</sub> = A<sub>p</sub> + E<sub>f</sub> + Ø<sub>tierra</sub>
					{/if}
				</p>
			</div>
			<p class="mt-3 text-xs text-muted-foreground">
				Factor triangular = 2.15 | DMG = 1.0 | Sin factor de agrupamiento | Tabla de ampacidad:
				310-15(b)(20)
			</p>
		</div>

		<!-- Conductores en la Instalación -->
		<div class="mb-6 overflow-hidden rounded border border-border">
			<table class="w-full text-sm">
				<thead class="bg-muted">
					<tr>
						<th class="px-4 py-2 text-left font-medium text-foreground">Conductor</th>
						<th class="px-4 py-2 text-left font-medium text-foreground">Especificación</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-border">
					<!-- Sistema -->
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Sistema eléctrico</td>
						<td class="px-4 py-2 font-medium text-foreground">{sistemaLabel}</td>
					</tr>
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Hilos por fase</td>
						<td class="px-4 py-2 font-medium text-foreground">{memoria.hilos_por_fase}</td>
					</tr>
					<!-- Separador visual -->
					<tr class="bg-muted/40">
						<td
							class="px-4 py-1 text-xs font-semibold tracking-wide text-muted-foreground uppercase"
							>Conductores físicos</td
						>
						<td class="px-4 py-1"></td>
					</tr>
					<!-- Fase -->
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Fase</td>
						<td class="px-4 py-2 text-foreground">
							<span class="font-mono font-medium">{memoria.conductor_alimentacion.Calibre}</span>
							<span class="ml-1 text-muted-foreground">{materialLabel}</span>
							{#if detalle}
								<span class="ml-2 text-xs text-muted-foreground"
									>Ø {detalle.diametro_fase_mm.toFixed(2)} mm</span
								>
							{/if}
						</td>
					</tr>
					<!-- Tierra -->
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Tierra</td>
						<td class="px-4 py-2 text-foreground">
							<span class="font-mono font-medium">{memoria.conductor_tierra.Calibre}</span>
							<span class="ml-1 text-muted-foreground">{materialLabel}</span>
							{#if detalle}
								<span class="ml-2 text-xs text-muted-foreground"
									>Ø {detalle.diametro_tierra_mm.toFixed(2)} mm</span
								>
							{/if}
							<span
								class="ml-2 rounded bg-muted px-1 py-0.5 text-xs font-medium text-muted-foreground"
								>Desnudo</span
							>
						</td>
					</tr>
					<!-- Control (opcional) -->
					{#if tieneControl && detalle?.diametro_control_mm}
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Control</td>
							<td class="px-4 py-2 text-foreground">
								<span class="text-xs text-muted-foreground"
									>Ø {detalle.diametro_control_mm.toFixed(2)} mm</span
								>
								<span
									class="ml-2 rounded bg-muted px-1 py-0.5 text-xs font-medium text-muted-foreground"
									>Factor 2.15 a cada lado</span
								>
							</td>
						</tr>
					{/if}
				</tbody>
			</table>
		</div>

		<!-- Desarrollo con números reales -->
		{#if detalle}
			<div class="mb-6">
				<h3 class="mb-2 font-semibold text-foreground">Desarrollo</h3>
				<div class="space-y-1.5 rounded bg-muted/30 p-3 font-mono text-sm">
					<p class="text-foreground">
						A<sub>p</sub> = 2 × {detalle.diametro_fase_mm.toFixed(2)} mm × {memoria.hilos_por_fase}
						hilos = <strong>{detalle.ancho_potencia_mm?.toFixed(2) ?? '—'} mm</strong>
					</p>
					<p class="text-foreground">
						E<sub>f</sub> = ({memoria.hilos_por_fase} − 1) × {detalle.factor_triangular?.toFixed(2)} ×
						{detalle.diametro_fase_mm.toFixed(2)} mm =
						<strong>{detalle.espacio_fuerza_mm.toFixed(2)} mm</strong>
					</p>
					{#if tieneControl && detalle.diametro_control_mm}
						<p class="text-foreground">
							E<sub>c</sub> = {detalle.factor_triangular?.toFixed(2)} × {detalle.diametro_control_mm.toFixed(
								2
							)} mm = <strong>{detalle.espacio_control_mm?.toFixed(2) ?? '—'} mm</strong>
						</p>
						<p class="text-foreground">
							A<sub>c</sub> = <strong>{detalle.ancho_control_mm?.toFixed(2) ?? '—'} mm</strong>
						</p>
					{/if}
					<p class="text-foreground">
						Ø<sub>tierra</sub> = <strong>{detalle.ancho_tierra_mm.toFixed(2)} mm</strong>
					</p>
					<hr class="my-2 border-border/60" />
					{#if tieneControl && detalle.diametro_control_mm}
						<p class="text-muted-foreground">
							A<sub>req</sub> = {detalle.ancho_potencia_mm?.toFixed(2)} + {detalle.espacio_fuerza_mm.toFixed(
								2
							)} + {detalle.espacio_control_mm?.toFixed(2)} + {detalle.ancho_control_mm?.toFixed(2)} +
							{detalle.ancho_tierra_mm.toFixed(2)}
						</p>
					{:else}
						<p class="text-muted-foreground">
							A<sub>req</sub> = {detalle.ancho_potencia_mm?.toFixed(2)} + {detalle.espacio_fuerza_mm.toFixed(
								2
							)} + {detalle.ancho_tierra_mm.toFixed(2)}
						</p>
					{/if}
					<p class="text-lg font-bold text-primary">
						A<sub>req</sub> = {canalizacion?.AreaRequeridaMM2?.toFixed(2) ?? '—'} mm
					</p>
				</div>
			</div>
		{/if}

		<!-- Resultado charola triangular -->
		<div class="mb-6">
			<h3 class="mb-2 font-semibold text-foreground">Charola Seleccionada</h3>
			<div class="overflow-hidden rounded border border-border">
				<table class="w-full text-sm">
					<thead class="bg-muted">
						<tr>
							<th class="px-4 py-2 text-left font-medium text-foreground">Parámetro</th>
							<th class="px-4 py-2 text-left font-medium text-foreground">Valor</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-border">
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Factor Triangular (NOM)</td>
							<td class="px-4 py-2 font-mono text-foreground">2.15</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Ancho Requerido (A<sub>req</sub>)</td>
							<td class="px-4 py-2 font-mono text-foreground">
								{canalizacion?.AreaRequeridaMM2?.toFixed(1) ?? '—'} mm
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Ancho Comercial Seleccionado</td>
							<td class="px-4 py-2 font-mono font-bold text-primary">
								{canalizacion?.Tamano || '—'}
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<!-- Verificación -->
		<div class="rounded border border-success/30 bg-success/10 p-3">
			<p class="text-sm text-foreground">
				✓ Charola dimensionada conforme a NOM-001-SEDE-2012 Art. 392 / 310-15(b)(20)
			</p>
		</div>
	{:else}
		<!-- Fallback para tipos no reconocidos -->
		<p class="text-sm text-muted-foreground">
			Tipo de canalización no reconocido: <span class="font-mono">{tipo}</span>
		</p>
	{/if}
</section>
