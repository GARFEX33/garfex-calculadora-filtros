// internal/pdf/geometry/tuberia.go
package geometry

import (
	"math"
)

// ============================================================================
// CÁLCULO DE POSICIONES PARA TUBERÍA (CONDUIT)
// ============================================================================

// ParametrosTuberia contiene los parámetros para el cálculo de posiciones
// de conductores dentro de una tubería.
type ParametrosTuberia struct {
	// Diámetro interior de la tubería en mm.
	DiametroInteriorMM float64
	// Diámetro exterior de la tubería en mm.
	DiametroExteriorMM float64
	// Área del conductor de fase en mm².
	AreaFaseMM2 float64
	// Área opcional del conductor neutro en mm².
	AreaNeutroMM2 *float64
	// Área del conductor de tierra en mm².
	AreaTierraMM2 float64
	// Número de fases por tubo.
	NumFasesPorTubo int
	// Número de neutros por tubo.
	NumNeutrosPorTubo int
	// Número de tierras por tubo.
	NumTierras int
	// Sistema eléctrico (DELTA, ESTRELLA, BIFASICO, MONOFASICO).
	SistemaElectrico SistemaElectrico
}

// CalcularPosicionesTuberia calcula las posiciones de los conductores
// para una sección transversal de tubería (conduit).
//
// Reglas geométricas:
//   - El tubo es un círculo grueso que representa el espesor de la pared
//   - El centro del tubo está en (0, 0) — el componente traducirá
//   - Radios de cables calculados desde área: r = √(area / π)
//   - Los cables se apoyan en la PARED INFERIOR del tubo debido a la gravedad
//   - Conductores fase+neutro empacados en filas en la parte inferior
//   - Tierra (T): TOCANDO el cable derecho de la fila inferior
//
// Sistema de coordenadas: posiciones relativas al centro del tubo (0, 0)
//   - Dirección +Y = abajo (convención SVG)
//   - Fondo del tubo = (0, radioInterior)
//   - Centros de cables de fila inferior = (cx, radioInterior - radioFase)
//
// Args:
//   - params: Parámetros de configuración.
//
// Returns:
//   - Slice de ConductorPosicion con las posiciones calculadas.
//
// Ejemplo de uso:
//
//	params := geometry.ParametrosTuberia{
//	    DiametroInteriorMM: 25.0,
//	    DiametroExteriorMM: 32.0,
//	    AreaFaseMM2:        150.0,
//	    AreaTierraMM2:      50.0,
//	    NumFasesPorTubo:    3,
//	    NumNeutrosPorTubo:  1,
//	    NumTierras:         1,
//	    SistemaElectrico:   geometry.SistemaElectricoEstrella,
//	}
//	posiciones := geometry.CalcularPosicionesTuberia(params)
func CalcularPosicionesTuberia(params ParametrosTuberia) []ConductorPosicion {
	R := params.DiametroInteriorMM / 2            // Radio interior del tubo
	rF := math.Sqrt(params.AreaFaseMM2 / math.Pi) // Radio del cable de fase (desde área)
	rN := rF                                      // Por defecto, neutro tiene mismo radio que fase
	if params.AreaNeutroMM2 != nil {
		rN = math.Sqrt(*params.AreaNeutroMM2 / math.Pi)
	}
	rT := math.Sqrt(params.AreaTierraMM2 / math.Pi) // Radio de tierra

	posiciones := make([]ConductorPosicion, 0, 10)

	if params.NumFasesPorTubo == 0 && params.NumNeutrosPorTubo == 0 && params.NumTierras == 0 {
		return posiciones
	}

	// Helper: posicionar un cable en la pared interna del tubo en ángulo theta
	// theta = 0 = fondo, positivo = horario
	// Retorna {cx, cy} relativo al centro del tubo, +Y = abajo
	posOnWall := func(r float64, theta float64) (cx, cy float64) {
		d := R - r                                      // distancia del centro del tubo al centro del cable
		return d * math.Sin(theta), d * math.Cos(theta) // +Y = abajo, cos(0) = fondo
	}

	// Helper: ángulo para dos cables que se tocan en la pared
	// Dos cables de radio r1, r2 apoyados en la pared interna de radio R, tocándose
	// Semi-ángulo desde el eje vertical al centro de cada cable
	halfAngleTwoTouching := func(r1, r2 float64) float64 {
		d1 := R - r1
		d2 := R - r2
		touchDist := r1 + r2
		// Ley de cosenos: touchDist² = d1² + d2² - 2·d1·d2·cos(angle)
		cosAngle := (d1*d1 + d2*d2 - touchDist*touchDist) / (2 * d1 * d2)
		fullAngle := math.Acos(max(-1, min(1, cosAngle)))
		// Retorna SEMI-ángulo (cada cable está desplazado medio ángulo del eje central)
		return fullAngle / 2
	}

	// ========================================================================
	// ESTRATEGIA POR TIPO DE SISTEMA
	// ========================================================================
	// Sistema de coordenadas: centro del tubo = (0,0), +Y = abajo (SVG)
	// Fondo del tubo = (0, R)

	switch params.SistemaElectrico {
	case SistemaElectricoMonofasico:
		// 1 fase (A) en el fondo centro sobre la pared + 1 neutro (N) sobre A tocándose
		posA_x, posA_y := posOnWall(rF, 0) // fondo centro
		posiciones = append(posiciones, ConductorPosicion{
			CX:       posA_x,
			CY:       posA_y,
			Radio:    rF,
			Color:    Colores["fase"],
			Etiqueta: "A",
			Tipo:     TipoConductorFase,
		})

		if params.NumNeutrosPorTubo > 0 {
			// N directamente sobre A, tocando A (no necesariamente en la pared)
			posiciones = append(posiciones, ConductorPosicion{
				CX:       posA_x,
				CY:       posA_y - rF - rN,
				Radio:    rN,
				Color:    Colores["neutro"],
				Etiqueta: "N",
				Tipo:     TipoConductorNeutro,
			})
		}

	case SistemaElectricoBifasico:
		// 2 fases (A, B) en el fondo tocándose + 1 neutro (N) sobre ellos
		alpha := halfAngleTwoTouching(rF, rF)
		posA_x, posA_y := posOnWall(rF, -alpha) // izquierda
		posB_x, posB_y := posOnWall(rF, alpha)  // derecha

		posiciones = append(posiciones, ConductorPosicion{
			CX:       posA_x,
			CY:       posA_y,
			Radio:    rF,
			Color:    Colores["fase"],
			Etiqueta: "A",
			Tipo:     TipoConductorFase,
		})
		posiciones = append(posiciones, ConductorPosicion{
			CX:       posB_x,
			CY:       posB_y,
			Radio:    rF,
			Color:    Colores["fase"],
			Etiqueta: "B",
			Tipo:     TipoConductorFase,
		})

		if params.NumNeutrosPorTubo > 0 {
			// N sobre A y B, centrado
			// Más preciso: N en cx=0, cy = A.cy - sqrt((rF+rN)² - (posA.x)²)
			dxN := 0 - posA_x
			contactDist := rF + rN
			dyN := math.Sqrt(max(0, contactDist*contactDist-dxN*dxN))
			posiciones = append(posiciones, ConductorPosicion{
				CX:       0,
				CY:       posA_y - dyN,
				Radio:    rN,
				Color:    Colores["neutro"],
				Etiqueta: "N",
				Tipo:     TipoConductorNeutro,
			})
		}

	case SistemaElectricoDelta:
		// 3 fases (A, B, C): A,B en el fondo tocándose, C sobre ellos
		alpha := halfAngleTwoTouching(rF, rF)
		posA_x, posA_y := posOnWall(rF, -alpha)
		posB_x, posB_y := posOnWall(rF, alpha)

		posiciones = append(posiciones, ConductorPosicion{
			CX:       posA_x,
			CY:       posA_y,
			Radio:    rF,
			Color:    Colores["fase"],
			Etiqueta: "A",
			Tipo:     TipoConductorFase,
		})
		posiciones = append(posiciones, ConductorPosicion{
			CX:       posB_x,
			CY:       posB_y,
			Radio:    rF,
			Color:    Colores["fase"],
			Etiqueta: "B",
			Tipo:     TipoConductorFase,
		})

		// C sobre A y B (triángulo equilátero, todos mismo radio)
		// Distancia de C al centro de A debe = 2*rF
		// C en cx=0: sqrt(A.cx² + (C.cy - A.cy)²) = 2*rF
		// C.cy = A.cy - sqrt((2*rF)² - A.cx²)
		cCy := posA_y - math.Sqrt(max(0, 2*rF*(2*rF)-posA_x*posA_x))
		posiciones = append(posiciones, ConductorPosicion{
			CX:       0,
			CY:       cCy,
			Radio:    rF,
			Color:    Colores["fase"],
			Etiqueta: "C",
			Tipo:     TipoConductorFase,
		})

	case SistemaElectricoEstrella:
		// 4 cables: A,B en fondo, C arriba-izquierda, N arriba-derecha (patrón diamante)
		alpha := halfAngleTwoTouching(rF, rF)
		posA_x, posA_y := posOnWall(rF, -alpha)
		posB_x, posB_y := posOnWall(rF, alpha)

		posiciones = append(posiciones, ConductorPosicion{
			CX:       posA_x,
			CY:       posA_y,
			Radio:    rF,
			Color:    Colores["fase"],
			Etiqueta: "A",
			Tipo:     TipoConductorFase,
		})
		posiciones = append(posiciones, ConductorPosicion{
			CX:       posB_x,
			CY:       posB_y,
			Radio:    rF,
			Color:    Colores["fase"],
			Etiqueta: "B",
			Tipo:     TipoConductorFase,
		})

		// C: sobre A (tocando A, posicionado directamente arriba-izquierda)
		// Para diamante: C sobre A, N sobre B
		posiciones = append(posiciones, ConductorPosicion{
			CX:       posA_x,
			CY:       posA_y - 2*rF,
			Radio:    rF,
			Color:    Colores["fase"],
			Etiqueta: "C",
			Tipo:     TipoConductorFase,
		})
		posiciones = append(posiciones, ConductorPosicion{
			CX:       posB_x,
			CY:       posB_y - 2*rN,
			Radio:    rN,
			Color:    Colores["neutro"],
			Etiqueta: "N",
			Tipo:     TipoConductorNeutro,
		})
	}

	// ========================================================================
	// TIERRA (T): se apoya en la pared del tubo, tocando el cable derecho de la fila inferior
	// ========================================================================
	if params.NumTierras > 0 && len(posiciones) > 0 {
		// Encontrar el cable derecho de la fila inferior
		// (el de mayor cx entre los cables apoyados en la pared)
		// Para todos los sistemas: B es el derecho (o A si es MONOFASICO)
		var rightmostBottom ConductorPosicion
		if params.SistemaElectrico == SistemaElectricoMonofasico {
			rightmostBottom = posiciones[0] // A
		} else {
			rightmostBottom = posiciones[1] // B
		}

		// T se apoya en la pared y toca el cable derecho
		// Centro de T a distancia (R - rT) del centro del tubo
		// Distancia entre centro de T y centro del cable derecho = rightmost.radio + rT
		// Usar geometría: encontrar ángulo de T en la pared tal que toque el cable derecho
		dRight := math.Sqrt(rightmostBottom.CX*rightmostBottom.CX + rightmostBottom.CY*rightmostBottom.CY)
		dT := R - rT
		touchDist := rightmostBottom.Radio + rT
		// Ley de cosenos para encontrar ángulo entre rightmostBottom y T (desde centro del tubo)
		cosGamma := (dRight*dRight + dT*dT - touchDist*touchDist) / (2 * dRight * dT)
		gamma := math.Acos(max(-1, min(1, cosGamma)))
		// Ángulo del cable derecho desde el eje +Y (fondo)
		angleRight := math.Atan2(rightmostBottom.CX, rightmostBottom.CY)
		// T se coloca en sentido horario desde el cable derecho (a la derecha y ligeramente arriba)
		angleT := angleRight + gamma
		posTierra_x, posTierra_y := posOnWall(rT, angleT)

		posiciones = append(posiciones, ConductorPosicion{
			CX:       posTierra_x,
			CY:       posTierra_y,
			Radio:    rT,
			Color:    Colores["tierra"],
			Etiqueta: "T",
			Tipo:     TipoConductorTierra,
		})
	}

	return posiciones
}

// max returns the larger of two float64 values.
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two float64 values.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
