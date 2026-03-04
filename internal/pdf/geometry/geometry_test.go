// internal/pdf/geometry/geometry_test.go
package geometry

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to find conductor by etiqueta
func findConductor(posiciones []ConductorPosicion, etiqueta string) *ConductorPosicion {
	for i := range posiciones {
		if posiciones[i].Etiqueta == etiqueta {
			return &posiciones[i]
		}
	}
	return nil
}

// Helper to calculate distance between two points
func distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// Helper to calculate distance from origin
func distFromOrigin(cx, cy float64) float64 {
	return math.Sqrt(cx*cx + cy*cy)
}

// ============================================================================
// Tests: CalcularPosicionesCharolaEspaciada
// ============================================================================

func TestCalcularPosicionesCharolaEspaciada_DELTA(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoDelta,
		HilosPorFase:     1,
		AnchoComercialMM: anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	// Should return 4 conductors: A, B, C, T
	require.Len(t, result, 4, "DELTA should have 4 conductors")

	etiquetas := make([]string, len(result))
	for i, c := range result {
		etiquetas[i] = c.Etiqueta
	}
	assert.Contains(t, etiquetas, "A")
	assert.Contains(t, etiquetas, "B")
	assert.Contains(t, etiquetas, "C")
	assert.Contains(t, etiquetas, "T")
}

func TestCalcularPosicionesCharolaEspaciada_DELTA_Spacing(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoDelta,
		HilosPorFase:     1,
		AnchoComercialMM: anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	a := findConductor(result, "A")
	b := findConductor(result, "B")
	c := findConductor(result, "C")

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.NotNil(t, c)

	// Distance A to B should be exactly 2 * diametro
	assert.InDelta(t, b.CX-a.CX, 2*diametroFase, 0.1, "A-B spacing should be 2*diametro")
	assert.InDelta(t, c.CX-b.CX, 2*diametroFase, 0.1, "B-C spacing should be 2*diametro")
}

func TestCalcularPosicionesCharolaEspaciada_DELTA_TierraTouching(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoDelta,
		HilosPorFase:     1,
		AnchoComercialMM: anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	c := findConductor(result, "C")
	tierra := findConductor(result, "T")

	require.NotNil(t, c)
	require.NotNil(t, tierra)

	// T should touch C: distance between centers = radioC + radioT
	distCT := tierra.CX - c.CX
	expectedDist := c.Radio + tierra.Radio

	assert.InDelta(t, distCT, expectedDist, 0.1, "T should touch C")
}

func TestCalcularPosicionesCharolaEspaciada_DELTA_Centered(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoDelta,
		HilosPorFase:     1,
		AnchoComercialMM: anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	// Find min and max positions (including radius)
	minCX := math.MaxFloat64
	maxCX := -math.MaxFloat64
	for _, c := range result {
		if c.CX-c.Radio < minCX {
			minCX = c.CX - c.Radio
		}
		if c.CX+c.Radio > maxCX {
			maxCX = c.CX + c.Radio
		}
	}

	totalWidth := maxCX - minCX
	centerOfCables := minCX + totalWidth/2
	centerOfCharola := anchoComercial / 2

	assert.InDelta(t, centerOfCables, centerOfCharola, 0.1, "Cables should be centered in charola")
}

func TestCalcularPosicionesCharolaEspaciada_ESTRELLA(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoEstrella,
		HilosPorFase:     1,
		AnchoComercialMM: anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	// Should return 5 conductors: A, B, C, N, T
	require.Len(t, result, 5, "ESTRELLA should have 5 conductors")

	etiquetas := make([]string, len(result))
	for i, c := range result {
		etiquetas[i] = c.Etiqueta
	}
	assert.Contains(t, etiquetas, "A")
	assert.Contains(t, etiquetas, "B")
	assert.Contains(t, etiquetas, "C")
	assert.Contains(t, etiquetas, "N")
	assert.Contains(t, etiquetas, "T")
}

func TestCalcularPosicionesCharolaEspaciada_MONOFASICO(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoMonofasico,
		HilosPorFase:     1,
		AnchoComercialMM: anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	// Should return 3 conductors: A, N, T
	require.Len(t, result, 3, "MONOFASICO should have 3 conductors")

	etiquetas := make([]string, len(result))
	for i, c := range result {
		etiquetas[i] = c.Etiqueta
	}
	assert.Contains(t, etiquetas, "A")
	assert.Contains(t, etiquetas, "N")
	assert.Contains(t, etiquetas, "T")
}

