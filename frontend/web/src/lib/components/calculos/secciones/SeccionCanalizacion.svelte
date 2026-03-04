<script lang="ts">
	import type { MemoriaOutput } from '$lib/types/calculos.types';
	import DiagramaCable from './DiagramaCable.svelte';

	interface Props {
		memoria: MemoriaOutput;
	}

	let { memoria }: Props = $props();

	let canalizacion = $derived(memoria.canalizacion);
	let tipo = $derived(memoria.instalacion.tipo_canalizacion);

	// Clasificación del tipo de canalización
	let esTuberia = $derived(
		tipo === 'TUBERIA_PVC' ||
			tipo === 'TUBERIA_ALUMINIO' ||
			tipo === 'TUBERIA_ACERO_PG' ||
			tipo === 'TUBERIA_ACERO_PD'
	);

	let esCharolaEspaciado = $derived(tipo === 'CHAROLA_CABLE_ESPACIADO');
	let esCharolaTriangular = $derived(tipo === 'CHAROLA_CABLE_TRIANGULAR');

	// Fill factor como porcentaje (solo tubería)
	let fillFactorPorcentaje = $derived((memoria.canalizacion.fill_factor * 100).toFixed(0));

	// Control conductor — optional
	let tieneControl = $derived(
		!!memoria.canalizacion.detalle_charola?.diametro_control_mm &&
			memoria.canalizacion.detalle_charola.diametro_control_mm > 0
	);

	// Detalle charola — valores intermedios para el desarrollo con números reales
	let detalle = $derived(memoria.canalizacion.detalle_charola);

	// Detalle tubería — valores intermedios para el desarrollo con números reales
	let detalleTuberia = $derived(memoria.canalizacion.detalle_tuberia);

	// Factor de control — usa el valor del API o default 1.0
	let factorControl = $derived(detalle?.factor_control ?? 1);

	// Props para DiagramaCable — solo incluir diagrama si existe valor
	let propsDiagrama = $derived(
		memoria.canalizacion.detalle_charola?.diagrama
			? { diagrama: memoria.canalizacion.detalle_charola.diagrama }
			: memoria.canalizacion.detalle_tuberia?.diagrama
				? { diagrama: memoria.canalizacion.detalle_tuberia.diagrama }
				: {}
	);

	// Labels legibles de material y sistema eléctrico
	let materialLabel = $derived(
		memoria.cable_fase.material?.toUpperCase() === 'CU' ? 'Cobre (Cu)' : 'Aluminio (Al)'
	);
	let materialTierraLabel = $derived(
		memoria.cable_tierra.material?.toUpperCase() === 'CU' ? 'Cobre (Cu)' : 'Aluminio (Al)'
	);
	let sistemaLabel = $derived.by(() => {
		switch (memoria.instalacion.sistema_electrico) {
			case 'DELTA':
				return 'Trifásico Delta (3F-3H)';
			case 'ESTRELLA':
				return 'Trifásico Estrella (3F-4H)';
			case 'BIFASICO':
				return 'Bifásico (2F-3H)';
			case 'MONOFASICO':
				return 'Monofásico (1F-2H)';
			default:
				return memoria.instalacion.sistema_electrico;
		}
	});
</script>

