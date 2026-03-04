// internal/pdf/geometry/charola_espaciada.go
package geometry

import (
	"fmt"
	"math"
)

// ============================================================================
// CÁLCULO DE POSICIONES PARA CHAROLA ESPACIADA
// ============================================================================

// CalcularPosicionesCharolaEspaciada calcula las posiciones de los conductores
// para una charola con distribución espaciada.
//
// Reglas geométricas:
//   - La charola se dibuja en forma de U (parte superior abierta): fondo + pared izquierda + pared derecha + bridas
//   - Todos los conductores se apoyan sobre el fondo de la charola
//   - Conductores de fase: espaciados 1 diámetro entre sí (centro-a-centro = 2 * diámetro)
//   - Tierra (T): TOUCHING al último conductor (sin espacio)
//   - Cables de control: posicionados después de tierra con 1 diámetro de espaciado
//
// Los cables se CENTRAN horizontalmente dentro del ancho comercial de la charola.
//
// Args:
//   - params: Parámetros de configuración (diámetros, sistema eléctrico, etc.)
//
// Returns:
//   - Slice de ConductorPosicion con las posiciones calculadas.
//     CY es siempre 0; el componente calculador debe sumar el offset del fondo.
//
// Ejemplo de uso:
//
//	params := geometry.ParametrosCharolaBase{
//	    DiametroFaseMM:     15.0,
//	    DiametroTierraMM:   10.0,
//	    NumHilosControl:    2,
//	    SistemaElectrico:   geometry.SistemaElectricoEstrella,
//	    HilosPorFase:       1,
//	    AnchoComercialMM:   300.0,
//	}
//	posiciones := geometry.CalcularPosicionesCharolaEspaciada(params)
func CalcularPosicionesCharolaEspaciada(params ParametrosCharolaBase) []ConductorPosicion {
	// Radios de los conductores
	radioFase := params.DiametroFaseMM / 2
	radioTierra := params.DiametroTierraMM / 2
	var radioControl float64
	if params.DiametroControlMM != nil {
		radioControl = *params.DiametroControlMM / 2
	}

	// Determinar número de fases y presencia de neutro según sistema eléctrico
	numFases := params.SistemaElectrico.CantidadFases()
	tieneNeutro := params.SistemaElectrico.NecesitaNeutro()

	// Construir lista de conductores de fase y neutro (para posicionamiento basado en índice)
	type wireInfo struct {
		etiqueta string
		tipo     TipoConductor
	}
	wires := make([]wireInfo, 0, numFases*params.HilosPorFase*2)

	// Etiquetas de fase
	faseLabels := []string{"A", "B", "C"}

	// Agregar conductores de fase
	for faseIdx := 0; faseIdx < numFases; faseIdx++ {
		labelBase := faseLabels[faseIdx]
		for hiloIdx := 0; hiloIdx < params.HilosPorFase; hiloIdx++ {
			etiqueta := labelBase
			if params.HilosPorFase > 1 {
				etiqueta = labelBase + fmt.Sprintf("%d", hiloIdx+1)
			}
			wires = append(wires, wireInfo{etiqueta: etiqueta, tipo: TipoConductorFase})
		}
	}

	// Agregar conductores neutro si el sistema lo requiere
	if tieneNeutro {
		for hiloIdx := 0; hiloIdx < params.HilosPorFase; hiloIdx++ {
			etiqueta := "N"
			if params.HilosPorFase > 1 {
				etiqueta = fmt.Sprintf("N%d", hiloIdx+1)
			}
			wires = append(wires, wireInfo{etiqueta: etiqueta, tipo: TipoConductorNeutro})
		}
	}

	// Calcular ancho total ocupado por conductores fase+neutro:
	// - Desde el borde izquierdo del primer cable hasta el borde derecho del último
	// - Cada cable ocupa 1 diámetro, cada espacio ocupa 1 diámetro
	// - Total: (2 * numWires - 1) * diámetro
	numPhaseNeutroWires := len(wires)
	totalPhaseNeutroWidth := float64(2*numPhaseNeutroWires-1) * params.DiametroFaseMM

	// Ancho total: fase+neutro + diámetro de tierra (desde borde derecho del último cable)
	totalWidth := totalPhaseNeutroWidth + params.DiametroTierraMM

	// Agregar cables de control: 1 diámetro de espacio después de tierra,
	// luego cables de control con espaciado de 1 diámetro
	if params.DiametroControlMM != nil && params.NumHilosControl > 0 {
		diametroControlMM := *params.DiametroControlMM
		// Espacio después de tierra: 1 diámetro de control
		// Cables de control: (2 * numHilos - 1) * diametroControl
		totalWidth += diametroControlMM + float64(2*params.NumHilosControl-1)*diametroControlMM
	}

	// Calcular offset de centrado: el centro de los cables debe estar en anchoComercialMM / 2
	offsetX := (params.AnchoComercialMM - totalWidth) / 2

	// Posición Y: todos los cables se apoyan en el fondo (cy se calculará con el offset del fondo)
	floorY := 0.0

	posiciones := make([]ConductorPosicion, 0, len(wires)+1+params.NumHilosControl)

	// Colocar conductores de fase y neutro con espaciado 2*diámetro centro-a-centro
	for i := 0; i < len(wires); i++ {
		wire := wires[i]
		cx := offsetX + radioFase + float64(i)*2*params.DiametroFaseMM

		color := Colores["fase"]
		if wire.tipo == TipoConductorNeutro {
			color = Colores["neutro"]
		}

		posiciones = append(posiciones, ConductorPosicion{
			CX:       cx,
			CY:       floorY,
			Radio:    radioFase,
			Color:    color,
			Etiqueta: wire.etiqueta,
			Tipo:     wire.tipo,
		})
	}

	// Conductor de tierra: TOUCHING al último conductor fase/neutro
	// Centro del último cable = offsetX + radioFase + (wires.length - 1) * 2 * diametroFaseMM
	// Borde derecho del último cable = lastCableCenter + radioFase
	// Centro de tierra = lastCableRightEdge + radioTierra
	lastWireCX := offsetX + radioFase + float64(len(wires)-1)*2*params.DiametroFaseMM
	tierraCX := lastWireCX + radioFase + radioTierra

	posiciones = append(posiciones, ConductorPosicion{
		CX:       tierraCX,
		CY:       floorY,
		Radio:    radioTierra,
		Color:    Colores["tierra"],
		Etiqueta: "T",
		Tipo:     TipoConductorTierra,
	})

	// Cables de control (si existen): posicionados después de tierra con 1 diámetro de espaciado
	if params.DiametroControlMM != nil && params.NumHilosControl > 0 {
		diametroControlMM := *params.DiametroControlMM
		controlStartCX := tierraCX + radioTierra + diametroControlMM

		for i := 0; i < params.NumHilosControl; i++ {
			etiqueta := "Ctrl"
			if params.NumHilosControl > 1 {
				etiqueta = fmt.Sprintf("Ctrl%d", i+1)
			}

			posiciones = append(posiciones, ConductorPosicion{
				CX:       controlStartCX + float64(i)*2*diametroControlMM,
				CY:       floorY,
				Radio:    radioControl,
				Color:    Colores["control"],
				Etiqueta: etiqueta,
				Tipo:     TipoConductorControl,
			})
		}
	}

	return posiciones
}

// CalcularAnchoOcupadoCharola calcula el ancho total ocupado por los conductores en la charola.
func CalcularAnchoOcupadoCharola(posiciones []ConductorPosicion) float64 {
	if len(posiciones) == 0 {
		return 0
	}

	minCX := math.MaxFloat64
	maxCX := -math.MaxFloat64

	for _, pos := range posiciones {
		izquierda := pos.CX - pos.Radio
		derecha := pos.CX + pos.Radio

		if izquierda < minCX {
			minCX = izquierda
		}
		if derecha > maxCX {
			maxCX = derecha
		}
	}

	return maxCX - minCX
}
