/**
 * Utilidades geométricas para diagramas SVG de arreglo de cables.
 *
 * Funciones PURAS para calcular posiciones de conductores, viewBox y cotas.
 * Sin dependencias Svelte, sin side effects.
 */

import type { SistemaElectrico } from '$lib/features/calculos/domain/types/calculo.enums.js';

// ============================================================================
// CONSTANTES
// ============================================================================

/** Standard tray wall height: 70mm (common commercial size) */
export const PERALTE_CHAROLA_MM = 70;

/** Floor thickness for SVG rendering (visual only) */
export const ESPESOR_FONDO_CHAROLA_MM = 3;

/** Wall thickness for SVG rendering (visual only) */
export const ESPESOR_PARED_CHAROLA_MM = 3;

/** Flange width for charola (cosmetic, extends outward) */
export const ANCHO_BRIDA_CHAROLA_MM = 15;

/** CAD-style color palette */
export const COLORES = {
	fase: '#D4AF37',
	tierra: '#28A745',
	neutro: '#6B7280',
	control: '#3B82F6',
	charolaStroke: '#004085',
	charolaFill: 'url(#charola-hatch)',
	tuboStroke: '#555555',
	tuboFill: '#E8E8E8',
	cotaLinea: '#374151',
	cotaTexto: '#6B7280',
	titleBlock: '#1F2937',
	fondo: '#FFFFFF'
} as const;

// ============================================================================
// TIPOS
// ============================================================================

/** Position of a conductor in the SVG diagram */
export interface ConductorPosicion {
	cx: number; // center X (mm)
	cy: number; // center Y (mm)
	radio: number; // radius (mm)
	color: string; // fill color hex
	etiqueta: string; // label: "A", "B", "C", "N", "T", "C1"
	tipo: 'fase' | 'neutro' | 'tierra' | 'control';
}

/** SVG viewBox calculation result */
export interface ViewBoxResult {
	viewBox: string; // "minX minY width height"
	ancho: number; // total width
	alto: number; // total height
}

/** Dimension line data for SVG rendering */
export interface LineaCota {
	x1: number;
	y1: number;
	x2: number;
	y2: number;
	valor: number; // dimension value in mm
	texto: string; // display text (e.g., "152.4 mm (6\")")
	posicionTexto: 'arriba' | 'abajo';
}

// ============================================================================
// FUNCIONES DE CÁLCULO DE POSICIONES
// ============================================================================

/**
 * Calculates conductor positions for spaced cable tray (charola espaciada).
 *
 * Geometry rules:
 * - The charola is drawn as U-shape (open top): bottom wall + left wall + right wall + flanges
 * - All conductors sit ON the bottom floor of the charola
 * - Phase conductors: spaced 1 diameter apart (center-to-center = 2 * diametro)
 * - Ground (tierra): TOUCHING the last conductor (no gap)
 * - Control cables: positioned after ground with 1 diameter spacing
 *
 * The cables are CENTERED horizontally within the charola's commercial width.
 *
 * @param params - Configuration parameters
 * @returns Array of conductor positions with cx relative to inner left wall (0 = left wall)
 *          cy is always 0; the component calculates actual Y position based on floor offset
 */
