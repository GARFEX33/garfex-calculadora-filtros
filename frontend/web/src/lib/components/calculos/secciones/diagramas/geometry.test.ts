/**
 * Unit tests for geometry.ts - pure geometric calculation functions
 * for SVG cable tray and conduit diagrams.
 *
 * No DOM, no Svelte - pure math verification.
 */

import { describe, it, expect } from 'vitest';
import {
	calcularPosicionesCharolaEspaciada,
	calcularPosicionesCharolaTriangular,
	calcularPosicionesTuberia,
	calcularViewBox,
	calcularAnchoOcupadoCharola
} from './geometry.js';
import type { SistemaElectrico } from '$lib/features/calculos/domain/types/calculo.enums.js';

describe('calcularPosicionesCharolaEspaciada', () => {
	// Common test parameters
	const diametroFase = 21.2; // AWG 12 ~ 2.05mm², but let's use actual mm for cables
	const diametroTierra = 21.2;
	const anchoComercial = 150;

	describe('DELTA system (3-phase, no neutral)', () => {
		it('should return 4 conductors: A, B, C, T', () => {
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			expect(result).toHaveLength(4);

			const etiquetas = result.map((c) => c.etiqueta);
			expect(etiquetas).toContain('A');
			expect(etiquetas).toContain('B');
			expect(etiquetas).toContain('C');
			expect(etiquetas).toContain('T');
		});

		it('should space phase cables 1 diameter apart (center-to-center = 2*diametro)', () => {
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			// Find A, B, C positions
			const a = result.find((c) => c.etiqueta === 'A')!;
			const b = result.find((c) => c.etiqueta === 'B')!;
			const c = result.find((c) => c.etiqueta === 'C')!;

			// Distance A to B should be exactly 2 * diametro
			expect(b.cx - a.cx).toBeCloseTo(2 * diametroFase, 2);
			expect(c.cx - b.cx).toBeCloseTo(2 * diametroFase, 2);
		});

		it('should place tierra touching the last phase cable (no gap)', () => {
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			const c = result.find((c) => c.etiqueta === 'C')!;
			const t = result.find((c) => c.etiqueta === 'T')!;

			// T should touch C: distance between centers = radioC + radioT
			const distCT = t.cx - c.cx;
			const expectedDist = c.radio + t.radio;

			expect(distCT).toBeCloseTo(expectedDist, 2);
		});

		it('should center cables within the charola width', () => {
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			// Find min and max positions (including radius)
			const minCX = Math.min(...result.map((c) => c.cx - c.radio));
			const maxCX = Math.max(...result.map((c) => c.cx + c.radio));

			const totalWidth = maxCX - minCX;
			const centerOfCables = minCX + totalWidth / 2;
			const centerOfCharola = anchoComercial / 2;

			// Center of all cables should be at center of charola
			expect(centerOfCables).toBeCloseTo(centerOfCharola, 2);
		});
	});

	describe('ESTRELLA system (3-phase with neutral)', () => {
		it('should return 5 conductors: A, B, C, N, T', () => {
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'ESTRELLA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			expect(result).toHaveLength(5);

			const etiquetas = result.map((c) => c.etiqueta);
			expect(etiquetas).toContain('A');
			expect(etiquetas).toContain('B');
			expect(etiquetas).toContain('C');
			expect(etiquetas).toContain('N');
			expect(etiquetas).toContain('T');
		});
	});

	describe('MONOFASICO system (1-phase with neutral)', () => {
		it('should return 3 conductors: A, N, T', () => {
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'MONOFASICO' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			expect(result).toHaveLength(3);

			const etiquetas = result.map((c) => c.etiqueta);
			expect(etiquetas).toContain('A');
			expect(etiquetas).toContain('N');
			expect(etiquetas).toContain('T');
		});
	});

	describe('with control cable', () => {
		it('should add Ctrl cable after tierra with correct spacing', () => {
			const diametroControl = 15;
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				diametroControlMM: diametroControl,
				numHilosTotal: 1,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			const t = result.find((c) => c.etiqueta === 'T')!;
			const ctrl = result.find((c) => c.etiqueta === 'Ctrl')!;

			// Gap from tierra center to Ctrl center = radioTierra + diametroControl
			// (radioTierra is half of tierra diameter)
			const expectedGap = diametroTierra / 2 + diametroControl;
			expect(ctrl.cx - t.cx).toBeCloseTo(expectedGap, 2);
		});

		it('should handle multiple control cables', () => {
			const diametroControl = 15;
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				diametroControlMM: diametroControl,
				numHilosTotal: 2,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			const ctrl1 = result.find((c) => c.etiqueta === 'Ctrl1')!;
			const ctrl2 = result.find((c) => c.etiqueta === 'Ctrl2')!;

			// Control cables spaced 1 diameter apart
			expect(ctrl2.cx - ctrl1.cx).toBeCloseTo(2 * diametroControl, 2);
		});
	});
});