func TestCalcularPosicionesCharolaEspaciada_BIFASICO(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoBifasico,
		HilosPorFase:     1,
		AnchoComercialMM: anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	// Should return 4 conductors: A, B, N, T
	require.Len(t, result, 4, "BIFASICO should have 4 conductors")

	etiquetas := make([]string, len(result))
	for i, c := range result {
		etiquetas[i] = c.Etiqueta
	}
	assert.Contains(t, etiquetas, "A")
	assert.Contains(t, etiquetas, "B")
	assert.Contains(t, etiquetas, "N")
	assert.Contains(t, etiquetas, "T")
}

func TestCalcularPosicionesCharolaEspaciada_BIFASICO_Spacing(t *testing.T) {
	diametroFase := 21.2
	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: 21.2,
		SistemaElectrico: SistemaElectricoBifasico,
		HilosPorFase:     1,
		AnchoComercialMM: 150.0,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	a := findConductor(result, "A")
	b := findConductor(result, "B")
	n := findConductor(result, "N")
	tierra := findConductor(result, "T")

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.NotNil(t, n)
	require.NotNil(t, tierra)

	// A-B spacing: 2*diametro
	assert.InDelta(t, b.CX-a.CX, 2*diametroFase, 0.1, "A-B spacing should be 2*diametro")
	// B-N spacing: 2*diametro
	assert.InDelta(t, n.CX-b.CX, 2*diametroFase, 0.1, "B-N spacing should be 2*diametro")
	// T touching N: distance = rN + rT
	assert.InDelta(t, tierra.CX-n.CX, n.Radio+tierra.Radio, 0.1, "T should touch N")
}

func TestCalcularPosicionesCharolaEspaciada_WithControlCable(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	diametroControl := 15.0
	anchoComercial := 150.0

	diametroControlParam := diametroControl

	params := ParametrosCharolaBase{
		DiametroFaseMM:    diametroFase,
		DiametroTierraMM:  diametroTierra,
		DiametroControlMM: &diametroControlParam,
		NumHilosControl:   1,
		SistemaElectrico:  SistemaElectricoDelta,
		HilosPorFase:      1,
		AnchoComercialMM:  anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	tierra := findConductor(result, "T")
	ctrl := findConductor(result, "Ctrl")

	require.NotNil(t, tierra)
	require.NotNil(t, ctrl)

	// Gap from tierra center to Ctrl center = radioTierra + diametroControl
	expectedGap := diametroTierra/2 + diametroControl
	assert.InDelta(t, ctrl.CX-tierra.CX, expectedGap, 0.1, "Control cable spacing")
}

func TestCalcularPosicionesCharolaEspaciada_MultipleControlCables(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	diametroControl := 15.0
	anchoComercial := 150.0

	diametroControlParam := diametroControl

	params := ParametrosCharolaBase{
		DiametroFaseMM:    diametroFase,
		DiametroTierraMM:  diametroTierra,
		DiametroControlMM: &diametroControlParam,
		NumHilosControl:   2,
		SistemaElectrico:  SistemaElectricoDelta,
		HilosPorFase:      1,
		AnchoComercialMM:  anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	ctrl1 := findConductor(result, "Ctrl1")
	ctrl2 := findConductor(result, "Ctrl2")

	require.NotNil(t, ctrl1)
	require.NotNil(t, ctrl2)

	// Control cables spaced 1 diameter apart
	assert.InDelta(t, ctrl2.CX-ctrl1.CX, 2*diametroControl, 0.1, "Control cables spacing")
}

func TestCalcularPosicionesCharolaEspaciada_MultiHilo(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoDelta,
		HilosPorFase:     2,
		AnchoComercialMM: anchoComercial,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	// Should have 2 * 3 fases + 1 tierra = 7 conductors
	require.Len(t, result, 7, "DELTA with 2 hilos should have 7 conductors")

	etiquetas := make([]string, len(result))
	for i, c := range result {
		etiquetas[i] = c.Etiqueta
	}
	assert.Contains(t, etiquetas, "A1")
	assert.Contains(t, etiquetas, "A2")
	assert.Contains(t, etiquetas, "B1")
	assert.Contains(t, etiquetas, "B2")
	assert.Contains(t, etiquetas, "C1")
	assert.Contains(t, etiquetas, "C2")
	assert.Contains(t, etiquetas, "T")
}

func TestCalcularPosicionesCharolaEspaciada_MultiHilo_Spacing(t *testing.T) {
	diametroFase := 21.2
	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: 21.2,
		SistemaElectrico: SistemaElectricoDelta,
		HilosPorFase:     2,
		AnchoComercialMM: 150.0,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	a1 := findConductor(result, "A1")
	a2 := findConductor(result, "A2")

	require.NotNil(t, a1)
	require.NotNil(t, a2)

	// A2 should be 2*diametro away from A1 (group step)
	assert.InDelta(t, a2.CX-a1.CX, 2*diametroFase, 0.1, "Multi-hilo spacing")
}

func TestCalcularPosicionesCharolaEspaciada_VerySmallDiameter(t *testing.T) {
	params := ParametrosCharolaBase{
		DiametroFaseMM:   1.0,
		DiametroTierraMM: 1.0,
		SistemaElectrico: SistemaElectricoDelta,
		HilosPorFase:     1,
		AnchoComercialMM: 100.0,
	}

	result := CalcularPosicionesCharolaEspaciada(params)

	require.Len(t, result, 4)

	// All cables should still be within charola
	maxRight := -math.MaxFloat64
	for _, c := range result {
		if c.CX+c.Radio > maxRight {
			maxRight = c.CX + c.Radio
		}
	}
	assert.LessOrEqual(t, maxRight, 100.0, "Cables should fit in charola")
}

// ============================================================================
// Tests: CalcularPosicionesCharolaTriangular
// ============================================================================

func TestCalcularPosicionesCharolaTriangular_DELTA(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     1,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoDelta,
	}

	result := CalcularPosicionesCharolaTriangular(params)

	a := findConductor(result, "A")
	b := findConductor(result, "B")
	c := findConductor(result, "C")

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.NotNil(t, c)

	// A and B at cy = 0 (on floor)
	assert.Equal(t, 0.0, a.CY, "A should be on floor")
	assert.Equal(t, 0.0, b.CY, "B should be on floor")

	// C at equilateral triangle height
	sin60 := math.Sqrt(3) / 2
	expectedHeight := diametroFase * sin60
	assert.InDelta(t, c.CY, expectedHeight, 0.1, "C should be at triangle height")

	// C centered between A and B
	expectedCX := a.CX + diametroFase/2
	assert.InDelta(t, c.CX, expectedCX, 0.1, "C should be centered between A and B")
}

func TestCalcularPosicionesCharolaTriangular_DELTA_GroupSpacing(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     2,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoDelta,
	}

	result := CalcularPosicionesCharolaTriangular(params)

	a1 := findConductor(result, "A1")
	a2 := findConductor(result, "A2")

	require.NotNil(t, a1)
	require.NotNil(t, a2)

	groupWidth := 2 * diametroFase // A and B
	expectedStep := groupWidth + 1.0*diametroFase

	assert.InDelta(t, a2.CX-a1.CX, expectedStep, 0.1, "Group spacing")
}

