// internal/pdf/geometry/geometry.go
package geometry

import "fmt"

// ============================================================================
// CONSTANTES
// ============================================================================

// PeralteCharolaMM es la altura estándar de pared de charola: 70mm (tamaño comercial común).
const PeralteCharolaMM = 70

// EspesorFondoCharolaMM es el espesor del fondo de la charola para renderizado SVG (visual).
const EspesorFondoCharolaMM = 3

// EspesorParedCharolaMM es el espesor de la pared de la charola para renderizado SVG (visual).
const EspesorParedCharolaMM = 3

// AnchoBridaCharolaMM es el ancho de la brida de la charola (cosmético, se extiende hacia afuera).
const AnchoBridaCharolaMM = 15

// Colores es la paleta de colores estilo CAD para diagramas SVG.
var Colores = map[string]string{
	"fase":          "#D4AF37",
	"tierra":        "#28A745",
	"neutro":        "#6B7280",
	"control":       "#3B82F6",
	"charolaStroke": "#004085",
	"charolaFill":   "url(#charola-hatch)",
	"tuboStroke":    "#555555",
	"tuboFill":      "#E8E8E8",
	"cotaLinea":     "#374151",
	"cotaTexto":     "#6B7280",
	"titleBlock":    "#1F2937",
	"fondo":         "#FFFFFF",
}

// ============================================================================
// TIPOS
// ============================================================================

// SistemaElectrico representa el tipo de sistema eléctrico.
// Valores: DELTA, ESTRELLA, BIFASICO, MONOFASICO.
type SistemaElectrico string

// Valores válidos de SistemaElectrico.
const (
	SistemaElectricoDelta      SistemaElectrico = "DELTA"
	SistemaElectricoEstrella   SistemaElectrico = "ESTRELLA"
	SistemaElectricoBifasico   SistemaElectrico = "BIFASICO"
	SistemaElectricoMonofasico SistemaElectrico = "MONOFASICO"
)

// TipoConductor representa el tipo de conductor en el diagrama.
type TipoConductor string

// Valores válidos de TipoConductor.
const (
	TipoConductorFase    TipoConductor = "fase"
	TipoConductorTierra  TipoConductor = "tierra"
	TipoConductorNeutro  TipoConductor = "neutro"
	TipoConductorControl TipoConductor = "control"
)

// ErrSistemaElectricoInvalido es el error returned when SistemaElectrico es inválido.
var ErrSistemaElectricoInvalido = fmt.Errorf("sistema eléctrico no válido")

// ParseSistemaElectrico convierte un string a SistemaElectrico.
func ParseSistemaElectrico(s string) (SistemaElectrico, error) {
	switch s {
	case string(SistemaElectricoDelta):
		return SistemaElectricoDelta, nil
	case string(SistemaElectricoEstrella):
		return SistemaElectricoEstrella, nil
	case string(SistemaElectricoBifasico):
		return SistemaElectricoBifasico, nil
	case string(SistemaElectricoMonofasico):
		return SistemaElectricoMonofasico, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrSistemaElectricoInvalido, s)
	}
}

// NecesitaNeutro devuelve true si el sistema eléctrico requiere conductor neutro.
func (s SistemaElectrico) NecesitaNeutro() bool {
	switch s {
	case SistemaElectricoDelta:
		return false
	case SistemaElectricoEstrella, SistemaElectricoBifasico, SistemaElectricoMonofasico:
		return true
	default:
		return false
	}
}

// CantidadFases devuelve la cantidad de fases del sistema.
func (s SistemaElectrico) CantidadFases() int {
	switch s {
	case SistemaElectricoDelta, SistemaElectricoEstrella:
		return 3
	case SistemaElectricoBifasico:
		return 2
	case SistemaElectricoMonofasico:
		return 1
	default:
		return 3
	}
}

// ConductorPosicion representa la posición de un conductor en el diagrama SVG.
type ConductorPosicion struct {
	// CX es la posición X del centro del conductor (mm).
	CX float64 `json:"cx"`
	// CY es la posición Y del centro del conductor (mm).
	CY float64 `json:"cy"`
	// Radio es el radio del conductor (mm).
	Radio float64 `json:"radio"`
	// Color es el color de relleno en formato hex.
	Color string `json:"color"`
	// Etiqueta es la etiqueta del conductor: "A", "B", "C", "N", "T", "Ctrl1".
	Etiqueta string `json:"etiqueta"`
	// Tipo es el tipo de conductor (fase, neutro, tierra, control).
	Tipo TipoConductor `json:"tipo"`
}

// ViewBoxResult contiene el resultado del cálculo del viewBox SVG.
type ViewBoxResult struct {
	// ViewBox es el string del viewBox en formato "minX minY width height".
	ViewBox string `json:"viewBox"`
	// Ancho es el ancho total del contenido (mm).
	Ancho float64 `json:"ancho"`
	// Alto es el alto total del contenido (mm).
	Alto float64 `json:"alto"`
}

// LineaCota representa una línea de dimensión para el diagrama SVG.
type LineaCota struct {
	// X1 es la posición X inicial.
	X1 float64 `json:"x1"`
	// Y1 es la posición Y inicial.
	Y1 float64 `json:"y1"`
	// X2 es la posición X final.
	X2 float64 `json:"x2"`
	// Y2 es la posición Y final.
	Y2 float64 `json:"y2"`
	// Valor es el valor de la dimensión en mm.
	Valor float64 `json:"valor"`
	// Texto es el texto a mostrar (ej: "152.4 mm (6\")").
	Texto string `json:"texto"`
	// PosicionTexto indica si el texto se muestra arriba o abajo de la línea.
	PosicionTexto string `json:"posicionTexto"` // "arriba" | "abajo"
}

// ParametrosCharolaBase contiene los parámetros base para cálculos de charola.
type ParametrosCharolaBase struct {
	// Diámetro del conductor de fase en mm.
	DiametroFaseMM float64
	// Diámetro del conductor de tierra en mm.
	DiametroTierraMM float64
	// Diámetro opcional del conductor de control en mm.
	DiametroControlMM *float64
	// Número total de hilos de control.
	NumHilosControl int
	// Sistema eléctrico (DELTA, ESTRELLA, BIFASICO, MONOFASICO).
	SistemaElectrico SistemaElectrico
	// Número de hilos por fase (conductores en paralelo).
	HilosPorFase int
	// Ancho comercial de la charola en mm (para centrado).
	AnchoComercialMM float64
}