describe('calcularPosicionesCharolaTriangular', () => {
	const diametroFase = 21.2;
	const diametroTierra = 21.2;
	const anchoComercial = 150;
	const factorTriangular = 1.0;

	describe('DELTA system', () => {
		it('should place A and B at bottom, C at top (equilateral triangle)', () => {
			const result = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 1,
				factorTriangular,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			const a = result.find((c) => c.etiqueta === 'A')!;
			const b = result.find((c) => c.etiqueta === 'B')!;
			const c = result.find((c) => c.etiqueta === 'C')!;

			// A and B at cy = 0 (on floor)
			expect(a.cy).toBe(0);
			expect(b.cy).toBe(0);

			// C at equilateral triangle height
			const sin60 = Math.sqrt(3) / 2;
			const expectedHeight = diametroFase * sin60;
			expect(c.cy).toBeCloseTo(expectedHeight, 2);

			// C centered between A and B
			const expectedCX = a.cx + diametroFase / 2;
			expect(c.cx).toBeCloseTo(expectedCX, 2);
		});

		it('should space groups by factorTriangular × diametro', () => {
			const result = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 2,
				factorTriangular,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			// A1 and A2 should be separated by group step
			const a1 = result.find((c) => c.etiqueta === 'A1')!;
			const a2 = result.find((c) => c.etiqueta === 'A2')!;

			const groupWidth = 2 * diametroFase; // A and B
			const expectedStep = groupWidth + factorTriangular * diametroFase;

			expect(a2.cx - a1.cx).toBeCloseTo(expectedStep, 2);
		});
	});

	describe('ESTRELLA system', () => {
		it('should place A and B at bottom, C and N at top (diamond)', () => {
			const result = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 1,
				factorTriangular,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'ESTRELLA' as SistemaElectrico
			});

			const a = result.find((c) => c.etiqueta === 'A')!;
			const b = result.find((c) => c.etiqueta === 'B')!;
			const c = result.find((c) => c.etiqueta === 'C')!;
			const n = result.find((c) => c.etiqueta === 'N')!;

			// A and B at cy = 0
			expect(a.cy).toBe(0);
			expect(b.cy).toBe(0);

			// C and N at triangle height
			const sin60 = Math.sqrt(3) / 2;
			const expectedHeight = diametroFase * sin60;
			expect(c.cy).toBeCloseTo(expectedHeight, 2);
			expect(n.cy).toBeCloseTo(expectedHeight, 2);

			// C above A, N above B
			expect(c.cx).toBeCloseTo(a.cx, 2);
			expect(n.cx).toBeCloseTo(b.cx, 2);
		});
	});

	describe('MONOFASICO system', () => {
		it('should place A at bottom, N stacked on top', () => {
			const result = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 1,
				factorTriangular,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'MONOFASICO' as SistemaElectrico
			});

			const a = result.find((c) => c.etiqueta === 'A')!;
			const n = result.find((c) => c.etiqueta === 'N')!;

			// A at floor
			expect(a.cy).toBe(0);

			// N stacked directly above A (1 diameter height for stacked arrangement)
			expect(n.cx).toBeCloseTo(a.cx, 2);
			expect(n.cy).toBeCloseTo(diametroFase, 2); // Full diameter, not triangle height
		});
	});

	describe('tierra placement', () => {
		it('should place tierra touching the last bottom cable', () => {
			const result = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 1,
				factorTriangular,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			const b = result.find((c) => c.etiqueta === 'B')!;
			const t = result.find((c) => c.etiqueta === 'T')!;

			// T should touch B
			const distBT = t.cx - b.cx;
			const expectedDist = b.radio + t.radio;

			expect(distBT).toBeCloseTo(expectedDist, 2);
		});
	});
});

