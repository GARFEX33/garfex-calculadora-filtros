package entity

import (
	"errors"
	"fmt"
)

// TipoCanalizacion identifies the wiring method used for conductor installation.
// It determines which NOM ampacity table is used for conductor selection.
type TipoCanalizacion string

const (
	// TipoCanalizacionTuberiaConduit represents conductors installed in conduit.
	// NOM reference table: 310-15(b)(16) → data/tablas_nom/310-15-b-16.csv
	TipoCanalizacionTuberiaConduit TipoCanalizacion = "TUBERIA_CONDUIT"

	// TipoCanalizacionCharolaCableEspaciado represents conductors on a cable tray
	// with spacing between cables.
	// NOM reference table: 310-15(b)(17) → data/tablas_nom/310-15-b-17.csv
	TipoCanalizacionCharolaCableEspaciado TipoCanalizacion = "CHAROLA_CABLE_ESPACIADO"

	// TipoCanalizacionCharolaCableTriangular represents conductors on a cable tray
	// in triangular arrangement (touching). No 60°C column — minimum column is 75°C.
	// NOM reference table: 310-15(b)(20) → data/tablas_nom/310-15-b-20.csv
	TipoCanalizacionCharolaCableTriangular TipoCanalizacion = "CHAROLA_CABLE_TRIANGULAR"
)

// ErrTipoCanalizacionInvalido is returned when an unknown TipoCanalizacion is provided.
var ErrTipoCanalizacionInvalido = errors.New("tipo de canalización inválido")

// ValidarTipoCanalizacion returns an error if tc is not a recognized canalization type.
func ValidarTipoCanalizacion(tc TipoCanalizacion) error {
	switch tc {
	case TipoCanalizacionTuberiaConduit,
		TipoCanalizacionCharolaCableEspaciado,
		TipoCanalizacionCharolaCableTriangular:
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrTipoCanalizacionInvalido, tc)
	}
}
