<script lang="ts">
	import type { DetalleTuberia, ResultadoCanalizacion } from '$lib/types/calculos.types';
	import type { SistemaElectrico } from '$lib/features/calculos/domain/types/calculo.enums.js';
	import type { DiagramaOutput } from '$lib/features/calculos/domain/types/index.js';
	import { calcularPosicionesTuberia } from './geometry.js';

	interface Props {
		detalle: DetalleTuberia;
		resultado: ResultadoCanalizacion;
		sistemaElectrico: SistemaElectrico;
		calibreFase?: string;
		calibreTierra?: string;
		diagrama?: DiagramaOutput;
	}

	// Props destructured
	let { detalle, resultado, sistemaElectrico, diagrama = undefined }: Props = $props();

	// Tube dimensions
	let diametroInterior = $derived(detalle.diametro_interior_mm ?? 0);
	let diametroExterior = $derived(detalle.diametro_exterior_mm ?? 0);
	let radioInterior = $derived(diametroInterior / 2);
	let radioExterior = $derived(diametroExterior / 2);
	let numeroTubos = $derived(resultado.numero_de_tubos);

	// Wall thickness: difference between exterior and interior radii
	let wallThickness = $derived((diametroExterior - diametroInterior) / 2);

	// Conductor positions (centered at tube center 0,0)
	let paramsTuberia = $derived.by(() => {
		const base = {
			diametroInteriorMM: diametroInterior,
			diametroExteriorMM: diametroExterior,
			areaFaseMM2: detalle.area_fase_mm2,
			areaTierraMM2: detalle.area_tierra_mm2,
			numFasesPorTubo: detalle.num_fases_por_tubo,
			numNeutrosPorTubo: detalle.num_neutros_por_tubo,
			numTierras: detalle.num_tierras,
			sistemaElectrico
		};
		if (detalle.area_neutro_mm2 !== undefined) {
			return { ...base, areaNeutroMM2: detalle.area_neutro_mm2 };
		}
		return base;
	});

	let posiciones = $derived(calcularPosicionesTuberia(paramsTuberia));

	// Layout constants
	const MARGEN = 30;
	const ESPACIO_ENTRE_TUBOS = 30;
	const MARGEN_COTA = 25;
	const MARGEN_ETIQUETA = 15;

	// Total width for all tubes (using exterior diameter for spacing)
	let anchoTubos = $derived(
		numeroTubos * diametroExterior + (numeroTubos - 1) * ESPACIO_ENTRE_TUBOS
	);

	// ViewBox calculation
	let svgWidth = $derived(anchoTubos + 2 * MARGEN);
	let svgHeight = $derived(
		diametroExterior + 2 * MARGEN + MARGEN_COTA + (numeroTubos > 1 ? MARGEN_ETIQUETA : 0)
	);
	let viewBoxStr = $derived(`0 0 ${svgWidth} ${svgHeight}`);

	// Center Y for all tubes (same vertical center)
	let centerY = $derived(MARGEN + radioExterior);
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
			aria-label="Sección transversal de tubería"
		>
			{#each Array(numeroTubos) as _, idx}
				{@const centerX = MARGEN + radioExterior + idx * (diametroExterior + ESPACIO_ENTRE_TUBOS)}
				{@const dimY = centerY + radioExterior + 15}

				<g>
					<!-- SINGLE thick circle for the tube (like real pipe cross-section) -->
					<!-- r = midpoint of wall, stroke-width = actual wall thickness -->
					<!-- Using midline radius ensures stroke edges align with interior and exterior -->
					<circle
						cx={centerX}
						cy={centerY}
						r={(radioInterior + radioExterior) / 2}
						stroke="black"
						stroke-width={wallThickness}
						fill="none"
					/>

					<!-- Conductors at the bottom of the tube -->
					<g stroke="black" stroke-width="1" fill="none">
						{#each posiciones as conductor}
							<circle cx={centerX + conductor.cx} cy={centerY + conductor.cy} r={conductor.radio} />
							<text
								x={centerX + conductor.cx}
								y={centerY + conductor.cy}
								font-size={Math.max(conductor.radio * 0.9, 4)}
								text-anchor="middle"
								dominant-baseline="central"
								stroke="none"
								fill="black"
								font-weight="bold"
								font-family="Arial, sans-serif"
							>
								{conductor.etiqueta}
							</text>
						{/each}
					</g>

					<!-- Diameter dimension line (horizontal across exterior) -->
					<line
						x1={centerX - radioExterior}
						y1={dimY}
						x2={centerX + radioExterior}
						y2={dimY}
						stroke="black"
						stroke-width="1.5"
					/>
					<!-- Ticks -->
					<line
						x1={centerX - radioExterior}
						y1={dimY - 5}
						x2={centerX - radioExterior}
						y2={dimY + 5}
						stroke="black"
						stroke-width="1.5"
					/>
					<line
						x1={centerX + radioExterior}
						y1={dimY - 5}
						x2={centerX + radioExterior}
						y2={dimY + 5}
						stroke="black"
						stroke-width="1.5"
					/>
					<!-- Dimension label -->
					<text
						x={centerX}
						y={dimY + 12}
						font-size="6"
						text-anchor="middle"
						fill="black"
						font-family="Arial, sans-serif"
					>
						Ø {diametroExterior.toFixed(1)} mm
					</text>

					<!-- Tube label (only if multiple tubes) -->
					{#if numeroTubos > 1}
						<text
							x={centerX}
							y={dimY + 22}
							font-size="5"
							text-anchor="middle"
							fill="black"
							font-family="Arial, sans-serif"
						>
							Tubo {idx + 1} de {numeroTubos}
						</text>
					{/if}
				</g>
			{/each}
		</svg>
	{/if}
</figure>
