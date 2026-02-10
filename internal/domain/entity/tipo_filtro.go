// internal/domain/entity/tipo_filtro.go
package entity

import "fmt"

type TipoFiltro string

const (
	TipoFiltroActivo  TipoFiltro = "ACTIVO"
	TipoFiltroRechazo TipoFiltro = "RECHAZO"
)

func ParseTipoFiltro(s string) (TipoFiltro, error) {
	switch s {
	case string(TipoFiltroActivo):
		return TipoFiltroActivo, nil
	case string(TipoFiltroRechazo):
		return TipoFiltroRechazo, nil
	default:
		return "", fmt.Errorf("%w: '%s'", ErrTipoFiltroInvalido, s)
	}
}

func (t TipoFiltro) String() string {
	return string(t)
}