export function calcularPosicionesCharolaEspaciada(params: {
	diametroFaseMM: number;
	diametroTierraMM: number;
	diametroControlMM?: number;
	numHilosTotal: number;
	sistemaElectrico: SistemaElectrico;
	hilosPorFase: number;
	anchoComercialMM: number; // needed for centering
}): ConductorPosicion[] {
	const {
		diametroFaseMM,
		diametroTierraMM,
		diametroControlMM,
		numHilosTotal,
		sistemaElectrico,
		hilosPorFase,
		anchoComercialMM
	} = params;

	const radioFase = diametroFaseMM / 2;
	const radioTierra = diametroTierraMM / 2;
	const radioControl = diametroControlMM ? diametroControlMM / 2 : 0;

	const posiciones: ConductorPosicion[] = [];

	// Determine number of phases and neutral based on electrical system
	let numFases: number;
	let tieneNeutro: boolean;

	switch (sistemaElectrico) {
		case 'DELTA':
			numFases = 3;
			tieneNeutro = false;
			break;
		case 'ESTRELLA':
			numFases = 3;
			tieneNeutro = true;
			break;
		case 'BIFASICO':
			numFases = 2;
			tieneNeutro = true;
			break;
		case 'MONOFASICO':
			numFases = 1;
			tieneNeutro = true;
			break;
		default:
			numFases = 3;
			tieneNeutro = false;
	}

	// Build list of all phase+neutral wires first (for index-based positioning)
	type WireInfo = { etiqueta: string; tipo: 'fase' | 'neutro' };
	const wires: WireInfo[] = [];

	const faseLabels = ['A', 'B', 'C'];

	// Add phase wires
	for (let faseIdx = 0; faseIdx < numFases; faseIdx++) {
		const labelBase = faseLabels[faseIdx]!;
		for (let hiloIdx = 0; hiloIdx < hilosPorFase; hiloIdx++) {
			const etiqueta = hilosPorFase > 1 ? `${labelBase}${hiloIdx + 1}` : labelBase;
			wires.push({ etiqueta, tipo: 'fase' });
		}
	}

	// Add neutral wires
	if (tieneNeutro) {
		for (let hiloIdx = 0; hiloIdx < hilosPorFase; hiloIdx++) {
			const etiqueta = hilosPorFase > 1 ? `N${hiloIdx + 1}` : 'N';
			wires.push({ etiqueta, tipo: 'neutro' });
		}
	}

	// Calculate number of phase+neutral cables
	const numPhaseNeutroWires = wires.length;

	// Calculate total width occupied by phase+neutral cables:
	// - First cable left edge to last cable right edge
	// - Each cable takes 1 diameter, each gap takes 1 diameter
	// - Total: (2 * numWires - 1) * diametro
	const totalPhaseNeutroWidth = (2 * numPhaseNeutroWires - 1) * diametroFaseMM;

	// Total width: phase+neutral width + tierra diameter (from last cable right edge to tierra right edge)
	let totalWidth = totalPhaseNeutroWidth + diametroTierraMM;

	// Add control cables: 1 diameter gap after tierra, then control cables
	if (diametroControlMM && numHilosTotal > 0) {
		// Space after tierra: 1 control diameter
		// Control cables: (2 * numHilosTotal - 1) * diametroControl
		totalWidth += diametroControlMM + (2 * numHilosTotal - 1) * diametroControlMM;
	}

	// Calculate centering offset: center of cables should be at anchoComercialMM / 2
	const offsetX = (anchoComercialMM - totalWidth) / 2;

	// Y position: all cables sit on the floor (cy will be calculated by component)
	const floorY = 0;

	// Place phase+neutral cables with 2*diametro center-to-center spacing
	for (let i = 0; i < wires.length; i++) {
		const wire = wires[i]!;
		const cx = offsetX + radioFase + i * 2 * diametroFaseMM;

		posiciones.push({
			cx,
			cy: floorY,
			radio: radioFase,
			color: wire.tipo === 'fase' ? COLORES.fase : COLORES.neutro,
			etiqueta: wire.etiqueta,
			tipo: wire.tipo
		});
	}

	// Ground conductor: TOUCHING the last phase/neutral cable
	// Last cable center = offsetX + radioFase + (wires.length - 1) * 2 * diametroFaseMM
	// Last cable right edge = lastCableCenter + radioFase
	// Tierra center = lastCableRightEdge + radioTierra
	const lastWireCX = offsetX + radioFase + (wires.length - 1) * 2 * diametroFaseMM;
	const tierraCX = lastWireCX + radioFase + radioTierra;

	posiciones.push({
		cx: tierraCX,
		cy: floorY,
		radio: radioTierra,
		color: COLORES.tierra,
		etiqueta: 'T',
		tipo: 'tierra'
	});

	// Control cables (if exists): positioned after ground with 1 diameter spacing
	if (diametroControlMM && numHilosTotal > 0) {
		const controlStartCX = tierraCX + radioTierra + diametroControlMM;

		for (let i = 0; i < numHilosTotal; i++) {
			const etiqueta = numHilosTotal > 1 ? `Ctrl${i + 1}` : 'Ctrl';

			posiciones.push({
				cx: controlStartCX + i * 2 * diametroControlMM,
				cy: floorY,
				radio: radioControl,
				color: COLORES.control,
				etiqueta,
				tipo: 'control'
			});
		}
	}

	return posiciones;
}

