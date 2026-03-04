// internal/pdf/geometry/svg.go
package geometry

import (
	"fmt"
	"strings"
)

// ============================================================================
// GENERADORES SVG - Phase 2
// ============================================================================

// ParametrosSVGCharola contiene los parámetros para generar el SVG completo
// de una charola con conductores.
type ParametrosSVGCharola struct {
	// Posiciones de los conductores (calculadas por las funciones de geometría).
	Posiciones []ConductorPosicion
	// Ancho comercial de la charola en mm.
	AnchoComercialMM float64
	// Peralte (altura) de la charola en mm.
	PeralteMM float64
	// Tipo de distribución ("espaciada" o "triangular").
	TipoDistribucion string
	// Offset del fondo donde se apoyan los cables.
	EspesorPisoMM float64
	// Margen izquierdo antes de la brida.
	MargenIzqMM float64
	// Ancho de la brida.
	AnchoBridaMM float64
	// Margen para cotas a la derecha.
	MargenCotaDerMM float64
	// Margen arriba para dimensión W.
	MargenArribaMM float64
	// Margen abajo.
	MargenAbajoMM float64
}

// Valores por defecto para ParametrosSVGCharola.
var DefaultParametrosSVGCharola = ParametrosSVGCharola{
	PeralteMM:       70,
	EspesorPisoMM:   7,
	MargenIzqMM:     30,
	AnchoBridaMM:    25,
	MargenCotaDerMM: 60,
	MargenArribaMM:  35,
	MargenAbajoMM:   20,
}

// ParametrosSVGTuberia contiene los parámetros para generar el SVG de tubería.
type ParametrosSVGTuberia struct {
	// Posiciones de los conductores dentro del tubo.
	Posiciones []ConductorPosicion
	// Diámetro interior del tubo en mm.
	DiametroInteriorMM float64
	// Diámetro exterior del tubo en mm.
	DiametroExteriorMM float64
	// Número de tubos.
	NumTubos int
	// Margen alrededor del diagrama.
	MargenMM float64
	// Espacio entre tubos.
	EspacioEntreTubosMM float64
	// Margen para cotas.
	MargenCotaMM float64
}

// Valores por defecto para ParametrosSVGTuberia.
var DefaultParametrosSVGTuberia = ParametrosSVGTuberia{
	MargenMM:            15,
	EspacioEntreTubosMM: 30,
	MargenCotaMM:        10,
}

