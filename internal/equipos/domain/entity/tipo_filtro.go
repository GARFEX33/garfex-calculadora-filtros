// internal/equipos/domain/entity/tipo_filtro.go
package entity

import "fmt"

// TipoFiltro identifies the category of filter equipment.
// Maps to the PostgreSQL enum public.tipo_filtro.
// Values match exactly the DB enum: 'A', 'KVA', 'KVAR'
type TipoFiltro string

const (
	TipoFiltroA    TipoFiltro = "A"    // Filtro activo — rated in Amperes
	TipoFiltroKVA  TipoFiltro = "KVA"  // Filtro rated in KVA
	TipoFiltroKVAR TipoFiltro = "KVAR" // Filtro de rechazo — rated in KVAR (reactive)
)

// ParseTipoFiltro converts a string (e.g., from the database or HTTP request)
// to a TipoFiltro. Returns ErrTipoFiltroInvalido if the value is not recognized.
func ParseTipoFiltro(s string) (TipoFiltro, error) {
	switch s {
	case string(TipoFiltroA):
		return TipoFiltroA, nil
	case string(TipoFiltroKVA):
		return TipoFiltroKVA, nil
	case string(TipoFiltroKVAR):
		return TipoFiltroKVAR, nil
	default:
		return "", fmt.Errorf("%w: '%s' — valores válidos: A, KVA, KVAR", ErrTipoFiltroInvalido, s)
	}
}

// String returns the string representation of the TipoFiltro.
func (t TipoFiltro) String() string {
	return string(t)
}