/**
 * Calculates conductor positions for triangular cable tray arrangement.
 *
 * Geometry rules vary by electrical system:
 *
 * DELTA (3-phase, no neutral) — 1 hilo = 3 cables in equilateral triangle:
 *       C
 *      / \
 *    A --- B
 *
 * ESTRELLA (3-phase with neutral) — 1 hilo = 4 cables in diamond/rectangle:
 *   C --- N
 *    |  X  |
 *   A --- B
 *
 * BIFASICO (2-phase with neutral) — 1 hilo = 3 cables (A, B, N in triangle):
 *       N
 *      / \
 *    A --- B
 *
 * MONOFASICO (1-phase with neutral) — 1 hilo = 2 cables stacked:
 *     N
 *     |
 *     A
 *
 * - Cables within a group touch each other
 * - Groups spaced by factorTriangular × diametro center-to-center (NOM standard)
 * - Ground (T): TOUCHING the last bottom-row cable (B of last group)
 * - Control cables: positioned after ground with 1 diameter spacing
 *
 * The function returns positions RELATIVE to the inner left wall of the charola.
 * - cx: horizontal offset from inner left wall (0 = at the wall)
 * - cy: vertical offset from floor (0 = on floor, positive = above floor)
 *
 * The component calculates actual SVG coordinates:
 * - SVG_cx = charolaLeft + conductor.cx
 * - SVG_cy = charolaBottom - ESPESOR_PISO - conductor.radio - conductor.cy
 */