// GenerarSVGCharola genera el SVG para una charola con distribución de cables.
//
// Args:
//   - posiciones: Slice de ConductorPosicion con las posiciones de los cables.
//   - ancho: Ancho comercial de la charola en mm.
//   - peralte: Peralte (altura) de la charola en mm.
//   - tipo: Tipo de distribución ("espaciada" o "triangular").
//
// Returns:
//   - String con el contenido SVG de la charola (sin wrapper <svg>).
func GenerarSVGCharola(posiciones []ConductorPosicion, ancho, peralte float64, tipo string) string {
	var sb strings.Builder

	// Parámetros de layout
	brida := 25.0
	espesorPiso := 7.0

	// Coordenadas de la charola
	charolaLeft := 30.0 + brida // margen izq (30) + brida
	charolaRight := charolaLeft + ancho
	charolaTop := 35.0 // margenArriba
	charolaBottom := charolaTop + peralte

	// Bridas
	flangeLeft := charolaLeft - brida
	flangeRight := charolaRight + brida

	// Línea W (horizontal arriba)
	wLineY := charolaTop - 15
	wMidX := charolaLeft + ancho/2

	// Línea P (vertical derecha)
	pLineX := flangeRight + 20
	pMidY := charolaTop + peralte/2

	// Líneas del piso
	floorSolidY := charolaBottom - 5
	floorDashedY := charolaBottom - 3

	// === GENERAR CONTENIDO SVG ===

	// Perfil U de la charola con bridas
	sb.WriteString(fmt.Sprintf(
		`<path d="M %.2f,%.2f L %.2f,%.2f L %.2f,%.2f L %.2f,%.2f L %.2f,%.2f L %.2f,%.2f" stroke="%s" stroke-width="2.5" fill="none" stroke-linejoin="round"/>`,
		flangeLeft, charolaTop,
		charolaLeft, charolaTop,
		charolaLeft, charolaBottom,
		charolaRight, charolaBottom,
		charolaRight, charolaTop,
		flangeRight, charolaTop,
		Colores["charolaStroke"],
	))

	// Dimensión W (horizontal arriba)
	sb.WriteString(fmt.Sprintf(
		`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
		charolaLeft, wLineY, charolaRight, wLineY, Colores["cotaLinea"],
	))
	// Ticks W
	sb.WriteString(fmt.Sprintf(
		`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
		charolaLeft, wLineY-5, charolaLeft, wLineY+5, Colores["cotaLinea"],
	))
	sb.WriteString(fmt.Sprintf(
		`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
		charolaRight, wLineY-5, charolaRight, wLineY+5, Colores["cotaLinea"],
	))
	// Texto W
	sb.WriteString(fmt.Sprintf(
		`<text x="%.2f" y="%.2f" font-size="8" text-anchor="middle" fill="%s" font-family="Arial, sans-serif">%.1f mm</text>`,
		wMidX, wLineY-6, Colores["cotaTexto"], ancho,
	))

	// Dimensión P (vertical derecha)
	sb.WriteString(fmt.Sprintf(
		`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
		pLineX, charolaTop, pLineX, charolaBottom, Colores["cotaLinea"],
	))
	// Ticks P
	sb.WriteString(fmt.Sprintf(
		`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
		pLineX-5, charolaTop, pLineX+5, charolaTop, Colores["cotaLinea"],
	))
	sb.WriteString(fmt.Sprintf(
		`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
		pLineX-5, charolaBottom, pLineX+5, charolaBottom, Colores["cotaLinea"],
	))
	// Texto P
	sb.WriteString(fmt.Sprintf(
		`<text x="%.2f" y="%.2f" font-size="8" text-anchor="start" fill="%s" font-family="Arial, sans-serif">%.0f mm</text>`,
		pLineX+8, pMidY+3, Colores["cotaTexto"], peralte,
	))

	// Líneas del piso
	sb.WriteString(fmt.Sprintf(
		`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="2"/>`,
		charolaLeft, floorSolidY, charolaRight, floorSolidY, Colores["charolaStroke"],
	))
	sb.WriteString(fmt.Sprintf(
		`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="2" stroke-dasharray="8,5"/>`,
		charolaLeft, floorDashedY, charolaRight, floorDashedY, Colores["charolaStroke"],
	))

	// Conductores
	for _, cond := range posiciones {
		var cx, cy float64

		if tipo == "triangular" {
			// Para triangular: CY es offset vertical desde el fondo
			cx = charolaLeft + cond.CX
			cy = charolaBottom - espesorPiso - cond.Radio - cond.CY
		} else {
			// Para espaciada: CY es 0 (todos en el fondo)
			cx = charolaLeft + cond.CX
			cy = charolaBottom - espesorPiso - cond.Radio
		}

		// Círculo del conductor
		sb.WriteString(fmt.Sprintf(
			`<circle cx="%.2f" cy="%.2f" r="%.2f" stroke="%s" stroke-width="2" fill="none"/>`,
			cx, cy, cond.Radio, Colores["charolaStroke"],
		))

		// Etiqueta dentro del círculo
		fontSize := cond.Radio * 0.8
		if fontSize < 4 {
			fontSize = 4
		}
		textY := cy + cond.Radio*0.35

		sb.WriteString(fmt.Sprintf(
			`<text x="%.2f" y="%.2f" font-size="%.1f" text-anchor="middle" stroke="none" fill="%s" font-weight="bold" font-family="Arial, sans-serif">%s</text>`,
			cx, textY, fontSize, Colores["charolaStroke"], cond.Etiqueta,
		))
	}

	return sb.String()
}

