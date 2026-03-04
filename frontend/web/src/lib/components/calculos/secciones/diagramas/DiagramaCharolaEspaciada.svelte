<script lang="ts">
	import type { DetalleCharola, ResultadoCanalizacion } from '$lib/types/calculos.types';
	import type { SistemaElectrico } from '$lib/features/calculos/domain/types/calculo.enums.js';
	import type { DiagramaOutput } from '$lib/features/calculos/domain/types/index.js';
	import { calcularPosicionesCharolaEspaciada, PERALTE_CHAROLA_MM } from './geometry.js';

	interface Props {
		detalle: DetalleCharola;
		resultado: ResultadoCanalizacion;
		sistemaElectrico: SistemaElectrico;
		hilosPorFase: number;
		calibreFase?: string;
		calibreTierra?: string;
		diagrama?: DiagramaOutput;
	}

	// Props destructured
	let {
		detalle,
		resultado,
		sistemaElectrico,
		hilosPorFase,
		diagrama = undefined
	}: Props = $props();

	// Constants for layout
	const MARGEN_IZQ = 30; // margin before left flange
	const BRIDA = 25; // flange width extending outward
	const MARGEN_COTA_DER = 60; // space for P dimension on right
	const MARGEN_ARRIBA = 35; // space above for W dimension
	const MARGEN_ABAJO = 20; // space below charola
	const ESPESOR_PISO = 7; // floor offset (where cables sit above bottom)

	let peralte = PERALTE_CHAROLA_MM; // 70mm
	let anchoComercial = $derived(resultado.ancho_comercial_mm ?? 150);

	// Charola geometry (inner coordinates)
	let charolaLeft = $derived(MARGEN_IZQ + BRIDA); // inner left wall x
	let charolaRight = $derived(charolaLeft + anchoComercial); // inner right wall x
	let charolaTop = $derived(MARGEN_ARRIBA); // top of walls
	let charolaBottom = $derived(charolaTop + peralte); // bottom floor

	// Flanges
	let flangeLeft = $derived(charolaLeft - BRIDA);
	let flangeRight = $derived(charolaRight + BRIDA);

	// Build params for geometry function
	let paramsCharola = $derived.by(() => {
		const anchoComercialMM = anchoComercial;
		const base = {
			diametroFaseMM: detalle.diametro_fase_mm,
			diametroTierraMM: detalle.diametro_tierra_mm,
			numHilosTotal: detalle.num_hilos_total ?? 0,
			sistemaElectrico,
			hilosPorFase,
			anchoComercialMM
		};
		if (detalle.diametro_control_mm && detalle.diametro_control_mm > 0) {
			return { ...base, diametroControlMM: detalle.diametro_control_mm };
		}
		return base;
	});

	// Conductor positions (cx relative to inner left wall)
	let posiciones = $derived(calcularPosicionesCharolaEspaciada(paramsCharola));

	// ViewBox
	let svgWidth = $derived(charolaRight + BRIDA + MARGEN_COTA_DER);
	let svgHeight = $derived(charolaBottom + MARGEN_ABAJO);
	let viewBoxStr = $derived(`0 0 ${svgWidth} ${svgHeight}`);

	// W dimension line positions
	let wLineY = $derived(charolaTop - 15);
	let wMidX = $derived(charolaLeft + anchoComercial / 2);

	// P dimension line positions
	let pLineX = $derived(flangeRight + 20);
	let pMidY = $derived(charolaTop + peralte / 2);

	// Floor lines
	let floorSolidY = $derived(charolaBottom - 5);
	let floorDashedY = $derived(charolaBottom - 3);
</script>

<figure class="my-4">
	{#if diagrama?.svg}
		{@html diagrama.svg}
	{:else}
		<svg
			class="mx-auto h-auto w-full max-w-2xl"
			viewBox={viewBoxStr}
			preserveAspectRatio="xMidYMid meet"
			role="img"
			aria-label="Diagrama de arreglo de cables en charola espaciada"
		>
			<!-- Charola U-profile with flanges -->
			<path
				d="M {flangeLeft},{charolaTop} L {charolaLeft},{charolaTop} L {charolaLeft},{charolaBottom} L {charolaRight},{charolaBottom} L {charolaRight},{charolaTop} L {flangeRight},{charolaTop}"
				stroke="black"
				stroke-width="2.5"
				fill="none"
				stroke-linejoin="round"
			/>

			<!-- W dimension: horizontal line above charola -->
			<line
				x1={charolaLeft}
				y1={wLineY}
				x2={charolaRight}
				y2={wLineY}
				stroke="black"
				stroke-width="1.5"
			/>
			<!-- W ticks -->
			<line
				x1={charolaLeft}
				y1={wLineY - 5}
				x2={charolaLeft}
				y2={wLineY + 5}
				stroke="black"
				stroke-width="1.5"
			/>
			<line
				x1={charolaRight}
				y1={wLineY - 5}
				x2={charolaRight}
				y2={wLineY + 5}
				stroke="black"
				stroke-width="1.5"
			/>
			<!-- W label -->
			<text
				x={wMidX}
				y={wLineY - 6}
				font-size="8"
				text-anchor="middle"
				fill="black"
				font-family="Arial, sans-serif"
			>
				{anchoComercial.toFixed(1)} mm
			</text>

			<!-- P dimension: vertical line to the right -->
			<line
				x1={pLineX}
				y1={charolaTop}
				x2={pLineX}
				y2={charolaBottom}
				stroke="black"
				stroke-width="1.5"
			/>
			<!-- P ticks -->
			<line
				x1={pLineX - 5}
				y1={charolaTop}
				x2={pLineX + 5}
				y2={charolaTop}
				stroke="black"
				stroke-width="1.5"
			/>
			<line
				x1={pLineX - 5}
				y1={charolaBottom}
				x2={pLineX + 5}
				y2={charolaBottom}
				stroke="black"
				stroke-width="1.5"
			/>
			<!-- P label -->
			<text
				x={pLineX + 8}
				y={pMidY + 3}
				font-size="8"
				text-anchor="start"
				fill="black"
				font-family="Arial, sans-serif"
			>
				{peralte} mm
			</text>

			<!-- Floor pattern lines -->
			<line
				x1={charolaLeft}
				y1={floorSolidY}
				x2={charolaRight}
				y2={floorSolidY}
				stroke="black"
				stroke-width="2"
			/>
			<line
				x1={charolaLeft}
				y1={floorDashedY}
				x2={charolaRight}
				y2={floorDashedY}
				stroke="black"
				stroke-width="2"
				stroke-dasharray="8,5"
			/>

			<!-- Conductors -->
			<g stroke="black" stroke-width="2" fill="none">
				{#each posiciones as conductor}
					<!-- Cable circle (cx relative to inner left wall + charolaLeft offset) -->
					<circle
						cx={charolaLeft + conductor.cx}
						cy={charolaBottom - ESPESOR_PISO - conductor.radio}
						r={conductor.radio}
					/>
					<!-- Label inside circle -->
					<text
						x={charolaLeft + conductor.cx}
						y={charolaBottom - ESPESOR_PISO - conductor.radio + conductor.radio * 0.35}
						font-size={Math.max(conductor.radio * 0.8, 4)}
						text-anchor="middle"
						stroke="none"
						fill="black"
						font-weight="bold"
					>
						{conductor.etiqueta}
					</text>
				{/each}
			</g>
		</svg>
	{/if}
</figure>