func TestCalcularPosicionesCharolaTriangular_ESTRELLA(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     1,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoEstrella,
	}

	result := CalcularPosicionesCharolaTriangular(params)

	a := findConductor(result, "A")
	b := findConductor(result, "B")
	c := findConductor(result, "C")
	n := findConductor(result, "N")

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.NotNil(t, c)
	require.NotNil(t, n)

	// A and B at cy = 0
	assert.Equal(t, 0.0, a.CY, "A should be on floor")
	assert.Equal(t, 0.0, b.CY, "B should be on floor")

	// C and N at triangle height
	sin60 := math.Sqrt(3) / 2
	expectedHeight := diametroFase * sin60
	assert.InDelta(t, c.CY, expectedHeight, 0.1, "C should be at triangle height")
	assert.InDelta(t, n.CY, expectedHeight, 0.1, "N should be at triangle height")

	// C above A, N above B
	assert.InDelta(t, c.CX, a.CX, 0.1, "C should be above A")
	assert.InDelta(t, n.CX, b.CX, 0.1, "N should be above B")
}

func TestCalcularPosicionesCharolaTriangular_MONOFASICO(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     1,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoMonofasico,
	}

	result := CalcularPosicionesCharolaTriangular(params)

	a := findConductor(result, "A")
	n := findConductor(result, "N")

	require.NotNil(t, a)
	require.NotNil(t, n)

	// A at floor
	assert.Equal(t, 0.0, a.CY, "A should be on floor")

	// N stacked directly above A (1 diameter height for stacked arrangement)
	assert.InDelta(t, n.CX, a.CX, 0.1, "N should be above A")
	assert.InDelta(t, n.CY, diametroFase, 0.1, "N should be at full diameter height")
}