export function calcularPosicionesCharolaTriangular(params: {
	diametroFaseMM: number;
	diametroTierraMM: number;
	diametroControlMM?: number;
	hilosPorFase: number;
	factorTriangular: number;
	anchoComercialMM: number;
	sistemaElectrico: SistemaElectrico;
}): ConductorPosicion[] {
	const {
		diametroFaseMM,
		diametroTierraMM,
		diametroControlMM,
		hilosPorFase,
		factorTriangular,
		anchoComercialMM,
		sistemaElectrico
	} = params;

	const radioFase = diametroFaseMM / 2;
	const radioTierra = diametroTierraMM / 2;
	const radioControl = diametroControlMM ? diametroControlMM / 2 : 0;

	// sin(60°) for equilateral triangle height
	const sin60 = Math.sqrt(3) / 2;
	const alturaTriangulo = diametroFaseMM * sin60;

	// Determine geometry based on electrical system type
	switch (sistemaElectrico) {
		case 'DELTA':
			// 3 conductors: A, B, C
			break;
		case 'ESTRELLA':
			// 4 conductors: A, B, C, N
			break;
		case 'BIFASICO':
			// 3 conductors: A, B, N
			break;
		case 'MONOFASICO':
			// 2 conductors: A, N (stacked)
			break;
		default:
			// Default to DELTA behavior
			break;
	}

	// Calculate group width (horizontal extent) per system type
	// DELTA, ESTRELLA, BIFASICO: 2 × Ø (A and B define bottom width)
	// MONOFASICO: 1 × Ø (just A, N stacked)
	const groupWidth = sistemaElectrico === 'MONOFASICO' ? diametroFaseMM : 2 * diametroFaseMM;

	// Calculate total width for centering BEFORE placing conductors
	// For N groups: totalPhaseWidth = N * groupWidth + (N - 1) * factorTriangular * diametroFaseMM
	const totalPhaseWidth =
		hilosPorFase * groupWidth + (hilosPorFase - 1) * factorTriangular * diametroFaseMM;

	// Total width including tierra (touching last group)
	let totalWidth = totalPhaseWidth + diametroTierraMM;

	// Add control cables if present (1 diameter gap from tierra + 1 diameter for cable)
	let numControlCables = 0;
	if (diametroControlMM && radioControl > 0) {
		numControlCables = 1;
		totalWidth += diametroControlMM + diametroControlMM; // gap + cable
	}

	// Calculate centering offset: center of cables should be at anchoComercialMM / 2
	const offsetX = (anchoComercialMM - totalWidth) / 2;

	const posiciones: ConductorPosicion[] = [];

	// Place phase conductors in groups
	for (let hiloIdx = 0; hiloIdx < hilosPorFase; hiloIdx++) {
		// Group offset: stride = groupWidth + spacing between groups
		const groupStep = groupWidth + factorTriangular * diametroFaseMM;
		const groupOffsetX = hiloIdx * groupStep;

		// A: bottom left (cy = 0 = on floor)
		posiciones.push({
			cx: offsetX + groupOffsetX + radioFase,
			cy: 0,
			radio: radioFase,
			color: COLORES.fase,
			etiqueta: hilosPorFase > 1 ? `A${hiloIdx + 1}` : 'A',
			tipo: 'fase'
		});

		// B: bottom right — only for systems with 2+ phases (not MONOFASICO)
		if (sistemaElectrico !== 'MONOFASICO') {
			posiciones.push({
				cx: offsetX + groupOffsetX + radioFase + diametroFaseMM,
				cy: 0,
				radio: radioFase,
				color: COLORES.fase,
				etiqueta: hilosPorFase > 1 ? `B${hiloIdx + 1}` : 'B',
				tipo: 'fase'
			});
		}

		// Third (and fourth) conductors depend on system type
		if (sistemaElectrico === 'DELTA') {
			// C: top center, centered between A and B (equilateral triangle)
			posiciones.push({
				cx: offsetX + groupOffsetX + diametroFaseMM,
				cy: alturaTriangulo,
				radio: radioFase,
				color: COLORES.fase,
				etiqueta: hilosPorFase > 1 ? `C${hiloIdx + 1}` : 'C',
				tipo: 'fase'
			});
		} else if (sistemaElectrico === 'ESTRELLA') {
			// C: top left (above A), N: top right (above B) — both at triangle height
			// C aligned with A (left edge of group), N aligned with B (right edge of group)
			posiciones.push({
				cx: offsetX + groupOffsetX + radioFase,
				cy: alturaTriangulo,
				radio: radioFase,
				color: COLORES.fase,
				etiqueta: hilosPorFase > 1 ? `C${hiloIdx + 1}` : 'C',
				tipo: 'fase'
			});
			posiciones.push({
				cx: offsetX + groupOffsetX + radioFase + diametroFaseMM,
				cy: alturaTriangulo,
				radio: radioFase,
				color: COLORES.neutro,
				etiqueta: hilosPorFase > 1 ? `N${hiloIdx + 1}` : 'N',
				tipo: 'neutro'
			});
		} else if (sistemaElectrico === 'BIFASICO') {
			// N: top center, like DELTA's C position
			posiciones.push({
				cx: offsetX + groupOffsetX + diametroFaseMM,
				cy: alturaTriangulo,
				radio: radioFase,
				color: COLORES.neutro,
				etiqueta: hilosPorFase > 1 ? `N${hiloIdx + 1}` : 'N',
				tipo: 'neutro'
			});
		} else if (sistemaElectrico === 'MONOFASICO') {
			// N: directly above A (stacked, touching)
			posiciones.push({
				cx: offsetX + groupOffsetX + radioFase,
				cy: diametroFaseMM, // Full diameter height (not triangle height) for stacked arrangement
				radio: radioFase,
				color: COLORES.neutro,
				etiqueta: hilosPorFase > 1 ? `N${hiloIdx + 1}` : 'N',
				tipo: 'neutro'
			});
		}
	}

	// Ground conductor: TOUCHING the last bottom-row cable
	// For MONOFASICO: last A cable (no B exists)
	// For other systems: last B cable (rightmost bottom conductor)
	const lastGroupStep = (hilosPorFase - 1) * (groupWidth + factorTriangular * diametroFaseMM);
	let lastBottomCX: number;
	if (sistemaElectrico === 'MONOFASICO') {
		// Last A cable center (no B in MONOFASICO)
		lastBottomCX = offsetX + lastGroupStep + radioFase;
	} else {
		// Last B cable center
		lastBottomCX = offsetX + lastGroupStep + radioFase + diametroFaseMM;
	}
	const tierraCX = lastBottomCX + radioFase + radioTierra;

	posiciones.push({
		cx: tierraCX,
		cy: 0,
		radio: radioTierra,
		color: COLORES.tierra,
		etiqueta: 'T',
		tipo: 'tierra'
	});

	// Control cables (if exists): after ground with 1 diameter spacing
	// Gap = 1 full control diameter + cable itself
	if (diametroControlMM && numControlCables > 0) {
		const controlStartCX = tierraCX + radioTierra + diametroControlMM + radioControl;

		posiciones.push({
			cx: controlStartCX,
			cy: 0,
			radio: radioControl,
			color: COLORES.control,
			etiqueta: 'Ctrl',
			tipo: 'control'
		});
	}

	return posiciones;
}