// GenerarSVGTuberia genera el SVG para una sección transversal de tubería.
//
// Args:
//   - posiciones: Slice de ConductorPosicion con las posiciones de los cables.
//   - diametroInterior: Diámetro interior del tubo en mm.
//   - diametroExterior: Diámetro exterior del tubo en mm.
//
// Returns:
//   - String con el contenido SVG de la tubería (sin wrapper <svg>).
func GenerarSVGTuberia(posiciones []ConductorPosicion, diametroInterior, diametroExterior float64) string {
	var sb strings.Builder

	radioInterior := diametroInterior / 2
	radioExterior := diametroExterior / 2
	wallThickness := (diametroExterior - diametroInterior) / 2

	// Centro del tubo en (0, 0) relativo
	centerX := 0.0
	centerY := 0.0

	// Círculo principal del tubo (espesor de pared)
	midRadius := (radioInterior + radioExterior) / 2
	sb.WriteString(fmt.Sprintf(
		`<circle cx="%.2f" cy="%.2f" r="%.2f" stroke="%s" stroke-width="%.2f" fill="none"/>`,
		centerX, centerY, midRadius, Colores["tuboStroke"], wallThickness,
	))

	// Conductores dentro del tubo
	for _, cond := range posiciones {
		// Los valores de CX y CY ya vienen relativos al centro del tubo
		cx := centerX + cond.CX
		cy := centerY + cond.CY

		// Círculo del conductor
		sb.WriteString(fmt.Sprintf(
			`<circle cx="%.2f" cy="%.2f" r="%.2f" stroke="%s" stroke-width="1" fill="none"/>`,
			cx, cy, cond.Radio, Colores["tuboStroke"],
		))

		// Etiqueta dentro del círculo
		fontSize := cond.Radio * 0.9
		if fontSize < 4 {
			fontSize = 4
		}

		sb.WriteString(fmt.Sprintf(
			`<text x="%.2f" y="%.2f" font-size="%.1f" text-anchor="middle" dominant-baseline="central" stroke="none" fill="%s" font-weight="bold" font-family="Arial, sans-serif">%s</text>`,
			cx, cy, fontSize, Colores["tuboStroke"], cond.Etiqueta,
		))
	}

	return sb.String()
}