func TestCalcularPosicionesCharolaTriangular_TierraPlacement(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     1,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoDelta,
	}

	result := CalcularPosicionesCharolaTriangular(params)

	b := findConductor(result, "B")
	tierra := findConductor(result, "T")

	require.NotNil(t, b)
	require.NotNil(t, tierra)

	// T should touch B
	distBT := tierra.CX - b.CX
	expectedDist := b.Radio + tierra.Radio

	assert.InDelta(t, distBT, expectedDist, 0.1, "T should touch B")
}

func TestCalcularPosicionesCharolaTriangular_BIFASICO(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     1,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoBifasico,
	}

	result := CalcularPosicionesCharolaTriangular(params)

	// Should return 4 conductors: A, B, N, T
	require.Len(t, result, 4, "BIFASICO should have 4 conductors")

	etiquetas := make([]string, len(result))
	for i, c := range result {
		etiquetas[i] = c.Etiqueta
	}
	assert.Contains(t, etiquetas, "A")
	assert.Contains(t, etiquetas, "B")
	assert.Contains(t, etiquetas, "N")
	assert.Contains(t, etiquetas, "T")
}

func TestCalcularPosicionesCharolaTriangular_BIFASICO_Positions(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     1,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoBifasico,
	}

	result := CalcularPosicionesCharolaTriangular(params)

	a := findConductor(result, "A")
	b := findConductor(result, "B")
	n := findConductor(result, "N")

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.NotNil(t, n)

	// A and B at bottom (cy = 0)
	assert.Equal(t, 0.0, a.CY, "A should be on floor")
	assert.Equal(t, 0.0, b.CY, "B should be on floor")

	// N at triangle height (equilateral triangle position)
	sin60 := math.Sqrt(3) / 2
	expectedHeight := diametroFase * sin60
	assert.InDelta(t, n.CY, expectedHeight, 0.1, "N should be at triangle height")

	// N centered between A and B
	expectedCX := a.CX + diametroFase/2
	assert.InDelta(t, n.CX, expectedCX, 0.1, "N should be centered between A and B")
}

func TestCalcularPosicionesCharolaTriangular_MultiHilo(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     2,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoDelta,
	}

	result := CalcularPosicionesCharolaTriangular(params)

	// 2 groups * 3 conductors (A, B, C) + 1 tierra = 7
	require.Len(t, result, 7)

	etiquetas := make([]string, len(result))
	for i, c := range result {
		etiquetas[i] = c.Etiqueta
	}
	assert.Contains(t, etiquetas, "A1")
	assert.Contains(t, etiquetas, "A2")
	assert.Contains(t, etiquetas, "B1")
	assert.Contains(t, etiquetas, "B2")
	assert.Contains(t, etiquetas, "C1")
	assert.Contains(t, etiquetas, "C2")
	assert.Contains(t, etiquetas, "T")
}

// ============================================================================
// Tests: CalcularPosicionesTuberia
// ============================================================================

func TestCalcularPosicionesTuberia_DELTA(t *testing.T) {
	diametroInterior := 50.0
	diametroExterior := 60.0
	R := diametroInterior / 2
	radioFase := 10.0
	radioTierra := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaTierra := math.Pi * radioTierra * radioTierra

	params := ParametrosTuberia{
		DiametroInteriorMM: diametroInterior,
		DiametroExteriorMM: diametroExterior,
		AreaFaseMM2:        areaFase,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    3,
		NumNeutrosPorTubo:  0,
		NumTierras:         1,
		SistemaElectrico:   SistemaElectricoDelta,
	}

	result := CalcularPosicionesTuberia(params)

	a := findConductor(result, "A")
	b := findConductor(result, "B")
	c := findConductor(result, "C")
	tierra := findConductor(result, "T")

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.NotNil(t, c)
	require.NotNil(t, tierra)

	// Verify A and B are on tube wall: distance from center = R - rF
	distA := distFromOrigin(a.CX, a.CY)
	distB := distFromOrigin(b.CX, b.CY)
	assert.InDelta(t, distA, R-radioFase, 0.1, "A should be on wall")
	assert.InDelta(t, distB, R-radioFase, 0.1, "B should be on wall")

	// Distance between A and B should be 2*rF (touching)
	distAB := distance(a.CX, a.CY, b.CX, b.CY)
	assert.InDelta(t, distAB, 2*radioFase, 0.1, "A and B should touch")

	// C should touch both A and B
	distAC := distance(a.CX, a.CY, c.CX, c.CY)
	distBC := distance(b.CX, b.CY, c.CX, c.CY)
	assert.InDelta(t, distAC, 2*radioFase, 0.1, "C should touch A")
	assert.InDelta(t, distBC, 2*radioFase, 0.1, "C should touch B")

	// T should be on tube wall
	distT := distFromOrigin(tierra.CX, tierra.CY)
	assert.InDelta(t, distT, R-radioTierra, 0.1, "T should be on wall")
}

