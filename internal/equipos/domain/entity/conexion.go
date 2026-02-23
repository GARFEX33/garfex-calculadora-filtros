// internal/equipos/domain/entity/conexion.go
package entity

import "fmt"

// Conexion identifies the electrical connection type of a filter equipment.
// Maps to the PostgreSQL enum public.conexion.
// Values match exactly the DB enum.
type Conexion string

const (
	ConexionDelta      Conexion = "DELTA"      // Three-phase delta (∆) connection
	ConexionEstrella   Conexion = "ESTRELLA"   // Three-phase star/wye (Y) connection
	ConexionMonofasico Conexion = "MONOFASICO" // Single-phase connection
	ConexionBifasico   Conexion = "BIFASICO"   // Two-phase connection
)

// ParseConexion converts a string (e.g., from the database or HTTP request)
// to a Conexion. Returns ErrConexionInvalida if the value is not recognized.
func ParseConexion(s string) (Conexion, error) {
	switch s {
	case string(ConexionDelta):
		return ConexionDelta, nil
	case string(ConexionEstrella):
		return ConexionEstrella, nil
	case string(ConexionMonofasico):
		return ConexionMonofasico, nil
	case string(ConexionBifasico):
		return ConexionBifasico, nil
	default:
		return "", fmt.Errorf("%w: '%s' — valores válidos: DELTA, ESTRELLA, MONOFASICO, BIFASICO", ErrConexionInvalida, s)
	}
}

// String returns the string representation of the Conexion.
func (c Conexion) String() string {
	return string(c)
}