describe('calcularPosicionesTuberia', () => {
	describe('DELTA system', () => {
		it('should place 3 phases + tierra, all touching tube wall', () => {
			const diametroInterior = 50;
			const diametroExterior = 60;
			const R = diametroInterior / 2;
			const areaFase = Math.PI * 10 * 10; // radius 10mm
			const areaTierra = Math.PI * 10 * 10;

			const result = calcularPosicionesTuberia({
				diametroInteriorMM: diametroInterior,
				diametroExteriorMM: diametroExterior,
				areaFaseMM2: areaFase,
				areaTierraMM2: areaTierra,
				numFasesPorTubo: 3,
				numNeutrosPorTubo: 0,
				numTierras: 1,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			const rF = Math.sqrt(areaFase / Math.PI); // should be 10
			const rT = Math.sqrt(areaTierra / Math.PI); // should be 10

			const a = result.find((c) => c.etiqueta === 'A')!;
			const b = result.find((c) => c.etiqueta === 'B')!;
			const c = result.find((c) => c.etiqueta === 'C')!;
			const t = result.find((c) => c.etiqueta === 'T')!;
			// Note: N is on wall but not asserted in this test

			// Verify A and B are on tube wall: distance from center = R - rF
			const distA = Math.sqrt(a.cx * a.cx + a.cy * a.cy);
			const distB = Math.sqrt(b.cx * b.cx + b.cy * b.cy);
			expect(distA).toBeCloseTo(R - rF, 2);
			expect(distB).toBeCloseTo(R - rF, 2);

			// Distance between A and B should be 2*rF (touching)
			const distAB = Math.sqrt(Math.pow(b.cx - a.cx, 2) + Math.pow(b.cy - a.cy, 2));
			expect(distAB).toBeCloseTo(2 * rF, 2);

			// C should touch both A and B
			const distAC = Math.sqrt(Math.pow(c.cx - a.cx, 2) + Math.pow(c.cy - a.cy, 2));
			const distBC = Math.sqrt(Math.pow(c.cx - b.cx, 2) + Math.pow(c.cy - b.cy, 2));
			expect(distAC).toBeCloseTo(2 * rF, 2);
			expect(distBC).toBeCloseTo(2 * rF, 2);

			// T should be on tube wall
			const distT = Math.sqrt(t.cx * t.cx + t.cy * t.cy);
			expect(distT).toBeCloseTo(R - rT, 2);
		});
	});

	describe('MONOFASICO system', () => {
		it('should place A at bottom, N above A', () => {
			const diametroInterior = 50;
			const diametroExterior = 60;
			const areaFase = Math.PI * 10 * 10;
			const areaNeutro = Math.PI * 10 * 10;
			const areaTierra = Math.PI * 10 * 10;

			const result = calcularPosicionesTuberia({
				diametroInteriorMM: diametroInterior,
				diametroExteriorMM: diametroExterior,
				areaFaseMM2: areaFase,
				areaNeutroMM2: areaNeutro,
				areaTierraMM2: areaTierra,
				numFasesPorTubo: 1,
				numNeutrosPorTubo: 1,
				numTierras: 0,
				sistemaElectrico: 'MONOFASICO' as SistemaElectrico
			});

			const rF = 10;
			const rN = 10;

			const a = result.find((c) => c.etiqueta === 'A')!;
			const n = result.find((c) => c.etiqueta === 'N')!;

			// A at bottom center (cx ≈ 0, positive cy)
			expect(Math.abs(a.cx)).toBeLessThan(0.1);
			expect(a.cy).toBeGreaterThan(0);

			// N directly above A
			expect(Math.abs(n.cx - a.cx)).toBeLessThan(0.1);

			// N touches A: distance = rF + rN
			const distAN = Math.abs(n.cy - a.cy);
			expect(distAN).toBeCloseTo(rF + rN, 2);
		});
	});
});

describe('calcularViewBox', () => {
	it('should return correct viewBox string with default margin', () => {
		const result = calcularViewBox({
			anchoContenido: 100,
			altoContenido: 50
		});

		expect(result.viewBox).toBe('-20 -20 140 90');
		expect(result.ancho).toBe(140);
		expect(result.alto).toBe(90);
	});

	it('should return correct viewBox string with custom margin', () => {
		const result = calcularViewBox({
			anchoContenido: 100,
			altoContenido: 50,
			margen: 10
		});

		expect(result.viewBox).toBe('-10 -10 120 70');
		expect(result.ancho).toBe(120);
		expect(result.alto).toBe(70);
	});

	it('should handle zero content dimensions', () => {
		const result = calcularViewBox({
			anchoContenido: 0,
			altoContenido: 0
		});

		expect(result.viewBox).toBe('-20 -20 40 40');
	});
});

describe('calcularAnchoOcupadoCharola', () => {
	it('should return 0 for empty array', () => {
		const result = calcularAnchoOcupadoCharola([]);
		expect(result).toBe(0);
	});

	it('should return correct width for single conductor', () => {
		const posiciones = [
			{ cx: 50, cy: 0, radio: 10, color: '#fff', etiqueta: 'A', tipo: 'fase' as const }
		];

		const result = calcularAnchoOcupadoCharola(posiciones);
		expect(result).toBe(20); // radio * 2
	});

	it('should return correct width for multiple conductors', () => {
		const posiciones = [
			{ cx: 10, cy: 0, radio: 10, color: '#fff', etiqueta: 'A', tipo: 'fase' as const },
			{ cx: 50, cy: 0, radio: 10, color: '#fff', etiqueta: 'B', tipo: 'fase' as const }
		];

		const result = calcularAnchoOcupadoCharola(posiciones);
		// Min left edge = 10 - 10 = 0
		// Max right edge = 50 + 10 = 60
		// Width = 60 - 0 = 60
		expect(result).toBe(60);
	});
});

// ============================================================================
// Phase 5.3 — Integration Tests (cross-function)
// ============================================================================

describe('Phase 5.3: Integration Tests - calcularAnchoOcupadoCharola with position functions', () => {
	const diametroFase = 21.2;
	const diametroTierra = 21.2;
	const anchoComercial = 150;

	describe('Integration: calcularAnchoOcupadoCharola + calcularPosicionesCharolaEspaciada', () => {
		it('should calculate correct width for DELTA system output', () => {
			const posiciones = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			const anchoOcupado = calcularAnchoOcupadoCharola(posiciones);

			// DELTA: A, B, C (3 fases) + T (tierra)
			// Width = 4 conductors with spacing: first gap to last gap
			// Total conductors: 3 fases + 1 tierra = 4
			// Total width = (4 * diametro) + (3 * diametro) = 7 * diametro
			// Actually: 3 fases spaced 1 diameter apart + tierra touching last fase
			// 3 fases = 3 * 21.2 = 63.6mm (center to center spans 2*diametro for each pair)
			// A-B gap: 21.2, B-C gap: 21.2 = 42.4mm span between A center and C center
			// A left edge: cx - r = offset + r - r = offset (actually cx starts at offset + r)
			// Let's calculate: cx for A = offset + r, cx for B = offset + r + 2*diametro
			// cx for C = offset + r + 4*diametro, cx for T = cx(C) + r + rT
			// Width = T.cx + rT - (A.cx - r) = (offset + r + 4*diametro + r + rT) - (offset + r - r) = 4*diametro + 2*r + rT + r = 4*diametro + 3*r + rT
			// For same diameter: 4*21.2 + 3*10.6 + 10.6 = 84.8 + 31.8 + 10.6 = 127.2mm
			expect(anchoOcupado).toBeGreaterThan(0);
			expect(anchoOcupado).toBeCloseTo(127.2, 1);
		});

		it('should calculate correct width for ESTRELLA system output', () => {
			const posiciones = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'ESTRELLA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			const anchoOcupado = calcularAnchoOcupadoCharola(posiciones);

			// ESTRELLA: A, B, C (3 fases) + N (neutro) + T (tierra) = 5 conductors
			// A, B, C, N in a row with 2*diametro spacing
			// A at cx = r, B at cx = r+2d, C at cx = r+4d, N at cx = r+6d
			// T touches N: T.cx = N.cx + r + rT
			// Left edge = A.cx - r = offset
			// Right edge = T.cx + rT = (r+6d) + r + rT + rT = r + 6d + 2*rT = 3d + 6d + 2d = 11d
			// Wait, r = d/2, so: d/2 + 6d + d = 7.5d = 7.5 * 21.2 = 159mm (from left edge to T center)
			// Actually: Right edge = T.cx + rT = (r+6d+r+rT) + rT = r + 6d + r + 2*rT = 2r + 6d + 2rT = d + 6d + d = 8d = 169.6mm
			expect(anchoOcupado).toBeCloseTo(169.6, 1);
		});
	});

	describe('Integration: calcularAnchoOcupadoCharola + calcularPosicionesCharolaTriangular', () => {
		it('should calculate correct width for DELTA triangular arrangement', () => {
			const posiciones = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 1,
				factorTriangular: 1.0,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			const anchoOcupado = calcularAnchoOcupadoCharola(posiciones);

			// DELTA triangular: A, B at bottom, C at top, T touching B
			// Group width = 2*diametro (A to B spans 2*diametro)
			// T touches B: T.cx = B.cx + r + rT
			// Left edge = A.cx - r, Right edge = T.cx + rT
			// Width = T.cx + rT - (A.cx - r) = (cxB + r + rT) + rT - cxA + r
			// cxB - cxA = diametro, so: diametro + 2*r + 2*rT = diametro + diametro + diametro = 3*diametro = 63.6mm
			expect(anchoOcupado).toBeCloseTo(63.6, 1);
		});

		it('should calculate correct width for ESTRELLA triangular arrangement', () => {
			const posiciones = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 1,
				factorTriangular: 1.0,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'ESTRELLA' as SistemaElectrico
			});

			const anchoOcupado = calcularAnchoOcupadoCharola(posiciones);

			// ESTRELLA triangular: A, B at bottom, C, N at top, T touching B
			// Same width as DELTA since A-B-T same positions
			expect(anchoOcupado).toBeCloseTo(63.6, 1);
		});
	});
});