func TestCalcularPosicionesTuberia_MONOFASICO(t *testing.T) {
	diametroInterior := 50.0
	diametroExterior := 60.0
	radioFase := 10.0
	radioN := 10.0
	radioTierra := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaNeutro := math.Pi * radioN * radioN
	areaTierra := math.Pi * radioTierra * radioTierra

	params := ParametrosTuberia{
		DiametroInteriorMM: diametroInterior,
		DiametroExteriorMM: diametroExterior,
		AreaFaseMM2:        areaFase,
		AreaNeutroMM2:      &areaNeutro,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    1,
		NumNeutrosPorTubo:  1,
		NumTierras:         0,
		SistemaElectrico:   SistemaElectricoMonofasico,
	}

	result := CalcularPosicionesTuberia(params)

	a := findConductor(result, "A")
	n := findConductor(result, "N")

	require.NotNil(t, a)
	require.NotNil(t, n)

	// A at bottom center (cx ≈ 0, positive cy)
	assert.True(t, math.Abs(a.CX) < 0.1, "A should be centered")
	assert.Greater(t, a.CY, 0.0, "A should be below center")

	// N directly above A
	assert.InDelta(t, n.CX, a.CX, 0.1, "N should be above A")

	// N touches A: distance = rF + rN
	distAN := math.Abs(n.CY - a.CY)
	assert.InDelta(t, distAN, radioFase+radioN, 0.1, "N should touch A")
}

func TestCalcularPosicionesTuberia_ESTRELLA(t *testing.T) {
	diametroInterior := 50.0
	radioFase := 10.0
	radioTierra := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaNeutro := math.Pi * radioFase * radioFase
	areaTierra := math.Pi * radioTierra * radioTierra

	params := ParametrosTuberia{
		DiametroInteriorMM: diametroInterior,
		DiametroExteriorMM: 60.0,
		AreaFaseMM2:        areaFase,
		AreaNeutroMM2:      &areaNeutro,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    3,
		NumNeutrosPorTubo:  1,
		NumTierras:         1,
		SistemaElectrico:   SistemaElectricoEstrella,
	}

	result := CalcularPosicionesTuberia(params)

	// Should return 5 positions: A, B, C, N, T
	require.Len(t, result, 5)

	etiquetas := make([]string, len(result))
	for i, c := range result {
		etiquetas[i] = c.Etiqueta
	}
	assert.Contains(t, etiquetas, "A")
	assert.Contains(t, etiquetas, "B")
	assert.Contains(t, etiquetas, "C")
	assert.Contains(t, etiquetas, "N")
	assert.Contains(t, etiquetas, "T")
}

func TestCalcularPosicionesTuberia_ESTRELLA_Positions(t *testing.T) {
	diametroInterior := 50.0
	R := diametroInterior / 2
	radioFase := 10.0
	radioTierra := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaNeutro := math.Pi * radioFase * radioFase
	areaTierra := math.Pi * radioTierra * radioTierra

	params := ParametrosTuberia{
		DiametroInteriorMM: diametroInterior,
		DiametroExteriorMM: 60.0,
		AreaFaseMM2:        areaFase,
		AreaNeutroMM2:      &areaNeutro,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    3,
		NumNeutrosPorTubo:  1,
		NumTierras:         1,
		SistemaElectrico:   SistemaElectricoEstrella,
	}

	result := CalcularPosicionesTuberia(params)

	a := findConductor(result, "A")
	b := findConductor(result, "B")
	c := findConductor(result, "C")
	n := findConductor(result, "N")

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.NotNil(t, c)
	require.NotNil(t, n)

	// A and B on wall
	distA := distFromOrigin(a.CX, a.CY)
	distB := distFromOrigin(b.CX, b.CY)
	assert.InDelta(t, distA, R-radioFase, 0.1, "A should be on wall")
	assert.InDelta(t, distB, R-radioFase, 0.1, "B should be on wall")

	// C above A (touching)
	distAC := distance(a.CX, a.CY, c.CX, c.CY)
	assert.InDelta(t, distAC, 2*radioFase, 0.1, "C should touch A")

	// N above B (touching)
	distBN := distance(b.CX, b.CY, n.CX, n.CY)
	assert.InDelta(t, distBN, 2*radioFase, 0.1, "N should touch B")
}

