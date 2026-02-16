// internal/calculos/domain/entity/tipo_canalizacion.go
package entity

import (
	"errors"
	"fmt"
)

// TipoCanalizacion identifies the wiring method and conduit material used for
// conductor installation. It determines:
//   - Which NOM ampacity table to use for conductor selection
//   - Which resistance column to use from Tabla 9 for voltage drop calculation
//   - The DMG factor for inductive reactance calculation
type TipoCanalizacion string

const (
	// TipoCanalizacionTuberiaPVC represents conductors installed in PVC conduit.
	// Ampacity table: 310-15(b)(16) → data/tablas_nom/310-15-b-16.csv
	// Tabla 9 resistance column: res_{material}_pvc
	// DMG factor: 1.0
	TipoCanalizacionTuberiaPVC TipoCanalizacion = "TUBERIA_PVC"

	// TipoCanalizacionTuberiaAluminio represents conductors installed in aluminum conduit.
	// Ampacity table: 310-15(b)(16) → data/tablas_nom/310-15-b-16.csv
	// Tabla 9 resistance column: res_{material}_al
	// DMG factor: 1.0
	TipoCanalizacionTuberiaAluminio TipoCanalizacion = "TUBERIA_ALUMINIO"

	// TipoCanalizacionTuberiaAceroPG represents conductors installed in rigid steel
	// conduit (pared gruesa / thick wall).
	// Ampacity table: 310-15(b)(16) → data/tablas_nom/310-15-b-16.csv
	// Tabla 9 resistance column: res_{material}_acero
	// DMG factor: 1.0
	TipoCanalizacionTuberiaAceroPG TipoCanalizacion = "TUBERIA_ACERO_PG"

	// TipoCanalizacionTuberiaAceroPD represents conductors installed in electrical
	// metallic tubing / thin wall steel conduit (pared delgada).
	// Ampacity table: 310-15(b)(16) → data/tablas_nom/310-15-b-16.csv
	// Tabla 9 resistance column: res_{material}_acero
	// DMG factor: 1.0
	TipoCanalizacionTuberiaAceroPD TipoCanalizacion = "TUBERIA_ACERO_PD"

	// TipoCanalizacionCharolaCableEspaciado represents conductors on a cable tray
	// with spacing between cables (one diameter apart).
	// Ampacity table: 310-15(b)(17) → data/tablas_nom/310-15-b-17.csv
	// Tabla 9 resistance column: res_{material}_pvc (no metallic conduit effect)
	// DMG factor: 2.0
	TipoCanalizacionCharolaCableEspaciado TipoCanalizacion = "CHAROLA_CABLE_ESPACIADO"

	// TipoCanalizacionCharolaCableTriangular represents conductors on a cable tray
	// in triangular arrangement (touching). No 60°C column — minimum column is 75°C.
	// Ampacity table: 310-15(b)(20) → data/tablas_nom/310-15-b-20.csv
	// Tabla 9 resistance column: res_{material}_pvc (no metallic conduit effect)
	// DMG factor: 1.0
	TipoCanalizacionCharolaCableTriangular TipoCanalizacion = "CHAROLA_CABLE_TRIANGULAR"
)

// ErrTipoCanalizacionInvalido is returned when an unknown TipoCanalizacion is provided.
var ErrTipoCanalizacionInvalido = errors.New("tipo de canalización inválido")

// ValidarTipoCanalizacion returns an error if tc is not a recognized canalization type.
func ValidarTipoCanalizacion(tc TipoCanalizacion) error {
	switch tc {
	case TipoCanalizacionTuberiaPVC,
		TipoCanalizacionTuberiaAluminio,
		TipoCanalizacionTuberiaAceroPG,
		TipoCanalizacionTuberiaAceroPD,
		TipoCanalizacionCharolaCableEspaciado,
		TipoCanalizacionCharolaCableTriangular:
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrTipoCanalizacionInvalido, tc)
	}
}

// ParseTipoCanalizacion converts a string to a TipoCanalizacion.
func ParseTipoCanalizacion(s string) (TipoCanalizacion, error) {
	switch s {
	case string(TipoCanalizacionTuberiaPVC):
		return TipoCanalizacionTuberiaPVC, nil
	case string(TipoCanalizacionTuberiaAluminio):
		return TipoCanalizacionTuberiaAluminio, nil
	case string(TipoCanalizacionTuberiaAceroPG):
		return TipoCanalizacionTuberiaAceroPG, nil
	case string(TipoCanalizacionTuberiaAceroPD):
		return TipoCanalizacionTuberiaAceroPD, nil
	case string(TipoCanalizacionCharolaCableEspaciado):
		return TipoCanalizacionCharolaCableEspaciado, nil
	case string(TipoCanalizacionCharolaCableTriangular):
		return TipoCanalizacionCharolaCableTriangular, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrTipoCanalizacionInvalido, s)
	}
}
