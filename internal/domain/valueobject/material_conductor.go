// internal/domain/valueobject/material_conductor.go
package valueobject

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