func TestCalcularPosicionesTuberia_BIFASICO(t *testing.T) {
	diametroInterior := 50.0
	radioFase := 10.0
	radioN := 10.0
	radioTierra := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaNeutro := math.Pi * radioN * radioN
	areaTierra := math.Pi * radioTierra * radioTierra

	params := ParametrosTuberia{
		DiametroInteriorMM: diametroInterior,
		DiametroExteriorMM: 60.0,
		AreaFaseMM2:        areaFase,
		AreaNeutroMM2:      &areaNeutro,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    2,
		NumNeutrosPorTubo:  1,
		NumTierras:         1,
		SistemaElectrico:   SistemaElectricoBifasico,
	}

	result := CalcularPosicionesTuberia(params)

	// Should return 4 positions: A, B, N, T
	require.Len(t, result, 4)

	etiquetas := make([]string, len(result))
	for i, c := range result {
		etiquetas[i] = c.Etiqueta
	}
	assert.Contains(t, etiquetas, "A")
	assert.Contains(t, etiquetas, "B")
	assert.Contains(t, etiquetas, "N")
	assert.Contains(t, etiquetas, "T")
}

func TestCalcularPosicionesTuberia_BIFASICO_Positions(t *testing.T) {
	radioFase := 10.0
	radioN := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaNeutro := math.Pi * radioN * radioN
	areaTierra := math.Pi * radioFase * radioFase

	params := ParametrosTuberia{
		DiametroInteriorMM: 50.0,
		DiametroExteriorMM: 60.0,
		AreaFaseMM2:        areaFase,
		AreaNeutroMM2:      &areaNeutro,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    2,
		NumNeutrosPorTubo:  1,
		NumTierras:         1,
		SistemaElectrico:   SistemaElectricoBifasico,
	}

	result := CalcularPosicionesTuberia(params)

	a := findConductor(result, "A")
	b := findConductor(result, "B")
	n := findConductor(result, "N")

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.NotNil(t, n)

	// A-B touching: distance = 2*rF
	distAB := distance(a.CX, a.CY, b.CX, b.CY)
	assert.InDelta(t, distAB, 2*radioFase, 0.1, "A and B should touch")

	// N touches A and B
	distAN := distance(a.CX, a.CY, n.CX, n.CY)
	distBN := distance(b.CX, b.CY, n.CX, n.CY)
	assert.InDelta(t, distAN, radioFase+radioN, 0.1, "N should touch A")
	assert.InDelta(t, distBN, radioFase+radioN, 0.1, "N should touch B")
}

func TestCalcularPosicionesTuberia_Empty(t *testing.T) {
	areaFase := math.Pi * 10 * 10
	areaNeutro := math.Pi * 10 * 10
	areaTierra := math.Pi * 10 * 10

	params := ParametrosTuberia{
		DiametroInteriorMM: 50.0,
		DiametroExteriorMM: 60.0,
		AreaFaseMM2:        areaFase,
		AreaNeutroMM2:      &areaNeutro,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    0,
		NumNeutrosPorTubo:  0,
		NumTierras:         0,
		SistemaElectrico:   SistemaElectricoDelta,
	}

	result := CalcularPosicionesTuberia(params)

	require.Len(t, result, 0, "Empty tube should return empty array")
}

func TestCalcularPosicionesTuberia_AllInside_DELTA(t *testing.T) {
	diametroInterior := 50.0
	R := diametroInterior / 2
	radioFase := 10.0
	radioTierra := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaTierra := math.Pi * radioTierra * radioTierra

	params := ParametrosTuberia{
		DiametroInteriorMM: diametroInterior,
		DiametroExteriorMM: 60.0,
		AreaFaseMM2:        areaFase,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    3,
		NumNeutrosPorTubo:  0,
		NumTierras:         1,
		SistemaElectrico:   SistemaElectricoDelta,
	}

	result := CalcularPosicionesTuberia(params)

	for _, conductor := range result {
		distFromCenter := distFromOrigin(conductor.CX, conductor.CY)
		cableEdgeDist := distFromCenter + conductor.Radio
		assert.LessOrEqual(t, cableEdgeDist, R+0.01, "All conductors should be inside tube")
	}
}

