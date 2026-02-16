// internal/shared/kernel/valueobject/material_conductor.go
package valueobject

import (
	"fmt"
	"strings"
)

// ErrMaterialConductorInvalido is returned when an unknown material string is provided.
var ErrMaterialConductorInvalido = fmt.Errorf("material conductor inválido (esperado Cu o Al)")

// MaterialConductor represents the conductor material (Cu or Al).
type MaterialConductor int

const (
	MaterialCobre MaterialConductor = iota
	MaterialAluminio
)

// String returns the NOM standard abbreviation for the material.
func (m MaterialConductor) String() string {
	switch m {
	case MaterialCobre:
		return "CU"
	case MaterialAluminio:
		return "AL"
	default:
		return "UNKNOWN"
	}
}

// ParseMaterialConductor converts a string to MaterialConductor.
// Accepts case-insensitive: "Cu", "CU", "cu", "cobre", "Al", "AL", "al", "aluminio".
// Returns ErrMaterialConductorInvalido if the value is unrecognized.
func ParseMaterialConductor(s string) (MaterialConductor, error) {
	switch strings.ToUpper(s) {
	case "CU", "COBRE":
		return MaterialCobre, nil
	case "AL", "ALUMINIO":
		return MaterialAluminio, nil
	default:
		return MaterialCobre, fmt.Errorf("%w: %q", ErrMaterialConductorInvalido, s)
	}
}
