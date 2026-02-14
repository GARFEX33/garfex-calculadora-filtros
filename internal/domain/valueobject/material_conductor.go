// internal/domain/valueobject/material_conductor.go
package valueobject

import (
	"fmt"
	"strings"
)

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

// MarshalJSON serializes MaterialConductor as JSON string.
func (m MaterialConductor) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", m.String())), nil
}

// UnmarshalJSON parses JSON string into MaterialConductor.
// Accepts: "Cu", "CU", "cu", "cobre", "Al", "AL", "al", "aluminio"
func (m *MaterialConductor) UnmarshalJSON(data []byte) error {
	// Remove quotes
	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// Case-insensitive comparison
	s = strings.ToUpper(s)

	switch s {
	case "CU", "COBRE", "0":
		*m = MaterialCobre
	case "AL", "ALUMINIO", "1":
		*m = MaterialAluminio
	default:
		return fmt.Errorf("invalid material: %s (expected Cu or Al)", s)
	}
	return nil
}
