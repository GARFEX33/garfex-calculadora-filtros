// internal/pdf/geometry/viewbox.go
package geometry

import (
	"fmt"
	"math"
)

// ============================================================================
// CÁLCULO DE VIEWBOX Y COTAS
// ============================================================================

// CalcularViewBox calcula el viewBox SVG basado en las dimensiones del contenido.
//
// Args:
//   - anchoContenido: Ancho del contenido en mm.
//   - altoContenido: Alto del contenido en mm.
//   - margen: Margen opcional alrededor del contenido (default: 20 mm).
//   - numTubos: Número de tubos (para calcular ancho total de múltiples tubos).
//
// Returns:
//
//	ViewBoxResult con el string viewBox y dimensiones calculadas.
//
// Ejemplo de uso:
//
//	resultado := geometry.CalcularViewBox(300, 200, 25)
//	// resultado.ViewBox = "-25.00 -25.00 350.00 250.00"
func CalcularViewBox(anchoContenido, altoContenido, margen float64, numTubos ...int) ViewBoxResult {
	if margen <= 0 {
		margen = 20
	}

	// Calcular ancho total si hay múltiples tubos
	anchoTotal := anchoContenido
	if len(numTubos) > 0 && numTubos[0] > 1 {
		espacioEntreTubos := 30.0 // Default spacing between tubes
		anchoTotal = float64(numTubos[0])*anchoContenido + float64(numTubos[0]-1)*espacioEntreTubos
	}

	minX := -margen
	minY := -margen
	ancho := anchoTotal + 2*margen
	alto := altoContenido + 2*margen

	return ViewBoxResult{
		ViewBox: fmt.Sprintf("%.2f %.2f %.2f %.2f", minX, minY, ancho, alto),
		Ancho:   ancho,
		Alto:    alto,
	}
}

// CalcularCotasCharola calcula las líneas de dimensión para la charola.
//
// Retorna 2 líneas de cota:
//   - Superior: ancho comercial completo de la charola
//   - Inferior: ancho requerido (Area / Peralte)
//
// Args:
//   - anchoComercialMM: Ancho comercial de la charola en mm.
//   - areaRequeridaMM2: Área requerida para los conductores en mm².
//   - peralte: Peralte (altura) de la charola en mm.
//
// Returns:
//   - Slice de LineaCota con las dos líneas de dimensión.
//
// Ejemplo de uso:
//
//	cotas := geometry.CalcularCotasCharola(300, 4500, 70)
//	// cota[0]: línea superior, cota[1]: línea inferior
func CalcularCotasCharola(anchoComercialMM, areaRequeridaMM2, peralte float64) []LineaCota {
	// Calcular ancho requerido desde el área: A_req = area / peralte
	anchoRequerido := areaRequeridaMM2 / peralte

	// Convertir mm a pulgadas para mostrar
	comercialPulgadas := anchoComercialMM / 25.4
	requeridoPulgadas := anchoRequerido / 25.4

	// Redondear a 1 decimal para display
	comercialPulgadas = math.Round(comercialPulgadas*10) / 10
	requeridoPulgadas = math.Round(requeridoPulgadas*10) / 10

	cotas := make([]LineaCota, 0, 2)

	// Línea de cota superior: ancho comercial completo
	cotas = append(cotas, LineaCota{
		X1:            0,
		Y1:            -20,
		X2:            anchoComercialMM,
		Y2:            -20,
		Valor:         anchoComercialMM,
		Texto:         fmt.Sprintf("%.1f mm (%.1f\")", anchoComercialMM, comercialPulgadas),
		PosicionTexto: "arriba",
	})

	// Línea de cota inferior: ancho requerido
	cotas = append(cotas, LineaCota{
		X1:            0,
		Y1:            peralte + 20,
		X2:            anchoRequerido,
		Y2:            peralte + 20,
		Valor:         anchoRequerido,
		Texto:         fmt.Sprintf("%.1f mm (%.1f\")", anchoRequerido, requeridoPulgadas),
		PosicionTexto: "abajo",
	})

	return cotas
}