func TestCalcularPosicionesTuberia_AllInside_ESTRELLA(t *testing.T) {
	diametroInterior := 50.0
	R := diametroInterior / 2
	radioFase := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaNeutro := math.Pi * radioFase * radioFase
	areaTierra := math.Pi * radioFase * radioFase

	params := ParametrosTuberia{
		DiametroInteriorMM: diametroInterior,
		DiametroExteriorMM: 60.0,
		AreaFaseMM2:        areaFase,
		AreaNeutroMM2:      &areaNeutro,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    3,
		NumNeutrosPorTubo:  1,
		NumTierras:         1,
		SistemaElectrico:   SistemaElectricoEstrella,
	}

	result := CalcularPosicionesTuberia(params)

	for _, conductor := range result {
		distFromCenter := distFromOrigin(conductor.CX, conductor.CY)
		cableEdgeDist := distFromCenter + conductor.Radio
		assert.LessOrEqual(t, cableEdgeDist, R+0.01, "All conductors should be inside tube")
	}
}

func TestCalcularPosicionesTuberia_AllInside_BIFASICO(t *testing.T) {
	diametroInterior := 50.0
	R := diametroInterior / 2
	radioFase := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaNeutro := math.Pi * radioFase * radioFase
	areaTierra := math.Pi * radioFase * radioFase

	params := ParametrosTuberia{
		DiametroInteriorMM: diametroInterior,
		DiametroExteriorMM: 60.0,
		AreaFaseMM2:        areaFase,
		AreaNeutroMM2:      &areaNeutro,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    2,
		NumNeutrosPorTubo:  1,
		NumTierras:         1,
		SistemaElectrico:   SistemaElectricoBifasico,
	}

	result := CalcularPosicionesTuberia(params)

	for _, conductor := range result {
		distFromCenter := distFromOrigin(conductor.CX, conductor.CY)
		cableEdgeDist := distFromCenter + conductor.Radio
		assert.LessOrEqual(t, cableEdgeDist, R+0.01, "All conductors should be inside tube")
	}
}

func TestCalcularPosicionesTuberia_DELTA_TouchingValidation(t *testing.T) {
	radioFase := 10.0
	areaFase := math.Pi * radioFase * radioFase
	areaTierra := math.Pi * radioFase * radioFase

	params := ParametrosTuberia{
		DiametroInteriorMM: 50.0,
		DiametroExteriorMM: 60.0,
		AreaFaseMM2:        areaFase,
		AreaTierraMM2:      areaTierra,
		NumFasesPorTubo:    3,
		NumNeutrosPorTubo:  0,
		NumTierras:         1,
		SistemaElectrico:   SistemaElectricoDelta,
	}

	result := CalcularPosicionesTuberia(params)

	a := findConductor(result, "A")
	b := findConductor(result, "B")
	c := findConductor(result, "C")
	tierra := findConductor(result, "T")

	require.NotNil(t, a)
	require.NotNil(t, b)
	require.NotNil(t, c)
	require.NotNil(t, tierra)

	// A-B distance ≈ 2*rF (touching)
	distAB := distance(a.CX, a.CY, b.CX, b.CY)
	assert.InDelta(t, distAB, 2*radioFase, 0.1, "A-B should touch")

	// A-C distance ≈ 2*rF (touching)
	distAC := distance(a.CX, a.CY, c.CX, c.CY)
	assert.InDelta(t, distAC, 2*radioFase, 0.1, "A-C should touch")

	// B-C distance ≈ 2*rF (touching)
	distBC := distance(b.CX, b.CY, c.CX, c.CY)
	assert.InDelta(t, distBC, 2*radioFase, 0.1, "B-C should touch")

	// T touches B (distance between centers ≈ rF + rT)
	distBT := distance(b.CX, b.CY, tierra.CX, tierra.CY)
	assert.InDelta(t, distBT, radioFase+radioFase, 0.1, "T should touch B")
}

// ============================================================================
// Tests: CalcularViewBox
// ============================================================================

func TestCalcularViewBox_DefaultMargin(t *testing.T) {
	result := CalcularViewBox(100, 50, 20)

	assert.Equal(t, "-20.00 -20.00 140.00 90.00", result.ViewBox)
	assert.Equal(t, 140.0, result.Ancho)
	assert.Equal(t, 90.0, result.Alto)
}

func TestCalcularViewBox_CustomMargin(t *testing.T) {
	result := CalcularViewBox(100, 50, 10)

	assert.Equal(t, "-10.00 -10.00 120.00 70.00", result.ViewBox)
	assert.Equal(t, 120.0, result.Ancho)
	assert.Equal(t, 70.0, result.Alto)
}