// ============================================================================
// Phase 5.4 — Edge Case Tests
// ============================================================================

describe('Phase 5.4: Edge Cases - calcularPosicionesCharolaEspaciada', () => {
	const diametroTierra = 21.2;
	const anchoComercial = 150;

	describe('Multi-hilo (hilosPorFase = 2)', () => {
		it('should double the number of phase cables with A1,A2,B1,B2 labels', () => {
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: 21.2,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 2,
				anchoComercialMM: anchoComercial
			});

			// Should have 2 * 3 fases + 1 tierra = 7 conductors
			expect(result).toHaveLength(7);

			const etiquetas = result.map((c) => c.etiqueta);
			expect(etiquetas).toContain('A1');
			expect(etiquetas).toContain('A2');
			expect(etiquetas).toContain('B1');
			expect(etiquetas).toContain('B2');
			expect(etiquetas).toContain('C1');
			expect(etiquetas).toContain('C2');
			expect(etiquetas).toContain('T');
		});

		it('should space groups correctly by 2*diametro', () => {
			const diametroFase = 21.2;
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				numHilosTotal: 0,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 2,
				anchoComercialMM: anchoComercial
			});

			const a1 = result.find((c) => c.etiqueta === 'A1')!;
			const a2 = result.find((c) => c.etiqueta === 'A2')!;

			// A2 should be 2*diametro away from A1 (group step)
			expect(a2.cx - a1.cx).toBeCloseTo(2 * diametroFase, 1);
		});
	});

	describe('Very small cable diameters', () => {
		it('should handle diametroFaseMM = 1.0 correctly', () => {
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: 1.0,
				diametroTierraMM: 1.0,
				numHilosTotal: 0,
				sistemaElectrico: 'DELTA' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: 100
			});

			expect(result).toHaveLength(4);

			// All cables should still be within charola
			const maxRight = Math.max(...result.map((c) => c.cx + c.radio));
			expect(maxRight).toBeLessThanOrEqual(100);
		});
	});

	describe('BIFASICO system', () => {
		it('should return 4 conductors: A, B, N, T', () => {
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: 21.2,
				diametroTierraMM: 21.2,
				numHilosTotal: 0,
				sistemaElectrico: 'BIFASICO' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			expect(result).toHaveLength(4);

			const etiquetas = result.map((c) => c.etiqueta);
			expect(etiquetas).toContain('A');
			expect(etiquetas).toContain('B');
			expect(etiquetas).toContain('N');
			expect(etiquetas).toContain('T');
		});

		it('should space A-B-N correctly with tierra touching N', () => {
			const diametroFase = 21.2;
			const result = calcularPosicionesCharolaEspaciada({
				diametroFaseMM: diametroFase,
				diametroTierraMM: 21.2,
				numHilosTotal: 0,
				sistemaElectrico: 'BIFASICO' as SistemaElectrico,
				hilosPorFase: 1,
				anchoComercialMM: anchoComercial
			});

			const a = result.find((c) => c.etiqueta === 'A')!;
			const b = result.find((c) => c.etiqueta === 'B')!;
			const n = result.find((c) => c.etiqueta === 'N')!;
			const t = result.find((c) => c.etiqueta === 'T')!;

			// A-B spacing: 2*diametro
			expect(b.cx - a.cx).toBeCloseTo(2 * diametroFase, 1);
			// B-N spacing: 2*diametro
			expect(n.cx - b.cx).toBeCloseTo(2 * diametroFase, 1);
			// T touching N: distance = rN + rT
			expect(t.cx - n.cx).toBeCloseTo(n.radio + t.radio, 1);
		});
	});
});