// GenerarSVGCotas genera las líneas de cota para el diagrama.
//
// Args:
//   - cotas: Slice de LineaCota con las líneas de dimensión a dibujar.
//
// Returns:
//   - String con los elementos SVG de cotas.
func GenerarSVGCotas(cotas []LineaCota) string {
	var sb strings.Builder

	for _, cota := range cotas {
		// Línea principal de cota
		sb.WriteString(fmt.Sprintf(
			`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
			cota.X1, cota.Y1, cota.X2, cota.Y2, Colores["cotaLinea"],
		))

		// Determinar posición de ticks y texto
		esHorizontal := abs(cota.Y1-cota.Y2) < abs(cota.X2-cota.X1)

		if esHorizontal {
			// Cota horizontal
			// Ticks verticales en cada extremo
			sb.WriteString(fmt.Sprintf(
				`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
				cota.X1, cota.Y1-5, cota.X1, cota.Y1+5, Colores["cotaLinea"],
			))
			sb.WriteString(fmt.Sprintf(
				`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
				cota.X2, cota.Y2-5, cota.X2, cota.Y2+5, Colores["cotaLinea"],
			))

			// Texto
			textX := (cota.X1 + cota.X2) / 2
			var textY float64
			if cota.PosicionTexto == "arriba" {
				textY = cota.Y1 - 6
			} else {
				textY = cota.Y1 + 12
			}
			sb.WriteString(fmt.Sprintf(
				`<text x="%.2f" y="%.2f" font-size="6" text-anchor="middle" fill="%s" font-family="Arial, sans-serif">%s</text>`,
				textX, textY, Colores["cotaTexto"], cota.Texto,
			))
		} else {
			// Cota vertical
			// Ticks horizontales en cada extremo
			sb.WriteString(fmt.Sprintf(
				`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
				cota.X1-5, cota.Y1, cota.X1+5, cota.Y1, Colores["cotaLinea"],
			))
			sb.WriteString(fmt.Sprintf(
				`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
				cota.X2-5, cota.Y2, cota.X2+5, cota.Y2, Colores["cotaLinea"],
			))

			// Texto
			textY := (cota.Y1 + cota.Y2) / 2
			var textX float64
			if cota.PosicionTexto == "derecha" {
				textX = cota.X1 + 8
			} else {
				textX = cota.X1 - 8
			}
			sb.WriteString(fmt.Sprintf(
				`<text x="%.2f" y="%.2f" font-size="6" text-anchor="start" fill="%s" font-family="Arial, sans-serif">%s</text>`,
				textX, textY+3, Colores["cotaTexto"], cota.Texto,
			))
		}
	}

	return sb.String()
}

// GenerarSVGCompletoCharola genera el SVG completo de una charola con todos sus elementos.
//
// Args:
//   - params: Parámetros de configuración para el SVG.
//
// Returns:
//   - String completo con el SVG incluyendo wrapper <svg>, definiciones y contenido.
func GenerarSVGCompletoCharola(params ParametrosSVGCharola) string {
	var sb strings.Builder

	// Valores por defecto
	ancho := params.AnchoComercialMM
	if ancho == 0 {
		ancho = 150
	}
	peralte := params.PeralteMM
	if peralte == 0 {
		peralte = DefaultParametrosSVGCharola.PeralteMM
	}
	tipo := params.TipoDistribucion
	if tipo == "" {
		tipo = "espaciada"
	}

	// Parámetros de layout
	brida := params.AnchoBridaMM
	if brida == 0 {
		brida = DefaultParametrosSVGCharola.AnchoBridaMM
	}
	margenIzq := params.MargenIzqMM
	if margenIzq == 0 {
		margenIzq = DefaultParametrosSVGCharola.MargenIzqMM
	}
	margenCotaDer := params.MargenCotaDerMM
	if margenCotaDer == 0 {
		margenCotaDer = DefaultParametrosSVGCharola.MargenCotaDerMM
	}
	margenArriba := params.MargenArribaMM
	if margenArriba == 0 {
		margenArriba = DefaultParametrosSVGCharola.MargenArribaMM
	}
	margenAbajo := params.MargenAbajoMM
	if margenAbajo == 0 {
		margenAbajo = DefaultParametrosSVGCharola.MargenAbajoMM
	}

	// Coordenadas de la charola
	charolaLeft := margenIzq + brida
	charolaRight := charolaLeft + ancho
	charolaTop := margenArriba
	charolaBottom := charolaTop + peralte

	// Ancho y alto total del SVG
	svgWidth := charolaRight + brida + margenCotaDer
	svgHeight := charolaBottom + margenAbajo

	// === INICIO SVG ===
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.2f %.2f" preserveAspectRatio="xMidYMid meet">`, svgWidth, svgHeight))

	// === DEFINICIONES (PATRONES) ===
	sb.WriteString(`<defs>`)
	sb.WriteString(`<pattern id="charola-hatch" patternUnits="userSpaceOnUse" width="8" height="8" patternTransform="rotate(45)">`)
	sb.WriteString(`<line x1="0" y1="0" x2="0" y2="8" stroke="#004085" stroke-width="2"/>`)
	sb.WriteString(`</pattern>`)
	sb.WriteString(`</defs>`)

	// === CONTENIDO ===
	contenido := GenerarSVGCharola(params.Posiciones, ancho, peralte, tipo)
	sb.WriteString(contenido)

	// === CIERRE SVG ===
	sb.WriteString(`</svg>`)

	return sb.String()
}

// GenerarSVGCompletoTuberia genera el SVG completo de una tubería con todos sus elementos.
//
// Args:
//   - posiciones: Slice de ConductorPosicion con las posiciones de los cables.
//   - diametroInterior: Diámetro interior del tubo en mm.
//   - diametroExterior: Diámetro exterior del tubo en mm.
//   - numTubos: Número de tubos.
//   - espacioEntreTubos: Espacio entre tubos en mm.
//
// Returns:
//   - String completo con el SVG incluyendo wrapper <svg>, definiciones y contenido.
func GenerarSVGCompletoTuberia(posiciones []ConductorPosicion, diametroInterior, diametroExterior float64, numTubos int, espacioEntreTubos float64) string {
	var sb strings.Builder

	if numTubos <= 0 {
		numTubos = 1
	}
	if espacioEntreTubos <= 0 {
		espacioEntreTubos = DefaultParametrosSVGTuberia.EspacioEntreTubosMM
	}

	margen := DefaultParametrosSVGTuberia.MargenMM
	margenCota := DefaultParametrosSVGTuberia.MargenCotaMM

	radioExterior := diametroExterior / 2

	// Ancho total
	anchoTubos := float64(numTubos)*diametroExterior + float64(numTubos-1)*espacioEntreTubos
	svgWidth := anchoTubos + 2*margen
	svgHeight := diametroExterior + 2*margen + margenCota + 8

	// Centro Y para todos los tubos
	centerY := margen + radioExterior

	// === INICIO SVG ===
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.2f %.2f" preserveAspectRatio="xMidYMid meet">`, svgWidth, svgHeight))

	// === DEFINICIONES ===
	sb.WriteString(`<defs>`)
	sb.WriteString(`<pattern id="tubo-hatch" patternUnits="userSpaceOnUse" width="6" height="6" patternTransform="rotate(45)">`)
	sb.WriteString(`<line x1="0" y1="0" x2="0" y2="6" stroke="#555555" stroke-width="1.5"/>`)
	sb.WriteString(`</pattern>`)
	sb.WriteString(`</defs>`)

	// === CONTENIDO: TUBOS ===
	for idx := 0; idx < numTubos; idx++ {
		centerX := margen + radioExterior + float64(idx)*(diametroExterior+espacioEntreTubos)
		dimY := centerY + radioExterior + 15

		// Tubo (círculo con espesor)
		radioInterior := diametroInterior / 2
		midRadius := (radioInterior + radioExterior) / 2
		wallThickness := (diametroExterior - diametroInterior) / 2

		sb.WriteString(fmt.Sprintf(
			`<circle cx="%.2f" cy="%.2f" r="%.2f" stroke="%s" stroke-width="%.2f" fill="none"/>`,
			centerX, centerY, midRadius, Colores["tuboStroke"], wallThickness,
		))

		// Conductores
		for _, cond := range posiciones {
			cx := centerX + cond.CX
			cy := centerY + cond.CY

			sb.WriteString(fmt.Sprintf(
				`<circle cx="%.2f" cy="%.2f" r="%.2f" stroke="%s" stroke-width="1" fill="none"/>`,
				cx, cy, cond.Radio, Colores["tuboStroke"],
			))

			fontSize := cond.Radio * 0.9
			if fontSize < 4 {
				fontSize = 4
			}
			sb.WriteString(fmt.Sprintf(
				`<text x="%.2f" y="%.2f" font-size="%.1f" text-anchor="middle" dominant-baseline="central" stroke="none" fill="%s" font-weight="bold" font-family="Arial, sans-serif">%s</text>`,
				cx, cy, fontSize, Colores["tuboStroke"], cond.Etiqueta,
			))
		}

		// Dimensión del diámetro
		sb.WriteString(fmt.Sprintf(
			`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
			centerX-radioExterior, dimY, centerX+radioExterior, dimY, Colores["cotaLinea"],
		))
		sb.WriteString(fmt.Sprintf(
			`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
			centerX-radioExterior, dimY-5, centerX-radioExterior, dimY+5, Colores["cotaLinea"],
		))
		sb.WriteString(fmt.Sprintf(
			`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="1.5"/>`,
			centerX+radioExterior, dimY-5, centerX+radioExterior, dimY+5, Colores["cotaLinea"],
		))
		sb.WriteString(fmt.Sprintf(
			`<text x="%.2f" y="%.2f" font-size="6" text-anchor="middle" fill="%s" font-family="Arial, sans-serif">Ø %.1f mm</text>`,
			centerX, dimY+12, Colores["cotaTexto"], diametroExterior,
		))

		// Etiqueta del tubo (solo si hay más de uno)
		if numTubos > 1 {
			sb.WriteString(fmt.Sprintf(
				`<text x="%.2f" y="%.2f" font-size="5" text-anchor="middle" fill="%s" font-family="Arial, sans-serif">Tubo %d de %d</text>`,
				centerX, dimY+22, Colores["cotaTexto"], idx+1, numTubos,
			))
		}
	}

	// === CIERRE SVG ===
	sb.WriteString(`</svg>`)

	return sb.String()
}

// abs returns the absolute value of a float64.
func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