func TestCalcularViewBox_ZeroDimensions(t *testing.T) {
	result := CalcularViewBox(0, 0, 20)

	assert.Equal(t, "-20.00 -20.00 40.00 40.00", result.ViewBox)
}

func TestCalcularViewBox_NegativeMargin_DefaultsTo20(t *testing.T) {
	result := CalcularViewBox(100, 50, -5)

	assert.Equal(t, "-20.00 -20.00 140.00 90.00", result.ViewBox)
}

// ============================================================================
// Tests: CalcularAnchoOcupadoCharola
// ============================================================================

func TestCalcularAnchoOcupadoCharola_Empty(t *testing.T) {
	result := CalcularAnchoOcupadoCharola([]ConductorPosicion{})
	assert.Equal(t, 0.0, result)
}

func TestCalcularAnchoOcupadoCharola_SingleConductor(t *testing.T) {
	posiciones := []ConductorPosicion{
		{CX: 50, CY: 0, Radio: 10, Color: "#fff", Etiqueta: "A", Tipo: TipoConductorFase},
	}

	result := CalcularAnchoOcupadoCharola(posiciones)
	assert.Equal(t, 20.0, result, "Width should be radio * 2")
}

func TestCalcularAnchoOcupadoCharola_MultipleConductors(t *testing.T) {
	posiciones := []ConductorPosicion{
		{CX: 10, CY: 0, Radio: 10, Color: "#fff", Etiqueta: "A", Tipo: TipoConductorFase},
		{CX: 50, CY: 0, Radio: 10, Color: "#fff", Etiqueta: "B", Tipo: TipoConductorFase},
	}

	result := CalcularAnchoOcupadoCharola(posiciones)
	// Min left edge = 10 - 10 = 0
	// Max right edge = 50 + 10 = 60
	// Width = 60 - 0 = 60
	assert.Equal(t, 60.0, result)
}

// ============================================================================
// Integration Tests: calcularAnchoOcupadoCharola + position functions
// ============================================================================

func TestIntegration_CharolaEspaciadaAncho_DELTA(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoDelta,
		HilosPorFase:     1,
		AnchoComercialMM: anchoComercial,
	}

	posiciones := CalcularPosicionesCharolaEspaciada(params)
	anchoOcupado := CalcularAnchoOcupadoCharola(posiciones)

	// DELTA: A, B, C (3 fases) + T (tierra)
	// Width calculation: 4*diametro + 3*r + rT
	// For same diameter: 4*21.2 + 3*10.6 + 10.6 = 127.2mm
	assert.Greater(t, anchoOcupado, 0.0)
	assert.InDelta(t, anchoOcupado, 127.2, 1.0)
}

func TestIntegration_CharolaEspaciadaAncho_ESTRELLA(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaBase{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		SistemaElectrico: SistemaElectricoEstrella,
		HilosPorFase:     1,
		AnchoComercialMM: anchoComercial,
	}

	posiciones := CalcularPosicionesCharolaEspaciada(params)
	anchoOcupado := CalcularAnchoOcupadoCharola(posiciones)

	// ESTRELLA: A, B, C (3 fases) + N (neutro) + T (tierra) = 5 conductors
	// Width ≈ 169.6mm
	assert.InDelta(t, anchoOcupado, 169.6, 1.0)
}

func TestIntegration_CharolaTriangularAncho_DELTA(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     1,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoDelta,
	}

	posiciones := CalcularPosicionesCharolaTriangular(params)
	anchoOcupado := CalcularAnchoOcupadoCharola(posiciones)

	// DELTA triangular: A, B at bottom, C at top, T touching B
	// Width = 3*diametro = 63.6mm
	assert.InDelta(t, anchoOcupado, 63.6, 1.0)
}

func TestIntegration_CharolaTriangularAncho_ESTRELLA(t *testing.T) {
	diametroFase := 21.2
	diametroTierra := 21.2
	anchoComercial := 150.0

	params := ParametrosCharolaTriangular{
		DiametroFaseMM:   diametroFase,
		DiametroTierraMM: diametroTierra,
		HilosPorFase:     1,
		FactorTriangular: 1.0,
		AnchoComercialMM: anchoComercial,
		SistemaElectrico: SistemaElectricoEstrella,
	}

	posiciones := CalcularPosicionesCharolaTriangular(params)
	anchoOcupado := CalcularAnchoOcupadoCharola(posiciones)

	// ESTRELLA triangular: same width as DELTA since A-B-T same positions
	assert.InDelta(t, anchoOcupado, 63.6, 1.0)
}