describe('Phase 5.4: Edge Cases - calcularPosicionesCharolaTriangular', () => {
	const diametroFase = 21.2;
	const diametroTierra = 21.2;
	const anchoComercial = 150;
	const factorTriangular = 1.0;

	describe('Multi-hilo (hilosPorFase = 2)', () => {
		it('should create 2 groups with correct spacing', () => {
			const result = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 2,
				factorTriangular,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			// 2 groups * 3 conductors (A, B, C) + 1 tierra = 7
			expect(result).toHaveLength(7);

			const a1 = result.find((c) => c.etiqueta === 'A1')!;
			const a2 = result.find((c) => c.etiqueta === 'A2')!;

			// Group step = groupWidth + factorTriangular * diametro
			// groupWidth = 2 * diametro (A to B)
			// Expected: 2*diametro + 1.0*diametro = 3*diametro
			const expectedStep = 2 * diametroFase + factorTriangular * diametroFase;
			expect(a2.cx - a1.cx).toBeCloseTo(expectedStep, 1);
		});

		it('should have correct labels: A1,A2,B1,B2,C1,C2,T', () => {
			const result = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 2,
				factorTriangular,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			const etiquetas = result.map((c) => c.etiqueta).sort();
			expect(etiquetas).toEqual(['A1', 'A2', 'B1', 'B2', 'C1', 'C2', 'T'].sort());
		});
	});

	describe('BIFASICO system in triangular arrangement', () => {
		it('should place A,B at bottom, N at top (like DELTA but N instead of C)', () => {
			const result = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 1,
				factorTriangular,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'BIFASICO' as SistemaElectrico
			});

			const a = result.find((c) => c.etiqueta === 'A')!;
			const b = result.find((c) => c.etiqueta === 'B')!;
			const n = result.find((c) => c.etiqueta === 'N')!;

			// A and B at bottom (cy = 0)
			expect(a.cy).toBe(0);
			expect(b.cy).toBe(0);

			// N at triangle height (equilateral triangle position)
			const sin60 = Math.sqrt(3) / 2;
			const expectedHeight = diametroFase * sin60;
			expect(n.cy).toBeCloseTo(expectedHeight, 1);

			// N centered between A and B
			const expectedCX = a.cx + diametroFase / 2;
			expect(n.cx).toBeCloseTo(expectedCX, 1);
		});

		it('should return 4 conductors: A, B, N, T', () => {
			const result = calcularPosicionesCharolaTriangular({
				diametroFaseMM: diametroFase,
				diametroTierraMM: diametroTierra,
				hilosPorFase: 1,
				factorTriangular,
				anchoComercialMM: anchoComercial,
				sistemaElectrico: 'BIFASICO' as SistemaElectrico
			});

			expect(result).toHaveLength(4);

			const etiquetas = result.map((c) => c.etiqueta);
			expect(etiquetas).toContain('A');
			expect(etiquetas).toContain('B');
			expect(etiquetas).toContain('N');
			expect(etiquetas).toContain('T');
		});
	});
});