/**
 * Calculates conductor positions for conduit cross-section.
 *
 * Geometry rules:
 * - The tube is a single thick circle representing pipe wall thickness
 * - Center of tube is at (0, 0) — the component will translate
 * - Conductor radii calculated from area: r = √(area / π)
 * - Cables sit at the BOTTOM of the tube due to gravity (like triangular charola)
 * - Phase+neutral conductors packed in rows at the bottom
 * - Ground (T): TOUCHING the rightmost bottom-row conductor
 *
 * Coordinate system: positions are relative to tube center (0, 0)
 * - +Y direction = down (SVG convention)
 * - Bottom of tube interior = (0, radioInterior)
 * - Bottom row cable centers = (cx, radioInterior - radioFase)
 */
export function calcularPosicionesTuberia(params: {
	diametroInteriorMM: number;
	diametroExteriorMM: number;
	areaFaseMM2: number;
	areaNeutroMM2?: number;
	areaTierraMM2: number;
	numFasesPorTubo: number;
	numNeutrosPorTubo: number;
	numTierras: number;
	sistemaElectrico: SistemaElectrico;
}): ConductorPosicion[] {
	const {
		diametroInteriorMM,
		areaFaseMM2,
		areaNeutroMM2,
		areaTierraMM2,
		numFasesPorTubo,
		numNeutrosPorTubo,
		numTierras,
		sistemaElectrico
	} = params;

	const R = diametroInteriorMM / 2; // Inner radius of tube
	const rF = Math.sqrt(areaFaseMM2 / Math.PI); // Phase cable radius (from area)
	const rN = areaNeutroMM2 ? Math.sqrt(areaNeutroMM2 / Math.PI) : rF; // Neutral radius
	const rT = Math.sqrt(areaTierraMM2 / Math.PI); // Tierra radius

	const posiciones: ConductorPosicion[] = [];

	if (numFasesPorTubo === 0 && numNeutrosPorTubo === 0 && numTierras === 0) {
		return posiciones;
	}

	// Helper: place a cable on the tube inner wall at angle theta (0 = bottom, positive = clockwise)
	// Returns {cx, cy} relative to tube center, +Y = down
	function posOnWall(r: number, theta: number): { cx: number; cy: number } {
		const d = R - r; // distance from tube center to cable center
		return {
			cx: d * Math.sin(theta),
			cy: d * Math.cos(theta) // +Y = down, so cos(0) = bottom
		};
	}

	// Helper: angle for two touching cables on the wall
	// Two cables of radius r sitting on the inner wall of radius R, touching each other
	// Half-angle from vertical axis to each cable center
	function halfAngleTwoTouching(r1: number, r2: number): number {
		const d1 = R - r1;
		const d2 = R - r2;
		const touchDist = r1 + r2;
		// Law of cosines: touchDist² = d1² + d2² - 2·d1·d2·cos(angle)
		const cosAngle = (d1 * d1 + d2 * d2 - touchDist * touchDist) / (2 * d1 * d2);
		const fullAngle = Math.acos(Math.max(-1, Math.min(1, cosAngle)));
		// Return HALF angle (each cable is offset by halfAngle from center axis)
		return fullAngle / 2;
	}

	// ============================================================
	// STRATEGY BY SYSTEM TYPE
	// ============================================================
	// Coordinate system: center of tube = (0,0), +Y = down (SVG)
	// Bottom of tube = (0, R)

	if (sistemaElectrico === 'MONOFASICO') {
		// 1 phase (A) at bottom center on wall + 1 neutral (N) on top touching A
		const posA = posOnWall(rF, 0); // bottom center
		posiciones.push({
			cx: posA.cx,
			cy: posA.cy,
			radio: rF,
			color: COLORES.fase,
			etiqueta: 'A',
			tipo: 'fase'
		});

		if (numNeutrosPorTubo > 0) {
			// N directly above A, touching A (not necessarily on the wall)
			posiciones.push({
				cx: posA.cx,
				cy: posA.cy - rF - rN,
				radio: rN,
				color: COLORES.neutro,
				etiqueta: 'N',
				tipo: 'neutro'
			});
		}
	} else if (sistemaElectrico === 'BIFASICO') {
		// 2 phases (A, B) on bottom touching each other + 1 neutral (N) on top
		const alpha = halfAngleTwoTouching(rF, rF);
		const posA = posOnWall(rF, -alpha); // left
		const posB = posOnWall(rF, alpha); // right

		posiciones.push({
			cx: posA.cx,
			cy: posA.cy,
			radio: rF,
			color: COLORES.fase,
			etiqueta: 'A',
			tipo: 'fase'
		});
		posiciones.push({
			cx: posB.cx,
			cy: posB.cy,
			radio: rF,
			color: COLORES.fase,
			etiqueta: 'B',
			tipo: 'fase'
		});

		if (numNeutrosPorTubo > 0) {
			// N nestled on top of A and B, centered
			// triHeight (for reference): rF * Math.sqrt(3) — height of equilateral triangle with side 2*rF (if rN=rF)
			// More precisely: N at cx=0, cy = A.cy - sqrt((rF+rN)² - (posA.cx)²)
			const dxN = 0 - posA.cx;
			const contactDist = rF + rN;
			const dyN = Math.sqrt(Math.max(0, contactDist * contactDist - dxN * dxN));
			posiciones.push({
				cx: 0,
				cy: posA.cy - dyN,
				radio: rN,
				color: COLORES.neutro,
				etiqueta: 'N',
				tipo: 'neutro'
			});
		}
	} else if (sistemaElectrico === 'DELTA') {
		// 3 phases (A, B, C): A,B on bottom touching, C on top nestled
		const alpha = halfAngleTwoTouching(rF, rF);
		const posA = posOnWall(rF, -alpha);
		const posB = posOnWall(rF, alpha);

		posiciones.push({
			cx: posA.cx,
			cy: posA.cy,
			radio: rF,
			color: COLORES.fase,
			etiqueta: 'A',
			tipo: 'fase'
		});
		posiciones.push({
			cx: posB.cx,
			cy: posB.cy,
			radio: rF,
			color: COLORES.fase,
			etiqueta: 'B',
			tipo: 'fase'
		});

		// C nestled on top of A and B (equilateral triangle, all same radius)
		// Distance from C center to A center must = 2*rF
		// C is at cx=0: sqrt(A.cx² + (C.cy - A.cy)²) = 2*rF
		// C.cy = A.cy - sqrt((2*rF)² - A.cx²)
		const cCy = posA.cy - Math.sqrt(Math.max(0, 2 * rF * (2 * rF) - posA.cx * posA.cx));
		posiciones.push({
			cx: 0,
			cy: cCy,
			radio: rF,
			color: COLORES.fase,
			etiqueta: 'C',
			tipo: 'fase'
		});
	} else if (sistemaElectrico === 'ESTRELLA') {
		// 4 cables: A,B on bottom, C top-left, N top-right (diamond pattern)
		const alpha = halfAngleTwoTouching(rF, rF);
		const posA = posOnWall(rF, -alpha);
		const posB = posOnWall(rF, alpha);

		posiciones.push({
			cx: posA.cx,
			cy: posA.cy,
			radio: rF,
			color: COLORES.fase,
			etiqueta: 'A',
			tipo: 'fase'
		});
		posiciones.push({
			cx: posB.cx,
			cy: posB.cy,
			radio: rF,
			color: COLORES.fase,
			etiqueta: 'B',
			tipo: 'fase'
		});

		// C: on top of A (touching A, positioned directly above-left)
		// C touches A: distance = 2*rF. C_cx = A_cx, C_cy = A_cy - 2*rF
		// Actually for diamond: C above A, N above B
		posiciones.push({
			cx: posA.cx,
			cy: posA.cy - 2 * rF,
			radio: rF,
			color: COLORES.fase,
			etiqueta: 'C',
			tipo: 'fase'
		});
		posiciones.push({
			cx: posB.cx,
			cy: posB.cy - 2 * rN,
			radio: rN,
			color: COLORES.neutro,
			etiqueta: 'N',
			tipo: 'neutro'
		});
	}

	// ============================================================
	// TIERRA (T): sits on tube wall, touching the rightmost bottom cable
	// ============================================================
	if (numTierras > 0 && posiciones.length > 0) {
		// Find the rightmost bottom-row cable (the one with the largest cx among cables on the wall)
		// For all systems: B is rightmost (or A if MONOFASICO)
		let rightmostBottom: ConductorPosicion;
		if (sistemaElectrico === 'MONOFASICO') {
			rightmostBottom = posiciones[0]!; // A
		} else {
			rightmostBottom = posiciones[1]!; // B
		}

		// T sits on the wall and touches the rightmost cable
		// T center is at distance (R - rT) from tube center
		// Distance between T center and rightmost cable center = rightmost.radio + rT
		// Use geometry: find angle of T on the wall such that it touches the rightmost cable
		const dRight = Math.sqrt(
			rightmostBottom.cx * rightmostBottom.cx + rightmostBottom.cy * rightmostBottom.cy
		);
		const dT = R - rT;
		const touchDist = rightmostBottom.radio + rT;
		// Law of cosines to find angle between rightmostBottom and T (from tube center)
		const cosGamma = (dRight * dRight + dT * dT - touchDist * touchDist) / (2 * dRight * dT);
		const gamma = Math.acos(Math.max(-1, Math.min(1, cosGamma)));
		// Angle of rightmost cable from +Y axis (bottom)
		const angleRight = Math.atan2(rightmostBottom.cx, rightmostBottom.cy);
		// T is placed clockwise from rightmost cable (to the right and slightly up)
		const angleT = angleRight + gamma;
		const posTierra = posOnWall(rT, angleT);

		posiciones.push({
			cx: posTierra.cx,
			cy: posTierra.cy,
			radio: rT,
			color: COLORES.tierra,
			etiqueta: 'T',
			tipo: 'tierra'
		});
	}

	return posiciones;
}