<section class="rounded-lg border border-border bg-card p-6">
	<h2 class="mb-4 border-b border-border pb-2 text-xl font-semibold text-card-foreground">
		5. Cálculo de Canalización
	</h2>

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
					class:border-primary={memoria.canalizacion.fill_factor === 0.53}
					class:bg-primary={memoria.canalizacion.fill_factor === 0.53}
					class:text-primary-foreground={memoria.canalizacion.fill_factor === 0.53}
				>
					<p class="font-mono font-bold">53%</p>
					<p class="text-xs">1 conductor</p>
				</div>
				<div
					class="rounded border border-border p-2 text-center"
					class:border-primary={memoria.canalizacion.fill_factor === 0.31}
					class:bg-primary={memoria.canalizacion.fill_factor === 0.31}
					class:text-primary-foreground={memoria.canalizacion.fill_factor === 0.31}
				>
					<p class="font-mono font-bold">31%</p>
					<p class="text-xs">2 conductores</p>
				</div>
				<div
					class="rounded border border-border p-2 text-center"
					class:border-primary={memoria.canalizacion.fill_factor === 0.4}
					class:bg-primary={memoria.canalizacion.fill_factor === 0.4}
					class:text-primary-foreground={memoria.canalizacion.fill_factor === 0.4}
				>
					<p class="font-mono font-bold">40%</p>
					<p class="text-xs">3+ conductores</p>
				</div>
			</div>
			{#if detalleTuberia}
				<hr class="my-3 border-border/60" />
				<!-- Fórmula simbólica -->
				<p class="font-mono text-sm text-muted-foreground">
					{#if canalizacion.resultado.numero_de_tubos > 1}
						{#if detalleTuberia.area_neutro_mm2}
							A<sub>req</sub> = (N<sub>fases</sub>/N<sub>tubos</sub>) × A<sub>fase</sub> + (N<sub
								>neutros</sub
							>/N<sub>tubos</sub>) × A<sub>neutro</sub> + (N<sub>tierra</sub>/N<sub>tubos</sub>) × A<sub
								>tierra</sub
							>
						{:else}
							A<sub>req</sub> = (N<sub>fases</sub>/N<sub>tubos</sub>) × A<sub>fase</sub> + (N<sub
								>tierra</sub
							>/N<sub>tubos</sub>) × A<sub>tierra</sub>
						{/if}
					{:else if detalleTuberia.area_neutro_mm2}
						A<sub>req</sub> = N<sub>fases</sub> × A<sub>fase</sub> + N<sub>neutros</sub> × A<sub
							>neutro</sub
						>
						+ N<sub>tierra</sub> × A<sub>tierra</sub>
					{:else}
						A<sub>req</sub> = N<sub>fases</sub> × A<sub>fase</sub> + N<sub>tierra</sub> × A<sub
							>tierra</sub
						>
					{/if}
				</p>
			{/if}
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
						<td class="px-4 py-2 font-medium text-foreground"
							>{memoria.instalacion.hilos_por_fase}</td
						>
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
							<span class="font-mono font-medium">{memoria.cable_fase.calibre}</span>
							<span class="ml-1 text-muted-foreground">{materialLabel}</span>
						</td>
					</tr>
					<!-- Tierra -->
					<tr>
						<td class="px-4 py-2 text-muted-foreground">Tierra</td>
						<td class="px-4 py-2 text-foreground">
							<span class="font-mono font-medium">{memoria.cable_tierra.calibre}</span>
							<span class="ml-1 text-muted-foreground">{materialTierraLabel}</span>
							<span
								class="ml-2 rounded bg-muted px-1 py-0.5 text-xs font-medium text-muted-foreground"
								>Desnudo</span
							>
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Desarrollo — con detalle de áreas por conductor -->
		{#if detalleTuberia}
			<div class="mb-6">
				<h3 class="mb-2 font-semibold text-foreground">Desarrollo</h3>
				<div class="space-y-1.5 rounded bg-muted/30 p-3 font-mono text-sm">
					<!-- Sistema -->
					<p class="text-muted-foreground">
						{sistemaLabel} — {memoria.instalacion.hilos_por_fase} conductor(es) por fase — {canalizacion
							.resultado.numero_de_tubos}
						tubo(s)
					</p>

					<hr class="my-2 border-border/60" />

					<!-- Fases -->
					<p class="text-foreground">
						<strong>Fase</strong>: {detalleTuberia.num_fases_por_tubo} × {detalleTuberia.area_fase_mm2.toFixed(
							2
						)} mm²
						<span class="font-sans text-xs text-muted-foreground">
							({memoria.cable_fase.calibre} — Tabla 5 NOM)
						</span>
						=
						<strong
							>{(detalleTuberia.num_fases_por_tubo * detalleTuberia.area_fase_mm2).toFixed(2)} mm²</strong
						>
					</p>

					<!-- Neutro (si aplica) -->
					{#if detalleTuberia.area_neutro_mm2}
						<p class="text-foreground">
							<strong>Neutro</strong>: {detalleTuberia.num_neutros_por_tubo} × {detalleTuberia.area_neutro_mm2.toFixed(
								2
							)} mm²
							<span class="font-sans text-xs text-muted-foreground">
								({memoria.cable_fase.calibre} — Tabla 5 NOM)
							</span>
							=
							<strong
								>{(detalleTuberia.num_neutros_por_tubo * detalleTuberia.area_neutro_mm2).toFixed(2)} mm²</strong
							>
						</p>
					{/if}

					<!-- Tierra -->
					<p class="text-foreground">
						<strong>Tierra</strong>: {detalleTuberia.num_tierras} × {detalleTuberia.area_tierra_mm2.toFixed(
							2
						)} mm²
						<span class="font-sans text-xs text-muted-foreground">
							({memoria.cable_tierra.calibre} Desnudo — Tabla 8 NOM)
						</span>
						=
						<strong
							>{(detalleTuberia.num_tierras * detalleTuberia.area_tierra_mm2).toFixed(2)} mm²</strong
						>
					</p>

					<hr class="my-2 border-border/60" />

					<!-- Selección en tabla NOM -->
					<p class="text-muted-foreground">
						Tabla NOM Cap. 9 — seleccionar primer tubo donde Área<sub
							>ocup. {fillFactorPorcentaje}%</sub
						>
						≥ {canalizacion.resultado.area_total_mm2.toFixed(2)} mm²:
					</p>
					<p class="text-lg font-bold text-primary">
						Tubo: {canalizacion.resultado.tamano}" / {detalleTuberia.designacion_metrica} mm — Área<sub
							>ocup.</sub
						>
						= {detalleTuberia.area_ocupacion_tubo_mm2.toFixed(0)} mm²
					</p>
				</div>
			</div>
		{/if}

		<!-- Diagrama SVG de arreglo de cables -->
		<div class="mb-6">
			<DiagramaCable
				tipoCanalizacion={memoria.instalacion.tipo_canalizacion}
				detalleCharola={memoria.canalizacion.detalle_charola}
				detalleTuberia={memoria.canalizacion.detalle_tuberia}
				resultado={memoria.canalizacion.resultado}
				sistemaElectrico={memoria.instalacion.sistema_electrico}
				hilosPorFase={memoria.instalacion.hilos_por_fase}
				calibreFase={memoria.cable_fase.calibre}
				calibreTierra={memoria.cable_tierra.calibre}
				{...propsDiagrama}
			/>
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
								{canalizacion.resultado?.tamano || '—'}
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Designación Métrica</td>
							<td class="px-4 py-2 font-mono text-foreground">
								{detalleTuberia?.designacion_metrica
									? `${detalleTuberia.designacion_metrica} mm`
									: '—'}
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Número de Tubos</td>
							<td class="px-4 py-2 font-medium text-foreground">
								{canalizacion.resultado?.numero_de_tubos ?? 1}
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Área Total de Cables</td>
							<td class="px-4 py-2 font-mono text-foreground">
								{canalizacion.resultado?.area_total_mm2?.toFixed(2) ?? '—'} mm²
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Área de Ocupación al 40% (NOM)</td>
							<td class="px-4 py-2 font-mono text-foreground">
								{detalleTuberia?.area_ocupacion_tubo_mm2?.toFixed(0) ?? '—'} mm²
							</td>
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
						<strong>E<sub>c</sub></strong> = Espacio de control = f<sub>c</sub> × Ø<sub>control</sub
						>
						<span class="text-muted-foreground">(factor de control)</span>
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
						<td class="px-4 py-2 font-medium text-foreground"
							>{memoria.instalacion.hilos_por_fase}</td
						>
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
							<span class="font-mono font-medium">{memoria.cable_fase.calibre}</span>
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
							<span class="font-mono font-medium">{memoria.cable_tierra.calibre}</span>
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
							E<sub>c</sub> = f<sub>c</sub> × {detalle.diametro_control_mm.toFixed(2)} mm =
							<strong>{detalle.espacio_control_mm?.toFixed(2) ?? '—'} mm</strong>
							<span class="font-sans text-xs text-muted-foreground">
								(f_c = {factorControl.toFixed(2)})</span
							>
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
						A<sub>req</sub> = {canalizacion.resultado?.area_requerida_mm2?.toFixed(2) ?? '—'} mm
					</p>
				</div>
			</div>
		{/if}

		<!-- Diagrama SVG de arreglo de cables -->
		<div class="mb-6">
			<DiagramaCable
				tipoCanalizacion={memoria.instalacion.tipo_canalizacion}
				detalleCharola={memoria.canalizacion.detalle_charola}
				detalleTuberia={memoria.canalizacion.detalle_tuberia}
				resultado={memoria.canalizacion.resultado}
				sistemaElectrico={memoria.instalacion.sistema_electrico}
				hilosPorFase={memoria.instalacion.hilos_por_fase}
				calibreFase={memoria.cable_fase.calibre}
				calibreTierra={memoria.cable_tierra.calibre}
				{...propsDiagrama}
			/>
		</div>

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
								{canalizacion.resultado?.area_requerida_mm2?.toFixed(1) ?? '—'} mm
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Ancho Comercial Seleccionado</td>
							<td class="px-4 py-2 font-mono font-bold text-primary">
								{canalizacion.resultado?.tamano || '—'}
							</td>
						</tr>
						{#if canalizacion.resultado?.ancho_comercial_mm && canalizacion.resultado.ancho_comercial_mm > 0}
							<tr>
								<td class="px-4 py-2 text-muted-foreground">Ancho Comercial (mm)</td>
								<td class="px-4 py-2 font-mono font-bold text-primary">
									{canalizacion.resultado.ancho_comercial_mm.toFixed(1)} mm
								</td>
							</tr>
						{/if}
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
	     CHAROLA TRIANGULAR
	     ═══════════════════════════════════════════════════════════ -->
	{:else if esCharolaTriangular}
		<p class="mb-4 text-sm text-muted-foreground">
			Los cables se instalan en disposición triangular, tocándose entre sí. Se aplica un factor de
			espaciado triangular de 2.15 conforme a NOM-001-SEDE-2012 Art. 310-15(b)(20).
		</p>

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
						<strong>E<sub>c</sub></strong> = Espacio de control = f<sub>c</sub> × Ø<sub>control</sub
						>
						<span class="text-muted-foreground">(factor de control)</span>
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
						<td class="px-4 py-2 font-medium text-foreground"
							>{memoria.instalacion.hilos_por_fase}</td
						>
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
							<span class="font-mono font-medium">{memoria.cable_fase.calibre}</span>
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
							<span class="font-mono font-medium">{memoria.cable_tierra.calibre}</span>
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
									>Factor de control</span
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
						A<sub>p</sub> = 2 × {detalle.diametro_fase_mm.toFixed(2)} mm × {memoria.instalacion
							.hilos_por_fase}
						hilos = <strong>{detalle.ancho_potencia_mm?.toFixed(2) ?? '—'} mm</strong>
					</p>
					<p class="text-foreground">
						E<sub>f</sub> = ({memoria.instalacion.hilos_por_fase} − 1) × {detalle.factor_triangular?.toFixed(
							2
						)} ×
						{detalle.diametro_fase_mm.toFixed(2)} mm =
						<strong>{detalle.espacio_fuerza_mm.toFixed(2)} mm</strong>
					</p>
					{#if tieneControl && detalle.diametro_control_mm}
						<p class="text-foreground">
							E<sub>c</sub> = f<sub>c</sub> × {detalle.diametro_control_mm.toFixed(2)} mm =
							<strong>{detalle.espacio_control_mm?.toFixed(2) ?? '—'} mm</strong>
							<span class="font-sans text-xs text-muted-foreground">
								(f_c = {factorControl.toFixed(2)})</span
							>
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
						A<sub>req</sub> = {canalizacion.resultado?.area_requerida_mm2?.toFixed(2) ?? '—'} mm
					</p>
				</div>
			</div>
		{/if}

		<!-- Diagrama SVG de arreglo de cables -->
		<div class="mb-6">
			<DiagramaCable
				tipoCanalizacion={memoria.instalacion.tipo_canalizacion}
				detalleCharola={memoria.canalizacion.detalle_charola}
				detalleTuberia={memoria.canalizacion.detalle_tuberia}
				resultado={memoria.canalizacion.resultado}
				sistemaElectrico={memoria.instalacion.sistema_electrico}
				hilosPorFase={memoria.instalacion.hilos_por_fase}
				calibreFase={memoria.cable_fase.calibre}
				calibreTierra={memoria.cable_tierra.calibre}
				{...propsDiagrama}
			/>
		</div>

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
								{canalizacion.resultado?.area_requerida_mm2?.toFixed(1) ?? '—'} mm
							</td>
						</tr>
						<tr>
							<td class="px-4 py-2 text-muted-foreground">Ancho Comercial Seleccionado</td>
							<td class="px-4 py-2 font-mono font-bold text-primary">
								{canalizacion.resultado?.tamano || '—'}
							</td>
						</tr>
						{#if canalizacion.resultado?.ancho_comercial_mm && canalizacion.resultado.ancho_comercial_mm > 0}
							<tr>
								<td class="px-4 py-2 text-muted-foreground">Ancho Comercial (mm)</td>
								<td class="px-4 py-2 font-mono font-bold text-primary">
									{canalizacion.resultado.ancho_comercial_mm.toFixed(1)} mm
								</td>
							</tr>
						{/if}
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
