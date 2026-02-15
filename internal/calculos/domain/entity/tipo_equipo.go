// internal/calculos/domain/entity/tipo_equipo.go
package entity

import "fmt"

// TipoEquipo identifies the category of electrical equipment.
// Using a named string type allows the domain to remain descriptive and
// extensible: filters (active, rejection), boards, transformers, loads, etc.
type TipoEquipo string

const (
	TipoEquipoFiltroActivo  TipoEquipo = "FILTRO_ACTIVO"
	TipoEquipoFiltroRechazo TipoEquipo = "FILTRO_RECHAZO"
	TipoEquipoTransformador TipoEquipo = "TRANSFORMADOR"
	TipoEquipoCarga         TipoEquipo = "CARGA"
)

// ParseTipoEquipo converts a string (e.g., from the database) to a TipoEquipo.
func ParseTipoEquipo(s string) (TipoEquipo, error) {
	switch s {
	case string(TipoEquipoFiltroActivo):
		return TipoEquipoFiltroActivo, nil
	case string(TipoEquipoFiltroRechazo):
		return TipoEquipoFiltroRechazo, nil
	case string(TipoEquipoTransformador):
		return TipoEquipoTransformador, nil
	case string(TipoEquipoCarga):
		return TipoEquipoCarga, nil
	default:
		return "", fmt.Errorf("%w: '%s'", ErrTipoEquipoInvalido, s)
	}
}

func (t TipoEquipo) String() string {
	return string(t)
}
