// internal/pdf/geometry/charola_triangular.go
package geometry

import (
	"fmt"
	"math"
)

// ============================================================================
// CÁLCULO DE POSICIONES PARA CHAROLA TRIANGULAR
// ============================================================================

// ParametrosCharolaTriangular contiene los parámetros específicos para el
// cálculo de posiciones en distribución triangular.
type ParametrosCharolaTriangular struct {
	// Diámetro del conductor de fase en mm.
	DiametroFaseMM float64
	// Diámetro del conductor de tierra en mm.
	DiametroTierraMM float64
	// Diámetro opcional del conductor de control en mm.
	DiametroControlMM *float64
	// Número de hilos por fase (conductores en paralelo).
	HilosPorFase int
	// Factor triangular (estándar NOM: 2.15).
	FactorTriangular float64
	// Ancho comercial de la charola en mm (para centrado).
	AnchoComercialMM float64
	// Sistema eléctrico (DELTA, ESTRELLA, BIFASICO, MONOFASICO).
	SistemaElectrico SistemaElectrico
}

// CalcularPosicionesCharolaTriangular calcula las posiciones de los conductores
// para una charola con distribución triangular.
//
// Reglas geométricas según el sistema eléctrico:
//
// DELTA (3 fases, sin neutro) — 1 hilo = 3 cables en triángulo equilátero:
//
//	C
//	/ \
//	A---B
//
// ESTRELLA (3 fases con neutro) — 1 hilo = 4 cables en rombo:
//
//	C---N
//	|  X |
//	A---B
//
// BIFASICO (2 fases con neutro) — 1 hilo = 3 cables (A, B, N en triángulo):
//
//	N
//	/ \
//	A---B
//
// MONOFASICO (1 fase con neutro) — 1 hilo = 2 cables apilados:
//
//	N
//	|
//	A
//
// Reglas:
//   - Los cables dentro de un grupo se TOCAN entre sí
//   - Los grupos están espaciados por factorTriangular × diámetro centro-a-centro (estándar NOM)
//   - Tierra (T): TOCANDO el último cable de la fila inferior (B del último grupo)
//   - Cables de control: posicionados después de tierra con 1 diámetro de espaciado
//
// La función retorna posiciones RELATIVAS a la pared izquierda interna de la charola:
//   - CX: offset horizontal desde la pared izquierda (0 = en la pared)
//   - CY: offset vertical desde el fondo (0 = sobre el fondo, positivo = sobre el fondo)
//
// Args:
//   - params: Parámetros de configuración.
//
// Returns:
//   - Slice de ConductorPosicion con las posiciones calculadas.
//
// Ejemplo de uso:
//
//	params := geometry.ParametrosCharolaTriangular{
//	    DiametroFaseMM:     15.0,
//	    DiametroTierraMM:   10.0,
//	    HilosPorFase:       1,
//	    FactorTriangular:   2.15,
//	    AnchoComercialMM:   300.0,
//	    SistemaElectrico:   geometry.SistemaElectricoDelta,
//	}
//	posiciones := geometry.CalcularPosicionesCharolaTriangular(params)
func CalcularPosicionesCharolaTriangular(params ParametrosCharolaTriangular) []ConductorPosicion {
	radioFase := params.DiametroFaseMM / 2
	radioTierra := params.DiametroTierraMM / 2
	var radioControl float64
	if params.DiametroControlMM != nil {
		radioControl = *params.DiametroControlMM / 2
	}

	// sin(60°) para altura del triángulo equilátero
	sin60 := math.Sqrt(3) / 2
	alturaTriangulo := params.DiametroFaseMM * sin60

	// Calcular ancho del grupo (extensión horizontal) según el tipo de sistema
	// DELTA, ESTRELLA, BIFASICO: 2 × Ø (A y B definen el ancho inferior)
	// MONOFASICO: 1 × Ø (solo A, N apilado)
	var groupWidth float64
	if params.SistemaElectrico == SistemaElectricoMonofasico {
		groupWidth = params.DiametroFaseMM
	} else {
		groupWidth = 2 * params.DiametroFaseMM
	}

	// Calcular ancho total para centrado ANTES de posicionar conductores
	// Para N grupos: totalPhaseWidth = N * groupWidth + (N - 1) * factorTriangular * diametroFaseMM
	totalPhaseWidth := float64(params.HilosPorFase)*groupWidth +
		float64(params.HilosPorFase-1)*params.FactorTriangular*params.DiametroFaseMM

	// Ancho total incluyendo tierra (tocando el último grupo)
	totalWidth := totalPhaseWidth + params.DiametroTierraMM

	// Agregar cables de control si existen (1 diámetro de espacio desde tierra + 1 diámetro para el cable)
	numControlCables := 0
	if params.DiametroControlMM != nil && radioControl > 0 {
		numControlCables = 1
		totalWidth += *params.DiametroControlMM + *params.DiametroControlMM // gap + cable
	}

	// Calcular offset de centrado: el centro de los cables debe estar en anchoComercialMM / 2
	offsetX := (params.AnchoComercialMM - totalWidth) / 2

	posiciones := make([]ConductorPosicion, 0, 10)

	// Colocar conductores de fase en grupos
	for hiloIdx := 0; hiloIdx < params.HilosPorFase; hiloIdx++ {
		// Offset del grupo: stride = groupWidth + espaciado entre grupos
		groupStep := groupWidth + params.FactorTriangular*params.DiametroFaseMM
		groupOffsetX := float64(hiloIdx) * groupStep

		// A: inferior izquierda (cy = 0 = sobre el fondo)
		posiciones = append(posiciones, ConductorPosicion{
			CX:       offsetX + groupOffsetX + radioFase,
			CY:       0,
			Radio:    radioFase,
			Color:    Colores["fase"],
			Etiqueta: etiquetaFase(params.HilosPorFase, hiloIdx, "A"),
			Tipo:     TipoConductorFase,
		})

		// B: inferior derecha — solo para sistemas con 2+ fases (no MONOFASICO)
		if params.SistemaElectrico != SistemaElectricoMonofasico {
			posiciones = append(posiciones, ConductorPosicion{
				CX:       offsetX + groupOffsetX + radioFase + params.DiametroFaseMM,
				CY:       0,
				Radio:    radioFase,
				Color:    Colores["fase"],
				Etiqueta: etiquetaFase(params.HilosPorFase, hiloIdx, "B"),
				Tipo:     TipoConductorFase,
			})
		}

		// Tercer (y cuarto) conductor depende del tipo de sistema
		switch params.SistemaElectrico {
		case SistemaElectricoDelta:
			// C: centro superior, centrado entre A y B (triángulo equilátero)
			posiciones = append(posiciones, ConductorPosicion{
				CX:       offsetX + groupOffsetX + params.DiametroFaseMM,
				CY:       alturaTriangulo,
				Radio:    radioFase,
				Color:    Colores["fase"],
				Etiqueta: etiquetaFase(params.HilosPorFase, hiloIdx, "C"),
				Tipo:     TipoConductorFase,
			})

		case SistemaElectricoEstrella:
			// C: superior izquierda (sobre A), N: superior derecha (sobre B) — ambos a altura de triángulo
			// C alineado con A (borde izquierdo del grupo), N alineado con B (borde derecho del grupo)
			posiciones = append(posiciones, ConductorPosicion{
				CX:       offsetX + groupOffsetX + radioFase,
				CY:       alturaTriangulo,
				Radio:    radioFase,
				Color:    Colores["fase"],
				Etiqueta: etiquetaFase(params.HilosPorFase, hiloIdx, "C"),
				Tipo:     TipoConductorFase,
			})
			posiciones = append(posiciones, ConductorPosicion{
				CX:       offsetX + groupOffsetX + radioFase + params.DiametroFaseMM,
				CY:       alturaTriangulo,
				Radio:    radioFase,
				Color:    Colores["neutro"],
				Etiqueta: etiquetaNeutro(params.HilosPorFase, hiloIdx),
				Tipo:     TipoConductorNeutro,
			})

		case SistemaElectricoBifasico:
			// N: centro superior, como la posición C de DELTA
			posiciones = append(posiciones, ConductorPosicion{
				CX:       offsetX + groupOffsetX + params.DiametroFaseMM,
				CY:       alturaTriangulo,
				Radio:    radioFase,
				Color:    Colores["neutro"],
				Etiqueta: etiquetaNeutro(params.HilosPorFase, hiloIdx),
				Tipo:     TipoConductorNeutro,
			})

		case SistemaElectricoMonofasico:
			// N: directamente sobre A (apilado, tocándose)
			posiciones = append(posiciones, ConductorPosicion{
				CX:       offsetX + groupOffsetX + radioFase,
				CY:       params.DiametroFaseMM, // Altura completa del diámetro para arreglo apilado
				Radio:    radioFase,
				Color:    Colores["neutro"],
				Etiqueta: etiquetaNeutro(params.HilosPorFase, hiloIdx),
				Tipo:     TipoConductorNeutro,
			})
		}
	}

	// Conductor de tierra: TOCANDO el último cable de la fila inferior
	// Para MONOFASICO: último cable A (no existe B)
	// Para otros sistemas: último cable B (conductor inferior derecho)
	var lastBottomCX float64
	if params.HilosPorFase > 0 {
		lastGroupStep := float64(params.HilosPorFase-1) * (groupWidth + params.FactorTriangular*params.DiametroFaseMM)
		if params.SistemaElectrico == SistemaElectricoMonofasico {
			// Último cable A (no hay B en MONOFASICO)
			lastBottomCX = offsetX + lastGroupStep + radioFase
		} else {
			// Último cable B
			lastBottomCX = offsetX + lastGroupStep + radioFase + params.DiametroFaseMM
		}
	}
	tierraCX := lastBottomCX + radioFase + radioTierra

	posiciones = append(posiciones, ConductorPosicion{
		CX:       tierraCX,
		CY:       0,
		Radio:    radioTierra,
		Color:    Colores["tierra"],
		Etiqueta: "T",
		Tipo:     TipoConductorTierra,
	})

	// Cables de control (si existen): después de tierra con 1 diámetro de espaciado
	// Gap = 1 diámetro completo de control + el cable mismo
	if params.DiametroControlMM != nil && numControlCables > 0 {
		diametroControlMM := *params.DiametroControlMM
		controlStartCX := tierraCX + radioTierra + diametroControlMM + radioControl

		posiciones = append(posiciones, ConductorPosicion{
			CX:       controlStartCX,
			CY:       0,
			Radio:    radioControl,
			Color:    Colores["control"],
			Etiqueta: "Ctrl",
			Tipo:     TipoConductorControl,
		})
	}

	return posiciones
}

// etiquetaFase genera la etiqueta para un conductor de fase.
func etiquetaFase(hilosPorFase, hiloIdx int, base string) string {
	if hilosPorFase > 1 {
		return fmt.Sprintf("%s%d", base, hiloIdx+1)
	}
	return base
}

// etiquetaNeutro genera la etiqueta para un conductor neutro.
func etiquetaNeutro(hilosPorFase, hiloIdx int) string {
	if hilosPorFase > 1 {
		return fmt.Sprintf("N%d", hiloIdx+1)
	}
	return "N"
}