// ============================================================================
// FUNCIONES DE VIEWBOX Y COTAS
// ============================================================================

/**
 * Calculate SVG viewBox based on content dimensions.
 */
export function calcularViewBox(params: {
	anchoContenido: number;
	altoContenido: number;
	margen?: number;
}): ViewBoxResult {
	const { anchoContenido, altoContenido, margen = 20 } = params;

	const minX = -margen;
	const minY = -margen;
	const ancho = anchoContenido + 2 * margen;
	const alto = altoContenido + 2 * margen;

	return {
		viewBox: `${minX} ${minY} ${ancho} ${alto}`,
		ancho,
		alto
	};
}

/**
 * Calculate dimension lines for charola (cotas).
 *
 * Returns 2 dimension lines:
 * - Top: full commercial width of the tray
 * - Bottom: required width (A_req)
 */
export function calcularCotasCharola(params: {
	anchoComercialMM: number;
	areaRequeridaMM: number;
	peralte: number;
}): LineaCota[] {
	const { anchoComercialMM, areaRequeridaMM, peralte } = params;

	// Calculate required width from area: A_req = area / peralte
	const anchoRequerido = areaRequeridaMM / peralte;

	// Convert mm to inches for display
	const comercialPulgadas = anchoComercialMM / 25.4;
	const requeridoPulgadas = anchoRequerido / 25.4;

	const cotitas: LineaCota[] = [];

	// Top dimension line: full commercial width
	cotitas.push({
		x1: 0,
		y1: -20,
		x2: anchoComercialMM,
		y2: -20,
		valor: anchoComercialMM,
		texto: `${anchoComercialMM.toFixed(1)} mm (${comercialPulgadas.toFixed(1)}")`,
		posicionTexto: 'arriba'
	});

	// Bottom dimension line: required width
	cotitas.push({
		x1: 0,
		y1: peralte + 20,
		x2: anchoRequerido,
		y2: peralte + 20,
		valor: anchoRequerido,
		texto: `${anchoRequerido.toFixed(1)} mm (${requeridoPulgadas.toFixed(1)}")`,
		posicionTexto: 'abajo'
	});

	return cotitas;
}

/**
 * Calculate total width occupied by conductors in charola.
 */
export function calcularAnchoOcupadoCharola(posiciones: ConductorPosicion[]): number {
	if (posiciones.length === 0) {
		return 0;
	}

	// Find min cx and max cx (including radius)
	let minCX = Infinity;
	let maxCX = -Infinity;

	for (const pos of posiciones) {
		const izquierda = pos.cx - pos.radio;
		const derecha = pos.cx + pos.radio;

		if (izquierda < minCX) {
			minCX = izquierda;
		}
		if (derecha > maxCX) {
			maxCX = derecha;
		}
	}

	return maxCX - minCX;
}