describe('Phase 5.4: Edge Cases - calcularPosicionesTuberia', () => {
	describe('ESTRELLA system', () => {
		it('should return 5 positions: A, B, C, N, T', () => {
			const diametroInterior = 50;
			const diametroExterior = 60;
			const areaFase = Math.PI * 10 * 10;
			const areaNeutro = Math.PI * 10 * 10;
			const areaTierra = Math.PI * 10 * 10;

			const result = calcularPosicionesTuberia({
				diametroInteriorMM: diametroInterior,
				diametroExteriorMM: diametroExterior,
				areaFaseMM2: areaFase,
				areaNeutroMM2: areaNeutro,
				areaTierraMM2: areaTierra,
				numFasesPorTubo: 3,
				numNeutrosPorTubo: 1,
				numTierras: 1,
				sistemaElectrico: 'ESTRELLA' as SistemaElectrico
			});

			expect(result).toHaveLength(5);

			const etiquetas = result.map((c) => c.etiqueta);
			expect(etiquetas).toContain('A');
			expect(etiquetas).toContain('B');
			expect(etiquetas).toContain('C');
			expect(etiquetas).toContain('N');
			expect(etiquetas).toContain('T');
		});

		it('should place A,B on bottom, C above A, N above B, T touching B on wall', () => {
			const diametroInterior = 50;
			const R = diametroInterior / 2;
			const rF = 10;
			const rT = 10;

			const result = calcularPosicionesTuberia({
				diametroInteriorMM: diametroInterior,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * rF * rF,
				areaNeutroMM2: Math.PI * rF * rF,
				areaTierraMM2: Math.PI * rT * rT,
				numFasesPorTubo: 3,
				numNeutrosPorTubo: 1,
				numTierras: 1,
				sistemaElectrico: 'ESTRELLA' as SistemaElectrico
			});

			const a = result.find((c) => c.etiqueta === 'A')!;
			const b = result.find((c) => c.etiqueta === 'B')!;
			const c = result.find((c) => c.etiqueta === 'C')!;
			const n = result.find((c) => c.etiqueta === 'N')!;
			// Note: T is on wall but not asserted in this test

			// A and B on wall
			const distA = Math.sqrt(a.cx * a.cx + a.cy * a.cy);
			const distB = Math.sqrt(b.cx * b.cx + b.cy * b.cy);
			expect(distA).toBeCloseTo(R - rF, 1);
			expect(distB).toBeCloseTo(R - rF, 1);

			// C above A (touching)
			const distAC = Math.sqrt(Math.pow(c.cx - a.cx, 2) + Math.pow(c.cy - a.cy, 2));
			expect(distAC).toBeCloseTo(2 * rF, 1);

			// N above B (touching)
			const distBN = Math.sqrt(Math.pow(n.cx - b.cx, 2) + Math.pow(n.cy - b.cy, 2));
			expect(distBN).toBeCloseTo(2 * rF, 1);
		});
	});

	describe('BIFASICO system', () => {
		it('should return 4 positions: A, B, N, T', () => {
			const diametroInterior = 50;
			const diametroExterior = 60;
			const areaFase = Math.PI * 10 * 10;
			const areaNeutro = Math.PI * 10 * 10;
			const areaTierra = Math.PI * 10 * 10;

			const result = calcularPosicionesTuberia({
				diametroInteriorMM: diametroInterior,
				diametroExteriorMM: diametroExterior,
				areaFaseMM2: areaFase,
				areaNeutroMM2: areaNeutro,
				areaTierraMM2: areaTierra,
				numFasesPorTubo: 2,
				numNeutrosPorTubo: 1,
				numTierras: 1,
				sistemaElectrico: 'BIFASICO' as SistemaElectrico
			});

			expect(result).toHaveLength(4);

			const etiquetas = result.map((c) => c.etiqueta);
			expect(etiquetas).toContain('A');
			expect(etiquetas).toContain('B');
			expect(etiquetas).toContain('N');
			expect(etiquetas).toContain('T');
		});

		it('should place A,B at bottom touching, N on top touching both, T on wall touching B', () => {
			const diametroInterior = 50;
			const rF = 10;
			const rN = 10;
			const rT = 10;

			const result = calcularPosicionesTuberia({
				diametroInteriorMM: diametroInterior,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * rF * rF,
				areaNeutroMM2: Math.PI * rN * rN,
				areaTierraMM2: Math.PI * rT * rT,
				numFasesPorTubo: 2,
				numNeutrosPorTubo: 1,
				numTierras: 1,
				sistemaElectrico: 'BIFASICO' as SistemaElectrico
			});

			const a = result.find((c) => c.etiqueta === 'A')!;
			const b = result.find((c) => c.etiqueta === 'B')!;
			const n = result.find((c) => c.etiqueta === 'N')!;
			// Note: T is on wall but not asserted in this test

			// A-B touching: distance = 2*rF
			const distAB = Math.sqrt(Math.pow(b.cx - a.cx, 2) + Math.pow(b.cy - a.cy, 2));
			expect(distAB).toBeCloseTo(2 * rF, 1);

			// N touches A and B
			const distAN = Math.sqrt(Math.pow(n.cx - a.cx, 2) + Math.pow(n.cy - a.cy, 2));
			const distBN = Math.sqrt(Math.pow(n.cx - b.cx, 2) + Math.pow(n.cy - b.cy, 2));
			expect(distAN).toBeCloseTo(rF + rN, 1);
			expect(distBN).toBeCloseTo(rF + rN, 1);
		});
	});

	describe('Empty tube', () => {
		it('should return empty array for 0 phases, 0 neutros, 0 tierras', () => {
			const result = calcularPosicionesTuberia({
				diametroInteriorMM: 50,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * 10 * 10,
				areaNeutroMM2: Math.PI * 10 * 10,
				areaTierraMM2: Math.PI * 10 * 10,
				numFasesPorTubo: 0,
				numNeutrosPorTubo: 0,
				numTierras: 0,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			expect(result).toHaveLength(0);
		});
	});

	describe('All cables within tube validation', () => {
		it('should keep all conductors inside the tube for DELTA', () => {
			const diametroInterior = 50;
			const R = diametroInterior / 2;

			const result = calcularPosicionesTuberia({
				diametroInteriorMM: diametroInterior,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * 10 * 10,
				areaTierraMM2: Math.PI * 10 * 10,
				numFasesPorTubo: 3,
				numNeutrosPorTubo: 0,
				numTierras: 1,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			// For every conductor: distance from tube center to cable edge <= R
			// i.e., sqrt(cx² + cy²) + radio <= R
			for (const conductor of result) {
				const distFromCenter = Math.sqrt(conductor.cx * conductor.cx + conductor.cy * conductor.cy);
				const cableEdgeDist = distFromCenter + conductor.radio;
				expect(cableEdgeDist).toBeLessThanOrEqual(R + 0.01); // small tolerance for float
			}
		});

		it('should keep all conductors inside the tube for ESTRELLA', () => {
			const diametroInterior = 50;
			const R = diametroInterior / 2;

			const result = calcularPosicionesTuberia({
				diametroInteriorMM: diametroInterior,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * 10 * 10,
				areaNeutroMM2: Math.PI * 10 * 10,
				areaTierraMM2: Math.PI * 10 * 10,
				numFasesPorTubo: 3,
				numNeutrosPorTubo: 1,
				numTierras: 1,
				sistemaElectrico: 'ESTRELLA' as SistemaElectrico
			});

			for (const conductor of result) {
				const distFromCenter = Math.sqrt(conductor.cx * conductor.cx + conductor.cy * conductor.cy);
				const cableEdgeDist = distFromCenter + conductor.radio;
				expect(cableEdgeDist).toBeLessThanOrEqual(R + 0.01);
			}
		});

		it('should keep all conductors inside the tube for BIFASICO', () => {
			const diametroInterior = 50;
			const R = diametroInterior / 2;

			const result = calcularPosicionesTuberia({
				diametroInteriorMM: diametroInterior,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * 10 * 10,
				areaNeutroMM2: Math.PI * 10 * 10,
				areaTierraMM2: Math.PI * 10 * 10,
				numFasesPorTubo: 2,
				numNeutrosPorTubo: 1,
				numTierras: 1,
				sistemaElectrico: 'BIFASICO' as SistemaElectrico
			});

			for (const conductor of result) {
				const distFromCenter = Math.sqrt(conductor.cx * conductor.cx + conductor.cy * conductor.cy);
				const cableEdgeDist = distFromCenter + conductor.radio;
				expect(cableEdgeDist).toBeLessThanOrEqual(R + 0.01);
			}
		});
	});

	describe('DELTA touching validation', () => {
		it('should verify A-B distance ≈ 2*rF (touching)', () => {
			const rF = 10;
			const result = calcularPosicionesTuberia({
				diametroInteriorMM: 50,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * rF * rF,
				areaTierraMM2: Math.PI * rF * rF,
				numFasesPorTubo: 3,
				numNeutrosPorTubo: 0,
				numTierras: 1,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			const a = result.find((c) => c.etiqueta === 'A')!;
			const b = result.find((c) => c.etiqueta === 'B')!;

			const distAB = Math.sqrt(Math.pow(b.cx - a.cx, 2) + Math.pow(b.cy - a.cy, 2));
			expect(distAB).toBeCloseTo(2 * rF, 1);
		});

		it('should verify A-C distance ≈ 2*rF (touching)', () => {
			const rF = 10;
			const result = calcularPosicionesTuberia({
				diametroInteriorMM: 50,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * rF * rF,
				areaTierraMM2: Math.PI * rF * rF,
				numFasesPorTubo: 3,
				numNeutrosPorTubo: 0,
				numTierras: 1,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			const a = result.find((c) => c.etiqueta === 'A')!;
			const c = result.find((c) => c.etiqueta === 'C')!;

			const distAC = Math.sqrt(Math.pow(c.cx - a.cx, 2) + Math.pow(c.cy - a.cy, 2));
			expect(distAC).toBeCloseTo(2 * rF, 1);
		});

		it('should verify B-C distance ≈ 2*rF (touching)', () => {
			const rF = 10;
			const result = calcularPosicionesTuberia({
				diametroInteriorMM: 50,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * rF * rF,
				areaTierraMM2: Math.PI * rF * rF,
				numFasesPorTubo: 3,
				numNeutrosPorTubo: 0,
				numTierras: 1,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			const b = result.find((c) => c.etiqueta === 'B')!;
			const c = result.find((c) => c.etiqueta === 'C')!;

			const distBC = Math.sqrt(Math.pow(c.cx - b.cx, 2) + Math.pow(c.cy - b.cy, 2));
			expect(distBC).toBeCloseTo(2 * rF, 1);
		});

		it('should verify T touches B (distance between centers ≈ rF + rT)', () => {
			const rF = 10;
			const rT = 10;
			const result = calcularPosicionesTuberia({
				diametroInteriorMM: 50,
				diametroExteriorMM: 60,
				areaFaseMM2: Math.PI * rF * rF,
				areaTierraMM2: Math.PI * rT * rT,
				numFasesPorTubo: 3,
				numNeutrosPorTubo: 0,
				numTierras: 1,
				sistemaElectrico: 'DELTA' as SistemaElectrico
			});

			const b = result.find((c) => c.etiqueta === 'B')!;
			const t = result.find((c) => c.etiqueta === 'T')!;

			const distBT = Math.sqrt(Math.pow(t.cx - b.cx, 2) + Math.pow(t.cy - b.cy, 2));
			expect(distBT).toBeCloseTo(rF + rT, 1);
		});
	});
});
